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

// ModularComplianceProxyMetaData contains all meta data concerning the ModularComplianceProxy contract.
var ModularComplianceProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementationAuthority\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_implementationAuthority\",\"type\":\"address\"}],\"name\":\"ImplementationAuthoritySet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"getImplementationAuthority\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newImplementationAuthority\",\"type\":\"address\"}],\"name\":\"setImplementationAuthority\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161093438038061093483398101604081905261002f9161022a565b6001600160a01b03811661008a5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064015b60405180910390fd5b6100a08160008051602061091483398151915255565b6040516001600160a01b038216907f3b1074392ed8e8424715d0dda2197eede67080b377fc8370e26f3e882207f6b890600090a260006100ec6000805160206109148339815191525490565b6001600160a01b03166361f898256040518163ffffffff1660e01b8152600401602060405180830381865afa158015610129573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061014d919061022a565b60408051600481526024810182526020810180516001600160e01b03166370e39c9560e11b17905290519192506000916001600160a01b038416916101919161025a565b600060405180830381855af49150503d80600081146101cc576040519150601f19603f3d011682016040523d82523d6000602084013e6101d1565b606091505b50509050806102225760405162461bcd60e51b815260206004820152601660248201527f496e697469616c697a6174696f6e206661696c65642e000000000000000000006044820152606401610081565b505050610289565b60006020828403121561023c57600080fd5b81516001600160a01b038116811461025357600080fd5b9392505050565b6000825160005b8181101561027b5760208186018101518583015201610261565b506000920191825250919050565b61067c806102986000396000f3fe6080604052600436106100295760003560e01c80632d5f1187146100e157806392dd9d651461012c575b60006100537f821f3e4d3d679f19eacc940c87acf846ea6eae24a63058ea750304437a62aafc5490565b6001600160a01b03166361f898256040518163ffffffff1660e01b8152600401602060405180830381865afa158015610090573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100b49190610605565b90503660008037600080366000846127105a03f43d806000803e8180156100da57816000f35b816000fd5b005b3480156100ed57600080fd5b507f821f3e4d3d679f19eacc940c87acf846ea6eae24a63058ea750304437a62aafc546040516001600160a01b03909116815260200160405180910390f35b34801561013857600080fd5b506100df610147366004610629565b7f821f3e4d3d679f19eacc940c87acf846ea6eae24a63058ea750304437a62aafc546001600160a01b0316336001600160a01b0316146101f45760405162461bcd60e51b815260206004820152602d60248201527f6f6e6c792063757272656e7420696d706c656d656e746174696f6e417574686f60448201527f726974792063616e2063616c6c0000000000000000000000000000000000000060648201526084015b60405180910390fd5b6001600160a01b03811661024a5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064016101eb565b60006001600160a01b0316816001600160a01b031663709bc7f36040518163ffffffff1660e01b8152600401602060405180830381865afa158015610293573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102b79190610605565b6001600160a01b031614158015610342575060006001600160a01b0316816001600160a01b0316636ff6e83f6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610312573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103369190610605565b6001600160a01b031614155b80156103c2575060006001600160a01b0316816001600160a01b0316631ee9ce8b6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610392573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103b69190610605565b6001600160a01b031614155b8015610442575060006001600160a01b0316816001600160a01b0316639e3e7bb96040518163ffffffff1660e01b8152600401602060405180830381865afa158015610412573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104369190610605565b6001600160a01b031614155b80156104c2575060006001600160a01b0316816001600160a01b03166361f898256040518163ffffffff1660e01b8152600401602060405180830381865afa158015610492573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104b69190610605565b6001600160a01b031614155b8015610542575060006001600160a01b0316816001600160a01b031663fedcc0526040518163ffffffff1660e01b8152600401602060405180830381865afa158015610512573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105369190610605565b6001600160a01b031614155b61058e5760405162461bcd60e51b815260206004820181905260248201527f696e76616c696420496d706c656d656e746174696f6e20417574686f7269747960448201526064016101eb565b6105b6817f821f3e4d3d679f19eacc940c87acf846ea6eae24a63058ea750304437a62aafc55565b6040516001600160a01b038216907f3b1074392ed8e8424715d0dda2197eede67080b377fc8370e26f3e882207f6b890600090a250565b6001600160a01b038116811461060257600080fd5b50565b60006020828403121561061757600080fd5b8151610622816105ed565b9392505050565b60006020828403121561063b57600080fd5b8135610622816105ed56fea2646970667358221220e0a5dcae35d416dd7692c26cec96a7b4e384e3a05b95c86501e5a6315d45274564736f6c63430008110033821f3e4d3d679f19eacc940c87acf846ea6eae24a63058ea750304437a62aafc",
}

// ModularComplianceProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use ModularComplianceProxyMetaData.ABI instead.
var ModularComplianceProxyABI = ModularComplianceProxyMetaData.ABI

// ModularComplianceProxyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ModularComplianceProxyMetaData.Bin instead.
var ModularComplianceProxyBin = ModularComplianceProxyMetaData.Bin

// DeployModularComplianceProxy deploys a new Ethereum contract, binding an instance of ModularComplianceProxy to it.
func DeployModularComplianceProxy(auth *bind.TransactOpts, backend bind.ContractBackend, implementationAuthority common.Address) (common.Address, *types.Transaction, *ModularComplianceProxy, error) {
	parsed, err := ModularComplianceProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ModularComplianceProxyBin), backend, implementationAuthority)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ModularComplianceProxy{ModularComplianceProxyCaller: ModularComplianceProxyCaller{contract: contract}, ModularComplianceProxyTransactor: ModularComplianceProxyTransactor{contract: contract}, ModularComplianceProxyFilterer: ModularComplianceProxyFilterer{contract: contract}}, nil
}

// ModularComplianceProxy is an auto generated Go binding around an Ethereum contract.
type ModularComplianceProxy struct {
	ModularComplianceProxyCaller     // Read-only binding to the contract
	ModularComplianceProxyTransactor // Write-only binding to the contract
	ModularComplianceProxyFilterer   // Log filterer for contract events
}

// ModularComplianceProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type ModularComplianceProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ModularComplianceProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ModularComplianceProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ModularComplianceProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ModularComplianceProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ModularComplianceProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ModularComplianceProxySession struct {
	Contract     *ModularComplianceProxy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ModularComplianceProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ModularComplianceProxyCallerSession struct {
	Contract *ModularComplianceProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// ModularComplianceProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ModularComplianceProxyTransactorSession struct {
	Contract     *ModularComplianceProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// ModularComplianceProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type ModularComplianceProxyRaw struct {
	Contract *ModularComplianceProxy // Generic contract binding to access the raw methods on
}

// ModularComplianceProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ModularComplianceProxyCallerRaw struct {
	Contract *ModularComplianceProxyCaller // Generic read-only contract binding to access the raw methods on
}

// ModularComplianceProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ModularComplianceProxyTransactorRaw struct {
	Contract *ModularComplianceProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewModularComplianceProxy creates a new instance of ModularComplianceProxy, bound to a specific deployed contract.
func NewModularComplianceProxy(address common.Address, backend bind.ContractBackend) (*ModularComplianceProxy, error) {
	contract, err := bindModularComplianceProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ModularComplianceProxy{ModularComplianceProxyCaller: ModularComplianceProxyCaller{contract: contract}, ModularComplianceProxyTransactor: ModularComplianceProxyTransactor{contract: contract}, ModularComplianceProxyFilterer: ModularComplianceProxyFilterer{contract: contract}}, nil
}

// NewModularComplianceProxyCaller creates a new read-only instance of ModularComplianceProxy, bound to a specific deployed contract.
func NewModularComplianceProxyCaller(address common.Address, caller bind.ContractCaller) (*ModularComplianceProxyCaller, error) {
	contract, err := bindModularComplianceProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ModularComplianceProxyCaller{contract: contract}, nil
}

// NewModularComplianceProxyTransactor creates a new write-only instance of ModularComplianceProxy, bound to a specific deployed contract.
func NewModularComplianceProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*ModularComplianceProxyTransactor, error) {
	contract, err := bindModularComplianceProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ModularComplianceProxyTransactor{contract: contract}, nil
}

// NewModularComplianceProxyFilterer creates a new log filterer instance of ModularComplianceProxy, bound to a specific deployed contract.
func NewModularComplianceProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*ModularComplianceProxyFilterer, error) {
	contract, err := bindModularComplianceProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ModularComplianceProxyFilterer{contract: contract}, nil
}

