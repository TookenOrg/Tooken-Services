package services

import "context"

type Service struct {
	minter TokenMinter
}

func NewService(minter TokenMinter) *Service {
	return &Service{minter: minter}
}

// Deliver runs the automatic delivery of a paid order. Meant for a goroutine:
// it reports through the order status and the logs, not through its caller.
func (s *Service) Deliver(ctx context.Context, orderRef string) {
	// Implementation of delivery logic goes here
}

// TokenMinter is what delivery needs from the blockchain module.
type TokenMinter interface {
	// SendMintShares signs the mint, calls record with its hash, and sends it only
	// if record returned nil. sent=false means nothing reached the network.
	SendMintShares(ctx context.Context, tokenAddress, wallet string, shares int,
		record func(txHash string) error) (txHash string, sent bool, err error)

	// MintReceipt waits until the mint is mined or the wait times out (step 7).
	// succeeded: status 1 and Transfer(0x0 → wallet, amount) in its logs (L5).
	// found=false: still not mined when the wait timed out — not a failure (§4.5).
	MintReceipt(ctx context.Context, txHash, wallet string, shares int) (
		found bool, succeeded bool, err error)
}
