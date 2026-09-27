package services

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	blkServices "github.com/TookenOrg/tooken-services/internal/blockchain/services"
	"github.com/TookenOrg/tooken-services/internal/users/database"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestKYCOnChainSyncBuildsPreciseResponseAndCallsChainInOrder(t *testing.T) {
	ctx := context.Background()
	wallet := common.HexToAddress("0x1111111111111111111111111111111111111111")
	identity := common.HexToAddress("0x2222222222222222222222222222222222222222")
	identityTx := "0xidentity"
	registerTx := "0xregister"
	claimTx := types.NewTransaction(7, identity, big.NewInt(0), 21_000, big.NewInt(1), nil)
	var calls []string

	syncSvc := testKYCOnChainSyncService()
	syncSvc.ensureWalletAndIdentity = func(context.Context, int) (common.Address, common.Address, *string, error) {
		calls = append(calls, "ensure")
		return wallet, identity, &identityTx, nil
	}
	syncSvc.addClaimToIdentity = func(_ context.Context, req server.AddClaimRequest) (*types.Transaction, error) {
		calls = append(calls, "claim")
		if req.UserId != 12 || req.ClaimTopic != kycClaimTopic {
			t.Fatalf("claim request = %+v", req)
		}
		return claimTx, nil
	}
	syncSvc.registerIdentityInSharedRegistry = func(_ context.Context, gotWallet, gotIdentity common.Address, country int) (*string, error) {
		calls = append(calls, "register")
		if gotWallet != wallet || gotIdentity != identity || country != 250 {
			t.Fatalf("register args = %s %s %d", gotWallet.Hex(), gotIdentity.Hex(), country)
		}
		return &registerTx, nil
	}
	syncSvc.fetchIsVerified = func(_ context.Context, gotWallet string) (bool, error) {
		calls = append(calls, "verify")
		if gotWallet != wallet.Hex() {
			t.Fatalf("verify wallet = %s, want %s", gotWallet, wallet.Hex())
		}
		return true, nil
	}

	response, verified, err := syncSvc.SyncApprovedKYCOnChain(ctx, server.KycVerification{
		Id:                  42,
		UserId:              12,
		DeclaredCountryCode: "FR",
	})
	if err != nil {
		t.Fatalf("SyncApprovedKYCOnChain error = %v", err)
	}
	if !verified {
		t.Fatal("SyncApprovedKYCOnChain verified = false, want true")
	}
	if !reflect.DeepEqual(calls, []string{"ensure", "claim", "register", "verify"}) {
		t.Fatalf("calls = %v", calls)
	}
	if response.UserId != 12 ||
		response.WalletAddress != wallet.Hex() ||
		response.IdentityAddress != identity.Hex() ||
		response.CountryCode != "FR" ||
		response.CountryNumeric != 250 {
		t.Fatalf("response identity fields = %+v", response)
	}
	if response.Transactions.Identity == nil || *response.Transactions.Identity != identityTx {
		t.Fatalf("identity tx = %v, want %s", response.Transactions.Identity, identityTx)
	}
	if response.Transactions.Claim == nil || *response.Transactions.Claim != claimTx.Hash().Hex() {
		t.Fatalf("claim tx = %v, want %s", response.Transactions.Claim, claimTx.Hash().Hex())
	}
	if response.Transactions.RegisterIdentity == nil || *response.Transactions.RegisterIdentity != registerTx {
		t.Fatalf("register tx = %v, want %s", response.Transactions.RegisterIdentity, registerTx)
	}
}

func TestKYCOnChainSyncRejectsInvalidCountryBeforeChainCalls(t *testing.T) {
	syncSvc := testKYCOnChainSyncService()
	called := false
	syncSvc.ensureWalletAndIdentity = func(context.Context, int) (common.Address, common.Address, *string, error) {
		called = true
		return common.Address{}, common.Address{}, nil, nil
	}

	_, _, err := syncSvc.SyncApprovedKYCOnChain(context.Background(), server.KycVerification{
		UserId:              12,
		DeclaredCountryCode: "FRA",
	})
	if !errors.Is(err, ErrMessageInvalidCountryCode) {
		t.Fatalf("err = %v, want ErrMessageInvalidCountryCode", err)
	}
	if called {
		t.Fatal("invalid country must fail before any blockchain call")
	}
}

