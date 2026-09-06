package database

import (
	"context"
	"database/sql"

	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
)

// TODO: add user isFavorite
// https://github.com/TookenOrg/Tooken-Services/issues/31

type rowScanner interface {
	Scan(dest ...any) error
}

// scanRealEstate fills a DTO from the column table.
//
// The 1-1 sub-structs are allocated BEFORE the scan (the targets point into
// them), then reset to nil when the join returned nothing. This is the only
// place where the absence of a row is interpreted.
func scanRealEstate(scanner rowScanner, cols []realEstateColumn) (RealEstateDTO, error) {
	dto := RealEstateDTO{
		Specification: &RealEstateSpecificationDTO{},
		SharesConfig:  &RealEstateSharesConfigDTO{},
		Location:      &RealEstateLocationDTO{},
	}

	if err := scanner.Scan(scanTargets(&dto, cols)...); err != nil {
		return RealEstateDTO{}, err
	}

	if !dto.specificationPresent {
		dto.Specification = nil
	}
	if !dto.sharesConfigPresent {
		dto.SharesConfig = nil
	}
	if !dto.locationPresent {
		dto.Location = nil
	}

	return dto, nil
}

// GetActiveRealEstates serves the list: no address, no issuer, no token, no
// gallery. A single query, whatever the number of assets.
//
// The ORDER BY is not cosmetic. Without it PostgreSQL guarantees nothing: rows
// come back in whatever order the plan happens to produce, which is physical
// heap order for a sequential scan. Any UPDATE rewrites the row at the end of
// the heap, so a single migration is enough to silently reshuffle the list —
// and a parallel scan would do the same. It is also a prerequisite for the
// pagination to come: without a total order, LIMIT/OFFSET can skip a row or
// serve it twice.
//
// includeUnpublished widens the filter to the assets that are not published
// yet, for a MANAGER or an ADMIN. Deleted assets stay out either way: they are
// kept for the audit trail, not to be browsed.
func GetRealEstates(ctx context.Context, includeUnpublished bool) ([]RealEstateDTO, error) {
	cols := realEstateCoreColumns
	query := buildRealEstateQuery(cols,
		"WHERE re.deleted_at IS NULL AND (re.active = true OR $1)\nORDER BY re.id")

	rows, err := globals.DB.QueryContext(ctx, query, includeUnpublished)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var estates []RealEstateDTO
	for rows.Next() {
		e, err := scanRealEstate(rows, cols)
		if err != nil {
			return nil, err
		}
		estates = append(estates, e)
	}

	return estates, rows.Err()
}

// GetActiveRealEstateById serves the detail: the full asset, address and
// issuer included.
//
// The remaining 1-1 tables and the gallery are loaded by dedicated queries
// rather than by extra LEFT JOINs. On a single record the cost is negligible,
// and each DTO keeps the exact typing of its schema instead of being made
// nullable by the join.
//
// includeUnpublished lets a MANAGER or an ADMIN open an asset that is not
// published yet. For anyone else such an asset yields sql.ErrNoRows, so a draft
// answers 404 exactly like an identifier that never existed: its existence is
// not disclosed.
func GetRealEstateById(ctx context.Context, id int, includeUnpublished bool) (RealEstateDTO, error) {
	cols := realEstateCoreColumns
	query := buildRealEstateQuery(cols,
		"WHERE re.id = $1 AND re.deleted_at IS NULL AND (re.active = true OR $2)")

	row := globals.DB.QueryRowContext(ctx, query, id, includeUnpublished)

	estate, err := scanRealEstate(row, cols)
	if err != nil {
		return RealEstateDTO{}, err
	}

	if estate.Address, err = getRealEstateAddress(ctx, id); err != nil {
		return RealEstateDTO{}, err
	}
	if estate.Issuer, err = getRealEstateIssuer(ctx, estate.IssuerId); err != nil {
		return RealEstateDTO{}, err
	}
	if estate.Token, err = getRealEstateToken(ctx, estate.TokenId); err != nil {
		return RealEstateDTO{}, err
	}
	if estate.Media, err = getRealEstateMedia(ctx, id); err != nil {
		return RealEstateDTO{}, err
	}

	logger.LogDebug("Real Estate found!")

	return estate, nil
}

// getRealEstateAddress returns nil when the asset has no address recorded: a
// draft record is a normal case, not an error.
func getRealEstateAddress(ctx context.Context, realEstateID int) (*RealEstateAddressDTO, error) {
	const query = `
SELECT street, street_complement, postal_code, city, region,
       country_code, latitude, longitude, created_at, updated_at
FROM ass.real_estate_address
WHERE real_estate_id = $1
`

	var a RealEstateAddressDTO
	err := globals.DB.QueryRowContext(ctx, query, realEstateID).Scan(
		&a.Street, &a.StreetComplement, &a.PostalCode, &a.City, &a.Region,
		&a.CountryCode, &a.Latitude, &a.Longitude, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// getRealEstateIssuer: issuer_id stays nullable as long as the asset is not
// attached to any legal vehicle.
func getRealEstateIssuer(ctx context.Context, issuerID *int) (*RealEstateIssuerDTO, error) {
	if issuerID == nil {
		return nil, nil
	}

	const query = `
SELECT i.id, i.name, i.legal_form, i.registration_number,
       i.country_code, i.lei_code, i.status_id, ist.code
FROM ass.issuer i
JOIN ass.issuer_status ist
    ON ist.id = i.status_id
WHERE i.id = $1
`

	var s RealEstateIssuerDTO
	err := globals.DB.QueryRowContext(ctx, query, *issuerID).Scan(
		&s.Id, &s.Name, &s.LegalForm, &s.RegistrationNumber,
		&s.CountryCode, &s.LeiCode, &s.StatusId, &s.StatusCode,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// getRealEstateToken: a NULL token_id means "not tokenized yet"
// (migration 000008).
func getRealEstateToken(ctx context.Context, tokenID *int) (*RealEstateTokenDTO, error) {
	if tokenID == nil {
		return nil, nil
	}

	const query = `
SELECT id, symbol, token_name, address, nb_decimal,
       modular_compliance_addr, created_at
FROM blk.token
WHERE id = $1
`

	var t RealEstateTokenDTO
	err := globals.DB.QueryRowContext(ctx, query, *tokenID).Scan(
		&t.Id, &t.Symbol, &t.TokenName, &t.Address, &t.NbDecimal,
		&t.ModularComplianceAddr, &t.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// getRealEstateMedia returns the ordered gallery, cover first.
func getRealEstateMedia(ctx context.Context, realEstateID int) ([]RealEstateMediaDTO, error) {
	const query = `
SELECT id, url, alt_text, media_type, position, is_cover, created_at
FROM ass.real_estate_media
WHERE real_estate_id = $1
ORDER BY is_cover DESC, position, id
`

	rows, err := globals.DB.QueryContext(ctx, query, realEstateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var media []RealEstateMediaDTO
	for rows.Next() {
		var m RealEstateMediaDTO
		if err := rows.Scan(
			&m.Id, &m.Url, &m.AltText, &m.MediaType,
			&m.Position, &m.IsCover, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		media = append(media, m)
	}

	return media, rows.Err()
}
