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

// IdentityRegistryStorageMetaData contains all meta data concerning the IdentityRegistryStorage contract.
var IdentityRegistryStorageMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"}],\"name\":\"AgentAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"}],\"name\":\"AgentRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"investorAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint16\",\"name\":\"country\",\"type\":\"uint16\"}],\"name\":\"CountryModified\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIIdentity\",\"name\":\"oldIdentity\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIIdentity\",\"name\":\"newIdentity\",\"type\":\"address\"}],\"name\":\"IdentityModified\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"identityRegistry\",\"type\":\"address\"}],\"name\":\"IdentityRegistryBound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"identityRegistry\",\"type\":\"address\"}],\"name\":\"IdentityRegistryUnbound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"investorAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIIdentity\",\"name\":\"identity\",\"type\":\"address\"}],\"name\":\"IdentityStored\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"investorAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIIdentity\",\"name\":\"identity\",\"type\":\"address\"}],\"name\":\"IdentityUnstored\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"}],\"name\":\"addAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userAddress\",\"type\":\"address\"},{\"internalType\":\"contractIIdentity\",\"name\":\"_identity\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"_country\",\"type\":\"uint16\"}],\"name\":\"addIdentityToStorage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_identityRegistry\",\"type\":\"address\"}],\"name\":\"bindIdentityRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"}],\"name\":\"isAgent\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkedIdentityRegistries\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userAddress\",\"type\":\"address\"},{\"internalType\":\"contractIIdentity\",\"name\":\"_identity\",\"type\":\"address\"}],\"name\":\"modifyStoredIdentity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userAddress\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"_country\",\"type\":\"uint16\"}],\"name\":\"modifyStoredInvestorCountry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"}],\"name\":\"removeAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userAddress\",\"type\":\"address\"}],\"name\":\"removeIdentityFromStorage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userAddress\",\"type\":\"address\"}],\"name\":\"storedIdentity\",\"outputs\":[{\"internalType\":\"contractIIdentity\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userAddress\",\"type\":\"address\"}],\"name\":\"storedInvestorCountry\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_identityRegistry\",\"type\":\"address\"}],\"name\":\"unbindIdentityRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061154e806100206000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c806397a6278e11610097578063cf191bcd11610066578063cf191bcd1461025a578063e1c7392a1461026d578063e805cf8614610275578063f2fde38b1461028857600080fd5b806397a6278e1461020c5780639f3418d51461021f578063a53410dd14610232578063bf9eb9591461024557600080fd5b80637988d3a5116100d35780637988d3a51461019157806384e79842146101d55780638da5cb5b146101e857806397a012f7146101f957600080fd5b80631ffbb06414610105578063690a49f91461012d578063715018a614610142578063727e13bc1461014a575b600080fd5b61011861011336600461136d565b61029b565b60405190151581526020015b60405180910390f35b61014061013b36600461136d565b6102ae565b005b610140610410565b61017e61015836600461136d565b6001600160a01b0316600090815260666020526040902054600160a01b900461ffff1690565b60405161ffff9091168152602001610124565b6101bd61019f36600461136d565b6001600160a01b039081166000908152606660205260409020541690565b6040516001600160a01b039091168152602001610124565b6101406101e336600461136d565b610424565b6033546001600160a01b03166101bd565b61014061020736600461136d565b6104c4565b61014061021a36600461136d565b6106be565b61014061022d3660046113a8565b61075e565b6101406102403660046113dd565b610904565b61024d610ab7565b6040516101249190611424565b61014061026836600461136d565b610b19565b610140610caa565b610140610283366004611471565b610dca565b61014061029636600461136d565b610f78565b60006102a8606583611005565b92915050565b6001600160a01b0381166103095760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064015b60405180910390fd5b60675461012c116103825760405162461bcd60e51b815260206004820152602560248201527f63616e6e6f742062696e64206d6f7265207468616e2033303020495220746f2060448201527f31204952530000000000000000000000000000000000000000000000000000006064820152608401610300565b61038b81610424565b6067805460018101825560009182527f9787eeb91fe3101235e4a76063c7023ecb40f923f97916639c598592fa30d6ae01805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b03841690811790915560405190917f500c250171aa20e861b680f93502547b9d436eda7d4c537fc360db6e0c6eedfb91a250565b6104186110a3565b61042260006110fd565b565b61042c6110a3565b6001600160a01b0381166104825760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f2061646472657373006044820152606401610300565b61048d60658261115c565b6040516001600160a01b038216907ff68e73cec97f2d70aa641fb26e87a4383686e2efacb648f2165aeb02ac562ec590600090a250565b6001600160a01b03811661051a5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f2061646472657373006044820152606401610300565b6067546105695760405162461bcd60e51b815260206004820152601f60248201527f6964656e74697479207265676973747279206973206e6f742073746f726564006044820152606401610300565b60675460005b8181101561067c57826001600160a01b031660678281548110610594576105946114aa565b6000918252602090912001546001600160a01b03160361066a5760676105bb6001846114d6565b815481106105cb576105cb6114aa565b600091825260209091200154606780546001600160a01b0390921691839081106105f7576105f76114aa565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b031602179055506067805480610636576106366114e9565b6000828152602090208101600019908101805473ffffffffffffffffffffffffffffffffffffffff1916905501905561067c565b80610674816114ff565b91505061056f565b50610686826106be565b6040516001600160a01b038316907f51f353eb5801583fdf2706e43c045b62fdf6b1566820b349390616363ecf72c990600090a25050565b6106c66110a3565b6001600160a01b03811661071c5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f2061646472657373006044820152606401610300565b6107276065826111d8565b6040516001600160a01b038216907fed9c8ad8d5a0a66898ea49d2956929c93ae2e8bd50281b2ed897c5d1a6737e0b90600090a250565b6107673361029b565b6107ca5760405162461bcd60e51b815260206004820152602e60248201527f4167656e74526f6c653a2063616c6c657220646f6573206e6f7420686176652060448201526d746865204167656e7420726f6c6560901b6064820152608401610300565b6001600160a01b0382166108205760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f2061646472657373006044820152606401610300565b6001600160a01b03828116600090815260666020526040902054166108875760405162461bcd60e51b815260206004820152601660248201527f61646472657373206e6f742073746f72656420796574000000000000000000006044820152606401610300565b6001600160a01b03821660008181526066602052604080822080547fffffffffffffffffffff0000ffffffffffffffffffffffffffffffffffffffff16600160a01b61ffff87169081029190911790915590519092917f20965fcdc6eed7ae398065b40ece4e732ba667992ca819fc54e80e9f2047c4cf91a35050565b61090d3361029b565b6109705760405162461bcd60e51b815260206004820152602e60248201527f4167656e74526f6c653a2063616c6c657220646f6573206e6f7420686176652060448201526d746865204167656e7420726f6c6560901b6064820152608401610300565b6001600160a01b0383161580159061099057506001600160a01b03821615155b6109dc5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f2061646472657373006044820152606401610300565b6001600160a01b038381166000908152606660205260409020541615610a445760405162461bcd60e51b815260206004820152601660248201527f616464726573732073746f72656420616c7265616479000000000000000000006044820152606401610300565b6001600160a01b03838116600081815260666020526040808220805494871675ffffffffffffffffffffffffffffffffffffffffffff199095168517600160a01b61ffff881602179055517e30dea7e9c9afaa2e3c9810f2fc9b5181f1bad74ca5a8db85f746a33585e7479190a3505050565b60606067805480602002602001604051908101604052809291908181526020018280548015610b0f57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610af1575b5050505050905090565b610b223361029b565b610b855760405162461bcd60e51b815260206004820152602e60248201527f4167656e74526f6c653a2063616c6c657220646f6573206e6f7420686176652060448201526d746865204167656e7420726f6c6560901b6064820152608401610300565b6001600160a01b038116610bdb5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f2061646472657373006044820152606401610300565b6001600160a01b0381811660009081526066602052604090205416610c425760405162461bcd60e51b815260206004820152601660248201527f61646472657373206e6f742073746f72656420796574000000000000000000006044820152606401610300565b6001600160a01b03808216600081815260666020526040808220805475ffffffffffffffffffffffffffffffffffffffffffff19811690915590519316928392917fca6a4c3370b859312246e7f086284076e557997e10d856b716c23ab67067790b91a35050565b600054610100900460ff1615808015610cca5750600054600160ff909116105b80610ce45750303b158015610ce4575060005460ff166001145b610d565760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a65640000000000000000000000000000000000006064820152608401610300565b6000805460ff191660011790558015610d79576000805461ff0019166101001790555b610d81611276565b8015610dc7576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50565b610dd33361029b565b610e365760405162461bcd60e51b815260206004820152602e60248201527f4167656e74526f6c653a2063616c6c657220646f6573206e6f7420686176652060448201526d746865204167656e7420726f6c6560901b6064820152608401610300565b6001600160a01b03821615801590610e5657506001600160a01b03811615155b610ea25760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f2061646472657373006044820152606401610300565b6001600160a01b0382811660009081526066602052604090205416610f095760405162461bcd60e51b815260206004820152601660248201527f61646472657373206e6f742073746f72656420796574000000000000000000006044820152606401610300565b6001600160a01b03808316600090815260666020526040808220805485851673ffffffffffffffffffffffffffffffffffffffff1982168117909255915191909316929183917f556ce885dfcea52155c773f1ed2e58781c51945c13030ab8f793c61f51d1b8089190a3505050565b610f806110a3565b6001600160a01b038116610ffc5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610300565b610dc7816110fd565b60006001600160a01b0382166110835760405162461bcd60e51b815260206004820152602260248201527f526f6c65733a206163636f756e7420697320746865207a65726f20616464726560448201527f73730000000000000000000000000000000000000000000000000000000000006064820152608401610300565b506001600160a01b03166000908152602091909152604090205460ff1690565b6033546001600160a01b031633146104225760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610300565b603380546001600160a01b0383811673ffffffffffffffffffffffffffffffffffffffff19831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6111668282611005565b156111b35760405162461bcd60e51b815260206004820152601f60248201527f526f6c65733a206163636f756e7420616c72656164792068617320726f6c65006044820152606401610300565b6001600160a01b0316600090815260209190915260409020805460ff19166001179055565b6111e28282611005565b6112545760405162461bcd60e51b815260206004820152602160248201527f526f6c65733a206163636f756e7420646f6573206e6f74206861766520726f6c60448201527f65000000000000000000000000000000000000000000000000000000000000006064820152608401610300565b6001600160a01b0316600090815260209190915260409020805460ff19169055565b600054610100900460ff166112e15760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b6064820152608401610300565b610422600054610100900460ff1661134f5760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b6064820152608401610300565b610422336110fd565b6001600160a01b0381168114610dc757600080fd5b60006020828403121561137f57600080fd5b813561138a81611358565b9392505050565b803561ffff811681146113a357600080fd5b919050565b600080604083850312156113bb57600080fd5b82356113c681611358565b91506113d460208401611391565b90509250929050565b6000806000606084860312156113f257600080fd5b83356113fd81611358565b9250602084013561140d81611358565b915061141b60408501611391565b90509250925092565b6020808252825182820181905260009190848201906040850190845b818110156114655783516001600160a01b031683529284019291840191600101611440565b50909695505050505050565b6000806040838503121561148457600080fd5b823561148f81611358565b9150602083013561149f81611358565b809150509250929050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b818103818111156102a8576102a86114c0565b634e487b7160e01b600052603160045260246000fd5b600060018201611511576115116114c0565b506001019056fea2646970667358221220392f78d77d545514e90cbf67dafa50fcbfa89964058cfe620a1f0a0384c1e66564736f6c63430008110033",
}

