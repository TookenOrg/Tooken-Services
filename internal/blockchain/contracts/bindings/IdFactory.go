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

// IdFactoryMetaData contains all meta data concerning the IdFactory contract.
var IdFactoryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementationAuthority\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"Deployed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"factory\",\"type\":\"address\"}],\"name\":\"TokenFactoryAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"factory\",\"type\":\"address\"}],\"name\":\"TokenFactoryRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"identity\",\"type\":\"address\"}],\"name\":\"TokenLinked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"identity\",\"type\":\"address\"}],\"name\":\"WalletLinked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"identity\",\"type\":\"address\"}],\"name\":\"WalletUnlinked\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factory\",\"type\":\"address\"}],\"name\":\"addTokenFactory\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wallet\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_salt\",\"type\":\"string\"}],\"name\":\"createIdentity\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wallet\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_salt\",\"type\":\"string\"},{\"internalType\":\"bytes32[]\",\"name\":\"_managementKeys\",\"type\":\"bytes32[]\"}],\"name\":\"createIdentityWithManagementKeys\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_tokenOwner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_salt\",\"type\":\"string\"}],\"name\":\"createTokenIdentity\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wallet\",\"type\":\"address\"}],\"name\":\"getIdentity\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_identity\",\"type\":\"address\"}],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_identity\",\"type\":\"address\"}],\"name\":\"getWallets\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"implementationAuthority\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_salt\",\"type\":\"string\"}],\"name\":\"isSaltTaken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factory\",\"type\":\"address\"}],\"name\":\"isTokenFactory\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newWallet\",\"type\":\"address\"}],\"name\":\"linkWallet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factory\",\"type\":\"address\"}],\"name\":\"removeTokenFactory\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oldWallet\",\"type\":\"address\"}],\"name\":\"unlinkWallet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620026b2380380620026b28339810160408190526200003491620000fc565b6200003f33620000ac565b6001600160a01b0381166200009a5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f206164647265737300604482015260640160405180910390fd5b6001600160a01b03166080526200012e565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b6000602082840312156200010f57600080fd5b81516001600160a01b03811681146200012757600080fd5b9392505050565b6080516125536200015f6000396000818161011c0152818161068701528181610d23015261163601526125536000f3fe60806040523480156200001157600080fd5b5060043610620001155760003560e01c8063715018a611620000a35780639ce19365116200006e5780639ce193651462000296578063b8bb812614620002ad578063f2fde38b14620002c4578063fe5cd59a14620002db57600080fd5b8063715018a6146200024c5780638da5cb5b14620002565780638e952bfe1462000268578063937529ef146200027f57600080fd5b80633e3bc3d711620000e45780633e3bc3d714620001af578063422c29a414620001de5780635027dbe2146200020457806359770438146200021d57600080fd5b80632307f882146200011a5780632fea7b8114620001595780633a50045114620001705780633d56ff661462000198575b600080fd5b7f00000000000000000000000000000000000000000000000000000000000000005b6040516001600160a01b0390911681526020015b60405180910390f35b6200013c6200016a36600462001b38565b620002f2565b620001876200018136600462001b5d565b62000358565b604051901515815260200162000150565b6200013c620001a936600462001c95565b6200038a565b62000187620001c036600462001b38565b6001600160a01b031660009081526001602052604090205460ff1690565b620001f5620001ef36600462001b38565b62000755565b60405162000150919062001cfb565b6200021b6200021536600462001b38565b620007cd565b005b6200013c6200022e36600462001b38565b6001600160a01b039081166000908152600660205260409020541690565b6200021b62000aef565b6000546001600160a01b03166200013c565b6200013c6200027936600462001d4a565b62000b07565b6200021b6200029036600462001b38565b62000e01565b6200021b620002a736600462001b38565b62000f16565b6200021b620002be36600462001b38565b62001032565b6200021b620002d536600462001b38565b6200130b565b6200013c620002ec36600462001d9e565b620013a1565b6001600160a01b03818116600090815260056020526040812054909116156200033457506001600160a01b039081166000908152600560205260409020541690565b506001600160a01b039081166000908152600360205260409020541690565b919050565b6000600283836040516200036e92919062001e8b565b9081526040519081900360200190205460ff1690505b92915050565b3360009081526001602052604081205460ff1680620003b357506000546001600160a01b031633145b620004055760405162461bcd60e51b815260206004820152601e60248201527f6f6e6c7920466163746f7279206f72206f776e65722063616e2063616c6c000060448201526064015b60405180910390fd5b6001600160a01b0384166200045d5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f2061646472657373006044820152606401620003fc565b6001600160a01b038316620004b55760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f2061646472657373006044820152606401620003fc565b604051602001620004d190602080825260009082015260400190565b6040516020818303038152906040528051906020012082604051602001620004fa919062001ec1565b60405160208183030381529060405280519060200120036200055f5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d20656d70747920737472696e67006044820152606401620003fc565b60008260405160200162000574919062001ef6565b604051602081830303815290604052905060028160405162000597919062001f3d565b9081526040519081900360200190205460ff1615620005ee5760405162461bcd60e51b815260206004820152601260248201527139b0b63a1030b63932b0b23c903a30b5b2b760711b6044820152606401620003fc565b6001600160a01b0385811660009081526005602052604090205416156200067e5760405162461bcd60e51b815260206004820152602360248201527f746f6b656e20616c7265616479206c696e6b656420746f20616e206964656e7460448201527f69747900000000000000000000000000000000000000000000000000000000006064820152608401620003fc565b6000620006ad827f0000000000000000000000000000000000000000000000000000000000000000876200194e565b90506001600283604051620006c3919062001f3d565b9081526040805160209281900383018120805460ff1916941515949094179093556001600160a01b0389811660008181526005855283812080546001600160a01b0319908116948916948517909155838252600690955292832080549094168117909355927fa8261d398ddc63db24cc53cd0c63c9464cabad1bc478ede2107b32c1c4010b7a9190a395945050505050565b6001600160a01b038116600090815260046020908152604091829020805483518184028101840190945280845260609392830182828015620007c157602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311620007a2575b50505050509050919050565b6001600160a01b038116620008255760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f2061646472657373006044820152606401620003fc565b336001600160a01b03821603620008a55760405162461bcd60e51b815260206004820152602260248201527f63616e6e6f742062652063616c6c6564206f6e2073656e64657220616464726560448201527f73730000000000000000000000000000000000000000000000000000000000006064820152608401620003fc565b6001600160a01b03808216600090815260036020526040808220543383529120548216911614620009195760405162461bcd60e51b815260206004820152601f60248201527f6f6e6c792061206c696e6b65642077616c6c65742063616e20756e6c696e6b006044820152606401620003fc565b6001600160a01b03808216600090815260036020908152604080832080546001600160a01b03198116909155909316808352600490915291812054905b8181101562000aa9576001600160a01b0383811660009081526004602052604090208054918616918390811062000991576200099162001f5b565b6000918252602090912001546001600160a01b03160362000a94576001600160a01b0383166000908152600460205260409020620009d160018462001f87565b81548110620009e457620009e462001f5b565b60009182526020808320909101546001600160a01b0386811684526004909252604090922080549190921691908390811062000a245762000a2462001f5b565b600091825260208083209190910180546001600160a01b0319166001600160a01b03948516179055918516815260049091526040902080548062000a6c5762000a6c62001f9d565b600082815260209020810160001990810180546001600160a01b031916905501905562000aa9565b8062000aa08162001fb3565b91505062000956565b50816001600160a01b0316836001600160a01b03167f35e6fc363a4bf723d53b26c1a751674aca9c3ead425f0591f84f5540ede86f1260405160405180910390a3505050565b62000af9620019de565b62000b05600062001a3a565b565b600062000b13620019de565b6001600160a01b03831662000b6b5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f2061646472657373006044820152606401620003fc565b60405160200162000b8790602080825260009082015260400190565b604051602081830303815290604052805190602001208260405160200162000bb0919062001ec1565b604051602081830303815290604052805190602001200362000c155760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d20656d70747920737472696e67006044820152606401620003fc565b60008260405160200162000c2a919062001fcf565b604051602081830303815290604052905060028160405162000c4d919062001f3d565b9081526040519081900360200190205460ff161562000ca45760405162461bcd60e51b815260206004820152601260248201527139b0b63a1030b63932b0b23c903a30b5b2b760711b6044820152606401620003fc565b6001600160a01b03848116600090815260036020526040902054161562000d1a5760405162461bcd60e51b8152602060048201526024808201527f77616c6c657420616c7265616479206c696e6b656420746f20616e206964656e6044820152637469747960e01b6064820152608401620003fc565b600062000d49827f0000000000000000000000000000000000000000000000000000000000000000876200194e565b9050600160028360405162000d5f919062001f3d565b9081526040805160209281900383018120805460ff1916941515949094179093556001600160a01b0388811660008181526003855283812080546001600160a01b0319908116948916948517909155838252600486529381208054600181018255908252948120909401805490931681179092559290917f8e0c709111388f5480579514d86663489ab1f206fe6da1a0c4d03ac8c318b3c691a3949350505050565b62000e0b620019de565b6001600160a01b03811662000e635760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f2061646472657373006044820152606401620003fc565b6001600160a01b03811660009081526001602052604090205460ff1662000ecd5760405162461bcd60e51b815260206004820152600d60248201527f6e6f74206120666163746f7279000000000000000000000000000000000000006044820152606401620003fc565b6001600160a01b038116600081815260016020526040808220805460ff19169055517fd1fd5274f793d20291c0abfe42e1ef63213a11b34996d485f7afb8fe014248519190a250565b62000f20620019de565b6001600160a01b03811662000f785760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f2061646472657373006044820152606401620003fc565b6001600160a01b03811660009081526001602052604090205460ff161562000fe35760405162461bcd60e51b815260206004820152601160248201527f616c7265616479206120666163746f72790000000000000000000000000000006044820152606401620003fc565b6001600160a01b0381166000818152600160208190526040808320805460ff1916909217909155517f45eb8ac5344d2d3f306550fe6e969ca4190526313c512afed851d052bf2ab2fd9190a250565b6001600160a01b0381166200108a5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f2061646472657373006044820152606401620003fc565b336000908152600360205260409020546001600160a01b0316620011175760405162461bcd60e51b815260206004820152602960248201527f77616c6c6574206e6f74206c696e6b656420746f20616e206964656e7469747960448201527f20636f6e747261637400000000000000000000000000000000000000000000006064820152608401620003fc565b6001600160a01b038181166000908152600360205260409020541615620011815760405162461bcd60e51b815260206004820152601960248201527f6e65772077616c6c657420616c7265616479206c696e6b6564000000000000006044820152606401620003fc565b6001600160a01b038181166000908152600560205260409020541615620011eb5760405162461bcd60e51b815260206004820181905260248201527f696e76616c696420617267756d656e74202d20746f6b656e20616464726573736044820152606401620003fc565b336000908152600360209081526040808320546001600160a01b03168084526004909252909120546065116200128a5760405162461bcd60e51b815260206004820152602560248201527f6d617820616d6f756e74206f662077616c6c657473207065722049442065786360448201527f65656465640000000000000000000000000000000000000000000000000000006064820152608401620003fc565b6001600160a01b03808316600081815260036020908152604080832080549587166001600160a01b031996871681179091558084526004835281842080546001810182559085529284209092018054909516841790945592517f8e0c709111388f5480579514d86663489ab1f206fe6da1a0c4d03ac8c318b3c69190a35050565b62001315620019de565b6001600160a01b038116620013935760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401620003fc565b6200139e8162001a3a565b50565b6000620013ad620019de565b6001600160a01b038416620014055760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f2061646472657373006044820152606401620003fc565b6040516020016200142190602080825260009082015260400190565b60405160208183030381529060405280519060200120836040516020016200144a919062001ec1565b6040516020818303038152906040528051906020012003620014af5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d20656d70747920737472696e67006044820152606401620003fc565b600083604051602001620014c4919062001fcf565b6040516020818303038152906040529050600281604051620014e7919062001f3d565b9081526040519081900360200190205460ff16156200153e5760405162461bcd60e51b815260206004820152601260248201527139b0b63a1030b63932b0b23c903a30b5b2b760711b6044820152606401620003fc565b6001600160a01b038581166000908152600360205260409020541615620015b45760405162461bcd60e51b8152602060048201526024808201527f77616c6c657420616c7265616479206c696e6b656420746f20616e206964656e6044820152637469747960e01b6064820152608401620003fc565b60008351116200162d5760405162461bcd60e51b815260206004820152602560248201527f696e76616c696420617267756d656e74202d20656d707479206c697374206f6660448201527f206b6579730000000000000000000000000000000000000000000000000000006064820152608401620003fc565b60006200165c827f0000000000000000000000000000000000000000000000000000000000000000306200194e565b905060005b8451811015620017e257604080516001600160a01b03891660208201520160405160208183030381529060405280519060200120858281518110620016aa57620016aa62001f5b565b602002602001015103620017275760405162461bcd60e51b815260206004820152603b60248201527f696e76616c696420617267756d656e74202d2077616c6c657420697320616c7360448201527f6f206c697374656420696e206d616e6167656d656e74206b65797300000000006064820152608401620003fc565b816001600160a01b0316631d3812408683815181106200174b576200174b62001f5b565b60200260200101516001806040518463ffffffff1660e01b815260040162001786939291909283526020830191909152604082015260600190565b6020604051808303816000875af1158015620017a6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620017cc919062002016565b5080620017d98162001fb3565b91505062001661565b50604080513060208201526001600160a01b038316916353d413c5910160408051601f198184030181529082905280516020909101207fffffffff0000000000000000000000000000000000000000000000000000000060e084901b1682526004820152600160248201526044016020604051808303816000875af115801562001870573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001896919062002016565b506001600283604051620018ab919062001f3d565b9081526040805160209281900383018120805460ff1916941515949094179093556001600160a01b0389811660008181526003855283812080546001600160a01b0319908116948916948517909155838252600486529381208054600181018255908252948120909401805490931681179092559290917f8e0c709111388f5480579514d86663489ab1f206fe6da1a0c4d03ac8c318b3c691a395945050505050565b60008060405180602001620019639062001b12565b601f1982820381018352601f9091011660408181526001600160a01b03878116602084015286168183015280518083038201815260608301909152919250600090620019b690849084906080016200203a565b6040516020818303038152906040529050620019d3878262001a8a565b979650505050505050565b6000546001600160a01b0316331462000b055760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401620003fc565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b6000808360405160200162001aa0919062001f3d565b60405160208183030381529060405280519060200120905060008360200184518381836000f592505050803b62001ad657600080fd5b6040516001600160a01b038216907ff40fcec21964ffb566044d083b4073f29f7f7929110ea19e1b3ebe375d89055e90600090a2949350505050565b6104b0806200206e83390190565b80356001600160a01b03811681146200035357600080fd5b60006020828403121562001b4b57600080fd5b62001b568262001b20565b9392505050565b6000806020838503121562001b7157600080fd5b823567ffffffffffffffff8082111562001b8a57600080fd5b818501915085601f83011262001b9f57600080fd5b81358181111562001baf57600080fd5b86602082850101111562001bc257600080fd5b60209290920196919550909350505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff8111828210171562001c165762001c1662001bd4565b604052919050565b600082601f83011262001c3057600080fd5b813567ffffffffffffffff81111562001c4d5762001c4d62001bd4565b62001c62601f8201601f191660200162001bea565b81815284602083860101111562001c7857600080fd5b816020850160208301376000918101602001919091529392505050565b60008060006060848603121562001cab57600080fd5b62001cb68462001b20565b925062001cc66020850162001b20565b9150604084013567ffffffffffffffff81111562001ce357600080fd5b62001cf18682870162001c1e565b9150509250925092565b6020808252825182820181905260009190848201906040850190845b8181101562001d3e5783516001600160a01b03168352928401929184019160010162001d17565b50909695505050505050565b6000806040838503121562001d5e57600080fd5b62001d698362001b20565b9150602083013567ffffffffffffffff81111562001d8657600080fd5b62001d948582860162001c1e565b9150509250929050565b60008060006060848603121562001db457600080fd5b62001dbf8462001b20565b925060208085013567ffffffffffffffff8082111562001dde57600080fd5b62001dec8883890162001c1e565b9450604087013591508082111562001e0357600080fd5b818701915087601f83011262001e1857600080fd5b81358181111562001e2d5762001e2d62001bd4565b8060051b915062001e4084830162001bea565b818152918301840191848101908a84111562001e5b57600080fd5b938501935b8385101562001e7b5784358252938501939085019062001e60565b8096505050505050509250925092565b8183823760009101908152919050565b60005b8381101562001eb857818101518382015260200162001e9e565b50506000910152565b602081526000825180602084015262001ee281604085016020870162001e9b565b601f01601f19169190910160400192915050565b7f546f6b656e00000000000000000000000000000000000000000000000000000081526000825162001f3081600585016020870162001e9b565b9190910160050192915050565b6000825162001f5181846020870162001e9b565b9190910192915050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b8181038181111562000384576200038462001f71565b634e487b7160e01b600052603160045260246000fd5b60006001820162001fc85762001fc862001f71565b5060010190565b7f4f494400000000000000000000000000000000000000000000000000000000008152600082516200200981600385016020870162001e9b565b9190910160030192915050565b6000602082840312156200202957600080fd5b8151801515811462001b5657600080fd5b600083516200204e81846020880162001e9b565b8351908301906200206481836020880162001e9b565b0194935050505056fe608060405234801561001057600080fd5b506040516104b03803806104b083398101604081905261002f91610271565b6001600160a01b03821661008a5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064015b60405180910390fd5b6001600160a01b0381166100e05760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f2061646472657373006044820152606401610081565b817f821f3e4d3d679f19eacc940c87acf846ea6eae24a63058ea750304437a62aafc556000826001600160a01b031663aaf10f426040518163ffffffff1660e01b8152600401602060405180830381865afa158015610143573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061016791906102a4565b6040516001600160a01b03848116602483015291925060009183169060440160408051601f198184030181529181526020820180516001600160e01b031663189acdbd60e31b179052516101bb91906102c6565b600060405180830381855af49150503d80600081146101f6576040519150601f19603f3d011682016040523d82523d6000602084013e6101fb565b606091505b505090508061024c5760405162461bcd60e51b815260206004820152601660248201527f496e697469616c697a6174696f6e206661696c65642e000000000000000000006044820152606401610081565b505050506102f5565b80516001600160a01b038116811461026c57600080fd5b919050565b6000806040838503121561028457600080fd5b61028d83610255565b915061029b60208401610255565b90509250929050565b6000602082840312156102b657600080fd5b6102bf82610255565b9392505050565b6000825160005b818110156102e757602081860181015185830152016102cd565b506000920191825250919050565b6101ac806103046000396000f3fe60806040526004361061001e5760003560e01c80632307f882146100e1575b60006100487f821f3e4d3d679f19eacc940c87acf846ea6eae24a63058ea750304437a62aafc5490565b73ffffffffffffffffffffffffffffffffffffffff1663aaf10f426040518163ffffffff1660e01b8152600401602060405180830381865afa158015610092573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100b69190610139565b90503660008037600080366000846127105a03f43d806000803e8180156100dc57816000f35b816000fd5b3480156100ed57600080fd5b507f821f3e4d3d679f19eacc940c87acf846ea6eae24a63058ea750304437a62aafc5460405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b60006020828403121561014b57600080fd5b815173ffffffffffffffffffffffffffffffffffffffff8116811461016f57600080fd5b939250505056fea2646970667358221220dc10dba4dcb99f75cb91819000d74353c98e3dd7471a2af2095fedc6a70516b664736f6c63430008110033a26469706673582212202de3b6b1c54bff10f48451f6928e6622f2eae107fcb0b3678807da568266f01f64736f6c63430008110033",
}

// IdFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use IdFactoryMetaData.ABI instead.
var IdFactoryABI = IdFactoryMetaData.ABI

// IdFactoryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use IdFactoryMetaData.Bin instead.
var IdFactoryBin = IdFactoryMetaData.Bin

// DeployIdFactory deploys a new Ethereum contract, binding an instance of IdFactory to it.
func DeployIdFactory(auth *bind.TransactOpts, backend bind.ContractBackend, implementationAuthority common.Address) (common.Address, *types.Transaction, *IdFactory, error) {
	parsed, err := IdFactoryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(IdFactoryBin), backend, implementationAuthority)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &IdFactory{IdFactoryCaller: IdFactoryCaller{contract: contract}, IdFactoryTransactor: IdFactoryTransactor{contract: contract}, IdFactoryFilterer: IdFactoryFilterer{contract: contract}}, nil
}

// IdFactory is an auto generated Go binding around an Ethereum contract.
type IdFactory struct {
	IdFactoryCaller     // Read-only binding to the contract
	IdFactoryTransactor // Write-only binding to the contract
	IdFactoryFilterer   // Log filterer for contract events
}

// IdFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type IdFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IdFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IdFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IdFactorySession struct {
	Contract     *IdFactory        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IdFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IdFactoryCallerSession struct {
	Contract *IdFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// IdFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IdFactoryTransactorSession struct {
	Contract     *IdFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// IdFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type IdFactoryRaw struct {
	Contract *IdFactory // Generic contract binding to access the raw methods on
}

// IdFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IdFactoryCallerRaw struct {
	Contract *IdFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// IdFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IdFactoryTransactorRaw struct {
	Contract *IdFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIdFactory creates a new instance of IdFactory, bound to a specific deployed contract.
func NewIdFactory(address common.Address, backend bind.ContractBackend) (*IdFactory, error) {
	contract, err := bindIdFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IdFactory{IdFactoryCaller: IdFactoryCaller{contract: contract}, IdFactoryTransactor: IdFactoryTransactor{contract: contract}, IdFactoryFilterer: IdFactoryFilterer{contract: contract}}, nil
}

// NewIdFactoryCaller creates a new read-only instance of IdFactory, bound to a specific deployed contract.
func NewIdFactoryCaller(address common.Address, caller bind.ContractCaller) (*IdFactoryCaller, error) {
	contract, err := bindIdFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IdFactoryCaller{contract: contract}, nil
}

// NewIdFactoryTransactor creates a new write-only instance of IdFactory, bound to a specific deployed contract.
func NewIdFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*IdFactoryTransactor, error) {
	contract, err := bindIdFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IdFactoryTransactor{contract: contract}, nil
}

// NewIdFactoryFilterer creates a new log filterer instance of IdFactory, bound to a specific deployed contract.
func NewIdFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*IdFactoryFilterer, error) {
	contract, err := bindIdFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IdFactoryFilterer{contract: contract}, nil
}

