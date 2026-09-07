package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/assets_managements/database"
	"github.com/lib/pq"
)

// Issuer service: the legal vehicle a token represents shares of.
//
// Since §11.18 an asset cannot be published without an issuer_id, but nothing
// could create an issuer other than a direct SQL INSERT. The rule was right and
// the tooling to satisfy it was missing; this file provides it.
//
// The three sentinels mirror the real estate ones for the same reason: the
// handler picks a status code without knowing a single business rule.
var (
	// ErrInvalidIssuer: the payload cannot produce a valid issuer. 400.
	ErrInvalidIssuer = errors.New("invalid issuer payload")
	// ErrIssuerNotFound: no such issuer. 404.
	ErrIssuerNotFound = errors.New("issuer not found")
	// ErrIssuerConflict: the payload is valid, but the current state forbids
	// the change. 409.
	ErrIssuerConflict = errors.New("issuer state forbids this change")
)

// leiCodeRe mirrors constraint issuer_lei_ck. Validating here turns a raw
// constraint name into a message naming the field.
var leiCodeRe = regexp.MustCompile(`^[A-Z0-9]{20}$`)

const (
	maxIssuerNameLength         = 255
	maxIssuerLegalFormLength    = 100
	maxIssuerRegistrationLength = 64
)

// issuerStatusLabels is the set seeded by migration 000007. A status outside it
// is refused by name rather than by the foreign key, whose message means
// nothing to a client.
var issuerStatusLabels = map[int]string{
	database.IssuerStatusDraft:     "draft",
	database.IssuerStatusActive:    "active",
	database.IssuerStatusSuspended: "suspended",
	database.IssuerStatusDissolved: "dissolved",
}

// activeRequirements is the single definition of what an issuer must carry to
// be active — the status at which assets can be attached and published under
// it.
//
// It is consulted from the two places that need it: the creation, which refuses
// to produce an active vehicle without it, and the patch, which refuses to
// activate a draft that still lacks it. One list, like publicationRequirements
// on the asset side, because a second one would drift.
//
// registration_number, and only it:
//
//   - a vehicle assets are published under must be identifiable in its register
//     (the RCS number in Luxembourg). It is an investor's only way to check
//     that the issuer exists at all, and a manager has no reason to re-verify a
//     name that looks plausible;
//   - a LEI is deliberately absent. It is mandatory only for entities under
//     market reporting obligations or whose securities are admitted to trading,
//     so demanding it would block vehicles placed privately — and it lapses
//     every year, which would turn a renewal delay into an unusable issuer.
func activeRequirements(in database.IssuerWriteDTO) []string {
	var missing []string

	if in.RegistrationNumber == nil {
		missing = append(missing, "registration_number")
	}

	return missing
}

func invalidIssuer(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrInvalidIssuer, fmt.Sprintf(format, args...))
}

func issuerConflict(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrIssuerConflict, fmt.Sprintf(format, args...))
}

// GetIssuers lists the issuers for the back office.
func (s *Service) GetIssuers(ctx context.Context) ([]server.Issuer, error) {
	dtos, err := database.GetIssuers(ctx)
	if err != nil {
		return nil, err
	}

	issuers := make([]server.Issuer, 0, len(dtos))
	for _, dto := range dtos {
		issuers = append(issuers, toServerIssuer(dto))
	}

	return issuers, nil
}

func (s *Service) GetIssuerById(ctx context.Context, id int) (server.Issuer, error) {
	dto, err := database.GetIssuerById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return server.Issuer{}, ErrIssuerNotFound
		}
		return server.Issuer{}, err
	}

	return toServerIssuer(dto), nil
}

// CreateIssuer validates the payload, writes the issuer and returns it as it
// now reads.
//
// The stored row is read back rather than echoed: the table fills the
// identifier, the timestamps and the status label, and the caller is better
// served by what was actually stored than by what he sent.
func (s *Service) CreateIssuer(ctx context.Context, req server.IssuerWriteRequest) (server.Issuer, error) {
	in, err := toIssuerWriteDTO(req)
	if err != nil {
		return server.Issuer{}, err
	}

	if err := checkIssuerCreatable(in); err != nil {
		return server.Issuer{}, err
	}

	id, err := database.CreateIssuer(ctx, in)
	if err != nil {
		return server.Issuer{}, translateIssuerWriteError(err)
	}

	return s.GetIssuerById(ctx, id)
}

