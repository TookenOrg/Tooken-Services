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

// TokenProxyMetaData contains all meta data concerning the TokenProxy contract.
var TokenProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementationAuthority\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_identityRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"_onchainID\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_implementationAuthority\",\"type\":\"address\"}],\"name\":\"ImplementationAuthoritySet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"getImplementationAuthority\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newImplementationAuthority\",\"type\":\"address\"}],\"name\":\"setImplementationAuthority\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162000d3d38038062000d3d8339810160408190526200003491620004dd565b6001600160a01b038716158015906200005557506001600160a01b03861615155b80156200006a57506001600160a01b03851615155b620000bc5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064015b60405180910390fd5b604051602001620000d890602080825260009082015260400190565b6040516020818303038152906040528051906020012084604051602001620001019190620005da565b60405160208183030381529060405280519060200120141580156200017f57506040516020016200013d90602080825260009082015260400190565b6040516020818303038152906040528051906020012083604051602001620001669190620005da565b6040516020818303038152906040528051906020012014155b620001cd5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d20656d70747920737472696e67006044820152606401620000b3565b60128260ff161115620002235760405162461bcd60e51b815260206004820152601960248201527f646563696d616c73206265747765656e203020616e64203138000000000000006044820152606401620000b3565b6200023b8760008051602062000d1d83398151915255565b6040516001600160a01b038816907f3b1074392ed8e8424715d0dda2197eede67080b377fc8370e26f3e882207f6b890600090a260006200028960008051602062000d1d8339815191525490565b6001600160a01b031663709bc7f36040518163ffffffff1660e01b8152600401602060405180830381865afa158015620002c7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002ed9190620005f6565b90506000816001600160a01b0316888888888888604051602401620003189695949392919062000614565b60408051601f198184030181529181526020820180516001600160e01b0316633e46d86760e21b179052516200034f919062000673565b600060405180830381855af49150503d80600081146200038c576040519150601f19603f3d011682016040523d82523d6000602084013e62000391565b606091505b5050905080620003e45760405162461bcd60e51b815260206004820152601660248201527f496e697469616c697a6174696f6e206661696c65642e000000000000000000006044820152606401620000b3565b50505050505050505062000691565b80516001600160a01b03811681146200040b57600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b60005b838110156200044357818101518382015260200162000429565b50506000910152565b600082601f8301126200045e57600080fd5b81516001600160401b03808211156200047b576200047b62000410565b604051601f8301601f19908116603f01168101908282118183101715620004a657620004a662000410565b81604052838152866020858801011115620004c057600080fd5b620004d384602083016020890162000426565b9695505050505050565b600080600080600080600060e0888a031215620004f957600080fd5b6200050488620003f3565b96506200051460208901620003f3565b95506200052460408901620003f3565b60608901519095506001600160401b03808211156200054257600080fd5b620005508b838c016200044c565b955060808a01519150808211156200056757600080fd5b50620005768a828b016200044c565b93505060a088015160ff811681146200058e57600080fd5b91506200059e60c08901620003f3565b905092959891949750929550565b60008151808452620005c681602086016020860162000426565b601f01601f19169290920160200192915050565b602081526000620005ef6020830184620005ac565b9392505050565b6000602082840312156200060957600080fd5b620005ef82620003f3565b600060018060a01b038089168352808816602084015260c060408401526200064060c0840188620005ac565b8381036060850152620006548188620005ac565b60ff969096166080850152509290921660a09091015250949350505050565b600082516200068781846020870162000426565b9190910192915050565b61067c80620006a16000396000f3fe6080604052600436106100295760003560e01c80632d5f1187146100e157806392dd9d651461012c575b60006100537f821f3e4d3d679f19eacc940c87acf846ea6eae24a63058ea750304437a62aafc5490565b6001600160a01b031663709bc7f36040518163ffffffff1660e01b8152600401602060405180830381865afa158015610090573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100b49190610605565b90503660008037600080366000846127105a03f43d806000803e8180156100da57816000f35b816000fd5b005b3480156100ed57600080fd5b507f821f3e4d3d679f19eacc940c87acf846ea6eae24a63058ea750304437a62aafc546040516001600160a01b03909116815260200160405180910390f35b34801561013857600080fd5b506100df610147366004610629565b7f821f3e4d3d679f19eacc940c87acf846ea6eae24a63058ea750304437a62aafc546001600160a01b0316336001600160a01b0316146101f45760405162461bcd60e51b815260206004820152602d60248201527f6f6e6c792063757272656e7420696d706c656d656e746174696f6e417574686f60448201527f726974792063616e2063616c6c0000000000000000000000000000000000000060648201526084015b60405180910390fd5b6001600160a01b03811661024a5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064016101eb565b60006001600160a01b0316816001600160a01b031663709bc7f36040518163ffffffff1660e01b8152600401602060405180830381865afa158015610293573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102b79190610605565b6001600160a01b031614158015610342575060006001600160a01b0316816001600160a01b0316636ff6e83f6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610312573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103369190610605565b6001600160a01b031614155b80156103c2575060006001600160a01b0316816001600160a01b0316631ee9ce8b6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610392573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103b69190610605565b6001600160a01b031614155b8015610442575060006001600160a01b0316816001600160a01b0316639e3e7bb96040518163ffffffff1660e01b8152600401602060405180830381865afa158015610412573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104369190610605565b6001600160a01b031614155b80156104c2575060006001600160a01b0316816001600160a01b03166361f898256040518163ffffffff1660e01b8152600401602060405180830381865afa158015610492573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104b69190610605565b6001600160a01b031614155b8015610542575060006001600160a01b0316816001600160a01b031663fedcc0526040518163ffffffff1660e01b8152600401602060405180830381865afa158015610512573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105369190610605565b6001600160a01b031614155b61058e5760405162461bcd60e51b815260206004820181905260248201527f696e76616c696420496d706c656d656e746174696f6e20417574686f7269747960448201526064016101eb565b6105b6817f821f3e4d3d679f19eacc940c87acf846ea6eae24a63058ea750304437a62aafc55565b6040516001600160a01b038216907f3b1074392ed8e8424715d0dda2197eede67080b377fc8370e26f3e882207f6b890600090a250565b6001600160a01b038116811461060257600080fd5b50565b60006020828403121561061757600080fd5b8151610622816105ed565b9392505050565b60006020828403121561063b57600080fd5b8135610622816105ed56fea26469706673582212202bd77fed6b471b5b325753b3412937bfc7470420aefad62d3ed0305ea220e80f64736f6c63430008110033821f3e4d3d679f19eacc940c87acf846ea6eae24a63058ea750304437a62aafc",
}

// TokenProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use TokenProxyMetaData.ABI instead.
var TokenProxyABI = TokenProxyMetaData.ABI

// TokenProxyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TokenProxyMetaData.Bin instead.
var TokenProxyBin = TokenProxyMetaData.Bin

// DeployTokenProxy deploys a new Ethereum contract, binding an instance of TokenProxy to it.
func DeployTokenProxy(auth *bind.TransactOpts, backend bind.ContractBackend, implementationAuthority common.Address, _identityRegistry common.Address, _compliance common.Address, _name string, _symbol string, _decimals uint8, _onchainID common.Address) (common.Address, *types.Transaction, *TokenProxy, error) {
	parsed, err := TokenProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TokenProxyBin), backend, implementationAuthority, _identityRegistry, _compliance, _name, _symbol, _decimals, _onchainID)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TokenProxy{TokenProxyCaller: TokenProxyCaller{contract: contract}, TokenProxyTransactor: TokenProxyTransactor{contract: contract}, TokenProxyFilterer: TokenProxyFilterer{contract: contract}}, nil
}

// TokenProxy is an auto generated Go binding around an Ethereum contract.
type TokenProxy struct {
	TokenProxyCaller     // Read-only binding to the contract
	TokenProxyTransactor // Write-only binding to the contract
	TokenProxyFilterer   // Log filterer for contract events
}

// TokenProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenProxySession struct {
	Contract     *TokenProxy       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenProxyCallerSession struct {
	Contract *TokenProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// TokenProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenProxyTransactorSession struct {
	Contract     *TokenProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// TokenProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenProxyRaw struct {
	Contract *TokenProxy // Generic contract binding to access the raw methods on
}

// TokenProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenProxyCallerRaw struct {
	Contract *TokenProxyCaller // Generic read-only contract binding to access the raw methods on
}

// TokenProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenProxyTransactorRaw struct {
	Contract *TokenProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenProxy creates a new instance of TokenProxy, bound to a specific deployed contract.
