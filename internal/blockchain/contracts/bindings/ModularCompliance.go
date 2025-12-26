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

// ModularComplianceMetaData contains all meta data concerning the ModularCompliance contract.
var ModularComplianceMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_module\",\"type\":\"address\"}],\"name\":\"ModuleAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"selector\",\"type\":\"bytes4\"}],\"name\":\"ModuleInteraction\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_module\",\"type\":\"address\"}],\"name\":\"ModuleRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"TokenBound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"TokenUnbound\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_module\",\"type\":\"address\"}],\"name\":\"addModule\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"bindToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"_module\",\"type\":\"address\"}],\"name\":\"callModuleFunction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"canTransfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"created\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"destroyed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getModules\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTokenBound\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_module\",\"type\":\"address\"}],\"name\":\"isModuleBound\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_module\",\"type\":\"address\"}],\"name\":\"removeModule\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferred\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"unbindToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061192d806100206000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c80638da5cb5b11610097578063e1c7392a11610066578063e1c7392a14610220578063e46638e614610228578063efb22d331461023b578063f2fde38b1461024e57600080fd5b80638da5cb5b146101ab578063a0632461146101bc578063a446d49f146101cf578063b2494df31461020b57600080fd5b80636a3edf28116100d35780636a3edf2814610153578063715018a61461017d5780638baf29b4146101855780638d2ea7721461019857600080fd5b80631ed86f19146101055780633ff5aa021461011a57806340db3b501461012d5780635f8dead314610140575b600080fd5b610118610113366004611718565b610261565b005b610118610128366004611718565b6105e0565b61011861013b366004611718565b61071b565b61011861014e366004611733565b610893565b6065546001600160a01b03165b6040516001600160a01b0390911681526020015b60405180910390f35b610118610a89565b61011861019336600461175d565b610a9d565b6101186101a6366004611733565b610cd7565b6033546001600160a01b0316610160565b6101186101ca366004611718565b610ec7565b6101fb6101dd366004611718565b6001600160a01b031660009081526067602052604090205460ff1690565b6040519015158152602001610174565b610213611155565b6040516101749190611799565b6101186111b7565b6101fb61023636600461175d565b6112d3565b6101186102493660046117e6565b6113ca565b61011861025c366004611718565b6114c8565b610269611555565b6001600160a01b0381166102c45760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064015b60405180910390fd5b6001600160a01b03811660009081526067602052604090205460ff161561032d5760405162461bcd60e51b815260206004820152601460248201527f6d6f64756c6520616c726561647920626f756e6400000000000000000000000060448201526064016102bb565b606654601810156103805760405162461bcd60e51b815260206004820152601f60248201527f63616e6e6f7420616464206d6f7265207468616e203235206d6f64756c65730060448201526064016102bb565b6000819050806001600160a01b031663e6f5e8076040518163ffffffff1660e01b8152600401602060405180830381865afa1580156103c3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103e79190611867565b6104de576040517fbcc210530000000000000000000000000000000000000000000000000000000081523060048201526001600160a01b0382169063bcc2105390602401602060405180830381865afa158015610448573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061046c9190611867565b6104de5760405162461bcd60e51b815260206004820152603460248201527f636f6d706c69616e6365206973206e6f74207375697461626c6520666f72206260448201527f696e64696e6720746f20746865206d6f64756c6500000000000000000000000060648201526084016102bb565b6040517f4a9325440000000000000000000000000000000000000000000000000000000081523060048201526001600160a01b03821690634a93254490602401600060405180830381600087803b15801561053857600080fd5b505af115801561054c573d6000803e3d6000fd5b50506066805460018082019092557f46501879b8ca8525e8c2fd519e2fbfcfa2ebea26501294aa02cbfcfb12e943540180546001600160a01b0319166001600160a01b038716908117909155600081815260676020526040808220805460ff191690941790935591519093507fead6a006345da1073a106d5f32372d2d2204f46cb0b4bca8f5ebafcbbed12b8a9250a25050565b336105f36033546001600160a01b031690565b6001600160a01b0316148061062457506065546001600160a01b03161580156106245750336001600160a01b038216145b6106705760405162461bcd60e51b815260206004820152601c60248201527f6f6e6c79206f776e6572206f7220746f6b656e2063616e2063616c6c0000000060448201526064016102bb565b6001600160a01b0381166106c65760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064016102bb565b606580546001600160a01b0319166001600160a01b0383169081179091556040519081527f2de35142b19ed5a07796cf30791959c592018f70b1d2d7c460eef8ffe713692b906020015b60405180910390a150565b3361072e6033546001600160a01b031690565b6001600160a01b0316148061074b5750336001600160a01b038216145b6107975760405162461bcd60e51b815260206004820152601c60248201527f6f6e6c79206f776e6572206f7220746f6b656e2063616e2063616c6c0000000060448201526064016102bb565b6065546001600160a01b038281169116146107f45760405162461bcd60e51b815260206004820152601760248201527f5468697320746f6b656e206973206e6f7420626f756e6400000000000000000060448201526064016102bb565b6001600160a01b03811661084a5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064016102bb565b606580546001600160a01b03191690556040516001600160a01b03821681527f28a4ca7134a3b3f9aff286e79ad3daadb4a06d1b43d037a3a98bdc074edd9b7a90602001610710565b6065546001600160a01b031633146109215760405162461bcd60e51b8152602060048201526044602482018190527f6572726f72203a20746869732061646472657373206973206e6f74206120746f908201527f6b656e20626f756e6420746f2074686520636f6d706c69616e636520636f6e746064820152631c9858dd60e21b608482015260a4016102bb565b6001600160a01b0382166109775760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064016102bb565b600081116109c75760405162461bcd60e51b815260206004820181905260248201527f696e76616c696420617267756d656e74202d206e6f2076616c7565206d696e7460448201526064016102bb565b60665460005b81811015610a8357606681815481106109e8576109e8611889565b6000918252602090912001546040517ff104a8c90000000000000000000000000000000000000000000000000000000081526001600160a01b038681166004830152602482018690529091169063f104a8c990604401600060405180830381600087803b158015610a5857600080fd5b505af1158015610a6c573d6000803e3d6000fd5b505050508080610a7b906118b5565b9150506109cd565b50505050565b610a91611555565b610a9b60006115af565b565b6065546001600160a01b03163314610b2b5760405162461bcd60e51b8152602060048201526044602482018190527f6572726f72203a20746869732061646472657373206973206e6f74206120746f908201527f6b656e20626f756e6420746f2074686520636f6d706c69616e636520636f6e746064820152631c9858dd60e21b608482015260a4016102bb565b6001600160a01b03831615801590610b4b57506001600160a01b03821615155b610b975760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064016102bb565b60008111610c0c5760405162461bcd60e51b8152602060048201526024808201527f696e76616c696420617267756d656e74202d206e6f2076616c7565207472616e60448201527f736665720000000000000000000000000000000000000000000000000000000060648201526084016102bb565b60665460005b81811015610cd05760668181548110610c2d57610c2d611889565b6000918252602090912001546040517f2cb7e1ec0000000000000000000000000000000000000000000000000000000081526001600160a01b03878116600483015286811660248301526044820186905290911690632cb7e1ec90606401600060405180830381600087803b158015610ca557600080fd5b505af1158015610cb9573d6000803e3d6000fd5b505050508080610cc8906118b5565b915050610c12565b5050505050565b6065546001600160a01b03163314610d655760405162461bcd60e51b8152602060048201526044602482018190527f6572726f72203a20746869732061646472657373206973206e6f74206120746f908201527f6b656e20626f756e6420746f2074686520636f6d706c69616e636520636f6e746064820152631c9858dd60e21b608482015260a4016102bb565b6001600160a01b038216610dbb5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064016102bb565b60008111610e0b5760405162461bcd60e51b815260206004820181905260248201527f696e76616c696420617267756d656e74202d206e6f2076616c7565206275726e60448201526064016102bb565b60665460005b81811015610a835760668181548110610e2c57610e2c611889565b6000918252602090912001546040517f372491a20000000000000000000000000000000000000000000000000000000081526001600160a01b038681166004830152602482018690529091169063372491a290604401600060405180830381600087803b158015610e9c57600080fd5b505af1158015610eb0573d6000803e3d6000fd5b505050508080610ebf906118b5565b915050610e11565b610ecf611555565b6001600160a01b038116610f255760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064016102bb565b6001600160a01b03811660009081526067602052604090205460ff16610f8d5760405162461bcd60e51b815260206004820152601060248201527f6d6f64756c65206e6f7420626f756e640000000000000000000000000000000060448201526064016102bb565b60665460005b8181101561115057826001600160a01b031660668281548110610fb857610fb8611889565b6000918252602090912001546001600160a01b03160361113e576040517f0694a5fb0000000000000000000000000000000000000000000000000000000081523060048201526001600160a01b03841690630694a5fb90602401600060405180830381600087803b15801561102c57600080fd5b505af1158015611040573d6000803e3d6000fd5b50505050606660018361105391906118ce565b8154811061106357611063611889565b600091825260209091200154606680546001600160a01b03909216918390811061108f5761108f611889565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060668054806110ce576110ce6118e1565b60008281526020808220830160001990810180546001600160a01b03191690559092019092556001600160a01b03851680835260679091526040808320805460ff191690555190917f0a1ee69f55c33d8467c69ca59ce2007a737a88603d75392972520bf67cb513b891a2505050565b80611148816118b5565b915050610f93565b505050565b606060668054806020026020016040519081016040528092919081815260200182805480156111ad57602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161118f575b5050505050905090565b600054610100900460ff16158080156111d75750600054600160ff909116105b806111f15750303b1580156111f1575060005460ff166001145b6112635760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084016102bb565b6000805460ff191660011790558015611286576000805461ff0019166101001790555b61128e611601565b80156112d0576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb384740249890602001610710565b50565b606654600090815b818110156113bc57606681815481106112f6576112f6611889565b6000918252602090912001546040517f013b7ce40000000000000000000000000000000000000000000000000000000081526001600160a01b0388811660048301528781166024830152604482018790523060648301529091169063013b7ce490608401602060405180830381865afa158015611377573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061139b9190611867565b6113aa576000925050506113c3565b806113b4816118b5565b9150506112db565b5060019150505b9392505050565b6113d2611555565b6001600160a01b03811660009081526067602052604090205460ff1661143a5760405162461bcd60e51b815260206004820152601960248201527f63616c6c206f6e6c79206f6e20626f756e64206d6f64756c650000000000000060448201526064016102bb565b6040518284823760008084836000865af1611459573d6000803e3d6000fd5b50806001600160a01b03167f20d79de70adcc6e9353d8a9a5646b46dc352710d0a310b1ad1f67faeca7ef89161148f8585611674565b6040517fffffffff00000000000000000000000000000000000000000000000000000000909116815260200160405180910390a2505050565b6114d0611555565b6001600160a01b03811661154c5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f646472657373000000000000000000000000000000000000000000000000000060648201526084016102bb565b6112d0816115af565b6033546001600160a01b03163314610a9b5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016102bb565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff1661166c5760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b60648201526084016102bb565b610a9b611688565b600060048210611682575081355b92915050565b600054610100900460ff166116f35760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b60648201526084016102bb565b610a9b336115af565b80356001600160a01b038116811461171357600080fd5b919050565b60006020828403121561172a57600080fd5b6113c3826116fc565b6000806040838503121561174657600080fd5b61174f836116fc565b946020939093013593505050565b60008060006060848603121561177257600080fd5b61177b846116fc565b9250611789602085016116fc565b9150604084013590509250925092565b6020808252825182820181905260009190848201906040850190845b818110156117da5783516001600160a01b0316835292840192918401916001016117b5565b50909695505050505050565b6000806000604084860312156117fb57600080fd5b833567ffffffffffffffff8082111561181357600080fd5b818601915086601f83011261182757600080fd5b81358181111561183657600080fd5b87602082850101111561184857600080fd5b60209283019550935061185e91860190506116fc565b90509250925092565b60006020828403121561187957600080fd5b815180151581146113c357600080fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000600182016118c7576118c761189f565b5060010190565b818103818111156116825761168261189f565b634e487b7160e01b600052603160045260246000fdfea26469706673582212206558c2a5df5646d2d2d63e2246e2455a753ab34dee3cb4066a209f3bad2d7d6e64736f6c63430008110033",
}