// PatchIssuer applies a partial update.
//
// As for an asset: the stored row is read, the patch is laid over it, and the
// RESULT is validated as a whole. Validating only the fields that changed would
// let a patch produce a state a creation would have refused.
func (s *Service) PatchIssuer(ctx context.Context, id int, patch server.IssuerPatchRequest) (server.Issuer, error) {
	state, err := database.GetIssuerGuardState(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return server.Issuer{}, ErrIssuerNotFound
		}
		return server.Issuer{}, err
	}

	current, err := database.GetIssuerById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return server.Issuer{}, ErrIssuerNotFound
		}
		return server.Issuer{}, err
	}

	merged := applyIssuerPatch(toIssuerWriteRequest(current), patch)

	in, err := toIssuerWriteDTO(merged)
	if err != nil {
		return server.Issuer{}, err
	}

	if err := checkIssuerStatusChange(state, in); err != nil {
		return server.Issuer{}, err
	}

	if err := checkIssuerStillIdentifiable(storedIssuerWriteDTO(current), in); err != nil {
		return server.Issuer{}, err
	}

	if err := database.UpdateIssuer(ctx, id, in); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return server.Issuer{}, ErrIssuerNotFound
		}
		return server.Issuer{}, translateIssuerWriteError(err)
	}

	return s.GetIssuerById(ctx, id)
}

// DeleteIssuer dissolves the vehicle instead of removing it.
//
// An issuer that still carries assets is refused: the foreign key is declared
// ON DELETE RESTRICT precisely because those assets describe shares of it, and
// leaving them pointing at a dissolved vehicle would publish an offer nothing
// stands behind.
//
// Dissolving an already dissolved issuer is not an error: it changes nothing,
// so two managers clicking the same button do not produce a failure.
func (s *Service) DeleteIssuer(ctx context.Context, id int) error {
	state, err := database.GetIssuerGuardState(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrIssuerNotFound
		}
		return err
	}

	if state.StatusId == database.IssuerStatusDissolved {
		return nil
	}

	if state.AttachedCount > 0 {
		return issuerConflict(
			"this issuer still carries %d asset(s); detach or delete them first",
			state.AttachedCount)
	}

	if err := database.DissolveIssuer(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrIssuerNotFound
		}
		return err
	}

	return nil
}

// checkIssuerCreatable refuses the states a creation must not produce.
//
// toIssuerWriteDTO has already resolved an absent status to active: a manager
// filling this form declares a vehicle that exists, and parking it in a draft
// nobody can attach an asset to would be a step nobody asked for. That default
// is precisely why the identification is demanded here by default too.
//
// Both refusals are a 400 and not a 409: nothing about a stored state forbids
// them, the payload is simply incomplete for what it asks to create.
func checkIssuerCreatable(in database.IssuerWriteDTO) error {
	// Dissolved is final: creating a row already in that state would produce an
	// issuer that can never be used and never be reopened.
	if in.StatusId == database.IssuerStatusDissolved {
		return invalidIssuer("an issuer cannot be created dissolved")
	}

	if in.StatusId == database.IssuerStatusActive {
		if missing := activeRequirements(in); len(missing) > 0 {
			return invalidIssuer(
				"an active issuer requires %s; create it as a draft (status_id 1) while it is still being incorporated",
				strings.Join(missing, ", "))
		}
	}

	return nil
}

// checkIssuerStatusChange guards the transitions that would contradict what
// investors are being shown.
//
// Three rules, deliberately separate:
//
//   - dissolved is final (is_final in ass.issuer_status). Reopening it would
//     make the flag a decoration.
//   - withdrawing a vehicle from service while it carries published assets
//     leaves an offer standing on nothing. Unpublishing those assets
//     automatically would be worse: a change to one row silently taking others
//     off the market.
//   - activating a draft that is not identifiable yet. This one is a 409 and
//     not a 400 for the same reason publishing an incomplete asset is: what is
//     refused is the transition against the record as it stands, not the shape
//     of the payload. Supplying the number in the same request lifts it.
//
// Everything here is keyed on a CHANGE of status. A patch that does not touch
// the status never triggers it, so a legacy row stored without a registration
// number stays editable — the freeze trap already met on the address, on
// estate_type and on the payment rhythm (§11.20).
func checkIssuerStatusChange(state database.IssuerGuardStateDTO, in database.IssuerWriteDTO) error {
	target := in.StatusId

	if state.StatusId == target {
		return nil
	}

	if state.StatusId == database.IssuerStatusDissolved {
		return issuerConflict("a dissolved issuer cannot be reactivated")
	}

	if target == database.IssuerStatusActive {
		if missing := activeRequirements(in); len(missing) > 0 {
			return issuerConflict(
				"this issuer cannot be activated while it lacks %s",
				strings.Join(missing, ", "))
		}
	}

	// "Withdrawn" is every status that is not active, not the two named ones.
	// draft is just as unable to carry an offer as suspended, and naming the
	// statuses instead of the condition left active -> draft unguarded: an
	// issuer carrying published assets could be parked in a draft, producing
	// the exact state this whole guard exists to prevent. The asset side has
	// always tested the condition (issuerStandsBehind compares to active); this
	// side now does too.
	withdrawn := target != database.IssuerStatusActive
	if withdrawn && state.PublishedCount > 0 {
		return issuerConflict(
			"this issuer carries %d published asset(s); unpublish them before setting it to %s",
			state.PublishedCount, issuerStatusLabels[target])
	}

	if target == database.IssuerStatusDissolved && state.AttachedCount > 0 {
		return issuerConflict(
			"this issuer still carries %d asset(s); detach or delete them first",
			state.AttachedCount)
	}

	return nil
}

