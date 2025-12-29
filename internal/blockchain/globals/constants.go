package globals

import (
	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var EthClient *ethclient.Client

// Implementation contract names
const (
	ImplClaimTopicRegistryName    = "ClaimTopicsRegistry"
	ImplTrustedIssuerRegistryName = "TrustedIssuerRegistry"
	ImplIdentityRegistryStorage   = "IdentityRegistryStorage"
	ImplIdentityRegistryName      = "IdentityRegistry"
	ImplModularComplianceName     = "ModularCompliance"
	ImplTokenName                 = "Token"
	ImplIdentityName              = "Identity"
	ImplIdentityAuthorityName     = "IdentityAuthority"
	ImplTrexAuthorityName         = "TREXImplementationAuthority"
)

const (
	IdentityFactoryName = "IdentityFactory"
	ClaimIssuerName     = "ClaimIssuer"
	TrexFactoryName     = "TrexFactory"
	TrexSuiteName       = "TREXSuite"
)

const (
	TransferRestrictionModuleName = "TransferRestrictionModule"
)

var TrexFactoryAddress common.Address
var ImplIdentityAuthorityAddress common.Address
var IdentityRegistryAddress common.Address
var IdentityRegistryInstance *contracts.IdentityRegistry