func TestKYCOnChainSyncMapsIdentityWithoutWallet(t *testing.T) {
	syncSvc := testKYCOnChainSyncService()
	syncSvc.ensureWalletAndIdentity = func(context.Context, int) (common.Address, common.Address, *string, error) {
		return common.Address{}, common.Address{}, nil, blkServices.ErrIdentityWithoutWallet
	}

	_, _, err := syncSvc.SyncApprovedKYCOnChain(context.Background(), server.KycVerification{
		UserId:              12,
		DeclaredCountryCode: "FR",
	})
	if !errors.Is(err, ErrMessageIncoherentUserState) {
		t.Fatalf("err = %v, want ErrMessageIncoherentUserState", err)
	}
}

func TestKYCOnChainSyncAlreadyVerifiedReturnsEarlyWhenChainConfirms(t *testing.T) {
	ctx := context.Background()
	wallet := common.HexToAddress("0x1111111111111111111111111111111111111111")
	identity := common.HexToAddress("0x2222222222222222222222222222222222222222")
	status := database.StatusVerified
	var calls []string

	syncSvc := testKYCOnChainSyncService()
	syncSvc.ensureWalletAndIdentity = func(context.Context, int) (common.Address, common.Address, *string, error) {
		calls = append(calls, "ensure")
		return wallet, identity, nil, nil
	}
	syncSvc.getWalletByUserID = func(context.Context, int) (*common.Address, error) {
		calls = append(calls, "get-wallet")
		return &wallet, nil
	}
	syncSvc.fetchIsVerified = func(context.Context, string) (bool, error) {
		calls = append(calls, "verify")
		return true, nil
	}

	response, verified, err := syncSvc.SyncApprovedKYCOnChain(ctx, server.KycVerification{
		Id:                  42,
		UserId:              12,
		DeclaredCountryCode: "FR",
		KycStatus:           &status,
	})
	if err != nil {
		t.Fatalf("SyncApprovedKYCOnChain error = %v", err)
	}
	if !verified {
		t.Fatal("verified = false, want true")
	}
	if !reflect.DeepEqual(calls, []string{"ensure", "get-wallet", "verify"}) {
		t.Fatalf("calls = %v", calls)
	}
	if response.Transactions.Claim != nil || response.Transactions.RegisterIdentity != nil {
		t.Fatalf("idempotent response must not invent tx hashes: %+v", response.Transactions)
	}
}

func TestKYCOnChainSyncRepairsWhenDatabaseVerifiedButChainIsNot(t *testing.T) {
	ctx := context.Background()
	wallet := common.HexToAddress("0x1111111111111111111111111111111111111111")
	identity := common.HexToAddress("0x2222222222222222222222222222222222222222")
	status := database.StatusVerified
	var calls []string

	syncSvc := testKYCOnChainSyncService()
	syncSvc.ensureWalletAndIdentity = func(context.Context, int) (common.Address, common.Address, *string, error) {
		calls = append(calls, "ensure")
		return wallet, identity, nil, nil
	}
	syncSvc.getWalletByUserID = func(context.Context, int) (*common.Address, error) {
		calls = append(calls, "get-wallet")
		return &wallet, nil
	}
	syncSvc.fetchIsVerified = func(context.Context, string) (bool, error) {
		calls = append(calls, "verify")
		return len(calls) > 4, nil
	}
	syncSvc.addClaimToIdentity = func(context.Context, server.AddClaimRequest) (*types.Transaction, error) {
		calls = append(calls, "claim")
		return types.NewTransaction(1, identity, big.NewInt(0), 21_000, big.NewInt(1), nil), nil
	}
	syncSvc.registerIdentityInSharedRegistry = func(context.Context, common.Address, common.Address, int) (*string, error) {
		calls = append(calls, "register")
		return nil, nil
	}

	_, verified, err := syncSvc.SyncApprovedKYCOnChain(ctx, server.KycVerification{
		Id:                  42,
		UserId:              12,
		DeclaredCountryCode: "FR",
		KycStatus:           &status,
	})
	if err != nil {
		t.Fatalf("SyncApprovedKYCOnChain error = %v", err)
	}
	if !verified {
		t.Fatal("verified = false, want true after repair")
	}
	if !reflect.DeepEqual(calls, []string{"ensure", "get-wallet", "verify", "claim", "register", "verify"}) {
		t.Fatalf("calls = %v", calls)
	}
}

