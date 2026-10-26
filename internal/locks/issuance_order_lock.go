package locks

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

var (
	ErrOrderAssetNotFound = errors.New("real estate not found") // 404
	ErrOrderNotFound      = errors.New("order not found")       // 404
)

// LockedIssuanceOrderDTO is the order as it is once the lock is held.
type LockedIssuanceOrderDTO struct {
	Id                   int64
	UserId               int
	AssetId              int
	Quantity             int
	StatusId             int
	OrderReference       string
	UnitPrice            decimal.Decimal
	CurrencyCode         string
	GrossAmount          decimal.Decimal
	FeeAmount            decimal.Decimal
	AmountDue            decimal.Decimal
	ReservationExpiresAt sql.NullTime
	ReservationExpired   bool           // computed by PostgreSQL, see ReadLockedIssuanceOrder
	DeliveryTxHash       sql.NullString // the mint that was broadcast, NULL while none was
}

// LockIssuanceOrder takes the row lock (FOR UPDATE) on the order, held until
// tx ends. It only locks: read the order afterwards with
// ReadLockedIssuanceOrder. When both are needed, lock the real estate first
// (LockRealEstate), then the order, to keep a single lock order and avoid
// deadlocks. Returns ErrOrderNotFound for an unknown reference.
func LockIssuanceOrder(ctx context.Context, tx *sql.Tx, orderReference string) error {
	const query = `SELECT 1 FROM iss.issuance_orders WHERE order_reference = $1 FOR UPDATE`

	var one int
	if err := tx.QueryRowContext(ctx, query, orderReference).Scan(&one); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrOrderNotFound
		}
		return fmt.Errorf("lock issuance order %s: %w", orderReference, err)
	}

	return nil
}

// ReadLockedIssuanceOrder reads the order in a statement of its own, issued
// after LockIssuanceOrder, so it sees every change committed before the lock
// was granted. ReservationExpired is judged by PostgreSQL on clock_timestamp()
// (the real time, not the start of the transaction); false when there is no
// reservation date. An unknown reference returns a wrapped sql.ErrNoRows.
func ReadLockedIssuanceOrder(ctx context.Context, tx *sql.Tx, orderReference string) (LockedIssuanceOrderDTO, error) {
	const query = `
SELECT
    io.id,
    io.user_id,
    io.asset_id,
    io.quantity,
    io.status_id,
    io.order_reference,
    io.unit_price,
    io.currency_code,
    io.gross_amount,
    io.fee_amount,
    io.amount_due,
    io.reservation_expires_at,
	COALESCE(io.reservation_expires_at <= clock_timestamp(), false) AS reservation_expired,
    io.delivery_tx_hash
FROM iss.issuance_orders io
WHERE io.order_reference = $1
`

	var o LockedIssuanceOrderDTO
	err := tx.QueryRowContext(ctx, query, orderReference).Scan(
		&o.Id,
		&o.UserId,
		&o.AssetId,
		&o.Quantity,
		&o.StatusId,
		&o.OrderReference,
		&o.UnitPrice,
		&o.CurrencyCode,
		&o.GrossAmount,
		&o.FeeAmount,
		&o.AmountDue,
		&o.ReservationExpiresAt,
		&o.ReservationExpired,
		&o.DeliveryTxHash,
	)
	if err != nil {
		return LockedIssuanceOrderDTO{}, fmt.Errorf("read locked issuance order %s: %w", orderReference, err)
	}

	return o, nil
}
