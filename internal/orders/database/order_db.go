package database

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	ord.status_id,
	sts.code,
	sts.label,
	sts.is_final,
	ass.title
FROM  iss.issuance_orders ord
LEFT JOIN iss.issuance_order_statuses sts on sts.id = ord.status_id
LEFT JOIN  ass.real_estate ass ON ass.id = ord.asset_id
    %s
`

func scanIssuanceOrder(scanner interface {
	Scan(dest ...any) error
}) (server.IssuanceOrder, error) {

	var e server.IssuanceOrder

	var (
		title *string
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
		&e.StatusIfFinal,
		&title,
	)
	if err != nil {
		return e, err
	}

	return e, nil
}

func InsertIssuranceOrder(ctx context.Context, realEstateId, quantity, userId int, orderReference string) (id int, createdAt time.Time, err error) {
	query := `
		INSERT INTO iss.issuance_orders (
			user_id,
			asset_id,
			quantity,
			status_id,
			order_reference
		)
		VALUES (
			$1,  
			$2, 
			$3,  
			$4,
			$5
		)
		RETURNING id, created_at;
    `

	err = globals.DB.QueryRowContext(ctx, query, userId, realEstateId, quantity, 1, orderReference).Scan(&id, &createdAt)
	if err != nil {
		if isDuplicateReferenceErr(err) {
			return 0, time.Time{}, ErrDuplicateOrderReference
		}
		return 0, time.Time{}, fmt.Errorf("failed to insert issuance order: %w", err)
	}

	logger.LogInfo("Issuance order inserted. ID=%d Ref=%s", id, orderReference)

	return
}

func GetIssuanceOrderByRef(ctx context.Context, orderRef string) (order server.IssuanceOrder, err error) {
	query := fmt.Sprintf(
		baseIssuanceOrderQuery,
		"WHERE ord.order_reference = $1",
	)

	row := globals.DB.QueryRowContext(ctx, query, orderRef)

	order, err = scanIssuanceOrder(row)
	if err != nil {
		return server.IssuanceOrder{}, err
	}

	logger.LogDebug("Real Estate found [%s]", utils.Dump(order))

	return order, nil
}