func TestKYCOnChainSyncStopsOnClaimErrorBeforeRegistering(t *testing.T) {
	ctx := context.Background()
	wallet := common.HexToAddress("0x1111111111111111111111111111111111111111")
	identity := common.HexToAddress("0x2222222222222222222222222222222222222222")
	claimErr := errors.New("claim failed")
	registered := false

	syncSvc := testKYCOnChainSyncService()
	syncSvc.ensureWalletAndIdentity = func(context.Context, int) (common.Address, common.Address, *string, error) {
		return wallet, identity, nil, nil
	}
	syncSvc.addClaimToIdentity = func(context.Context, server.AddClaimRequest) (*types.Transaction, error) {
		return nil, claimErr
	}
	syncSvc.registerIdentityInSharedRegistry = func(context.Context, common.Address, common.Address, int) (*string, error) {
		registered = true
		return nil, nil
	}

	_, _, err := syncSvc.SyncApprovedKYCOnChain(ctx, server.KycVerification{
		UserId:              12,
		DeclaredCountryCode: "FR",
	})
	if !errors.Is(err, claimErr) {
		t.Fatalf("err = %v, want claimErr", err)
	}
	if registered {
		t.Fatal("registerIdentity must not run after AddClaimToIdentity fails")
	}
}

func TestKYCOnChainSyncReturnsFalseWhenFinalVerificationFails(t *testing.T) {
	ctx := context.Background()
	wallet := common.HexToAddress("0x1111111111111111111111111111111111111111")
	identity := common.HexToAddress("0x2222222222222222222222222222222222222222")

	syncSvc := testKYCOnChainSyncService()
	syncSvc.ensureWalletAndIdentity = func(context.Context, int) (common.Address, common.Address, *string, error) {
		return wallet, identity, nil, nil
	}
	syncSvc.addClaimToIdentity = func(context.Context, server.AddClaimRequest) (*types.Transaction, error) {
		return types.NewTransaction(1, identity, big.NewInt(0), 21_000, big.NewInt(1), nil), nil
	}
	syncSvc.registerIdentityInSharedRegistry = func(context.Context, common.Address, common.Address, int) (*string, error) {
		return nil, nil
	}
	syncSvc.fetchIsVerified = func(context.Context, string) (bool, error) {
		return false, nil
	}

	response, verified, err := syncSvc.SyncApprovedKYCOnChain(ctx, server.KycVerification{
		UserId:              12,
		DeclaredCountryCode: "FR",
	})
	if err != nil {
		t.Fatalf("SyncApprovedKYCOnChain error = %v", err)
	}
	if verified {
		t.Fatal("verified = true, want false")
	}
	if response.WalletAddress != wallet.Hex() || response.IdentityAddress != identity.Hex() {
		t.Fatalf("response should still explain what happened: %+v", response)
	}
}

func testKYCOnChainSyncService() *kycOnChainSyncService {
	t := func(context.Context, server.AddClaimRequest) (*types.Transaction, error) {
		panic("unexpected AddClaimToIdentity call")
	}
	return &kycOnChainSyncService{
		addClaimToIdentity: t,
		ensureWalletAndIdentity: func(context.Context, int) (common.Address, common.Address, *string, error) {
			panic("unexpected EnsureWalletAndIdentity call")
		},
		fetchIsVerified: func(context.Context, string) (bool, error) {
			panic("unexpected FetchIsVerified call")
		},
		getWalletByUserID: func(context.Context, int) (*common.Address, error) {
			panic("unexpected GetWalletByUserID call")
		},
		registerIdentityInSharedRegistry: func(context.Context, common.Address, common.Address, int) (*string, error) {
			panic("unexpected RegisterIdentityInSharedRegistry call")
		},
	}
}
