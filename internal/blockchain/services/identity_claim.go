package services

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/database"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/models"
	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	claimData = "OK"
)

func (s *Service) AddClaimToIdentity(ctx context.Context, addClaimRequest server.AddClaimRequest) (tx *types.Transaction, err error) {

	// 1 - Get user's Identity contract
	userIdentityAddrPtr, err := database.GetIdentityAddrByUserId(ctx, addClaimRequest.UserId)
	if err != nil {
		return nil, err
	}
	if userIdentityAddrPtr == nil {
		return nil, fmt.Errorf("user Identity not found")
	}
	userIdentityAddr := *userIdentityAddrPtr

	// 2 - Generate Signature
	signaturePtr, err := generateSignatureAddClaim(userIdentityAddr, int64(addClaimRequest.ClaimTopic))
	if err != nil {
		return nil, err
	}
	if signaturePtr == nil {
		return nil, fmt.Errorf("generated signature is nil")
	}
	signature := *signaturePtr

	// 3 - Create Identity Instance
	identityInstance, err := getIdentityInstance(userIdentityAddr)
	if err != nil {
		return nil, err
	}

	// 4 - Add Claim
	tx, err = addClaim(ctx, identityInstance, addClaimRequest, signature)
	if err != nil {
		return nil, err
	}

	// 5 - DB

	return
}

func generateSignatureAddClaim(identityToClaim common.Address, claimTopic int64) (sign *models.SignatureResult, err error) {
	// Get private key from environment
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		return nil, fmt.Errorf("PRIVATE_KEY environment variable not set")
	}

	// Remove 0x prefix if present
	if len(privateKeyHex) >= 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Get wallet address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to cast public key to ECDSA")
	}
	walletAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Convert data to bytes (equivalent to ethers.utils.toUtf8Bytes)
	dataBytes := []byte(claimData)

	// Create ABI encoder types
	addressType, _ := abi.NewType("address", "", nil)
	uint256Type, _ := abi.NewType("uint256", "", nil)
	bytesType, _ := abi.NewType("bytes", "", nil)

	// Create arguments for encoding
	arguments := abi.Arguments{
		{Type: addressType},
		{Type: uint256Type},
		{Type: bytesType},
	}

	// Encode data (equivalent to ethers.utils.defaultAbiCoder.encode)
	topicBig := big.NewInt(claimTopic)
	encodedData, err := arguments.Pack(
		identityToClaim,
		topicBig,
		dataBytes,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode data: %w", err)
	}

	// Calculate data hash (equivalent to ethers.utils.keccak256)
	dataHash := crypto.Keccak256Hash(encodedData)

	// Create the Ethereum message hash (equivalent to ethers.utils.hashMessage)
	message := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(dataHash.Bytes()))
	ethSignedHash := crypto.Keccak256Hash([]byte(message), dataHash.Bytes())

	// Sign the ethSignedHash (equivalent to wallet.signMessage)
	signature, err := crypto.Sign(ethSignedHash.Bytes(), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}

	// Extract v, r, s from signature
	v := signature[64]
	r := common.BytesToHash(signature[0:32])
	s := common.BytesToHash(signature[32:64])

	// Ethereum signatures have recovery ID adjustment
	if v < 27 {
		v += 27
	}

	// Recover address for verification (use ethSignedHash for verification)
	recoveredPubKey, err := crypto.SigToPub(ethSignedHash.Bytes(), signature)
	if err != nil {
		return nil, fmt.Errorf("failed to recover public key: %w", err)
	}
	recoveredAddr := crypto.PubkeyToAddress(*recoveredPubKey)

	// Adjust the signature for Ethereum compatibility (after recovery)
	signature[64] = v

	// Calculate key (equivalent to ethers.utils.keccak256(ethers.utils.zeroPad(...)))
	paddedAddr := common.LeftPadBytes(recoveredAddr.Bytes(), 32)
	key := crypto.Keccak256Hash(paddedAddr)

	result := &models.SignatureResult{
		Wallet:          walletAddress.Hex(),
		IdentityToClaim: identityToClaim.Hex(),
		Topic:           topicBig,
		DataBytes:       dataBytes,
		DataHash:        dataHash,
		EthSignedHash:   ethSignedHash,
		Signature:       "0x" + common.Bytes2Hex(signature),
		SignatureBytes:  signature,
		V:               v,
		R:               r,
		S:               s,
		RecoveredAddr:   recoveredAddr,
		Key:             key,
	}

	logger.LogTrace("🔒 Wallet: %s", result.Wallet)
	logger.LogTrace("👤 Identity to claim: %s", result.IdentityToClaim)
	logger.LogTrace("📛 Topic: %s", result.Topic.String())
	logger.LogTrace("📦 Data (bytes): %s", common.Bytes2Hex(result.DataBytes))
	logger.LogTrace("🧩 DataHash: %s", result.DataHash.Hex())
	logger.LogTrace("🔁 EthSignedHash: %s", result.EthSignedHash.Hex())
	logger.LogTrace("✍️ Signature: %s", result.Signature)
	logger.LogTrace("↪️ v: %d | r: %s | s: %s", result.V, result.R.Hex(), result.S.Hex())
	logger.LogTrace("✅ Signer address (recovered): %s", result.RecoveredAddr.Hex())
	logger.LogTrace("🗝️ getKey(...) value: %s", result.Key.Hex())

	return result, nil
}

func addClaim(ctx context.Context, identityInstance *contracts.Identity, addClaimRequest server.AddClaimRequest, signature models.SignatureResult) (tx *types.Transaction, err error) {

	// 1 - Get issuer address
	issuerAddressDetails, err := database.GetContractRoleByName(ctx, globals.ClaimIssuerName)
	if err != nil {
		return
	}

	issuer := common.HexToAddress(issuerAddressDetails.Address)

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Adding claim on Identity...")

	claimTopic := big.NewInt(int64(addClaimRequest.ClaimTopic))
	scheme := big.NewInt(1)
	signatureBytes := signature.SignatureBytes
	dataBytes := []byte(claimData)

	tx, err = identityInstance.AddClaim(auth, claimTopic, scheme, issuer, signatureBytes, dataBytes, "")
	if err != nil {
		return
	}

	deployedTxDetails, err := utils.WaitDeployedTransaction(ctx, tx, false)
	if err != nil {
		return
	}

	logger.LogInfo("📬 Claim added on Identity with transaction: %s", tx.Hash().Hex())

	_, err = database.InsertEthTransaction(ctx, deployedTxDetails.Tx.Hash().Hex(), "ADD_CLAIMS_INDENTITY", deployedTxDetails.Tx.To().Hex(), deployedTxDetails.BlockNumber.Int64(), big.Int{})
	if err != nil {
		return
	}

	return
}
