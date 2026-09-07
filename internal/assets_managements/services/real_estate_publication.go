package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/assets_managements/database"
)

// publicationRequirements is the single definition of a record complete enough
// to be shown to an investor.
//
// It is deliberately one list, consulted from the three places that need it —
// the creation, which derives the initial status from it; the publication,
// which refuses without it; and the patch, which refuses to strip a published
// asset of it. A second list would drift from this one, and the drift would
// show up as an asset visible on the site that the tokenization then refuses.
//
// What it demands, and why:
//
//   - the economics (shares, price, currency, yield and the rhythm it is paid
//     at): without them there is no
//     offer, only a photograph. An investor cannot evaluate, and the front has
//     nothing to put on the card;
//   - the issuer: an offer is made by a legal vehicle. Publishing one without
//     naming it is the kind of shortcut a regulator asks about;
//   - a description and a visual: a listing with neither is not browsable.
//
// Naming an issuer is not enough — it must also still stand behind the offer.
// That condition is checked by issuerStandsBehind rather than here: it is not a
// field the record is missing but the state of another record, and reading it
// costs a round trip this pure function is called too often to afford.
//
// Surface is deliberately absent: it helps compare assets, it does not prevent
// the offer from being understood, and demanding it would push a manager to
// invent a number to get past the check.
func publicationRequirements(in database.RealEstateWriteDTO) []string {
	var missing []string

	if in.IssuerId == nil {
		missing = append(missing, "issuer_id")
	}

	if strings.TrimSpace(derefString(in.Description)) == "" {
		missing = append(missing, "description")
	}

	// Either source of a visual is accepted: the gallery is the model going
	// forward, imageurl is what the assets created before it still carry.
	if len(in.Media) == 0 && strings.TrimSpace(derefString(in.Imageurl)) == "" {
		missing = append(missing, "media")
	}

	c := in.SharesConfig
	if c == nil {
		return append(missing,
			"configuration.total_shares",
			"configuration.price_per_share",
			"configuration.currency_code",
			"configuration.yield",
			"configuration.payment_frequency",
			"configuration.payment_frequency_type_id",
		)
	}

	// The values are already known to be positive when present — toWriteDTO
	// refuses zero and negative. What is checked here is that they are there.
	if !c.TotalShares.IsPositive() {
		missing = append(missing, "configuration.total_shares")
	}
	if !c.PricePerShare.IsPositive() {
		missing = append(missing, "configuration.price_per_share")
	}
	if strings.TrimSpace(c.CurrencyCode) == "" {
		missing = append(missing, "configuration.currency_code")
	}
	if !c.Yield.Valid {
		missing = append(missing, "configuration.yield")
	}
	// A yield with no rhythm cannot be compared to another offer, so the two
	// travel together: "4.25 %" is only meaningful next to "paid once a year".
	if c.PaymentFrequency == nil {
		missing = append(missing, "configuration.payment_frequency")
	}
	if c.PaymentFrequencyTypeId == nil {
		missing = append(missing, "configuration.payment_frequency_type_id")
	}

	return missing
}

