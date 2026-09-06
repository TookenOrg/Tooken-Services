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
	"github.com/shopspring/decimal"
)

// Write side of the real estate service: validation, business guards and
// orchestration of the transactional writes.
//
// The three sentinel errors below are what lets the handler pick a status code
// without knowing anything about the rules. Everything else is a 500: an
// unexpected failure must never be reported as the caller's fault.
var (
	// ErrInvalidRealEstate: the payload cannot produce a valid asset. 400.
	ErrInvalidRealEstate = errors.New("invalid real estate payload")
	// ErrRealEstateNotFound: no such asset, or it was deleted. 404.
	ErrRealEstateNotFound = errors.New("real estate not found")
	// ErrRealEstateConflict: the payload is valid, but the current state of the
	// asset forbids the change. 409.
	ErrRealEstateConflict = errors.New("real estate state forbids this change")
)

var (
	countryCodeRe  = regexp.MustCompile(`^[A-Z]{2}$`)
	currencyCodeRe = regexp.MustCompile(`^[A-Z]{3}$`)
	energyClassRe  = regexp.MustCompile(`^[A-G]$`)
)

// allowedMediaTypes mirrors the CHECK constraint of migration 000005. Rejecting
// here turns a raw constraint violation into a message the caller can act on.
var allowedMediaTypes = map[string]bool{
	"image":        true,
	"video":        true,
	"floorplan":    true,
	"virtual_tour": true,
}

// writeMode tells a creation from a patch, which is what decides whether the
// fields that describe the asset itself are mandatory.
//
// A creation must carry an address and an estate type: an asset nobody can
// locate, or that belongs to no category, is not a registry entry. A patch must
// not require them, or a row inherited without one could never be corrected —
// not even to fix its title.
//
// The distinction is only about what is *demanded*. What is *validated* is
// always the whole merged result, so a patch can never leave the asset in a
// shape a creation would have refused.
type writeMode bool

const (
	writeCreate writeMode = true
	writePatch  writeMode = false
)

// isBlankAddress reports an address the caller never filled, as opposed to one
// filled with wrong values — the latter must still be rejected.
func isBlankAddress(a server.RealEstateAddress) bool {
	for _, value := range []string{a.Street, a.PostalCode, a.City, a.CountryCode} {
		if strings.TrimSpace(value) != "" {
			return false
		}
	}

	for _, value := range []*string{a.StreetComplement, a.Region, a.Latitude, a.Longitude} {
		if value != nil && strings.TrimSpace(*value) != "" {
			return false
		}
	}

	return true
}

func invalid(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrInvalidRealEstate, fmt.Sprintf(format, args...))
}

func conflict(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrRealEstateConflict, fmt.Sprintf(format, args...))
}

// CreateRealEstate validates the payload, writes the asset in one transaction
// and returns it as it now reads.
//
// The created asset is read back rather than echoed: the database fills the
// identifiers, the timestamps, the derived valuation and the status, and the
// caller is better served by what was actually stored than by what was sent.
func (s *Service) CreateRealEstate(ctx context.Context, req server.RealEstateWriteRequest) (server.RealEstate, error) {
	in, err := toWriteDTO(req, writeCreate)
	if err != nil {
		return server.RealEstate{}, err
	}

	// A complete asset goes online immediately, an incomplete one waits as a
	// draft. See publicationRequirements for what "complete" means and why the
	// definition lives in one place.
	id, err := database.CreateRealEstate(ctx, in, initialStatus(in))
	if err != nil {
		return server.RealEstate{}, translateWriteError(err)
	}

	return s.GetRealEstateById(ctx, id, true)
}