// IdentityRegistryStorageABI is the input ABI used to generate the binding from.
// Deprecated: Use IdentityRegistryStorageMetaData.ABI instead.
var IdentityRegistryStorageABI = IdentityRegistryStorageMetaData.ABI

// IdentityRegistryStorageBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use IdentityRegistryStorageMetaData.Bin instead.
var IdentityRegistryStorageBin = IdentityRegistryStorageMetaData.Bin

// DeployIdentityRegistryStorage deploys a new Ethereum contract, binding an instance of IdentityRegistryStorage to it.
func DeployIdentityRegistryStorage(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *IdentityRegistryStorage, error) {
	parsed, err := IdentityRegistryStorageMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(IdentityRegistryStorageBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &IdentityRegistryStorage{IdentityRegistryStorageCaller: IdentityRegistryStorageCaller{contract: contract}, IdentityRegistryStorageTransactor: IdentityRegistryStorageTransactor{contract: contract}, IdentityRegistryStorageFilterer: IdentityRegistryStorageFilterer{contract: contract}}, nil
}

// IdentityRegistryStorage is an auto generated Go binding around an Ethereum contract.
type IdentityRegistryStorage struct {
	IdentityRegistryStorageCaller     // Read-only binding to the contract
	IdentityRegistryStorageTransactor // Write-only binding to the contract
	IdentityRegistryStorageFilterer   // Log filterer for contract events
}

// IdentityRegistryStorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type IdentityRegistryStorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityRegistryStorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IdentityRegistryStorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityRegistryStorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IdentityRegistryStorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityRegistryStorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IdentityRegistryStorageSession struct {
	Contract     *IdentityRegistryStorage // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// IdentityRegistryStorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IdentityRegistryStorageCallerSession struct {
	Contract *IdentityRegistryStorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// IdentityRegistryStorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IdentityRegistryStorageTransactorSession struct {
	Contract     *IdentityRegistryStorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// IdentityRegistryStorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type IdentityRegistryStorageRaw struct {
	Contract *IdentityRegistryStorage // Generic contract binding to access the raw methods on
}

// IdentityRegistryStorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IdentityRegistryStorageCallerRaw struct {
	Contract *IdentityRegistryStorageCaller // Generic read-only contract binding to access the raw methods on
}

// IdentityRegistryStorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IdentityRegistryStorageTransactorRaw struct {
	Contract *IdentityRegistryStorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIdentityRegistryStorage creates a new instance of IdentityRegistryStorage, bound to a specific deployed contract.
