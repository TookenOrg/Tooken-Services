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

// IdentityRegistryMetaData contains all meta data concerning the IdentityRegistry contract.
var IdentityRegistryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"}],\"name\":\"AgentAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"}],\"name\":\"AgentRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"claimTopicsRegistry\",\"type\":\"address\"}],\"name\":\"ClaimTopicsRegistrySet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"investorAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint16\",\"name\":\"country\",\"type\":\"uint16\"}],\"name\":\"CountryUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"investorAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIIdentity\",\"name\":\"identity\",\"type\":\"address\"}],\"name\":\"IdentityRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"investorAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIIdentity\",\"name\":\"identity\",\"type\":\"address\"}],\"name\":\"IdentityRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"identityStorage\",\"type\":\"address\"}],\"name\":\"IdentityStorageSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIIdentity\",\"name\":\"oldIdentity\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIIdentity\",\"name\":\"newIdentity\",\"type\":\"address\"}],\"name\":\"IdentityUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"trustedIssuersRegistry\",\"type\":\"address\"}],\"name\":\"TrustedIssuersRegistrySet\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"}],\"name\":\"addAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_userAddresses\",\"type\":\"address[]\"},{\"internalType\":\"contractIIdentity[]\",\"name\":\"_identities\",\"type\":\"address[]\"},{\"internalType\":\"uint16[]\",\"name\":\"_countries\",\"type\":\"uint16[]\"}],\"name\":\"batchRegisterIdentity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userAddress\",\"type\":\"address\"}],\"name\":\"contains\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userAddress\",\"type\":\"address\"}],\"name\":\"deleteIdentity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userAddress\",\"type\":\"address\"}],\"name\":\"identity\",\"outputs\":[{\"internalType\":\"contractIIdentity\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"identityStorage\",\"outputs\":[{\"internalType\":\"contractIIdentityRegistryStorage\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_trustedIssuersRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_claimTopicsRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_identityStorage\",\"type\":\"address\"}],\"name\":\"init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userAddress\",\"type\":\"address\"}],\"name\":\"investorCountry\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"}],\"name\":\"isAgent\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userAddress\",\"type\":\"address\"}],\"name\":\"isVerified\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"issuersRegistry\",\"outputs\":[{\"internalType\":\"contractITrustedIssuersRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userAddress\",\"type\":\"address\"},{\"internalType\":\"contractIIdentity\",\"name\":\"_identity\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"_country\",\"type\":\"uint16\"}],\"name\":\"registerIdentity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"}],\"name\":\"removeAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_claimTopicsRegistry\",\"type\":\"address\"}],\"name\":\"setClaimTopicsRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_identityRegistryStorage\",\"type\":\"address\"}],\"name\":\"setIdentityRegistryStorage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_trustedIssuersRegistry\",\"type\":\"address\"}],\"name\":\"setTrustedIssuersRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"topicsRegistry\",\"outputs\":[{\"internalType\":\"contractIClaimTopicsRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userAddress\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"_country\",\"type\":\"uint16\"}],\"name\":\"updateCountry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userAddress\",\"type\":\"address\"},{\"internalType\":\"contractIIdentity\",\"name\":\"_identity\",\"type\":\"address\"}],\"name\":\"updateIdentity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611daf806100206000396000f3fe608060405234801561001057600080fd5b50600436106101825760003560e01c806384e79842116100d8578063b4f3fcb71161008c578063f0eb5e5411610066578063f0eb5e541461031d578063f11abfd814610330578063f2fde38b1461034157600080fd5b8063b4f3fcb7146102e6578063b9209e33146102f7578063e744d7891461030a57600080fd5b80638e098ca1116100bd5780638e098ca1146102ad57806397a6278e146102c0578063a8d29d1d146102d357600080fd5b806384e79842146102895780638da5cb5b1461029c57600080fd5b8063454a03e01161013a578063670af6a911610114578063670af6a914610248578063715018a61461025b5780637e42683b1461026357600080fd5b8063454a03e01461020f5780635dbe47e814610222578063653dc9f11461023557600080fd5b806326d941ae1161016b57806326d941ae146101c45780633b239a7f146101d75780633b3e12f4146101ea57600080fd5b8063184b9559146101875780631ffbb0641461019c575b600080fd5b61019a610195366004611751565b610354565b005b6101af6101aa36600461179c565b6105c8565b60405190151581526020015b60405180910390f35b61019a6101d236600461179c565b6105db565b61019a6101e53660046117d0565b61062d565b6066546001600160a01b03165b6040516001600160a01b0390911681526020016101bb565b61019a61021d366004611809565b610759565b6101af61023036600461179c565b61088c565b61019a610243366004611895565b6108b6565b61019a61025636600461179c565b610954565b61019a6109a6565b61027661027136600461179c565b6109ba565b60405161ffff90911681526020016101bb565b61019a61029736600461179c565b610a42565b6033546001600160a01b03166101f7565b61019a6102bb36600461192f565b610ae2565b61019a6102ce36600461179c565b610c18565b61019a6102e136600461179c565b610cb8565b6067546001600160a01b03166101f7565b6101af61030536600461179c565b610de5565b61019a61031836600461179c565b611287565b6101f761032b36600461179c565b6112d9565b6068546001600160a01b03166101f7565b61019a61034f36600461179c565b611361565b600054610100900460ff16158080156103745750600054600160ff909116105b8061038e5750303b15801561038e575060005460ff166001145b6104055760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084015b60405180910390fd5b6000805460ff191660011790558015610428576000805461ff0019166101001790555b6001600160a01b0384161580159061044857506001600160a01b03831615155b801561045c57506001600160a01b03821615155b6104a85760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064016103fc565b606680546001600160a01b038086166001600160a01b031992831681179093556067805488831690841617905560688054918616919092161790556040517f7170bf15b246e880b2369cd7c67d057760d8a35149e8c64dde91efa22bcc76d090600090a26040516001600160a01b038516907f1b98cb79e6f73020175fe87333f1b91ad6a881519c0afe30340c2599b2b4bde090600090a26040516001600160a01b038316907f2fa8b95c1db7afe99e3398f3792f008135cedc1fa26b0bb2ecd2352cd166d53c90600090a261057c6113f1565b80156105c2576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50505050565b60006105d5606583611464565b92915050565b6105e3611502565b606880546001600160a01b0319166001600160a01b0383169081179091556040517f2fa8b95c1db7afe99e3398f3792f008135cedc1fa26b0bb2ecd2352cd166d53c90600090a250565b610636336105c8565b6106995760405162461bcd60e51b815260206004820152602e60248201527f4167656e74526f6c653a2063616c6c657220646f6573206e6f7420686176652060448201526d746865204167656e7420726f6c6560901b60648201526084016103fc565b6068546040517f9f3418d50000000000000000000000000000000000000000000000000000000081526001600160a01b03848116600483015261ffff8416602483015290911690639f3418d590604401600060405180830381600087803b15801561070357600080fd5b505af1158015610717573d6000803e3d6000fd5b505060405161ffff841692506001600160a01b03851691507f04ed3b726495c2dca1ff1215d9ca54e1a4030abb5e82b0f6ce55702416cee85390600090a35050565b610762336105c8565b6107c55760405162461bcd60e51b815260206004820152602e60248201527f4167656e74526f6c653a2063616c6c657220646f6573206e6f7420686176652060448201526d746865204167656e7420726f6c6560901b60648201526084016103fc565b6068546040517fa53410dd0000000000000000000000000000000000000000000000000000000081526001600160a01b038581166004830152848116602483015261ffff841660448301529091169063a53410dd90606401600060405180830381600087803b15801561083757600080fd5b505af115801561084b573d6000803e3d6000fd5b50506040516001600160a01b038086169350861691507f6ae73635c50d24a45af6fbd5e016ac4bed179addbc8bf24e04ff0fcc6d33af1990600090a3505050565b600080610898836112d9565b6001600160a01b0316036108ae57506000919050565b506001919050565b60005b8581101561094b576109398787838181106108d6576108d661195d565b90506020020160208101906108eb919061179c565b8686848181106108fd576108fd61195d565b9050602002016020810190610912919061179c565b8585858181106109245761092461195d565b905060200201602081019061021d9190611973565b80610943816119a6565b9150506108b9565b50505050505050565b61095c611502565b606680546001600160a01b0319166001600160a01b0383169081179091556040517f7170bf15b246e880b2369cd7c67d057760d8a35149e8c64dde91efa22bcc76d090600090a250565b6109ae611502565b6109b8600061155c565b565b6068546040517f727e13bc0000000000000000000000000000000000000000000000000000000081526001600160a01b038381166004830152600092169063727e13bc90602401602060405180830381865afa158015610a1e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105d591906119bf565b610a4a611502565b6001600160a01b038116610aa05760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064016103fc565b610aab6065826115ae565b6040516001600160a01b038216907ff68e73cec97f2d70aa641fb26e87a4383686e2efacb648f2165aeb02ac562ec590600090a250565b610aeb336105c8565b610b4e5760405162461bcd60e51b815260206004820152602e60248201527f4167656e74526f6c653a2063616c6c657220646f6573206e6f7420686176652060448201526d746865204167656e7420726f6c6560901b60648201526084016103fc565b6000610b59836112d9565b6068546040517fe805cf860000000000000000000000000000000000000000000000000000000081526001600160a01b038681166004830152858116602483015292935091169063e805cf8690604401600060405180830381600087803b158015610bc357600080fd5b505af1158015610bd7573d6000803e3d6000fd5b50506040516001600160a01b038086169350841691507fe98082932c8056a0f514da9104e4a66bc2cbaef102ad59d90c4b24220ebf601090600090a3505050565b610c20611502565b6001600160a01b038116610c765760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064016103fc565b610c8160658261162a565b6040516001600160a01b038216907fed9c8ad8d5a0a66898ea49d2956929c93ae2e8bd50281b2ed897c5d1a6737e0b90600090a250565b610cc1336105c8565b610d245760405162461bcd60e51b815260206004820152602e60248201527f4167656e74526f6c653a2063616c6c657220646f6573206e6f7420686176652060448201526d746865204167656e7420726f6c6560901b60648201526084016103fc565b6000610d2f826112d9565b6068546040517fcf191bcd0000000000000000000000000000000000000000000000000000000081526001600160a01b03858116600483015292935091169063cf191bcd90602401600060405180830381600087803b158015610d9157600080fd5b505af1158015610da5573d6000803e3d6000fd5b50506040516001600160a01b038085169350851691507f59d6590e225b81befe259af056324092801080acbb7feab310eb34678871f32790600090a35050565b600080610df1836112d9565b6001600160a01b031603610e0757506000919050565b606654604080517fdf09d60400000000000000000000000000000000000000000000000000000000815290516000926001600160a01b03169163df09d60491600480830192869291908290030181865afa158015610e69573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610e919190810190611a47565b90508051600003610ea55750600192915050565b600080600060608060005b86518110156112785760675487516000916001600160a01b0316906352c111d1908a9085908110610ee357610ee361195d565b60200260200101516040518263ffffffff1660e01b8152600401610f0991815260200190565b600060405180830381865afa158015610f26573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610f4e9190810190611add565b90508051600003610f69575060009998505050505050505050565b6000815167ffffffffffffffff811115610f8557610f856119dc565b604051908082528060200260200182016040528015610fae578160200160208202803683370190505b50905060005b825181101561105b57828181518110610fcf57610fcf61195d565b60200260200101518a8581518110610fe957610fe961195d565b60200260200101516040516020016110169291906001600160a01b03929092168252602082015260400190565b6040516020818303038152906040528051906020012082828151811061103e5761103e61195d565b602090810291909101015280611053816119a6565b915050610fb4565b5060005b8151811015611262576110718c6112d9565b6001600160a01b031663c9100bcb8383815181106110915761109161195d565b60200260200101516040518263ffffffff1660e01b81526004016110b791815260200190565b600060405180830381865afa1580156110d4573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526110fc9190810190611bff565b508e51949d50929b50909950975095508a908590811061111e5761111e61195d565b6020026020010151890361122a57866001600160a01b031663c0969a6e6111448e6112d9565b8c87815181106111565761115661195d565b602002602001015189896040518563ffffffff1660e01b815260040161117f9493929190611ced565b602060405180830381865afa9250505080156111b8575060408051601f3d908101601f191682019092526111b591810190611d27565b60015b6111e757600182516111ca9190611d49565b81036111e2575060009b9a5050505050505050505050565b611250565b80156111f257825191505b8015801561120c5750600183516112099190611d49565b82145b15611224575060009c9b505050505050505050505050565b50611250565b600182516112389190611d49565b8103611250575060009b9a5050505050505050505050565b8061125a816119a6565b91505061105f565b5050508080611270906119a6565b915050610eb0565b50600198975050505050505050565b61128f611502565b606780546001600160a01b0319166001600160a01b0383169081179091556040517f1b98cb79e6f73020175fe87333f1b91ad6a881519c0afe30340c2599b2b4bde090600090a250565b6068546040517f7988d3a50000000000000000000000000000000000000000000000000000000081526001600160a01b0383811660048301526000921690637988d3a590602401602060405180830381865afa15801561133d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105d59190611d5c565b611369611502565b6001600160a01b0381166113e55760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f646472657373000000000000000000000000000000000000000000000000000060648201526084016103fc565b6113ee8161155c565b50565b600054610100900460ff1661145c5760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b60648201526084016103fc565b6109b86116c8565b60006001600160a01b0382166114e25760405162461bcd60e51b815260206004820152602260248201527f526f6c65733a206163636f756e7420697320746865207a65726f20616464726560448201527f737300000000000000000000000000000000000000000000000000000000000060648201526084016103fc565b506001600160a01b03166000908152602091909152604090205460ff1690565b6033546001600160a01b031633146109b85760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016103fc565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6115b88282611464565b156116055760405162461bcd60e51b815260206004820152601f60248201527f526f6c65733a206163636f756e7420616c72656164792068617320726f6c650060448201526064016103fc565b6001600160a01b0316600090815260209190915260409020805460ff19166001179055565b6116348282611464565b6116a65760405162461bcd60e51b815260206004820152602160248201527f526f6c65733a206163636f756e7420646f6573206e6f74206861766520726f6c60448201527f650000000000000000000000000000000000000000000000000000000000000060648201526084016103fc565b6001600160a01b0316600090815260209190915260409020805460ff19169055565b600054610100900460ff166117335760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b60648201526084016103fc565b6109b83361155c565b6001600160a01b03811681146113ee57600080fd5b60008060006060848603121561176657600080fd5b83356117718161173c565b925060208401356117818161173c565b915060408401356117918161173c565b809150509250925092565b6000602082840312156117ae57600080fd5b81356117b98161173c565b9392505050565b61ffff811681146113ee57600080fd5b600080604083850312156117e357600080fd5b82356117ee8161173c565b915060208301356117fe816117c0565b809150509250929050565b60008060006060848603121561181e57600080fd5b83356118298161173c565b925060208401356118398161173c565b91506040840135611791816117c0565b60008083601f84011261185b57600080fd5b50813567ffffffffffffffff81111561187357600080fd5b6020830191508360208260051b850101111561188e57600080fd5b9250929050565b600080600080600080606087890312156118ae57600080fd5b863567ffffffffffffffff808211156118c657600080fd5b6118d28a838b01611849565b909850965060208901359150808211156118eb57600080fd5b6118f78a838b01611849565b9096509450604089013591508082111561191057600080fd5b5061191d89828a01611849565b979a9699509497509295939492505050565b6000806040838503121561194257600080fd5b823561194d8161173c565b915060208301356117fe8161173c565b634e487b7160e01b600052603260045260246000fd5b60006020828403121561198557600080fd5b81356117b9816117c0565b634e487b7160e01b600052601160045260246000fd5b6000600182016119b8576119b8611990565b5060010190565b6000602082840312156119d157600080fd5b81516117b9816117c0565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715611a1b57611a1b6119dc565b604052919050565b600067ffffffffffffffff821115611a3d57611a3d6119dc565b5060051b60200190565b60006020808385031215611a5a57600080fd5b825167ffffffffffffffff811115611a7157600080fd5b8301601f81018513611a8257600080fd5b8051611a95611a9082611a23565b6119f2565b81815260059190911b82018301908381019087831115611ab457600080fd5b928401925b82841015611ad257835182529284019290840190611ab9565b979650505050505050565b60006020808385031215611af057600080fd5b825167ffffffffffffffff811115611b0757600080fd5b8301601f81018513611b1857600080fd5b8051611b26611a9082611a23565b81815260059190911b82018301908381019087831115611b4557600080fd5b928401925b82841015611ad2578351611b5d8161173c565b82529284019290840190611b4a565b60005b83811015611b87578181015183820152602001611b6f565b50506000910152565b600067ffffffffffffffff831115611baa57611baa6119dc565b611bbd601f8401601f19166020016119f2565b9050828152838383011115611bd157600080fd5b6117b9836020830184611b6c565b600082601f830112611bf057600080fd5b6117b983835160208501611b90565b60008060008060008060c08789031215611c1857600080fd5b86519550602087015194506040870151611c318161173c565b606088015190945067ffffffffffffffff80821115611c4f57600080fd5b611c5b8a838b01611bdf565b94506080890151915080821115611c7157600080fd5b611c7d8a838b01611bdf565b935060a0890151915080821115611c9357600080fd5b508701601f81018913611ca557600080fd5b611cb489825160208401611b90565b9150509295509295509295565b60008151808452611cd9816020860160208601611b6c565b601f01601f19169290920160200192915050565b6001600160a01b0385168152836020820152608060408201526000611d156080830185611cc1565b8281036060840152611ad28185611cc1565b600060208284031215611d3957600080fd5b815180151581146117b957600080fd5b818103818111156105d5576105d5611990565b600060208284031215611d6e57600080fd5b81516117b98161173c56fea264697066735822122007f535ef37df51033c230c3cce306d1011026544569db92212b9d18bb8c7fe3864736f6c63430008110033",
}

// IdentityRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use IdentityRegistryMetaData.ABI instead.
var IdentityRegistryABI = IdentityRegistryMetaData.ABI

// IdentityRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use IdentityRegistryMetaData.Bin instead.
var IdentityRegistryBin = IdentityRegistryMetaData.Bin

// DeployIdentityRegistry deploys a new Ethereum contract, binding an instance of IdentityRegistry to it.
func DeployIdentityRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *IdentityRegistry, error) {
	parsed, err := IdentityRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(IdentityRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &IdentityRegistry{IdentityRegistryCaller: IdentityRegistryCaller{contract: contract}, IdentityRegistryTransactor: IdentityRegistryTransactor{contract: contract}, IdentityRegistryFilterer: IdentityRegistryFilterer{contract: contract}}, nil
}

// IdentityRegistry is an auto generated Go binding around an Ethereum contract.
type IdentityRegistry struct {
	IdentityRegistryCaller     // Read-only binding to the contract
	IdentityRegistryTransactor // Write-only binding to the contract
	IdentityRegistryFilterer   // Log filterer for contract events
}

// IdentityRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type IdentityRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IdentityRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IdentityRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IdentityRegistrySession struct {
	Contract     *IdentityRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IdentityRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IdentityRegistryCallerSession struct {
	Contract *IdentityRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// IdentityRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IdentityRegistryTransactorSession struct {
	Contract     *IdentityRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// IdentityRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type IdentityRegistryRaw struct {
	Contract *IdentityRegistry // Generic contract binding to access the raw methods on
}

// IdentityRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IdentityRegistryCallerRaw struct {
	Contract *IdentityRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// IdentityRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IdentityRegistryTransactorRaw struct {
	Contract *IdentityRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIdentityRegistry creates a new instance of IdentityRegistry, bound to a specific deployed contract.
func NewIdentityRegistry(address common.Address, backend bind.ContractBackend) (*IdentityRegistry, error) {
	contract, err := bindIdentityRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistry{IdentityRegistryCaller: IdentityRegistryCaller{contract: contract}, IdentityRegistryTransactor: IdentityRegistryTransactor{contract: contract}, IdentityRegistryFilterer: IdentityRegistryFilterer{contract: contract}}, nil
}

// NewIdentityRegistryCaller creates a new read-only instance of IdentityRegistry, bound to a specific deployed contract.
func NewIdentityRegistryCaller(address common.Address, caller bind.ContractCaller) (*IdentityRegistryCaller, error) {
	contract, err := bindIdentityRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryCaller{contract: contract}, nil
}

// NewIdentityRegistryTransactor creates a new write-only instance of IdentityRegistry, bound to a specific deployed contract.
func NewIdentityRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*IdentityRegistryTransactor, error) {
	contract, err := bindIdentityRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryTransactor{contract: contract}, nil
}