// PatchRealEstate applies a partial update.
//
// The current asset is read, the patch is laid over it, and the RESULT is
// validated as a whole. Validating only the fields that changed would let a
// patch produce a state a create would have refused — lowering a maximum below
// an untouched minimum, for instance.
//
// The merged asset is then written by the same code path as a create, so there
// is a single definition of what a stored asset looks like.
func (s *Service) PatchRealEstate(ctx context.Context, id int, patch server.RealEstatePatchRequest) (server.RealEstate, error) {
	current, err := database.GetRealEstateById(ctx, id, true)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return server.RealEstate{}, ErrRealEstateNotFound
		}
		return server.RealEstate{}, err
	}

	stored := toWriteRequest(current)
	merged := applyPatch(stored, patch)

	// An asset stored before this endpoint existed may carry no address. A
	// patch must still be able to fix its title, so the address is only
	// demanded when the merged request actually contains one.
	in, err := toWriteDTO(merged, writePatch)
	if err != nil {
		return server.RealEstate{}, err
	}

	state, err := database.GetRealEstateGuardState(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return server.RealEstate{}, ErrRealEstateNotFound
		}
		return server.RealEstate{}, err
	}

	if state.Deleted {
		return server.RealEstate{}, ErrRealEstateNotFound
	}

	if err := checkSharesConfigChange(state, in.SharesConfig); err != nil {
		return server.RealEstate{}, err
	}

	// The stored asset is converted with the same rules as the merged one, so
	// the comparison of what is missing before and after is meaningful.
	was, err := toWriteDTO(stored, writePatch)
	if err != nil {
		return server.RealEstate{}, err
	}

	if err := checkStillPublishable(state.StatusId, was, in); err != nil {
		return server.RealEstate{}, err
	}

	if err := database.UpdateRealEstate(ctx, id, in); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return server.RealEstate{}, ErrRealEstateNotFound
		}
		return server.RealEstate{}, translateWriteError(err)
	}

	return s.GetRealEstateById(ctx, id, true)
}

// toWriteRequest renders a stored asset in the shape of a write request, so a
// patch can be merged onto it with the create rules applying unchanged.
//
// Amounts go back out as decimal strings: the same exact form they were
// accepted in, never a float.
func toWriteRequest(dto database.RealEstateDTO) server.RealEstateWriteRequest {
	req := server.RealEstateWriteRequest{
		Description: dto.Description,
		Imageurl:    dto.Imageurl,
		IssuerId:    dto.IssuerId,
	}

	// The stored value cannot be NULL — the column is NOT NULL since migration
	// 000012 — but a row written before it can be, and a patch on such a row
	// must not crash. It simply has to name a type, like a creation does.
	if dto.EstateTypeId != nil {
		req.EstateTypeId = *dto.EstateTypeId
	}

	if dto.Title != nil {
		req.Title = *dto.Title
	}

	if a := dto.Address; a != nil {
		req.Address = server.RealEstateAddress{
			Street:           a.Street,
			StreetComplement: a.StreetComplement,
			PostalCode:       a.PostalCode,
			City:             a.City,
			Region:           a.Region,
			CountryCode:      a.CountryCode,
			Latitude:         nullDecimalToString(a.Latitude),
			Longitude:        nullDecimalToString(a.Longitude),
		}
	}

	if sp := dto.Specification; sp != nil {
		req.Specification = &server.RealEstateSpecificationInput{
			SurfaceArea:    nullDecimalToString(sp.SurfaceArea),
			LotSize:        nullDecimalToString(sp.LotSize),
			PoolSize:       nullDecimalToString(sp.PoolSize),
			TerraceSize:    nullDecimalToString(sp.TerraceSize),
			BedroomNumber:  sp.BedroomNumber,
			BathroomNumber: sp.BathroomNumber,
			BuildYear:      sp.BuiltYear,
			EnergyClass:    sp.EnergyClass,
			GesClass:       sp.GesClass,
		}
	}

	if c := dto.SharesConfig; c != nil {
		config := &server.RealEstateConfigurationInput{
			CurrencyCode:           c.CurrencyCode,
			Yield:                  nullDecimalToString(c.Yield),
			PaymentFrequency:       c.PaymentFrequency,
			PaymentFrequencyTypeId: c.PaymentFrequencyTypeId,
			MinInvestment:          nullDecimalToString(c.MinInvestment),
			MaxInvestment:          nullDecimalToString(c.MaxInvestment),
			EntryFeeRate:           nullDecimalToString(c.EntryFeeRate),
			MgmtFeeRate:            nullDecimalToString(c.MgmtFeeRate),
			ExitFeeRate:            nullDecimalToString(c.ExitFeeRate),
			CompartmentRef:         c.CompartmentRef,
		}
		if c.TotalShares.Valid {
			config.TotalShares = c.TotalShares.Decimal.String()
		}
		if c.PricePerShare.Valid {
			config.PricePerShare = c.PricePerShare.Decimal.String()
		}
		req.Configuration = config
	}

	if len(dto.Media) > 0 {
		media := make([]server.RealEstateMediaInput, 0, len(dto.Media))
		for _, m := range dto.Media {
			position := m.Position
			isCover := m.IsCover
			mediaType := m.MediaType
			media = append(media, server.RealEstateMediaInput{
				Url:       m.Url,
				AltText:   m.AltText,
				MediaType: &mediaType,
				Position:  &position,
				IsCover:   &isCover,
			})
		}
		req.Media = &media
	}

	return req
}

