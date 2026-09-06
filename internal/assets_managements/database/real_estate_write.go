package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/shopspring/decimal"
)

// Write side of ass.real_estate.
//
// An asset is spread over five tables. Writing them one statement at a time
// would let a failure halfway through leave an asset without its address or,
// worse, without its share configuration: readable, priced at nothing, and
// impossible to tell apart from a legitimately incomplete draft. Every write
// here runs in a single transaction.
//
// Amounts and measures are decimal.Decimal all the way to the driver. Passing
// a float64 would put back the rounding that migration 000001 and the DTOs were
// written to remove.

// StatusPublished is the status an asset is created with.
//
// A manager who clicks "create" expects the asset on the site, not a draft
// nobody asked for. The trigger of migration 000006 derives active = true from
// it, so nothing else has to be set.
const StatusPublished = 3

// StatusCancelled is the status a soft-deleted asset carries. Constraint
// real_estate_deleted_status_ck (migration 000011) forbids any other.
const StatusCancelled = 7

// RealEstateWriteDTO is the whole asset as accepted by a write.
//
// Every optional section is a pointer: nil means "this asset has no such
// section", which a write must translate into the absence of a row, not into a
// row full of NULLs.
type RealEstateWriteDTO struct {
	Title        string
	Description  *string
	Imageurl     *string
	EstateTypeId *int
	IssuerId     *int

	// nil means the asset has no address row. A create always supplies one;
	// a patch may leave an inherited asset without any, and refusing the write
	// would make such an asset impossible to edit at all.
	Address       *RealEstateAddressWriteDTO
	Specification *RealEstateSpecificationWriteDTO
	SharesConfig  *RealEstateSharesConfigWriteDTO
	Media         []RealEstateMediaWriteDTO
}

type RealEstateAddressWriteDTO struct {
	Street           string
	StreetComplement *string
	PostalCode       string
	City             string
	Region           *string
	CountryCode      string
	Latitude         decimal.NullDecimal
	Longitude        decimal.NullDecimal
}

type RealEstateSpecificationWriteDTO struct {
	SurfaceArea    decimal.NullDecimal
	LotSize        decimal.NullDecimal
	PoolSize       decimal.NullDecimal
	TerraceSize    decimal.NullDecimal
	BedroomNumber  *int
	BathroomNumber *int
	BuiltYear      *int
	EnergyClass    *string
	GesClass       *string
}

// RealEstateSharesConfigWriteDTO carries no TotalValuation: the column is
// GENERATED ALWAYS and PostgreSQL rejects any attempt to write it. The
// derivation lives in one place, which is why the three values can never
// disagree.
type RealEstateSharesConfigWriteDTO struct {
	TotalShares            decimal.Decimal
	PricePerShare          decimal.Decimal
	CurrencyCode           string
	Yield                  decimal.NullDecimal
	PaymentFrequency       *int
	PaymentFrequencyTypeId *int
	MinInvestment          decimal.NullDecimal
	MaxInvestment          decimal.NullDecimal
	EntryFeeRate           decimal.NullDecimal
	MgmtFeeRate            decimal.NullDecimal
	ExitFeeRate            decimal.NullDecimal
	CompartmentRef         *string
}

type RealEstateMediaWriteDTO struct {
	Url       string
	AltText   *string
	MediaType string
	Position  int
	IsCover   bool
}

// RealEstateGuardStateDTO is everything a business rule needs to decide whether
// a write is allowed, read in one query.
//
// Deleted is separate from StatusId because a cancelled asset and a deleted one
// are not the same thing: the first can be reopened, the second is gone.
type RealEstateGuardStateDTO struct {
	Id             int
	StatusId       int
	TokenId        *int
	Deleted        bool
	ReservedShares int64
	TotalShares    decimal.NullDecimal
}

