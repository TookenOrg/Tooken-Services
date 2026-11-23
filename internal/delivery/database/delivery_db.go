package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/internal/locks"
	ordrDb "github.com/TookenOrg/tooken-services/internal/orders/database"
	usrDb "github.com/TookenOrg/tooken-services/internal/users/database"
	"github.com/TookenOrg/tooken-services/internal/utils"
)

var (
	// ErrOrderNotDeliverable means the order is not in PAYMENT_CONFIRMED any
	// more: nothing was written, and the caller must stop without raising an
	// alarm. This is what makes a second Deliver on the same order a no-op.
	ErrOrderNotDeliverable = errors.New("order is not deliverable")

	// ErrDeliveryRefused means the order was moved to DELIVERY_FAILED before any
	// on-chain action. The transaction committed, so the failure is
	// recorded; the reason is wrapped in the returned error.
	ErrDeliveryRefused = errors.New("delivery refused")

	// ErrDeliveryAlreadyRecorded means another delivery of the same order got
	// its hash in first. Nothing must be sent: that other one owns the mint.
	ErrDeliveryAlreadyRecorded = errors.New("a delivery transaction is already recorded for this order")
)

type DeliveryTargetDTO struct {
	OrderID      int64
	UserID       int
	Shares       int    // iss.issuance_orders.quantity
	Wallet       string // blk.user_wallet.wallet_address (is_active)
	TokenAddress string // ass.real_estate.contract_address
	RealEstateID int    // iss.issuance_orders.asset_id
	TotalShares  int64  // ass.real_estate_shares_config.total_shares
}

// StartDelivery runs transaction A: it locks the order, checks it is still
// deliverable, reads the mint target, re-checks the KYC, and moves the order to
// DELIVERY_IN_PROGRESS.
//
// An order that is not in PAYMENT_CONFIRMED yields ErrOrderNotDeliverable and
// nothing is written. A refusal found before any on-chain action (no active
// wallet, asset without a token, invalid KYC) is recorded as DELIVERY_FAILED
// and yields ErrDeliveryRefused: the transaction must commit for that record to
// survive, so the refusal is carried out of the closure rather than returned
// from it. Any other error rolls the transaction back and leaves the order
// deliverable.
func StartDelivery(ctx context.Context, orderRef string) (DeliveryTargetDTO, error) {

	var deliveryTarget DeliveryTargetDTO
	var refusal string

	err := utils.InTransaction(ctx, func(tx *sql.Tx) error {
		// 1 - Lock issuance order
		if err := locks.LockIssuanceOrder(ctx, tx, orderRef); err != nil {
			return err
		}

		// 2 - Read locked order
		order, err := locks.ReadLockedIssuanceOrder(ctx, tx, orderRef)
		if err != nil {
			return err
		}

		// 2A - Only a confirmed payment starts a delivery. Anything else means
		// the order was already taken care of: leave it untouched.
		if order.StatusId != int(server.IssuanceOrderStatusIdPaymentConfirmed) {
			return ErrOrderNotDeliverable
		}

		// 3 - Read the active wallet and the token address of the asset.
		wallet, tokenAddress, err := readWalletAndTokenAddressForMint(ctx, tx, order.Id)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			refusal = "the investor has no active wallet"
		case err != nil:
			return err
		case tokenAddress == nil:
			refusal = "the asset carries no token contract address"
		}

		// 3A - Read the size of the offering, to know later whether it is sold out.
		var totalShares int64
		if refusal == "" {
			totalShares, err = readTotalShares(ctx, tx, order.AssetId)
			if err != nil {
				return err
			}
		}

		// 4 - Check KYC
		if refusal == "" {
			kycVerified, err := usrDb.CheckKYCisVerified(ctx, tx, order.UserId)
			if err != nil {
				return err
			}
			if !kycVerified {
				refusal = "the investor KYC is not verified or has expired"
			}
		}

		// 5 - A refusal found here happens before anything is signed status is set to DELIVERY_FAILED and return
		if refusal != "" {
			return moveOrder(ctx, tx, order,
				int(server.IssuanceOrderStatusIdDeliveryFailed), refusal, nil)
		}

		// 6 - Move the order to DELIVERY_IN_PROGRESS as the delivery is now starting.
		if err := moveOrder(ctx, tx, order,
			int(server.IssuanceOrderStatusIdDeliveryInProgress), "delivery started", nil); err != nil {
			return err
		}

		deliveryTarget = DeliveryTargetDTO{
			OrderID:      order.Id,
			UserID:       order.UserId,
			Shares:       order.Quantity,
			Wallet:       wallet,
			TokenAddress: *tokenAddress,
			RealEstateID: order.AssetId,
			TotalShares:  totalShares,
		}

		return nil
	})
	if err != nil {
		return DeliveryTargetDTO{}, err
	}

	if refusal != "" {
		return DeliveryTargetDTO{}, fmt.Errorf("%w: %s", ErrDeliveryRefused, refusal)
	}

	return deliveryTarget, nil
}