// cloneWriteRequest detaches every section a patch can write through, so the
// original stays exactly as it was read from the database.
func cloneWriteRequest(req server.RealEstateWriteRequest) server.RealEstateWriteRequest {
	if s := req.Specification; s != nil {
		copied := *s
		req.Specification = &copied
	}
	if c := req.Configuration; c != nil {
		copied := *c
		req.Configuration = &copied
	}
	if m := req.Media; m != nil {
		copied := append([]server.RealEstateMediaInput(nil), *m...)
		req.Media = &copied
	}

	return req
}

// applyPatch lays the supplied fields over the current asset.
//
// Only a field actually present in the JSON body is applied: everything else
// keeps the value it already had. An empty string is not an omission, it is the
// way to clear an optional value, and it travels down to the same code that
// turns a blank into a NULL on a create.
//
// The request is taken by value, but its sections are pointers: writing through
// them would edit the caller's copy of the stored asset as well. The caller
// compares the asset before and after the patch, so a shared section would make
// the two identical and the comparison blind. Hence the clone.
func applyPatch(current server.RealEstateWriteRequest, patch server.RealEstatePatchRequest) server.RealEstateWriteRequest {
	current = cloneWriteRequest(current)

	if patch.Title != nil {
		current.Title = *patch.Title
	}
	setIfPresent(&current.Description, patch.Description)
	setIfPresent(&current.Imageurl, patch.Imageurl)
	if patch.EstateTypeId != nil {
		current.EstateTypeId = *patch.EstateTypeId
	}
	setIfPresent(&current.IssuerId, patch.IssuerId)

	if a := patch.Address; a != nil {
		if a.Street != nil {
			current.Address.Street = *a.Street
		}
		if a.PostalCode != nil {
			current.Address.PostalCode = *a.PostalCode
		}
		if a.City != nil {
			current.Address.City = *a.City
		}
		if a.CountryCode != nil {
			current.Address.CountryCode = *a.CountryCode
		}
		setIfPresent(&current.Address.StreetComplement, a.StreetComplement)
		setIfPresent(&current.Address.Region, a.Region)
		setIfPresent(&current.Address.Latitude, a.Latitude)
		setIfPresent(&current.Address.Longitude, a.Longitude)
	}

	if sp := patch.Specification; sp != nil {
		if current.Specification == nil {
			current.Specification = &server.RealEstateSpecificationInput{}
		}
		target := current.Specification
		setIfPresent(&target.SurfaceArea, sp.SurfaceArea)
		setIfPresent(&target.LotSize, sp.LotSize)
		setIfPresent(&target.PoolSize, sp.PoolSize)
		setIfPresent(&target.TerraceSize, sp.TerraceSize)
		setIfPresent(&target.BedroomNumber, sp.BedroomNumber)
		setIfPresent(&target.BathroomNumber, sp.BathroomNumber)
		setIfPresent(&target.BuildYear, sp.BuildYear)
		setIfPresent(&target.EnergyClass, sp.EnergyClass)
		setIfPresent(&target.GesClass, sp.GesClass)
	}

	if c := patch.Configuration; c != nil {
		if current.Configuration == nil {
			current.Configuration = &server.RealEstateConfigurationInput{}
		}
		target := current.Configuration
		if c.TotalShares != nil {
			target.TotalShares = *c.TotalShares
		}
		if c.PricePerShare != nil {
			target.PricePerShare = *c.PricePerShare
		}
		setIfPresent(&target.CurrencyCode, c.CurrencyCode)
		setIfPresent(&target.Yield, c.Yield)
		setIfPresent(&target.PaymentFrequency, c.PaymentFrequency)
		setIfPresent(&target.PaymentFrequencyTypeId, c.PaymentFrequencyTypeId)
		setIfPresent(&target.MinInvestment, c.MinInvestment)
		setIfPresent(&target.MaxInvestment, c.MaxInvestment)
		setIfPresent(&target.EntryFeeRate, c.EntryFeeRate)
		setIfPresent(&target.MgmtFeeRate, c.MgmtFeeRate)
		setIfPresent(&target.ExitFeeRate, c.ExitFeeRate)
		setIfPresent(&target.CompartmentRef, c.CompartmentRef)
	}

	// A gallery is a list: there is no field to merge item by item, and the
	// items carry no identifier a patch could target. Supplying the array
	// replaces it, omitting it leaves it alone.
	if patch.Media != nil {
		current.Media = patch.Media
	}

	return current
}

