package database

import (
	"context"
	"fmt"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/internal/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
)

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

func InsertIssuranceOrder(ctx context.Context, realEstateId, quantity, userId int) (id int, createdAt time.Time, err error) {
	query := `
		INSERT INTO iss.issuance_orders (
			user_id,
			asset_id,
			quantity,
			status_id
		)
		VALUES (
			$1,  
			$2, 
			$3,  
			$4   
		)
		RETURNING id, created_at;
    `

	err = globals.DB.QueryRow(query, userId, realEstateId, quantity, 1).Scan(&id, &createdAt)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to insert issuance order: %w", err)
	}

	logger.LogInfo("Issuance order inserted. ID=%d", id)

	return
}

func UpdateReferenceIssuanceOrder(ctx context.Context, orderId int, orderReference string) (err error) {
	query := `
		UPDATE iss.issuance_orders
		SET order_reference = $1
		WHERE id = $2;
    `

	res, err := globals.DB.Exec(query, orderReference, orderId)
	if err != nil {
		return logger.LogError("failed to insert issuance order: %v", err)
	}

	updatedRow, err := res.RowsAffected()
	if err != nil {
		return
	}

	if updatedRow != 1 {
		return logger.LogError("No row updated or more than one raw updated")
	}

	logger.LogInfo("Issuance order updated with reference %s", orderReference)

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
