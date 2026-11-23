package services

import (
	"context"

	assetsDb "github.com/TookenOrg/tooken-services/internal/assets_managements/database"
	"github.com/TookenOrg/tooken-services/internal/delivery/database"
)

// TokenMinter is what delivery needs from the blockchain module.
type TokenMinter interface {
	// SendMintShares signs the mint, calls record with its hash, and sends it
	// only if record returned nil. sent=false means nothing reached the network.
	SendMintShares(ctx context.Context, tokenAddress, wallet string, shares int,
		record func(txHash string) error) (txHash string, sent bool, err error)

	// MintReceipt waits for the mint. succeeded needs status 1 and a
	// Transfer(0x0 → wallet, amount) in the logs. found=false means still not
	// mined — undecided, not failed.
	MintReceipt(ctx context.Context, txHash, wallet string, shares int) (
		found bool, succeeded bool, err error)

	// CompareTotalSupply reads totalSupply() of the token at the latest block
	// and compares it to shares whole shares in base units. It returns -1, 0 or
	// +1 like big.Int.Cmp, all in integers.
	CompareTotalSupply(ctx context.Context, tokenAddress string, shares int64) (int, error)
}

// The database work Deliver relies on, as function types.
type (
	startDeliveryFn    func(ctx context.Context, orderRef string) (database.DeliveryTargetDTO, error)
	recordTxHashFn     func(ctx context.Context, orderID int64, txHash string) error
	failDeliveryFn     func(ctx context.Context, orderRef, reason string, txHash *string, reachedNetwork bool) error
	completeDeliveryFn func(ctx context.Context, orderRef, txHash string) error
	markFundedFn       func(ctx context.Context, realEstateID int) error
)

type Service struct {
	minter TokenMinter

	// Held as fields rather than called directly, so the orchestration — where
	// the no-double-mint rule lives — can be tested without a database.
	startDelivery    startDeliveryFn
	recordTxHash     recordTxHashFn
	failDelivery     failDeliveryFn
	completeDelivery completeDeliveryFn
	markFunded       markFundedFn
}

func NewService(minter TokenMinter) *Service {
	return &Service{
		minter:           minter,
		startDelivery:    database.StartDelivery,
		recordTxHash:     database.RecordDeliveryTxHash,
		failDelivery:     database.FailDelivery,
		completeDelivery: database.CompleteDelivery,
		markFunded:       assetsDb.MarkRealEstateFunded,
	}
}