// setIfPresent overwrites the target only when the patch carried the field.
func setIfPresent[T any](target **T, patched *T) {
	if patched != nil {
		*target = patched
	}
}

// DeleteRealEstate soft-deletes the asset.
func (s *Service) DeleteRealEstate(ctx context.Context, id int) error {
	state, err := database.GetRealEstateGuardState(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRealEstateNotFound
		}
		return err
	}

	if state.Deleted {
		return ErrRealEstateNotFound
	}

	// The token exists on chain and investors hold it: nothing off-chain can
	// undo that, so the asset it represents has to stay readable.
	if state.TokenId != nil {
		return conflict("this asset is tokenized and cannot be deleted")
	}

	if err := database.SoftDeleteRealEstate(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRealEstateNotFound
		}
		return err
	}

	return nil
}

// Mirrors of the column lengths of ass.real_estate, as set by migration 000012
// and advertised by api/openapi.yaml. The three definitions have to move
// together; the database stays the one that enforces them.
const (
	maxTitleLength       = 255
	maxDescriptionLength = 2000
	maxImageURLLength    = 500
)

func checkLengths(in database.RealEstateWriteDTO) error {
	// Counted in runes: a 200-character description written in accented French
	// is not twice as long as the same text in English.
	if n := len([]rune(in.Title)); n > maxTitleLength {
		return invalid("title must be at most %d characters, got %d", maxTitleLength, n)
	}
	if in.Description != nil {
		if n := len([]rune(*in.Description)); n > maxDescriptionLength {
			return invalid("description must be at most %d characters, got %d", maxDescriptionLength, n)
		}
	}
	if in.Imageurl != nil {
		if n := len([]rune(*in.Imageurl)); n > maxImageURLLength {
			return invalid("imageurl must be at most %d characters, got %d", maxImageURLLength, n)
		}
	}

	return nil
}

// fkFields names the constraint that a payload can violate, and the field the
// caller has to fix. A raw constraint name means nothing to a client.
var fkFields = map[string]string{
	"real_estate_estate_type_fkey":                          "estate_type_id",
	"real_estate_issuer_id_fkey":                            "issuer_id",
	"real_estate_shares_config_payment_frequency_type_fkey": "configuration.payment_frequency_type_id",
}

