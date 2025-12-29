package globals

import (
	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var EthClient *ethclient.Client

// Implementation contract names
const (
	ImplClaimsTopicRegistryName   = "CLAIMS_TOPIC_REGISTRY"
	ImplTrustedIssuerRegistryName = "TRUSTED_ISSUER_REGISTRY"
	ImplIdentityRegistryStorage   = "IDENTITY_REGISTRY_STORAGE"
	ImplIdentityRegistryName      = "IDENTITY_REGISTRY"
	ImplModularComplianceName     = "MODULAR_COMPLIANCE"
	ImplTokenName                 = "TOKEN"
	ImplIdentityName              = "IDENTITY"
	ImplIdentityAuthorityName     = "IDENTITY_AUTHORITY"
	ImplTrexAuthorityName         = "TREX_AUTHORITY"
)

const (
	IdentityFactoryName = "IDENTITY_FACTORY"
	IdentityRegistry    = "IDENTITY_REGISTRY"
	ClaimIssuerName     = "CLAIM_ISSUER"
	TrexFactoryName     = "TREX_FACTORY"
	TrexSuiteName       = "TREX_SUITE"
)

const (
	TransferRestrictionModuleName = "TRANSFER_RESTRICTION_MODULE"
)

var TrexFactoryAddress common.Address
var ImplIdentityAuthorityAddress common.Address
var IdentityRegistryAddress common.Address
var IdentityRegistryInstance *contracts.IdentityRegistry
