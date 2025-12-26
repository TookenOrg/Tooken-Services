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

// IdentityProxyMetaData contains all meta data concerning the IdentityProxy contract.
var IdentityProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_implementationAuthority\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"initialManagementKey\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"implementationAuthority\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516104b03803806104b083398101604081905261002f91610271565b6001600160a01b03821661008a5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064015b60405180910390fd5b6001600160a01b0381166100e05760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f2061646472657373006044820152606401610081565b817f821f3e4d3d679f19eacc940c87acf846ea6eae24a63058ea750304437a62aafc556000826001600160a01b031663aaf10f426040518163ffffffff1660e01b8152600401602060405180830381865afa158015610143573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061016791906102a4565b6040516001600160a01b03848116602483015291925060009183169060440160408051601f198184030181529181526020820180516001600160e01b031663189acdbd60e31b179052516101bb91906102c6565b600060405180830381855af49150503d80600081146101f6576040519150601f19603f3d011682016040523d82523d6000602084013e6101fb565b606091505b505090508061024c5760405162461bcd60e51b815260206004820152601660248201527f496e697469616c697a6174696f6e206661696c65642e000000000000000000006044820152606401610081565b505050506102f5565b80516001600160a01b038116811461026c57600080fd5b919050565b6000806040838503121561028457600080fd5b61028d83610255565b915061029b60208401610255565b90509250929050565b6000602082840312156102b657600080fd5b6102bf82610255565b9392505050565b6000825160005b818110156102e757602081860181015185830152016102cd565b506000920191825250919050565b6101ac806103046000396000f3fe60806040526004361061001e5760003560e01c80632307f882146100e1575b60006100487f821f3e4d3d679f19eacc940c87acf846ea6eae24a63058ea750304437a62aafc5490565b73ffffffffffffffffffffffffffffffffffffffff1663aaf10f426040518163ffffffff1660e01b8152600401602060405180830381865afa158015610092573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100b69190610139565b90503660008037600080366000846127105a03f43d806000803e8180156100dc57816000f35b816000fd5b3480156100ed57600080fd5b507f821f3e4d3d679f19eacc940c87acf846ea6eae24a63058ea750304437a62aafc5460405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b60006020828403121561014b57600080fd5b815173ffffffffffffffffffffffffffffffffffffffff8116811461016f57600080fd5b939250505056fea2646970667358221220dc10dba4dcb99f75cb91819000d74353c98e3dd7471a2af2095fedc6a70516b664736f6c63430008110033",
}

// IdentityProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use IdentityProxyMetaData.ABI instead.
var IdentityProxyABI = IdentityProxyMetaData.ABI

// IdentityProxyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use IdentityProxyMetaData.Bin instead.
var IdentityProxyBin = IdentityProxyMetaData.Bin

// DeployIdentityProxy deploys a new Ethereum contract, binding an instance of IdentityProxy to it.
func DeployIdentityProxy(auth *bind.TransactOpts, backend bind.ContractBackend, _implementationAuthority common.Address, initialManagementKey common.Address) (common.Address, *types.Transaction, *IdentityProxy, error) {
	parsed, err := IdentityProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(IdentityProxyBin), backend, _implementationAuthority, initialManagementKey)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &IdentityProxy{IdentityProxyCaller: IdentityProxyCaller{contract: contract}, IdentityProxyTransactor: IdentityProxyTransactor{contract: contract}, IdentityProxyFilterer: IdentityProxyFilterer{contract: contract}}, nil
}

// IdentityProxy is an auto generated Go binding around an Ethereum contract.
type IdentityProxy struct {
	IdentityProxyCaller     // Read-only binding to the contract
	IdentityProxyTransactor // Write-only binding to the contract
	IdentityProxyFilterer   // Log filterer for contract events
}

// IdentityProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type IdentityProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IdentityProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IdentityProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IdentityProxySession struct {
	Contract     *IdentityProxy    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IdentityProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IdentityProxyCallerSession struct {
	Contract *IdentityProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// IdentityProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IdentityProxyTransactorSession struct {
	Contract     *IdentityProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// IdentityProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type IdentityProxyRaw struct {
	Contract *IdentityProxy // Generic contract binding to access the raw methods on
}

// IdentityProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IdentityProxyCallerRaw struct {
	Contract *IdentityProxyCaller // Generic read-only contract binding to access the raw methods on
}

// IdentityProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IdentityProxyTransactorRaw struct {
	Contract *IdentityProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIdentityProxy creates a new instance of IdentityProxy, bound to a specific deployed contract.
func NewIdentityProxy(address common.Address, backend bind.ContractBackend) (*IdentityProxy, error) {
	contract, err := bindIdentityProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IdentityProxy{IdentityProxyCaller: IdentityProxyCaller{contract: contract}, IdentityProxyTransactor: IdentityProxyTransactor{contract: contract}, IdentityProxyFilterer: IdentityProxyFilterer{contract: contract}}, nil
}

// NewIdentityProxyCaller creates a new read-only instance of IdentityProxy, bound to a specific deployed contract.
func NewIdentityProxyCaller(address common.Address, caller bind.ContractCaller) (*IdentityProxyCaller, error) {
	contract, err := bindIdentityProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IdentityProxyCaller{contract: contract}, nil
}

// NewIdentityProxyTransactor creates a new write-only instance of IdentityProxy, bound to a specific deployed contract.
func NewIdentityProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*IdentityProxyTransactor, error) {
	contract, err := bindIdentityProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IdentityProxyTransactor{contract: contract}, nil
}

// NewIdentityProxyFilterer creates a new log filterer instance of IdentityProxy, bound to a specific deployed contract.
func NewIdentityProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*IdentityProxyFilterer, error) {
	contract, err := bindIdentityProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IdentityProxyFilterer{contract: contract}, nil
}

// bindIdentityProxy binds a generic wrapper to an already deployed contract.
func bindIdentityProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IdentityProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdentityProxy *IdentityProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IdentityProxy.Contract.IdentityProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdentityProxy *IdentityProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityProxy.Contract.IdentityProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdentityProxy *IdentityProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdentityProxy.Contract.IdentityProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdentityProxy *IdentityProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IdentityProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdentityProxy *IdentityProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdentityProxy *IdentityProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdentityProxy.Contract.contract.Transact(opts, method, params...)
}

// ImplementationAuthority is a free data retrieval call binding the contract method 0x2307f882.
//
// Solidity: function implementationAuthority() view returns(address)
func (_IdentityProxy *IdentityProxyCaller) ImplementationAuthority(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IdentityProxy.contract.Call(opts, &out, "implementationAuthority")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ImplementationAuthority is a free data retrieval call binding the contract method 0x2307f882.
//
// Solidity: function implementationAuthority() view returns(address)
func (_IdentityProxy *IdentityProxySession) ImplementationAuthority() (common.Address, error) {
	return _IdentityProxy.Contract.ImplementationAuthority(&_IdentityProxy.CallOpts)
}

// ImplementationAuthority is a free data retrieval call binding the contract method 0x2307f882.
//
// Solidity: function implementationAuthority() view returns(address)
func (_IdentityProxy *IdentityProxyCallerSession) ImplementationAuthority() (common.Address, error) {
	return _IdentityProxy.Contract.ImplementationAuthority(&_IdentityProxy.CallOpts)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_IdentityProxy *IdentityProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _IdentityProxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_IdentityProxy *IdentityProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _IdentityProxy.Contract.Fallback(&_IdentityProxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_IdentityProxy *IdentityProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _IdentityProxy.Contract.Fallback(&_IdentityProxy.TransactOpts, calldata)
}
