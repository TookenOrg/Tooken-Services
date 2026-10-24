package services

import (
	"context"
	"errors"
)

// Mint_shares after payment confirmation.
func (s *Service) SendMintShares(ctx context.Context, tokenAddress, wallet string,
	shares int, record func(txHash string) error) (string, bool, error) {
	return "", false, errors.New("not implemented")
}

// MintReceipt checks the status of a mint transaction by its hash.
func (s *Service) MintReceipt(ctx context.Context, txHash, wallet string, shares int) (bool, bool, error) {
	return false, false, errors.New("not implemented")
}