// isPrimaryKey recognises a primary key from its name.
//
// The driver does not say which kind of constraint it hit, and asking the
// catalog on every failed write would put a query on the error path. The two
// naming styles present in this database are covered: `pk_...`, used by the
// tables created by hand, and `..._pkey`, the name PostgreSQL generates.
func isPrimaryKey(constraint string) bool {
	return strings.HasPrefix(constraint, "pk_") || strings.HasSuffix(constraint, "_pkey")
}

// translateWriteError turns the constraint violations a payload can cause into
// the sentinel the handler understands.
//
// Referring to a type or an issuer that does not exist is the caller's mistake,
// not a server failure: answering 500 would tell him to retry an identical
// request forever. Anything the database rejects for another reason stays a
// 500, because it means the code let through something it should have checked.
func translateWriteError(err error) error {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return err
	}

	switch pqErr.Code.Name() {
	case "foreign_key_violation":
		if field, ok := fkFields[pqErr.Constraint]; ok {
			return invalid("%s does not refer to an existing record", field)
		}
		return invalid("a referenced record does not exist (%s)", pqErr.Constraint)
	case "unique_violation":
		// A primary key collision is not the caller's doing: the id was handed
		// out by the table's own identity. It means the sequence sits behind
		// the data, which happens as soon as rows are inserted with an explicit
		// id. Answering 409 sent the manager looking for a duplicate value in a
		// payload that never carried one.
		if isPrimaryKey(pqErr.Constraint) {
			return fmt.Errorf(
				"the identity sequence of the table is behind its data (%s, %s): replay migration 000013 to resynchronise it",
				pqErr.Constraint, pqErr.Detail)
		}
		// A compartment reference identifies a ring-fenced estate: two assets
		// claiming the same one is a state conflict, not a typo.
		if pqErr.Constraint == "real_estate_shares_config_compartment_ref_uk" {
			return conflict("this compartment reference is already used by another asset")
		}
		return conflict("this value is already used (%s)", pqErr.Constraint)
	case "check_violation":
		return invalid("a value is out of the range allowed by the database (%s)", pqErr.Constraint)
	case "string_data_right_truncation":
		// The service checks the lengths it knows about; this catches the
		// column whose limit it does not, and still names the problem.
		return invalid("a value is longer than the database allows (%s)", pqErr.Column)
	case "not_null_violation":
		// Schema drift: a column the contract treats as optional is mandatory
		// in the table. Naming it turns an opaque 500 into something actionable.
		return invalid("%s cannot be empty in this database", pqErr.Column)
	}

	return err
}

// checkSharesConfigChange refuses the changes that would rewrite what an
// investor already owns.
//
// Two distinct reasons, deliberately not merged:
//
//   - the asset is tokenized: total_shares and price_per_share describe a
//     contract deployed on chain. Changing them off-chain would only make the
//     database disagree with the truth.
//   - shares are already reserved: cutting the total below what was sold would
//     leave more shares held than issued, and repricing them would silently
//     change the amount investors committed to.
func checkSharesConfigChange(state database.RealEstateGuardStateDTO, in *database.RealEstateSharesConfigWriteDTO) error {
	tokenized := state.TokenId != nil
	reserved := state.ReservedShares > 0

	if !tokenized && !reserved {
		return nil
	}

	if in == nil {
		return conflict("the share configuration cannot be removed once the asset is tokenized or has orders")
	}

	if !state.TotalShares.Valid {
		return nil
	}

	current := state.TotalShares.Decimal

	if tokenized {
		if !in.TotalShares.Equal(current) {
			return conflict("total_shares is fixed by the deployed token and cannot be changed")
		}
	} else if in.TotalShares.LessThan(decimal.NewFromInt(state.ReservedShares)) {
		return conflict("total_shares cannot be lowered to %s: %d shares are already reserved",
			in.TotalShares.String(), state.ReservedShares)
	}

	return nil
}

