package locks

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

// LockedRealEstateDTO is the asset as it is once the lock is held. Every
// business check of the order reads it, never an earlier read.
type LockedRealEstateDTO struct {
	StatusID             int
	Active               bool
	Deleted              bool
	Tokenized            bool
	HasSharesConfig      bool
	TotalShares          int64
	PricePerShare        decimal.Decimal
	CurrencyCode         string
	EntryFeeRate         decimal.NullDecimal // NULL: no entry fee
	MaxSharesPerInvestor sql.NullInt64       // NULL: no cap (U8)
}

// LockRealEstate takes the row lock on the asset.
//
// The query reads nothing else on purpose: when FOR UPDATE has to wait,
// PostgreSQL refreshes the locked row but not the rows of a joined table, which
// stay as they were when the statement started (proven by an integration test).
func LockRealEstate(ctx context.Context, tx *sql.Tx, realEstateID int) error {
	const query = `SELECT 1 FROM ass.real_estate WHERE id = $1 FOR UPDATE`

	var one int
	if err := tx.QueryRowContext(ctx, query, realEstateID).Scan(&one); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrOrderAssetNotFound
		}
		return fmt.Errorf("lock real estate %d: %w", realEstateID, err)
	}

	return nil
}

// ReadLockedRealEstate reads what the order checks, in a statement of its own
// issued after the lock, so it sees every change committed before the lock
// was granted. The join is safe here: no FOR UPDATE, fresh snapshot.
func ReadLockedRealEstate(ctx context.Context, tx *sql.Tx, realEstateID int) (LockedRealEstateDTO, error) {
	const query = `
SELECT
    re.status_id,
    re.active,
    re.deleted_at IS NOT NULL,
    re.token_id IS NOT NULL,
    sc.id IS NOT NULL,
    COALESCE(sc.total_shares, 0),
    COALESCE(sc.price_per_share, 0),
    COALESCE(sc.currency_code, ''),
    sc.entry_fee_rate,
    sc.max_shares_per_investor
FROM ass.real_estate re
LEFT JOIN ass.real_estate_shares_config sc
    ON sc.real_estate_id = re.id
WHERE re.id = $1
`

	var a LockedRealEstateDTO
	err := tx.QueryRowContext(ctx, query, realEstateID).Scan(
		&a.StatusID,
		&a.Active,
		&a.Deleted,
		&a.Tokenized,
		&a.HasSharesConfig,
		&a.TotalShares,
		&a.PricePerShare,
		&a.CurrencyCode,
		&a.EntryFeeRate,
		&a.MaxSharesPerInvestor,
	)
	if err != nil {
		return LockedRealEstateDTO{}, fmt.Errorf("read locked real estate %d: %w", realEstateID, err)
	}

	return a, nil
}