func NewIdentityRegistryStorage(address common.Address, backend bind.ContractBackend) (*IdentityRegistryStorage, error) {
	contract, err := bindIdentityRegistryStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryStorage{IdentityRegistryStorageCaller: IdentityRegistryStorageCaller{contract: contract}, IdentityRegistryStorageTransactor: IdentityRegistryStorageTransactor{contract: contract}, IdentityRegistryStorageFilterer: IdentityRegistryStorageFilterer{contract: contract}}, nil
}

// NewIdentityRegistryStorageCaller creates a new read-only instance of IdentityRegistryStorage, bound to a specific deployed contract.
func NewIdentityRegistryStorageCaller(address common.Address, caller bind.ContractCaller) (*IdentityRegistryStorageCaller, error) {
	contract, err := bindIdentityRegistryStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryStorageCaller{contract: contract}, nil
}

// NewIdentityRegistryStorageTransactor creates a new write-only instance of IdentityRegistryStorage, bound to a specific deployed contract.
func NewIdentityRegistryStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*IdentityRegistryStorageTransactor, error) {
	contract, err := bindIdentityRegistryStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryStorageTransactor{contract: contract}, nil
}

// NewIdentityRegistryStorageFilterer creates a new log filterer instance of IdentityRegistryStorage, bound to a specific deployed contract.
func NewIdentityRegistryStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*IdentityRegistryStorageFilterer, error) {
	contract, err := bindIdentityRegistryStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryStorageFilterer{contract: contract}, nil
}

// bindIdentityRegistryStorage binds a generic wrapper to an already deployed contract.
func bindIdentityRegistryStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IdentityRegistryStorageMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdentityRegistryStorage *IdentityRegistryStorageRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IdentityRegistryStorage.Contract.IdentityRegistryStorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdentityRegistryStorage *IdentityRegistryStorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.IdentityRegistryStorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdentityRegistryStorage *IdentityRegistryStorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.IdentityRegistryStorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdentityRegistryStorage *IdentityRegistryStorageCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IdentityRegistryStorage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.contract.Transact(opts, method, params...)
}

// IsAgent is a free data retrieval call binding the contract method 0x1ffbb064.
//
// Solidity: function isAgent(address _agent) view returns(bool)
func (_IdentityRegistryStorage *IdentityRegistryStorageCaller) IsAgent(opts *bind.CallOpts, _agent common.Address) (bool, error) {
	var out []interface{}
	err := _IdentityRegistryStorage.contract.Call(opts, &out, "isAgent", _agent)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAgent is a free data retrieval call binding the contract method 0x1ffbb064.
//
// Solidity: function isAgent(address _agent) view returns(bool)
func (_IdentityRegistryStorage *IdentityRegistryStorageSession) IsAgent(_agent common.Address) (bool, error) {
	return _IdentityRegistryStorage.Contract.IsAgent(&_IdentityRegistryStorage.CallOpts, _agent)
}

// IsAgent is a free data retrieval call binding the contract method 0x1ffbb064.
//
// Solidity: function isAgent(address _agent) view returns(bool)
func (_IdentityRegistryStorage *IdentityRegistryStorageCallerSession) IsAgent(_agent common.Address) (bool, error) {
	return _IdentityRegistryStorage.Contract.IsAgent(&_IdentityRegistryStorage.CallOpts, _agent)
}

// LinkedIdentityRegistries is a free data retrieval call binding the contract method 0xbf9eb959.
//
// Solidity: function linkedIdentityRegistries() view returns(address[])
func (_IdentityRegistryStorage *IdentityRegistryStorageCaller) LinkedIdentityRegistries(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _IdentityRegistryStorage.contract.Call(opts, &out, "linkedIdentityRegistries")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// LinkedIdentityRegistries is a free data retrieval call binding the contract method 0xbf9eb959.
//
// Solidity: function linkedIdentityRegistries() view returns(address[])
func (_IdentityRegistryStorage *IdentityRegistryStorageSession) LinkedIdentityRegistries() ([]common.Address, error) {
	return _IdentityRegistryStorage.Contract.LinkedIdentityRegistries(&_IdentityRegistryStorage.CallOpts)
}

// LinkedIdentityRegistries is a free data retrieval call binding the contract method 0xbf9eb959.
//
// Solidity: function linkedIdentityRegistries() view returns(address[])
func (_IdentityRegistryStorage *IdentityRegistryStorageCallerSession) LinkedIdentityRegistries() ([]common.Address, error) {
	return _IdentityRegistryStorage.Contract.LinkedIdentityRegistries(&_IdentityRegistryStorage.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IdentityRegistryStorage *IdentityRegistryStorageCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IdentityRegistryStorage.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IdentityRegistryStorage *IdentityRegistryStorageSession) Owner() (common.Address, error) {
	return _IdentityRegistryStorage.Contract.Owner(&_IdentityRegistryStorage.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IdentityRegistryStorage *IdentityRegistryStorageCallerSession) Owner() (common.Address, error) {
	return _IdentityRegistryStorage.Contract.Owner(&_IdentityRegistryStorage.CallOpts)
}

// StoredIdentity is a free data retrieval call binding the contract method 0x7988d3a5.
//
// Solidity: function storedIdentity(address _userAddress) view returns(address)
func (_IdentityRegistryStorage *IdentityRegistryStorageCaller) StoredIdentity(opts *bind.CallOpts, _userAddress common.Address) (common.Address, error) {
	var out []interface{}
	err := _IdentityRegistryStorage.contract.Call(opts, &out, "storedIdentity", _userAddress)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StoredIdentity is a free data retrieval call binding the contract method 0x7988d3a5.
//
// Solidity: function storedIdentity(address _userAddress) view returns(address)
func (_IdentityRegistryStorage *IdentityRegistryStorageSession) StoredIdentity(_userAddress common.Address) (common.Address, error) {
	return _IdentityRegistryStorage.Contract.StoredIdentity(&_IdentityRegistryStorage.CallOpts, _userAddress)
}

// StoredIdentity is a free data retrieval call binding the contract method 0x7988d3a5.
//
// Solidity: function storedIdentity(address _userAddress) view returns(address)
func (_IdentityRegistryStorage *IdentityRegistryStorageCallerSession) StoredIdentity(_userAddress common.Address) (common.Address, error) {
	return _IdentityRegistryStorage.Contract.StoredIdentity(&_IdentityRegistryStorage.CallOpts, _userAddress)
}

// StoredInvestorCountry is a free data retrieval call binding the contract method 0x727e13bc.
//
// Solidity: function storedInvestorCountry(address _userAddress) view returns(uint16)
func (_IdentityRegistryStorage *IdentityRegistryStorageCaller) StoredInvestorCountry(opts *bind.CallOpts, _userAddress common.Address) (uint16, error) {
	var out []interface{}
	err := _IdentityRegistryStorage.contract.Call(opts, &out, "storedInvestorCountry", _userAddress)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// StoredInvestorCountry is a free data retrieval call binding the contract method 0x727e13bc.
//
// Solidity: function storedInvestorCountry(address _userAddress) view returns(uint16)
func (_IdentityRegistryStorage *IdentityRegistryStorageSession) StoredInvestorCountry(_userAddress common.Address) (uint16, error) {
	return _IdentityRegistryStorage.Contract.StoredInvestorCountry(&_IdentityRegistryStorage.CallOpts, _userAddress)
}

// StoredInvestorCountry is a free data retrieval call binding the contract method 0x727e13bc.
//
// Solidity: function storedInvestorCountry(address _userAddress) view returns(uint16)
func (_IdentityRegistryStorage *IdentityRegistryStorageCallerSession) StoredInvestorCountry(_userAddress common.Address) (uint16, error) {
	return _IdentityRegistryStorage.Contract.StoredInvestorCountry(&_IdentityRegistryStorage.CallOpts, _userAddress)
}

// AddAgent is a paid mutator transaction binding the contract method 0x84e79842.
//
// Solidity: function addAgent(address _agent) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactor) AddAgent(opts *bind.TransactOpts, _agent common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.contract.Transact(opts, "addAgent", _agent)
}

// AddAgent is a paid mutator transaction binding the contract method 0x84e79842.
//
// Solidity: function addAgent(address _agent) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageSession) AddAgent(_agent common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.AddAgent(&_IdentityRegistryStorage.TransactOpts, _agent)
}

// AddAgent is a paid mutator transaction binding the contract method 0x84e79842.
//
// Solidity: function addAgent(address _agent) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactorSession) AddAgent(_agent common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.AddAgent(&_IdentityRegistryStorage.TransactOpts, _agent)
}

// AddIdentityToStorage is a paid mutator transaction binding the contract method 0xa53410dd.
//
// Solidity: function addIdentityToStorage(address _userAddress, address _identity, uint16 _country) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactor) AddIdentityToStorage(opts *bind.TransactOpts, _userAddress common.Address, _identity common.Address, _country uint16) (*types.Transaction, error) {
	return _IdentityRegistryStorage.contract.Transact(opts, "addIdentityToStorage", _userAddress, _identity, _country)
}

// AddIdentityToStorage is a paid mutator transaction binding the contract method 0xa53410dd.
//
// Solidity: function addIdentityToStorage(address _userAddress, address _identity, uint16 _country) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageSession) AddIdentityToStorage(_userAddress common.Address, _identity common.Address, _country uint16) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.AddIdentityToStorage(&_IdentityRegistryStorage.TransactOpts, _userAddress, _identity, _country)
}

// AddIdentityToStorage is a paid mutator transaction binding the contract method 0xa53410dd.
//
// Solidity: function addIdentityToStorage(address _userAddress, address _identity, uint16 _country) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactorSession) AddIdentityToStorage(_userAddress common.Address, _identity common.Address, _country uint16) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.AddIdentityToStorage(&_IdentityRegistryStorage.TransactOpts, _userAddress, _identity, _country)
}

