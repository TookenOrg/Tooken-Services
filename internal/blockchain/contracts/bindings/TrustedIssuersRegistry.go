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

// TrustedIssuersRegistryMetaData contains all meta data concerning the TrustedIssuersRegistry contract.
var TrustedIssuersRegistryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIClaimIssuer\",\"name\":\"trustedIssuer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"claimTopics\",\"type\":\"uint256[]\"}],\"name\":\"ClaimTopicsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIClaimIssuer\",\"name\":\"trustedIssuer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"claimTopics\",\"type\":\"uint256[]\"}],\"name\":\"TrustedIssuerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIClaimIssuer\",\"name\":\"trustedIssuer\",\"type\":\"address\"}],\"name\":\"TrustedIssuerRemoved\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"contractIClaimIssuer\",\"name\":\"_trustedIssuer\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"_claimTopics\",\"type\":\"uint256[]\"}],\"name\":\"addTrustedIssuer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIClaimIssuer\",\"name\":\"_trustedIssuer\",\"type\":\"address\"}],\"name\":\"getTrustedIssuerClaimTopics\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTrustedIssuers\",\"outputs\":[{\"internalType\":\"contractIClaimIssuer[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"claimTopic\",\"type\":\"uint256\"}],\"name\":\"getTrustedIssuersForClaimTopic\",\"outputs\":[{\"internalType\":\"contractIClaimIssuer[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_issuer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_claimTopic\",\"type\":\"uint256\"}],\"name\":\"hasClaimTopic\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_issuer\",\"type\":\"address\"}],\"name\":\"isTrustedIssuer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIClaimIssuer\",\"name\":\"_trustedIssuer\",\"type\":\"address\"}],\"name\":\"removeTrustedIssuer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIClaimIssuer\",\"name\":\"_trustedIssuer\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"_claimTopics\",\"type\":\"uint256[]\"}],\"name\":\"updateIssuerClaimTopics\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506115b9806100206000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c8063b93d28eb11610081578063e1c7392a1161005b578063e1c7392a146101a7578063ef2ed1a4146101af578063f2fde38b146101c257600080fd5b8063b93d28eb1461016c578063c28fb2781461017f578063d9dd24c51461019f57600080fd5b8063715018a6116100b2578063715018a6146101365780638da5cb5b1461013e5780639f63ea981461015957600080fd5b806304bc7e84146100d957806334a89987146100ee57806352c111d114610116575b600080fd5b6100ec6100e736600461134d565b6101d5565b005b6101016100fc3660046113d5565b6105e4565b60405190151581526020015b60405180910390f35b610129610124366004611401565b6106a4565b60405161010d919061141a565b6100ec610710565b6033546040516001600160a01b03909116815260200161010d565b6100ec61016736600461134d565b610724565b6100ec61017a366004611467565b610a47565b61019261018d366004611467565b610e21565b60405161010d919061148b565b610129610ef3565b6100ec610f55565b6101016101bd366004611467565b611075565b6100ec6101d0366004611467565b6110a3565b6101dd611130565b6001600160a01b0383166102385760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064015b60405180910390fd5b6001600160a01b038316600090815260666020526040812054900361029f5760405162461bcd60e51b815260206004820152601460248201527f4e4f542061207472757374656420697373756572000000000000000000000000604482015260640161022f565b600f8111156102fe5760405162461bcd60e51b815260206004820152602560248201527f63616e6e6f742068617665206d6f7265207468616e20313520636c61696d20746044820152646f7069637360d81b606482015260840161022f565b8061034b5760405162461bcd60e51b815260206004820152601c60248201527f636c61696d20746f706963732063616e6e6f7420626520656d70747900000000604482015260640161022f565b60005b6001600160a01b0384166000908152606660205260409020548110156104ff576001600160a01b0384166000908152606660205260408120805483908110610398576103986114c3565b600091825260208083209091015480835260679091526040822054909250905b818110156104e957600083815260676020526040902080546001600160a01b0389169190839081106103ec576103ec6114c3565b6000918252602090912001546001600160a01b0316036104d757600083815260676020526040902061041f6001846114ef565b8154811061042f5761042f6114c3565b60009182526020808320909101548583526067909152604090912080546001600160a01b039092169183908110610468576104686114c3565b600091825260208083209190910180546001600160a01b0319166001600160a01b0394909416939093179092558481526067909152604090208054806104b0576104b0611502565b600082815260209020810160001990810180546001600160a01b03191690550190556104e9565b806104e181611518565b9150506103b8565b50505080806104f790611518565b91505061034e565b506001600160a01b03831660009081526066602052604090206105239083836112be565b5060005b8181101561059b5760676000848484818110610545576105456114c3565b6020908102929092013583525081810192909252604001600090812080546001810182559082529190200180546001600160a01b0319166001600160a01b0386161790558061059381611518565b915050610527565b50826001600160a01b03167fec753cfc52044f61676f18a11e500093a9f2b1cd5e4942bc476f2b0438159bcf83836040516105d7929190611531565b60405180910390a2505050565b6001600160a01b038216600090815260666020908152604080832080548251818502810185019093528083529284929190849083018282801561064657602002820191906000526020600020905b815481526020019060010190808311610632575b5050505050905060005b82811015610696578482828151811061066b5761066b6114c3565b602002602001015103610684576001935050505061069e565b8061068e81611518565b915050610650565b506000925050505b92915050565b60008181526067602090815260409182902080548351818402810184019094528084526060939283018282801561070457602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116106e6575b50505050509050919050565b610718611130565b610722600061118a565b565b61072c611130565b6001600160a01b0383166107825760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f206164647265737300604482015260640161022f565b6001600160a01b038316600090815260666020526040902054156107e85760405162461bcd60e51b815260206004820152601d60248201527f747275737465642049737375657220616c726561647920657869737473000000604482015260640161022f565b8061085a5760405162461bcd60e51b8152602060048201526024808201527f7472757374656420636c61696d20746f706963732063616e6e6f74206265206560448201527f6d70747900000000000000000000000000000000000000000000000000000000606482015260840161022f565b600f8111156108b95760405162461bcd60e51b815260206004820152602560248201527f63616e6e6f742068617665206d6f7265207468616e20313520636c61696d20746044820152646f7069637360d81b606482015260840161022f565b6065546032116109315760405162461bcd60e51b815260206004820152602860248201527f63616e6e6f742068617665206d6f7265207468616e203530207472757374656460448201527f2069737375657273000000000000000000000000000000000000000000000000606482015260840161022f565b60658054600181019091557f8ff97419363ffd7000167f130ef7168fbea05faf9251824ca5043f113cc6a7c70180546001600160a01b0319166001600160a01b03851690811790915560009081526066602052604090206109939083836112be565b5060005b81811015610a0b57606760008484848181106109b5576109b56114c3565b6020908102929092013583525081810192909252604001600090812080546001810182559082529190200180546001600160a01b0319166001600160a01b03861617905580610a0381611518565b915050610997565b50826001600160a01b03167ffedc33fd34859594822c0ff6f3f4f9fc279cc6d5cae53068f706a088e450087283836040516105d7929190611531565b610a4f611130565b6001600160a01b038116610aa55760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f206164647265737300604482015260640161022f565b6001600160a01b0381166000908152606660205260408120549003610b0c5760405162461bcd60e51b815260206004820152601460248201527f4e4f542061207472757374656420697373756572000000000000000000000000604482015260640161022f565b60655460005b81811015610c1257826001600160a01b031660658281548110610b3757610b376114c3565b6000918252602090912001546001600160a01b031603610c00576065610b5e6001846114ef565b81548110610b6e57610b6e6114c3565b600091825260209091200154606580546001600160a01b039092169183908110610b9a57610b9a6114c3565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b031602179055506065805480610bd957610bd9611502565b600082815260209020810160001990810180546001600160a01b0319169055019055610c12565b80610c0a81611518565b915050610b12565b5060005b6001600160a01b038316600090815260666020526040902054811015610dc7576001600160a01b0383166000908152606660205260408120805483908110610c6057610c606114c3565b600091825260208083209091015480835260679091526040822054909250905b81811015610db157600083815260676020526040902080546001600160a01b038816919083908110610cb457610cb46114c3565b6000918252602090912001546001600160a01b031603610d9f576000838152606760205260409020610ce76001846114ef565b81548110610cf757610cf76114c3565b60009182526020808320909101548583526067909152604090912080546001600160a01b039092169183908110610d3057610d306114c3565b600091825260208083209190910180546001600160a01b0319166001600160a01b039490941693909317909255848152606790915260409020805480610d7857610d78611502565b600082815260209020810160001990810180546001600160a01b0319169055019055610db1565b80610da981611518565b915050610c80565b5050508080610dbf90611518565b915050610c16565b506001600160a01b0382166000908152606660205260408120610de991611309565b6040516001600160a01b038316907f2214ded40113cc3fb63fc206cafee88270b0a903dac7245d54efdde30ebb032190600090a25050565b6001600160a01b03811660009081526066602052604081205460609103610e8a5760405162461bcd60e51b815260206004820152601c60248201527f747275737465642049737375657220646f65736e277420657869737400000000604482015260640161022f565b6001600160a01b0382166000908152606660209081526040918290208054835181840281018401909452808452909183018282801561070457602002820191906000526020600020905b815481526020019060010190808311610ed45750505050509050919050565b60606065805480602002602001604051908101604052809291908181526020018280548015610f4b57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610f2d575b5050505050905090565b600054610100900460ff1615808015610f755750600054600160ff909116105b80610f8f5750303b158015610f8f575060005460ff166001145b6110015760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a6564000000000000000000000000000000000000606482015260840161022f565b6000805460ff191660011790558015611024576000805461ff0019166101001790555b61102c6111dc565b8015611072576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50565b6001600160a01b0381166000908152606660205260408120541561109b57506001919050565b506000919050565b6110ab611130565b6001600160a01b0381166111275760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f6464726573730000000000000000000000000000000000000000000000000000606482015260840161022f565b6110728161118a565b6033546001600160a01b031633146107225760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015260640161022f565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166112475760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b606482015260840161022f565b610722600054610100900460ff166112b55760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b606482015260840161022f565b6107223361118a565b8280548282559060005260206000209081019282156112f9579160200282015b828111156112f95782358255916020019190600101906112de565b50611305929150611323565b5090565b508054600082559060005260206000209081019061107291905b5b808211156113055760008155600101611324565b6001600160a01b038116811461107257600080fd5b60008060006040848603121561136257600080fd5b833561136d81611338565b9250602084013567ffffffffffffffff8082111561138a57600080fd5b818601915086601f83011261139e57600080fd5b8135818111156113ad57600080fd5b8760208260051b85010111156113c257600080fd5b6020830194508093505050509250925092565b600080604083850312156113e857600080fd5b82356113f381611338565b946020939093013593505050565b60006020828403121561141357600080fd5b5035919050565b6020808252825182820181905260009190848201906040850190845b8181101561145b5783516001600160a01b031683529284019291840191600101611436565b50909695505050505050565b60006020828403121561147957600080fd5b813561148481611338565b9392505050565b6020808252825182820181905260009190848201906040850190845b8181101561145b578351835292840192918401916001016114a7565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b8181038181111561069e5761069e6114d9565b634e487b7160e01b600052603160045260246000fd5b60006001820161152a5761152a6114d9565b5060010190565b6020815281602082015260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83111561156a57600080fd5b8260051b8085604085013791909101604001939250505056fea2646970667358221220e263289b5a88a67b98519b9adbc7b01b35a391035e7c524e9cf9863631e33da264736f6c63430008110033",
}

// TrustedIssuersRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use TrustedIssuersRegistryMetaData.ABI instead.
var TrustedIssuersRegistryABI = TrustedIssuersRegistryMetaData.ABI

// TrustedIssuersRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TrustedIssuersRegistryMetaData.Bin instead.
var TrustedIssuersRegistryBin = TrustedIssuersRegistryMetaData.Bin

// DeployTrustedIssuersRegistry deploys a new Ethereum contract, binding an instance of TrustedIssuersRegistry to it.
func DeployTrustedIssuersRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TrustedIssuersRegistry, error) {
	parsed, err := TrustedIssuersRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TrustedIssuersRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TrustedIssuersRegistry{TrustedIssuersRegistryCaller: TrustedIssuersRegistryCaller{contract: contract}, TrustedIssuersRegistryTransactor: TrustedIssuersRegistryTransactor{contract: contract}, TrustedIssuersRegistryFilterer: TrustedIssuersRegistryFilterer{contract: contract}}, nil
}

// TrustedIssuersRegistry is an auto generated Go binding around an Ethereum contract.
type TrustedIssuersRegistry struct {
	TrustedIssuersRegistryCaller     // Read-only binding to the contract
	TrustedIssuersRegistryTransactor // Write-only binding to the contract
	TrustedIssuersRegistryFilterer   // Log filterer for contract events
}

// TrustedIssuersRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type TrustedIssuersRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TrustedIssuersRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TrustedIssuersRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TrustedIssuersRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TrustedIssuersRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TrustedIssuersRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TrustedIssuersRegistrySession struct {
	Contract     *TrustedIssuersRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// TrustedIssuersRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TrustedIssuersRegistryCallerSession struct {
	Contract *TrustedIssuersRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// TrustedIssuersRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TrustedIssuersRegistryTransactorSession struct {
	Contract     *TrustedIssuersRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// TrustedIssuersRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type TrustedIssuersRegistryRaw struct {
	Contract *TrustedIssuersRegistry // Generic contract binding to access the raw methods on
}

// TrustedIssuersRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TrustedIssuersRegistryCallerRaw struct {
	Contract *TrustedIssuersRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// TrustedIssuersRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TrustedIssuersRegistryTransactorRaw struct {
	Contract *TrustedIssuersRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTrustedIssuersRegistry creates a new instance of TrustedIssuersRegistry, bound to a specific deployed contract.
func NewTrustedIssuersRegistry(address common.Address, backend bind.ContractBackend) (*TrustedIssuersRegistry, error) {
	contract, err := bindTrustedIssuersRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TrustedIssuersRegistry{TrustedIssuersRegistryCaller: TrustedIssuersRegistryCaller{contract: contract}, TrustedIssuersRegistryTransactor: TrustedIssuersRegistryTransactor{contract: contract}, TrustedIssuersRegistryFilterer: TrustedIssuersRegistryFilterer{contract: contract}}, nil
}

// NewTrustedIssuersRegistryCaller creates a new read-only instance of TrustedIssuersRegistry, bound to a specific deployed contract.
func NewTrustedIssuersRegistryCaller(address common.Address, caller bind.ContractCaller) (*TrustedIssuersRegistryCaller, error) {
	contract, err := bindTrustedIssuersRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TrustedIssuersRegistryCaller{contract: contract}, nil
}

// NewTrustedIssuersRegistryTransactor creates a new write-only instance of TrustedIssuersRegistry, bound to a specific deployed contract.
func NewTrustedIssuersRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*TrustedIssuersRegistryTransactor, error) {
	contract, err := bindTrustedIssuersRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TrustedIssuersRegistryTransactor{contract: contract}, nil
}

// NewTrustedIssuersRegistryFilterer creates a new log filterer instance of TrustedIssuersRegistry, bound to a specific deployed contract.
func NewTrustedIssuersRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*TrustedIssuersRegistryFilterer, error) {
	contract, err := bindTrustedIssuersRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TrustedIssuersRegistryFilterer{contract: contract}, nil
}

