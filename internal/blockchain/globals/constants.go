package globals

import (
	"database/sql"

	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var DB *sql.DB
var EthClient *ethclient.Client

// Implementation contract names
const (
	ImplClaimTopicRegistryName     = "ClaimTopicsRegistryImplementation"
	ImplTrustedIssuersRegistryName = "TrustedIssuersRegistryImplementation"
	ImplIdentityRegistryStorage    = "IdentityRegistryStorageImplementation"
	ImplIdentityRegistryName       = "IdentityRegistryImplementation"
	ImplModularComplianceName      = "ModularComplianceImplementation"
	ImplTokenName                  = "TokenImplementation"
	ImplIdentityName               = "IdentityImplementation"
	ImplIdentityAuthorityName      = "IdentityAuthorityImplementation"
	ImplTrexAuthorityName          = "TREXImplementationAuthority"
	ImplClaimIssuerName            = "ClaimIssuerImplementation"
	ImplIdentityFactoryName        = "IdentityFactoryImplementation"
	ImplTrexSuiteName              = "TREXSuite"
)

const (
	TransferRestrictionModuleName = "TransferRestrictionModule"
)

var TrexFactoryAddress common.Address
var ImplIdentityAuthorityAddress common.Address
var ClaimIssuerAddress common.Address
var IdentityRegistryAddress common.Address

var IdentityRegistryInstance *contracts.IdentityRegistry