// BindIdentityRegistry is a paid mutator transaction binding the contract method 0x690a49f9.
//
// Solidity: function bindIdentityRegistry(address _identityRegistry) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactor) BindIdentityRegistry(opts *bind.TransactOpts, _identityRegistry common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.contract.Transact(opts, "bindIdentityRegistry", _identityRegistry)
}

// BindIdentityRegistry is a paid mutator transaction binding the contract method 0x690a49f9.
//
// Solidity: function bindIdentityRegistry(address _identityRegistry) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageSession) BindIdentityRegistry(_identityRegistry common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.BindIdentityRegistry(&_IdentityRegistryStorage.TransactOpts, _identityRegistry)
}

// BindIdentityRegistry is a paid mutator transaction binding the contract method 0x690a49f9.
//
// Solidity: function bindIdentityRegistry(address _identityRegistry) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactorSession) BindIdentityRegistry(_identityRegistry common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.BindIdentityRegistry(&_IdentityRegistryStorage.TransactOpts, _identityRegistry)
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactor) Init(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityRegistryStorage.contract.Transact(opts, "init")
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageSession) Init() (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.Init(&_IdentityRegistryStorage.TransactOpts)
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactorSession) Init() (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.Init(&_IdentityRegistryStorage.TransactOpts)
}

// ModifyStoredIdentity is a paid mutator transaction binding the contract method 0xe805cf86.
//
// Solidity: function modifyStoredIdentity(address _userAddress, address _identity) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactor) ModifyStoredIdentity(opts *bind.TransactOpts, _userAddress common.Address, _identity common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.contract.Transact(opts, "modifyStoredIdentity", _userAddress, _identity)
}

// ModifyStoredIdentity is a paid mutator transaction binding the contract method 0xe805cf86.
//
// Solidity: function modifyStoredIdentity(address _userAddress, address _identity) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageSession) ModifyStoredIdentity(_userAddress common.Address, _identity common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.ModifyStoredIdentity(&_IdentityRegistryStorage.TransactOpts, _userAddress, _identity)
}

// ModifyStoredIdentity is a paid mutator transaction binding the contract method 0xe805cf86.
//
// Solidity: function modifyStoredIdentity(address _userAddress, address _identity) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactorSession) ModifyStoredIdentity(_userAddress common.Address, _identity common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.ModifyStoredIdentity(&_IdentityRegistryStorage.TransactOpts, _userAddress, _identity)
}

// ModifyStoredInvestorCountry is a paid mutator transaction binding the contract method 0x9f3418d5.
//
// Solidity: function modifyStoredInvestorCountry(address _userAddress, uint16 _country) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactor) ModifyStoredInvestorCountry(opts *bind.TransactOpts, _userAddress common.Address, _country uint16) (*types.Transaction, error) {
	return _IdentityRegistryStorage.contract.Transact(opts, "modifyStoredInvestorCountry", _userAddress, _country)
}

// ModifyStoredInvestorCountry is a paid mutator transaction binding the contract method 0x9f3418d5.
//
// Solidity: function modifyStoredInvestorCountry(address _userAddress, uint16 _country) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageSession) ModifyStoredInvestorCountry(_userAddress common.Address, _country uint16) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.ModifyStoredInvestorCountry(&_IdentityRegistryStorage.TransactOpts, _userAddress, _country)
}

// ModifyStoredInvestorCountry is a paid mutator transaction binding the contract method 0x9f3418d5.
//
// Solidity: function modifyStoredInvestorCountry(address _userAddress, uint16 _country) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactorSession) ModifyStoredInvestorCountry(_userAddress common.Address, _country uint16) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.ModifyStoredInvestorCountry(&_IdentityRegistryStorage.TransactOpts, _userAddress, _country)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x97a6278e.
//
// Solidity: function removeAgent(address _agent) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactor) RemoveAgent(opts *bind.TransactOpts, _agent common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.contract.Transact(opts, "removeAgent", _agent)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x97a6278e.
//
// Solidity: function removeAgent(address _agent) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageSession) RemoveAgent(_agent common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.RemoveAgent(&_IdentityRegistryStorage.TransactOpts, _agent)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x97a6278e.
//
// Solidity: function removeAgent(address _agent) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactorSession) RemoveAgent(_agent common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.RemoveAgent(&_IdentityRegistryStorage.TransactOpts, _agent)
}

// RemoveIdentityFromStorage is a paid mutator transaction binding the contract method 0xcf191bcd.
//
// Solidity: function removeIdentityFromStorage(address _userAddress) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactor) RemoveIdentityFromStorage(opts *bind.TransactOpts, _userAddress common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.contract.Transact(opts, "removeIdentityFromStorage", _userAddress)
}