// toWriteDTO validates the request and converts it into the database DTO.
//
// This is where the decimal strings of the contract become decimal.Decimal: a
// value that cannot be parsed is a 400, never a silent zero.
func toWriteDTO(req server.RealEstateWriteRequest, mode writeMode) (database.RealEstateWriteDTO, error) {
	in := database.RealEstateWriteDTO{
		Title:        strings.TrimSpace(req.Title),
		Description:  trimmedPtr(req.Description),
		Imageurl:     trimmedPtr(req.Imageurl),
		EstateTypeId: &req.EstateTypeId,
		IssuerId:     req.IssuerId,
	}

	if in.Title == "" {
		return in, invalid("title is required")
	}

	// The column lengths of ass.real_estate are mirrored here so an oversized
	// value comes back as a 400 naming the field, instead of the 500 that a
	// string_data_right_truncation would produce.
	if err := checkLengths(in); err != nil {
		return in, err
	}

	// estate_type is NOT NULL in the table: an asset without a type cannot be
	// stored at all, so the contract demands it rather than letting the insert
	// fail on a constraint the caller cannot read. Demanded on creation only,
	// for the reason given on writeMode.
	if mode == writeCreate && (in.EstateTypeId == nil || *in.EstateTypeId <= 0) {
		return in, invalid("estate_type_id is required")
	}

	if mode == writeCreate || !isBlankAddress(req.Address) {
		address, err := toAddressWriteDTO(req.Address)
		if err != nil {
			return in, err
		}
		in.Address = &address
	}

	if req.Specification != nil {
		spec, err := toSpecificationWriteDTO(*req.Specification)
		if err != nil {
			return in, err
		}
		in.Specification = &spec
	}

	if req.Configuration != nil {
		config, err := toSharesConfigWriteDTO(*req.Configuration)
		if err != nil {
			return in, err
		}
		in.SharesConfig = &config
	}

	if req.Media != nil {
		media, err := toMediaWriteDTO(*req.Media)
		if err != nil {
			return in, err
		}
		in.Media = media
	}

	return in, nil
}

func toAddressWriteDTO(a server.RealEstateAddress) (database.RealEstateAddressWriteDTO, error) {
	out := database.RealEstateAddressWriteDTO{
		Street:           strings.TrimSpace(a.Street),
		StreetComplement: trimmedPtr(a.StreetComplement),
		PostalCode:       strings.TrimSpace(a.PostalCode),
		City:             strings.TrimSpace(a.City),
		Region:           trimmedPtr(a.Region),
		CountryCode:      strings.ToUpper(strings.TrimSpace(a.CountryCode)),
	}

	for label, value := range map[string]string{
		"address.street":      out.Street,
		"address.postal_code": out.PostalCode,
		"address.city":        out.City,
	} {
		if value == "" {
			return out, invalid("%s is required", label)
		}
	}

	if !countryCodeRe.MatchString(out.CountryCode) {
		return out, invalid("address.country_code must be an ISO 3166-1 alpha-2 code, got %q", a.CountryCode)
	}

	lat, err := parseOptionalDecimal("address.latitude", a.Latitude)
	if err != nil {
		return out, err
	}
	lng, err := parseOptionalDecimal("address.longitude", a.Longitude)
	if err != nil {
		return out, err
	}

	// Constraint real_estate_address_coordinates_ck refuses a half pair. The
	// same rule is enforced here so the caller reads why rather than a raw
	// constraint name.
	if lat.Valid != lng.Valid {
		return out, invalid("address.latitude and address.longitude must be set together")
	}
	if lat.Valid {
		if lat.Decimal.Abs().GreaterThan(decimal.NewFromInt(90)) {
			return out, invalid("address.latitude must be between -90 and 90")
		}
		if lng.Decimal.Abs().GreaterThan(decimal.NewFromInt(180)) {
			return out, invalid("address.longitude must be between -180 and 180")
		}
	}

	out.Latitude = lat
	out.Longitude = lng

	return out, nil
}