func NewTokenProxy(address common.Address, backend bind.ContractBackend) (*TokenProxy, error) {
	contract, err := bindTokenProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenProxy{TokenProxyCaller: TokenProxyCaller{contract: contract}, TokenProxyTransactor: TokenProxyTransactor{contract: contract}, TokenProxyFilterer: TokenProxyFilterer{contract: contract}}, nil
}

// NewTokenProxyCaller creates a new read-only instance of TokenProxy, bound to a specific deployed contract.
func NewTokenProxyCaller(address common.Address, caller bind.ContractCaller) (*TokenProxyCaller, error) {
	contract, err := bindTokenProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenProxyCaller{contract: contract}, nil
}

// NewTokenProxyTransactor creates a new write-only instance of TokenProxy, bound to a specific deployed contract.
func NewTokenProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenProxyTransactor, error) {
	contract, err := bindTokenProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenProxyTransactor{contract: contract}, nil
}

// NewTokenProxyFilterer creates a new log filterer instance of TokenProxy, bound to a specific deployed contract.
func NewTokenProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenProxyFilterer, error) {
	contract, err := bindTokenProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenProxyFilterer{contract: contract}, nil
}

// bindTokenProxy binds a generic wrapper to an already deployed contract.
func bindTokenProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TokenProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenProxy *TokenProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenProxy.Contract.TokenProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenProxy *TokenProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenProxy.Contract.TokenProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenProxy *TokenProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenProxy.Contract.TokenProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenProxy *TokenProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenProxy *TokenProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenProxy *TokenProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenProxy.Contract.contract.Transact(opts, method, params...)
}

