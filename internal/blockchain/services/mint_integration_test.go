//go:build integration

package services

import (
	"context"
	"fmt"
)

// mintSharesForTest mints whole shares and waits for the chain, exactly the way
// the delivery does: sign, record, send, then read the receipt.
//
// It replaces the Service.Mint that was removed. The tests ask the same question
// as before — does this mint succeed, and does the balance move — but through
// the path the product actually takes, and without a float anywhere.
//
// The record callback does nothing here: storing the hash is the caller's job
// in production, and these tests have no order to store it against.
func mintSharesForTest(ctx context.Context, svc *Service, tokenAddr, wallet string, shares int) (string, error) {

	txHash, sent, err := svc.SendMintShares(ctx, tokenAddr, wallet, shares, func(string) error { return nil })
	if err != nil {
		return txHash, err
	}
	if !sent {
		return txHash, fmt.Errorf("mint of %d shares to %s never reached the network", shares, wallet)
	}

	found, succeeded, err := svc.MintReceipt(ctx, txHash, wallet, shares)
	if err != nil {
		return txHash, err
	}
	if !found {
		return txHash, fmt.Errorf("mint %s was not mined in time", txHash)
	}
	if !succeeded {
		return txHash, fmt.Errorf("mint %s was mined but did not create the shares", txHash)
	}

	return txHash, nil
}