// GetRealEstateGuardState reads the state a write must be checked against.
//
// It ignores active and status: a manager may perfectly well update a draft,
// and the caller needs to know that a deleted asset exists in order to answer
// 404 rather than pretend the write succeeded.
func GetRealEstateGuardState(ctx context.Context, id int) (RealEstateGuardStateDTO, error) {
	const query = `
SELECT re.id,
       re.status_id,
       re.token_id,
       re.deleted_at IS NOT NULL AS deleted,
       COALESCE(sold.reserved, 0)::bigint AS reserved,
       conf.total_shares
FROM ass.real_estate re
LEFT JOIN ass.real_estate_shares_config conf
    ON conf.real_estate_id = re.id
LEFT JOIN (
    SELECT o.asset_id, SUM(o.quantity) AS reserved
    FROM iss.issuance_orders o
    JOIN iss.issuance_order_statuses s
        ON s.id = o.status_id
    WHERE s.counts_as_reserved
    GROUP BY o.asset_id
) sold
    ON sold.asset_id = re.id
WHERE re.id = $1
`

	var s RealEstateGuardStateDTO
	err := globals.DB.QueryRowContext(ctx, query, id).Scan(
		&s.Id, &s.StatusId, &s.TokenId, &s.Deleted, &s.ReservedShares, &s.TotalShares,
	)
	if err != nil {
		return RealEstateGuardStateDTO{}, err
	}

	return s, nil
}