// moveOrder writes a status change and the history row that goes with it. The
// actor is NULL: the delivery is automatic, no human asked for it. txHash is
// the transaction the transition is about, nil when there is none.
func moveOrder(ctx context.Context, tx *sql.Tx, order locks.LockedIssuanceOrderDTO,
	toStatusID int, reason string, txHash *string) error {

	if err := ordrDb.UpdateIssuanceOrderStatus(ctx, tx, order.Id, toStatusID); err != nil {
		return err
	}

	return ordrDb.InsertIssuanceOrderHistory(ctx, tx, order.Id, order.StatusId, toStatusID, nil, reason, txHash)
}

// FailDelivery records a delivery that failed after transaction A, once the
// order was already DELIVERY_IN_PROGRESS and the chain had the last word. It
// opens its own transaction and takes its own lock.
//
// txHash is the transaction to put on the timeline, nil when nothing was ever
// signed. reachedNetwork says whether that transaction actually left for the
// node: when it did not, iss.issuance_orders.delivery_tx_hash is
// cleared, so the column keeps answering one question only — what was
// broadcast — and nothing later waits on a transaction the chain never saw.
//
// An order that is no longer DELIVERY_IN_PROGRESS yields
// ErrOrderNotDeliverable and nothing is written: someone else settled it.
func FailDelivery(ctx context.Context, orderRef, reason string, txHash *string, reachedNetwork bool) error {

	return utils.InTransaction(ctx, func(tx *sql.Tx) error {
		if err := locks.LockIssuanceOrder(ctx, tx, orderRef); err != nil {
			return err
		}

		order, err := locks.ReadLockedIssuanceOrder(ctx, tx, orderRef)
		if err != nil {
			return err
		}

		if order.StatusId != int(server.IssuanceOrderStatusIdDeliveryInProgress) {
			return ErrOrderNotDeliverable
		}

		if err := moveOrder(ctx, tx, order,
			int(server.IssuanceOrderStatusIdDeliveryFailed), reason, txHash); err != nil {
			return err
		}

		if !reachedNetwork {
			return clearDeliveryTxHash(ctx, tx, order.Id)
		}

		return nil
	})
}

// clearDeliveryTxHash takes back the hash recorded before sending, for the one
// case where the send was refused and the node never saw the transaction. The
// history row keeps it, so the attempt is not lost.
func clearDeliveryTxHash(ctx context.Context, tx *sql.Tx, orderID int64) error {
	const query = `
UPDATE iss.issuance_orders
SET delivery_tx_hash = NULL,
    updated_at = now()
WHERE id = $1
`

	if _, err := tx.ExecContext(ctx, query, orderID); err != nil {
		return fmt.Errorf("clear delivery tx hash of order %d: %w", orderID, err)
	}

	return nil
}

