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

// ImplementationAuthorityMetaData contains all meta data concerning the ImplementationAuthority contract.
var ImplementationAuthorityMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"UpdatedImplementation\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getImplementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newImplementation\",\"type\":\"address\"}],\"name\":\"updateImplementation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516104d33803806104d383398101604081905261002f91610139565b610038336100e9565b6001600160a01b0381166100925760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f206164647265737300604482015260640160405180910390fd5b600180546001600160a01b0319166001600160a01b0383169081179091556040519081527f87c4e67a766ffddda27f441d63853a36ae64fbb07775a7c59d395e064b204eeb9060200160405180910390a150610169565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b60006020828403121561014b57600080fd5b81516001600160a01b038116811461016257600080fd5b9392505050565b61035b806101786000396000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c80638da5cb5b116100505780638da5cb5b14610089578063aaf10f42146100b2578063f2fde38b146100c357600080fd5b8063025b22bc1461006c578063715018a614610081575b600080fd5b61007f61007a3660046102f5565b6100d6565b005b61007f61019a565b6000546001600160a01b03165b6040516001600160a01b03909116815260200160405180910390f35b6001546001600160a01b0316610096565b61007f6100d13660046102f5565b6101ae565b6100de61023e565b6001600160a01b0381166101395760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064015b60405180910390fd5b6001805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0383169081179091556040519081527f87c4e67a766ffddda27f441d63853a36ae64fbb07775a7c59d395e064b204eeb9060200160405180910390a150565b6101a261023e565b6101ac6000610298565b565b6101b661023e565b6001600160a01b0381166102325760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610130565b61023b81610298565b50565b6000546001600160a01b031633146101ac5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610130565b600080546001600160a01b0383811673ffffffffffffffffffffffffffffffffffffffff19831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b60006020828403121561030757600080fd5b81356001600160a01b038116811461031e57600080fd5b939250505056fea26469706673582212200a1dba1eb343b4059d1a95d018f9b13319b86a4091805729f70018dd46a40dc664736f6c63430008110033",
}

// ImplementationAuthorityABI is the input ABI used to generate the binding from.
// Deprecated: Use ImplementationAuthorityMetaData.ABI instead.
var ImplementationAuthorityABI = ImplementationAuthorityMetaData.ABI

// ImplementationAuthorityBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ImplementationAuthorityMetaData.Bin instead.
var ImplementationAuthorityBin = ImplementationAuthorityMetaData.Bin

// DeployImplementationAuthority deploys a new Ethereum contract, binding an instance of ImplementationAuthority to it.
func DeployImplementationAuthority(auth *bind.TransactOpts, backend bind.ContractBackend, implementation common.Address) (common.Address, *types.Transaction, *ImplementationAuthority, error) {
	parsed, err := ImplementationAuthorityMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ImplementationAuthorityBin), backend, implementation)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ImplementationAuthority{ImplementationAuthorityCaller: ImplementationAuthorityCaller{contract: contract}, ImplementationAuthorityTransactor: ImplementationAuthorityTransactor{contract: contract}, ImplementationAuthorityFilterer: ImplementationAuthorityFilterer{contract: contract}}, nil
}

// ImplementationAuthority is an auto generated Go binding around an Ethereum contract.
type ImplementationAuthority struct {
	ImplementationAuthorityCaller     // Read-only binding to the contract
	ImplementationAuthorityTransactor // Write-only binding to the contract
	ImplementationAuthorityFilterer   // Log filterer for contract events
}

// ImplementationAuthorityCaller is an auto generated read-only Go binding around an Ethereum contract.
type ImplementationAuthorityCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ImplementationAuthorityTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ImplementationAuthorityTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ImplementationAuthorityFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ImplementationAuthorityFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ImplementationAuthoritySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ImplementationAuthoritySession struct {
	Contract     *ImplementationAuthority // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// ImplementationAuthorityCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ImplementationAuthorityCallerSession struct {
	Contract *ImplementationAuthorityCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// ImplementationAuthorityTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ImplementationAuthorityTransactorSession struct {
	Contract     *ImplementationAuthorityTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// ImplementationAuthorityRaw is an auto generated low-level Go binding around an Ethereum contract.
type ImplementationAuthorityRaw struct {
	Contract *ImplementationAuthority // Generic contract binding to access the raw methods on
}

// ImplementationAuthorityCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ImplementationAuthorityCallerRaw struct {
	Contract *ImplementationAuthorityCaller // Generic read-only contract binding to access the raw methods on
}

// ImplementationAuthorityTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ImplementationAuthorityTransactorRaw struct {
	Contract *ImplementationAuthorityTransactor // Generic write-only contract binding to access the raw methods on
}

