package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/internal/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/lib/pq"
)

// ErrDuplicateOrderReference is returned when the generated order reference
// collides with the UNIQUE constraint; callers should regenerate and retry.
var ErrDuplicateOrderReference = errors.New("duplicate order reference")

// Name of constraint is PostgreSQL
const orderReferenceConstraint = "issuance_orders_order_reference_uk"

// isDuplicateReferenceErr reports whether err is a Postgres unique-violation
// (SQLSTATE 23505) on the order_reference constraint.
func isDuplicateReferenceErr(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505" && pqErr.Constraint == orderReferenceConstraint
	}
	return false
}

const baseIssuanceOrderQuery = `
SELECT 
	ord.id,
	ord.user_id,
	ord.asset_id,
	ord.quantity,
	ord.created_at,
	ord.updated_at,
	ord.order_reference,
	sts.id,
	sts.code,
	sts.label,
	sts.is_final,
	ass.title,
	ord.unit_price,
	ord.currency_code,
	ord.gross_amount,
	ord.fee_amount,
	ord.amount_due,
	ord.reservation_expires_at
FROM iss.issuance_orders ord
CROSS JOIN LATERAL (
     SELECT CASE
         WHEN ord.status_id = 2 AND ord.reservation_expires_at <= now() THEN 7
         ELSE ord.status_id
     END AS status_id
 ) shown
 LEFT JOIN iss.issuance_order_statuses sts ON sts.id = shown.status_id
 LEFT JOIN ass.real_estate ass ON ass.id = ord.asset_id
     %s
`

func scanIssuanceOrder(scanner interface {
	Scan(dest ...any) error
}) (server.IssuanceOrder, error) {

	var e server.IssuanceOrder

	var (
		title     *string
		isFinal   *bool
		amounts   [5]sql.NullString // unit_price, currency_code, gross, fee, due
		expiresAt sql.NullTime
	)

	err := scanner.Scan(
		&e.Id,
		&e.UserId,
		&e.RealEstateId,
		&e.TokenQuantity,
		&e.CreatedAt,
		&e.UpdatedAt,
		&e.OrderRef,
		&e.StatusId,
		&e.StatusCode,
		&e.StatusLabel,
		&isFinal,
		&title,
		&amounts[0], &amounts[1], &amounts[2], &amounts[3], &amounts[4],
		&expiresAt,
	)
	if err != nil {
		return e, err
	}

	// Real estate title
	e.RealEstateTitle = title

	// All five amounts are set or none is (issuance_orders_amounts_ck): an
	// order created before M3-3 has none and keeps empty strings.
	e.UnitPrice = amounts[0].String
	e.CurrencyCode = amounts[1].String
	e.GrossAmount = amounts[2].String
	e.FeeAmount = amounts[3].String
	e.AmountDue = amounts[4].String
	if expiresAt.Valid {
		e.ReservationExpiresAt = expiresAt.Time
	}

	return e, nil
}

// getIssuanceOrderFees reads the fee lines of an order. The result is never
// nil, so the API answers "fees": [] rather than null when there is no fee.
func getIssuanceOrderFees(ctx context.Context, orderID int) ([]server.IssuanceOrderFee, error) {
	const query = `
SELECT fee_code, rate, base_amount, amount
FROM iss.issuance_order_fees
WHERE order_id = $1
ORDER BY id
`

	rows, err := globals.DB.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("read fees of order %d: %w", orderID, err)
	}
	defer rows.Close()

	fees := []server.IssuanceOrderFee{}
	for rows.Next() {
		var f server.IssuanceOrderFee
		if err := rows.Scan(&f.Code, &f.Rate, &f.BaseAmount, &f.Amount); err != nil {
			return nil, fmt.Errorf("scan fee of order %d: %w", orderID, err)
		}
		fees = append(fees, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read fees of order %d: %w", orderID, err)
	}

	return fees, nil
}

func ListIssuanceOrders(ctx context.Context, viewAll bool, userId int) (orders []server.IssuanceOrder, err error) {

	var rows *sql.Rows

	query := baseIssuanceOrderQuery
	if !viewAll {
		query = fmt.Sprintf(
			query,
			"WHERE ord.user_id = $1 ORDER BY ord.created_at DESC, ord.id DESC",
		)
		rows, err = globals.DB.QueryContext(ctx, query, userId)
	} else {
		query = fmt.Sprintf(query, "ORDER BY ord.created_at DESC, ord.id DESC")
		rows, err = globals.DB.QueryContext(ctx, query)
	}

	if err != nil {
		return nil, fmt.Errorf("list issuance orders for user %d: %w", userId, err)
	}
	defer rows.Close()

	orders = []server.IssuanceOrder{}
	for rows.Next() {
		order, err := scanIssuanceOrder(rows)
		if err != nil {
			return nil, fmt.Errorf("scan issuance order for user %d: %w", userId, err)
		}
		if order.Id != nil {
			order.Fees, err = getIssuanceOrderFees(ctx, *order.Id)
			if err != nil {
				return nil, fmt.Errorf("read fees of order %d: %w", *order.Id, err)
			}
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate issuance orders for user %d: %w", userId, err)
	}

	return
}

func GetIssuanceOrderByRef(ctx context.Context, orderRef string, viewAll bool, userID int) (order server.IssuanceOrder, err error) {

	var row *sql.Row

	query := baseIssuanceOrderQuery

	if !viewAll {
		query = fmt.Sprintf(
			query,
			"WHERE ord.order_reference = $1 AND ord.user_id = $2",
		)
		row = globals.DB.QueryRowContext(ctx, query, orderRef, userID)
	} else {
		query = fmt.Sprintf(query, "WHERE ord.order_reference = $1 ")
		row = globals.DB.QueryRowContext(ctx, query, orderRef)
	}

	order, err = scanIssuanceOrder(row)
	if err != nil {
		return server.IssuanceOrder{}, err
	}

	if order.Id != nil {
		order.Fees, err = getIssuanceOrderFees(ctx, *order.Id)
		if err != nil {
			return server.IssuanceOrder{}, err
		}
	}

	logger.LogDebug("Issuance order found [%s]", utils.Dump(order))

	return order, nil
}