// checkIssuerStillIdentifiable refuses a patch that would strip an active
// issuer of what makes it identifiable.
//
// checkIssuerStatusChange returns early when the status does not change, which
// is what keeps legacy rows editable — and which also left the invariant of
// §11.22 unenforced for every patch that keeps the issuer active. Clearing
// registration_number was then accepted on a live vehicle, producing exactly
// the state a creation and an activation both refuse.
//
// The way out is the one the asset side already uses: compare BEFORE and AFTER
// and refuse only a LOSS. A row stored without a number stays editable — it
// loses nothing — while a row that has one cannot be emptied. The freeze trap
// is avoided without trading the invariant away for it.
//
// The stored state is mapped field by field rather than run through
// toIssuerWriteDTO: a legacy row that fails today's validation must still be
// correctable, and validating it here would refuse the very patch that fixes
// it.
func checkIssuerStillIdentifiable(before, after database.IssuerWriteDTO) error {
	if after.StatusId != database.IssuerStatusActive {
		return nil
	}

	alreadyMissing := make(map[string]bool)
	for _, field := range activeRequirements(before) {
		alreadyMissing[field] = true
	}

	var lost []string
	for _, field := range activeRequirements(after) {
		if !alreadyMissing[field] {
			lost = append(lost, field)
		}
	}

	if len(lost) > 0 {
		return issuerConflict(
			"an active issuer cannot lose what identifies it, missing: %s",
			strings.Join(lost, ", "))
	}

	return nil
}

// storedIssuerWriteDTO renders a stored row as a write DTO without validating
// it, so it can serve as the "before" of a comparison.
func storedIssuerWriteDTO(dto database.IssuerDTO) database.IssuerWriteDTO {
	return database.IssuerWriteDTO{
		Name:               dto.Name,
		LegalForm:          dto.LegalForm,
		RegistrationNumber: dto.RegistrationNumber,
		CountryCode:        dto.CountryCode,
		LeiCode:            dto.LeiCode,
		StatusId:           dto.StatusId,
	}
}

// toIssuerWriteRequest renders a stored issuer in the shape of a write request,
// so a patch merges onto it with the creation rules applying unchanged.
func toIssuerWriteRequest(dto database.IssuerDTO) server.IssuerWriteRequest {
	statusID := dto.StatusId

	return server.IssuerWriteRequest{
		Name:               dto.Name,
		LegalForm:          dto.LegalForm,
		RegistrationNumber: dto.RegistrationNumber,
		CountryCode:        dto.CountryCode,
		LeiCode:            dto.LeiCode,
		StatusId:           &statusID,
	}
}

// applyIssuerPatch lays the supplied fields over the current state.
//
// An absent field keeps its value — that is the whole point of a patch — while
// an optional field sent empty is cleared: without that, a LEI typed by mistake
// could never be removed. name, legal_form and country_code have no empty form:
// they are required, so an empty one is a mistake the validation names.
func applyIssuerPatch(current server.IssuerWriteRequest, patch server.IssuerPatchRequest) server.IssuerWriteRequest {
	merged := current

	if patch.Name != nil {
		merged.Name = *patch.Name
	}
	if patch.LegalForm != nil {
		merged.LegalForm = *patch.LegalForm
	}
	if patch.CountryCode != nil {
		merged.CountryCode = *patch.CountryCode
	}
	if patch.RegistrationNumber != nil {
		merged.RegistrationNumber = patch.RegistrationNumber
	}
	if patch.LeiCode != nil {
		merged.LeiCode = patch.LeiCode
	}
	if patch.StatusId != nil {
		merged.StatusId = patch.StatusId
	}

	return merged
}