// RemoveIdentityFromStorage is a paid mutator transaction binding the contract method 0xcf191bcd.
//
// Solidity: function removeIdentityFromStorage(address _userAddress) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageSession) RemoveIdentityFromStorage(_userAddress common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.RemoveIdentityFromStorage(&_IdentityRegistryStorage.TransactOpts, _userAddress)
}

// RemoveIdentityFromStorage is a paid mutator transaction binding the contract method 0xcf191bcd.
//
// Solidity: function removeIdentityFromStorage(address _userAddress) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactorSession) RemoveIdentityFromStorage(_userAddress common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.RemoveIdentityFromStorage(&_IdentityRegistryStorage.TransactOpts, _userAddress)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityRegistryStorage.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageSession) RenounceOwnership() (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.RenounceOwnership(&_IdentityRegistryStorage.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.RenounceOwnership(&_IdentityRegistryStorage.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.TransferOwnership(&_IdentityRegistryStorage.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.TransferOwnership(&_IdentityRegistryStorage.TransactOpts, newOwner)
}

// UnbindIdentityRegistry is a paid mutator transaction binding the contract method 0x97a012f7.
//
// Solidity: function unbindIdentityRegistry(address _identityRegistry) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactor) UnbindIdentityRegistry(opts *bind.TransactOpts, _identityRegistry common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.contract.Transact(opts, "unbindIdentityRegistry", _identityRegistry)
}

// UnbindIdentityRegistry is a paid mutator transaction binding the contract method 0x97a012f7.
//
// Solidity: function unbindIdentityRegistry(address _identityRegistry) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageSession) UnbindIdentityRegistry(_identityRegistry common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.UnbindIdentityRegistry(&_IdentityRegistryStorage.TransactOpts, _identityRegistry)
}

// UnbindIdentityRegistry is a paid mutator transaction binding the contract method 0x97a012f7.
//
// Solidity: function unbindIdentityRegistry(address _identityRegistry) returns()
func (_IdentityRegistryStorage *IdentityRegistryStorageTransactorSession) UnbindIdentityRegistry(_identityRegistry common.Address) (*types.Transaction, error) {
	return _IdentityRegistryStorage.Contract.UnbindIdentityRegistry(&_IdentityRegistryStorage.TransactOpts, _identityRegistry)
}

