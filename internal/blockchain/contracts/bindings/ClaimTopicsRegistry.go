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

// ClaimTopicsRegistryMetaData contains all meta data concerning the ClaimTopicsRegistry contract.
var ClaimTopicsRegistryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"claimTopic\",\"type\":\"uint256\"}],\"name\":\"ClaimTopicAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"claimTopic\",\"type\":\"uint256\"}],\"name\":\"ClaimTopicRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_claimTopic\",\"type\":\"uint256\"}],\"name\":\"addClaimTopic\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getClaimTopics\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_claimTopic\",\"type\":\"uint256\"}],\"name\":\"removeClaimTopic\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610859806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063c7b225511161005b578063c7b22551146100bf578063df09d604146100d2578063e1c7392a146100e7578063f2fde38b146100ef57600080fd5b80630829784614610082578063715018a6146100975780638da5cb5b1461009f575b600080fd5b610095610090366004610722565b610102565b005b6100956101f3565b6033546040516001600160a01b0390911681526020015b60405180910390f35b6100956100cd366004610722565b610207565b6100da610377565b6040516100b6919061073b565b6100956103cf565b6100956100fd36600461077f565b6104ef565b61010a61057c565b60655460005b818110156101ee57826065828154811061012c5761012c6107af565b9060005260206000200154036101dc5760656101496001846107db565b81548110610159576101596107af565b906000526020600020015460658281548110610177576101776107af565b6000918252602090912001556065805480610194576101946107f4565b60019003818190600052602060002001600090559055827f0b1381093c776453c1bbe54fd68be1b235c65db61d099cb50d194b2991e0eec560405160405180910390a2505050565b806101e68161080a565b915050610110565b505050565b6101fb61057c565b61020560006105d6565b565b61020f61057c565b606554600f811061028d5760405162461bcd60e51b815260206004820152602260248201527f63616e6e6f742072657175697265206d6f7265207468616e20313520746f706960448201527f637300000000000000000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b60005b818110156103165782606582815481106102ac576102ac6107af565b9060005260206000200154036103045760405162461bcd60e51b815260206004820152601960248201527f636c61696d546f70696320616c726561647920657869737473000000000000006044820152606401610284565b8061030e8161080a565b915050610290565b506065805460018101825560009182527f8ff97419363ffd7000167f130ef7168fbea05faf9251824ca5043f113cc6a7c70183905560405183917f01c928b7f7ade2949e92366aa9454dbef3a416b731cf6ec786ba9595bbd814d691a25050565b606060658054806020026020016040519081016040528092919081815260200182805480156103c557602002820191906000526020600020905b8154815260200190600101908083116103b1575b5050505050905090565b600054610100900460ff16158080156103ef5750600054600160ff909116105b806104095750303b158015610409575060005460ff166001145b61047b5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a65640000000000000000000000000000000000006064820152608401610284565b6000805460ff19166001179055801561049e576000805461ff0019166101001790555b6104a6610640565b80156104ec576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50565b6104f761057c565b6001600160a01b0381166105735760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610284565b6104ec816105d6565b6033546001600160a01b031633146102055760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610284565b603380546001600160a01b038381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166106ab5760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b6064820152608401610284565b610205600054610100900460ff166107195760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b6064820152608401610284565b610205336105d6565b60006020828403121561073457600080fd5b5035919050565b6020808252825182820181905260009190848201906040850190845b8181101561077357835183529284019291840191600101610757565b50909695505050505050565b60006020828403121561079157600080fd5b81356001600160a01b03811681146107a857600080fd5b9392505050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b818103818111156107ee576107ee6107c5565b92915050565b634e487b7160e01b600052603160045260246000fd5b60006001820161081c5761081c6107c5565b506001019056fea26469706673582212202c57d4f18f25c5ab1213c615d871f30f62457e6664e41e392338b914b819131764736f6c63430008110033",
}

// ClaimTopicsRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ClaimTopicsRegistryMetaData.ABI instead.
var ClaimTopicsRegistryABI = ClaimTopicsRegistryMetaData.ABI

// ClaimTopicsRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ClaimTopicsRegistryMetaData.Bin instead.
var ClaimTopicsRegistryBin = ClaimTopicsRegistryMetaData.Bin

// DeployClaimTopicsRegistry deploys a new Ethereum contract, binding an instance of ClaimTopicsRegistry to it.
func DeployClaimTopicsRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ClaimTopicsRegistry, error) {
	parsed, err := ClaimTopicsRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ClaimTopicsRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ClaimTopicsRegistry{ClaimTopicsRegistryCaller: ClaimTopicsRegistryCaller{contract: contract}, ClaimTopicsRegistryTransactor: ClaimTopicsRegistryTransactor{contract: contract}, ClaimTopicsRegistryFilterer: ClaimTopicsRegistryFilterer{contract: contract}}, nil
}

// ClaimTopicsRegistry is an auto generated Go binding around an Ethereum contract.
type ClaimTopicsRegistry struct {
	ClaimTopicsRegistryCaller     // Read-only binding to the contract
	ClaimTopicsRegistryTransactor // Write-only binding to the contract
	ClaimTopicsRegistryFilterer   // Log filterer for contract events
}

// ClaimTopicsRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ClaimTopicsRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClaimTopicsRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ClaimTopicsRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClaimTopicsRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ClaimTopicsRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClaimTopicsRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ClaimTopicsRegistrySession struct {
	Contract     *ClaimTopicsRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ClaimTopicsRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ClaimTopicsRegistryCallerSession struct {
	Contract *ClaimTopicsRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// ClaimTopicsRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ClaimTopicsRegistryTransactorSession struct {
	Contract     *ClaimTopicsRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ClaimTopicsRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ClaimTopicsRegistryRaw struct {
	Contract *ClaimTopicsRegistry // Generic contract binding to access the raw methods on
}

// ClaimTopicsRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ClaimTopicsRegistryCallerRaw struct {
	Contract *ClaimTopicsRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ClaimTopicsRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ClaimTopicsRegistryTransactorRaw struct {
	Contract *ClaimTopicsRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewClaimTopicsRegistry creates a new instance of ClaimTopicsRegistry, bound to a specific deployed contract.
func NewClaimTopicsRegistry(address common.Address, backend bind.ContractBackend) (*ClaimTopicsRegistry, error) {
	contract, err := bindClaimTopicsRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ClaimTopicsRegistry{ClaimTopicsRegistryCaller: ClaimTopicsRegistryCaller{contract: contract}, ClaimTopicsRegistryTransactor: ClaimTopicsRegistryTransactor{contract: contract}, ClaimTopicsRegistryFilterer: ClaimTopicsRegistryFilterer{contract: contract}}, nil
}

// NewClaimTopicsRegistryCaller creates a new read-only instance of ClaimTopicsRegistry, bound to a specific deployed contract.
func NewClaimTopicsRegistryCaller(address common.Address, caller bind.ContractCaller) (*ClaimTopicsRegistryCaller, error) {
	contract, err := bindClaimTopicsRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ClaimTopicsRegistryCaller{contract: contract}, nil
}

// NewClaimTopicsRegistryTransactor creates a new write-only instance of ClaimTopicsRegistry, bound to a specific deployed contract.
func NewClaimTopicsRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ClaimTopicsRegistryTransactor, error) {
	contract, err := bindClaimTopicsRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ClaimTopicsRegistryTransactor{contract: contract}, nil
}

// NewClaimTopicsRegistryFilterer creates a new log filterer instance of ClaimTopicsRegistry, bound to a specific deployed contract.
func NewClaimTopicsRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ClaimTopicsRegistryFilterer, error) {
	contract, err := bindClaimTopicsRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ClaimTopicsRegistryFilterer{contract: contract}, nil
}