func toSpecificationWriteDTO(s server.RealEstateSpecificationInput) (database.RealEstateSpecificationWriteDTO, error) {
	out := database.RealEstateSpecificationWriteDTO{
		BedroomNumber:  s.BedroomNumber,
		BathroomNumber: s.BathroomNumber,
		BuiltYear:      s.BuildYear,
	}

	measures := []struct {
		label  string
		raw    *string
		target *decimal.NullDecimal
	}{
		{"specification.surface_area", s.SurfaceArea, &out.SurfaceArea},
		{"specification.lot_size", s.LotSize, &out.LotSize},
		{"specification.pool_size", s.PoolSize, &out.PoolSize},
		{"specification.terrace_size", s.TerraceSize, &out.TerraceSize},
	}

	for _, m := range measures {
		value, err := parseOptionalDecimal(m.label, m.raw)
		if err != nil {
			return out, err
		}
		if value.Valid && value.Decimal.IsNegative() {
			return out, invalid("%s cannot be negative", m.label)
		}
		*m.target = value
	}

	for label, count := range map[string]*int{
		"specification.bedroom_number":  s.BedroomNumber,
		"specification.bathroom_number": s.BathroomNumber,
	} {
		if count != nil && *count < 0 {
			return out, invalid("%s cannot be negative", label)
		}
	}

	for label, raw := range map[string]*string{
		"specification.energy_class": s.EnergyClass,
		"specification.ges_class":    s.GesClass,
	} {
		if raw == nil {
			continue
		}
		class := strings.ToUpper(strings.TrimSpace(*raw))
		if class == "" {
			continue
		}
		if !energyClassRe.MatchString(class) {
			return out, invalid("%s must be a single letter from A to G, got %q", label, *raw)
		}
		if label == "specification.energy_class" {
			out.EnergyClass = &class
		} else {
			out.GesClass = &class
		}
	}

	return out, nil
}

func toSharesConfigWriteDTO(c server.RealEstateConfigurationInput) (database.RealEstateSharesConfigWriteDTO, error) {
	out := database.RealEstateSharesConfigWriteDTO{
		CurrencyCode:           "EUR",
		PaymentFrequency:       c.PaymentFrequency,
		PaymentFrequencyTypeId: c.PaymentFrequencyTypeId,
		CompartmentRef:         trimmedPtr(c.CompartmentRef),
	}

	totalShares, err := decimal.NewFromString(strings.TrimSpace(c.TotalShares))
	if err != nil {
		return out, invalid("configuration.total_shares must be a decimal string, got %q", c.TotalShares)
	}
	// A share is the unit an investor buys: half of one cannot be issued, and
	// NUMERIC(20,0) would round it away without saying so.
	if !totalShares.Equal(totalShares.Truncate(0)) {
		return out, invalid("configuration.total_shares must be a whole number, got %q", c.TotalShares)
	}
	if totalShares.LessThanOrEqual(decimal.Zero) {
		return out, invalid("configuration.total_shares must be greater than zero")
	}
	out.TotalShares = totalShares

	pricePerShare, err := decimal.NewFromString(strings.TrimSpace(c.PricePerShare))
	if err != nil {
		return out, invalid("configuration.price_per_share must be a decimal string, got %q", c.PricePerShare)
	}
	if pricePerShare.LessThanOrEqual(decimal.Zero) {
		return out, invalid("configuration.price_per_share must be greater than zero")
	}
	// The column is NUMERIC(20,8): anything finer would be rounded on write,
	// and the total valuation would then disagree with the price shown.
	if !pricePerShare.Equal(pricePerShare.Truncate(8)) {
		return out, invalid("configuration.price_per_share cannot hold more than 8 decimals")
	}
	out.PricePerShare = pricePerShare

	if c.CurrencyCode != nil && strings.TrimSpace(*c.CurrencyCode) != "" {
		currency := strings.ToUpper(strings.TrimSpace(*c.CurrencyCode))
		if !currencyCodeRe.MatchString(currency) {
			return out, invalid("configuration.currency_code must be an ISO 4217 code, got %q", *c.CurrencyCode)
		}
		out.CurrencyCode = currency
	}

	rates := []struct {
		label  string
		raw    *string
		target *decimal.NullDecimal
	}{
		{"configuration.yield", c.Yield, &out.Yield},
		{"configuration.entry_fee_rate", c.EntryFeeRate, &out.EntryFeeRate},
		{"configuration.mgmt_fee_rate", c.MgmtFeeRate, &out.MgmtFeeRate},
		{"configuration.exit_fee_rate", c.ExitFeeRate, &out.ExitFeeRate},
	}

	for _, r := range rates {
		value, err := parseOptionalDecimal(r.label, r.raw)
		if err != nil {
			return out, err
		}
		if value.Valid && (value.Decimal.IsNegative() || value.Decimal.GreaterThan(decimal.NewFromInt(100))) {
			return out, invalid("%s must be a percentage between 0 and 100", r.label)
		}
		*r.target = value
	}

	minInvestment, err := parseOptionalDecimal("configuration.min_investment", c.MinInvestment)
	if err != nil {
		return out, err
	}
	maxInvestment, err := parseOptionalDecimal("configuration.max_investment", c.MaxInvestment)
	if err != nil {
		return out, err
	}
	if minInvestment.Valid && maxInvestment.Valid && minInvestment.Decimal.GreaterThan(maxInvestment.Decimal) {
		return out, invalid("configuration.min_investment cannot exceed configuration.max_investment")
	}
	out.MinInvestment = minInvestment
	out.MaxInvestment = maxInvestment

	if c.PaymentFrequency != nil && *c.PaymentFrequency <= 0 {
		return out, invalid("configuration.payment_frequency must be greater than zero")
	}

	return out, nil
}