// bindIdFactory binds a generic wrapper to an already deployed contract.
func bindIdFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IdFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdFactory *IdFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IdFactory.Contract.IdFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdFactory *IdFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdFactory.Contract.IdFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdFactory *IdFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdFactory.Contract.IdFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdFactory *IdFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IdFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdFactory *IdFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdFactory *IdFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdFactory.Contract.contract.Transact(opts, method, params...)
}

// GetIdentity is a free data retrieval call binding the contract method 0x2fea7b81.
//
// Solidity: function getIdentity(address _wallet) view returns(address)
func (_IdFactory *IdFactoryCaller) GetIdentity(opts *bind.CallOpts, _wallet common.Address) (common.Address, error) {
	var out []interface{}
	err := _IdFactory.contract.Call(opts, &out, "getIdentity", _wallet)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetIdentity is a free data retrieval call binding the contract method 0x2fea7b81.
//
// Solidity: function getIdentity(address _wallet) view returns(address)
func (_IdFactory *IdFactorySession) GetIdentity(_wallet common.Address) (common.Address, error) {
	return _IdFactory.Contract.GetIdentity(&_IdFactory.CallOpts, _wallet)
}

// GetIdentity is a free data retrieval call binding the contract method 0x2fea7b81.
//
// Solidity: function getIdentity(address _wallet) view returns(address)
func (_IdFactory *IdFactoryCallerSession) GetIdentity(_wallet common.Address) (common.Address, error) {
	return _IdFactory.Contract.GetIdentity(&_IdFactory.CallOpts, _wallet)
}

// GetToken is a free data retrieval call binding the contract method 0x59770438.
//
// Solidity: function getToken(address _identity) view returns(address)
func (_IdFactory *IdFactoryCaller) GetToken(opts *bind.CallOpts, _identity common.Address) (common.Address, error) {
	var out []interface{}
	err := _IdFactory.contract.Call(opts, &out, "getToken", _identity)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetToken is a free data retrieval call binding the contract method 0x59770438.
//
// Solidity: function getToken(address _identity) view returns(address)
func (_IdFactory *IdFactorySession) GetToken(_identity common.Address) (common.Address, error) {
	return _IdFactory.Contract.GetToken(&_IdFactory.CallOpts, _identity)
}

// GetToken is a free data retrieval call binding the contract method 0x59770438.
//
// Solidity: function getToken(address _identity) view returns(address)
func (_IdFactory *IdFactoryCallerSession) GetToken(_identity common.Address) (common.Address, error) {
	return _IdFactory.Contract.GetToken(&_IdFactory.CallOpts, _identity)
}

// GetWallets is a free data retrieval call binding the contract method 0x422c29a4.
//
// Solidity: function getWallets(address _identity) view returns(address[])
func (_IdFactory *IdFactoryCaller) GetWallets(opts *bind.CallOpts, _identity common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _IdFactory.contract.Call(opts, &out, "getWallets", _identity)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetWallets is a free data retrieval call binding the contract method 0x422c29a4.
//
// Solidity: function getWallets(address _identity) view returns(address[])
func (_IdFactory *IdFactorySession) GetWallets(_identity common.Address) ([]common.Address, error) {
	return _IdFactory.Contract.GetWallets(&_IdFactory.CallOpts, _identity)
}

// GetWallets is a free data retrieval call binding the contract method 0x422c29a4.
//
// Solidity: function getWallets(address _identity) view returns(address[])
func (_IdFactory *IdFactoryCallerSession) GetWallets(_identity common.Address) ([]common.Address, error) {
	return _IdFactory.Contract.GetWallets(&_IdFactory.CallOpts, _identity)
}

// ImplementationAuthority is a free data retrieval call binding the contract method 0x2307f882.
//
// Solidity: function implementationAuthority() view returns(address)
func (_IdFactory *IdFactoryCaller) ImplementationAuthority(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IdFactory.contract.Call(opts, &out, "implementationAuthority")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ImplementationAuthority is a free data retrieval call binding the contract method 0x2307f882.
//
// Solidity: function implementationAuthority() view returns(address)
func (_IdFactory *IdFactorySession) ImplementationAuthority() (common.Address, error) {
	return _IdFactory.Contract.ImplementationAuthority(&_IdFactory.CallOpts)
}

// ImplementationAuthority is a free data retrieval call binding the contract method 0x2307f882.
//
// Solidity: function implementationAuthority() view returns(address)
func (_IdFactory *IdFactoryCallerSession) ImplementationAuthority() (common.Address, error) {
	return _IdFactory.Contract.ImplementationAuthority(&_IdFactory.CallOpts)
}

// IsSaltTaken is a free data retrieval call binding the contract method 0x3a500451.
//
// Solidity: function isSaltTaken(string _salt) view returns(bool)
func (_IdFactory *IdFactoryCaller) IsSaltTaken(opts *bind.CallOpts, _salt string) (bool, error) {
	var out []interface{}
	err := _IdFactory.contract.Call(opts, &out, "isSaltTaken", _salt)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsSaltTaken is a free data retrieval call binding the contract method 0x3a500451.
//
// Solidity: function isSaltTaken(string _salt) view returns(bool)
func (_IdFactory *IdFactorySession) IsSaltTaken(_salt string) (bool, error) {
	return _IdFactory.Contract.IsSaltTaken(&_IdFactory.CallOpts, _salt)
}

// IsSaltTaken is a free data retrieval call binding the contract method 0x3a500451.
//
// Solidity: function isSaltTaken(string _salt) view returns(bool)
func (_IdFactory *IdFactoryCallerSession) IsSaltTaken(_salt string) (bool, error) {
	return _IdFactory.Contract.IsSaltTaken(&_IdFactory.CallOpts, _salt)
}

// IsTokenFactory is a free data retrieval call binding the contract method 0x3e3bc3d7.
//
// Solidity: function isTokenFactory(address _factory) view returns(bool)
func (_IdFactory *IdFactoryCaller) IsTokenFactory(opts *bind.CallOpts, _factory common.Address) (bool, error) {
	var out []interface{}
	err := _IdFactory.contract.Call(opts, &out, "isTokenFactory", _factory)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsTokenFactory is a free data retrieval call binding the contract method 0x3e3bc3d7.
//
// Solidity: function isTokenFactory(address _factory) view returns(bool)
func (_IdFactory *IdFactorySession) IsTokenFactory(_factory common.Address) (bool, error) {
	return _IdFactory.Contract.IsTokenFactory(&_IdFactory.CallOpts, _factory)
}

// IsTokenFactory is a free data retrieval call binding the contract method 0x3e3bc3d7.
//
// Solidity: function isTokenFactory(address _factory) view returns(bool)
func (_IdFactory *IdFactoryCallerSession) IsTokenFactory(_factory common.Address) (bool, error) {
	return _IdFactory.Contract.IsTokenFactory(&_IdFactory.CallOpts, _factory)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IdFactory *IdFactoryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IdFactory.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IdFactory *IdFactorySession) Owner() (common.Address, error) {
	return _IdFactory.Contract.Owner(&_IdFactory.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IdFactory *IdFactoryCallerSession) Owner() (common.Address, error) {
	return _IdFactory.Contract.Owner(&_IdFactory.CallOpts)
}

// AddTokenFactory is a paid mutator transaction binding the contract method 0x9ce19365.
//
// Solidity: function addTokenFactory(address _factory) returns()
func (_IdFactory *IdFactoryTransactor) AddTokenFactory(opts *bind.TransactOpts, _factory common.Address) (*types.Transaction, error) {
	return _IdFactory.contract.Transact(opts, "addTokenFactory", _factory)
}

// AddTokenFactory is a paid mutator transaction binding the contract method 0x9ce19365.
//
// Solidity: function addTokenFactory(address _factory) returns()
func (_IdFactory *IdFactorySession) AddTokenFactory(_factory common.Address) (*types.Transaction, error) {
	return _IdFactory.Contract.AddTokenFactory(&_IdFactory.TransactOpts, _factory)
}

// AddTokenFactory is a paid mutator transaction binding the contract method 0x9ce19365.
//
// Solidity: function addTokenFactory(address _factory) returns()
func (_IdFactory *IdFactoryTransactorSession) AddTokenFactory(_factory common.Address) (*types.Transaction, error) {
	return _IdFactory.Contract.AddTokenFactory(&_IdFactory.TransactOpts, _factory)
}

// CreateIdentity is a paid mutator transaction binding the contract method 0x8e952bfe.
//
// Solidity: function createIdentity(address _wallet, string _salt) returns(address)
func (_IdFactory *IdFactoryTransactor) CreateIdentity(opts *bind.TransactOpts, _wallet common.Address, _salt string) (*types.Transaction, error) {
	return _IdFactory.contract.Transact(opts, "createIdentity", _wallet, _salt)
}

// CreateIdentity is a paid mutator transaction binding the contract method 0x8e952bfe.
//
// Solidity: function createIdentity(address _wallet, string _salt) returns(address)
func (_IdFactory *IdFactorySession) CreateIdentity(_wallet common.Address, _salt string) (*types.Transaction, error) {
	return _IdFactory.Contract.CreateIdentity(&_IdFactory.TransactOpts, _wallet, _salt)
}

// CreateIdentity is a paid mutator transaction binding the contract method 0x8e952bfe.
//
// Solidity: function createIdentity(address _wallet, string _salt) returns(address)
func (_IdFactory *IdFactoryTransactorSession) CreateIdentity(_wallet common.Address, _salt string) (*types.Transaction, error) {
	return _IdFactory.Contract.CreateIdentity(&_IdFactory.TransactOpts, _wallet, _salt)
}

// CreateIdentityWithManagementKeys is a paid mutator transaction binding the contract method 0xfe5cd59a.
//
// Solidity: function createIdentityWithManagementKeys(address _wallet, string _salt, bytes32[] _managementKeys) returns(address)
func (_IdFactory *IdFactoryTransactor) CreateIdentityWithManagementKeys(opts *bind.TransactOpts, _wallet common.Address, _salt string, _managementKeys [][32]byte) (*types.Transaction, error) {
	return _IdFactory.contract.Transact(opts, "createIdentityWithManagementKeys", _wallet, _salt, _managementKeys)
}

// CreateIdentityWithManagementKeys is a paid mutator transaction binding the contract method 0xfe5cd59a.
//
// Solidity: function createIdentityWithManagementKeys(address _wallet, string _salt, bytes32[] _managementKeys) returns(address)
func (_IdFactory *IdFactorySession) CreateIdentityWithManagementKeys(_wallet common.Address, _salt string, _managementKeys [][32]byte) (*types.Transaction, error) {
	return _IdFactory.Contract.CreateIdentityWithManagementKeys(&_IdFactory.TransactOpts, _wallet, _salt, _managementKeys)
}

// CreateIdentityWithManagementKeys is a paid mutator transaction binding the contract method 0xfe5cd59a.
//
// Solidity: function createIdentityWithManagementKeys(address _wallet, string _salt, bytes32[] _managementKeys) returns(address)
func (_IdFactory *IdFactoryTransactorSession) CreateIdentityWithManagementKeys(_wallet common.Address, _salt string, _managementKeys [][32]byte) (*types.Transaction, error) {
	return _IdFactory.Contract.CreateIdentityWithManagementKeys(&_IdFactory.TransactOpts, _wallet, _salt, _managementKeys)
}

// CreateTokenIdentity is a paid mutator transaction binding the contract method 0x3d56ff66.
//
// Solidity: function createTokenIdentity(address _token, address _tokenOwner, string _salt) returns(address)
func (_IdFactory *IdFactoryTransactor) CreateTokenIdentity(opts *bind.TransactOpts, _token common.Address, _tokenOwner common.Address, _salt string) (*types.Transaction, error) {
	return _IdFactory.contract.Transact(opts, "createTokenIdentity", _token, _tokenOwner, _salt)
}

// CreateTokenIdentity is a paid mutator transaction binding the contract method 0x3d56ff66.
//
// Solidity: function createTokenIdentity(address _token, address _tokenOwner, string _salt) returns(address)
func (_IdFactory *IdFactorySession) CreateTokenIdentity(_token common.Address, _tokenOwner common.Address, _salt string) (*types.Transaction, error) {
	return _IdFactory.Contract.CreateTokenIdentity(&_IdFactory.TransactOpts, _token, _tokenOwner, _salt)
}

// CreateTokenIdentity is a paid mutator transaction binding the contract method 0x3d56ff66.
//
// Solidity: function createTokenIdentity(address _token, address _tokenOwner, string _salt) returns(address)
func (_IdFactory *IdFactoryTransactorSession) CreateTokenIdentity(_token common.Address, _tokenOwner common.Address, _salt string) (*types.Transaction, error) {
	return _IdFactory.Contract.CreateTokenIdentity(&_IdFactory.TransactOpts, _token, _tokenOwner, _salt)
}

// LinkWallet is a paid mutator transaction binding the contract method 0xb8bb8126.
//
// Solidity: function linkWallet(address _newWallet) returns()
func (_IdFactory *IdFactoryTransactor) LinkWallet(opts *bind.TransactOpts, _newWallet common.Address) (*types.Transaction, error) {
	return _IdFactory.contract.Transact(opts, "linkWallet", _newWallet)
}

// LinkWallet is a paid mutator transaction binding the contract method 0xb8bb8126.
//
// Solidity: function linkWallet(address _newWallet) returns()
func (_IdFactory *IdFactorySession) LinkWallet(_newWallet common.Address) (*types.Transaction, error) {
	return _IdFactory.Contract.LinkWallet(&_IdFactory.TransactOpts, _newWallet)
}

// LinkWallet is a paid mutator transaction binding the contract method 0xb8bb8126.
//
// Solidity: function linkWallet(address _newWallet) returns()
func (_IdFactory *IdFactoryTransactorSession) LinkWallet(_newWallet common.Address) (*types.Transaction, error) {
	return _IdFactory.Contract.LinkWallet(&_IdFactory.TransactOpts, _newWallet)
}

// RemoveTokenFactory is a paid mutator transaction binding the contract method 0x937529ef.
//
// Solidity: function removeTokenFactory(address _factory) returns()
func (_IdFactory *IdFactoryTransactor) RemoveTokenFactory(opts *bind.TransactOpts, _factory common.Address) (*types.Transaction, error) {
	return _IdFactory.contract.Transact(opts, "removeTokenFactory", _factory)
}

// RemoveTokenFactory is a paid mutator transaction binding the contract method 0x937529ef.
//
// Solidity: function removeTokenFactory(address _factory) returns()
func (_IdFactory *IdFactorySession) RemoveTokenFactory(_factory common.Address) (*types.Transaction, error) {
	return _IdFactory.Contract.RemoveTokenFactory(&_IdFactory.TransactOpts, _factory)
}

// RemoveTokenFactory is a paid mutator transaction binding the contract method 0x937529ef.
//
// Solidity: function removeTokenFactory(address _factory) returns()
func (_IdFactory *IdFactoryTransactorSession) RemoveTokenFactory(_factory common.Address) (*types.Transaction, error) {
	return _IdFactory.Contract.RemoveTokenFactory(&_IdFactory.TransactOpts, _factory)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IdFactory *IdFactoryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdFactory.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IdFactory *IdFactorySession) RenounceOwnership() (*types.Transaction, error) {
	return _IdFactory.Contract.RenounceOwnership(&_IdFactory.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IdFactory *IdFactoryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _IdFactory.Contract.RenounceOwnership(&_IdFactory.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IdFactory *IdFactoryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _IdFactory.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IdFactory *IdFactorySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IdFactory.Contract.TransferOwnership(&_IdFactory.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IdFactory *IdFactoryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IdFactory.Contract.TransferOwnership(&_IdFactory.TransactOpts, newOwner)
}

// UnlinkWallet is a paid mutator transaction binding the contract method 0x5027dbe2.
//
// Solidity: function unlinkWallet(address _oldWallet) returns()
func (_IdFactory *IdFactoryTransactor) UnlinkWallet(opts *bind.TransactOpts, _oldWallet common.Address) (*types.Transaction, error) {
	return _IdFactory.contract.Transact(opts, "unlinkWallet", _oldWallet)
}

// UnlinkWallet is a paid mutator transaction binding the contract method 0x5027dbe2.
//
// Solidity: function unlinkWallet(address _oldWallet) returns()
func (_IdFactory *IdFactorySession) UnlinkWallet(_oldWallet common.Address) (*types.Transaction, error) {
	return _IdFactory.Contract.UnlinkWallet(&_IdFactory.TransactOpts, _oldWallet)
}

// UnlinkWallet is a paid mutator transaction binding the contract method 0x5027dbe2.
//
// Solidity: function unlinkWallet(address _oldWallet) returns()
func (_IdFactory *IdFactoryTransactorSession) UnlinkWallet(_oldWallet common.Address) (*types.Transaction, error) {
	return _IdFactory.Contract.UnlinkWallet(&_IdFactory.TransactOpts, _oldWallet)
}

// IdFactoryDeployedIterator is returned from FilterDeployed and is used to iterate over the raw logs and unpacked data for Deployed events raised by the IdFactory contract.
type IdFactoryDeployedIterator struct {
	Event *IdFactoryDeployed // Event containing the contract specifics and raw log

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
func (it *IdFactoryDeployedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdFactoryDeployed)
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
		it.Event = new(IdFactoryDeployed)
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
func (it *IdFactoryDeployedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdFactoryDeployedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdFactoryDeployed represents a Deployed event raised by the IdFactory contract.
type IdFactoryDeployed struct {
	Addr common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterDeployed is a free log retrieval operation binding the contract event 0xf40fcec21964ffb566044d083b4073f29f7f7929110ea19e1b3ebe375d89055e.
//
// Solidity: event Deployed(address indexed _addr)
func (_IdFactory *IdFactoryFilterer) FilterDeployed(opts *bind.FilterOpts, _addr []common.Address) (*IdFactoryDeployedIterator, error) {

	var _addrRule []interface{}
	for _, _addrItem := range _addr {
		_addrRule = append(_addrRule, _addrItem)
	}

	logs, sub, err := _IdFactory.contract.FilterLogs(opts, "Deployed", _addrRule)
	if err != nil {
		return nil, err
	}
	return &IdFactoryDeployedIterator{contract: _IdFactory.contract, event: "Deployed", logs: logs, sub: sub}, nil
}

// WatchDeployed is a free log subscription operation binding the contract event 0xf40fcec21964ffb566044d083b4073f29f7f7929110ea19e1b3ebe375d89055e.
//
// Solidity: event Deployed(address indexed _addr)
func (_IdFactory *IdFactoryFilterer) WatchDeployed(opts *bind.WatchOpts, sink chan<- *IdFactoryDeployed, _addr []common.Address) (event.Subscription, error) {

	var _addrRule []interface{}
	for _, _addrItem := range _addr {
		_addrRule = append(_addrRule, _addrItem)
	}

	logs, sub, err := _IdFactory.contract.WatchLogs(opts, "Deployed", _addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdFactoryDeployed)
				if err := _IdFactory.contract.UnpackLog(event, "Deployed", log); err != nil {
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

// ParseDeployed is a log parse operation binding the contract event 0xf40fcec21964ffb566044d083b4073f29f7f7929110ea19e1b3ebe375d89055e.
//
// Solidity: event Deployed(address indexed _addr)
func (_IdFactory *IdFactoryFilterer) ParseDeployed(log types.Log) (*IdFactoryDeployed, error) {
	event := new(IdFactoryDeployed)
	if err := _IdFactory.contract.UnpackLog(event, "Deployed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdFactoryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the IdFactory contract.
type IdFactoryOwnershipTransferredIterator struct {
	Event *IdFactoryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *IdFactoryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdFactoryOwnershipTransferred)
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
		it.Event = new(IdFactoryOwnershipTransferred)
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
func (it *IdFactoryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdFactoryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdFactoryOwnershipTransferred represents a OwnershipTransferred event raised by the IdFactory contract.
type IdFactoryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IdFactory *IdFactoryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*IdFactoryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IdFactory.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &IdFactoryOwnershipTransferredIterator{contract: _IdFactory.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IdFactory *IdFactoryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IdFactoryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IdFactory.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdFactoryOwnershipTransferred)
				if err := _IdFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_IdFactory *IdFactoryFilterer) ParseOwnershipTransferred(log types.Log) (*IdFactoryOwnershipTransferred, error) {
	event := new(IdFactoryOwnershipTransferred)
	if err := _IdFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdFactoryTokenFactoryAddedIterator is returned from FilterTokenFactoryAdded and is used to iterate over the raw logs and unpacked data for TokenFactoryAdded events raised by the IdFactory contract.
type IdFactoryTokenFactoryAddedIterator struct {
	Event *IdFactoryTokenFactoryAdded // Event containing the contract specifics and raw log

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
func (it *IdFactoryTokenFactoryAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdFactoryTokenFactoryAdded)
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
		it.Event = new(IdFactoryTokenFactoryAdded)
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
func (it *IdFactoryTokenFactoryAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdFactoryTokenFactoryAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdFactoryTokenFactoryAdded represents a TokenFactoryAdded event raised by the IdFactory contract.
type IdFactoryTokenFactoryAdded struct {
	Factory common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTokenFactoryAdded is a free log retrieval operation binding the contract event 0x45eb8ac5344d2d3f306550fe6e969ca4190526313c512afed851d052bf2ab2fd.
//
// Solidity: event TokenFactoryAdded(address indexed factory)
func (_IdFactory *IdFactoryFilterer) FilterTokenFactoryAdded(opts *bind.FilterOpts, factory []common.Address) (*IdFactoryTokenFactoryAddedIterator, error) {

	var factoryRule []interface{}
	for _, factoryItem := range factory {
		factoryRule = append(factoryRule, factoryItem)
	}

	logs, sub, err := _IdFactory.contract.FilterLogs(opts, "TokenFactoryAdded", factoryRule)
	if err != nil {
		return nil, err
	}
	return &IdFactoryTokenFactoryAddedIterator{contract: _IdFactory.contract, event: "TokenFactoryAdded", logs: logs, sub: sub}, nil
}

// WatchTokenFactoryAdded is a free log subscription operation binding the contract event 0x45eb8ac5344d2d3f306550fe6e969ca4190526313c512afed851d052bf2ab2fd.
//
// Solidity: event TokenFactoryAdded(address indexed factory)
func (_IdFactory *IdFactoryFilterer) WatchTokenFactoryAdded(opts *bind.WatchOpts, sink chan<- *IdFactoryTokenFactoryAdded, factory []common.Address) (event.Subscription, error) {

	var factoryRule []interface{}
	for _, factoryItem := range factory {
		factoryRule = append(factoryRule, factoryItem)
	}

	logs, sub, err := _IdFactory.contract.WatchLogs(opts, "TokenFactoryAdded", factoryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdFactoryTokenFactoryAdded)
				if err := _IdFactory.contract.UnpackLog(event, "TokenFactoryAdded", log); err != nil {
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

// ParseTokenFactoryAdded is a log parse operation binding the contract event 0x45eb8ac5344d2d3f306550fe6e969ca4190526313c512afed851d052bf2ab2fd.
//
// Solidity: event TokenFactoryAdded(address indexed factory)
func (_IdFactory *IdFactoryFilterer) ParseTokenFactoryAdded(log types.Log) (*IdFactoryTokenFactoryAdded, error) {
	event := new(IdFactoryTokenFactoryAdded)
	if err := _IdFactory.contract.UnpackLog(event, "TokenFactoryAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdFactoryTokenFactoryRemovedIterator is returned from FilterTokenFactoryRemoved and is used to iterate over the raw logs and unpacked data for TokenFactoryRemoved events raised by the IdFactory contract.
type IdFactoryTokenFactoryRemovedIterator struct {
	Event *IdFactoryTokenFactoryRemoved // Event containing the contract specifics and raw log

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
func (it *IdFactoryTokenFactoryRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdFactoryTokenFactoryRemoved)
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
		it.Event = new(IdFactoryTokenFactoryRemoved)
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
func (it *IdFactoryTokenFactoryRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdFactoryTokenFactoryRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdFactoryTokenFactoryRemoved represents a TokenFactoryRemoved event raised by the IdFactory contract.
type IdFactoryTokenFactoryRemoved struct {
	Factory common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTokenFactoryRemoved is a free log retrieval operation binding the contract event 0xd1fd5274f793d20291c0abfe42e1ef63213a11b34996d485f7afb8fe01424851.
//
// Solidity: event TokenFactoryRemoved(address indexed factory)
func (_IdFactory *IdFactoryFilterer) FilterTokenFactoryRemoved(opts *bind.FilterOpts, factory []common.Address) (*IdFactoryTokenFactoryRemovedIterator, error) {

	var factoryRule []interface{}
	for _, factoryItem := range factory {
		factoryRule = append(factoryRule, factoryItem)
	}

	logs, sub, err := _IdFactory.contract.FilterLogs(opts, "TokenFactoryRemoved", factoryRule)
	if err != nil {
		return nil, err
	}
	return &IdFactoryTokenFactoryRemovedIterator{contract: _IdFactory.contract, event: "TokenFactoryRemoved", logs: logs, sub: sub}, nil
}

// WatchTokenFactoryRemoved is a free log subscription operation binding the contract event 0xd1fd5274f793d20291c0abfe42e1ef63213a11b34996d485f7afb8fe01424851.
//
// Solidity: event TokenFactoryRemoved(address indexed factory)
func (_IdFactory *IdFactoryFilterer) WatchTokenFactoryRemoved(opts *bind.WatchOpts, sink chan<- *IdFactoryTokenFactoryRemoved, factory []common.Address) (event.Subscription, error) {

	var factoryRule []interface{}
	for _, factoryItem := range factory {
		factoryRule = append(factoryRule, factoryItem)
	}

	logs, sub, err := _IdFactory.contract.WatchLogs(opts, "TokenFactoryRemoved", factoryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdFactoryTokenFactoryRemoved)
				if err := _IdFactory.contract.UnpackLog(event, "TokenFactoryRemoved", log); err != nil {
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

// ParseTokenFactoryRemoved is a log parse operation binding the contract event 0xd1fd5274f793d20291c0abfe42e1ef63213a11b34996d485f7afb8fe01424851.
//
// Solidity: event TokenFactoryRemoved(address indexed factory)
func (_IdFactory *IdFactoryFilterer) ParseTokenFactoryRemoved(log types.Log) (*IdFactoryTokenFactoryRemoved, error) {
	event := new(IdFactoryTokenFactoryRemoved)
	if err := _IdFactory.contract.UnpackLog(event, "TokenFactoryRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdFactoryTokenLinkedIterator is returned from FilterTokenLinked and is used to iterate over the raw logs and unpacked data for TokenLinked events raised by the IdFactory contract.
type IdFactoryTokenLinkedIterator struct {
	Event *IdFactoryTokenLinked // Event containing the contract specifics and raw log

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
func (it *IdFactoryTokenLinkedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdFactoryTokenLinked)
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
		it.Event = new(IdFactoryTokenLinked)
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
func (it *IdFactoryTokenLinkedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdFactoryTokenLinkedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdFactoryTokenLinked represents a TokenLinked event raised by the IdFactory contract.
type IdFactoryTokenLinked struct {
	Token    common.Address
	Identity common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTokenLinked is a free log retrieval operation binding the contract event 0xa8261d398ddc63db24cc53cd0c63c9464cabad1bc478ede2107b32c1c4010b7a.
//
// Solidity: event TokenLinked(address indexed token, address indexed identity)
func (_IdFactory *IdFactoryFilterer) FilterTokenLinked(opts *bind.FilterOpts, token []common.Address, identity []common.Address) (*IdFactoryTokenLinkedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _IdFactory.contract.FilterLogs(opts, "TokenLinked", tokenRule, identityRule)
	if err != nil {
		return nil, err
	}
	return &IdFactoryTokenLinkedIterator{contract: _IdFactory.contract, event: "TokenLinked", logs: logs, sub: sub}, nil
}

// WatchTokenLinked is a free log subscription operation binding the contract event 0xa8261d398ddc63db24cc53cd0c63c9464cabad1bc478ede2107b32c1c4010b7a.
//
// Solidity: event TokenLinked(address indexed token, address indexed identity)
func (_IdFactory *IdFactoryFilterer) WatchTokenLinked(opts *bind.WatchOpts, sink chan<- *IdFactoryTokenLinked, token []common.Address, identity []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _IdFactory.contract.WatchLogs(opts, "TokenLinked", tokenRule, identityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdFactoryTokenLinked)
				if err := _IdFactory.contract.UnpackLog(event, "TokenLinked", log); err != nil {
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

// ParseTokenLinked is a log parse operation binding the contract event 0xa8261d398ddc63db24cc53cd0c63c9464cabad1bc478ede2107b32c1c4010b7a.
//
// Solidity: event TokenLinked(address indexed token, address indexed identity)
func (_IdFactory *IdFactoryFilterer) ParseTokenLinked(log types.Log) (*IdFactoryTokenLinked, error) {
	event := new(IdFactoryTokenLinked)
	if err := _IdFactory.contract.UnpackLog(event, "TokenLinked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdFactoryWalletLinkedIterator is returned from FilterWalletLinked and is used to iterate over the raw logs and unpacked data for WalletLinked events raised by the IdFactory contract.
type IdFactoryWalletLinkedIterator struct {
	Event *IdFactoryWalletLinked // Event containing the contract specifics and raw log

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
func (it *IdFactoryWalletLinkedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdFactoryWalletLinked)
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
		it.Event = new(IdFactoryWalletLinked)
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
func (it *IdFactoryWalletLinkedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdFactoryWalletLinkedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdFactoryWalletLinked represents a WalletLinked event raised by the IdFactory contract.
type IdFactoryWalletLinked struct {
	Wallet   common.Address
	Identity common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterWalletLinked is a free log retrieval operation binding the contract event 0x8e0c709111388f5480579514d86663489ab1f206fe6da1a0c4d03ac8c318b3c6.
//
// Solidity: event WalletLinked(address indexed wallet, address indexed identity)
func (_IdFactory *IdFactoryFilterer) FilterWalletLinked(opts *bind.FilterOpts, wallet []common.Address, identity []common.Address) (*IdFactoryWalletLinkedIterator, error) {

	var walletRule []interface{}
	for _, walletItem := range wallet {
		walletRule = append(walletRule, walletItem)
	}
	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _IdFactory.contract.FilterLogs(opts, "WalletLinked", walletRule, identityRule)
	if err != nil {
		return nil, err
	}
	return &IdFactoryWalletLinkedIterator{contract: _IdFactory.contract, event: "WalletLinked", logs: logs, sub: sub}, nil
}

// WatchWalletLinked is a free log subscription operation binding the contract event 0x8e0c709111388f5480579514d86663489ab1f206fe6da1a0c4d03ac8c318b3c6.
//
// Solidity: event WalletLinked(address indexed wallet, address indexed identity)
func (_IdFactory *IdFactoryFilterer) WatchWalletLinked(opts *bind.WatchOpts, sink chan<- *IdFactoryWalletLinked, wallet []common.Address, identity []common.Address) (event.Subscription, error) {

	var walletRule []interface{}
	for _, walletItem := range wallet {
		walletRule = append(walletRule, walletItem)
	}
	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _IdFactory.contract.WatchLogs(opts, "WalletLinked", walletRule, identityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdFactoryWalletLinked)
				if err := _IdFactory.contract.UnpackLog(event, "WalletLinked", log); err != nil {
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

// ParseWalletLinked is a log parse operation binding the contract event 0x8e0c709111388f5480579514d86663489ab1f206fe6da1a0c4d03ac8c318b3c6.
//
// Solidity: event WalletLinked(address indexed wallet, address indexed identity)
func (_IdFactory *IdFactoryFilterer) ParseWalletLinked(log types.Log) (*IdFactoryWalletLinked, error) {
	event := new(IdFactoryWalletLinked)
	if err := _IdFactory.contract.UnpackLog(event, "WalletLinked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdFactoryWalletUnlinkedIterator is returned from FilterWalletUnlinked and is used to iterate over the raw logs and unpacked data for WalletUnlinked events raised by the IdFactory contract.
type IdFactoryWalletUnlinkedIterator struct {
	Event *IdFactoryWalletUnlinked // Event containing the contract specifics and raw log

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
func (it *IdFactoryWalletUnlinkedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdFactoryWalletUnlinked)
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
		it.Event = new(IdFactoryWalletUnlinked)
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
func (it *IdFactoryWalletUnlinkedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdFactoryWalletUnlinkedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdFactoryWalletUnlinked represents a WalletUnlinked event raised by the IdFactory contract.
type IdFactoryWalletUnlinked struct {
	Wallet   common.Address
	Identity common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterWalletUnlinked is a free log retrieval operation binding the contract event 0x35e6fc363a4bf723d53b26c1a751674aca9c3ead425f0591f84f5540ede86f12.
//
// Solidity: event WalletUnlinked(address indexed wallet, address indexed identity)
func (_IdFactory *IdFactoryFilterer) FilterWalletUnlinked(opts *bind.FilterOpts, wallet []common.Address, identity []common.Address) (*IdFactoryWalletUnlinkedIterator, error) {

	var walletRule []interface{}
	for _, walletItem := range wallet {
		walletRule = append(walletRule, walletItem)
	}
	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _IdFactory.contract.FilterLogs(opts, "WalletUnlinked", walletRule, identityRule)
	if err != nil {
		return nil, err
	}
	return &IdFactoryWalletUnlinkedIterator{contract: _IdFactory.contract, event: "WalletUnlinked", logs: logs, sub: sub}, nil
}

// WatchWalletUnlinked is a free log subscription operation binding the contract event 0x35e6fc363a4bf723d53b26c1a751674aca9c3ead425f0591f84f5540ede86f12.
//
// Solidity: event WalletUnlinked(address indexed wallet, address indexed identity)
func (_IdFactory *IdFactoryFilterer) WatchWalletUnlinked(opts *bind.WatchOpts, sink chan<- *IdFactoryWalletUnlinked, wallet []common.Address, identity []common.Address) (event.Subscription, error) {

	var walletRule []interface{}
	for _, walletItem := range wallet {
		walletRule = append(walletRule, walletItem)
	}
	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _IdFactory.contract.WatchLogs(opts, "WalletUnlinked", walletRule, identityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdFactoryWalletUnlinked)
				if err := _IdFactory.contract.UnpackLog(event, "WalletUnlinked", log); err != nil {
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

// ParseWalletUnlinked is a log parse operation binding the contract event 0x35e6fc363a4bf723d53b26c1a751674aca9c3ead425f0591f84f5540ede86f12.
//
// Solidity: event WalletUnlinked(address indexed wallet, address indexed identity)
func (_IdFactory *IdFactoryFilterer) ParseWalletUnlinked(log types.Log) (*IdFactoryWalletUnlinked, error) {
	event := new(IdFactoryWalletUnlinked)
	if err := _IdFactory.contract.UnpackLog(event, "WalletUnlinked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