// IdentityRegistryStorageAgentAddedIterator is returned from FilterAgentAdded and is used to iterate over the raw logs and unpacked data for AgentAdded events raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageAgentAddedIterator struct {
	Event *IdentityRegistryStorageAgentAdded // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryStorageAgentAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryStorageAgentAdded)
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
		it.Event = new(IdentityRegistryStorageAgentAdded)
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
func (it *IdentityRegistryStorageAgentAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryStorageAgentAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryStorageAgentAdded represents a AgentAdded event raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageAgentAdded struct {
	Agent common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterAgentAdded is a free log retrieval operation binding the contract event 0xf68e73cec97f2d70aa641fb26e87a4383686e2efacb648f2165aeb02ac562ec5.
//
// Solidity: event AgentAdded(address indexed _agent)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) FilterAgentAdded(opts *bind.FilterOpts, _agent []common.Address) (*IdentityRegistryStorageAgentAddedIterator, error) {

	var _agentRule []interface{}
	for _, _agentItem := range _agent {
		_agentRule = append(_agentRule, _agentItem)
	}

	logs, sub, err := _IdentityRegistryStorage.contract.FilterLogs(opts, "AgentAdded", _agentRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryStorageAgentAddedIterator{contract: _IdentityRegistryStorage.contract, event: "AgentAdded", logs: logs, sub: sub}, nil
}

// WatchAgentAdded is a free log subscription operation binding the contract event 0xf68e73cec97f2d70aa641fb26e87a4383686e2efacb648f2165aeb02ac562ec5.
//
// Solidity: event AgentAdded(address indexed _agent)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) WatchAgentAdded(opts *bind.WatchOpts, sink chan<- *IdentityRegistryStorageAgentAdded, _agent []common.Address) (event.Subscription, error) {

	var _agentRule []interface{}
	for _, _agentItem := range _agent {
		_agentRule = append(_agentRule, _agentItem)
	}

	logs, sub, err := _IdentityRegistryStorage.contract.WatchLogs(opts, "AgentAdded", _agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryStorageAgentAdded)
				if err := _IdentityRegistryStorage.contract.UnpackLog(event, "AgentAdded", log); err != nil {
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

// ParseAgentAdded is a log parse operation binding the contract event 0xf68e73cec97f2d70aa641fb26e87a4383686e2efacb648f2165aeb02ac562ec5.
//
// Solidity: event AgentAdded(address indexed _agent)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) ParseAgentAdded(log types.Log) (*IdentityRegistryStorageAgentAdded, error) {
	event := new(IdentityRegistryStorageAgentAdded)
	if err := _IdentityRegistryStorage.contract.UnpackLog(event, "AgentAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryStorageAgentRemovedIterator is returned from FilterAgentRemoved and is used to iterate over the raw logs and unpacked data for AgentRemoved events raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageAgentRemovedIterator struct {
	Event *IdentityRegistryStorageAgentRemoved // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryStorageAgentRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryStorageAgentRemoved)
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
		it.Event = new(IdentityRegistryStorageAgentRemoved)
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
func (it *IdentityRegistryStorageAgentRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryStorageAgentRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryStorageAgentRemoved represents a AgentRemoved event raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageAgentRemoved struct {
	Agent common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterAgentRemoved is a free log retrieval operation binding the contract event 0xed9c8ad8d5a0a66898ea49d2956929c93ae2e8bd50281b2ed897c5d1a6737e0b.
//
// Solidity: event AgentRemoved(address indexed _agent)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) FilterAgentRemoved(opts *bind.FilterOpts, _agent []common.Address) (*IdentityRegistryStorageAgentRemovedIterator, error) {

	var _agentRule []interface{}
	for _, _agentItem := range _agent {
		_agentRule = append(_agentRule, _agentItem)
	}

	logs, sub, err := _IdentityRegistryStorage.contract.FilterLogs(opts, "AgentRemoved", _agentRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryStorageAgentRemovedIterator{contract: _IdentityRegistryStorage.contract, event: "AgentRemoved", logs: logs, sub: sub}, nil
}

// WatchAgentRemoved is a free log subscription operation binding the contract event 0xed9c8ad8d5a0a66898ea49d2956929c93ae2e8bd50281b2ed897c5d1a6737e0b.
//
// Solidity: event AgentRemoved(address indexed _agent)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) WatchAgentRemoved(opts *bind.WatchOpts, sink chan<- *IdentityRegistryStorageAgentRemoved, _agent []common.Address) (event.Subscription, error) {

	var _agentRule []interface{}
	for _, _agentItem := range _agent {
		_agentRule = append(_agentRule, _agentItem)
	}

	logs, sub, err := _IdentityRegistryStorage.contract.WatchLogs(opts, "AgentRemoved", _agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryStorageAgentRemoved)
				if err := _IdentityRegistryStorage.contract.UnpackLog(event, "AgentRemoved", log); err != nil {
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

// ParseAgentRemoved is a log parse operation binding the contract event 0xed9c8ad8d5a0a66898ea49d2956929c93ae2e8bd50281b2ed897c5d1a6737e0b.
//
// Solidity: event AgentRemoved(address indexed _agent)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) ParseAgentRemoved(log types.Log) (*IdentityRegistryStorageAgentRemoved, error) {
	event := new(IdentityRegistryStorageAgentRemoved)
	if err := _IdentityRegistryStorage.contract.UnpackLog(event, "AgentRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryStorageCountryModifiedIterator is returned from FilterCountryModified and is used to iterate over the raw logs and unpacked data for CountryModified events raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageCountryModifiedIterator struct {
	Event *IdentityRegistryStorageCountryModified // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryStorageCountryModifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryStorageCountryModified)
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
		it.Event = new(IdentityRegistryStorageCountryModified)
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
func (it *IdentityRegistryStorageCountryModifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryStorageCountryModifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryStorageCountryModified represents a CountryModified event raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageCountryModified struct {
	InvestorAddress common.Address
	Country         uint16
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterCountryModified is a free log retrieval operation binding the contract event 0x20965fcdc6eed7ae398065b40ece4e732ba667992ca819fc54e80e9f2047c4cf.
//
// Solidity: event CountryModified(address indexed investorAddress, uint16 indexed country)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) FilterCountryModified(opts *bind.FilterOpts, investorAddress []common.Address, country []uint16) (*IdentityRegistryStorageCountryModifiedIterator, error) {

	var investorAddressRule []interface{}
	for _, investorAddressItem := range investorAddress {
		investorAddressRule = append(investorAddressRule, investorAddressItem)
	}
	var countryRule []interface{}
	for _, countryItem := range country {
		countryRule = append(countryRule, countryItem)
	}

	logs, sub, err := _IdentityRegistryStorage.contract.FilterLogs(opts, "CountryModified", investorAddressRule, countryRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryStorageCountryModifiedIterator{contract: _IdentityRegistryStorage.contract, event: "CountryModified", logs: logs, sub: sub}, nil
}

// WatchCountryModified is a free log subscription operation binding the contract event 0x20965fcdc6eed7ae398065b40ece4e732ba667992ca819fc54e80e9f2047c4cf.
//
// Solidity: event CountryModified(address indexed investorAddress, uint16 indexed country)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) WatchCountryModified(opts *bind.WatchOpts, sink chan<- *IdentityRegistryStorageCountryModified, investorAddress []common.Address, country []uint16) (event.Subscription, error) {

	var investorAddressRule []interface{}
	for _, investorAddressItem := range investorAddress {
		investorAddressRule = append(investorAddressRule, investorAddressItem)
	}
	var countryRule []interface{}
	for _, countryItem := range country {
		countryRule = append(countryRule, countryItem)
	}

	logs, sub, err := _IdentityRegistryStorage.contract.WatchLogs(opts, "CountryModified", investorAddressRule, countryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryStorageCountryModified)
				if err := _IdentityRegistryStorage.contract.UnpackLog(event, "CountryModified", log); err != nil {
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

// ParseCountryModified is a log parse operation binding the contract event 0x20965fcdc6eed7ae398065b40ece4e732ba667992ca819fc54e80e9f2047c4cf.
//
// Solidity: event CountryModified(address indexed investorAddress, uint16 indexed country)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) ParseCountryModified(log types.Log) (*IdentityRegistryStorageCountryModified, error) {
	event := new(IdentityRegistryStorageCountryModified)
	if err := _IdentityRegistryStorage.contract.UnpackLog(event, "CountryModified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryStorageIdentityModifiedIterator is returned from FilterIdentityModified and is used to iterate over the raw logs and unpacked data for IdentityModified events raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageIdentityModifiedIterator struct {
	Event *IdentityRegistryStorageIdentityModified // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryStorageIdentityModifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryStorageIdentityModified)
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
		it.Event = new(IdentityRegistryStorageIdentityModified)
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
func (it *IdentityRegistryStorageIdentityModifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryStorageIdentityModifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryStorageIdentityModified represents a IdentityModified event raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageIdentityModified struct {
	OldIdentity common.Address
	NewIdentity common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterIdentityModified is a free log retrieval operation binding the contract event 0x556ce885dfcea52155c773f1ed2e58781c51945c13030ab8f793c61f51d1b808.
//
// Solidity: event IdentityModified(address indexed oldIdentity, address indexed newIdentity)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) FilterIdentityModified(opts *bind.FilterOpts, oldIdentity []common.Address, newIdentity []common.Address) (*IdentityRegistryStorageIdentityModifiedIterator, error) {

	var oldIdentityRule []interface{}
	for _, oldIdentityItem := range oldIdentity {
		oldIdentityRule = append(oldIdentityRule, oldIdentityItem)
	}
	var newIdentityRule []interface{}
	for _, newIdentityItem := range newIdentity {
		newIdentityRule = append(newIdentityRule, newIdentityItem)
	}

	logs, sub, err := _IdentityRegistryStorage.contract.FilterLogs(opts, "IdentityModified", oldIdentityRule, newIdentityRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryStorageIdentityModifiedIterator{contract: _IdentityRegistryStorage.contract, event: "IdentityModified", logs: logs, sub: sub}, nil
}

// WatchIdentityModified is a free log subscription operation binding the contract event 0x556ce885dfcea52155c773f1ed2e58781c51945c13030ab8f793c61f51d1b808.
//
// Solidity: event IdentityModified(address indexed oldIdentity, address indexed newIdentity)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) WatchIdentityModified(opts *bind.WatchOpts, sink chan<- *IdentityRegistryStorageIdentityModified, oldIdentity []common.Address, newIdentity []common.Address) (event.Subscription, error) {

	var oldIdentityRule []interface{}
	for _, oldIdentityItem := range oldIdentity {
		oldIdentityRule = append(oldIdentityRule, oldIdentityItem)
	}
	var newIdentityRule []interface{}
	for _, newIdentityItem := range newIdentity {
		newIdentityRule = append(newIdentityRule, newIdentityItem)
	}

	logs, sub, err := _IdentityRegistryStorage.contract.WatchLogs(opts, "IdentityModified", oldIdentityRule, newIdentityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryStorageIdentityModified)
				if err := _IdentityRegistryStorage.contract.UnpackLog(event, "IdentityModified", log); err != nil {
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

// ParseIdentityModified is a log parse operation binding the contract event 0x556ce885dfcea52155c773f1ed2e58781c51945c13030ab8f793c61f51d1b808.
//
// Solidity: event IdentityModified(address indexed oldIdentity, address indexed newIdentity)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) ParseIdentityModified(log types.Log) (*IdentityRegistryStorageIdentityModified, error) {
	event := new(IdentityRegistryStorageIdentityModified)
	if err := _IdentityRegistryStorage.contract.UnpackLog(event, "IdentityModified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryStorageIdentityRegistryBoundIterator is returned from FilterIdentityRegistryBound and is used to iterate over the raw logs and unpacked data for IdentityRegistryBound events raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageIdentityRegistryBoundIterator struct {
	Event *IdentityRegistryStorageIdentityRegistryBound // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryStorageIdentityRegistryBoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryStorageIdentityRegistryBound)
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
		it.Event = new(IdentityRegistryStorageIdentityRegistryBound)
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
func (it *IdentityRegistryStorageIdentityRegistryBoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryStorageIdentityRegistryBoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryStorageIdentityRegistryBound represents a IdentityRegistryBound event raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageIdentityRegistryBound struct {
	IdentityRegistry common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterIdentityRegistryBound is a free log retrieval operation binding the contract event 0x500c250171aa20e861b680f93502547b9d436eda7d4c537fc360db6e0c6eedfb.
//
// Solidity: event IdentityRegistryBound(address indexed identityRegistry)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) FilterIdentityRegistryBound(opts *bind.FilterOpts, identityRegistry []common.Address) (*IdentityRegistryStorageIdentityRegistryBoundIterator, error) {

	var identityRegistryRule []interface{}
	for _, identityRegistryItem := range identityRegistry {
		identityRegistryRule = append(identityRegistryRule, identityRegistryItem)
	}

	logs, sub, err := _IdentityRegistryStorage.contract.FilterLogs(opts, "IdentityRegistryBound", identityRegistryRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryStorageIdentityRegistryBoundIterator{contract: _IdentityRegistryStorage.contract, event: "IdentityRegistryBound", logs: logs, sub: sub}, nil
}

// WatchIdentityRegistryBound is a free log subscription operation binding the contract event 0x500c250171aa20e861b680f93502547b9d436eda7d4c537fc360db6e0c6eedfb.
//
// Solidity: event IdentityRegistryBound(address indexed identityRegistry)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) WatchIdentityRegistryBound(opts *bind.WatchOpts, sink chan<- *IdentityRegistryStorageIdentityRegistryBound, identityRegistry []common.Address) (event.Subscription, error) {

	var identityRegistryRule []interface{}
	for _, identityRegistryItem := range identityRegistry {
		identityRegistryRule = append(identityRegistryRule, identityRegistryItem)
	}

	logs, sub, err := _IdentityRegistryStorage.contract.WatchLogs(opts, "IdentityRegistryBound", identityRegistryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryStorageIdentityRegistryBound)
				if err := _IdentityRegistryStorage.contract.UnpackLog(event, "IdentityRegistryBound", log); err != nil {
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

// ParseIdentityRegistryBound is a log parse operation binding the contract event 0x500c250171aa20e861b680f93502547b9d436eda7d4c537fc360db6e0c6eedfb.
//
// Solidity: event IdentityRegistryBound(address indexed identityRegistry)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) ParseIdentityRegistryBound(log types.Log) (*IdentityRegistryStorageIdentityRegistryBound, error) {
	event := new(IdentityRegistryStorageIdentityRegistryBound)
	if err := _IdentityRegistryStorage.contract.UnpackLog(event, "IdentityRegistryBound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryStorageIdentityRegistryUnboundIterator is returned from FilterIdentityRegistryUnbound and is used to iterate over the raw logs and unpacked data for IdentityRegistryUnbound events raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageIdentityRegistryUnboundIterator struct {
	Event *IdentityRegistryStorageIdentityRegistryUnbound // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryStorageIdentityRegistryUnboundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryStorageIdentityRegistryUnbound)
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
		it.Event = new(IdentityRegistryStorageIdentityRegistryUnbound)
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
func (it *IdentityRegistryStorageIdentityRegistryUnboundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryStorageIdentityRegistryUnboundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryStorageIdentityRegistryUnbound represents a IdentityRegistryUnbound event raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageIdentityRegistryUnbound struct {
	IdentityRegistry common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterIdentityRegistryUnbound is a free log retrieval operation binding the contract event 0x51f353eb5801583fdf2706e43c045b62fdf6b1566820b349390616363ecf72c9.
//
// Solidity: event IdentityRegistryUnbound(address indexed identityRegistry)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) FilterIdentityRegistryUnbound(opts *bind.FilterOpts, identityRegistry []common.Address) (*IdentityRegistryStorageIdentityRegistryUnboundIterator, error) {

	var identityRegistryRule []interface{}
	for _, identityRegistryItem := range identityRegistry {
		identityRegistryRule = append(identityRegistryRule, identityRegistryItem)
	}

	logs, sub, err := _IdentityRegistryStorage.contract.FilterLogs(opts, "IdentityRegistryUnbound", identityRegistryRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryStorageIdentityRegistryUnboundIterator{contract: _IdentityRegistryStorage.contract, event: "IdentityRegistryUnbound", logs: logs, sub: sub}, nil
}

// WatchIdentityRegistryUnbound is a free log subscription operation binding the contract event 0x51f353eb5801583fdf2706e43c045b62fdf6b1566820b349390616363ecf72c9.
//
// Solidity: event IdentityRegistryUnbound(address indexed identityRegistry)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) WatchIdentityRegistryUnbound(opts *bind.WatchOpts, sink chan<- *IdentityRegistryStorageIdentityRegistryUnbound, identityRegistry []common.Address) (event.Subscription, error) {

	var identityRegistryRule []interface{}
	for _, identityRegistryItem := range identityRegistry {
		identityRegistryRule = append(identityRegistryRule, identityRegistryItem)
	}

	logs, sub, err := _IdentityRegistryStorage.contract.WatchLogs(opts, "IdentityRegistryUnbound", identityRegistryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryStorageIdentityRegistryUnbound)
				if err := _IdentityRegistryStorage.contract.UnpackLog(event, "IdentityRegistryUnbound", log); err != nil {
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

// ParseIdentityRegistryUnbound is a log parse operation binding the contract event 0x51f353eb5801583fdf2706e43c045b62fdf6b1566820b349390616363ecf72c9.
//
// Solidity: event IdentityRegistryUnbound(address indexed identityRegistry)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) ParseIdentityRegistryUnbound(log types.Log) (*IdentityRegistryStorageIdentityRegistryUnbound, error) {
	event := new(IdentityRegistryStorageIdentityRegistryUnbound)
	if err := _IdentityRegistryStorage.contract.UnpackLog(event, "IdentityRegistryUnbound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryStorageIdentityStoredIterator is returned from FilterIdentityStored and is used to iterate over the raw logs and unpacked data for IdentityStored events raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageIdentityStoredIterator struct {
	Event *IdentityRegistryStorageIdentityStored // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryStorageIdentityStoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryStorageIdentityStored)
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
		it.Event = new(IdentityRegistryStorageIdentityStored)
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
func (it *IdentityRegistryStorageIdentityStoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryStorageIdentityStoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryStorageIdentityStored represents a IdentityStored event raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageIdentityStored struct {
	InvestorAddress common.Address
	Identity        common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterIdentityStored is a free log retrieval operation binding the contract event 0x0030dea7e9c9afaa2e3c9810f2fc9b5181f1bad74ca5a8db85f746a33585e747.
//
// Solidity: event IdentityStored(address indexed investorAddress, address indexed identity)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) FilterIdentityStored(opts *bind.FilterOpts, investorAddress []common.Address, identity []common.Address) (*IdentityRegistryStorageIdentityStoredIterator, error) {

	var investorAddressRule []interface{}
	for _, investorAddressItem := range investorAddress {
		investorAddressRule = append(investorAddressRule, investorAddressItem)
	}
	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _IdentityRegistryStorage.contract.FilterLogs(opts, "IdentityStored", investorAddressRule, identityRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryStorageIdentityStoredIterator{contract: _IdentityRegistryStorage.contract, event: "IdentityStored", logs: logs, sub: sub}, nil
}

// WatchIdentityStored is a free log subscription operation binding the contract event 0x0030dea7e9c9afaa2e3c9810f2fc9b5181f1bad74ca5a8db85f746a33585e747.
//
// Solidity: event IdentityStored(address indexed investorAddress, address indexed identity)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) WatchIdentityStored(opts *bind.WatchOpts, sink chan<- *IdentityRegistryStorageIdentityStored, investorAddress []common.Address, identity []common.Address) (event.Subscription, error) {

	var investorAddressRule []interface{}
	for _, investorAddressItem := range investorAddress {
		investorAddressRule = append(investorAddressRule, investorAddressItem)
	}
	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _IdentityRegistryStorage.contract.WatchLogs(opts, "IdentityStored", investorAddressRule, identityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryStorageIdentityStored)
				if err := _IdentityRegistryStorage.contract.UnpackLog(event, "IdentityStored", log); err != nil {
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

// ParseIdentityStored is a log parse operation binding the contract event 0x0030dea7e9c9afaa2e3c9810f2fc9b5181f1bad74ca5a8db85f746a33585e747.
//
// Solidity: event IdentityStored(address indexed investorAddress, address indexed identity)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) ParseIdentityStored(log types.Log) (*IdentityRegistryStorageIdentityStored, error) {
	event := new(IdentityRegistryStorageIdentityStored)
	if err := _IdentityRegistryStorage.contract.UnpackLog(event, "IdentityStored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryStorageIdentityUnstoredIterator is returned from FilterIdentityUnstored and is used to iterate over the raw logs and unpacked data for IdentityUnstored events raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageIdentityUnstoredIterator struct {
	Event *IdentityRegistryStorageIdentityUnstored // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryStorageIdentityUnstoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryStorageIdentityUnstored)
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
		it.Event = new(IdentityRegistryStorageIdentityUnstored)
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
func (it *IdentityRegistryStorageIdentityUnstoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryStorageIdentityUnstoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryStorageIdentityUnstored represents a IdentityUnstored event raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageIdentityUnstored struct {
	InvestorAddress common.Address
	Identity        common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterIdentityUnstored is a free log retrieval operation binding the contract event 0xca6a4c3370b859312246e7f086284076e557997e10d856b716c23ab67067790b.
//
// Solidity: event IdentityUnstored(address indexed investorAddress, address indexed identity)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) FilterIdentityUnstored(opts *bind.FilterOpts, investorAddress []common.Address, identity []common.Address) (*IdentityRegistryStorageIdentityUnstoredIterator, error) {

	var investorAddressRule []interface{}
	for _, investorAddressItem := range investorAddress {
		investorAddressRule = append(investorAddressRule, investorAddressItem)
	}
	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _IdentityRegistryStorage.contract.FilterLogs(opts, "IdentityUnstored", investorAddressRule, identityRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryStorageIdentityUnstoredIterator{contract: _IdentityRegistryStorage.contract, event: "IdentityUnstored", logs: logs, sub: sub}, nil
}

// WatchIdentityUnstored is a free log subscription operation binding the contract event 0xca6a4c3370b859312246e7f086284076e557997e10d856b716c23ab67067790b.
//
// Solidity: event IdentityUnstored(address indexed investorAddress, address indexed identity)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) WatchIdentityUnstored(opts *bind.WatchOpts, sink chan<- *IdentityRegistryStorageIdentityUnstored, investorAddress []common.Address, identity []common.Address) (event.Subscription, error) {

	var investorAddressRule []interface{}
	for _, investorAddressItem := range investorAddress {
		investorAddressRule = append(investorAddressRule, investorAddressItem)
	}
	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _IdentityRegistryStorage.contract.WatchLogs(opts, "IdentityUnstored", investorAddressRule, identityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryStorageIdentityUnstored)
				if err := _IdentityRegistryStorage.contract.UnpackLog(event, "IdentityUnstored", log); err != nil {
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

// ParseIdentityUnstored is a log parse operation binding the contract event 0xca6a4c3370b859312246e7f086284076e557997e10d856b716c23ab67067790b.
//
// Solidity: event IdentityUnstored(address indexed investorAddress, address indexed identity)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) ParseIdentityUnstored(log types.Log) (*IdentityRegistryStorageIdentityUnstored, error) {
	event := new(IdentityRegistryStorageIdentityUnstored)
	if err := _IdentityRegistryStorage.contract.UnpackLog(event, "IdentityUnstored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryStorageInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageInitializedIterator struct {
	Event *IdentityRegistryStorageInitialized // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryStorageInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryStorageInitialized)
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
		it.Event = new(IdentityRegistryStorageInitialized)
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
func (it *IdentityRegistryStorageInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryStorageInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryStorageInitialized represents a Initialized event raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) FilterInitialized(opts *bind.FilterOpts) (*IdentityRegistryStorageInitializedIterator, error) {

	logs, sub, err := _IdentityRegistryStorage.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryStorageInitializedIterator{contract: _IdentityRegistryStorage.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *IdentityRegistryStorageInitialized) (event.Subscription, error) {

	logs, sub, err := _IdentityRegistryStorage.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryStorageInitialized)
				if err := _IdentityRegistryStorage.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) ParseInitialized(log types.Log) (*IdentityRegistryStorageInitialized, error) {
	event := new(IdentityRegistryStorageInitialized)
	if err := _IdentityRegistryStorage.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryStorageOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageOwnershipTransferredIterator struct {
	Event *IdentityRegistryStorageOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryStorageOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryStorageOwnershipTransferred)
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
		it.Event = new(IdentityRegistryStorageOwnershipTransferred)
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
func (it *IdentityRegistryStorageOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryStorageOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryStorageOwnershipTransferred represents a OwnershipTransferred event raised by the IdentityRegistryStorage contract.
type IdentityRegistryStorageOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*IdentityRegistryStorageOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IdentityRegistryStorage.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryStorageOwnershipTransferredIterator{contract: _IdentityRegistryStorage.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IdentityRegistryStorageOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IdentityRegistryStorage.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryStorageOwnershipTransferred)
				if err := _IdentityRegistryStorage.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_IdentityRegistryStorage *IdentityRegistryStorageFilterer) ParseOwnershipTransferred(log types.Log) (*IdentityRegistryStorageOwnershipTransferred, error) {
	event := new(IdentityRegistryStorageOwnershipTransferred)
	if err := _IdentityRegistryStorage.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