// bindModularComplianceProxy binds a generic wrapper to an already deployed contract.
func bindModularComplianceProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ModularComplianceProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ModularComplianceProxy *ModularComplianceProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ModularComplianceProxy.Contract.ModularComplianceProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ModularComplianceProxy *ModularComplianceProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ModularComplianceProxy.Contract.ModularComplianceProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ModularComplianceProxy *ModularComplianceProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ModularComplianceProxy.Contract.ModularComplianceProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ModularComplianceProxy *ModularComplianceProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ModularComplianceProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ModularComplianceProxy *ModularComplianceProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ModularComplianceProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ModularComplianceProxy *ModularComplianceProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ModularComplianceProxy.Contract.contract.Transact(opts, method, params...)
}

// GetImplementationAuthority is a free data retrieval call binding the contract method 0x2d5f1187.
//
// Solidity: function getImplementationAuthority() view returns(address)
func (_ModularComplianceProxy *ModularComplianceProxyCaller) GetImplementationAuthority(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ModularComplianceProxy.contract.Call(opts, &out, "getImplementationAuthority")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetImplementationAuthority is a free data retrieval call binding the contract method 0x2d5f1187.
//
// Solidity: function getImplementationAuthority() view returns(address)
func (_ModularComplianceProxy *ModularComplianceProxySession) GetImplementationAuthority() (common.Address, error) {
	return _ModularComplianceProxy.Contract.GetImplementationAuthority(&_ModularComplianceProxy.CallOpts)
}

// GetImplementationAuthority is a free data retrieval call binding the contract method 0x2d5f1187.
//
// Solidity: function getImplementationAuthority() view returns(address)
func (_ModularComplianceProxy *ModularComplianceProxyCallerSession) GetImplementationAuthority() (common.Address, error) {
	return _ModularComplianceProxy.Contract.GetImplementationAuthority(&_ModularComplianceProxy.CallOpts)
}

// SetImplementationAuthority is a paid mutator transaction binding the contract method 0x92dd9d65.
//
// Solidity: function setImplementationAuthority(address _newImplementationAuthority) returns()
func (_ModularComplianceProxy *ModularComplianceProxyTransactor) SetImplementationAuthority(opts *bind.TransactOpts, _newImplementationAuthority common.Address) (*types.Transaction, error) {
	return _ModularComplianceProxy.contract.Transact(opts, "setImplementationAuthority", _newImplementationAuthority)
}

// SetImplementationAuthority is a paid mutator transaction binding the contract method 0x92dd9d65.
//
// Solidity: function setImplementationAuthority(address _newImplementationAuthority) returns()
func (_ModularComplianceProxy *ModularComplianceProxySession) SetImplementationAuthority(_newImplementationAuthority common.Address) (*types.Transaction, error) {
	return _ModularComplianceProxy.Contract.SetImplementationAuthority(&_ModularComplianceProxy.TransactOpts, _newImplementationAuthority)
}

// SetImplementationAuthority is a paid mutator transaction binding the contract method 0x92dd9d65.
//
// Solidity: function setImplementationAuthority(address _newImplementationAuthority) returns()
func (_ModularComplianceProxy *ModularComplianceProxyTransactorSession) SetImplementationAuthority(_newImplementationAuthority common.Address) (*types.Transaction, error) {
	return _ModularComplianceProxy.Contract.SetImplementationAuthority(&_ModularComplianceProxy.TransactOpts, _newImplementationAuthority)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ModularComplianceProxy *ModularComplianceProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _ModularComplianceProxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ModularComplianceProxy *ModularComplianceProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ModularComplianceProxy.Contract.Fallback(&_ModularComplianceProxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ModularComplianceProxy *ModularComplianceProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ModularComplianceProxy.Contract.Fallback(&_ModularComplianceProxy.TransactOpts, calldata)
}

// ModularComplianceProxyImplementationAuthoritySetIterator is returned from FilterImplementationAuthoritySet and is used to iterate over the raw logs and unpacked data for ImplementationAuthoritySet events raised by the ModularComplianceProxy contract.
type ModularComplianceProxyImplementationAuthoritySetIterator struct {
	Event *ModularComplianceProxyImplementationAuthoritySet // Event containing the contract specifics and raw log

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
func (it *ModularComplianceProxyImplementationAuthoritySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ModularComplianceProxyImplementationAuthoritySet)
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
		it.Event = new(ModularComplianceProxyImplementationAuthoritySet)
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
func (it *ModularComplianceProxyImplementationAuthoritySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ModularComplianceProxyImplementationAuthoritySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ModularComplianceProxyImplementationAuthoritySet represents a ImplementationAuthoritySet event raised by the ModularComplianceProxy contract.
type ModularComplianceProxyImplementationAuthoritySet struct {
	ImplementationAuthority common.Address
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterImplementationAuthoritySet is a free log retrieval operation binding the contract event 0x3b1074392ed8e8424715d0dda2197eede67080b377fc8370e26f3e882207f6b8.
//
// Solidity: event ImplementationAuthoritySet(address indexed _implementationAuthority)
func (_ModularComplianceProxy *ModularComplianceProxyFilterer) FilterImplementationAuthoritySet(opts *bind.FilterOpts, _implementationAuthority []common.Address) (*ModularComplianceProxyImplementationAuthoritySetIterator, error) {

	var _implementationAuthorityRule []interface{}
	for _, _implementationAuthorityItem := range _implementationAuthority {
		_implementationAuthorityRule = append(_implementationAuthorityRule, _implementationAuthorityItem)
	}

	logs, sub, err := _ModularComplianceProxy.contract.FilterLogs(opts, "ImplementationAuthoritySet", _implementationAuthorityRule)
	if err != nil {
		return nil, err
	}
	return &ModularComplianceProxyImplementationAuthoritySetIterator{contract: _ModularComplianceProxy.contract, event: "ImplementationAuthoritySet", logs: logs, sub: sub}, nil
}

// WatchImplementationAuthoritySet is a free log subscription operation binding the contract event 0x3b1074392ed8e8424715d0dda2197eede67080b377fc8370e26f3e882207f6b8.
//
// Solidity: event ImplementationAuthoritySet(address indexed _implementationAuthority)
func (_ModularComplianceProxy *ModularComplianceProxyFilterer) WatchImplementationAuthoritySet(opts *bind.WatchOpts, sink chan<- *ModularComplianceProxyImplementationAuthoritySet, _implementationAuthority []common.Address) (event.Subscription, error) {

	var _implementationAuthorityRule []interface{}
	for _, _implementationAuthorityItem := range _implementationAuthority {
		_implementationAuthorityRule = append(_implementationAuthorityRule, _implementationAuthorityItem)
	}

	logs, sub, err := _ModularComplianceProxy.contract.WatchLogs(opts, "ImplementationAuthoritySet", _implementationAuthorityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ModularComplianceProxyImplementationAuthoritySet)
				if err := _ModularComplianceProxy.contract.UnpackLog(event, "ImplementationAuthoritySet", log); err != nil {
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

// ParseImplementationAuthoritySet is a log parse operation binding the contract event 0x3b1074392ed8e8424715d0dda2197eede67080b377fc8370e26f3e882207f6b8.
//
// Solidity: event ImplementationAuthoritySet(address indexed _implementationAuthority)
func (_ModularComplianceProxy *ModularComplianceProxyFilterer) ParseImplementationAuthoritySet(log types.Log) (*ModularComplianceProxyImplementationAuthoritySet, error) {
	event := new(ModularComplianceProxyImplementationAuthoritySet)
	if err := _ModularComplianceProxy.contract.UnpackLog(event, "ImplementationAuthoritySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ModularComplianceProxyInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ModularComplianceProxy contract.
type ModularComplianceProxyInitializedIterator struct {
	Event *ModularComplianceProxyInitialized // Event containing the contract specifics and raw log

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
func (it *ModularComplianceProxyInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ModularComplianceProxyInitialized)
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
		it.Event = new(ModularComplianceProxyInitialized)
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
func (it *ModularComplianceProxyInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ModularComplianceProxyInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ModularComplianceProxyInitialized represents a Initialized event raised by the ModularComplianceProxy contract.
type ModularComplianceProxyInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ModularComplianceProxy *ModularComplianceProxyFilterer) FilterInitialized(opts *bind.FilterOpts) (*ModularComplianceProxyInitializedIterator, error) {

	logs, sub, err := _ModularComplianceProxy.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ModularComplianceProxyInitializedIterator{contract: _ModularComplianceProxy.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ModularComplianceProxy *ModularComplianceProxyFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ModularComplianceProxyInitialized) (event.Subscription, error) {

	logs, sub, err := _ModularComplianceProxy.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ModularComplianceProxyInitialized)
				if err := _ModularComplianceProxy.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ModularComplianceProxy *ModularComplianceProxyFilterer) ParseInitialized(log types.Log) (*ModularComplianceProxyInitialized, error) {
	event := new(ModularComplianceProxyInitialized)
	if err := _ModularComplianceProxy.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