// issuerStandsBehind reports whether the vehicle named by an asset is still
// active, and therefore able to carry an offer.
//
// An asset with no issuer at all is not this function's business:
// publicationRequirements already reports issuer_id as missing, and saying it
// twice in two different words would only make the answer harder to act on.
//
// An issuer_id pointing at nothing is treated as "does not stand behind"
// rather than as a failure: the foreign key makes it impossible, and if the row
// vanished anyway, refusing to publish is the safe reading.
func issuerStandsBehind(ctx context.Context, issuerID *int) (bool, error) {
	if issuerID == nil {
		return true, nil
	}

	status, err := database.GetIssuerStatus(ctx, *issuerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return status == database.IssuerStatusActive, nil
}

// checkIssuerStandsBehind refuses to move a visible asset onto a vehicle that
// is not active.
//
// This closes the symmetry of the issuer guards: withdrawing a vehicle that
// carries published assets was already refused, but nothing stopped a published
// asset from being pointed at a vehicle already withdrawn — the same forbidden
// state reached from the other side.
//
// Only a CHANGE of issuer is checked, for the same reason checkStillPublishable
// compares before and after: an asset published under a vehicle that has since
// been suspended must stay editable, or every patch would be refused —
// including the one moving it to a vehicle that does stand behind it.
//
// That asset also stays visible, which is a deliberate gap: the invariant held
// here is the one on ENTRY. Taking already published assets off the market when
// their issuer is withdrawn is wanted, but not yet decided — see PROGRESS.md
// §11.23 for what it would have to settle first (ongoing fundraisings, and
// whether the effect is a read-time mask or a status change).
func checkIssuerStandsBehind(ctx context.Context, statusID int, before, after database.RealEstateWriteDTO) error {
	if !publicStatuses[statusID] {
		return nil
	}

	if samePtrInt(before.IssuerId, after.IssuerId) {
		return nil
	}

	standing, err := issuerStandsBehind(ctx, after.IssuerId)
	if err != nil {
		return err
	}

	if !standing {
		return conflict("a published asset cannot be moved onto an issuer that is not active")
	}

	return nil
}

func samePtrInt(a, b *int) bool {
	if a == nil || b == nil {
		return a == b
	}

	return *a == *b
}

func derefString(v *string) string {
	if v == nil {
		return ""
	}

	return *v
}

// initialStatus decides where a newly created asset starts.
//
// A manager who fills everything in one go gets his asset online immediately,
// which is what "create this property" is expected to do. One who starts a
// stub gets a draft: invisible to investors, editable, and publishable later.
// Nothing incomplete ever reaches the public listing, and nothing complete
// waits for a step the manager did not ask for.
func initialStatus(in database.RealEstateWriteDTO) int {
	if len(publicationRequirements(in)) > 0 {
		return database.StatusDraft
	}

	return database.StatusPublished
}

// publicStatuses are the states in which investors can see the asset. Taken
// from is_public of migration 000006; a CHECK cannot query the referential, and
// neither can this code without a round trip on every write.
var publicStatuses = map[int]bool{3: true, 4: true, 5: true, 6: true}

// checkStillPublishable refuses a change that would take away from a visible
// asset something investors rely on.
//
// Refusing is the deliberate choice over silently sending the asset back to
// draft: an asset can be under an ongoing fundraising, and taking it off the
// site because a field was cleared by mistake is far more brutal than asking
// for the field back.
//
// What is compared is before and after, not the result alone. Assets published
// before these requirements existed do not satisfy them, and judging only the
// result would freeze them: every patch would be refused, including the one
// that completes them. Only a requirement that was met and no longer is counts
// as a loss.
func checkStillPublishable(statusID int, before, after database.RealEstateWriteDTO) error {
	if !publicStatuses[statusID] {
		return nil
	}

	alreadyMissing := make(map[string]bool)
	for _, field := range publicationRequirements(before) {
		alreadyMissing[field] = true
	}

	var lost []string
	for _, field := range publicationRequirements(after) {
		if !alreadyMissing[field] {
			lost = append(lost, field)
		}
	}

	if len(lost) > 0 {
		return conflict(
			"a published asset cannot lose what investors rely on, missing: %s",
			strings.Join(lost, ", "))
	}

	return nil
}

// PublishRealEstate makes a draft visible to investors.
//
// The completeness of the STORED asset is what is checked, not of a payload:
// publication is a decision taken on the record as it is, possibly long after
// it was filled in, possibly by someone else.
func (s *Service) PublishRealEstate(ctx context.Context, id int) (server.RealEstate, error) {
	current, err := database.GetRealEstateById(ctx, id, true)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return server.RealEstate{}, ErrRealEstateNotFound
		}
		return server.RealEstate{}, err
	}

	state, err := database.GetRealEstateGuardState(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return server.RealEstate{}, ErrRealEstateNotFound
		}
		return server.RealEstate{}, err
	}

	// A deleted asset does not exist for anyone, so it cannot be resurrected by
	// publishing it.
	if state.Deleted {
		return server.RealEstate{}, ErrRealEstateNotFound
	}

	// Reusing the write path to read the stored asset keeps one definition of
	// what an asset is made of. writePatch, not writeCreate: an asset stored
	// before those rules existed must still be publishable once completed.
	in, err := toWriteDTO(toWriteRequest(current), writePatch)
	if err != nil {
		return server.RealEstate{}, err
	}

	if missing := publicationRequirements(in); len(missing) > 0 {
		return server.RealEstate{}, conflict(
			"this asset is not complete enough to be published, missing: %s",
			strings.Join(missing, ", "))
	}

	// Naming a vehicle is not enough: it has to still stand behind the offer.
	// Publishing under a suspended or dissolved issuer would show investors an
	// offer nothing backs — the mirror image of the guard that refuses to
	// withdraw an issuer carrying published assets.
	standing, err := issuerStandsBehind(ctx, in.IssuerId)
	if err != nil {
		return server.RealEstate{}, err
	}
	if !standing {
		return server.RealEstate{}, conflict(
			"the issuer of this asset is not active; activate it before publishing")
	}

	// Publishing an already published asset changes nothing and is not an
	// error: two managers clicking the same button must not produce a failure.
	if err := database.PublishRealEstate(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return server.RealEstate{}, ErrRealEstateNotFound
		}
		return server.RealEstate{}, translateWriteError(err)
	}

	return s.GetRealEstateById(ctx, id, true)
}