// bindClaimTopicsRegistry binds a generic wrapper to an already deployed contract.
func bindClaimTopicsRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ClaimTopicsRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ClaimTopicsRegistry *ClaimTopicsRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ClaimTopicsRegistry.Contract.ClaimTopicsRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ClaimTopicsRegistry *ClaimTopicsRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClaimTopicsRegistry.Contract.ClaimTopicsRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ClaimTopicsRegistry *ClaimTopicsRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ClaimTopicsRegistry.Contract.ClaimTopicsRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ClaimTopicsRegistry *ClaimTopicsRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ClaimTopicsRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ClaimTopicsRegistry *ClaimTopicsRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClaimTopicsRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ClaimTopicsRegistry *ClaimTopicsRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ClaimTopicsRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetClaimTopics is a free data retrieval call binding the contract method 0xdf09d604.
//
// Solidity: function getClaimTopics() view returns(uint256[])
func (_ClaimTopicsRegistry *ClaimTopicsRegistryCaller) GetClaimTopics(opts *bind.CallOpts) ([]*big.Int, error) {
	var out []interface{}
	err := _ClaimTopicsRegistry.contract.Call(opts, &out, "getClaimTopics")

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetClaimTopics is a free data retrieval call binding the contract method 0xdf09d604.
//
// Solidity: function getClaimTopics() view returns(uint256[])
func (_ClaimTopicsRegistry *ClaimTopicsRegistrySession) GetClaimTopics() ([]*big.Int, error) {
	return _ClaimTopicsRegistry.Contract.GetClaimTopics(&_ClaimTopicsRegistry.CallOpts)
}

// GetClaimTopics is a free data retrieval call binding the contract method 0xdf09d604.
//
// Solidity: function getClaimTopics() view returns(uint256[])
func (_ClaimTopicsRegistry *ClaimTopicsRegistryCallerSession) GetClaimTopics() ([]*big.Int, error) {
	return _ClaimTopicsRegistry.Contract.GetClaimTopics(&_ClaimTopicsRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ClaimTopicsRegistry *ClaimTopicsRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ClaimTopicsRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ClaimTopicsRegistry *ClaimTopicsRegistrySession) Owner() (common.Address, error) {
	return _ClaimTopicsRegistry.Contract.Owner(&_ClaimTopicsRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ClaimTopicsRegistry *ClaimTopicsRegistryCallerSession) Owner() (common.Address, error) {
	return _ClaimTopicsRegistry.Contract.Owner(&_ClaimTopicsRegistry.CallOpts)
}

// AddClaimTopic is a paid mutator transaction binding the contract method 0xc7b22551.
//
// Solidity: function addClaimTopic(uint256 _claimTopic) returns()
func (_ClaimTopicsRegistry *ClaimTopicsRegistryTransactor) AddClaimTopic(opts *bind.TransactOpts, _claimTopic *big.Int) (*types.Transaction, error) {
	return _ClaimTopicsRegistry.contract.Transact(opts, "addClaimTopic", _claimTopic)
}

// AddClaimTopic is a paid mutator transaction binding the contract method 0xc7b22551.
//
// Solidity: function addClaimTopic(uint256 _claimTopic) returns()
func (_ClaimTopicsRegistry *ClaimTopicsRegistrySession) AddClaimTopic(_claimTopic *big.Int) (*types.Transaction, error) {
	return _ClaimTopicsRegistry.Contract.AddClaimTopic(&_ClaimTopicsRegistry.TransactOpts, _claimTopic)
}

// AddClaimTopic is a paid mutator transaction binding the contract method 0xc7b22551.
//
// Solidity: function addClaimTopic(uint256 _claimTopic) returns()
func (_ClaimTopicsRegistry *ClaimTopicsRegistryTransactorSession) AddClaimTopic(_claimTopic *big.Int) (*types.Transaction, error) {
	return _ClaimTopicsRegistry.Contract.AddClaimTopic(&_ClaimTopicsRegistry.TransactOpts, _claimTopic)
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_ClaimTopicsRegistry *ClaimTopicsRegistryTransactor) Init(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClaimTopicsRegistry.contract.Transact(opts, "init")
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_ClaimTopicsRegistry *ClaimTopicsRegistrySession) Init() (*types.Transaction, error) {
	return _ClaimTopicsRegistry.Contract.Init(&_ClaimTopicsRegistry.TransactOpts)
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_ClaimTopicsRegistry *ClaimTopicsRegistryTransactorSession) Init() (*types.Transaction, error) {
	return _ClaimTopicsRegistry.Contract.Init(&_ClaimTopicsRegistry.TransactOpts)
}

// RemoveClaimTopic is a paid mutator transaction binding the contract method 0x08297846.
//
// Solidity: function removeClaimTopic(uint256 _claimTopic) returns()
func (_ClaimTopicsRegistry *ClaimTopicsRegistryTransactor) RemoveClaimTopic(opts *bind.TransactOpts, _claimTopic *big.Int) (*types.Transaction, error) {
	return _ClaimTopicsRegistry.contract.Transact(opts, "removeClaimTopic", _claimTopic)
}

// RemoveClaimTopic is a paid mutator transaction binding the contract method 0x08297846.
//
// Solidity: function removeClaimTopic(uint256 _claimTopic) returns()
func (_ClaimTopicsRegistry *ClaimTopicsRegistrySession) RemoveClaimTopic(_claimTopic *big.Int) (*types.Transaction, error) {
	return _ClaimTopicsRegistry.Contract.RemoveClaimTopic(&_ClaimTopicsRegistry.TransactOpts, _claimTopic)
}

// RemoveClaimTopic is a paid mutator transaction binding the contract method 0x08297846.
//
// Solidity: function removeClaimTopic(uint256 _claimTopic) returns()
func (_ClaimTopicsRegistry *ClaimTopicsRegistryTransactorSession) RemoveClaimTopic(_claimTopic *big.Int) (*types.Transaction, error) {
	return _ClaimTopicsRegistry.Contract.RemoveClaimTopic(&_ClaimTopicsRegistry.TransactOpts, _claimTopic)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ClaimTopicsRegistry *ClaimTopicsRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClaimTopicsRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ClaimTopicsRegistry *ClaimTopicsRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _ClaimTopicsRegistry.Contract.RenounceOwnership(&_ClaimTopicsRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ClaimTopicsRegistry *ClaimTopicsRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ClaimTopicsRegistry.Contract.RenounceOwnership(&_ClaimTopicsRegistry.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ClaimTopicsRegistry *ClaimTopicsRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ClaimTopicsRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ClaimTopicsRegistry *ClaimTopicsRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ClaimTopicsRegistry.Contract.TransferOwnership(&_ClaimTopicsRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ClaimTopicsRegistry *ClaimTopicsRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ClaimTopicsRegistry.Contract.TransferOwnership(&_ClaimTopicsRegistry.TransactOpts, newOwner)
}

// ClaimTopicsRegistryClaimTopicAddedIterator is returned from FilterClaimTopicAdded and is used to iterate over the raw logs and unpacked data for ClaimTopicAdded events raised by the ClaimTopicsRegistry contract.
type ClaimTopicsRegistryClaimTopicAddedIterator struct {
	Event *ClaimTopicsRegistryClaimTopicAdded // Event containing the contract specifics and raw log

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
func (it *ClaimTopicsRegistryClaimTopicAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClaimTopicsRegistryClaimTopicAdded)
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
		it.Event = new(ClaimTopicsRegistryClaimTopicAdded)
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
func (it *ClaimTopicsRegistryClaimTopicAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClaimTopicsRegistryClaimTopicAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClaimTopicsRegistryClaimTopicAdded represents a ClaimTopicAdded event raised by the ClaimTopicsRegistry contract.
type ClaimTopicsRegistryClaimTopicAdded struct {
	ClaimTopic *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterClaimTopicAdded is a free log retrieval operation binding the contract event 0x01c928b7f7ade2949e92366aa9454dbef3a416b731cf6ec786ba9595bbd814d6.
//
// Solidity: event ClaimTopicAdded(uint256 indexed claimTopic)
func (_ClaimTopicsRegistry *ClaimTopicsRegistryFilterer) FilterClaimTopicAdded(opts *bind.FilterOpts, claimTopic []*big.Int) (*ClaimTopicsRegistryClaimTopicAddedIterator, error) {

	var claimTopicRule []interface{}
	for _, claimTopicItem := range claimTopic {
		claimTopicRule = append(claimTopicRule, claimTopicItem)
	}

	logs, sub, err := _ClaimTopicsRegistry.contract.FilterLogs(opts, "ClaimTopicAdded", claimTopicRule)
	if err != nil {
		return nil, err
	}
	return &ClaimTopicsRegistryClaimTopicAddedIterator{contract: _ClaimTopicsRegistry.contract, event: "ClaimTopicAdded", logs: logs, sub: sub}, nil
}

// WatchClaimTopicAdded is a free log subscription operation binding the contract event 0x01c928b7f7ade2949e92366aa9454dbef3a416b731cf6ec786ba9595bbd814d6.
//
// Solidity: event ClaimTopicAdded(uint256 indexed claimTopic)
func (_ClaimTopicsRegistry *ClaimTopicsRegistryFilterer) WatchClaimTopicAdded(opts *bind.WatchOpts, sink chan<- *ClaimTopicsRegistryClaimTopicAdded, claimTopic []*big.Int) (event.Subscription, error) {

	var claimTopicRule []interface{}
	for _, claimTopicItem := range claimTopic {
		claimTopicRule = append(claimTopicRule, claimTopicItem)
	}

	logs, sub, err := _ClaimTopicsRegistry.contract.WatchLogs(opts, "ClaimTopicAdded", claimTopicRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClaimTopicsRegistryClaimTopicAdded)
				if err := _ClaimTopicsRegistry.contract.UnpackLog(event, "ClaimTopicAdded", log); err != nil {
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

// ParseClaimTopicAdded is a log parse operation binding the contract event 0x01c928b7f7ade2949e92366aa9454dbef3a416b731cf6ec786ba9595bbd814d6.
//
// Solidity: event ClaimTopicAdded(uint256 indexed claimTopic)
func (_ClaimTopicsRegistry *ClaimTopicsRegistryFilterer) ParseClaimTopicAdded(log types.Log) (*ClaimTopicsRegistryClaimTopicAdded, error) {
	event := new(ClaimTopicsRegistryClaimTopicAdded)
	if err := _ClaimTopicsRegistry.contract.UnpackLog(event, "ClaimTopicAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClaimTopicsRegistryClaimTopicRemovedIterator is returned from FilterClaimTopicRemoved and is used to iterate over the raw logs and unpacked data for ClaimTopicRemoved events raised by the ClaimTopicsRegistry contract.
type ClaimTopicsRegistryClaimTopicRemovedIterator struct {
	Event *ClaimTopicsRegistryClaimTopicRemoved // Event containing the contract specifics and raw log

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
func (it *ClaimTopicsRegistryClaimTopicRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClaimTopicsRegistryClaimTopicRemoved)
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
		it.Event = new(ClaimTopicsRegistryClaimTopicRemoved)
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
func (it *ClaimTopicsRegistryClaimTopicRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClaimTopicsRegistryClaimTopicRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClaimTopicsRegistryClaimTopicRemoved represents a ClaimTopicRemoved event raised by the ClaimTopicsRegistry contract.
type ClaimTopicsRegistryClaimTopicRemoved struct {
	ClaimTopic *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterClaimTopicRemoved is a free log retrieval operation binding the contract event 0x0b1381093c776453c1bbe54fd68be1b235c65db61d099cb50d194b2991e0eec5.
//
// Solidity: event ClaimTopicRemoved(uint256 indexed claimTopic)
func (_ClaimTopicsRegistry *ClaimTopicsRegistryFilterer) FilterClaimTopicRemoved(opts *bind.FilterOpts, claimTopic []*big.Int) (*ClaimTopicsRegistryClaimTopicRemovedIterator, error) {

	var claimTopicRule []interface{}
	for _, claimTopicItem := range claimTopic {
		claimTopicRule = append(claimTopicRule, claimTopicItem)
	}

	logs, sub, err := _ClaimTopicsRegistry.contract.FilterLogs(opts, "ClaimTopicRemoved", claimTopicRule)
	if err != nil {
		return nil, err
	}
	return &ClaimTopicsRegistryClaimTopicRemovedIterator{contract: _ClaimTopicsRegistry.contract, event: "ClaimTopicRemoved", logs: logs, sub: sub}, nil
}

// WatchClaimTopicRemoved is a free log subscription operation binding the contract event 0x0b1381093c776453c1bbe54fd68be1b235c65db61d099cb50d194b2991e0eec5.
//
// Solidity: event ClaimTopicRemoved(uint256 indexed claimTopic)
func (_ClaimTopicsRegistry *ClaimTopicsRegistryFilterer) WatchClaimTopicRemoved(opts *bind.WatchOpts, sink chan<- *ClaimTopicsRegistryClaimTopicRemoved, claimTopic []*big.Int) (event.Subscription, error) {

	var claimTopicRule []interface{}
	for _, claimTopicItem := range claimTopic {
		claimTopicRule = append(claimTopicRule, claimTopicItem)
	}

	logs, sub, err := _ClaimTopicsRegistry.contract.WatchLogs(opts, "ClaimTopicRemoved", claimTopicRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClaimTopicsRegistryClaimTopicRemoved)
				if err := _ClaimTopicsRegistry.contract.UnpackLog(event, "ClaimTopicRemoved", log); err != nil {
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

// ParseClaimTopicRemoved is a log parse operation binding the contract event 0x0b1381093c776453c1bbe54fd68be1b235c65db61d099cb50d194b2991e0eec5.
//
// Solidity: event ClaimTopicRemoved(uint256 indexed claimTopic)
func (_ClaimTopicsRegistry *ClaimTopicsRegistryFilterer) ParseClaimTopicRemoved(log types.Log) (*ClaimTopicsRegistryClaimTopicRemoved, error) {
	event := new(ClaimTopicsRegistryClaimTopicRemoved)
	if err := _ClaimTopicsRegistry.contract.UnpackLog(event, "ClaimTopicRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClaimTopicsRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ClaimTopicsRegistry contract.
type ClaimTopicsRegistryInitializedIterator struct {
	Event *ClaimTopicsRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *ClaimTopicsRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClaimTopicsRegistryInitialized)
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
		it.Event = new(ClaimTopicsRegistryInitialized)
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
func (it *ClaimTopicsRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClaimTopicsRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClaimTopicsRegistryInitialized represents a Initialized event raised by the ClaimTopicsRegistry contract.
type ClaimTopicsRegistryInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ClaimTopicsRegistry *ClaimTopicsRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*ClaimTopicsRegistryInitializedIterator, error) {

	logs, sub, err := _ClaimTopicsRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ClaimTopicsRegistryInitializedIterator{contract: _ClaimTopicsRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ClaimTopicsRegistry *ClaimTopicsRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ClaimTopicsRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _ClaimTopicsRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClaimTopicsRegistryInitialized)
				if err := _ClaimTopicsRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ClaimTopicsRegistry *ClaimTopicsRegistryFilterer) ParseInitialized(log types.Log) (*ClaimTopicsRegistryInitialized, error) {
	event := new(ClaimTopicsRegistryInitialized)
	if err := _ClaimTopicsRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClaimTopicsRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ClaimTopicsRegistry contract.
type ClaimTopicsRegistryOwnershipTransferredIterator struct {
	Event *ClaimTopicsRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ClaimTopicsRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClaimTopicsRegistryOwnershipTransferred)
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
		it.Event = new(ClaimTopicsRegistryOwnershipTransferred)
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
func (it *ClaimTopicsRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClaimTopicsRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClaimTopicsRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the ClaimTopicsRegistry contract.
type ClaimTopicsRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ClaimTopicsRegistry *ClaimTopicsRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ClaimTopicsRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ClaimTopicsRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ClaimTopicsRegistryOwnershipTransferredIterator{contract: _ClaimTopicsRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ClaimTopicsRegistry *ClaimTopicsRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ClaimTopicsRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ClaimTopicsRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClaimTopicsRegistryOwnershipTransferred)
				if err := _ClaimTopicsRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ClaimTopicsRegistry *ClaimTopicsRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*ClaimTopicsRegistryOwnershipTransferred, error) {
	event := new(ClaimTopicsRegistryOwnershipTransferred)
	if err := _ClaimTopicsRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
