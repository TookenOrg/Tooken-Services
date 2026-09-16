package services

import (
	"context"
	"math"
	"math/big"

	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/database"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// registerIdentity records the wallet -> ONCHAINID link in the shared
// IdentityRegistryStorage: the single investor whitelist that every token reads
// through its own IdentityRegistry.
//
// It writes into the storage directly rather than calling
// IdentityRegistry.registerIdentity, because an IdentityRegistry is a per-token facade
// while a KYC is a platform-level fact. Going through one would make onboarding an
// investor depend on a property being tokenised first, which the product does not
// require. registerIdentity is itself a thin wrapper around addIdentityToStorage, so
// both write the same data and isVerified reads it either way — only the emitted event
// differs (IdentityStored instead of IdentityRegistered).
func registerIdentity(ctx context.Context, userWallet, identityAddress common.Address, countryCode int) (tx *types.Transaction, err error) {

	// A country is stored on-chain as a uint16, and zero is a *valid* value there. A
	// silent truncation would register a wrong nationality, permanently.
	if countryCode < 0 || countryCode > math.MaxUint16 {
		err = logger.LogError("country code %d is out of range for an on-chain uint16", countryCode)
		return
	}

	irsInstance, err := resolveSharedIRSInstance(ctx)
	if err != nil {
		return
	}

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Registration of new identity in the shared investor whitelist...")
	tx, err = irsInstance.AddIdentityToStorage(auth, userWallet, identityAddress, uint16(countryCode))
	if err != nil {
		return
	}
	deployedTxDetails, err := utils.WaitDeployedTransaction(ctx, tx, false)
	if err != nil {
		return
	}

	logger.LogInfo("📬 Identity registred on transaction: %s", tx.Hash().Hex())

	err = database.InsertEthTransaction(ctx, deployedTxDetails.Tx.Hash().Hex(), "REGISTER_IDENTITY", deployedTxDetails.Tx.To().Hex(), deployedTxDetails.BlockNumber.Int64(), big.Int{})
	if err != nil {
		return
	}

	return
}

// resolveSharedIRSInstance binds the shared IdentityRegistryStorage and checks the
// platform is allowed to write into it.
//
// The IRS address is a singleton persisted under SHARED_IRS in blk.contract_role —
// unlike the per-token IdentityRegistry addresses, which are read from the chain.
func resolveSharedIRSInstance(ctx context.Context) (irsInstance *contracts.IdentityRegistryStorage, err error) {

	// resolveSharedIRS answers with the zero address when no IRS has been recorded yet.
	// That is a legitimate answer when deploying the first token, never here: binding
	// the zero address would let the registration fail far away from its cause.
	irsAddress, err := resolveSharedIRS(ctx)
	if err != nil {
		return nil, err
	}
	if irsAddress == (common.Address{}) {
		return nil, logger.LogError("no shared IdentityRegistryStorage recorded yet: deploy a token suite before registering investors")
	}

	irsInstance, err = contracts.NewIdentityRegistryStorage(irsAddress, globals.EthClient)
	if err != nil {
		return nil, err
	}

	// addIdentityToStorage is agent-restricted. Asking first turns an opaque revert
	// into a message that names the fix.
	platform := utils.GetEthFrom()
	isAgent, err := irsInstance.IsAgent(&bind.CallOpts{Context: ctx}, platform)
	if err != nil {
		return nil, logger.LogError("could not check whether %s is an agent of the shared IRS %s: %s", platform.Hex(), irsAddress.Hex(), err.Error())
	}
	if !isAgent {
		return nil, logger.LogError("%s is not an agent of the shared IRS %s: run the IRS agent setup first", platform.Hex(), irsAddress.Hex())
	}

	return irsInstance, nil
}