// bindTrustedIssuersRegistry binds a generic wrapper to an already deployed contract.
func bindTrustedIssuersRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TrustedIssuersRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TrustedIssuersRegistry *TrustedIssuersRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TrustedIssuersRegistry.Contract.TrustedIssuersRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TrustedIssuersRegistry *TrustedIssuersRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TrustedIssuersRegistry.Contract.TrustedIssuersRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TrustedIssuersRegistry *TrustedIssuersRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TrustedIssuersRegistry.Contract.TrustedIssuersRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TrustedIssuersRegistry *TrustedIssuersRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TrustedIssuersRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TrustedIssuersRegistry *TrustedIssuersRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TrustedIssuersRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TrustedIssuersRegistry *TrustedIssuersRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TrustedIssuersRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetTrustedIssuerClaimTopics is a free data retrieval call binding the contract method 0xc28fb278.
//
// Solidity: function getTrustedIssuerClaimTopics(address _trustedIssuer) view returns(uint256[])
func (_TrustedIssuersRegistry *TrustedIssuersRegistryCaller) GetTrustedIssuerClaimTopics(opts *bind.CallOpts, _trustedIssuer common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _TrustedIssuersRegistry.contract.Call(opts, &out, "getTrustedIssuerClaimTopics", _trustedIssuer)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetTrustedIssuerClaimTopics is a free data retrieval call binding the contract method 0xc28fb278.
//
// Solidity: function getTrustedIssuerClaimTopics(address _trustedIssuer) view returns(uint256[])
func (_TrustedIssuersRegistry *TrustedIssuersRegistrySession) GetTrustedIssuerClaimTopics(_trustedIssuer common.Address) ([]*big.Int, error) {
	return _TrustedIssuersRegistry.Contract.GetTrustedIssuerClaimTopics(&_TrustedIssuersRegistry.CallOpts, _trustedIssuer)
}

// GetTrustedIssuerClaimTopics is a free data retrieval call binding the contract method 0xc28fb278.
//
// Solidity: function getTrustedIssuerClaimTopics(address _trustedIssuer) view returns(uint256[])
func (_TrustedIssuersRegistry *TrustedIssuersRegistryCallerSession) GetTrustedIssuerClaimTopics(_trustedIssuer common.Address) ([]*big.Int, error) {
	return _TrustedIssuersRegistry.Contract.GetTrustedIssuerClaimTopics(&_TrustedIssuersRegistry.CallOpts, _trustedIssuer)
}

// GetTrustedIssuers is a free data retrieval call binding the contract method 0xd9dd24c5.
//
// Solidity: function getTrustedIssuers() view returns(address[])
func (_TrustedIssuersRegistry *TrustedIssuersRegistryCaller) GetTrustedIssuers(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _TrustedIssuersRegistry.contract.Call(opts, &out, "getTrustedIssuers")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetTrustedIssuers is a free data retrieval call binding the contract method 0xd9dd24c5.
//
// Solidity: function getTrustedIssuers() view returns(address[])
func (_TrustedIssuersRegistry *TrustedIssuersRegistrySession) GetTrustedIssuers() ([]common.Address, error) {
	return _TrustedIssuersRegistry.Contract.GetTrustedIssuers(&_TrustedIssuersRegistry.CallOpts)
}

// GetTrustedIssuers is a free data retrieval call binding the contract method 0xd9dd24c5.
//
// Solidity: function getTrustedIssuers() view returns(address[])
func (_TrustedIssuersRegistry *TrustedIssuersRegistryCallerSession) GetTrustedIssuers() ([]common.Address, error) {
	return _TrustedIssuersRegistry.Contract.GetTrustedIssuers(&_TrustedIssuersRegistry.CallOpts)
}

// GetTrustedIssuersForClaimTopic is a free data retrieval call binding the contract method 0x52c111d1.
//
// Solidity: function getTrustedIssuersForClaimTopic(uint256 claimTopic) view returns(address[])
func (_TrustedIssuersRegistry *TrustedIssuersRegistryCaller) GetTrustedIssuersForClaimTopic(opts *bind.CallOpts, claimTopic *big.Int) ([]common.Address, error) {
	var out []interface{}
	err := _TrustedIssuersRegistry.contract.Call(opts, &out, "getTrustedIssuersForClaimTopic", claimTopic)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetTrustedIssuersForClaimTopic is a free data retrieval call binding the contract method 0x52c111d1.
//
// Solidity: function getTrustedIssuersForClaimTopic(uint256 claimTopic) view returns(address[])
func (_TrustedIssuersRegistry *TrustedIssuersRegistrySession) GetTrustedIssuersForClaimTopic(claimTopic *big.Int) ([]common.Address, error) {
	return _TrustedIssuersRegistry.Contract.GetTrustedIssuersForClaimTopic(&_TrustedIssuersRegistry.CallOpts, claimTopic)
}

// GetTrustedIssuersForClaimTopic is a free data retrieval call binding the contract method 0x52c111d1.
//
// Solidity: function getTrustedIssuersForClaimTopic(uint256 claimTopic) view returns(address[])
func (_TrustedIssuersRegistry *TrustedIssuersRegistryCallerSession) GetTrustedIssuersForClaimTopic(claimTopic *big.Int) ([]common.Address, error) {
	return _TrustedIssuersRegistry.Contract.GetTrustedIssuersForClaimTopic(&_TrustedIssuersRegistry.CallOpts, claimTopic)
}

// HasClaimTopic is a free data retrieval call binding the contract method 0x34a89987.
//
// Solidity: function hasClaimTopic(address _issuer, uint256 _claimTopic) view returns(bool)
func (_TrustedIssuersRegistry *TrustedIssuersRegistryCaller) HasClaimTopic(opts *bind.CallOpts, _issuer common.Address, _claimTopic *big.Int) (bool, error) {
	var out []interface{}
	err := _TrustedIssuersRegistry.contract.Call(opts, &out, "hasClaimTopic", _issuer, _claimTopic)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasClaimTopic is a free data retrieval call binding the contract method 0x34a89987.
//
// Solidity: function hasClaimTopic(address _issuer, uint256 _claimTopic) view returns(bool)
func (_TrustedIssuersRegistry *TrustedIssuersRegistrySession) HasClaimTopic(_issuer common.Address, _claimTopic *big.Int) (bool, error) {
	return _TrustedIssuersRegistry.Contract.HasClaimTopic(&_TrustedIssuersRegistry.CallOpts, _issuer, _claimTopic)
}

// HasClaimTopic is a free data retrieval call binding the contract method 0x34a89987.
//
// Solidity: function hasClaimTopic(address _issuer, uint256 _claimTopic) view returns(bool)
func (_TrustedIssuersRegistry *TrustedIssuersRegistryCallerSession) HasClaimTopic(_issuer common.Address, _claimTopic *big.Int) (bool, error) {
	return _TrustedIssuersRegistry.Contract.HasClaimTopic(&_TrustedIssuersRegistry.CallOpts, _issuer, _claimTopic)
}

// IsTrustedIssuer is a free data retrieval call binding the contract method 0xef2ed1a4.
//
// Solidity: function isTrustedIssuer(address _issuer) view returns(bool)
func (_TrustedIssuersRegistry *TrustedIssuersRegistryCaller) IsTrustedIssuer(opts *bind.CallOpts, _issuer common.Address) (bool, error) {
	var out []interface{}
	err := _TrustedIssuersRegistry.contract.Call(opts, &out, "isTrustedIssuer", _issuer)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsTrustedIssuer is a free data retrieval call binding the contract method 0xef2ed1a4.
//
// Solidity: function isTrustedIssuer(address _issuer) view returns(bool)
func (_TrustedIssuersRegistry *TrustedIssuersRegistrySession) IsTrustedIssuer(_issuer common.Address) (bool, error) {
	return _TrustedIssuersRegistry.Contract.IsTrustedIssuer(&_TrustedIssuersRegistry.CallOpts, _issuer)
}

// IsTrustedIssuer is a free data retrieval call binding the contract method 0xef2ed1a4.
//
// Solidity: function isTrustedIssuer(address _issuer) view returns(bool)
func (_TrustedIssuersRegistry *TrustedIssuersRegistryCallerSession) IsTrustedIssuer(_issuer common.Address) (bool, error) {
	return _TrustedIssuersRegistry.Contract.IsTrustedIssuer(&_TrustedIssuersRegistry.CallOpts, _issuer)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TrustedIssuersRegistry *TrustedIssuersRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TrustedIssuersRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TrustedIssuersRegistry *TrustedIssuersRegistrySession) Owner() (common.Address, error) {
	return _TrustedIssuersRegistry.Contract.Owner(&_TrustedIssuersRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TrustedIssuersRegistry *TrustedIssuersRegistryCallerSession) Owner() (common.Address, error) {
	return _TrustedIssuersRegistry.Contract.Owner(&_TrustedIssuersRegistry.CallOpts)
}

// AddTrustedIssuer is a paid mutator transaction binding the contract method 0x9f63ea98.
//
// Solidity: function addTrustedIssuer(address _trustedIssuer, uint256[] _claimTopics) returns()
func (_TrustedIssuersRegistry *TrustedIssuersRegistryTransactor) AddTrustedIssuer(opts *bind.TransactOpts, _trustedIssuer common.Address, _claimTopics []*big.Int) (*types.Transaction, error) {
	return _TrustedIssuersRegistry.contract.Transact(opts, "addTrustedIssuer", _trustedIssuer, _claimTopics)
}

// AddTrustedIssuer is a paid mutator transaction binding the contract method 0x9f63ea98.
//
// Solidity: function addTrustedIssuer(address _trustedIssuer, uint256[] _claimTopics) returns()
func (_TrustedIssuersRegistry *TrustedIssuersRegistrySession) AddTrustedIssuer(_trustedIssuer common.Address, _claimTopics []*big.Int) (*types.Transaction, error) {
	return _TrustedIssuersRegistry.Contract.AddTrustedIssuer(&_TrustedIssuersRegistry.TransactOpts, _trustedIssuer, _claimTopics)
}

// AddTrustedIssuer is a paid mutator transaction binding the contract method 0x9f63ea98.
//
// Solidity: function addTrustedIssuer(address _trustedIssuer, uint256[] _claimTopics) returns()
func (_TrustedIssuersRegistry *TrustedIssuersRegistryTransactorSession) AddTrustedIssuer(_trustedIssuer common.Address, _claimTopics []*big.Int) (*types.Transaction, error) {
	return _TrustedIssuersRegistry.Contract.AddTrustedIssuer(&_TrustedIssuersRegistry.TransactOpts, _trustedIssuer, _claimTopics)
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_TrustedIssuersRegistry *TrustedIssuersRegistryTransactor) Init(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TrustedIssuersRegistry.contract.Transact(opts, "init")
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_TrustedIssuersRegistry *TrustedIssuersRegistrySession) Init() (*types.Transaction, error) {
	return _TrustedIssuersRegistry.Contract.Init(&_TrustedIssuersRegistry.TransactOpts)
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_TrustedIssuersRegistry *TrustedIssuersRegistryTransactorSession) Init() (*types.Transaction, error) {
	return _TrustedIssuersRegistry.Contract.Init(&_TrustedIssuersRegistry.TransactOpts)
}

// RemoveTrustedIssuer is a paid mutator transaction binding the contract method 0xb93d28eb.
//
// Solidity: function removeTrustedIssuer(address _trustedIssuer) returns()
func (_TrustedIssuersRegistry *TrustedIssuersRegistryTransactor) RemoveTrustedIssuer(opts *bind.TransactOpts, _trustedIssuer common.Address) (*types.Transaction, error) {
	return _TrustedIssuersRegistry.contract.Transact(opts, "removeTrustedIssuer", _trustedIssuer)
}

// RemoveTrustedIssuer is a paid mutator transaction binding the contract method 0xb93d28eb.
//
// Solidity: function removeTrustedIssuer(address _trustedIssuer) returns()
func (_TrustedIssuersRegistry *TrustedIssuersRegistrySession) RemoveTrustedIssuer(_trustedIssuer common.Address) (*types.Transaction, error) {
	return _TrustedIssuersRegistry.Contract.RemoveTrustedIssuer(&_TrustedIssuersRegistry.TransactOpts, _trustedIssuer)
}

// RemoveTrustedIssuer is a paid mutator transaction binding the contract method 0xb93d28eb.
//
// Solidity: function removeTrustedIssuer(address _trustedIssuer) returns()
func (_TrustedIssuersRegistry *TrustedIssuersRegistryTransactorSession) RemoveTrustedIssuer(_trustedIssuer common.Address) (*types.Transaction, error) {
	return _TrustedIssuersRegistry.Contract.RemoveTrustedIssuer(&_TrustedIssuersRegistry.TransactOpts, _trustedIssuer)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TrustedIssuersRegistry *TrustedIssuersRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TrustedIssuersRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TrustedIssuersRegistry *TrustedIssuersRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _TrustedIssuersRegistry.Contract.RenounceOwnership(&_TrustedIssuersRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TrustedIssuersRegistry *TrustedIssuersRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TrustedIssuersRegistry.Contract.RenounceOwnership(&_TrustedIssuersRegistry.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TrustedIssuersRegistry *TrustedIssuersRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TrustedIssuersRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TrustedIssuersRegistry *TrustedIssuersRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TrustedIssuersRegistry.Contract.TransferOwnership(&_TrustedIssuersRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TrustedIssuersRegistry *TrustedIssuersRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TrustedIssuersRegistry.Contract.TransferOwnership(&_TrustedIssuersRegistry.TransactOpts, newOwner)
}

// UpdateIssuerClaimTopics is a paid mutator transaction binding the contract method 0x04bc7e84.
//
// Solidity: function updateIssuerClaimTopics(address _trustedIssuer, uint256[] _claimTopics) returns()
func (_TrustedIssuersRegistry *TrustedIssuersRegistryTransactor) UpdateIssuerClaimTopics(opts *bind.TransactOpts, _trustedIssuer common.Address, _claimTopics []*big.Int) (*types.Transaction, error) {
	return _TrustedIssuersRegistry.contract.Transact(opts, "updateIssuerClaimTopics", _trustedIssuer, _claimTopics)
}

// UpdateIssuerClaimTopics is a paid mutator transaction binding the contract method 0x04bc7e84.
//
// Solidity: function updateIssuerClaimTopics(address _trustedIssuer, uint256[] _claimTopics) returns()
func (_TrustedIssuersRegistry *TrustedIssuersRegistrySession) UpdateIssuerClaimTopics(_trustedIssuer common.Address, _claimTopics []*big.Int) (*types.Transaction, error) {
	return _TrustedIssuersRegistry.Contract.UpdateIssuerClaimTopics(&_TrustedIssuersRegistry.TransactOpts, _trustedIssuer, _claimTopics)
}

// UpdateIssuerClaimTopics is a paid mutator transaction binding the contract method 0x04bc7e84.
//
// Solidity: function updateIssuerClaimTopics(address _trustedIssuer, uint256[] _claimTopics) returns()
func (_TrustedIssuersRegistry *TrustedIssuersRegistryTransactorSession) UpdateIssuerClaimTopics(_trustedIssuer common.Address, _claimTopics []*big.Int) (*types.Transaction, error) {
	return _TrustedIssuersRegistry.Contract.UpdateIssuerClaimTopics(&_TrustedIssuersRegistry.TransactOpts, _trustedIssuer, _claimTopics)
}

// TrustedIssuersRegistryClaimTopicsUpdatedIterator is returned from FilterClaimTopicsUpdated and is used to iterate over the raw logs and unpacked data for ClaimTopicsUpdated events raised by the TrustedIssuersRegistry contract.
type TrustedIssuersRegistryClaimTopicsUpdatedIterator struct {
	Event *TrustedIssuersRegistryClaimTopicsUpdated // Event containing the contract specifics and raw log

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
func (it *TrustedIssuersRegistryClaimTopicsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrustedIssuersRegistryClaimTopicsUpdated)
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
		it.Event = new(TrustedIssuersRegistryClaimTopicsUpdated)
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
func (it *TrustedIssuersRegistryClaimTopicsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrustedIssuersRegistryClaimTopicsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrustedIssuersRegistryClaimTopicsUpdated represents a ClaimTopicsUpdated event raised by the TrustedIssuersRegistry contract.
type TrustedIssuersRegistryClaimTopicsUpdated struct {
	TrustedIssuer common.Address
	ClaimTopics   []*big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterClaimTopicsUpdated is a free log retrieval operation binding the contract event 0xec753cfc52044f61676f18a11e500093a9f2b1cd5e4942bc476f2b0438159bcf.
//
// Solidity: event ClaimTopicsUpdated(address indexed trustedIssuer, uint256[] claimTopics)
func (_TrustedIssuersRegistry *TrustedIssuersRegistryFilterer) FilterClaimTopicsUpdated(opts *bind.FilterOpts, trustedIssuer []common.Address) (*TrustedIssuersRegistryClaimTopicsUpdatedIterator, error) {

	var trustedIssuerRule []interface{}
	for _, trustedIssuerItem := range trustedIssuer {
		trustedIssuerRule = append(trustedIssuerRule, trustedIssuerItem)
	}

	logs, sub, err := _TrustedIssuersRegistry.contract.FilterLogs(opts, "ClaimTopicsUpdated", trustedIssuerRule)
	if err != nil {
		return nil, err
	}
	return &TrustedIssuersRegistryClaimTopicsUpdatedIterator{contract: _TrustedIssuersRegistry.contract, event: "ClaimTopicsUpdated", logs: logs, sub: sub}, nil
}

// WatchClaimTopicsUpdated is a free log subscription operation binding the contract event 0xec753cfc52044f61676f18a11e500093a9f2b1cd5e4942bc476f2b0438159bcf.
//
// Solidity: event ClaimTopicsUpdated(address indexed trustedIssuer, uint256[] claimTopics)
func (_TrustedIssuersRegistry *TrustedIssuersRegistryFilterer) WatchClaimTopicsUpdated(opts *bind.WatchOpts, sink chan<- *TrustedIssuersRegistryClaimTopicsUpdated, trustedIssuer []common.Address) (event.Subscription, error) {

	var trustedIssuerRule []interface{}
	for _, trustedIssuerItem := range trustedIssuer {
		trustedIssuerRule = append(trustedIssuerRule, trustedIssuerItem)
	}

	logs, sub, err := _TrustedIssuersRegistry.contract.WatchLogs(opts, "ClaimTopicsUpdated", trustedIssuerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrustedIssuersRegistryClaimTopicsUpdated)
				if err := _TrustedIssuersRegistry.contract.UnpackLog(event, "ClaimTopicsUpdated", log); err != nil {
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

// ParseClaimTopicsUpdated is a log parse operation binding the contract event 0xec753cfc52044f61676f18a11e500093a9f2b1cd5e4942bc476f2b0438159bcf.
//
// Solidity: event ClaimTopicsUpdated(address indexed trustedIssuer, uint256[] claimTopics)
func (_TrustedIssuersRegistry *TrustedIssuersRegistryFilterer) ParseClaimTopicsUpdated(log types.Log) (*TrustedIssuersRegistryClaimTopicsUpdated, error) {
	event := new(TrustedIssuersRegistryClaimTopicsUpdated)
	if err := _TrustedIssuersRegistry.contract.UnpackLog(event, "ClaimTopicsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrustedIssuersRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the TrustedIssuersRegistry contract.
type TrustedIssuersRegistryInitializedIterator struct {
	Event *TrustedIssuersRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *TrustedIssuersRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrustedIssuersRegistryInitialized)
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
		it.Event = new(TrustedIssuersRegistryInitialized)
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
func (it *TrustedIssuersRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrustedIssuersRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrustedIssuersRegistryInitialized represents a Initialized event raised by the TrustedIssuersRegistry contract.
type TrustedIssuersRegistryInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_TrustedIssuersRegistry *TrustedIssuersRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*TrustedIssuersRegistryInitializedIterator, error) {

	logs, sub, err := _TrustedIssuersRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &TrustedIssuersRegistryInitializedIterator{contract: _TrustedIssuersRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_TrustedIssuersRegistry *TrustedIssuersRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *TrustedIssuersRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _TrustedIssuersRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrustedIssuersRegistryInitialized)
				if err := _TrustedIssuersRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_TrustedIssuersRegistry *TrustedIssuersRegistryFilterer) ParseInitialized(log types.Log) (*TrustedIssuersRegistryInitialized, error) {
	event := new(TrustedIssuersRegistryInitialized)
	if err := _TrustedIssuersRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrustedIssuersRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TrustedIssuersRegistry contract.
type TrustedIssuersRegistryOwnershipTransferredIterator struct {
	Event *TrustedIssuersRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TrustedIssuersRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrustedIssuersRegistryOwnershipTransferred)
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
		it.Event = new(TrustedIssuersRegistryOwnershipTransferred)
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
func (it *TrustedIssuersRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrustedIssuersRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrustedIssuersRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the TrustedIssuersRegistry contract.
type TrustedIssuersRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TrustedIssuersRegistry *TrustedIssuersRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TrustedIssuersRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TrustedIssuersRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TrustedIssuersRegistryOwnershipTransferredIterator{contract: _TrustedIssuersRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TrustedIssuersRegistry *TrustedIssuersRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TrustedIssuersRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TrustedIssuersRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrustedIssuersRegistryOwnershipTransferred)
				if err := _TrustedIssuersRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TrustedIssuersRegistry *TrustedIssuersRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*TrustedIssuersRegistryOwnershipTransferred, error) {
	event := new(TrustedIssuersRegistryOwnershipTransferred)
	if err := _TrustedIssuersRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrustedIssuersRegistryTrustedIssuerAddedIterator is returned from FilterTrustedIssuerAdded and is used to iterate over the raw logs and unpacked data for TrustedIssuerAdded events raised by the TrustedIssuersRegistry contract.
type TrustedIssuersRegistryTrustedIssuerAddedIterator struct {
	Event *TrustedIssuersRegistryTrustedIssuerAdded // Event containing the contract specifics and raw log

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
func (it *TrustedIssuersRegistryTrustedIssuerAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrustedIssuersRegistryTrustedIssuerAdded)
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
		it.Event = new(TrustedIssuersRegistryTrustedIssuerAdded)
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
func (it *TrustedIssuersRegistryTrustedIssuerAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrustedIssuersRegistryTrustedIssuerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrustedIssuersRegistryTrustedIssuerAdded represents a TrustedIssuerAdded event raised by the TrustedIssuersRegistry contract.
type TrustedIssuersRegistryTrustedIssuerAdded struct {
	TrustedIssuer common.Address
	ClaimTopics   []*big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterTrustedIssuerAdded is a free log retrieval operation binding the contract event 0xfedc33fd34859594822c0ff6f3f4f9fc279cc6d5cae53068f706a088e4500872.
//
// Solidity: event TrustedIssuerAdded(address indexed trustedIssuer, uint256[] claimTopics)
func (_TrustedIssuersRegistry *TrustedIssuersRegistryFilterer) FilterTrustedIssuerAdded(opts *bind.FilterOpts, trustedIssuer []common.Address) (*TrustedIssuersRegistryTrustedIssuerAddedIterator, error) {

	var trustedIssuerRule []interface{}
	for _, trustedIssuerItem := range trustedIssuer {
		trustedIssuerRule = append(trustedIssuerRule, trustedIssuerItem)
	}

	logs, sub, err := _TrustedIssuersRegistry.contract.FilterLogs(opts, "TrustedIssuerAdded", trustedIssuerRule)
	if err != nil {
		return nil, err
	}
	return &TrustedIssuersRegistryTrustedIssuerAddedIterator{contract: _TrustedIssuersRegistry.contract, event: "TrustedIssuerAdded", logs: logs, sub: sub}, nil
}

// WatchTrustedIssuerAdded is a free log subscription operation binding the contract event 0xfedc33fd34859594822c0ff6f3f4f9fc279cc6d5cae53068f706a088e4500872.
//
// Solidity: event TrustedIssuerAdded(address indexed trustedIssuer, uint256[] claimTopics)
func (_TrustedIssuersRegistry *TrustedIssuersRegistryFilterer) WatchTrustedIssuerAdded(opts *bind.WatchOpts, sink chan<- *TrustedIssuersRegistryTrustedIssuerAdded, trustedIssuer []common.Address) (event.Subscription, error) {

	var trustedIssuerRule []interface{}
	for _, trustedIssuerItem := range trustedIssuer {
		trustedIssuerRule = append(trustedIssuerRule, trustedIssuerItem)
	}

	logs, sub, err := _TrustedIssuersRegistry.contract.WatchLogs(opts, "TrustedIssuerAdded", trustedIssuerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrustedIssuersRegistryTrustedIssuerAdded)
				if err := _TrustedIssuersRegistry.contract.UnpackLog(event, "TrustedIssuerAdded", log); err != nil {
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

// ParseTrustedIssuerAdded is a log parse operation binding the contract event 0xfedc33fd34859594822c0ff6f3f4f9fc279cc6d5cae53068f706a088e4500872.
//
// Solidity: event TrustedIssuerAdded(address indexed trustedIssuer, uint256[] claimTopics)
func (_TrustedIssuersRegistry *TrustedIssuersRegistryFilterer) ParseTrustedIssuerAdded(log types.Log) (*TrustedIssuersRegistryTrustedIssuerAdded, error) {
	event := new(TrustedIssuersRegistryTrustedIssuerAdded)
	if err := _TrustedIssuersRegistry.contract.UnpackLog(event, "TrustedIssuerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrustedIssuersRegistryTrustedIssuerRemovedIterator is returned from FilterTrustedIssuerRemoved and is used to iterate over the raw logs and unpacked data for TrustedIssuerRemoved events raised by the TrustedIssuersRegistry contract.
type TrustedIssuersRegistryTrustedIssuerRemovedIterator struct {
	Event *TrustedIssuersRegistryTrustedIssuerRemoved // Event containing the contract specifics and raw log

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
func (it *TrustedIssuersRegistryTrustedIssuerRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrustedIssuersRegistryTrustedIssuerRemoved)
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
		it.Event = new(TrustedIssuersRegistryTrustedIssuerRemoved)
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
func (it *TrustedIssuersRegistryTrustedIssuerRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrustedIssuersRegistryTrustedIssuerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrustedIssuersRegistryTrustedIssuerRemoved represents a TrustedIssuerRemoved event raised by the TrustedIssuersRegistry contract.
type TrustedIssuersRegistryTrustedIssuerRemoved struct {
	TrustedIssuer common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterTrustedIssuerRemoved is a free log retrieval operation binding the contract event 0x2214ded40113cc3fb63fc206cafee88270b0a903dac7245d54efdde30ebb0321.
//
// Solidity: event TrustedIssuerRemoved(address indexed trustedIssuer)
func (_TrustedIssuersRegistry *TrustedIssuersRegistryFilterer) FilterTrustedIssuerRemoved(opts *bind.FilterOpts, trustedIssuer []common.Address) (*TrustedIssuersRegistryTrustedIssuerRemovedIterator, error) {

	var trustedIssuerRule []interface{}
	for _, trustedIssuerItem := range trustedIssuer {
		trustedIssuerRule = append(trustedIssuerRule, trustedIssuerItem)
	}

	logs, sub, err := _TrustedIssuersRegistry.contract.FilterLogs(opts, "TrustedIssuerRemoved", trustedIssuerRule)
	if err != nil {
		return nil, err
	}
	return &TrustedIssuersRegistryTrustedIssuerRemovedIterator{contract: _TrustedIssuersRegistry.contract, event: "TrustedIssuerRemoved", logs: logs, sub: sub}, nil
}

// WatchTrustedIssuerRemoved is a free log subscription operation binding the contract event 0x2214ded40113cc3fb63fc206cafee88270b0a903dac7245d54efdde30ebb0321.
//
// Solidity: event TrustedIssuerRemoved(address indexed trustedIssuer)
func (_TrustedIssuersRegistry *TrustedIssuersRegistryFilterer) WatchTrustedIssuerRemoved(opts *bind.WatchOpts, sink chan<- *TrustedIssuersRegistryTrustedIssuerRemoved, trustedIssuer []common.Address) (event.Subscription, error) {

	var trustedIssuerRule []interface{}
	for _, trustedIssuerItem := range trustedIssuer {
		trustedIssuerRule = append(trustedIssuerRule, trustedIssuerItem)
	}

	logs, sub, err := _TrustedIssuersRegistry.contract.WatchLogs(opts, "TrustedIssuerRemoved", trustedIssuerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrustedIssuersRegistryTrustedIssuerRemoved)
				if err := _TrustedIssuersRegistry.contract.UnpackLog(event, "TrustedIssuerRemoved", log); err != nil {
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

// ParseTrustedIssuerRemoved is a log parse operation binding the contract event 0x2214ded40113cc3fb63fc206cafee88270b0a903dac7245d54efdde30ebb0321.
//
// Solidity: event TrustedIssuerRemoved(address indexed trustedIssuer)
func (_TrustedIssuersRegistry *TrustedIssuersRegistryFilterer) ParseTrustedIssuerRemoved(log types.Log) (*TrustedIssuersRegistryTrustedIssuerRemoved, error) {
	event := new(TrustedIssuersRegistryTrustedIssuerRemoved)
	if err := _TrustedIssuersRegistry.contract.UnpackLog(event, "TrustedIssuerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