// NewIdentityRegistryFilterer creates a new log filterer instance of IdentityRegistry, bound to a specific deployed contract.
func NewIdentityRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*IdentityRegistryFilterer, error) {
	contract, err := bindIdentityRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryFilterer{contract: contract}, nil
}

// bindIdentityRegistry binds a generic wrapper to an already deployed contract.
func bindIdentityRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IdentityRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdentityRegistry *IdentityRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IdentityRegistry.Contract.IdentityRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdentityRegistry *IdentityRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.IdentityRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdentityRegistry *IdentityRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.IdentityRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdentityRegistry *IdentityRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IdentityRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdentityRegistry *IdentityRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdentityRegistry *IdentityRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.contract.Transact(opts, method, params...)
}

// Contains is a free data retrieval call binding the contract method 0x5dbe47e8.
//
// Solidity: function contains(address _userAddress) view returns(bool)
func (_IdentityRegistry *IdentityRegistryCaller) Contains(opts *bind.CallOpts, _userAddress common.Address) (bool, error) {
	var out []interface{}
	err := _IdentityRegistry.contract.Call(opts, &out, "contains", _userAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Contains is a free data retrieval call binding the contract method 0x5dbe47e8.
//
// Solidity: function contains(address _userAddress) view returns(bool)
func (_IdentityRegistry *IdentityRegistrySession) Contains(_userAddress common.Address) (bool, error) {
	return _IdentityRegistry.Contract.Contains(&_IdentityRegistry.CallOpts, _userAddress)
}

// Contains is a free data retrieval call binding the contract method 0x5dbe47e8.
//
// Solidity: function contains(address _userAddress) view returns(bool)
func (_IdentityRegistry *IdentityRegistryCallerSession) Contains(_userAddress common.Address) (bool, error) {
	return _IdentityRegistry.Contract.Contains(&_IdentityRegistry.CallOpts, _userAddress)
}

// Identity is a free data retrieval call binding the contract method 0xf0eb5e54.
//
// Solidity: function identity(address _userAddress) view returns(address)
func (_IdentityRegistry *IdentityRegistryCaller) Identity(opts *bind.CallOpts, _userAddress common.Address) (common.Address, error) {
	var out []interface{}
	err := _IdentityRegistry.contract.Call(opts, &out, "identity", _userAddress)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Identity is a free data retrieval call binding the contract method 0xf0eb5e54.
//
// Solidity: function identity(address _userAddress) view returns(address)
func (_IdentityRegistry *IdentityRegistrySession) Identity(_userAddress common.Address) (common.Address, error) {
	return _IdentityRegistry.Contract.Identity(&_IdentityRegistry.CallOpts, _userAddress)
}

// Identity is a free data retrieval call binding the contract method 0xf0eb5e54.
//
// Solidity: function identity(address _userAddress) view returns(address)
func (_IdentityRegistry *IdentityRegistryCallerSession) Identity(_userAddress common.Address) (common.Address, error) {
	return _IdentityRegistry.Contract.Identity(&_IdentityRegistry.CallOpts, _userAddress)
}

// IdentityStorage is a free data retrieval call binding the contract method 0xf11abfd8.
//
// Solidity: function identityStorage() view returns(address)
func (_IdentityRegistry *IdentityRegistryCaller) IdentityStorage(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IdentityRegistry.contract.Call(opts, &out, "identityStorage")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// IdentityStorage is a free data retrieval call binding the contract method 0xf11abfd8.
//
// Solidity: function identityStorage() view returns(address)
func (_IdentityRegistry *IdentityRegistrySession) IdentityStorage() (common.Address, error) {
	return _IdentityRegistry.Contract.IdentityStorage(&_IdentityRegistry.CallOpts)
}

// IdentityStorage is a free data retrieval call binding the contract method 0xf11abfd8.
//
// Solidity: function identityStorage() view returns(address)
func (_IdentityRegistry *IdentityRegistryCallerSession) IdentityStorage() (common.Address, error) {
	return _IdentityRegistry.Contract.IdentityStorage(&_IdentityRegistry.CallOpts)
}

// InvestorCountry is a free data retrieval call binding the contract method 0x7e42683b.
//
// Solidity: function investorCountry(address _userAddress) view returns(uint16)
func (_IdentityRegistry *IdentityRegistryCaller) InvestorCountry(opts *bind.CallOpts, _userAddress common.Address) (uint16, error) {
	var out []interface{}
	err := _IdentityRegistry.contract.Call(opts, &out, "investorCountry", _userAddress)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// InvestorCountry is a free data retrieval call binding the contract method 0x7e42683b.
//
// Solidity: function investorCountry(address _userAddress) view returns(uint16)
func (_IdentityRegistry *IdentityRegistrySession) InvestorCountry(_userAddress common.Address) (uint16, error) {
	return _IdentityRegistry.Contract.InvestorCountry(&_IdentityRegistry.CallOpts, _userAddress)
}

// InvestorCountry is a free data retrieval call binding the contract method 0x7e42683b.
//
// Solidity: function investorCountry(address _userAddress) view returns(uint16)
func (_IdentityRegistry *IdentityRegistryCallerSession) InvestorCountry(_userAddress common.Address) (uint16, error) {
	return _IdentityRegistry.Contract.InvestorCountry(&_IdentityRegistry.CallOpts, _userAddress)
}

// IsAgent is a free data retrieval call binding the contract method 0x1ffbb064.
//
// Solidity: function isAgent(address _agent) view returns(bool)
func (_IdentityRegistry *IdentityRegistryCaller) IsAgent(opts *bind.CallOpts, _agent common.Address) (bool, error) {
	var out []interface{}
	err := _IdentityRegistry.contract.Call(opts, &out, "isAgent", _agent)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAgent is a free data retrieval call binding the contract method 0x1ffbb064.
//
// Solidity: function isAgent(address _agent) view returns(bool)
func (_IdentityRegistry *IdentityRegistrySession) IsAgent(_agent common.Address) (bool, error) {
	return _IdentityRegistry.Contract.IsAgent(&_IdentityRegistry.CallOpts, _agent)
}

// IsAgent is a free data retrieval call binding the contract method 0x1ffbb064.
//
// Solidity: function isAgent(address _agent) view returns(bool)
func (_IdentityRegistry *IdentityRegistryCallerSession) IsAgent(_agent common.Address) (bool, error) {
	return _IdentityRegistry.Contract.IsAgent(&_IdentityRegistry.CallOpts, _agent)
}

// IsVerified is a free data retrieval call binding the contract method 0xb9209e33.
//
// Solidity: function isVerified(address _userAddress) view returns(bool)
func (_IdentityRegistry *IdentityRegistryCaller) IsVerified(opts *bind.CallOpts, _userAddress common.Address) (bool, error) {
	var out []interface{}
	err := _IdentityRegistry.contract.Call(opts, &out, "isVerified", _userAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsVerified is a free data retrieval call binding the contract method 0xb9209e33.
//
// Solidity: function isVerified(address _userAddress) view returns(bool)
func (_IdentityRegistry *IdentityRegistrySession) IsVerified(_userAddress common.Address) (bool, error) {
	return _IdentityRegistry.Contract.IsVerified(&_IdentityRegistry.CallOpts, _userAddress)
}

// IsVerified is a free data retrieval call binding the contract method 0xb9209e33.
//
// Solidity: function isVerified(address _userAddress) view returns(bool)
func (_IdentityRegistry *IdentityRegistryCallerSession) IsVerified(_userAddress common.Address) (bool, error) {
	return _IdentityRegistry.Contract.IsVerified(&_IdentityRegistry.CallOpts, _userAddress)
}

// IssuersRegistry is a free data retrieval call binding the contract method 0xb4f3fcb7.
//
// Solidity: function issuersRegistry() view returns(address)
func (_IdentityRegistry *IdentityRegistryCaller) IssuersRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IdentityRegistry.contract.Call(opts, &out, "issuersRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// IssuersRegistry is a free data retrieval call binding the contract method 0xb4f3fcb7.
//
// Solidity: function issuersRegistry() view returns(address)
func (_IdentityRegistry *IdentityRegistrySession) IssuersRegistry() (common.Address, error) {
	return _IdentityRegistry.Contract.IssuersRegistry(&_IdentityRegistry.CallOpts)
}

// IssuersRegistry is a free data retrieval call binding the contract method 0xb4f3fcb7.
//
// Solidity: function issuersRegistry() view returns(address)
func (_IdentityRegistry *IdentityRegistryCallerSession) IssuersRegistry() (common.Address, error) {
	return _IdentityRegistry.Contract.IssuersRegistry(&_IdentityRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IdentityRegistry *IdentityRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IdentityRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IdentityRegistry *IdentityRegistrySession) Owner() (common.Address, error) {
	return _IdentityRegistry.Contract.Owner(&_IdentityRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IdentityRegistry *IdentityRegistryCallerSession) Owner() (common.Address, error) {
	return _IdentityRegistry.Contract.Owner(&_IdentityRegistry.CallOpts)
}

// TopicsRegistry is a free data retrieval call binding the contract method 0x3b3e12f4.
//
// Solidity: function topicsRegistry() view returns(address)
func (_IdentityRegistry *IdentityRegistryCaller) TopicsRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IdentityRegistry.contract.Call(opts, &out, "topicsRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TopicsRegistry is a free data retrieval call binding the contract method 0x3b3e12f4.
//
// Solidity: function topicsRegistry() view returns(address)
func (_IdentityRegistry *IdentityRegistrySession) TopicsRegistry() (common.Address, error) {
	return _IdentityRegistry.Contract.TopicsRegistry(&_IdentityRegistry.CallOpts)
}

// TopicsRegistry is a free data retrieval call binding the contract method 0x3b3e12f4.
//
// Solidity: function topicsRegistry() view returns(address)
func (_IdentityRegistry *IdentityRegistryCallerSession) TopicsRegistry() (common.Address, error) {
	return _IdentityRegistry.Contract.TopicsRegistry(&_IdentityRegistry.CallOpts)
}

// AddAgent is a paid mutator transaction binding the contract method 0x84e79842.
//
// Solidity: function addAgent(address _agent) returns()
func (_IdentityRegistry *IdentityRegistryTransactor) AddAgent(opts *bind.TransactOpts, _agent common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.contract.Transact(opts, "addAgent", _agent)
}

// AddAgent is a paid mutator transaction binding the contract method 0x84e79842.
//
// Solidity: function addAgent(address _agent) returns()
func (_IdentityRegistry *IdentityRegistrySession) AddAgent(_agent common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.AddAgent(&_IdentityRegistry.TransactOpts, _agent)
}

// AddAgent is a paid mutator transaction binding the contract method 0x84e79842.
//
// Solidity: function addAgent(address _agent) returns()
func (_IdentityRegistry *IdentityRegistryTransactorSession) AddAgent(_agent common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.AddAgent(&_IdentityRegistry.TransactOpts, _agent)
}

// BatchRegisterIdentity is a paid mutator transaction binding the contract method 0x653dc9f1.
//
// Solidity: function batchRegisterIdentity(address[] _userAddresses, address[] _identities, uint16[] _countries) returns()
func (_IdentityRegistry *IdentityRegistryTransactor) BatchRegisterIdentity(opts *bind.TransactOpts, _userAddresses []common.Address, _identities []common.Address, _countries []uint16) (*types.Transaction, error) {
	return _IdentityRegistry.contract.Transact(opts, "batchRegisterIdentity", _userAddresses, _identities, _countries)
}

// BatchRegisterIdentity is a paid mutator transaction binding the contract method 0x653dc9f1.
//
// Solidity: function batchRegisterIdentity(address[] _userAddresses, address[] _identities, uint16[] _countries) returns()
func (_IdentityRegistry *IdentityRegistrySession) BatchRegisterIdentity(_userAddresses []common.Address, _identities []common.Address, _countries []uint16) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.BatchRegisterIdentity(&_IdentityRegistry.TransactOpts, _userAddresses, _identities, _countries)
}

// BatchRegisterIdentity is a paid mutator transaction binding the contract method 0x653dc9f1.
//
// Solidity: function batchRegisterIdentity(address[] _userAddresses, address[] _identities, uint16[] _countries) returns()
func (_IdentityRegistry *IdentityRegistryTransactorSession) BatchRegisterIdentity(_userAddresses []common.Address, _identities []common.Address, _countries []uint16) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.BatchRegisterIdentity(&_IdentityRegistry.TransactOpts, _userAddresses, _identities, _countries)
}

// DeleteIdentity is a paid mutator transaction binding the contract method 0xa8d29d1d.
//
// Solidity: function deleteIdentity(address _userAddress) returns()
func (_IdentityRegistry *IdentityRegistryTransactor) DeleteIdentity(opts *bind.TransactOpts, _userAddress common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.contract.Transact(opts, "deleteIdentity", _userAddress)
}

// DeleteIdentity is a paid mutator transaction binding the contract method 0xa8d29d1d.
//
// Solidity: function deleteIdentity(address _userAddress) returns()
func (_IdentityRegistry *IdentityRegistrySession) DeleteIdentity(_userAddress common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.DeleteIdentity(&_IdentityRegistry.TransactOpts, _userAddress)
}

// DeleteIdentity is a paid mutator transaction binding the contract method 0xa8d29d1d.
//
// Solidity: function deleteIdentity(address _userAddress) returns()
func (_IdentityRegistry *IdentityRegistryTransactorSession) DeleteIdentity(_userAddress common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.DeleteIdentity(&_IdentityRegistry.TransactOpts, _userAddress)
}

// Init is a paid mutator transaction binding the contract method 0x184b9559.
//
// Solidity: function init(address _trustedIssuersRegistry, address _claimTopicsRegistry, address _identityStorage) returns()
func (_IdentityRegistry *IdentityRegistryTransactor) Init(opts *bind.TransactOpts, _trustedIssuersRegistry common.Address, _claimTopicsRegistry common.Address, _identityStorage common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.contract.Transact(opts, "init", _trustedIssuersRegistry, _claimTopicsRegistry, _identityStorage)
}

// Init is a paid mutator transaction binding the contract method 0x184b9559.
//
// Solidity: function init(address _trustedIssuersRegistry, address _claimTopicsRegistry, address _identityStorage) returns()
func (_IdentityRegistry *IdentityRegistrySession) Init(_trustedIssuersRegistry common.Address, _claimTopicsRegistry common.Address, _identityStorage common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.Init(&_IdentityRegistry.TransactOpts, _trustedIssuersRegistry, _claimTopicsRegistry, _identityStorage)
}

// Init is a paid mutator transaction binding the contract method 0x184b9559.
//
// Solidity: function init(address _trustedIssuersRegistry, address _claimTopicsRegistry, address _identityStorage) returns()
func (_IdentityRegistry *IdentityRegistryTransactorSession) Init(_trustedIssuersRegistry common.Address, _claimTopicsRegistry common.Address, _identityStorage common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.Init(&_IdentityRegistry.TransactOpts, _trustedIssuersRegistry, _claimTopicsRegistry, _identityStorage)
}

// RegisterIdentity is a paid mutator transaction binding the contract method 0x454a03e0.
//
// Solidity: function registerIdentity(address _userAddress, address _identity, uint16 _country) returns()
func (_IdentityRegistry *IdentityRegistryTransactor) RegisterIdentity(opts *bind.TransactOpts, _userAddress common.Address, _identity common.Address, _country uint16) (*types.Transaction, error) {
	return _IdentityRegistry.contract.Transact(opts, "registerIdentity", _userAddress, _identity, _country)
}

// RegisterIdentity is a paid mutator transaction binding the contract method 0x454a03e0.
//
// Solidity: function registerIdentity(address _userAddress, address _identity, uint16 _country) returns()
func (_IdentityRegistry *IdentityRegistrySession) RegisterIdentity(_userAddress common.Address, _identity common.Address, _country uint16) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.RegisterIdentity(&_IdentityRegistry.TransactOpts, _userAddress, _identity, _country)
}

// RegisterIdentity is a paid mutator transaction binding the contract method 0x454a03e0.
//
// Solidity: function registerIdentity(address _userAddress, address _identity, uint16 _country) returns()
func (_IdentityRegistry *IdentityRegistryTransactorSession) RegisterIdentity(_userAddress common.Address, _identity common.Address, _country uint16) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.RegisterIdentity(&_IdentityRegistry.TransactOpts, _userAddress, _identity, _country)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x97a6278e.
//
// Solidity: function removeAgent(address _agent) returns()
func (_IdentityRegistry *IdentityRegistryTransactor) RemoveAgent(opts *bind.TransactOpts, _agent common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.contract.Transact(opts, "removeAgent", _agent)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x97a6278e.
//
// Solidity: function removeAgent(address _agent) returns()
func (_IdentityRegistry *IdentityRegistrySession) RemoveAgent(_agent common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.RemoveAgent(&_IdentityRegistry.TransactOpts, _agent)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x97a6278e.
//
// Solidity: function removeAgent(address _agent) returns()
func (_IdentityRegistry *IdentityRegistryTransactorSession) RemoveAgent(_agent common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.RemoveAgent(&_IdentityRegistry.TransactOpts, _agent)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IdentityRegistry *IdentityRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IdentityRegistry *IdentityRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _IdentityRegistry.Contract.RenounceOwnership(&_IdentityRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IdentityRegistry *IdentityRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _IdentityRegistry.Contract.RenounceOwnership(&_IdentityRegistry.TransactOpts)
}

// SetClaimTopicsRegistry is a paid mutator transaction binding the contract method 0x670af6a9.
//
// Solidity: function setClaimTopicsRegistry(address _claimTopicsRegistry) returns()
func (_IdentityRegistry *IdentityRegistryTransactor) SetClaimTopicsRegistry(opts *bind.TransactOpts, _claimTopicsRegistry common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.contract.Transact(opts, "setClaimTopicsRegistry", _claimTopicsRegistry)
}

// SetClaimTopicsRegistry is a paid mutator transaction binding the contract method 0x670af6a9.
//
// Solidity: function setClaimTopicsRegistry(address _claimTopicsRegistry) returns()
func (_IdentityRegistry *IdentityRegistrySession) SetClaimTopicsRegistry(_claimTopicsRegistry common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.SetClaimTopicsRegistry(&_IdentityRegistry.TransactOpts, _claimTopicsRegistry)
}

// SetClaimTopicsRegistry is a paid mutator transaction binding the contract method 0x670af6a9.
//
// Solidity: function setClaimTopicsRegistry(address _claimTopicsRegistry) returns()
func (_IdentityRegistry *IdentityRegistryTransactorSession) SetClaimTopicsRegistry(_claimTopicsRegistry common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.SetClaimTopicsRegistry(&_IdentityRegistry.TransactOpts, _claimTopicsRegistry)
}

// SetIdentityRegistryStorage is a paid mutator transaction binding the contract method 0x26d941ae.
//
// Solidity: function setIdentityRegistryStorage(address _identityRegistryStorage) returns()
func (_IdentityRegistry *IdentityRegistryTransactor) SetIdentityRegistryStorage(opts *bind.TransactOpts, _identityRegistryStorage common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.contract.Transact(opts, "setIdentityRegistryStorage", _identityRegistryStorage)
}

// SetIdentityRegistryStorage is a paid mutator transaction binding the contract method 0x26d941ae.
//
// Solidity: function setIdentityRegistryStorage(address _identityRegistryStorage) returns()
func (_IdentityRegistry *IdentityRegistrySession) SetIdentityRegistryStorage(_identityRegistryStorage common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.SetIdentityRegistryStorage(&_IdentityRegistry.TransactOpts, _identityRegistryStorage)
}

// SetIdentityRegistryStorage is a paid mutator transaction binding the contract method 0x26d941ae.
//
// Solidity: function setIdentityRegistryStorage(address _identityRegistryStorage) returns()
func (_IdentityRegistry *IdentityRegistryTransactorSession) SetIdentityRegistryStorage(_identityRegistryStorage common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.SetIdentityRegistryStorage(&_IdentityRegistry.TransactOpts, _identityRegistryStorage)
}

// SetTrustedIssuersRegistry is a paid mutator transaction binding the contract method 0xe744d789.
//
// Solidity: function setTrustedIssuersRegistry(address _trustedIssuersRegistry) returns()
func (_IdentityRegistry *IdentityRegistryTransactor) SetTrustedIssuersRegistry(opts *bind.TransactOpts, _trustedIssuersRegistry common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.contract.Transact(opts, "setTrustedIssuersRegistry", _trustedIssuersRegistry)
}

// SetTrustedIssuersRegistry is a paid mutator transaction binding the contract method 0xe744d789.
//
// Solidity: function setTrustedIssuersRegistry(address _trustedIssuersRegistry) returns()
func (_IdentityRegistry *IdentityRegistrySession) SetTrustedIssuersRegistry(_trustedIssuersRegistry common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.SetTrustedIssuersRegistry(&_IdentityRegistry.TransactOpts, _trustedIssuersRegistry)
}

// SetTrustedIssuersRegistry is a paid mutator transaction binding the contract method 0xe744d789.
//
// Solidity: function setTrustedIssuersRegistry(address _trustedIssuersRegistry) returns()
func (_IdentityRegistry *IdentityRegistryTransactorSession) SetTrustedIssuersRegistry(_trustedIssuersRegistry common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.SetTrustedIssuersRegistry(&_IdentityRegistry.TransactOpts, _trustedIssuersRegistry)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IdentityRegistry *IdentityRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IdentityRegistry *IdentityRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.TransferOwnership(&_IdentityRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IdentityRegistry *IdentityRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.TransferOwnership(&_IdentityRegistry.TransactOpts, newOwner)
}

// UpdateCountry is a paid mutator transaction binding the contract method 0x3b239a7f.
//
// Solidity: function updateCountry(address _userAddress, uint16 _country) returns()
func (_IdentityRegistry *IdentityRegistryTransactor) UpdateCountry(opts *bind.TransactOpts, _userAddress common.Address, _country uint16) (*types.Transaction, error) {
	return _IdentityRegistry.contract.Transact(opts, "updateCountry", _userAddress, _country)
}

// UpdateCountry is a paid mutator transaction binding the contract method 0x3b239a7f.
//
// Solidity: function updateCountry(address _userAddress, uint16 _country) returns()
func (_IdentityRegistry *IdentityRegistrySession) UpdateCountry(_userAddress common.Address, _country uint16) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.UpdateCountry(&_IdentityRegistry.TransactOpts, _userAddress, _country)
}

// UpdateCountry is a paid mutator transaction binding the contract method 0x3b239a7f.
//
// Solidity: function updateCountry(address _userAddress, uint16 _country) returns()
func (_IdentityRegistry *IdentityRegistryTransactorSession) UpdateCountry(_userAddress common.Address, _country uint16) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.UpdateCountry(&_IdentityRegistry.TransactOpts, _userAddress, _country)
}

// UpdateIdentity is a paid mutator transaction binding the contract method 0x8e098ca1.
//
// Solidity: function updateIdentity(address _userAddress, address _identity) returns()
func (_IdentityRegistry *IdentityRegistryTransactor) UpdateIdentity(opts *bind.TransactOpts, _userAddress common.Address, _identity common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.contract.Transact(opts, "updateIdentity", _userAddress, _identity)
}

// UpdateIdentity is a paid mutator transaction binding the contract method 0x8e098ca1.
//
// Solidity: function updateIdentity(address _userAddress, address _identity) returns()
func (_IdentityRegistry *IdentityRegistrySession) UpdateIdentity(_userAddress common.Address, _identity common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.UpdateIdentity(&_IdentityRegistry.TransactOpts, _userAddress, _identity)
}

// UpdateIdentity is a paid mutator transaction binding the contract method 0x8e098ca1.
//
// Solidity: function updateIdentity(address _userAddress, address _identity) returns()
func (_IdentityRegistry *IdentityRegistryTransactorSession) UpdateIdentity(_userAddress common.Address, _identity common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.UpdateIdentity(&_IdentityRegistry.TransactOpts, _userAddress, _identity)
}

// IdentityRegistryAgentAddedIterator is returned from FilterAgentAdded and is used to iterate over the raw logs and unpacked data for AgentAdded events raised by the IdentityRegistry contract.
type IdentityRegistryAgentAddedIterator struct {
	Event *IdentityRegistryAgentAdded // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryAgentAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryAgentAdded)
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
		it.Event = new(IdentityRegistryAgentAdded)
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
func (it *IdentityRegistryAgentAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryAgentAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryAgentAdded represents a AgentAdded event raised by the IdentityRegistry contract.
type IdentityRegistryAgentAdded struct {
	Agent common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterAgentAdded is a free log retrieval operation binding the contract event 0xf68e73cec97f2d70aa641fb26e87a4383686e2efacb648f2165aeb02ac562ec5.
//
// Solidity: event AgentAdded(address indexed _agent)
func (_IdentityRegistry *IdentityRegistryFilterer) FilterAgentAdded(opts *bind.FilterOpts, _agent []common.Address) (*IdentityRegistryAgentAddedIterator, error) {

	var _agentRule []interface{}
	for _, _agentItem := range _agent {
		_agentRule = append(_agentRule, _agentItem)
	}

	logs, sub, err := _IdentityRegistry.contract.FilterLogs(opts, "AgentAdded", _agentRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryAgentAddedIterator{contract: _IdentityRegistry.contract, event: "AgentAdded", logs: logs, sub: sub}, nil
}

// WatchAgentAdded is a free log subscription operation binding the contract event 0xf68e73cec97f2d70aa641fb26e87a4383686e2efacb648f2165aeb02ac562ec5.
//
// Solidity: event AgentAdded(address indexed _agent)
func (_IdentityRegistry *IdentityRegistryFilterer) WatchAgentAdded(opts *bind.WatchOpts, sink chan<- *IdentityRegistryAgentAdded, _agent []common.Address) (event.Subscription, error) {

	var _agentRule []interface{}
	for _, _agentItem := range _agent {
		_agentRule = append(_agentRule, _agentItem)
	}

	logs, sub, err := _IdentityRegistry.contract.WatchLogs(opts, "AgentAdded", _agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryAgentAdded)
				if err := _IdentityRegistry.contract.UnpackLog(event, "AgentAdded", log); err != nil {
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
func (_IdentityRegistry *IdentityRegistryFilterer) ParseAgentAdded(log types.Log) (*IdentityRegistryAgentAdded, error) {
	event := new(IdentityRegistryAgentAdded)
	if err := _IdentityRegistry.contract.UnpackLog(event, "AgentAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryAgentRemovedIterator is returned from FilterAgentRemoved and is used to iterate over the raw logs and unpacked data for AgentRemoved events raised by the IdentityRegistry contract.
type IdentityRegistryAgentRemovedIterator struct {
	Event *IdentityRegistryAgentRemoved // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryAgentRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryAgentRemoved)
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
		it.Event = new(IdentityRegistryAgentRemoved)
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
func (it *IdentityRegistryAgentRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryAgentRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryAgentRemoved represents a AgentRemoved event raised by the IdentityRegistry contract.
type IdentityRegistryAgentRemoved struct {
	Agent common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterAgentRemoved is a free log retrieval operation binding the contract event 0xed9c8ad8d5a0a66898ea49d2956929c93ae2e8bd50281b2ed897c5d1a6737e0b.
//
// Solidity: event AgentRemoved(address indexed _agent)
func (_IdentityRegistry *IdentityRegistryFilterer) FilterAgentRemoved(opts *bind.FilterOpts, _agent []common.Address) (*IdentityRegistryAgentRemovedIterator, error) {

	var _agentRule []interface{}
	for _, _agentItem := range _agent {
		_agentRule = append(_agentRule, _agentItem)
	}

	logs, sub, err := _IdentityRegistry.contract.FilterLogs(opts, "AgentRemoved", _agentRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryAgentRemovedIterator{contract: _IdentityRegistry.contract, event: "AgentRemoved", logs: logs, sub: sub}, nil
}

// WatchAgentRemoved is a free log subscription operation binding the contract event 0xed9c8ad8d5a0a66898ea49d2956929c93ae2e8bd50281b2ed897c5d1a6737e0b.
//
// Solidity: event AgentRemoved(address indexed _agent)
func (_IdentityRegistry *IdentityRegistryFilterer) WatchAgentRemoved(opts *bind.WatchOpts, sink chan<- *IdentityRegistryAgentRemoved, _agent []common.Address) (event.Subscription, error) {

	var _agentRule []interface{}
	for _, _agentItem := range _agent {
		_agentRule = append(_agentRule, _agentItem)
	}

	logs, sub, err := _IdentityRegistry.contract.WatchLogs(opts, "AgentRemoved", _agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryAgentRemoved)
				if err := _IdentityRegistry.contract.UnpackLog(event, "AgentRemoved", log); err != nil {
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
func (_IdentityRegistry *IdentityRegistryFilterer) ParseAgentRemoved(log types.Log) (*IdentityRegistryAgentRemoved, error) {
	event := new(IdentityRegistryAgentRemoved)
	if err := _IdentityRegistry.contract.UnpackLog(event, "AgentRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryClaimTopicsRegistrySetIterator is returned from FilterClaimTopicsRegistrySet and is used to iterate over the raw logs and unpacked data for ClaimTopicsRegistrySet events raised by the IdentityRegistry contract.
type IdentityRegistryClaimTopicsRegistrySetIterator struct {
	Event *IdentityRegistryClaimTopicsRegistrySet // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryClaimTopicsRegistrySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryClaimTopicsRegistrySet)
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
		it.Event = new(IdentityRegistryClaimTopicsRegistrySet)
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
func (it *IdentityRegistryClaimTopicsRegistrySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryClaimTopicsRegistrySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryClaimTopicsRegistrySet represents a ClaimTopicsRegistrySet event raised by the IdentityRegistry contract.
type IdentityRegistryClaimTopicsRegistrySet struct {
	ClaimTopicsRegistry common.Address
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterClaimTopicsRegistrySet is a free log retrieval operation binding the contract event 0x7170bf15b246e880b2369cd7c67d057760d8a35149e8c64dde91efa22bcc76d0.
//
// Solidity: event ClaimTopicsRegistrySet(address indexed claimTopicsRegistry)
func (_IdentityRegistry *IdentityRegistryFilterer) FilterClaimTopicsRegistrySet(opts *bind.FilterOpts, claimTopicsRegistry []common.Address) (*IdentityRegistryClaimTopicsRegistrySetIterator, error) {

	var claimTopicsRegistryRule []interface{}
	for _, claimTopicsRegistryItem := range claimTopicsRegistry {
		claimTopicsRegistryRule = append(claimTopicsRegistryRule, claimTopicsRegistryItem)
	}

	logs, sub, err := _IdentityRegistry.contract.FilterLogs(opts, "ClaimTopicsRegistrySet", claimTopicsRegistryRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryClaimTopicsRegistrySetIterator{contract: _IdentityRegistry.contract, event: "ClaimTopicsRegistrySet", logs: logs, sub: sub}, nil
}

// WatchClaimTopicsRegistrySet is a free log subscription operation binding the contract event 0x7170bf15b246e880b2369cd7c67d057760d8a35149e8c64dde91efa22bcc76d0.
//
// Solidity: event ClaimTopicsRegistrySet(address indexed claimTopicsRegistry)
func (_IdentityRegistry *IdentityRegistryFilterer) WatchClaimTopicsRegistrySet(opts *bind.WatchOpts, sink chan<- *IdentityRegistryClaimTopicsRegistrySet, claimTopicsRegistry []common.Address) (event.Subscription, error) {

	var claimTopicsRegistryRule []interface{}
	for _, claimTopicsRegistryItem := range claimTopicsRegistry {
		claimTopicsRegistryRule = append(claimTopicsRegistryRule, claimTopicsRegistryItem)
	}

	logs, sub, err := _IdentityRegistry.contract.WatchLogs(opts, "ClaimTopicsRegistrySet", claimTopicsRegistryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryClaimTopicsRegistrySet)
				if err := _IdentityRegistry.contract.UnpackLog(event, "ClaimTopicsRegistrySet", log); err != nil {
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

// ParseClaimTopicsRegistrySet is a log parse operation binding the contract event 0x7170bf15b246e880b2369cd7c67d057760d8a35149e8c64dde91efa22bcc76d0.
//
// Solidity: event ClaimTopicsRegistrySet(address indexed claimTopicsRegistry)
func (_IdentityRegistry *IdentityRegistryFilterer) ParseClaimTopicsRegistrySet(log types.Log) (*IdentityRegistryClaimTopicsRegistrySet, error) {
	event := new(IdentityRegistryClaimTopicsRegistrySet)
	if err := _IdentityRegistry.contract.UnpackLog(event, "ClaimTopicsRegistrySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryCountryUpdatedIterator is returned from FilterCountryUpdated and is used to iterate over the raw logs and unpacked data for CountryUpdated events raised by the IdentityRegistry contract.
type IdentityRegistryCountryUpdatedIterator struct {
	Event *IdentityRegistryCountryUpdated // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryCountryUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryCountryUpdated)
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
		it.Event = new(IdentityRegistryCountryUpdated)
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
func (it *IdentityRegistryCountryUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryCountryUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryCountryUpdated represents a CountryUpdated event raised by the IdentityRegistry contract.
type IdentityRegistryCountryUpdated struct {
	InvestorAddress common.Address
	Country         uint16
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterCountryUpdated is a free log retrieval operation binding the contract event 0x04ed3b726495c2dca1ff1215d9ca54e1a4030abb5e82b0f6ce55702416cee853.
//
// Solidity: event CountryUpdated(address indexed investorAddress, uint16 indexed country)
func (_IdentityRegistry *IdentityRegistryFilterer) FilterCountryUpdated(opts *bind.FilterOpts, investorAddress []common.Address, country []uint16) (*IdentityRegistryCountryUpdatedIterator, error) {

	var investorAddressRule []interface{}
	for _, investorAddressItem := range investorAddress {
		investorAddressRule = append(investorAddressRule, investorAddressItem)
	}
	var countryRule []interface{}
	for _, countryItem := range country {
		countryRule = append(countryRule, countryItem)
	}

	logs, sub, err := _IdentityRegistry.contract.FilterLogs(opts, "CountryUpdated", investorAddressRule, countryRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryCountryUpdatedIterator{contract: _IdentityRegistry.contract, event: "CountryUpdated", logs: logs, sub: sub}, nil
}

// WatchCountryUpdated is a free log subscription operation binding the contract event 0x04ed3b726495c2dca1ff1215d9ca54e1a4030abb5e82b0f6ce55702416cee853.
//
// Solidity: event CountryUpdated(address indexed investorAddress, uint16 indexed country)
func (_IdentityRegistry *IdentityRegistryFilterer) WatchCountryUpdated(opts *bind.WatchOpts, sink chan<- *IdentityRegistryCountryUpdated, investorAddress []common.Address, country []uint16) (event.Subscription, error) {

	var investorAddressRule []interface{}
	for _, investorAddressItem := range investorAddress {
		investorAddressRule = append(investorAddressRule, investorAddressItem)
	}
	var countryRule []interface{}
	for _, countryItem := range country {
		countryRule = append(countryRule, countryItem)
	}

	logs, sub, err := _IdentityRegistry.contract.WatchLogs(opts, "CountryUpdated", investorAddressRule, countryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryCountryUpdated)
				if err := _IdentityRegistry.contract.UnpackLog(event, "CountryUpdated", log); err != nil {
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

// ParseCountryUpdated is a log parse operation binding the contract event 0x04ed3b726495c2dca1ff1215d9ca54e1a4030abb5e82b0f6ce55702416cee853.
//
// Solidity: event CountryUpdated(address indexed investorAddress, uint16 indexed country)
func (_IdentityRegistry *IdentityRegistryFilterer) ParseCountryUpdated(log types.Log) (*IdentityRegistryCountryUpdated, error) {
	event := new(IdentityRegistryCountryUpdated)
	if err := _IdentityRegistry.contract.UnpackLog(event, "CountryUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryIdentityRegisteredIterator is returned from FilterIdentityRegistered and is used to iterate over the raw logs and unpacked data for IdentityRegistered events raised by the IdentityRegistry contract.
type IdentityRegistryIdentityRegisteredIterator struct {
	Event *IdentityRegistryIdentityRegistered // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryIdentityRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryIdentityRegistered)
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
		it.Event = new(IdentityRegistryIdentityRegistered)
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
func (it *IdentityRegistryIdentityRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryIdentityRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryIdentityRegistered represents a IdentityRegistered event raised by the IdentityRegistry contract.
type IdentityRegistryIdentityRegistered struct {
	InvestorAddress common.Address
	Identity        common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterIdentityRegistered is a free log retrieval operation binding the contract event 0x6ae73635c50d24a45af6fbd5e016ac4bed179addbc8bf24e04ff0fcc6d33af19.
//
// Solidity: event IdentityRegistered(address indexed investorAddress, address indexed identity)
func (_IdentityRegistry *IdentityRegistryFilterer) FilterIdentityRegistered(opts *bind.FilterOpts, investorAddress []common.Address, identity []common.Address) (*IdentityRegistryIdentityRegisteredIterator, error) {

	var investorAddressRule []interface{}
	for _, investorAddressItem := range investorAddress {
		investorAddressRule = append(investorAddressRule, investorAddressItem)
	}
	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _IdentityRegistry.contract.FilterLogs(opts, "IdentityRegistered", investorAddressRule, identityRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryIdentityRegisteredIterator{contract: _IdentityRegistry.contract, event: "IdentityRegistered", logs: logs, sub: sub}, nil
}

// WatchIdentityRegistered is a free log subscription operation binding the contract event 0x6ae73635c50d24a45af6fbd5e016ac4bed179addbc8bf24e04ff0fcc6d33af19.
//
// Solidity: event IdentityRegistered(address indexed investorAddress, address indexed identity)
func (_IdentityRegistry *IdentityRegistryFilterer) WatchIdentityRegistered(opts *bind.WatchOpts, sink chan<- *IdentityRegistryIdentityRegistered, investorAddress []common.Address, identity []common.Address) (event.Subscription, error) {

	var investorAddressRule []interface{}
	for _, investorAddressItem := range investorAddress {
		investorAddressRule = append(investorAddressRule, investorAddressItem)
	}
	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _IdentityRegistry.contract.WatchLogs(opts, "IdentityRegistered", investorAddressRule, identityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryIdentityRegistered)
				if err := _IdentityRegistry.contract.UnpackLog(event, "IdentityRegistered", log); err != nil {
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

// ParseIdentityRegistered is a log parse operation binding the contract event 0x6ae73635c50d24a45af6fbd5e016ac4bed179addbc8bf24e04ff0fcc6d33af19.
//
// Solidity: event IdentityRegistered(address indexed investorAddress, address indexed identity)
func (_IdentityRegistry *IdentityRegistryFilterer) ParseIdentityRegistered(log types.Log) (*IdentityRegistryIdentityRegistered, error) {
	event := new(IdentityRegistryIdentityRegistered)
	if err := _IdentityRegistry.contract.UnpackLog(event, "IdentityRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryIdentityRemovedIterator is returned from FilterIdentityRemoved and is used to iterate over the raw logs and unpacked data for IdentityRemoved events raised by the IdentityRegistry contract.
type IdentityRegistryIdentityRemovedIterator struct {
	Event *IdentityRegistryIdentityRemoved // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryIdentityRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryIdentityRemoved)
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
		it.Event = new(IdentityRegistryIdentityRemoved)
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
func (it *IdentityRegistryIdentityRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryIdentityRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryIdentityRemoved represents a IdentityRemoved event raised by the IdentityRegistry contract.
type IdentityRegistryIdentityRemoved struct {
	InvestorAddress common.Address
	Identity        common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterIdentityRemoved is a free log retrieval operation binding the contract event 0x59d6590e225b81befe259af056324092801080acbb7feab310eb34678871f327.
//
// Solidity: event IdentityRemoved(address indexed investorAddress, address indexed identity)
func (_IdentityRegistry *IdentityRegistryFilterer) FilterIdentityRemoved(opts *bind.FilterOpts, investorAddress []common.Address, identity []common.Address) (*IdentityRegistryIdentityRemovedIterator, error) {

	var investorAddressRule []interface{}
	for _, investorAddressItem := range investorAddress {
		investorAddressRule = append(investorAddressRule, investorAddressItem)
	}
	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _IdentityRegistry.contract.FilterLogs(opts, "IdentityRemoved", investorAddressRule, identityRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryIdentityRemovedIterator{contract: _IdentityRegistry.contract, event: "IdentityRemoved", logs: logs, sub: sub}, nil
}

// WatchIdentityRemoved is a free log subscription operation binding the contract event 0x59d6590e225b81befe259af056324092801080acbb7feab310eb34678871f327.
//
// Solidity: event IdentityRemoved(address indexed investorAddress, address indexed identity)
func (_IdentityRegistry *IdentityRegistryFilterer) WatchIdentityRemoved(opts *bind.WatchOpts, sink chan<- *IdentityRegistryIdentityRemoved, investorAddress []common.Address, identity []common.Address) (event.Subscription, error) {

	var investorAddressRule []interface{}
	for _, investorAddressItem := range investorAddress {
		investorAddressRule = append(investorAddressRule, investorAddressItem)
	}
	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _IdentityRegistry.contract.WatchLogs(opts, "IdentityRemoved", investorAddressRule, identityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryIdentityRemoved)
				if err := _IdentityRegistry.contract.UnpackLog(event, "IdentityRemoved", log); err != nil {
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

// ParseIdentityRemoved is a log parse operation binding the contract event 0x59d6590e225b81befe259af056324092801080acbb7feab310eb34678871f327.
//
// Solidity: event IdentityRemoved(address indexed investorAddress, address indexed identity)
func (_IdentityRegistry *IdentityRegistryFilterer) ParseIdentityRemoved(log types.Log) (*IdentityRegistryIdentityRemoved, error) {
	event := new(IdentityRegistryIdentityRemoved)
	if err := _IdentityRegistry.contract.UnpackLog(event, "IdentityRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryIdentityStorageSetIterator is returned from FilterIdentityStorageSet and is used to iterate over the raw logs and unpacked data for IdentityStorageSet events raised by the IdentityRegistry contract.
type IdentityRegistryIdentityStorageSetIterator struct {
	Event *IdentityRegistryIdentityStorageSet // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryIdentityStorageSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryIdentityStorageSet)
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
		it.Event = new(IdentityRegistryIdentityStorageSet)
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
func (it *IdentityRegistryIdentityStorageSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryIdentityStorageSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryIdentityStorageSet represents a IdentityStorageSet event raised by the IdentityRegistry contract.
type IdentityRegistryIdentityStorageSet struct {
	IdentityStorage common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterIdentityStorageSet is a free log retrieval operation binding the contract event 0x2fa8b95c1db7afe99e3398f3792f008135cedc1fa26b0bb2ecd2352cd166d53c.
//
// Solidity: event IdentityStorageSet(address indexed identityStorage)
func (_IdentityRegistry *IdentityRegistryFilterer) FilterIdentityStorageSet(opts *bind.FilterOpts, identityStorage []common.Address) (*IdentityRegistryIdentityStorageSetIterator, error) {

	var identityStorageRule []interface{}
	for _, identityStorageItem := range identityStorage {
		identityStorageRule = append(identityStorageRule, identityStorageItem)
	}

	logs, sub, err := _IdentityRegistry.contract.FilterLogs(opts, "IdentityStorageSet", identityStorageRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryIdentityStorageSetIterator{contract: _IdentityRegistry.contract, event: "IdentityStorageSet", logs: logs, sub: sub}, nil
}

// WatchIdentityStorageSet is a free log subscription operation binding the contract event 0x2fa8b95c1db7afe99e3398f3792f008135cedc1fa26b0bb2ecd2352cd166d53c.
//
// Solidity: event IdentityStorageSet(address indexed identityStorage)
func (_IdentityRegistry *IdentityRegistryFilterer) WatchIdentityStorageSet(opts *bind.WatchOpts, sink chan<- *IdentityRegistryIdentityStorageSet, identityStorage []common.Address) (event.Subscription, error) {

	var identityStorageRule []interface{}
	for _, identityStorageItem := range identityStorage {
		identityStorageRule = append(identityStorageRule, identityStorageItem)
	}

	logs, sub, err := _IdentityRegistry.contract.WatchLogs(opts, "IdentityStorageSet", identityStorageRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryIdentityStorageSet)
				if err := _IdentityRegistry.contract.UnpackLog(event, "IdentityStorageSet", log); err != nil {
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

// ParseIdentityStorageSet is a log parse operation binding the contract event 0x2fa8b95c1db7afe99e3398f3792f008135cedc1fa26b0bb2ecd2352cd166d53c.
//
// Solidity: event IdentityStorageSet(address indexed identityStorage)
func (_IdentityRegistry *IdentityRegistryFilterer) ParseIdentityStorageSet(log types.Log) (*IdentityRegistryIdentityStorageSet, error) {
	event := new(IdentityRegistryIdentityStorageSet)
	if err := _IdentityRegistry.contract.UnpackLog(event, "IdentityStorageSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryIdentityUpdatedIterator is returned from FilterIdentityUpdated and is used to iterate over the raw logs and unpacked data for IdentityUpdated events raised by the IdentityRegistry contract.
type IdentityRegistryIdentityUpdatedIterator struct {
	Event *IdentityRegistryIdentityUpdated // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryIdentityUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryIdentityUpdated)
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
		it.Event = new(IdentityRegistryIdentityUpdated)
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
func (it *IdentityRegistryIdentityUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryIdentityUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryIdentityUpdated represents a IdentityUpdated event raised by the IdentityRegistry contract.
type IdentityRegistryIdentityUpdated struct {
	OldIdentity common.Address
	NewIdentity common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterIdentityUpdated is a free log retrieval operation binding the contract event 0xe98082932c8056a0f514da9104e4a66bc2cbaef102ad59d90c4b24220ebf6010.
//
// Solidity: event IdentityUpdated(address indexed oldIdentity, address indexed newIdentity)
func (_IdentityRegistry *IdentityRegistryFilterer) FilterIdentityUpdated(opts *bind.FilterOpts, oldIdentity []common.Address, newIdentity []common.Address) (*IdentityRegistryIdentityUpdatedIterator, error) {

	var oldIdentityRule []interface{}
	for _, oldIdentityItem := range oldIdentity {
		oldIdentityRule = append(oldIdentityRule, oldIdentityItem)
	}
	var newIdentityRule []interface{}
	for _, newIdentityItem := range newIdentity {
		newIdentityRule = append(newIdentityRule, newIdentityItem)
	}

	logs, sub, err := _IdentityRegistry.contract.FilterLogs(opts, "IdentityUpdated", oldIdentityRule, newIdentityRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryIdentityUpdatedIterator{contract: _IdentityRegistry.contract, event: "IdentityUpdated", logs: logs, sub: sub}, nil
}

// WatchIdentityUpdated is a free log subscription operation binding the contract event 0xe98082932c8056a0f514da9104e4a66bc2cbaef102ad59d90c4b24220ebf6010.
//
// Solidity: event IdentityUpdated(address indexed oldIdentity, address indexed newIdentity)
func (_IdentityRegistry *IdentityRegistryFilterer) WatchIdentityUpdated(opts *bind.WatchOpts, sink chan<- *IdentityRegistryIdentityUpdated, oldIdentity []common.Address, newIdentity []common.Address) (event.Subscription, error) {

	var oldIdentityRule []interface{}
	for _, oldIdentityItem := range oldIdentity {
		oldIdentityRule = append(oldIdentityRule, oldIdentityItem)
	}
	var newIdentityRule []interface{}
	for _, newIdentityItem := range newIdentity {
		newIdentityRule = append(newIdentityRule, newIdentityItem)
	}

	logs, sub, err := _IdentityRegistry.contract.WatchLogs(opts, "IdentityUpdated", oldIdentityRule, newIdentityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryIdentityUpdated)
				if err := _IdentityRegistry.contract.UnpackLog(event, "IdentityUpdated", log); err != nil {
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

// ParseIdentityUpdated is a log parse operation binding the contract event 0xe98082932c8056a0f514da9104e4a66bc2cbaef102ad59d90c4b24220ebf6010.
//
// Solidity: event IdentityUpdated(address indexed oldIdentity, address indexed newIdentity)
func (_IdentityRegistry *IdentityRegistryFilterer) ParseIdentityUpdated(log types.Log) (*IdentityRegistryIdentityUpdated, error) {
	event := new(IdentityRegistryIdentityUpdated)
	if err := _IdentityRegistry.contract.UnpackLog(event, "IdentityUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the IdentityRegistry contract.
type IdentityRegistryInitializedIterator struct {
	Event *IdentityRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryInitialized)
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
		it.Event = new(IdentityRegistryInitialized)
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
func (it *IdentityRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryInitialized represents a Initialized event raised by the IdentityRegistry contract.
type IdentityRegistryInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_IdentityRegistry *IdentityRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*IdentityRegistryInitializedIterator, error) {

	logs, sub, err := _IdentityRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryInitializedIterator{contract: _IdentityRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_IdentityRegistry *IdentityRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *IdentityRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _IdentityRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryInitialized)
				if err := _IdentityRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_IdentityRegistry *IdentityRegistryFilterer) ParseInitialized(log types.Log) (*IdentityRegistryInitialized, error) {
	event := new(IdentityRegistryInitialized)
	if err := _IdentityRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the IdentityRegistry contract.
type IdentityRegistryOwnershipTransferredIterator struct {
	Event *IdentityRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryOwnershipTransferred)
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
		it.Event = new(IdentityRegistryOwnershipTransferred)
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
func (it *IdentityRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the IdentityRegistry contract.
type IdentityRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IdentityRegistry *IdentityRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*IdentityRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IdentityRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryOwnershipTransferredIterator{contract: _IdentityRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IdentityRegistry *IdentityRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IdentityRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IdentityRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryOwnershipTransferred)
				if err := _IdentityRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_IdentityRegistry *IdentityRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*IdentityRegistryOwnershipTransferred, error) {
	event := new(IdentityRegistryOwnershipTransferred)
	if err := _IdentityRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryTrustedIssuersRegistrySetIterator is returned from FilterTrustedIssuersRegistrySet and is used to iterate over the raw logs and unpacked data for TrustedIssuersRegistrySet events raised by the IdentityRegistry contract.
type IdentityRegistryTrustedIssuersRegistrySetIterator struct {
	Event *IdentityRegistryTrustedIssuersRegistrySet // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryTrustedIssuersRegistrySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryTrustedIssuersRegistrySet)
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
		it.Event = new(IdentityRegistryTrustedIssuersRegistrySet)
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
func (it *IdentityRegistryTrustedIssuersRegistrySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryTrustedIssuersRegistrySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryTrustedIssuersRegistrySet represents a TrustedIssuersRegistrySet event raised by the IdentityRegistry contract.
type IdentityRegistryTrustedIssuersRegistrySet struct {
	TrustedIssuersRegistry common.Address
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterTrustedIssuersRegistrySet is a free log retrieval operation binding the contract event 0x1b98cb79e6f73020175fe87333f1b91ad6a881519c0afe30340c2599b2b4bde0.
//
// Solidity: event TrustedIssuersRegistrySet(address indexed trustedIssuersRegistry)
func (_IdentityRegistry *IdentityRegistryFilterer) FilterTrustedIssuersRegistrySet(opts *bind.FilterOpts, trustedIssuersRegistry []common.Address) (*IdentityRegistryTrustedIssuersRegistrySetIterator, error) {

	var trustedIssuersRegistryRule []interface{}
	for _, trustedIssuersRegistryItem := range trustedIssuersRegistry {
		trustedIssuersRegistryRule = append(trustedIssuersRegistryRule, trustedIssuersRegistryItem)
	}

	logs, sub, err := _IdentityRegistry.contract.FilterLogs(opts, "TrustedIssuersRegistrySet", trustedIssuersRegistryRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryTrustedIssuersRegistrySetIterator{contract: _IdentityRegistry.contract, event: "TrustedIssuersRegistrySet", logs: logs, sub: sub}, nil
}

// WatchTrustedIssuersRegistrySet is a free log subscription operation binding the contract event 0x1b98cb79e6f73020175fe87333f1b91ad6a881519c0afe30340c2599b2b4bde0.
//
// Solidity: event TrustedIssuersRegistrySet(address indexed trustedIssuersRegistry)
func (_IdentityRegistry *IdentityRegistryFilterer) WatchTrustedIssuersRegistrySet(opts *bind.WatchOpts, sink chan<- *IdentityRegistryTrustedIssuersRegistrySet, trustedIssuersRegistry []common.Address) (event.Subscription, error) {

	var trustedIssuersRegistryRule []interface{}
	for _, trustedIssuersRegistryItem := range trustedIssuersRegistry {
		trustedIssuersRegistryRule = append(trustedIssuersRegistryRule, trustedIssuersRegistryItem)
	}

	logs, sub, err := _IdentityRegistry.contract.WatchLogs(opts, "TrustedIssuersRegistrySet", trustedIssuersRegistryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryTrustedIssuersRegistrySet)
				if err := _IdentityRegistry.contract.UnpackLog(event, "TrustedIssuersRegistrySet", log); err != nil {
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

// ParseTrustedIssuersRegistrySet is a log parse operation binding the contract event 0x1b98cb79e6f73020175fe87333f1b91ad6a881519c0afe30340c2599b2b4bde0.
//
// Solidity: event TrustedIssuersRegistrySet(address indexed trustedIssuersRegistry)
func (_IdentityRegistry *IdentityRegistryFilterer) ParseTrustedIssuersRegistrySet(log types.Log) (*IdentityRegistryTrustedIssuersRegistrySet, error) {
	event := new(IdentityRegistryTrustedIssuersRegistrySet)
	if err := _IdentityRegistry.contract.UnpackLog(event, "TrustedIssuersRegistrySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
