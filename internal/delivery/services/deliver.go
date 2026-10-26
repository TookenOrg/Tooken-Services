package services

import (
	"context"
	"errors"

	"github.com/TookenOrg/tooken-services/internal/delivery/database"
	"github.com/TookenOrg/tooken-services/pkg/logger"
)

// Deliver turns a paid order into tokens in the investor's wallet. Meant for a
// goroutine detached from the HTTP request: it reports through the order
// status and the logs, not through a return value.
//
// It serves one rule — an order is never minted twice — and most of its
// branches do nothing on purpose. An order is only marked failed when nothing
// can have been broadcast, or when the chain said no. In every other case,
// including "not mined yet", it is left alone: marking it failed is what would
// let a second mint start.
func (s *Service) Deliver(ctx context.Context, orderRef string) {

	logger.LogInfo("🚀 Starting delivery of order %s", orderRef)

	defer func() {
		if p := recover(); p != nil {
			_ = logger.LogError("Delivery of order %s panicked, the order is left as it is: %v", orderRef, p)
		}
	}()

	target, err := s.startDelivery(ctx, orderRef)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrOrderNotDeliverable):
			logger.LogInfo("🆗 Order %s is not awaiting delivery any more", orderRef)
		case errors.Is(err, database.ErrDeliveryRefused):
			// StartDelivery already recorded the failure.
			logger.LogWarn("⚠️ Delivery of order %s refused: %s", orderRef, err.Error())
		default:
			_ = logger.LogError("Failed to start delivery of order %s: %s", orderRef, err.Error())
		}
		return
	}

	txHash, sent, err := s.minter.SendMintShares(ctx, target.TokenAddress, target.Wallet, target.Shares,
		func(hash string) error {
			return s.recordTxHash(ctx, target.OrderID, hash)
		})

	if !sent {
		s.reportMintNotSent(ctx, orderRef, txHash, err)
		return
	}

	found, succeeded, err := s.minter.MintReceipt(ctx, txHash, target.Wallet, target.Shares)
	if err != nil {
		// The chain had the last word and we failed to hear it. The hash is
		// stored and the order stays in progress: nothing settles it yet, it
		// has to be checked on chain by hand.
		_ = logger.LogError("Failed to read the receipt of mint %s: %s", txHash, err.Error())
		return
	}

	if !found {
		logger.LogWarn("⏳ Mint %s of order %s is not mined yet: the order waits", txHash, orderRef)
		return
	}

	if !succeeded {
		// The transaction is on the chain, so its hash stays on the order.
		s.fail(ctx, orderRef, "the mint was mined but did not create the shares", &txHash, true)
		return
	}

	if err := s.completeDelivery(ctx, orderRef, txHash); err != nil {
		if errors.Is(err, database.ErrOrderNotDeliverable) {
			logger.LogInfo("🆗 Order %s was already settled by someone else", orderRef)
			return
		}
		_ = logger.LogError("Mint %s succeeded but order %s could not be closed: %s", txHash, orderRef, err.Error())
		return
	}

	logger.LogInfo("✅ Order %s delivered and closed, %d shares minted to %s", orderRef, target.Shares, target.Wallet)
}

// reportMintNotSent decides what to do when the mint never reached the network.
func (s *Service) reportMintNotSent(ctx context.Context, orderRef, txHash string, err error) {

	if errors.Is(err, database.ErrDeliveryAlreadyRecorded) {
		// Another delivery owns this mint: touching the order would fight it.
		logger.LogInfo("🆗 Order %s is already being minted by another delivery", orderRef)
		return
	}

	reason := "the mint could not be sent"
	if err != nil {
		reason = err.Error()
	}

	// Signed but never broadcast: the timeline keeps the attempt, the order
	// column does not.
	var hash *string
	if txHash != "" {
		hash = &txHash
	}

	s.fail(ctx, orderRef, reason, hash, false)
}

// fail records a delivery failure and says so, since nobody is waiting on a
// return value.
func (s *Service) fail(ctx context.Context, orderRef, reason string, txHash *string, reachedNetwork bool) {

	if err := s.failDelivery(ctx, orderRef, reason, txHash, reachedNetwork); err != nil {
		if errors.Is(err, database.ErrOrderNotDeliverable) {
			logger.LogInfo("🆗 Order %s was already settled by someone else", orderRef)
			return
		}
		_ = logger.LogError("Failed to mark order %s as failed: %s", orderRef, err.Error())
		return
	}

	logger.LogWarn("⚠️ Delivery of order %s failed: %s", orderRef, reason)
}
