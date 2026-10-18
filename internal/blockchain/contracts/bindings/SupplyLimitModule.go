// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// SupplyLimitModuleMetaData contains all meta data concerning the SupplyLimitModule contract.
var SupplyLimitModuleMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"ComplianceBound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"ComplianceUnbound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_limit\",\"type\":\"uint256\"}],\"name\":\"SupplyLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"bindCompliance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"canComplianceBind\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"getSupplyLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"isComplianceBound\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isPlugAndPlay\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"moduleBurnAction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"moduleCheck\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"moduleMintAction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"moduleTransferAction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_limit\",\"type\":\"uint256\"}],\"name\":\"setSupplyLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"unbindCompliance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523060805234801561001457600080fd5b5060805161195161004c60003960008181610825015281816108aa01528181610bb501528181610c3a0152610d2401526119516000f3fe6080604052600436106101445760003560e01c80634cf4d295116100c05780638da5cb5b11610074578063e6f5e80711610059578063e6f5e807146103c3578063f104a8c91461024c578063f2fde38b146103d757600080fd5b80638da5cb5b1461037a578063bcc21053146103a257600080fd5b806352d1902d116100a557806352d1902d1461033b578063715018a6146103505780638129fc1c1461036557600080fd5b80634cf4d295146102d05780634f1ef2861461032857600080fd5b8063361fab2511610117578063372491a2116100fc578063372491a21461024c5780634a9325441461026c5780634ae65b5c1461028c57600080fd5b8063361fab251461020c5780633659cfe61461022c57600080fd5b8063013b7ce4146101495780630694a5fb1461017e57806306fdde03146101a05780632cb7e1ec146101ec575b600080fd5b34801561015557600080fd5b50610169610164366004611653565b6103f7565b60405190151581526020015b60405180910390f35b34801561018a57600080fd5b5061019e6101993660046116a6565b61051a565b005b3480156101ac57600080fd5b50604080518082018252601181527f537570706c794c696d69744d6f64756c650000000000000000000000000000006020820152905161017591906116ee565b3480156101f857600080fd5b5061019e610207366004611721565b6106c4565b34801561021857600080fd5b5061019e610227366004611762565b61074b565b34801561023857600080fd5b5061019e6102473660046116a6565b61081b565b34801561025857600080fd5b5061019e61026736600461177b565b610996565b34801561027857600080fd5b5061019e6102873660046116a6565b610a1c565b34801561029857600080fd5b506102c26102a73660046116a6565b6001600160a01b0316600090815260c9602052604090205490565b604051908152602001610175565b3480156102dc57600080fd5b506101696102eb3660046116a6565b6001600160a01b031660009081527ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae00602052604090205460ff1690565b61019e6103363660046117bd565b610bab565b34801561034757600080fd5b506102c2610d17565b34801561035c57600080fd5b5061019e610ddc565b34801561037157600080fd5b5061019e610df0565b34801561038657600080fd5b506033546040516001600160a01b039091168152602001610175565b3480156103ae57600080fd5b506101696103bd3660046116a6565b50600190565b3480156103cf57600080fd5b506001610169565b3480156103e357600080fd5b5061019e6103f23660046116a6565b610f0f565b60006001600160a01b038516158015610501575060c96000836001600160a01b03166001600160a01b031681526020019081526020016000205483836001600160a01b0316636a3edf286040518163ffffffff1660e01b8152600401602060405180830381865afa158015610470573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104949190611881565b6001600160a01b03166318160ddd6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156104d1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104f5919061189e565b6104ff91906118b7565b115b1561050e57506000610512565b5060015b949350505050565b3360009081527ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae00602081905260409091205460ff166105a05760405162461bcd60e51b815260206004820152601e60248201527f6f6e6c7920626f756e6420636f6d706c69616e63652063616e2063616c6c000060448201526064015b60405180910390fd5b7ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae006001600160a01b0383166106175760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f2061646472657373006044820152606401610597565b336001600160a01b038416146106795760405162461bcd60e51b815260206004820152602160248201527f6f6e6c7920636f6d706c69616e636520636f6e74726163742063616e2063616c6044820152601b60fa1b6064820152608401610597565b6001600160a01b038316600081815260208390526040808220805460ff19169055517f408b49d9be1c914c52a0227e18a077e5a892dddf32a26cfa94a5d9708fad77189190a2505050565b3360009081527ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae00602081905260409091205460ff166107455760405162461bcd60e51b815260206004820152601e60248201527f6f6e6c7920626f756e6420636f6d706c69616e63652063616e2063616c6c00006044820152606401610597565b50505050565b3360009081527ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae00602081905260409091205460ff166107cc5760405162461bcd60e51b815260206004820152601e60248201527f6f6e6c7920626f756e6420636f6d706c69616e63652063616e2063616c6c00006044820152606401610597565b33600081815260c96020908152604091829020859055815192835282018490527f2f1ac7412e0c36ff7b73f7c4f474022dfa5d3f1c09e3262f63c9eb4627a9f505910160405180910390a15050565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036108a85760405162461bcd60e51b815260206004820152602c60248201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060448201526b19195b1959d85d1958d85b1b60a21b6064820152608401610597565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166109037f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b6001600160a01b03161461096e5760405162461bcd60e51b815260206004820152602c60248201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060448201526b6163746976652070726f787960a01b6064820152608401610597565b61097781610f9c565b6040805160008082526020820190925261099391839190610fa4565b50565b3360009081527ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae00602081905260409091205460ff16610a175760405162461bcd60e51b815260206004820152601e60248201527f6f6e6c7920626f756e6420636f6d706c69616e63652063616e2063616c6c00006044820152606401610597565b505050565b7ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae006001600160a01b038216610a935760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f2061646472657373006044820152606401610597565b6001600160a01b03821660009081526020829052604090205460ff1615610afc5760405162461bcd60e51b815260206004820152601860248201527f636f6d706c69616e636520616c726561647920626f756e6400000000000000006044820152606401610597565b336001600160a01b03831614610b5e5760405162461bcd60e51b815260206004820152602160248201527f6f6e6c7920636f6d706c69616e636520636f6e74726163742063616e2063616c6044820152601b60fa1b6064820152608401610597565b6001600160a01b038216600081815260208390526040808220805460ff19166001179055517f1f7b76c58fb697eb53c6c7c1becb96911516a136e24d7ced386b2355358b75a39190a25050565b6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000163003610c385760405162461bcd60e51b815260206004820152602c60248201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060448201526b19195b1959d85d1958d85b1b60a21b6064820152608401610597565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610c937f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b6001600160a01b031614610cfe5760405162461bcd60e51b815260206004820152602c60248201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060448201526b6163746976652070726f787960a01b6064820152608401610597565b610d0782610f9c565b610d1382826001610fa4565b5050565b6000306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614610db75760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c00000000000000006064820152608401610597565b507f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc90565b610de4611144565b610dee600061119e565b565b600054610100900460ff1615808015610e105750600054600160ff909116105b80610e2a5750303b158015610e2a575060005460ff166001145b610e9c5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a65640000000000000000000000000000000000006064820152608401610597565b6000805460ff191660011790558015610ebf576000805461ff0019166101001790555b610ec76111fd565b8015610993576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a150565b610f17611144565b6001600160a01b038116610f935760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610597565b6109938161119e565b610993611144565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff1615610fd757610a1783611278565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015611031575060408051601f3d908101601f1916820190925261102e9181019061189e565b60015b6110a35760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201527f6f6e206973206e6f7420555550530000000000000000000000000000000000006064820152608401610597565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc81146111385760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f7860448201527f6961626c655555494400000000000000000000000000000000000000000000006064820152608401610597565b50610a17838383611343565b6033546001600160a01b03163314610dee5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610597565b603380546001600160a01b0383811673ffffffffffffffffffffffffffffffffffffffff19831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166112685760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b6064820152608401610597565b611270611368565b610dee6113db565b6001600160a01b0381163b6112f55760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201527f6f74206120636f6e7472616374000000000000000000000000000000000000006064820152608401610597565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0392909216919091179055565b61134c83611446565b6000825111806113595750805b15610a17576107458383611486565b600054610100900460ff166113d35760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b6064820152608401610597565b610dee6114b4565b600054610100900460ff16610dee5760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b6064820152608401610597565b61144f81611278565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b60606114ab83836040518060600160405280602781526020016118f560279139611528565b90505b92915050565b600054610100900460ff1661151f5760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b6064820152608401610597565b610dee3361119e565b6060600080856001600160a01b03168560405161154591906118d8565b600060405180830381855af49150503d8060008114611580576040519150601f19603f3d011682016040523d82523d6000602084013e611585565b606091505b5091509150611596868383876115a0565b9695505050505050565b6060831561160f578251600003611608576001600160a01b0385163b6116085760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610597565b5081610512565b61051283838151156116245781518083602001fd5b8060405162461bcd60e51b815260040161059791906116ee565b6001600160a01b038116811461099357600080fd5b6000806000806080858703121561166957600080fd5b84356116748161163e565b935060208501356116848161163e565b925060408501359150606085013561169b8161163e565b939692955090935050565b6000602082840312156116b857600080fd5b81356116c38161163e565b9392505050565b60005b838110156116e55781810151838201526020016116cd565b50506000910152565b602081526000825180602084015261170d8160408501602087016116ca565b601f01601f19169190910160400192915050565b60008060006060848603121561173657600080fd5b83356117418161163e565b925060208401356117518161163e565b929592945050506040919091013590565b60006020828403121561177457600080fd5b5035919050565b6000806040838503121561178e57600080fd5b82356117998161163e565b946020939093013593505050565b634e487b7160e01b600052604160045260246000fd5b600080604083850312156117d057600080fd5b82356117db8161163e565b9150602083013567ffffffffffffffff808211156117f857600080fd5b818501915085601f83011261180c57600080fd5b81358181111561181e5761181e6117a7565b604051601f8201601f19908116603f01168101908382118183101715611846576118466117a7565b8160405282815288602084870101111561185f57600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b60006020828403121561189357600080fd5b81516116c38161163e565b6000602082840312156118b057600080fd5b5051919050565b808201808211156114ae57634e487b7160e01b600052601160045260246000fd5b600082516118ea8184602087016116ca565b919091019291505056fe416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a2646970667358221220234abed121cc71346ebb057f0279c561c01dfe59fdb903ba5d055bc03090122464736f6c63430008110033",
}

