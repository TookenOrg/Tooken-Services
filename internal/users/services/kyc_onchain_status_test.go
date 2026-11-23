package services

import (
	"context"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestReadKYCOnChainStatus(t *testing.T) {
	wallet := common.HexToAddress("0x0000000000000000000000000000000000000042")
	dbError := errors.New("wallet database unavailable")
	for _, tc := range []struct {
		name      string
		wallet    *common.Address
		walletErr error
		verified  bool
		chainErr  error
		wantErr   error
	}{
		{name: "no wallet does not contact the chain"},
		{name: "verified wallet", wallet: &wallet, verified: true},
		{name: "unverified wallet", wallet: &wallet},
		{name: "RPC unavailable cannot report verified", wallet: &wallet, verified: true, chainErr: errors.New("RPC failed"), wantErr: ErrKYCOnChainUnavailable},
		{name: "registry missing", wallet: &wallet, chainErr: errors.New("no registry"), wantErr: ErrKYCOnChainUnavailable},
		{name: "wallet lookup failure", walletErr: dbError, wantErr: dbError},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			reader := kycOnChainStatusReader{
				getWalletByUserID: func(ctx context.Context, userID int) (*common.Address, error) {
					if userID != 42 {
						t.Fatalf("userID = %d", userID)
					}
					return tc.wallet, tc.walletErr
				},
				fetchIsVerified: func(ctx context.Context, address string) (bool, error) {
					calls++
					if address != wallet.Hex() {
						t.Fatalf("address = %s", address)
					}
					if _, ok := ctx.Deadline(); !ok {
						t.Fatal("chain call must have a deadline")
					}
					return tc.verified, tc.chainErr
				},
			}
			status, err := reader.read(t.Context(), 42)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("err = %v, want %v", err, tc.wantErr)
			}
			if status.UserId != 42 || status.IsVerified != (tc.verified && tc.wantErr == nil) {
				t.Fatalf("status = %+v", status)
			}
			if tc.wallet == nil {
				if calls != 0 || status.WalletAddress != nil {
					t.Fatalf("missing wallet caused a chain call or address: %+v, calls=%d", status, calls)
				}
			} else if calls != 1 || status.WalletAddress == nil || *status.WalletAddress != wallet.Hex() {
				t.Fatalf("wallet status = %+v, calls=%d", status, calls)
			}
		})
	}
}