// ModularComplianceABI is the input ABI used to generate the binding from.
// Deprecated: Use ModularComplianceMetaData.ABI instead.
var ModularComplianceABI = ModularComplianceMetaData.ABI

// ModularComplianceBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ModularComplianceMetaData.Bin instead.
var ModularComplianceBin = ModularComplianceMetaData.Bin

// DeployModularCompliance deploys a new Ethereum contract, binding an instance of ModularCompliance to it.
func DeployModularCompliance(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ModularCompliance, error) {
	parsed, err := ModularComplianceMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ModularComplianceBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ModularCompliance{ModularComplianceCaller: ModularComplianceCaller{contract: contract}, ModularComplianceTransactor: ModularComplianceTransactor{contract: contract}, ModularComplianceFilterer: ModularComplianceFilterer{contract: contract}}, nil
}

// ModularCompliance is an auto generated Go binding around an Ethereum contract.
type ModularCompliance struct {
	ModularComplianceCaller     // Read-only binding to the contract
	ModularComplianceTransactor // Write-only binding to the contract
	ModularComplianceFilterer   // Log filterer for contract events
}

// ModularComplianceCaller is an auto generated read-only Go binding around an Ethereum contract.
type ModularComplianceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ModularComplianceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ModularComplianceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ModularComplianceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ModularComplianceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ModularComplianceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ModularComplianceSession struct {
	Contract     *ModularCompliance // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ModularComplianceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ModularComplianceCallerSession struct {
	Contract *ModularComplianceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ModularComplianceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ModularComplianceTransactorSession struct {
	Contract     *ModularComplianceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ModularComplianceRaw is an auto generated low-level Go binding around an Ethereum contract.
type ModularComplianceRaw struct {
	Contract *ModularCompliance // Generic contract binding to access the raw methods on
}

// ModularComplianceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ModularComplianceCallerRaw struct {
	Contract *ModularComplianceCaller // Generic read-only contract binding to access the raw methods on
}

// ModularComplianceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ModularComplianceTransactorRaw struct {
	Contract *ModularComplianceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewModularCompliance creates a new instance of ModularCompliance, bound to a specific deployed contract.
func NewModularCompliance(address common.Address, backend bind.ContractBackend) (*ModularCompliance, error) {
	contract, err := bindModularCompliance(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ModularCompliance{ModularComplianceCaller: ModularComplianceCaller{contract: contract}, ModularComplianceTransactor: ModularComplianceTransactor{contract: contract}, ModularComplianceFilterer: ModularComplianceFilterer{contract: contract}}, nil
}

// NewModularComplianceCaller creates a new read-only instance of ModularCompliance, bound to a specific deployed contract.
func NewModularComplianceCaller(address common.Address, caller bind.ContractCaller) (*ModularComplianceCaller, error) {
	contract, err := bindModularCompliance(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ModularComplianceCaller{contract: contract}, nil
}

// NewModularComplianceTransactor creates a new write-only instance of ModularCompliance, bound to a specific deployed contract.
func NewModularComplianceTransactor(address common.Address, transactor bind.ContractTransactor) (*ModularComplianceTransactor, error) {
	contract, err := bindModularCompliance(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ModularComplianceTransactor{contract: contract}, nil
}

// NewModularComplianceFilterer creates a new log filterer instance of ModularCompliance, bound to a specific deployed contract.
func NewModularComplianceFilterer(address common.Address, filterer bind.ContractFilterer) (*ModularComplianceFilterer, error) {
	contract, err := bindModularCompliance(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ModularComplianceFilterer{contract: contract}, nil
}

// bindModularCompliance binds a generic wrapper to an already deployed contract.
func bindModularCompliance(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ModularComplianceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ModularCompliance *ModularComplianceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ModularCompliance.Contract.ModularComplianceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ModularCompliance *ModularComplianceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ModularCompliance.Contract.ModularComplianceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ModularCompliance *ModularComplianceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ModularCompliance.Contract.ModularComplianceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ModularCompliance *ModularComplianceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ModularCompliance.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ModularCompliance *ModularComplianceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ModularCompliance.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ModularCompliance *ModularComplianceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ModularCompliance.Contract.contract.Transact(opts, method, params...)
}

// CanTransfer is a free data retrieval call binding the contract method 0xe46638e6.
//
// Solidity: function canTransfer(address _from, address _to, uint256 _value) view returns(bool)
func (_ModularCompliance *ModularComplianceCaller) CanTransfer(opts *bind.CallOpts, _from common.Address, _to common.Address, _value *big.Int) (bool, error) {
	var out []interface{}
	err := _ModularCompliance.contract.Call(opts, &out, "canTransfer", _from, _to, _value)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanTransfer is a free data retrieval call binding the contract method 0xe46638e6.
//
// Solidity: function canTransfer(address _from, address _to, uint256 _value) view returns(bool)
func (_ModularCompliance *ModularComplianceSession) CanTransfer(_from common.Address, _to common.Address, _value *big.Int) (bool, error) {
	return _ModularCompliance.Contract.CanTransfer(&_ModularCompliance.CallOpts, _from, _to, _value)
}

// CanTransfer is a free data retrieval call binding the contract method 0xe46638e6.
//
// Solidity: function canTransfer(address _from, address _to, uint256 _value) view returns(bool)
func (_ModularCompliance *ModularComplianceCallerSession) CanTransfer(_from common.Address, _to common.Address, _value *big.Int) (bool, error) {
	return _ModularCompliance.Contract.CanTransfer(&_ModularCompliance.CallOpts, _from, _to, _value)
}

// GetModules is a free data retrieval call binding the contract method 0xb2494df3.
//
// Solidity: function getModules() view returns(address[])
func (_ModularCompliance *ModularComplianceCaller) GetModules(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _ModularCompliance.contract.Call(opts, &out, "getModules")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetModules is a free data retrieval call binding the contract method 0xb2494df3.
//
// Solidity: function getModules() view returns(address[])
func (_ModularCompliance *ModularComplianceSession) GetModules() ([]common.Address, error) {
	return _ModularCompliance.Contract.GetModules(&_ModularCompliance.CallOpts)
}

// GetModules is a free data retrieval call binding the contract method 0xb2494df3.
//
// Solidity: function getModules() view returns(address[])
func (_ModularCompliance *ModularComplianceCallerSession) GetModules() ([]common.Address, error) {
	return _ModularCompliance.Contract.GetModules(&_ModularCompliance.CallOpts)
}

// GetTokenBound is a free data retrieval call binding the contract method 0x6a3edf28.
//
// Solidity: function getTokenBound() view returns(address)
func (_ModularCompliance *ModularComplianceCaller) GetTokenBound(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ModularCompliance.contract.Call(opts, &out, "getTokenBound")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetTokenBound is a free data retrieval call binding the contract method 0x6a3edf28.
//
// Solidity: function getTokenBound() view returns(address)
func (_ModularCompliance *ModularComplianceSession) GetTokenBound() (common.Address, error) {
	return _ModularCompliance.Contract.GetTokenBound(&_ModularCompliance.CallOpts)
}

// GetTokenBound is a free data retrieval call binding the contract method 0x6a3edf28.
//
// Solidity: function getTokenBound() view returns(address)
func (_ModularCompliance *ModularComplianceCallerSession) GetTokenBound() (common.Address, error) {
	return _ModularCompliance.Contract.GetTokenBound(&_ModularCompliance.CallOpts)
}

// IsModuleBound is a free data retrieval call binding the contract method 0xa446d49f.
//
// Solidity: function isModuleBound(address _module) view returns(bool)
func (_ModularCompliance *ModularComplianceCaller) IsModuleBound(opts *bind.CallOpts, _module common.Address) (bool, error) {
	var out []interface{}
	err := _ModularCompliance.contract.Call(opts, &out, "isModuleBound", _module)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsModuleBound is a free data retrieval call binding the contract method 0xa446d49f.
//
// Solidity: function isModuleBound(address _module) view returns(bool)
func (_ModularCompliance *ModularComplianceSession) IsModuleBound(_module common.Address) (bool, error) {
	return _ModularCompliance.Contract.IsModuleBound(&_ModularCompliance.CallOpts, _module)
}

// IsModuleBound is a free data retrieval call binding the contract method 0xa446d49f.
//
// Solidity: function isModuleBound(address _module) view returns(bool)
func (_ModularCompliance *ModularComplianceCallerSession) IsModuleBound(_module common.Address) (bool, error) {
	return _ModularCompliance.Contract.IsModuleBound(&_ModularCompliance.CallOpts, _module)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ModularCompliance *ModularComplianceCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ModularCompliance.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ModularCompliance *ModularComplianceSession) Owner() (common.Address, error) {
	return _ModularCompliance.Contract.Owner(&_ModularCompliance.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ModularCompliance *ModularComplianceCallerSession) Owner() (common.Address, error) {
	return _ModularCompliance.Contract.Owner(&_ModularCompliance.CallOpts)
}

// AddModule is a paid mutator transaction binding the contract method 0x1ed86f19.
//
// Solidity: function addModule(address _module) returns()
func (_ModularCompliance *ModularComplianceTransactor) AddModule(opts *bind.TransactOpts, _module common.Address) (*types.Transaction, error) {
	return _ModularCompliance.contract.Transact(opts, "addModule", _module)
}

// AddModule is a paid mutator transaction binding the contract method 0x1ed86f19.
//
// Solidity: function addModule(address _module) returns()
func (_ModularCompliance *ModularComplianceSession) AddModule(_module common.Address) (*types.Transaction, error) {
	return _ModularCompliance.Contract.AddModule(&_ModularCompliance.TransactOpts, _module)
}

// AddModule is a paid mutator transaction binding the contract method 0x1ed86f19.
//
// Solidity: function addModule(address _module) returns()
func (_ModularCompliance *ModularComplianceTransactorSession) AddModule(_module common.Address) (*types.Transaction, error) {
	return _ModularCompliance.Contract.AddModule(&_ModularCompliance.TransactOpts, _module)
}

// BindToken is a paid mutator transaction binding the contract method 0x3ff5aa02.
//
// Solidity: function bindToken(address _token) returns()
func (_ModularCompliance *ModularComplianceTransactor) BindToken(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _ModularCompliance.contract.Transact(opts, "bindToken", _token)
}

// BindToken is a paid mutator transaction binding the contract method 0x3ff5aa02.
//
// Solidity: function bindToken(address _token) returns()
func (_ModularCompliance *ModularComplianceSession) BindToken(_token common.Address) (*types.Transaction, error) {
	return _ModularCompliance.Contract.BindToken(&_ModularCompliance.TransactOpts, _token)
}

// BindToken is a paid mutator transaction binding the contract method 0x3ff5aa02.
//
// Solidity: function bindToken(address _token) returns()
func (_ModularCompliance *ModularComplianceTransactorSession) BindToken(_token common.Address) (*types.Transaction, error) {
	return _ModularCompliance.Contract.BindToken(&_ModularCompliance.TransactOpts, _token)
}

// CallModuleFunction is a paid mutator transaction binding the contract method 0xefb22d33.
//
// Solidity: function callModuleFunction(bytes callData, address _module) returns()
func (_ModularCompliance *ModularComplianceTransactor) CallModuleFunction(opts *bind.TransactOpts, callData []byte, _module common.Address) (*types.Transaction, error) {
	return _ModularCompliance.contract.Transact(opts, "callModuleFunction", callData, _module)
}

// CallModuleFunction is a paid mutator transaction binding the contract method 0xefb22d33.
//
// Solidity: function callModuleFunction(bytes callData, address _module) returns()
func (_ModularCompliance *ModularComplianceSession) CallModuleFunction(callData []byte, _module common.Address) (*types.Transaction, error) {
	return _ModularCompliance.Contract.CallModuleFunction(&_ModularCompliance.TransactOpts, callData, _module)
}

// CallModuleFunction is a paid mutator transaction binding the contract method 0xefb22d33.
//
// Solidity: function callModuleFunction(bytes callData, address _module) returns()
func (_ModularCompliance *ModularComplianceTransactorSession) CallModuleFunction(callData []byte, _module common.Address) (*types.Transaction, error) {
	return _ModularCompliance.Contract.CallModuleFunction(&_ModularCompliance.TransactOpts, callData, _module)
}

// Created is a paid mutator transaction binding the contract method 0x5f8dead3.
//
// Solidity: function created(address _to, uint256 _value) returns()
func (_ModularCompliance *ModularComplianceTransactor) Created(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ModularCompliance.contract.Transact(opts, "created", _to, _value)
}

// Created is a paid mutator transaction binding the contract method 0x5f8dead3.
//
// Solidity: function created(address _to, uint256 _value) returns()
func (_ModularCompliance *ModularComplianceSession) Created(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ModularCompliance.Contract.Created(&_ModularCompliance.TransactOpts, _to, _value)
}

// Created is a paid mutator transaction binding the contract method 0x5f8dead3.
//
// Solidity: function created(address _to, uint256 _value) returns()
func (_ModularCompliance *ModularComplianceTransactorSession) Created(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ModularCompliance.Contract.Created(&_ModularCompliance.TransactOpts, _to, _value)
}

// Destroyed is a paid mutator transaction binding the contract method 0x8d2ea772.
//
// Solidity: function destroyed(address _from, uint256 _value) returns()
func (_ModularCompliance *ModularComplianceTransactor) Destroyed(opts *bind.TransactOpts, _from common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ModularCompliance.contract.Transact(opts, "destroyed", _from, _value)
}

// Destroyed is a paid mutator transaction binding the contract method 0x8d2ea772.
//
// Solidity: function destroyed(address _from, uint256 _value) returns()
func (_ModularCompliance *ModularComplianceSession) Destroyed(_from common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ModularCompliance.Contract.Destroyed(&_ModularCompliance.TransactOpts, _from, _value)
}

// Destroyed is a paid mutator transaction binding the contract method 0x8d2ea772.
//
// Solidity: function destroyed(address _from, uint256 _value) returns()
func (_ModularCompliance *ModularComplianceTransactorSession) Destroyed(_from common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ModularCompliance.Contract.Destroyed(&_ModularCompliance.TransactOpts, _from, _value)
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_ModularCompliance *ModularComplianceTransactor) Init(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ModularCompliance.contract.Transact(opts, "init")
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_ModularCompliance *ModularComplianceSession) Init() (*types.Transaction, error) {
	return _ModularCompliance.Contract.Init(&_ModularCompliance.TransactOpts)
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_ModularCompliance *ModularComplianceTransactorSession) Init() (*types.Transaction, error) {
	return _ModularCompliance.Contract.Init(&_ModularCompliance.TransactOpts)
}

// RemoveModule is a paid mutator transaction binding the contract method 0xa0632461.
//
// Solidity: function removeModule(address _module) returns()
func (_ModularCompliance *ModularComplianceTransactor) RemoveModule(opts *bind.TransactOpts, _module common.Address) (*types.Transaction, error) {
	return _ModularCompliance.contract.Transact(opts, "removeModule", _module)
}

// RemoveModule is a paid mutator transaction binding the contract method 0xa0632461.
//
// Solidity: function removeModule(address _module) returns()
func (_ModularCompliance *ModularComplianceSession) RemoveModule(_module common.Address) (*types.Transaction, error) {
	return _ModularCompliance.Contract.RemoveModule(&_ModularCompliance.TransactOpts, _module)
}

// RemoveModule is a paid mutator transaction binding the contract method 0xa0632461.
//
// Solidity: function removeModule(address _module) returns()
func (_ModularCompliance *ModularComplianceTransactorSession) RemoveModule(_module common.Address) (*types.Transaction, error) {
	return _ModularCompliance.Contract.RemoveModule(&_ModularCompliance.TransactOpts, _module)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ModularCompliance *ModularComplianceTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ModularCompliance.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ModularCompliance *ModularComplianceSession) RenounceOwnership() (*types.Transaction, error) {
	return _ModularCompliance.Contract.RenounceOwnership(&_ModularCompliance.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ModularCompliance *ModularComplianceTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ModularCompliance.Contract.RenounceOwnership(&_ModularCompliance.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ModularCompliance *ModularComplianceTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ModularCompliance.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ModularCompliance *ModularComplianceSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ModularCompliance.Contract.TransferOwnership(&_ModularCompliance.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ModularCompliance *ModularComplianceTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ModularCompliance.Contract.TransferOwnership(&_ModularCompliance.TransactOpts, newOwner)
}

// Transferred is a paid mutator transaction binding the contract method 0x8baf29b4.
//
// Solidity: function transferred(address _from, address _to, uint256 _value) returns()
func (_ModularCompliance *ModularComplianceTransactor) Transferred(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ModularCompliance.contract.Transact(opts, "transferred", _from, _to, _value)
}

// Transferred is a paid mutator transaction binding the contract method 0x8baf29b4.
//
// Solidity: function transferred(address _from, address _to, uint256 _value) returns()
func (_ModularCompliance *ModularComplianceSession) Transferred(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ModularCompliance.Contract.Transferred(&_ModularCompliance.TransactOpts, _from, _to, _value)
}

// Transferred is a paid mutator transaction binding the contract method 0x8baf29b4.
//
// Solidity: function transferred(address _from, address _to, uint256 _value) returns()
func (_ModularCompliance *ModularComplianceTransactorSession) Transferred(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ModularCompliance.Contract.Transferred(&_ModularCompliance.TransactOpts, _from, _to, _value)
}

// UnbindToken is a paid mutator transaction binding the contract method 0x40db3b50.
//
// Solidity: function unbindToken(address _token) returns()
func (_ModularCompliance *ModularComplianceTransactor) UnbindToken(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _ModularCompliance.contract.Transact(opts, "unbindToken", _token)
}

// UnbindToken is a paid mutator transaction binding the contract method 0x40db3b50.
//
// Solidity: function unbindToken(address _token) returns()
func (_ModularCompliance *ModularComplianceSession) UnbindToken(_token common.Address) (*types.Transaction, error) {
	return _ModularCompliance.Contract.UnbindToken(&_ModularCompliance.TransactOpts, _token)
}

// UnbindToken is a paid mutator transaction binding the contract method 0x40db3b50.
//
// Solidity: function unbindToken(address _token) returns()
func (_ModularCompliance *ModularComplianceTransactorSession) UnbindToken(_token common.Address) (*types.Transaction, error) {
	return _ModularCompliance.Contract.UnbindToken(&_ModularCompliance.TransactOpts, _token)
}

// ModularComplianceInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ModularCompliance contract.
type ModularComplianceInitializedIterator struct {
	Event *ModularComplianceInitialized // Event containing the contract specifics and raw log

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
func (it *ModularComplianceInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ModularComplianceInitialized)
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
		it.Event = new(ModularComplianceInitialized)
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
func (it *ModularComplianceInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ModularComplianceInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ModularComplianceInitialized represents a Initialized event raised by the ModularCompliance contract.
type ModularComplianceInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ModularCompliance *ModularComplianceFilterer) FilterInitialized(opts *bind.FilterOpts) (*ModularComplianceInitializedIterator, error) {

	logs, sub, err := _ModularCompliance.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ModularComplianceInitializedIterator{contract: _ModularCompliance.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ModularCompliance *ModularComplianceFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ModularComplianceInitialized) (event.Subscription, error) {

	logs, sub, err := _ModularCompliance.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ModularComplianceInitialized)
				if err := _ModularCompliance.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ModularCompliance *ModularComplianceFilterer) ParseInitialized(log types.Log) (*ModularComplianceInitialized, error) {
	event := new(ModularComplianceInitialized)
	if err := _ModularCompliance.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ModularComplianceModuleAddedIterator is returned from FilterModuleAdded and is used to iterate over the raw logs and unpacked data for ModuleAdded events raised by the ModularCompliance contract.
type ModularComplianceModuleAddedIterator struct {
	Event *ModularComplianceModuleAdded // Event containing the contract specifics and raw log

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
func (it *ModularComplianceModuleAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ModularComplianceModuleAdded)
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
		it.Event = new(ModularComplianceModuleAdded)
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
func (it *ModularComplianceModuleAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ModularComplianceModuleAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ModularComplianceModuleAdded represents a ModuleAdded event raised by the ModularCompliance contract.
type ModularComplianceModuleAdded struct {
	Module common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterModuleAdded is a free log retrieval operation binding the contract event 0xead6a006345da1073a106d5f32372d2d2204f46cb0b4bca8f5ebafcbbed12b8a.
//
// Solidity: event ModuleAdded(address indexed _module)
func (_ModularCompliance *ModularComplianceFilterer) FilterModuleAdded(opts *bind.FilterOpts, _module []common.Address) (*ModularComplianceModuleAddedIterator, error) {

	var _moduleRule []interface{}
	for _, _moduleItem := range _module {
		_moduleRule = append(_moduleRule, _moduleItem)
	}

	logs, sub, err := _ModularCompliance.contract.FilterLogs(opts, "ModuleAdded", _moduleRule)
	if err != nil {
		return nil, err
	}
	return &ModularComplianceModuleAddedIterator{contract: _ModularCompliance.contract, event: "ModuleAdded", logs: logs, sub: sub}, nil
}

// WatchModuleAdded is a free log subscription operation binding the contract event 0xead6a006345da1073a106d5f32372d2d2204f46cb0b4bca8f5ebafcbbed12b8a.
//
// Solidity: event ModuleAdded(address indexed _module)
func (_ModularCompliance *ModularComplianceFilterer) WatchModuleAdded(opts *bind.WatchOpts, sink chan<- *ModularComplianceModuleAdded, _module []common.Address) (event.Subscription, error) {

	var _moduleRule []interface{}
	for _, _moduleItem := range _module {
		_moduleRule = append(_moduleRule, _moduleItem)
	}

	logs, sub, err := _ModularCompliance.contract.WatchLogs(opts, "ModuleAdded", _moduleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ModularComplianceModuleAdded)
				if err := _ModularCompliance.contract.UnpackLog(event, "ModuleAdded", log); err != nil {
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

// ParseModuleAdded is a log parse operation binding the contract event 0xead6a006345da1073a106d5f32372d2d2204f46cb0b4bca8f5ebafcbbed12b8a.
//
// Solidity: event ModuleAdded(address indexed _module)
func (_ModularCompliance *ModularComplianceFilterer) ParseModuleAdded(log types.Log) (*ModularComplianceModuleAdded, error) {
	event := new(ModularComplianceModuleAdded)
	if err := _ModularCompliance.contract.UnpackLog(event, "ModuleAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ModularComplianceModuleInteractionIterator is returned from FilterModuleInteraction and is used to iterate over the raw logs and unpacked data for ModuleInteraction events raised by the ModularCompliance contract.
type ModularComplianceModuleInteractionIterator struct {
	Event *ModularComplianceModuleInteraction // Event containing the contract specifics and raw log

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
func (it *ModularComplianceModuleInteractionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ModularComplianceModuleInteraction)
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
		it.Event = new(ModularComplianceModuleInteraction)
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
func (it *ModularComplianceModuleInteractionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ModularComplianceModuleInteractionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ModularComplianceModuleInteraction represents a ModuleInteraction event raised by the ModularCompliance contract.
type ModularComplianceModuleInteraction struct {
	Target   common.Address
	Selector [4]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterModuleInteraction is a free log retrieval operation binding the contract event 0x20d79de70adcc6e9353d8a9a5646b46dc352710d0a310b1ad1f67faeca7ef891.
//
// Solidity: event ModuleInteraction(address indexed target, bytes4 selector)
func (_ModularCompliance *ModularComplianceFilterer) FilterModuleInteraction(opts *bind.FilterOpts, target []common.Address) (*ModularComplianceModuleInteractionIterator, error) {

	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}

	logs, sub, err := _ModularCompliance.contract.FilterLogs(opts, "ModuleInteraction", targetRule)
	if err != nil {
		return nil, err
	}
	return &ModularComplianceModuleInteractionIterator{contract: _ModularCompliance.contract, event: "ModuleInteraction", logs: logs, sub: sub}, nil
}

// WatchModuleInteraction is a free log subscription operation binding the contract event 0x20d79de70adcc6e9353d8a9a5646b46dc352710d0a310b1ad1f67faeca7ef891.
//
// Solidity: event ModuleInteraction(address indexed target, bytes4 selector)
func (_ModularCompliance *ModularComplianceFilterer) WatchModuleInteraction(opts *bind.WatchOpts, sink chan<- *ModularComplianceModuleInteraction, target []common.Address) (event.Subscription, error) {

	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}

	logs, sub, err := _ModularCompliance.contract.WatchLogs(opts, "ModuleInteraction", targetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ModularComplianceModuleInteraction)
				if err := _ModularCompliance.contract.UnpackLog(event, "ModuleInteraction", log); err != nil {
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

// ParseModuleInteraction is a log parse operation binding the contract event 0x20d79de70adcc6e9353d8a9a5646b46dc352710d0a310b1ad1f67faeca7ef891.
//
// Solidity: event ModuleInteraction(address indexed target, bytes4 selector)
func (_ModularCompliance *ModularComplianceFilterer) ParseModuleInteraction(log types.Log) (*ModularComplianceModuleInteraction, error) {
	event := new(ModularComplianceModuleInteraction)
	if err := _ModularCompliance.contract.UnpackLog(event, "ModuleInteraction", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ModularComplianceModuleRemovedIterator is returned from FilterModuleRemoved and is used to iterate over the raw logs and unpacked data for ModuleRemoved events raised by the ModularCompliance contract.
type ModularComplianceModuleRemovedIterator struct {
	Event *ModularComplianceModuleRemoved // Event containing the contract specifics and raw log

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
func (it *ModularComplianceModuleRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ModularComplianceModuleRemoved)
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
		it.Event = new(ModularComplianceModuleRemoved)
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
func (it *ModularComplianceModuleRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ModularComplianceModuleRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ModularComplianceModuleRemoved represents a ModuleRemoved event raised by the ModularCompliance contract.
type ModularComplianceModuleRemoved struct {
	Module common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterModuleRemoved is a free log retrieval operation binding the contract event 0x0a1ee69f55c33d8467c69ca59ce2007a737a88603d75392972520bf67cb513b8.
//
// Solidity: event ModuleRemoved(address indexed _module)
func (_ModularCompliance *ModularComplianceFilterer) FilterModuleRemoved(opts *bind.FilterOpts, _module []common.Address) (*ModularComplianceModuleRemovedIterator, error) {

	var _moduleRule []interface{}
	for _, _moduleItem := range _module {
		_moduleRule = append(_moduleRule, _moduleItem)
	}

	logs, sub, err := _ModularCompliance.contract.FilterLogs(opts, "ModuleRemoved", _moduleRule)
	if err != nil {
		return nil, err
	}
	return &ModularComplianceModuleRemovedIterator{contract: _ModularCompliance.contract, event: "ModuleRemoved", logs: logs, sub: sub}, nil
}

// WatchModuleRemoved is a free log subscription operation binding the contract event 0x0a1ee69f55c33d8467c69ca59ce2007a737a88603d75392972520bf67cb513b8.
//
// Solidity: event ModuleRemoved(address indexed _module)
func (_ModularCompliance *ModularComplianceFilterer) WatchModuleRemoved(opts *bind.WatchOpts, sink chan<- *ModularComplianceModuleRemoved, _module []common.Address) (event.Subscription, error) {

	var _moduleRule []interface{}
	for _, _moduleItem := range _module {
		_moduleRule = append(_moduleRule, _moduleItem)
	}

	logs, sub, err := _ModularCompliance.contract.WatchLogs(opts, "ModuleRemoved", _moduleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ModularComplianceModuleRemoved)
				if err := _ModularCompliance.contract.UnpackLog(event, "ModuleRemoved", log); err != nil {
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

// ParseModuleRemoved is a log parse operation binding the contract event 0x0a1ee69f55c33d8467c69ca59ce2007a737a88603d75392972520bf67cb513b8.
//
// Solidity: event ModuleRemoved(address indexed _module)
func (_ModularCompliance *ModularComplianceFilterer) ParseModuleRemoved(log types.Log) (*ModularComplianceModuleRemoved, error) {
	event := new(ModularComplianceModuleRemoved)
	if err := _ModularCompliance.contract.UnpackLog(event, "ModuleRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ModularComplianceOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ModularCompliance contract.
type ModularComplianceOwnershipTransferredIterator struct {
	Event *ModularComplianceOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ModularComplianceOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ModularComplianceOwnershipTransferred)
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
		it.Event = new(ModularComplianceOwnershipTransferred)
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
func (it *ModularComplianceOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ModularComplianceOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ModularComplianceOwnershipTransferred represents a OwnershipTransferred event raised by the ModularCompliance contract.
type ModularComplianceOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ModularCompliance *ModularComplianceFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ModularComplianceOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ModularCompliance.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ModularComplianceOwnershipTransferredIterator{contract: _ModularCompliance.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ModularCompliance *ModularComplianceFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ModularComplianceOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ModularCompliance.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ModularComplianceOwnershipTransferred)
				if err := _ModularCompliance.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ModularCompliance *ModularComplianceFilterer) ParseOwnershipTransferred(log types.Log) (*ModularComplianceOwnershipTransferred, error) {
	event := new(ModularComplianceOwnershipTransferred)
	if err := _ModularCompliance.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ModularComplianceTokenBoundIterator is returned from FilterTokenBound and is used to iterate over the raw logs and unpacked data for TokenBound events raised by the ModularCompliance contract.
type ModularComplianceTokenBoundIterator struct {
	Event *ModularComplianceTokenBound // Event containing the contract specifics and raw log

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
func (it *ModularComplianceTokenBoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ModularComplianceTokenBound)
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
		it.Event = new(ModularComplianceTokenBound)
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
func (it *ModularComplianceTokenBoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ModularComplianceTokenBoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ModularComplianceTokenBound represents a TokenBound event raised by the ModularCompliance contract.
type ModularComplianceTokenBound struct {
	Token common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTokenBound is a free log retrieval operation binding the contract event 0x2de35142b19ed5a07796cf30791959c592018f70b1d2d7c460eef8ffe713692b.
//
// Solidity: event TokenBound(address _token)
func (_ModularCompliance *ModularComplianceFilterer) FilterTokenBound(opts *bind.FilterOpts) (*ModularComplianceTokenBoundIterator, error) {

	logs, sub, err := _ModularCompliance.contract.FilterLogs(opts, "TokenBound")
	if err != nil {
		return nil, err
	}
	return &ModularComplianceTokenBoundIterator{contract: _ModularCompliance.contract, event: "TokenBound", logs: logs, sub: sub}, nil
}

// WatchTokenBound is a free log subscription operation binding the contract event 0x2de35142b19ed5a07796cf30791959c592018f70b1d2d7c460eef8ffe713692b.
//
// Solidity: event TokenBound(address _token)
func (_ModularCompliance *ModularComplianceFilterer) WatchTokenBound(opts *bind.WatchOpts, sink chan<- *ModularComplianceTokenBound) (event.Subscription, error) {

	logs, sub, err := _ModularCompliance.contract.WatchLogs(opts, "TokenBound")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ModularComplianceTokenBound)
				if err := _ModularCompliance.contract.UnpackLog(event, "TokenBound", log); err != nil {
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

// ParseTokenBound is a log parse operation binding the contract event 0x2de35142b19ed5a07796cf30791959c592018f70b1d2d7c460eef8ffe713692b.
//
// Solidity: event TokenBound(address _token)
func (_ModularCompliance *ModularComplianceFilterer) ParseTokenBound(log types.Log) (*ModularComplianceTokenBound, error) {
	event := new(ModularComplianceTokenBound)
	if err := _ModularCompliance.contract.UnpackLog(event, "TokenBound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ModularComplianceTokenUnboundIterator is returned from FilterTokenUnbound and is used to iterate over the raw logs and unpacked data for TokenUnbound events raised by the ModularCompliance contract.
type ModularComplianceTokenUnboundIterator struct {
	Event *ModularComplianceTokenUnbound // Event containing the contract specifics and raw log

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
func (it *ModularComplianceTokenUnboundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ModularComplianceTokenUnbound)
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
		it.Event = new(ModularComplianceTokenUnbound)
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
func (it *ModularComplianceTokenUnboundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ModularComplianceTokenUnboundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ModularComplianceTokenUnbound represents a TokenUnbound event raised by the ModularCompliance contract.
type ModularComplianceTokenUnbound struct {
	Token common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTokenUnbound is a free log retrieval operation binding the contract event 0x28a4ca7134a3b3f9aff286e79ad3daadb4a06d1b43d037a3a98bdc074edd9b7a.
//
// Solidity: event TokenUnbound(address _token)
func (_ModularCompliance *ModularComplianceFilterer) FilterTokenUnbound(opts *bind.FilterOpts) (*ModularComplianceTokenUnboundIterator, error) {

	logs, sub, err := _ModularCompliance.contract.FilterLogs(opts, "TokenUnbound")
	if err != nil {
		return nil, err
	}
	return &ModularComplianceTokenUnboundIterator{contract: _ModularCompliance.contract, event: "TokenUnbound", logs: logs, sub: sub}, nil
}

// WatchTokenUnbound is a free log subscription operation binding the contract event 0x28a4ca7134a3b3f9aff286e79ad3daadb4a06d1b43d037a3a98bdc074edd9b7a.
//
// Solidity: event TokenUnbound(address _token)
func (_ModularCompliance *ModularComplianceFilterer) WatchTokenUnbound(opts *bind.WatchOpts, sink chan<- *ModularComplianceTokenUnbound) (event.Subscription, error) {

	logs, sub, err := _ModularCompliance.contract.WatchLogs(opts, "TokenUnbound")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ModularComplianceTokenUnbound)
				if err := _ModularCompliance.contract.UnpackLog(event, "TokenUnbound", log); err != nil {
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

// ParseTokenUnbound is a log parse operation binding the contract event 0x28a4ca7134a3b3f9aff286e79ad3daadb4a06d1b43d037a3a98bdc074edd9b7a.
//
// Solidity: event TokenUnbound(address _token)
func (_ModularCompliance *ModularComplianceFilterer) ParseTokenUnbound(log types.Log) (*ModularComplianceTokenUnbound, error) {
	event := new(ModularComplianceTokenUnbound)
	if err := _ModularCompliance.contract.UnpackLog(event, "TokenUnbound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