// SupplyLimitModuleABI is the input ABI used to generate the binding from.
// Deprecated: Use SupplyLimitModuleMetaData.ABI instead.
var SupplyLimitModuleABI = SupplyLimitModuleMetaData.ABI

// SupplyLimitModuleBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SupplyLimitModuleMetaData.Bin instead.
var SupplyLimitModuleBin = SupplyLimitModuleMetaData.Bin

// DeploySupplyLimitModule deploys a new Ethereum contract, binding an instance of SupplyLimitModule to it.
func DeploySupplyLimitModule(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SupplyLimitModule, error) {
	parsed, err := SupplyLimitModuleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SupplyLimitModuleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SupplyLimitModule{SupplyLimitModuleCaller: SupplyLimitModuleCaller{contract: contract}, SupplyLimitModuleTransactor: SupplyLimitModuleTransactor{contract: contract}, SupplyLimitModuleFilterer: SupplyLimitModuleFilterer{contract: contract}}, nil
}

// SupplyLimitModule is an auto generated Go binding around an Ethereum contract.
type SupplyLimitModule struct {
	SupplyLimitModuleCaller     // Read-only binding to the contract
	SupplyLimitModuleTransactor // Write-only binding to the contract
	SupplyLimitModuleFilterer   // Log filterer for contract events
}