// toIssuerWriteDTO validates and normalises a write request.
//
// The same function serves a creation and the merged result of a patch: there
// is one definition of a valid issuer, so a patch cannot slip the row into a
// state a creation would have refused.
func toIssuerWriteDTO(req server.IssuerWriteRequest) (database.IssuerWriteDTO, error) {
	in := database.IssuerWriteDTO{
		Name:               strings.TrimSpace(req.Name),
		LegalForm:          strings.TrimSpace(req.LegalForm),
		RegistrationNumber: trimmedPtr(req.RegistrationNumber),
		CountryCode:        strings.ToUpper(strings.TrimSpace(req.CountryCode)),
		StatusId:           database.IssuerStatusActive,
	}

	if in.Name == "" {
		return in, invalidIssuer("name is required")
	}
	if in.LegalForm == "" {
		return in, invalidIssuer("legal_form is required")
	}

	// Counted in runes, like the asset lengths: a name written in accented
	// French is not longer than the same name in English.
	for _, field := range []struct {
		label string
		value string
		max   int
	}{
		{"name", in.Name, maxIssuerNameLength},
		{"legal_form", in.LegalForm, maxIssuerLegalFormLength},
	} {
		if n := len([]rune(field.value)); n > field.max {
			return in, invalidIssuer("%s must be at most %d characters, got %d", field.label, field.max, n)
		}
	}

	if in.RegistrationNumber != nil {
		if n := len([]rune(*in.RegistrationNumber)); n > maxIssuerRegistrationLength {
			return in, invalidIssuer(
				"registration_number must be at most %d characters, got %d",
				maxIssuerRegistrationLength, n)
		}
	}

	if !countryCodeRe.MatchString(in.CountryCode) {
		return in, invalidIssuer("country_code must be an ISO 3166-1 alpha-2 code, got %q", req.CountryCode)
	}

	// CHAR(20) pads whatever it is given to twenty characters, so a LEI that is
	// merely too short would be stored padded and pass its own CHECK on the way
	// back. It is validated here, on the exact value.
	if lei := trimmedPtr(req.LeiCode); lei != nil {
		upper := strings.ToUpper(*lei)
		if !leiCodeRe.MatchString(upper) {
			return in, invalidIssuer("lei_code must be 20 upper-case alphanumerics, got %q", *req.LeiCode)
		}
		in.LeiCode = &upper
	}

	if req.StatusId != nil {
		if _, known := issuerStatusLabels[*req.StatusId]; !known {
			return in, invalidIssuer("status_id %d does not refer to an existing status", *req.StatusId)
		}
		in.StatusId = *req.StatusId
	}

	return in, nil
}

func toServerIssuer(dto database.IssuerDTO) server.Issuer {
	count := dto.RealEstateCount
	createdAt := dto.CreatedAt
	updatedAt := dto.UpdatedAt

	return server.Issuer{
		Id:                 dto.Id,
		Name:               dto.Name,
		LegalForm:          dto.LegalForm,
		RegistrationNumber: dto.RegistrationNumber,
		CountryCode:        dto.CountryCode,
		LeiCode:            dto.LeiCode,
		StatusId:           dto.StatusId,
		Status:             dto.StatusCode,
		RealEstateCount:    &count,
		CreatedAt:          &createdAt,
		UpdatedAt:          &updatedAt,
	}
}

// translateIssuerWriteError turns the constraint violations a payload can cause
// into the sentinel the handler understands.
//
// Duplicating a registration number or a LEI is a state conflict, not a typo:
// the value is well formed, it simply already designates another vehicle.
func translateIssuerWriteError(err error) error {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return err
	}

	switch pqErr.Code.Name() {
	case "unique_violation":
		// Same reasoning as on the asset side: a primary key collision is the
		// sequence lagging behind the data, never something the caller sent.
		if isPrimaryKey(pqErr.Constraint) {
			return fmt.Errorf(
				"the identity sequence of the table is behind its data (%s, %s): replay migration 000013 to resynchronise it",
				pqErr.Constraint, pqErr.Detail)
		}
		switch pqErr.Constraint {
		case "issuer_registration_uk":
			return issuerConflict("this registration number is already used by another issuer in this country")
		case "issuer_lei_uk":
			return issuerConflict("this LEI code already identifies another issuer")
		}
		return issuerConflict("this value is already used (%s)", pqErr.Constraint)
	case "foreign_key_violation":
		if pqErr.Constraint == "issuer_status_id_fkey" {
			return invalidIssuer("status_id does not refer to an existing status")
		}
		return invalidIssuer("a referenced record does not exist (%s)", pqErr.Constraint)
	case "check_violation":
		switch pqErr.Constraint {
		case "issuer_country_code_ck":
			return invalidIssuer("country_code must be an ISO 3166-1 alpha-2 code")
		case "issuer_lei_ck":
			return invalidIssuer("lei_code must be 20 upper-case alphanumerics")
		}
		return invalidIssuer("a value is out of the range allowed by the database (%s)", pqErr.Constraint)
	case "string_data_right_truncation":
		return invalidIssuer("a value is longer than the database allows (%s)", pqErr.Column)
	case "not_null_violation":
		return invalidIssuer("%s cannot be empty in this database", pqErr.Column)
	}

	return err
}