// readWalletAndTokenAddressForMint reads, in one statement, where the tokens
// must go and which contract must mint them. sql.ErrNoRows means the investor
// has no active wallet. A nil token address means the asset was never
// tokenised: ass.real_estate.contract_address is nullable.
func readWalletAndTokenAddressForMint(ctx context.Context, tx *sql.Tx, orderID int64) (string, *string, error) {
	const query = `
SELECT w.wallet_address, re.contract_address
FROM iss.issuance_orders o
JOIN blk.user_wallet w  ON w.user_id = o.user_id AND w.is_active
JOIN ass.real_estate re ON re.id = o.asset_id
WHERE o.id = $1
`

	var wallet string
	var tokenAddress sql.NullString
	err := tx.QueryRowContext(ctx, query, orderID).Scan(&wallet, &tokenAddress)
	if err != nil {
		return "", nil, err
	}

	if !tokenAddress.Valid {
		return wallet, nil, nil
	}

	return wallet, &tokenAddress.String, nil
}

// readTotalShares reads how many shares the asset offers in total.
func readTotalShares(ctx context.Context, tx *sql.Tx, realEstateID int) (int64, error) {
	const query = `SELECT total_shares::bigint FROM ass.real_estate_shares_config WHERE real_estate_id = $1`

	var total int64
	if err := tx.QueryRowContext(ctx, query, realEstateID).Scan(&total); err != nil {
		return 0, fmt.Errorf("read total shares of real estate %d: %w", realEstateID, err)
	}

	return total, nil
}

// RecordDeliveryTxHash is transaction B: it stores the hash of a mint that is
// signed but not sent yet.
//
// It is the gate that makes a double mint impossible. The WHERE clause accepts
// the write only while the order is still DELIVERY_IN_PROGRESS and carries no
// hash, so of two deliveries racing on the same order exactly one gets through.
// The other gets ErrDeliveryAlreadyRecorded and must not send anything.
//
// One statement, no lock: the UPDATE is atomic on its own, and holding a lock
// here would keep it across the network call that follows.
func RecordDeliveryTxHash(ctx context.Context, orderID int64, txHash string) error {
	const query = `
UPDATE iss.issuance_orders
SET delivery_tx_hash = $2,
    updated_at = now()
WHERE id = $1
  AND status_id = $3
  AND delivery_tx_hash IS NULL
`

	result, err := globals.DB.ExecContext(ctx, query, orderID, txHash,
		int(server.IssuanceOrderStatusIdDeliveryInProgress))
	if err != nil {
		return fmt.Errorf("record delivery tx hash of order %d: %w", orderID, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("record delivery tx hash of order %d: %w", orderID, err)
	}

	if affected == 0 {
		return ErrDeliveryAlreadyRecorded
	}

	return nil
}

// CompleteDelivery is transaction C: the mint is mined and verified, so the
// order becomes DELIVERED then CLOSED, both in the same transaction.
//
// It only accepts an order that is still DELIVERY_IN_PROGRESS and carries this
// very hash. Anything else means the order moved on without us, and nothing is
// written.
func CompleteDelivery(ctx context.Context, orderRef, txHash string) error {

	return utils.InTransaction(ctx, func(tx *sql.Tx) error {
		if err := locks.LockIssuanceOrder(ctx, tx, orderRef); err != nil {
			return err
		}

		order, err := locks.ReadLockedIssuanceOrder(ctx, tx, orderRef)
		if err != nil {
			return err
		}

		if order.StatusId != int(server.IssuanceOrderStatusIdDeliveryInProgress) ||
			!order.DeliveryTxHash.Valid || order.DeliveryTxHash.String != txHash {
			return ErrOrderNotDeliverable
		}

		delivered := int(server.IssuanceOrderStatusIdDelivered)
		if err := moveOrder(ctx, tx, order, delivered, "shares minted and verified on chain", &txHash); err != nil {
			return err
		}

		// The next transition starts where the previous one landed.
		order.StatusId = delivered

		return moveOrder(ctx, tx, order, int(server.IssuanceOrderStatusIdClosed),
			"nothing left to do on this order", &txHash)
	})
}