// SupplyLimitModuleCaller is an auto generated read-only Go binding around an Ethereum contract.
type SupplyLimitModuleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SupplyLimitModuleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SupplyLimitModuleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SupplyLimitModuleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SupplyLimitModuleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SupplyLimitModuleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SupplyLimitModuleSession struct {
	Contract     *SupplyLimitModule // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// SupplyLimitModuleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SupplyLimitModuleCallerSession struct {
	Contract *SupplyLimitModuleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// SupplyLimitModuleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SupplyLimitModuleTransactorSession struct {
	Contract     *SupplyLimitModuleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// SupplyLimitModuleRaw is an auto generated low-level Go binding around an Ethereum contract.
type SupplyLimitModuleRaw struct {
	Contract *SupplyLimitModule // Generic contract binding to access the raw methods on
}

// SupplyLimitModuleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SupplyLimitModuleCallerRaw struct {
	Contract *SupplyLimitModuleCaller // Generic read-only contract binding to access the raw methods on
}

// SupplyLimitModuleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SupplyLimitModuleTransactorRaw struct {
	Contract *SupplyLimitModuleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSupplyLimitModule creates a new instance of SupplyLimitModule, bound to a specific deployed contract.
func NewSupplyLimitModule(address common.Address, backend bind.ContractBackend) (*SupplyLimitModule, error) {
	contract, err := bindSupplyLimitModule(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SupplyLimitModule{SupplyLimitModuleCaller: SupplyLimitModuleCaller{contract: contract}, SupplyLimitModuleTransactor: SupplyLimitModuleTransactor{contract: contract}, SupplyLimitModuleFilterer: SupplyLimitModuleFilterer{contract: contract}}, nil
}

// NewSupplyLimitModuleCaller creates a new read-only instance of SupplyLimitModule, bound to a specific deployed contract.
func NewSupplyLimitModuleCaller(address common.Address, caller bind.ContractCaller) (*SupplyLimitModuleCaller, error) {
	contract, err := bindSupplyLimitModule(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SupplyLimitModuleCaller{contract: contract}, nil
}

// NewSupplyLimitModuleTransactor creates a new write-only instance of SupplyLimitModule, bound to a specific deployed contract.
func NewSupplyLimitModuleTransactor(address common.Address, transactor bind.ContractTransactor) (*SupplyLimitModuleTransactor, error) {
	contract, err := bindSupplyLimitModule(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SupplyLimitModuleTransactor{contract: contract}, nil
}

// NewSupplyLimitModuleFilterer creates a new log filterer instance of SupplyLimitModule, bound to a specific deployed contract.
func NewSupplyLimitModuleFilterer(address common.Address, filterer bind.ContractFilterer) (*SupplyLimitModuleFilterer, error) {
	contract, err := bindSupplyLimitModule(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SupplyLimitModuleFilterer{contract: contract}, nil
}

// bindSupplyLimitModule binds a generic wrapper to an already deployed contract.
func bindSupplyLimitModule(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SupplyLimitModuleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SupplyLimitModule *SupplyLimitModuleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SupplyLimitModule.Contract.SupplyLimitModuleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SupplyLimitModule *SupplyLimitModuleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.SupplyLimitModuleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SupplyLimitModule *SupplyLimitModuleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.SupplyLimitModuleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SupplyLimitModule *SupplyLimitModuleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SupplyLimitModule.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SupplyLimitModule *SupplyLimitModuleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SupplyLimitModule *SupplyLimitModuleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.contract.Transact(opts, method, params...)
}

// CanComplianceBind is a free data retrieval call binding the contract method 0xbcc21053.
//
// Solidity: function canComplianceBind(address ) view returns(bool)
func (_SupplyLimitModule *SupplyLimitModuleCaller) CanComplianceBind(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _SupplyLimitModule.contract.Call(opts, &out, "canComplianceBind", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanComplianceBind is a free data retrieval call binding the contract method 0xbcc21053.
//
// Solidity: function canComplianceBind(address ) view returns(bool)
func (_SupplyLimitModule *SupplyLimitModuleSession) CanComplianceBind(arg0 common.Address) (bool, error) {
	return _SupplyLimitModule.Contract.CanComplianceBind(&_SupplyLimitModule.CallOpts, arg0)
}

// CanComplianceBind is a free data retrieval call binding the contract method 0xbcc21053.
//
// Solidity: function canComplianceBind(address ) view returns(bool)
func (_SupplyLimitModule *SupplyLimitModuleCallerSession) CanComplianceBind(arg0 common.Address) (bool, error) {
	return _SupplyLimitModule.Contract.CanComplianceBind(&_SupplyLimitModule.CallOpts, arg0)
}

// GetSupplyLimit is a free data retrieval call binding the contract method 0x4ae65b5c.
//
// Solidity: function getSupplyLimit(address _compliance) view returns(uint256)
func (_SupplyLimitModule *SupplyLimitModuleCaller) GetSupplyLimit(opts *bind.CallOpts, _compliance common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SupplyLimitModule.contract.Call(opts, &out, "getSupplyLimit", _compliance)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetSupplyLimit is a free data retrieval call binding the contract method 0x4ae65b5c.
//
// Solidity: function getSupplyLimit(address _compliance) view returns(uint256)
func (_SupplyLimitModule *SupplyLimitModuleSession) GetSupplyLimit(_compliance common.Address) (*big.Int, error) {
	return _SupplyLimitModule.Contract.GetSupplyLimit(&_SupplyLimitModule.CallOpts, _compliance)
}

// GetSupplyLimit is a free data retrieval call binding the contract method 0x4ae65b5c.
//
// Solidity: function getSupplyLimit(address _compliance) view returns(uint256)
func (_SupplyLimitModule *SupplyLimitModuleCallerSession) GetSupplyLimit(_compliance common.Address) (*big.Int, error) {
	return _SupplyLimitModule.Contract.GetSupplyLimit(&_SupplyLimitModule.CallOpts, _compliance)
}

// IsComplianceBound is a free data retrieval call binding the contract method 0x4cf4d295.
//
// Solidity: function isComplianceBound(address _compliance) view returns(bool)
func (_SupplyLimitModule *SupplyLimitModuleCaller) IsComplianceBound(opts *bind.CallOpts, _compliance common.Address) (bool, error) {
	var out []interface{}
	err := _SupplyLimitModule.contract.Call(opts, &out, "isComplianceBound", _compliance)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsComplianceBound is a free data retrieval call binding the contract method 0x4cf4d295.
//
// Solidity: function isComplianceBound(address _compliance) view returns(bool)
func (_SupplyLimitModule *SupplyLimitModuleSession) IsComplianceBound(_compliance common.Address) (bool, error) {
	return _SupplyLimitModule.Contract.IsComplianceBound(&_SupplyLimitModule.CallOpts, _compliance)
}

// IsComplianceBound is a free data retrieval call binding the contract method 0x4cf4d295.
//
// Solidity: function isComplianceBound(address _compliance) view returns(bool)
func (_SupplyLimitModule *SupplyLimitModuleCallerSession) IsComplianceBound(_compliance common.Address) (bool, error) {
	return _SupplyLimitModule.Contract.IsComplianceBound(&_SupplyLimitModule.CallOpts, _compliance)
}

// IsPlugAndPlay is a free data retrieval call binding the contract method 0xe6f5e807.
//
// Solidity: function isPlugAndPlay() pure returns(bool)
func (_SupplyLimitModule *SupplyLimitModuleCaller) IsPlugAndPlay(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SupplyLimitModule.contract.Call(opts, &out, "isPlugAndPlay")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPlugAndPlay is a free data retrieval call binding the contract method 0xe6f5e807.
//
// Solidity: function isPlugAndPlay() pure returns(bool)
func (_SupplyLimitModule *SupplyLimitModuleSession) IsPlugAndPlay() (bool, error) {
	return _SupplyLimitModule.Contract.IsPlugAndPlay(&_SupplyLimitModule.CallOpts)
}

// IsPlugAndPlay is a free data retrieval call binding the contract method 0xe6f5e807.
//
// Solidity: function isPlugAndPlay() pure returns(bool)
func (_SupplyLimitModule *SupplyLimitModuleCallerSession) IsPlugAndPlay() (bool, error) {
	return _SupplyLimitModule.Contract.IsPlugAndPlay(&_SupplyLimitModule.CallOpts)
}

// ModuleCheck is a free data retrieval call binding the contract method 0x013b7ce4.
//
// Solidity: function moduleCheck(address _from, address , uint256 _value, address _compliance) view returns(bool)
func (_SupplyLimitModule *SupplyLimitModuleCaller) ModuleCheck(opts *bind.CallOpts, _from common.Address, arg1 common.Address, _value *big.Int, _compliance common.Address) (bool, error) {
	var out []interface{}
	err := _SupplyLimitModule.contract.Call(opts, &out, "moduleCheck", _from, arg1, _value, _compliance)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ModuleCheck is a free data retrieval call binding the contract method 0x013b7ce4.
//
// Solidity: function moduleCheck(address _from, address , uint256 _value, address _compliance) view returns(bool)
func (_SupplyLimitModule *SupplyLimitModuleSession) ModuleCheck(_from common.Address, arg1 common.Address, _value *big.Int, _compliance common.Address) (bool, error) {
	return _SupplyLimitModule.Contract.ModuleCheck(&_SupplyLimitModule.CallOpts, _from, arg1, _value, _compliance)
}

// ModuleCheck is a free data retrieval call binding the contract method 0x013b7ce4.
//
// Solidity: function moduleCheck(address _from, address , uint256 _value, address _compliance) view returns(bool)
func (_SupplyLimitModule *SupplyLimitModuleCallerSession) ModuleCheck(_from common.Address, arg1 common.Address, _value *big.Int, _compliance common.Address) (bool, error) {
	return _SupplyLimitModule.Contract.ModuleCheck(&_SupplyLimitModule.CallOpts, _from, arg1, _value, _compliance)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string _name)
func (_SupplyLimitModule *SupplyLimitModuleCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SupplyLimitModule.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string _name)
func (_SupplyLimitModule *SupplyLimitModuleSession) Name() (string, error) {
	return _SupplyLimitModule.Contract.Name(&_SupplyLimitModule.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string _name)
func (_SupplyLimitModule *SupplyLimitModuleCallerSession) Name() (string, error) {
	return _SupplyLimitModule.Contract.Name(&_SupplyLimitModule.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SupplyLimitModule *SupplyLimitModuleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SupplyLimitModule.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SupplyLimitModule *SupplyLimitModuleSession) Owner() (common.Address, error) {
	return _SupplyLimitModule.Contract.Owner(&_SupplyLimitModule.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SupplyLimitModule *SupplyLimitModuleCallerSession) Owner() (common.Address, error) {
	return _SupplyLimitModule.Contract.Owner(&_SupplyLimitModule.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_SupplyLimitModule *SupplyLimitModuleCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _SupplyLimitModule.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_SupplyLimitModule *SupplyLimitModuleSession) ProxiableUUID() ([32]byte, error) {
	return _SupplyLimitModule.Contract.ProxiableUUID(&_SupplyLimitModule.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_SupplyLimitModule *SupplyLimitModuleCallerSession) ProxiableUUID() ([32]byte, error) {
	return _SupplyLimitModule.Contract.ProxiableUUID(&_SupplyLimitModule.CallOpts)
}

// BindCompliance is a paid mutator transaction binding the contract method 0x4a932544.
//
// Solidity: function bindCompliance(address _compliance) returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactor) BindCompliance(opts *bind.TransactOpts, _compliance common.Address) (*types.Transaction, error) {
	return _SupplyLimitModule.contract.Transact(opts, "bindCompliance", _compliance)
}

// BindCompliance is a paid mutator transaction binding the contract method 0x4a932544.
//
// Solidity: function bindCompliance(address _compliance) returns()
func (_SupplyLimitModule *SupplyLimitModuleSession) BindCompliance(_compliance common.Address) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.BindCompliance(&_SupplyLimitModule.TransactOpts, _compliance)
}

// BindCompliance is a paid mutator transaction binding the contract method 0x4a932544.
//
// Solidity: function bindCompliance(address _compliance) returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactorSession) BindCompliance(_compliance common.Address) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.BindCompliance(&_SupplyLimitModule.TransactOpts, _compliance)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SupplyLimitModule.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_SupplyLimitModule *SupplyLimitModuleSession) Initialize() (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.Initialize(&_SupplyLimitModule.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactorSession) Initialize() (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.Initialize(&_SupplyLimitModule.TransactOpts)
}

// ModuleBurnAction is a paid mutator transaction binding the contract method 0x372491a2.
//
// Solidity: function moduleBurnAction(address _from, uint256 _value) returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactor) ModuleBurnAction(opts *bind.TransactOpts, _from common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SupplyLimitModule.contract.Transact(opts, "moduleBurnAction", _from, _value)
}

// ModuleBurnAction is a paid mutator transaction binding the contract method 0x372491a2.
//
// Solidity: function moduleBurnAction(address _from, uint256 _value) returns()
func (_SupplyLimitModule *SupplyLimitModuleSession) ModuleBurnAction(_from common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.ModuleBurnAction(&_SupplyLimitModule.TransactOpts, _from, _value)
}

// ModuleBurnAction is a paid mutator transaction binding the contract method 0x372491a2.
//
// Solidity: function moduleBurnAction(address _from, uint256 _value) returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactorSession) ModuleBurnAction(_from common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.ModuleBurnAction(&_SupplyLimitModule.TransactOpts, _from, _value)
}

// ModuleMintAction is a paid mutator transaction binding the contract method 0xf104a8c9.
//
// Solidity: function moduleMintAction(address _to, uint256 _value) returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactor) ModuleMintAction(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SupplyLimitModule.contract.Transact(opts, "moduleMintAction", _to, _value)
}

// ModuleMintAction is a paid mutator transaction binding the contract method 0xf104a8c9.
//
// Solidity: function moduleMintAction(address _to, uint256 _value) returns()
func (_SupplyLimitModule *SupplyLimitModuleSession) ModuleMintAction(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.ModuleMintAction(&_SupplyLimitModule.TransactOpts, _to, _value)
}

// ModuleMintAction is a paid mutator transaction binding the contract method 0xf104a8c9.
//
// Solidity: function moduleMintAction(address _to, uint256 _value) returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactorSession) ModuleMintAction(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.ModuleMintAction(&_SupplyLimitModule.TransactOpts, _to, _value)
}

// ModuleTransferAction is a paid mutator transaction binding the contract method 0x2cb7e1ec.
//
// Solidity: function moduleTransferAction(address _from, address _to, uint256 _value) returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactor) ModuleTransferAction(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SupplyLimitModule.contract.Transact(opts, "moduleTransferAction", _from, _to, _value)
}

// ModuleTransferAction is a paid mutator transaction binding the contract method 0x2cb7e1ec.
//
// Solidity: function moduleTransferAction(address _from, address _to, uint256 _value) returns()
func (_SupplyLimitModule *SupplyLimitModuleSession) ModuleTransferAction(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.ModuleTransferAction(&_SupplyLimitModule.TransactOpts, _from, _to, _value)
}

// ModuleTransferAction is a paid mutator transaction binding the contract method 0x2cb7e1ec.
//
// Solidity: function moduleTransferAction(address _from, address _to, uint256 _value) returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactorSession) ModuleTransferAction(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.ModuleTransferAction(&_SupplyLimitModule.TransactOpts, _from, _to, _value)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SupplyLimitModule.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SupplyLimitModule *SupplyLimitModuleSession) RenounceOwnership() (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.RenounceOwnership(&_SupplyLimitModule.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.RenounceOwnership(&_SupplyLimitModule.TransactOpts)
}

// SetSupplyLimit is a paid mutator transaction binding the contract method 0x361fab25.
//
// Solidity: function setSupplyLimit(uint256 _limit) returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactor) SetSupplyLimit(opts *bind.TransactOpts, _limit *big.Int) (*types.Transaction, error) {
	return _SupplyLimitModule.contract.Transact(opts, "setSupplyLimit", _limit)
}

// SetSupplyLimit is a paid mutator transaction binding the contract method 0x361fab25.
//
// Solidity: function setSupplyLimit(uint256 _limit) returns()
func (_SupplyLimitModule *SupplyLimitModuleSession) SetSupplyLimit(_limit *big.Int) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.SetSupplyLimit(&_SupplyLimitModule.TransactOpts, _limit)
}

// SetSupplyLimit is a paid mutator transaction binding the contract method 0x361fab25.
//
// Solidity: function setSupplyLimit(uint256 _limit) returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactorSession) SetSupplyLimit(_limit *big.Int) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.SetSupplyLimit(&_SupplyLimitModule.TransactOpts, _limit)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SupplyLimitModule.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SupplyLimitModule *SupplyLimitModuleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.TransferOwnership(&_SupplyLimitModule.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.TransferOwnership(&_SupplyLimitModule.TransactOpts, newOwner)
}

// UnbindCompliance is a paid mutator transaction binding the contract method 0x0694a5fb.
//
// Solidity: function unbindCompliance(address _compliance) returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactor) UnbindCompliance(opts *bind.TransactOpts, _compliance common.Address) (*types.Transaction, error) {
	return _SupplyLimitModule.contract.Transact(opts, "unbindCompliance", _compliance)
}

// UnbindCompliance is a paid mutator transaction binding the contract method 0x0694a5fb.
//
// Solidity: function unbindCompliance(address _compliance) returns()
func (_SupplyLimitModule *SupplyLimitModuleSession) UnbindCompliance(_compliance common.Address) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.UnbindCompliance(&_SupplyLimitModule.TransactOpts, _compliance)
}

// UnbindCompliance is a paid mutator transaction binding the contract method 0x0694a5fb.
//
// Solidity: function unbindCompliance(address _compliance) returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactorSession) UnbindCompliance(_compliance common.Address) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.UnbindCompliance(&_SupplyLimitModule.TransactOpts, _compliance)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _SupplyLimitModule.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_SupplyLimitModule *SupplyLimitModuleSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.UpgradeTo(&_SupplyLimitModule.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.UpgradeTo(&_SupplyLimitModule.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _SupplyLimitModule.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_SupplyLimitModule *SupplyLimitModuleSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.UpgradeToAndCall(&_SupplyLimitModule.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_SupplyLimitModule *SupplyLimitModuleTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _SupplyLimitModule.Contract.UpgradeToAndCall(&_SupplyLimitModule.TransactOpts, newImplementation, data)
}

// SupplyLimitModuleAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the SupplyLimitModule contract.
type SupplyLimitModuleAdminChangedIterator struct {
	Event *SupplyLimitModuleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SupplyLimitModuleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SupplyLimitModuleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SupplyLimitModuleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SupplyLimitModuleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SupplyLimitModuleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SupplyLimitModuleAdminChanged represents a AdminChanged event raised by the SupplyLimitModule contract.
type SupplyLimitModuleAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*SupplyLimitModuleAdminChangedIterator, error) {

	logs, sub, err := _SupplyLimitModule.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &SupplyLimitModuleAdminChangedIterator{contract: _SupplyLimitModule.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *SupplyLimitModuleAdminChanged) (event.Subscription, error) {

	logs, sub, err := _SupplyLimitModule.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SupplyLimitModuleAdminChanged)
				if err := _SupplyLimitModule.contract.UnpackLog(event, "AdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) ParseAdminChanged(log types.Log) (*SupplyLimitModuleAdminChanged, error) {
	event := new(SupplyLimitModuleAdminChanged)
	if err := _SupplyLimitModule.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SupplyLimitModuleBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the SupplyLimitModule contract.
type SupplyLimitModuleBeaconUpgradedIterator struct {
	Event *SupplyLimitModuleBeaconUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SupplyLimitModuleBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SupplyLimitModuleBeaconUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SupplyLimitModuleBeaconUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SupplyLimitModuleBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SupplyLimitModuleBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SupplyLimitModuleBeaconUpgraded represents a BeaconUpgraded event raised by the SupplyLimitModule contract.
type SupplyLimitModuleBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*SupplyLimitModuleBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _SupplyLimitModule.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &SupplyLimitModuleBeaconUpgradedIterator{contract: _SupplyLimitModule.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *SupplyLimitModuleBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _SupplyLimitModule.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SupplyLimitModuleBeaconUpgraded)
				if err := _SupplyLimitModule.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) ParseBeaconUpgraded(log types.Log) (*SupplyLimitModuleBeaconUpgraded, error) {
	event := new(SupplyLimitModuleBeaconUpgraded)
	if err := _SupplyLimitModule.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SupplyLimitModuleComplianceBoundIterator is returned from FilterComplianceBound and is used to iterate over the raw logs and unpacked data for ComplianceBound events raised by the SupplyLimitModule contract.
type SupplyLimitModuleComplianceBoundIterator struct {
	Event *SupplyLimitModuleComplianceBound // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SupplyLimitModuleComplianceBoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SupplyLimitModuleComplianceBound)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SupplyLimitModuleComplianceBound)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SupplyLimitModuleComplianceBoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SupplyLimitModuleComplianceBoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SupplyLimitModuleComplianceBound represents a ComplianceBound event raised by the SupplyLimitModule contract.
type SupplyLimitModuleComplianceBound struct {
	Compliance common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterComplianceBound is a free log retrieval operation binding the contract event 0x1f7b76c58fb697eb53c6c7c1becb96911516a136e24d7ced386b2355358b75a3.
//
// Solidity: event ComplianceBound(address indexed _compliance)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) FilterComplianceBound(opts *bind.FilterOpts, _compliance []common.Address) (*SupplyLimitModuleComplianceBoundIterator, error) {

	var _complianceRule []interface{}
	for _, _complianceItem := range _compliance {
		_complianceRule = append(_complianceRule, _complianceItem)
	}

	logs, sub, err := _SupplyLimitModule.contract.FilterLogs(opts, "ComplianceBound", _complianceRule)
	if err != nil {
		return nil, err
	}
	return &SupplyLimitModuleComplianceBoundIterator{contract: _SupplyLimitModule.contract, event: "ComplianceBound", logs: logs, sub: sub}, nil
}

// WatchComplianceBound is a free log subscription operation binding the contract event 0x1f7b76c58fb697eb53c6c7c1becb96911516a136e24d7ced386b2355358b75a3.
//
// Solidity: event ComplianceBound(address indexed _compliance)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) WatchComplianceBound(opts *bind.WatchOpts, sink chan<- *SupplyLimitModuleComplianceBound, _compliance []common.Address) (event.Subscription, error) {

	var _complianceRule []interface{}
	for _, _complianceItem := range _compliance {
		_complianceRule = append(_complianceRule, _complianceItem)
	}

	logs, sub, err := _SupplyLimitModule.contract.WatchLogs(opts, "ComplianceBound", _complianceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SupplyLimitModuleComplianceBound)
				if err := _SupplyLimitModule.contract.UnpackLog(event, "ComplianceBound", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseComplianceBound is a log parse operation binding the contract event 0x1f7b76c58fb697eb53c6c7c1becb96911516a136e24d7ced386b2355358b75a3.
//
// Solidity: event ComplianceBound(address indexed _compliance)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) ParseComplianceBound(log types.Log) (*SupplyLimitModuleComplianceBound, error) {
	event := new(SupplyLimitModuleComplianceBound)
	if err := _SupplyLimitModule.contract.UnpackLog(event, "ComplianceBound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SupplyLimitModuleComplianceUnboundIterator is returned from FilterComplianceUnbound and is used to iterate over the raw logs and unpacked data for ComplianceUnbound events raised by the SupplyLimitModule contract.
type SupplyLimitModuleComplianceUnboundIterator struct {
	Event *SupplyLimitModuleComplianceUnbound // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SupplyLimitModuleComplianceUnboundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SupplyLimitModuleComplianceUnbound)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SupplyLimitModuleComplianceUnbound)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SupplyLimitModuleComplianceUnboundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SupplyLimitModuleComplianceUnboundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SupplyLimitModuleComplianceUnbound represents a ComplianceUnbound event raised by the SupplyLimitModule contract.
type SupplyLimitModuleComplianceUnbound struct {
	Compliance common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterComplianceUnbound is a free log retrieval operation binding the contract event 0x408b49d9be1c914c52a0227e18a077e5a892dddf32a26cfa94a5d9708fad7718.
//
// Solidity: event ComplianceUnbound(address indexed _compliance)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) FilterComplianceUnbound(opts *bind.FilterOpts, _compliance []common.Address) (*SupplyLimitModuleComplianceUnboundIterator, error) {

	var _complianceRule []interface{}
	for _, _complianceItem := range _compliance {
		_complianceRule = append(_complianceRule, _complianceItem)
	}

	logs, sub, err := _SupplyLimitModule.contract.FilterLogs(opts, "ComplianceUnbound", _complianceRule)
	if err != nil {
		return nil, err
	}
	return &SupplyLimitModuleComplianceUnboundIterator{contract: _SupplyLimitModule.contract, event: "ComplianceUnbound", logs: logs, sub: sub}, nil
}

// WatchComplianceUnbound is a free log subscription operation binding the contract event 0x408b49d9be1c914c52a0227e18a077e5a892dddf32a26cfa94a5d9708fad7718.
//
// Solidity: event ComplianceUnbound(address indexed _compliance)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) WatchComplianceUnbound(opts *bind.WatchOpts, sink chan<- *SupplyLimitModuleComplianceUnbound, _compliance []common.Address) (event.Subscription, error) {

	var _complianceRule []interface{}
	for _, _complianceItem := range _compliance {
		_complianceRule = append(_complianceRule, _complianceItem)
	}

	logs, sub, err := _SupplyLimitModule.contract.WatchLogs(opts, "ComplianceUnbound", _complianceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SupplyLimitModuleComplianceUnbound)
				if err := _SupplyLimitModule.contract.UnpackLog(event, "ComplianceUnbound", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseComplianceUnbound is a log parse operation binding the contract event 0x408b49d9be1c914c52a0227e18a077e5a892dddf32a26cfa94a5d9708fad7718.
//
// Solidity: event ComplianceUnbound(address indexed _compliance)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) ParseComplianceUnbound(log types.Log) (*SupplyLimitModuleComplianceUnbound, error) {
	event := new(SupplyLimitModuleComplianceUnbound)
	if err := _SupplyLimitModule.contract.UnpackLog(event, "ComplianceUnbound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SupplyLimitModuleInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the SupplyLimitModule contract.
type SupplyLimitModuleInitializedIterator struct {
	Event *SupplyLimitModuleInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SupplyLimitModuleInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SupplyLimitModuleInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SupplyLimitModuleInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SupplyLimitModuleInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SupplyLimitModuleInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SupplyLimitModuleInitialized represents a Initialized event raised by the SupplyLimitModule contract.
type SupplyLimitModuleInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) FilterInitialized(opts *bind.FilterOpts) (*SupplyLimitModuleInitializedIterator, error) {

	logs, sub, err := _SupplyLimitModule.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &SupplyLimitModuleInitializedIterator{contract: _SupplyLimitModule.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *SupplyLimitModuleInitialized) (event.Subscription, error) {

	logs, sub, err := _SupplyLimitModule.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SupplyLimitModuleInitialized)
				if err := _SupplyLimitModule.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) ParseInitialized(log types.Log) (*SupplyLimitModuleInitialized, error) {
	event := new(SupplyLimitModuleInitialized)
	if err := _SupplyLimitModule.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SupplyLimitModuleOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SupplyLimitModule contract.
type SupplyLimitModuleOwnershipTransferredIterator struct {
	Event *SupplyLimitModuleOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SupplyLimitModuleOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SupplyLimitModuleOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SupplyLimitModuleOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SupplyLimitModuleOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SupplyLimitModuleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SupplyLimitModuleOwnershipTransferred represents a OwnershipTransferred event raised by the SupplyLimitModule contract.
type SupplyLimitModuleOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SupplyLimitModuleOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SupplyLimitModule.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SupplyLimitModuleOwnershipTransferredIterator{contract: _SupplyLimitModule.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SupplyLimitModuleOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SupplyLimitModule.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SupplyLimitModuleOwnershipTransferred)
				if err := _SupplyLimitModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) ParseOwnershipTransferred(log types.Log) (*SupplyLimitModuleOwnershipTransferred, error) {
	event := new(SupplyLimitModuleOwnershipTransferred)
	if err := _SupplyLimitModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SupplyLimitModuleSupplyLimitSetIterator is returned from FilterSupplyLimitSet and is used to iterate over the raw logs and unpacked data for SupplyLimitSet events raised by the SupplyLimitModule contract.
type SupplyLimitModuleSupplyLimitSetIterator struct {
	Event *SupplyLimitModuleSupplyLimitSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SupplyLimitModuleSupplyLimitSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SupplyLimitModuleSupplyLimitSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SupplyLimitModuleSupplyLimitSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SupplyLimitModuleSupplyLimitSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SupplyLimitModuleSupplyLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SupplyLimitModuleSupplyLimitSet represents a SupplyLimitSet event raised by the SupplyLimitModule contract.
type SupplyLimitModuleSupplyLimitSet struct {
	Compliance common.Address
	Limit      *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSupplyLimitSet is a free log retrieval operation binding the contract event 0x2f1ac7412e0c36ff7b73f7c4f474022dfa5d3f1c09e3262f63c9eb4627a9f505.
//
// Solidity: event SupplyLimitSet(address _compliance, uint256 _limit)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) FilterSupplyLimitSet(opts *bind.FilterOpts) (*SupplyLimitModuleSupplyLimitSetIterator, error) {

	logs, sub, err := _SupplyLimitModule.contract.FilterLogs(opts, "SupplyLimitSet")
	if err != nil {
		return nil, err
	}
	return &SupplyLimitModuleSupplyLimitSetIterator{contract: _SupplyLimitModule.contract, event: "SupplyLimitSet", logs: logs, sub: sub}, nil
}

// WatchSupplyLimitSet is a free log subscription operation binding the contract event 0x2f1ac7412e0c36ff7b73f7c4f474022dfa5d3f1c09e3262f63c9eb4627a9f505.
//
// Solidity: event SupplyLimitSet(address _compliance, uint256 _limit)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) WatchSupplyLimitSet(opts *bind.WatchOpts, sink chan<- *SupplyLimitModuleSupplyLimitSet) (event.Subscription, error) {

	logs, sub, err := _SupplyLimitModule.contract.WatchLogs(opts, "SupplyLimitSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SupplyLimitModuleSupplyLimitSet)
				if err := _SupplyLimitModule.contract.UnpackLog(event, "SupplyLimitSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSupplyLimitSet is a log parse operation binding the contract event 0x2f1ac7412e0c36ff7b73f7c4f474022dfa5d3f1c09e3262f63c9eb4627a9f505.
//
// Solidity: event SupplyLimitSet(address _compliance, uint256 _limit)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) ParseSupplyLimitSet(log types.Log) (*SupplyLimitModuleSupplyLimitSet, error) {
	event := new(SupplyLimitModuleSupplyLimitSet)
	if err := _SupplyLimitModule.contract.UnpackLog(event, "SupplyLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SupplyLimitModuleUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the SupplyLimitModule contract.
type SupplyLimitModuleUpgradedIterator struct {
	Event *SupplyLimitModuleUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SupplyLimitModuleUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SupplyLimitModuleUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SupplyLimitModuleUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SupplyLimitModuleUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SupplyLimitModuleUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SupplyLimitModuleUpgraded represents a Upgraded event raised by the SupplyLimitModule contract.
type SupplyLimitModuleUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*SupplyLimitModuleUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _SupplyLimitModule.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &SupplyLimitModuleUpgradedIterator{contract: _SupplyLimitModule.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *SupplyLimitModuleUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _SupplyLimitModule.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SupplyLimitModuleUpgraded)
				if err := _SupplyLimitModule.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_SupplyLimitModule *SupplyLimitModuleFilterer) ParseUpgraded(log types.Log) (*SupplyLimitModuleUpgraded, error) {
	event := new(SupplyLimitModuleUpgraded)
	if err := _SupplyLimitModule.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