// GetImplementationAuthority is a free data retrieval call binding the contract method 0x2d5f1187.
//
// Solidity: function getImplementationAuthority() view returns(address)
func (_TokenProxy *TokenProxyCaller) GetImplementationAuthority(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenProxy.contract.Call(opts, &out, "getImplementationAuthority")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetImplementationAuthority is a free data retrieval call binding the contract method 0x2d5f1187.
//
// Solidity: function getImplementationAuthority() view returns(address)
func (_TokenProxy *TokenProxySession) GetImplementationAuthority() (common.Address, error) {
	return _TokenProxy.Contract.GetImplementationAuthority(&_TokenProxy.CallOpts)
}

// GetImplementationAuthority is a free data retrieval call binding the contract method 0x2d5f1187.
//
// Solidity: function getImplementationAuthority() view returns(address)
func (_TokenProxy *TokenProxyCallerSession) GetImplementationAuthority() (common.Address, error) {
	return _TokenProxy.Contract.GetImplementationAuthority(&_TokenProxy.CallOpts)
}

// SetImplementationAuthority is a paid mutator transaction binding the contract method 0x92dd9d65.
//
// Solidity: function setImplementationAuthority(address _newImplementationAuthority) returns()
func (_TokenProxy *TokenProxyTransactor) SetImplementationAuthority(opts *bind.TransactOpts, _newImplementationAuthority common.Address) (*types.Transaction, error) {
	return _TokenProxy.contract.Transact(opts, "setImplementationAuthority", _newImplementationAuthority)
}

// SetImplementationAuthority is a paid mutator transaction binding the contract method 0x92dd9d65.
//
// Solidity: function setImplementationAuthority(address _newImplementationAuthority) returns()
func (_TokenProxy *TokenProxySession) SetImplementationAuthority(_newImplementationAuthority common.Address) (*types.Transaction, error) {
	return _TokenProxy.Contract.SetImplementationAuthority(&_TokenProxy.TransactOpts, _newImplementationAuthority)
}

// SetImplementationAuthority is a paid mutator transaction binding the contract method 0x92dd9d65.
//
// Solidity: function setImplementationAuthority(address _newImplementationAuthority) returns()
func (_TokenProxy *TokenProxyTransactorSession) SetImplementationAuthority(_newImplementationAuthority common.Address) (*types.Transaction, error) {
	return _TokenProxy.Contract.SetImplementationAuthority(&_TokenProxy.TransactOpts, _newImplementationAuthority)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TokenProxy *TokenProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _TokenProxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TokenProxy *TokenProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TokenProxy.Contract.Fallback(&_TokenProxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TokenProxy *TokenProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TokenProxy.Contract.Fallback(&_TokenProxy.TransactOpts, calldata)
}

// TokenProxyImplementationAuthoritySetIterator is returned from FilterImplementationAuthoritySet and is used to iterate over the raw logs and unpacked data for ImplementationAuthoritySet events raised by the TokenProxy contract.
type TokenProxyImplementationAuthoritySetIterator struct {
	Event *TokenProxyImplementationAuthoritySet // Event containing the contract specifics and raw log

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
func (it *TokenProxyImplementationAuthoritySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenProxyImplementationAuthoritySet)
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
		it.Event = new(TokenProxyImplementationAuthoritySet)
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
func (it *TokenProxyImplementationAuthoritySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenProxyImplementationAuthoritySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenProxyImplementationAuthoritySet represents a ImplementationAuthoritySet event raised by the TokenProxy contract.
type TokenProxyImplementationAuthoritySet struct {
	ImplementationAuthority common.Address
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterImplementationAuthoritySet is a free log retrieval operation binding the contract event 0x3b1074392ed8e8424715d0dda2197eede67080b377fc8370e26f3e882207f6b8.
//
// Solidity: event ImplementationAuthoritySet(address indexed _implementationAuthority)
func (_TokenProxy *TokenProxyFilterer) FilterImplementationAuthoritySet(opts *bind.FilterOpts, _implementationAuthority []common.Address) (*TokenProxyImplementationAuthoritySetIterator, error) {

	var _implementationAuthorityRule []interface{}
	for _, _implementationAuthorityItem := range _implementationAuthority {
		_implementationAuthorityRule = append(_implementationAuthorityRule, _implementationAuthorityItem)
	}

	logs, sub, err := _TokenProxy.contract.FilterLogs(opts, "ImplementationAuthoritySet", _implementationAuthorityRule)
	if err != nil {
		return nil, err
	}
	return &TokenProxyImplementationAuthoritySetIterator{contract: _TokenProxy.contract, event: "ImplementationAuthoritySet", logs: logs, sub: sub}, nil
}

// WatchImplementationAuthoritySet is a free log subscription operation binding the contract event 0x3b1074392ed8e8424715d0dda2197eede67080b377fc8370e26f3e882207f6b8.
//
// Solidity: event ImplementationAuthoritySet(address indexed _implementationAuthority)
func (_TokenProxy *TokenProxyFilterer) WatchImplementationAuthoritySet(opts *bind.WatchOpts, sink chan<- *TokenProxyImplementationAuthoritySet, _implementationAuthority []common.Address) (event.Subscription, error) {

	var _implementationAuthorityRule []interface{}
	for _, _implementationAuthorityItem := range _implementationAuthority {
		_implementationAuthorityRule = append(_implementationAuthorityRule, _implementationAuthorityItem)
	}

	logs, sub, err := _TokenProxy.contract.WatchLogs(opts, "ImplementationAuthoritySet", _implementationAuthorityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenProxyImplementationAuthoritySet)
				if err := _TokenProxy.contract.UnpackLog(event, "ImplementationAuthoritySet", log); err != nil {
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
func (_TokenProxy *TokenProxyFilterer) ParseImplementationAuthoritySet(log types.Log) (*TokenProxyImplementationAuthoritySet, error) {
	event := new(TokenProxyImplementationAuthoritySet)
	if err := _TokenProxy.contract.UnpackLog(event, "ImplementationAuthoritySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenProxyInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the TokenProxy contract.
type TokenProxyInitializedIterator struct {
	Event *TokenProxyInitialized // Event containing the contract specifics and raw log

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
func (it *TokenProxyInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenProxyInitialized)
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
		it.Event = new(TokenProxyInitialized)
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
func (it *TokenProxyInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenProxyInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenProxyInitialized represents a Initialized event raised by the TokenProxy contract.
type TokenProxyInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_TokenProxy *TokenProxyFilterer) FilterInitialized(opts *bind.FilterOpts) (*TokenProxyInitializedIterator, error) {

	logs, sub, err := _TokenProxy.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &TokenProxyInitializedIterator{contract: _TokenProxy.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_TokenProxy *TokenProxyFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *TokenProxyInitialized) (event.Subscription, error) {

	logs, sub, err := _TokenProxy.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenProxyInitialized)
				if err := _TokenProxy.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_TokenProxy *TokenProxyFilterer) ParseInitialized(log types.Log) (*TokenProxyInitialized, error) {
	event := new(TokenProxyInitialized)
	if err := _TokenProxy.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