func toMediaWriteDTO(media []server.RealEstateMediaInput) ([]database.RealEstateMediaWriteDTO, error) {
	out := make([]database.RealEstateMediaWriteDTO, 0, len(media))
	covers := 0

	for i, m := range media {
		item := database.RealEstateMediaWriteDTO{
			Url:       strings.TrimSpace(m.Url),
			AltText:   trimmedPtr(m.AltText),
			MediaType: "image",
		}

		if item.Url == "" {
			return nil, invalid("media[%d].url is required", i)
		}

		if m.MediaType != nil && strings.TrimSpace(*m.MediaType) != "" {
			mediaType := strings.ToLower(strings.TrimSpace(*m.MediaType))
			if !allowedMediaTypes[mediaType] {
				return nil, invalid("media[%d].media_type must be image, video, floorplan or virtual_tour, got %q", i, *m.MediaType)
			}
			item.MediaType = mediaType
		}

		if m.Position != nil {
			if *m.Position < 0 {
				return nil, invalid("media[%d].position cannot be negative", i)
			}
			item.Position = *m.Position
		} else {
			item.Position = i
		}

		if m.IsCover != nil && *m.IsCover {
			item.IsCover = true
			covers++
		}

		out = append(out, item)
	}

	// Migration 000005 enforces this with a partial unique index. Catching it
	// here names the offending field instead of surfacing an index name.
	if covers > 1 {
		return nil, invalid("only one media can be the cover, got %d", covers)
	}

	return out, nil
}

// parseOptionalDecimal turns an absent or empty string into a NULL, and refuses
// anything that is not a number. An unparseable amount is never defaulted to
// zero: that would store a price nobody typed.
func parseOptionalDecimal(label string, raw *string) (decimal.NullDecimal, error) {
	if raw == nil {
		return decimal.NullDecimal{}, nil
	}

	trimmed := strings.TrimSpace(*raw)
	if trimmed == "" {
		return decimal.NullDecimal{}, nil
	}

	value, err := decimal.NewFromString(trimmed)
	if err != nil {
		return decimal.NullDecimal{}, invalid("%s must be a decimal string, got %q", label, *raw)
	}

	return decimal.NullDecimal{Decimal: value, Valid: true}, nil
}

func trimmedPtr(raw *string) *string {
	if raw == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*raw)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}