// CreateRealEstate inserts the asset and all its sections, and returns the new
// identifier.
func CreateRealEstate(ctx context.Context, in RealEstateWriteDTO) (id int, err error) {
	err = inTransaction(ctx, func(tx *sql.Tx) error {
		const query = `
INSERT INTO ass.real_estate (title, description, imageurl, estate_type, issuer_id, status_id, published_at)
VALUES ($1, $2, $3, $4, $5, $6, now())
RETURNING id
`

		if err := tx.QueryRowContext(ctx, query,
			in.Title, in.Description, in.Imageurl, in.EstateTypeId, in.IssuerId, StatusPublished,
		).Scan(&id); err != nil {
			return fmt.Errorf("insert real estate: %w", err)
		}

		return writeSections(ctx, tx, id, in)
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateRealEstate replaces the asset and its sections.
//
// It takes the asset in its entirety, already merged by the caller: the
// sections are deleted then reinserted rather than updated in place, which is
// the only form that can express "this asset now has no specification". The
// number of rows involved is a handful.
//
// It reports sql.ErrNoRows when the asset does not exist or was deleted, so a
// caller can answer 404 without a second round trip.
func UpdateRealEstate(ctx context.Context, id int, in RealEstateWriteDTO) error {
	return inTransaction(ctx, func(tx *sql.Tx) error {
		const query = `
UPDATE ass.real_estate
SET title = $2,
    description = $3,
    imageurl = $4,
    estate_type = $5,
    issuer_id = $6,
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
`

		res, err := tx.ExecContext(ctx, query,
			id, in.Title, in.Description, in.Imageurl, in.EstateTypeId, in.IssuerId,
		)
		if err != nil {
			return fmt.Errorf("update real estate: %w", err)
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return sql.ErrNoRows
		}

		for _, table := range []string{
			"ass.real_estate_address",
			"ass.real_estate_specification",
			"ass.real_estate_shares_config",
			"ass.real_estate_media",
		} {
			if _, err := tx.ExecContext(ctx,
				"DELETE FROM "+table+" WHERE real_estate_id = $1", id,
			); err != nil {
				return fmt.Errorf("clear %s: %w", table, err)
			}
		}

		return writeSections(ctx, tx, id, in)
	})
}

// SoftDeleteRealEstate cancels and dates the asset instead of removing it.
//
// Both columns are set in the same statement: constraint
// real_estate_deleted_status_ck rejects a row dated without being cancelled, so
// splitting them in two would fail on the first one.
//
// Deleting twice reports sql.ErrNoRows: the second call changes nothing, and
// saying so is more useful than pretending it deleted something.
func SoftDeleteRealEstate(ctx context.Context, id int) error {
	const query = `
UPDATE ass.real_estate
SET status_id = $2,
    deleted_at = now(),
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
`

	res, err := globals.DB.ExecContext(ctx, query, id, StatusCancelled)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// writeSections inserts the one-to-one tables and the gallery for an asset
// whose row already exists.
func writeSections(ctx context.Context, tx *sql.Tx, id int, in RealEstateWriteDTO) error {
	if err := insertAddress(ctx, tx, id, in.Address); err != nil {
		return err
	}
	if err := insertSpecification(ctx, tx, id, in.Specification); err != nil {
		return err
	}
	if err := insertSharesConfig(ctx, tx, id, in.SharesConfig); err != nil {
		return err
	}

	return insertMedia(ctx, tx, id, in.Media)
}

func insertAddress(ctx context.Context, tx *sql.Tx, id int, a *RealEstateAddressWriteDTO) error {
	if a == nil {
		return nil
	}

	const query = `
INSERT INTO ass.real_estate_address
    (real_estate_id, street, street_complement, postal_code, city, region, country_code, latitude, longitude)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`

	_, err := tx.ExecContext(ctx, query,
		id, a.Street, a.StreetComplement, a.PostalCode, a.City, a.Region,
		a.CountryCode, a.Latitude, a.Longitude,
	)
	if err != nil {
		return fmt.Errorf("insert address: %w", err)
	}

	return nil
}

func insertSpecification(ctx context.Context, tx *sql.Tx, id int, s *RealEstateSpecificationWriteDTO) error {
	if s == nil {
		return nil
	}

	const query = `
INSERT INTO ass.real_estate_specification
    (real_estate_id, surface_area, lot_size, pool_size, terrace_size,
     bedroom_number, bathroom_number, built_year, energy_class, ges_class)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`

	_, err := tx.ExecContext(ctx, query,
		id, s.SurfaceArea, s.LotSize, s.PoolSize, s.TerraceSize,
		s.BedroomNumber, s.BathroomNumber, s.BuiltYear, s.EnergyClass, s.GesClass,
	)
	if err != nil {
		return fmt.Errorf("insert specification: %w", err)
	}

	return nil
}

func insertSharesConfig(ctx context.Context, tx *sql.Tx, id int, c *RealEstateSharesConfigWriteDTO) error {
	if c == nil {
		return nil
	}

	const query = `
INSERT INTO ass.real_estate_shares_config
    (real_estate_id, total_shares, price_per_share, currency_code, yield,
     payment_frequency, payment_frequency_type, min_investment, max_investment,
     entry_fee_rate, mgmt_fee_rate, exit_fee_rate, compartment_ref)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
`

	_, err := tx.ExecContext(ctx, query,
		id, c.TotalShares, c.PricePerShare, c.CurrencyCode, c.Yield,
		c.PaymentFrequency, c.PaymentFrequencyTypeId, c.MinInvestment, c.MaxInvestment,
		c.EntryFeeRate, c.MgmtFeeRate, c.ExitFeeRate, c.CompartmentRef,
	)
	if err != nil {
		return fmt.Errorf("insert shares config: %w", err)
	}

	return nil
}

func insertMedia(ctx context.Context, tx *sql.Tx, id int, media []RealEstateMediaWriteDTO) error {
	const query = `
INSERT INTO ass.real_estate_media
    (real_estate_id, url, alt_text, media_type, position, is_cover)
VALUES ($1, $2, $3, $4, $5, $6)
`

	for _, m := range media {
		if _, err := tx.ExecContext(ctx, query,
			id, m.Url, m.AltText, m.MediaType, m.Position, m.IsCover,
		); err != nil {
			return fmt.Errorf("insert media %q: %w", m.Url, err)
		}
	}

	return nil
}

// inTransaction runs fn in a transaction and rolls back on any failure,
// including a panic: leaving a transaction open would hold locks on the asset
// until the connection is recycled.
func inTransaction(ctx context.Context, fn func(*sql.Tx) error) (err error) {
	tx, err := globals.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}
