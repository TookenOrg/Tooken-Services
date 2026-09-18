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

// registerIdentity records the wallet -> ONCHAINID link in the investor whitelist,
// through the IdentityRegistry of a deployed token.
//
// This is the path T-REX designed, and it needs no setup at all. When the factory
// deploys a suite it calls IdentityRegistryStorage.bindIdentityRegistry, which makes
// the new IdentityRegistry an agent of the storage; buildTokenDetails passes the
// platform in IrAgents, which makes the platform an agent of that IdentityRegistry.
// The authorisation chain therefore already exists end to end:
//
//	platform --agent of--> IdentityRegistry --agent of--> shared IdentityRegistryStorage
//
// Because every token created by the platform shares one storage, writing through any
// IdentityRegistry bound to it registers the investor for *all* of them — which is
// what lets a single KYC serve the whole platform.
func registerIdentity(ctx context.Context, userWallet, identityAddress common.Address, countryCode int) (tx *types.Transaction, err error) {

	// A country is stored on-chain as a uint16, and zero is a *valid* value there. A
	// silent truncation would register a wrong nationality, permanently.
	if countryCode < 0 || countryCode > math.MaxUint16 {
		err = logger.LogError("country code %d is out of range for an on-chain uint16", countryCode)
		return
	}

	irInstance, err := resolveIdentityRegistryInstance(ctx)
	if err != nil {
		return
	}

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Registration of new identity in the investor whitelist...")
	tx, err = irInstance.RegisterIdentity(auth, userWallet, identityAddress, uint16(countryCode))
	if err != nil {
		return
	}
	deployedTxDetails, err := utils.WaitDeployedTransaction(ctx, tx, false)
	if err != nil {
		return
	}

	logger.LogInfo("📬 Identity registred on transaction: %s", tx.Hash().Hex())

	err = database.InsertEthTransaction(ctx, deployedTxDetails.Tx.Hash().Hex(), "REGISTER_IDENTITY", deployedTxDetails.ToAddressHex(), deployedTxDetails.BlockNumber.Int64(), big.Int{})
	if err != nil {
		return
	}

	return
}

// resolveIdentityRegistryInstance picks an IdentityRegistry able to write into the
// shared investor whitelist, and asks the whitelist itself which ones those are.
//
// IdentityRegistryStorage.linkedIdentityRegistries() is the authoritative list: it is
// filled by bindIdentityRegistry, the very call that also makes each registry an agent
// of the storage. Reading it removes the need to guess — no token lookup, no heuristic
// "most recent one", and no check that the registry writes into the right ledger,
// since it was obtained *from* that ledger.
//
// Which registry is used does not matter: they all write the same row into the same
// storage, so one KYC serves every token of the platform. What matters is that the
// platform may write through it, and that is what the loop selects on.
func resolveIdentityRegistryInstance(ctx context.Context) (irInstance *contracts.IdentityRegistry, err error) {

	// resolveSharedIRS answers the zero address when nothing has been recorded yet.
	// The shared storage is created with the first token suite, so this means the
	// platform has no token at all.
	sharedIRS, err := resolveSharedIRS(ctx)
	if err != nil {
		return nil, err
	}
	if sharedIRS == (common.Address{}) {
		return nil, logger.LogError("no shared IdentityRegistryStorage recorded under %s: deploy a token suite before registering investors", globals.SharedIdentityRegistryStorageName)
	}

	irsInstance, err := contracts.NewIdentityRegistryStorage(sharedIRS, globals.EthClient)
	if err != nil {
		return nil, err
	}

	callOpts := &bind.CallOpts{Context: ctx}

	linkedRegistries, err := irsInstance.LinkedIdentityRegistries(callOpts)
	if err != nil {
		return nil, logger.LogError("could not list the registries bound to the shared whitelist %s: %s", sharedIRS.Hex(), err.Error())
	}
	if len(linkedRegistries) == 0 {
		return nil, logger.LogError("no IdentityRegistry is bound to the shared whitelist %s: deploy a token suite before registering investors", sharedIRS.Hex())
	}

	// registerIdentity is agent-restricted. buildTokenDetails puts the platform in
	// IrAgents of every suite it creates, so a usable registry is the normal case —
	// but a registry bound by someone else, or one the right was revoked on, must be
	// skipped rather than reverted through.
	platform := utils.GetEthFrom()
	for _, irAddress := range linkedRegistries {
		candidate, bindErr := contracts.NewIdentityRegistry(irAddress, globals.EthClient)
		if bindErr != nil {
			return nil, bindErr
		}

		isAgent, callErr := candidate.IsAgent(callOpts, platform)
		if callErr != nil {
			return nil, logger.LogError("could not check whether %s is an agent of the IdentityRegistry %s: %s", platform.Hex(), irAddress.Hex(), callErr.Error())
		}
		if isAgent {
			return candidate, nil
		}
	}

	return nil, logger.LogError("%s is an agent of none of the %d registries bound to the shared whitelist %s: no door to write the investor through", platform.Hex(), len(linkedRegistries), sharedIRS.Hex())
}
