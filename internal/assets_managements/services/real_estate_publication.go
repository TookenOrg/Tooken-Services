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
//   - the economics (shares, price, currency, yield): without them there is no
//     offer, only a photograph. An investor cannot evaluate, and the front has
//     nothing to put on the card;
//   - the issuer: an offer is made by a legal vehicle. Publishing one without
//     naming it is the kind of shortcut a regulator asks about;
//   - a description and a visual: a listing with neither is not browsable.
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

	return missing
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
