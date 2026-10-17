package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/locks"
	"github.com/TookenOrg/tooken-services/internal/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/shopspring/decimal"
)

var (
	ErrOrderNotFound           = locks.ErrOrderNotFound                      // 404
	ErrOrderNotAwaitingPayment = errors.New("order is not awaiting payment") // 409
	ErrOrderReservationExpired = errors.New("reservation expired")           // 409
)

func ConfirmIssuanceOrderPayment(ctx context.Context, realEstateId int, orderReference string, userID int) (err error) {
	err = utils.InTransaction(ctx, func(tx *sql.Tx) error {

		// 1 - Locks first
		if err := locks.LockRealEstate(ctx, tx, realEstateId); err != nil {
			return err
		}

		if err := locks.LockIssuanceOrder(ctx, tx, orderReference); err != nil {
			return err
		}

		// 2 - The order, read under the lock.
		order, err := locks.ReadLockedIssuanceOrder(ctx, tx, orderReference)
		if err != nil {
			return err
		}

		// The order must not be already paid for the confirmation to proceed.
		if order.StatusId == int(server.IssuanceOrderStatusIdPaymentConfirmed) {
			logger.LogInfo("Order %s has already been paid", orderReference)
			return nil
		}

		// 3 - Check if the order is awaiting payment.
		if err := checkOrderForPaymentConfirmation(order); err != nil {
			return err
		}

		// 4 - Insert payment
		if err := InsertPayment(ctx, tx, order.Id, userID, order.AmountDue, order.CurrencyCode); err != nil {
			return err
		}

		// 5 - Insert payment history
		if err := InsertOrderStatusHistory(ctx, tx, order.Id, int(server.IssuanceOrderStatusIdAwaitingPayment), int(server.IssuanceOrderStatusIdPaymentConfirmed), userID, "Payment confirmed"); err != nil {
			return err
		}

		// 6 - Update the order status to reflect the payment confirmation.
		err = UpdateIssuanceOrderStatusToPaymentConfirmed(ctx, tx, orderReference)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// checkOrderForPaymentConfirmation verifies that the order is in a state awaiting payment.
func checkOrderForPaymentConfirmation(a locks.LockedIssuanceOrderDTO) error {

	// Status must be 'awaiting payment' for the confirmation to proceed.
	if a.StatusId != int(server.IssuanceOrderStatusIdAwaitingPayment) {
		return ErrOrderNotAwaitingPayment
	}

	// The order must not be expired for the confirmation to proceed.
	if a.ReservationExpired {
		return ErrOrderReservationExpired
	}

	return nil
}

func UpdateIssuanceOrderStatusToPaymentConfirmed(ctx context.Context, tx *sql.Tx, orderReference string) error {
	const query = `
	UPDATE iss.issuance_orders 
	SET 
		status_id = $1,
		updated_at = now() 
	WHERE order_reference = $2`
	result, err := tx.ExecContext(ctx, query, int(server.IssuanceOrderStatusIdPaymentConfirmed), orderReference)
	if err != nil {
		return fmt.Errorf("update issuance order status to payment confirmed for %s: %w", orderReference, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update issuance order status to payment confirmed for %s: %w", orderReference, err)
	}
	if rowsAffected == 0 {
		return ErrOrderNotFound
	}
	return nil
}

func InsertPayment(ctx context.Context, tx *sql.Tx, orderID int64, userID int, amount decimal.Decimal, currencyCode string) error {
	const query = `
	INSERT INTO iss.issuance_order_payments (order_id, kind, provider, amount, currency_code, recorded_by_user_id, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, now())`
	_, err := tx.ExecContext(ctx, query, orderID, "PAYMENT", "MANUAL", amount, currencyCode, userID)
	if err != nil {
		return fmt.Errorf("insert payment for order %d: %w", orderID, err)
	}
	return nil
}

func InsertOrderStatusHistory(ctx context.Context, tx *sql.Tx, orderID int64, fromStatusID, toStatusID, userID int, reason string) error {
	const query = `
	INSERT INTO iss.issuance_order_status_history (order_id, from_status_id, to_status_id, actor_user_id, reason, created_at)
	VALUES ($1, $2, $3, $4, $5, now())`
	_, err := tx.ExecContext(ctx, query, orderID, fromStatusID, toStatusID, userID, reason)
	if err != nil {
		return fmt.Errorf("insert order status history for order %d: %w", orderID, err)
	}
	return nil
}