// NewImplementationAuthority creates a new instance of ImplementationAuthority, bound to a specific deployed contract.
func NewImplementationAuthority(address common.Address, backend bind.ContractBackend) (*ImplementationAuthority, error) {
	contract, err := bindImplementationAuthority(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ImplementationAuthority{ImplementationAuthorityCaller: ImplementationAuthorityCaller{contract: contract}, ImplementationAuthorityTransactor: ImplementationAuthorityTransactor{contract: contract}, ImplementationAuthorityFilterer: ImplementationAuthorityFilterer{contract: contract}}, nil
}

// NewImplementationAuthorityCaller creates a new read-only instance of ImplementationAuthority, bound to a specific deployed contract.
func NewImplementationAuthorityCaller(address common.Address, caller bind.ContractCaller) (*ImplementationAuthorityCaller, error) {
	contract, err := bindImplementationAuthority(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ImplementationAuthorityCaller{contract: contract}, nil
}

// NewImplementationAuthorityTransactor creates a new write-only instance of ImplementationAuthority, bound to a specific deployed contract.
func NewImplementationAuthorityTransactor(address common.Address, transactor bind.ContractTransactor) (*ImplementationAuthorityTransactor, error) {
	contract, err := bindImplementationAuthority(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ImplementationAuthorityTransactor{contract: contract}, nil
}

// NewImplementationAuthorityFilterer creates a new log filterer instance of ImplementationAuthority, bound to a specific deployed contract.
func NewImplementationAuthorityFilterer(address common.Address, filterer bind.ContractFilterer) (*ImplementationAuthorityFilterer, error) {
	contract, err := bindImplementationAuthority(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ImplementationAuthorityFilterer{contract: contract}, nil
}

// bindImplementationAuthority binds a generic wrapper to an already deployed contract.
func bindImplementationAuthority(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ImplementationAuthorityMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ImplementationAuthority *ImplementationAuthorityRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ImplementationAuthority.Contract.ImplementationAuthorityCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ImplementationAuthority *ImplementationAuthorityRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ImplementationAuthority.Contract.ImplementationAuthorityTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ImplementationAuthority *ImplementationAuthorityRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ImplementationAuthority.Contract.ImplementationAuthorityTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ImplementationAuthority *ImplementationAuthorityCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ImplementationAuthority.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ImplementationAuthority *ImplementationAuthorityTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ImplementationAuthority.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ImplementationAuthority *ImplementationAuthorityTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ImplementationAuthority.Contract.contract.Transact(opts, method, params...)
}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_ImplementationAuthority *ImplementationAuthorityCaller) GetImplementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ImplementationAuthority.contract.Call(opts, &out, "getImplementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_ImplementationAuthority *ImplementationAuthoritySession) GetImplementation() (common.Address, error) {
	return _ImplementationAuthority.Contract.GetImplementation(&_ImplementationAuthority.CallOpts)
}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_ImplementationAuthority *ImplementationAuthorityCallerSession) GetImplementation() (common.Address, error) {
	return _ImplementationAuthority.Contract.GetImplementation(&_ImplementationAuthority.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ImplementationAuthority *ImplementationAuthorityCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ImplementationAuthority.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ImplementationAuthority *ImplementationAuthoritySession) Owner() (common.Address, error) {
	return _ImplementationAuthority.Contract.Owner(&_ImplementationAuthority.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ImplementationAuthority *ImplementationAuthorityCallerSession) Owner() (common.Address, error) {
	return _ImplementationAuthority.Contract.Owner(&_ImplementationAuthority.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ImplementationAuthority *ImplementationAuthorityTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ImplementationAuthority.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ImplementationAuthority *ImplementationAuthoritySession) RenounceOwnership() (*types.Transaction, error) {
	return _ImplementationAuthority.Contract.RenounceOwnership(&_ImplementationAuthority.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ImplementationAuthority *ImplementationAuthorityTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ImplementationAuthority.Contract.RenounceOwnership(&_ImplementationAuthority.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ImplementationAuthority *ImplementationAuthorityTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ImplementationAuthority.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ImplementationAuthority *ImplementationAuthoritySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ImplementationAuthority.Contract.TransferOwnership(&_ImplementationAuthority.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ImplementationAuthority *ImplementationAuthorityTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ImplementationAuthority.Contract.TransferOwnership(&_ImplementationAuthority.TransactOpts, newOwner)
}

// UpdateImplementation is a paid mutator transaction binding the contract method 0x025b22bc.
//
// Solidity: function updateImplementation(address _newImplementation) returns()
func (_ImplementationAuthority *ImplementationAuthorityTransactor) UpdateImplementation(opts *bind.TransactOpts, _newImplementation common.Address) (*types.Transaction, error) {
	return _ImplementationAuthority.contract.Transact(opts, "updateImplementation", _newImplementation)
}

// UpdateImplementation is a paid mutator transaction binding the contract method 0x025b22bc.
//
// Solidity: function updateImplementation(address _newImplementation) returns()
func (_ImplementationAuthority *ImplementationAuthoritySession) UpdateImplementation(_newImplementation common.Address) (*types.Transaction, error) {
	return _ImplementationAuthority.Contract.UpdateImplementation(&_ImplementationAuthority.TransactOpts, _newImplementation)
}

// UpdateImplementation is a paid mutator transaction binding the contract method 0x025b22bc.
//
// Solidity: function updateImplementation(address _newImplementation) returns()
func (_ImplementationAuthority *ImplementationAuthorityTransactorSession) UpdateImplementation(_newImplementation common.Address) (*types.Transaction, error) {
	return _ImplementationAuthority.Contract.UpdateImplementation(&_ImplementationAuthority.TransactOpts, _newImplementation)
}

// ImplementationAuthorityOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ImplementationAuthority contract.
type ImplementationAuthorityOwnershipTransferredIterator struct {
	Event *ImplementationAuthorityOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ImplementationAuthorityOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ImplementationAuthorityOwnershipTransferred)
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
		it.Event = new(ImplementationAuthorityOwnershipTransferred)
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
func (it *ImplementationAuthorityOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ImplementationAuthorityOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ImplementationAuthorityOwnershipTransferred represents a OwnershipTransferred event raised by the ImplementationAuthority contract.
type ImplementationAuthorityOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ImplementationAuthority *ImplementationAuthorityFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ImplementationAuthorityOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ImplementationAuthority.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ImplementationAuthorityOwnershipTransferredIterator{contract: _ImplementationAuthority.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ImplementationAuthority *ImplementationAuthorityFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ImplementationAuthorityOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ImplementationAuthority.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ImplementationAuthorityOwnershipTransferred)
				if err := _ImplementationAuthority.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ImplementationAuthority *ImplementationAuthorityFilterer) ParseOwnershipTransferred(log types.Log) (*ImplementationAuthorityOwnershipTransferred, error) {
	event := new(ImplementationAuthorityOwnershipTransferred)
	if err := _ImplementationAuthority.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ImplementationAuthorityUpdatedImplementationIterator is returned from FilterUpdatedImplementation and is used to iterate over the raw logs and unpacked data for UpdatedImplementation events raised by the ImplementationAuthority contract.
type ImplementationAuthorityUpdatedImplementationIterator struct {
	Event *ImplementationAuthorityUpdatedImplementation // Event containing the contract specifics and raw log

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
func (it *ImplementationAuthorityUpdatedImplementationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ImplementationAuthorityUpdatedImplementation)
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
		it.Event = new(ImplementationAuthorityUpdatedImplementation)
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
func (it *ImplementationAuthorityUpdatedImplementationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ImplementationAuthorityUpdatedImplementationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ImplementationAuthorityUpdatedImplementation represents a UpdatedImplementation event raised by the ImplementationAuthority contract.
type ImplementationAuthorityUpdatedImplementation struct {
	NewAddress common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUpdatedImplementation is a free log retrieval operation binding the contract event 0x87c4e67a766ffddda27f441d63853a36ae64fbb07775a7c59d395e064b204eeb.
//
// Solidity: event UpdatedImplementation(address newAddress)
func (_ImplementationAuthority *ImplementationAuthorityFilterer) FilterUpdatedImplementation(opts *bind.FilterOpts) (*ImplementationAuthorityUpdatedImplementationIterator, error) {

	logs, sub, err := _ImplementationAuthority.contract.FilterLogs(opts, "UpdatedImplementation")
	if err != nil {
		return nil, err
	}
	return &ImplementationAuthorityUpdatedImplementationIterator{contract: _ImplementationAuthority.contract, event: "UpdatedImplementation", logs: logs, sub: sub}, nil
}

// WatchUpdatedImplementation is a free log subscription operation binding the contract event 0x87c4e67a766ffddda27f441d63853a36ae64fbb07775a7c59d395e064b204eeb.
//
// Solidity: event UpdatedImplementation(address newAddress)
func (_ImplementationAuthority *ImplementationAuthorityFilterer) WatchUpdatedImplementation(opts *bind.WatchOpts, sink chan<- *ImplementationAuthorityUpdatedImplementation) (event.Subscription, error) {

	logs, sub, err := _ImplementationAuthority.contract.WatchLogs(opts, "UpdatedImplementation")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ImplementationAuthorityUpdatedImplementation)
				if err := _ImplementationAuthority.contract.UnpackLog(event, "UpdatedImplementation", log); err != nil {
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

// ParseUpdatedImplementation is a log parse operation binding the contract event 0x87c4e67a766ffddda27f441d63853a36ae64fbb07775a7c59d395e064b204eeb.
//
// Solidity: event UpdatedImplementation(address newAddress)
func (_ImplementationAuthority *ImplementationAuthorityFilterer) ParseUpdatedImplementation(log types.Log) (*ImplementationAuthorityUpdatedImplementation, error) {
	event := new(ImplementationAuthorityUpdatedImplementation)
	if err := _ImplementationAuthority.contract.UnpackLog(event, "UpdatedImplementation", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
