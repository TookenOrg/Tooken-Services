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

// TransferRestrictModuleMetaData contains all meta data concerning the TransferRestrictModule contract.
var TransferRestrictModuleMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"ComplianceBound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"ComplianceUnbound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_userAddress\",\"type\":\"address\"}],\"name\":\"UserAllowed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_userAddress\",\"type\":\"address\"}],\"name\":\"UserDisallowed\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userAddress\",\"type\":\"address\"}],\"name\":\"allowUser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_userAddresses\",\"type\":\"address[]\"}],\"name\":\"batchAllowUsers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_userAddresses\",\"type\":\"address[]\"}],\"name\":\"batchDisallowUsers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"bindCompliance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"canComplianceBind\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userAddress\",\"type\":\"address\"}],\"name\":\"disallowUser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"isComplianceBound\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isPlugAndPlay\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_userAddress\",\"type\":\"address\"}],\"name\":\"isUserAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"moduleBurnAction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"moduleCheck\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"moduleMintAction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"moduleTransferAction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"unbindCompliance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523060805234801561001457600080fd5b50608051611d7961004c6000396000818161088e0152818161091301528181610c1e01528181610ca30152610d8d0152611d796000f3fe6080604052600436106101755760003560e01c806358591fb7116100cb5780638da5cb5b1161007f578063e6f5e80711610059578063e6f5e80714610467578063f104a8c91461027d578063f2fde38b1461047b57600080fd5b80638da5cb5b146103d5578063a9fb3b35146103fd578063bcc210531461044657600080fd5b8063726bb465116100b0578063726bb46514610380578063771c456f146103a05780638129fc1c146103c057600080fd5b806358591fb71461034b578063715018a61461036b57600080fd5b80633659cfe61161012d5780634cf4d295116101075780634cf4d295146102bd5780634f1ef2861461031557806352d1902d1461032857600080fd5b80633659cfe61461025d578063372491a21461027d5780634a9325441461029d57600080fd5b806306fdde031161015e57806306fdde03146101d15780632cb7e1ec1461021d57806332bb1f161461023d57600080fd5b8063013b7ce41461017a5780630694a5fb146101af575b600080fd5b34801561018657600080fd5b5061019a6101953660046119b8565b61049b565b60405190151581526020015b60405180910390f35b3480156101bb57600080fd5b506101cf6101ca366004611a05565b610530565b005b3480156101dd57600080fd5b50604080518082018252601681527f5472616e7366657252657374726963744d6f64756c6500000000000000000000602082015290516101a69190611a44565b34801561022957600080fd5b506101cf610238366004611a77565b6106da565b34801561024957600080fd5b506101cf610258366004611afa565b610761565b34801561026957600080fd5b506101cf610278366004611a05565b610884565b34801561028957600080fd5b506101cf610298366004611ba7565b6109ff565b3480156102a957600080fd5b506101cf6102b8366004611a05565b610a85565b3480156102c957600080fd5b5061019a6102d8366004611a05565b6001600160a01b031660009081527ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae00602052604090205460ff1690565b6101cf610323366004611bd1565b610c14565b34801561033457600080fd5b5061033d610d80565b6040519081526020016101a6565b34801561035757600080fd5b506101cf610366366004611a05565b610e45565b34801561037757600080fd5b506101cf610f30565b34801561038c57600080fd5b506101cf61039b366004611afa565b610f44565b3480156103ac57600080fd5b506101cf6103bb366004611a05565b61106a565b3480156103cc57600080fd5b506101cf611150565b3480156103e157600080fd5b506033546040516001600160a01b0390911681526020016101a6565b34801561040957600080fd5b5061019a610418366004611c77565b6001600160a01b03918216600090815260c96020908152604080832093909416825291909152205460ff1690565b34801561045257600080fd5b5061019a610461366004611a05565b50600190565b34801561047357600080fd5b50600161019a565b34801561048757600080fd5b506101cf610496366004611a05565b61126f565b60006001600160a01b03851615806104ba57506001600160a01b038416155b156104c757506001610528565b6001600160a01b03808316600090815260c9602090815260408083209389168352929052205460ff16156104fd57506001610528565b506001600160a01b03808216600090815260c9602090815260408083209387168352929052205460ff165b949350505050565b3360009081527ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae00602081905260409091205460ff166105b65760405162461bcd60e51b815260206004820152601e60248201527f6f6e6c7920626f756e6420636f6d706c69616e63652063616e2063616c6c000060448201526064015b60405180910390fd5b7ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae006001600160a01b03831661062d5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064016105ad565b336001600160a01b0384161461068f5760405162461bcd60e51b815260206004820152602160248201527f6f6e6c7920636f6d706c69616e636520636f6e74726163742063616e2063616c6044820152601b60fa1b60648201526084016105ad565b6001600160a01b038316600081815260208390526040808220805460ff19169055517f408b49d9be1c914c52a0227e18a077e5a892dddf32a26cfa94a5d9708fad77189190a2505050565b3360009081527ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae00602081905260409091205460ff1661075b5760405162461bcd60e51b815260206004820152601e60248201527f6f6e6c7920626f756e6420636f6d706c69616e63652063616e2063616c6c000060448201526064016105ad565b50505050565b3360009081527ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae00602081905260409091205460ff166107e25760405162461bcd60e51b815260206004820152601e60248201527f6f6e6c7920626f756e6420636f6d706c69616e63652063616e2063616c6c000060448201526064016105ad565b815160005b8181101561075b57600084828151811061080357610803611caa565b60209081029190910181015133600081815260c9845260408082206001600160a01b03851680845290865291819020805460ff191690558051928352938201529092507fe9ab2574b9b111f06f42baaa274030c7228b1dad45b248037101cd11de348240910160405180910390a1508061087c81611cc0565b9150506107e7565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036109115760405162461bcd60e51b815260206004820152602c60248201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060448201526b19195b1959d85d1958d85b1b60a21b60648201526084016105ad565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031661096c7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b6001600160a01b0316146109d75760405162461bcd60e51b815260206004820152602c60248201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060448201526b6163746976652070726f787960a01b60648201526084016105ad565b6109e0816112fc565b604080516000808252602082019092526109fc91839190611304565b50565b3360009081527ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae00602081905260409091205460ff16610a805760405162461bcd60e51b815260206004820152601e60248201527f6f6e6c7920626f756e6420636f6d706c69616e63652063616e2063616c6c000060448201526064016105ad565b505050565b7ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae006001600160a01b038216610afc5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f20616464726573730060448201526064016105ad565b6001600160a01b03821660009081526020829052604090205460ff1615610b655760405162461bcd60e51b815260206004820152601860248201527f636f6d706c69616e636520616c726561647920626f756e64000000000000000060448201526064016105ad565b336001600160a01b03831614610bc75760405162461bcd60e51b815260206004820152602160248201527f6f6e6c7920636f6d706c69616e636520636f6e74726163742063616e2063616c6044820152601b60fa1b60648201526084016105ad565b6001600160a01b038216600081815260208390526040808220805460ff19166001179055517f1f7b76c58fb697eb53c6c7c1becb96911516a136e24d7ced386b2355358b75a39190a25050565b6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000163003610ca15760405162461bcd60e51b815260206004820152602c60248201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060448201526b19195b1959d85d1958d85b1b60a21b60648201526084016105ad565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610cfc7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b6001600160a01b031614610d675760405162461bcd60e51b815260206004820152602c60248201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060448201526b6163746976652070726f787960a01b60648201526084016105ad565b610d70826112fc565b610d7c82826001611304565b5050565b6000306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614610e205760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c000000000000000060648201526084016105ad565b507f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc90565b3360009081527ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae00602081905260409091205460ff16610ec65760405162461bcd60e51b815260206004820152601e60248201527f6f6e6c7920626f756e6420636f6d706c69616e63652063616e2063616c6c000060448201526064016105ad565b33600081815260c9602090815260408083206001600160a01b03871680855290835292819020805460ff191690558051938452908301919091527fe9ab2574b9b111f06f42baaa274030c7228b1dad45b248037101cd11de34824091015b60405180910390a15050565b610f386114a4565b610f4260006114fe565b565b3360009081527ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae00602081905260409091205460ff16610fc55760405162461bcd60e51b815260206004820152601e60248201527f6f6e6c7920626f756e6420636f6d706c69616e63652063616e2063616c6c000060448201526064016105ad565b815160005b8181101561075b576000848281518110610fe657610fe6611caa565b60209081029190910181015133600081815260c9845260408082206001600160a01b03851680845290865291819020805460ff191660011790558051928352938201529092507f767240d3c58c1058ed268ca225eb1cc93d85cf46b6e85d8d120b39921475b1c1910160405180910390a1508061106281611cc0565b915050610fca565b3360009081527ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae00602081905260409091205460ff166110eb5760405162461bcd60e51b815260206004820152601e60248201527f6f6e6c7920626f756e6420636f6d706c69616e63652063616e2063616c6c000060448201526064016105ad565b33600081815260c9602090815260408083206001600160a01b03871680855290835292819020805460ff191660011790558051938452908301919091527f767240d3c58c1058ed268ca225eb1cc93d85cf46b6e85d8d120b39921475b1c19101610f24565b600054610100900460ff16158080156111705750600054600160ff909116105b8061118a5750303b15801561118a575060005460ff166001145b6111fc5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084016105ad565b6000805460ff19166001179055801561121f576000805461ff0019166101001790555b61122761155d565b80156109fc576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a150565b6112776114a4565b6001600160a01b0381166112f35760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f646472657373000000000000000000000000000000000000000000000000000060648201526084016105ad565b6109fc816114fe565b6109fc6114a4565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff161561133757610a80836115d8565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015611391575060408051601f3d908101601f1916820190925261138e91810190611ce7565b60015b6114035760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201527f6f6e206973206e6f74205555505300000000000000000000000000000000000060648201526084016105ad565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc81146114985760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f7860448201527f6961626c6555554944000000000000000000000000000000000000000000000060648201526084016105ad565b50610a808383836116a3565b6033546001600160a01b03163314610f425760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016105ad565b603380546001600160a01b0383811673ffffffffffffffffffffffffffffffffffffffff19831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166115c85760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b60648201526084016105ad565b6115d06116c8565b610f4261173b565b6001600160a01b0381163b6116555760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201527f6f74206120636f6e74726163740000000000000000000000000000000000000060648201526084016105ad565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0392909216919091179055565b6116ac836117a6565b6000825111806116b95750805b15610a805761075b83836117e6565b600054610100900460ff166117335760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b60648201526084016105ad565b610f42611812565b600054610100900460ff16610f425760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b60648201526084016105ad565b6117af816115d8565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b606061180b8383604051806060016040528060278152602001611d1d60279139611886565b9392505050565b600054610100900460ff1661187d5760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b60648201526084016105ad565b610f42336114fe565b6060600080856001600160a01b0316856040516118a39190611d00565b600060405180830381855af49150503d80600081146118de576040519150601f19603f3d011682016040523d82523d6000602084013e6118e3565b606091505b50915091506118f4868383876118fe565b9695505050505050565b6060831561196d578251600003611966576001600160a01b0385163b6119665760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016105ad565b5081610528565b61052883838151156119825781518083602001fd5b8060405162461bcd60e51b81526004016105ad9190611a44565b80356001600160a01b03811681146119b357600080fd5b919050565b600080600080608085870312156119ce57600080fd5b6119d78561199c565b93506119e56020860161199c565b9250604085013591506119fa6060860161199c565b905092959194509250565b600060208284031215611a1757600080fd5b61180b8261199c565b60005b83811015611a3b578181015183820152602001611a23565b50506000910152565b6020815260008251806020840152611a63816040850160208701611a20565b601f01601f19169190910160400192915050565b600080600060608486031215611a8c57600080fd5b611a958461199c565b9250611aa36020850161199c565b9150604084013590509250925092565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715611af257611af2611ab3565b604052919050565b60006020808385031215611b0d57600080fd5b823567ffffffffffffffff80821115611b2557600080fd5b818501915085601f830112611b3957600080fd5b813581811115611b4b57611b4b611ab3565b8060051b9150611b5c848301611ac9565b8181529183018401918481019088841115611b7657600080fd5b938501935b83851015611b9b57611b8c8561199c565b82529385019390850190611b7b565b98975050505050505050565b60008060408385031215611bba57600080fd5b611bc38361199c565b946020939093013593505050565b60008060408385031215611be457600080fd5b611bed8361199c565b915060208084013567ffffffffffffffff80821115611c0b57600080fd5b818601915086601f830112611c1f57600080fd5b813581811115611c3157611c31611ab3565b611c43601f8201601f19168501611ac9565b91508082528784828501011115611c5957600080fd5b80848401858401376000848284010152508093505050509250929050565b60008060408385031215611c8a57600080fd5b611c938361199c565b9150611ca16020840161199c565b90509250929050565b634e487b7160e01b600052603260045260246000fd5b600060018201611ce057634e487b7160e01b600052601160045260246000fd5b5060010190565b600060208284031215611cf957600080fd5b5051919050565b60008251611d12818460208701611a20565b919091019291505056fe416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a264697066735822122023d512e02d66c20938ab45325a31edb7f13c967f55c162d6898ed1431ae4406764736f6c63430008110033",
}

// TransferRestrictModuleABI is the input ABI used to generate the binding from.
// Deprecated: Use TransferRestrictModuleMetaData.ABI instead.
var TransferRestrictModuleABI = TransferRestrictModuleMetaData.ABI

// TransferRestrictModuleBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TransferRestrictModuleMetaData.Bin instead.
var TransferRestrictModuleBin = TransferRestrictModuleMetaData.Bin

// DeployTransferRestrictModule deploys a new Ethereum contract, binding an instance of TransferRestrictModule to it.
func DeployTransferRestrictModule(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TransferRestrictModule, error) {
	parsed, err := TransferRestrictModuleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TransferRestrictModuleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TransferRestrictModule{TransferRestrictModuleCaller: TransferRestrictModuleCaller{contract: contract}, TransferRestrictModuleTransactor: TransferRestrictModuleTransactor{contract: contract}, TransferRestrictModuleFilterer: TransferRestrictModuleFilterer{contract: contract}}, nil
}

// TransferRestrictModule is an auto generated Go binding around an Ethereum contract.
type TransferRestrictModule struct {
	TransferRestrictModuleCaller     // Read-only binding to the contract
	TransferRestrictModuleTransactor // Write-only binding to the contract
	TransferRestrictModuleFilterer   // Log filterer for contract events
}

// TransferRestrictModuleCaller is an auto generated read-only Go binding around an Ethereum contract.
type TransferRestrictModuleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferRestrictModuleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TransferRestrictModuleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferRestrictModuleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TransferRestrictModuleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferRestrictModuleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TransferRestrictModuleSession struct {
	Contract     *TransferRestrictModule // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// TransferRestrictModuleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TransferRestrictModuleCallerSession struct {
	Contract *TransferRestrictModuleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// TransferRestrictModuleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TransferRestrictModuleTransactorSession struct {
	Contract     *TransferRestrictModuleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// TransferRestrictModuleRaw is an auto generated low-level Go binding around an Ethereum contract.
type TransferRestrictModuleRaw struct {
	Contract *TransferRestrictModule // Generic contract binding to access the raw methods on
}

// TransferRestrictModuleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TransferRestrictModuleCallerRaw struct {
	Contract *TransferRestrictModuleCaller // Generic read-only contract binding to access the raw methods on
}

// TransferRestrictModuleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TransferRestrictModuleTransactorRaw struct {
	Contract *TransferRestrictModuleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTransferRestrictModule creates a new instance of TransferRestrictModule, bound to a specific deployed contract.
func NewTransferRestrictModule(address common.Address, backend bind.ContractBackend) (*TransferRestrictModule, error) {
	contract, err := bindTransferRestrictModule(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TransferRestrictModule{TransferRestrictModuleCaller: TransferRestrictModuleCaller{contract: contract}, TransferRestrictModuleTransactor: TransferRestrictModuleTransactor{contract: contract}, TransferRestrictModuleFilterer: TransferRestrictModuleFilterer{contract: contract}}, nil
}

// NewTransferRestrictModuleCaller creates a new read-only instance of TransferRestrictModule, bound to a specific deployed contract.
func NewTransferRestrictModuleCaller(address common.Address, caller bind.ContractCaller) (*TransferRestrictModuleCaller, error) {
	contract, err := bindTransferRestrictModule(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TransferRestrictModuleCaller{contract: contract}, nil
}

// NewTransferRestrictModuleTransactor creates a new write-only instance of TransferRestrictModule, bound to a specific deployed contract.
func NewTransferRestrictModuleTransactor(address common.Address, transactor bind.ContractTransactor) (*TransferRestrictModuleTransactor, error) {
	contract, err := bindTransferRestrictModule(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TransferRestrictModuleTransactor{contract: contract}, nil
}

// NewTransferRestrictModuleFilterer creates a new log filterer instance of TransferRestrictModule, bound to a specific deployed contract.
func NewTransferRestrictModuleFilterer(address common.Address, filterer bind.ContractFilterer) (*TransferRestrictModuleFilterer, error) {
	contract, err := bindTransferRestrictModule(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TransferRestrictModuleFilterer{contract: contract}, nil
}

// bindTransferRestrictModule binds a generic wrapper to an already deployed contract.
func bindTransferRestrictModule(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TransferRestrictModuleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransferRestrictModule *TransferRestrictModuleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransferRestrictModule.Contract.TransferRestrictModuleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransferRestrictModule *TransferRestrictModuleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.TransferRestrictModuleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransferRestrictModule *TransferRestrictModuleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.TransferRestrictModuleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransferRestrictModule *TransferRestrictModuleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransferRestrictModule.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransferRestrictModule *TransferRestrictModuleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransferRestrictModule *TransferRestrictModuleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.contract.Transact(opts, method, params...)
}

// CanComplianceBind is a free data retrieval call binding the contract method 0xbcc21053.
//
// Solidity: function canComplianceBind(address ) view returns(bool)
func (_TransferRestrictModule *TransferRestrictModuleCaller) CanComplianceBind(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _TransferRestrictModule.contract.Call(opts, &out, "canComplianceBind", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanComplianceBind is a free data retrieval call binding the contract method 0xbcc21053.
//
// Solidity: function canComplianceBind(address ) view returns(bool)
func (_TransferRestrictModule *TransferRestrictModuleSession) CanComplianceBind(arg0 common.Address) (bool, error) {
	return _TransferRestrictModule.Contract.CanComplianceBind(&_TransferRestrictModule.CallOpts, arg0)
}

// CanComplianceBind is a free data retrieval call binding the contract method 0xbcc21053.
//
// Solidity: function canComplianceBind(address ) view returns(bool)
func (_TransferRestrictModule *TransferRestrictModuleCallerSession) CanComplianceBind(arg0 common.Address) (bool, error) {
	return _TransferRestrictModule.Contract.CanComplianceBind(&_TransferRestrictModule.CallOpts, arg0)
}

// IsComplianceBound is a free data retrieval call binding the contract method 0x4cf4d295.
//
// Solidity: function isComplianceBound(address _compliance) view returns(bool)
func (_TransferRestrictModule *TransferRestrictModuleCaller) IsComplianceBound(opts *bind.CallOpts, _compliance common.Address) (bool, error) {
	var out []interface{}
	err := _TransferRestrictModule.contract.Call(opts, &out, "isComplianceBound", _compliance)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsComplianceBound is a free data retrieval call binding the contract method 0x4cf4d295.
//
// Solidity: function isComplianceBound(address _compliance) view returns(bool)
func (_TransferRestrictModule *TransferRestrictModuleSession) IsComplianceBound(_compliance common.Address) (bool, error) {
	return _TransferRestrictModule.Contract.IsComplianceBound(&_TransferRestrictModule.CallOpts, _compliance)
}

// IsComplianceBound is a free data retrieval call binding the contract method 0x4cf4d295.
//
// Solidity: function isComplianceBound(address _compliance) view returns(bool)
func (_TransferRestrictModule *TransferRestrictModuleCallerSession) IsComplianceBound(_compliance common.Address) (bool, error) {
	return _TransferRestrictModule.Contract.IsComplianceBound(&_TransferRestrictModule.CallOpts, _compliance)
}

// IsPlugAndPlay is a free data retrieval call binding the contract method 0xe6f5e807.
//
// Solidity: function isPlugAndPlay() pure returns(bool)
func (_TransferRestrictModule *TransferRestrictModuleCaller) IsPlugAndPlay(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TransferRestrictModule.contract.Call(opts, &out, "isPlugAndPlay")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPlugAndPlay is a free data retrieval call binding the contract method 0xe6f5e807.
//
// Solidity: function isPlugAndPlay() pure returns(bool)
func (_TransferRestrictModule *TransferRestrictModuleSession) IsPlugAndPlay() (bool, error) {
	return _TransferRestrictModule.Contract.IsPlugAndPlay(&_TransferRestrictModule.CallOpts)
}

// IsPlugAndPlay is a free data retrieval call binding the contract method 0xe6f5e807.
//
// Solidity: function isPlugAndPlay() pure returns(bool)
func (_TransferRestrictModule *TransferRestrictModuleCallerSession) IsPlugAndPlay() (bool, error) {
	return _TransferRestrictModule.Contract.IsPlugAndPlay(&_TransferRestrictModule.CallOpts)
}

// IsUserAllowed is a free data retrieval call binding the contract method 0xa9fb3b35.
//
// Solidity: function isUserAllowed(address _compliance, address _userAddress) view returns(bool)
func (_TransferRestrictModule *TransferRestrictModuleCaller) IsUserAllowed(opts *bind.CallOpts, _compliance common.Address, _userAddress common.Address) (bool, error) {
	var out []interface{}
	err := _TransferRestrictModule.contract.Call(opts, &out, "isUserAllowed", _compliance, _userAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsUserAllowed is a free data retrieval call binding the contract method 0xa9fb3b35.
//
// Solidity: function isUserAllowed(address _compliance, address _userAddress) view returns(bool)
func (_TransferRestrictModule *TransferRestrictModuleSession) IsUserAllowed(_compliance common.Address, _userAddress common.Address) (bool, error) {
	return _TransferRestrictModule.Contract.IsUserAllowed(&_TransferRestrictModule.CallOpts, _compliance, _userAddress)
}

// IsUserAllowed is a free data retrieval call binding the contract method 0xa9fb3b35.
//
// Solidity: function isUserAllowed(address _compliance, address _userAddress) view returns(bool)
func (_TransferRestrictModule *TransferRestrictModuleCallerSession) IsUserAllowed(_compliance common.Address, _userAddress common.Address) (bool, error) {
	return _TransferRestrictModule.Contract.IsUserAllowed(&_TransferRestrictModule.CallOpts, _compliance, _userAddress)
}

// ModuleCheck is a free data retrieval call binding the contract method 0x013b7ce4.
//
// Solidity: function moduleCheck(address _from, address _to, uint256 , address _compliance) view returns(bool)
func (_TransferRestrictModule *TransferRestrictModuleCaller) ModuleCheck(opts *bind.CallOpts, _from common.Address, _to common.Address, arg2 *big.Int, _compliance common.Address) (bool, error) {
	var out []interface{}
	err := _TransferRestrictModule.contract.Call(opts, &out, "moduleCheck", _from, _to, arg2, _compliance)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ModuleCheck is a free data retrieval call binding the contract method 0x013b7ce4.
//
// Solidity: function moduleCheck(address _from, address _to, uint256 , address _compliance) view returns(bool)
func (_TransferRestrictModule *TransferRestrictModuleSession) ModuleCheck(_from common.Address, _to common.Address, arg2 *big.Int, _compliance common.Address) (bool, error) {
	return _TransferRestrictModule.Contract.ModuleCheck(&_TransferRestrictModule.CallOpts, _from, _to, arg2, _compliance)
}

// ModuleCheck is a free data retrieval call binding the contract method 0x013b7ce4.
//
// Solidity: function moduleCheck(address _from, address _to, uint256 , address _compliance) view returns(bool)
func (_TransferRestrictModule *TransferRestrictModuleCallerSession) ModuleCheck(_from common.Address, _to common.Address, arg2 *big.Int, _compliance common.Address) (bool, error) {
	return _TransferRestrictModule.Contract.ModuleCheck(&_TransferRestrictModule.CallOpts, _from, _to, arg2, _compliance)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string _name)
func (_TransferRestrictModule *TransferRestrictModuleCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TransferRestrictModule.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string _name)
func (_TransferRestrictModule *TransferRestrictModuleSession) Name() (string, error) {
	return _TransferRestrictModule.Contract.Name(&_TransferRestrictModule.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string _name)
func (_TransferRestrictModule *TransferRestrictModuleCallerSession) Name() (string, error) {
	return _TransferRestrictModule.Contract.Name(&_TransferRestrictModule.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TransferRestrictModule *TransferRestrictModuleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TransferRestrictModule.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TransferRestrictModule *TransferRestrictModuleSession) Owner() (common.Address, error) {
	return _TransferRestrictModule.Contract.Owner(&_TransferRestrictModule.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TransferRestrictModule *TransferRestrictModuleCallerSession) Owner() (common.Address, error) {
	return _TransferRestrictModule.Contract.Owner(&_TransferRestrictModule.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_TransferRestrictModule *TransferRestrictModuleCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TransferRestrictModule.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_TransferRestrictModule *TransferRestrictModuleSession) ProxiableUUID() ([32]byte, error) {
	return _TransferRestrictModule.Contract.ProxiableUUID(&_TransferRestrictModule.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_TransferRestrictModule *TransferRestrictModuleCallerSession) ProxiableUUID() ([32]byte, error) {
	return _TransferRestrictModule.Contract.ProxiableUUID(&_TransferRestrictModule.CallOpts)
}

// AllowUser is a paid mutator transaction binding the contract method 0x771c456f.
//
// Solidity: function allowUser(address _userAddress) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactor) AllowUser(opts *bind.TransactOpts, _userAddress common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.contract.Transact(opts, "allowUser", _userAddress)
}

// AllowUser is a paid mutator transaction binding the contract method 0x771c456f.
//
// Solidity: function allowUser(address _userAddress) returns()
func (_TransferRestrictModule *TransferRestrictModuleSession) AllowUser(_userAddress common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.AllowUser(&_TransferRestrictModule.TransactOpts, _userAddress)
}

// AllowUser is a paid mutator transaction binding the contract method 0x771c456f.
//
// Solidity: function allowUser(address _userAddress) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactorSession) AllowUser(_userAddress common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.AllowUser(&_TransferRestrictModule.TransactOpts, _userAddress)
}

// BatchAllowUsers is a paid mutator transaction binding the contract method 0x726bb465.
//
// Solidity: function batchAllowUsers(address[] _userAddresses) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactor) BatchAllowUsers(opts *bind.TransactOpts, _userAddresses []common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.contract.Transact(opts, "batchAllowUsers", _userAddresses)
}

// BatchAllowUsers is a paid mutator transaction binding the contract method 0x726bb465.
//
// Solidity: function batchAllowUsers(address[] _userAddresses) returns()
func (_TransferRestrictModule *TransferRestrictModuleSession) BatchAllowUsers(_userAddresses []common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.BatchAllowUsers(&_TransferRestrictModule.TransactOpts, _userAddresses)
}

// BatchAllowUsers is a paid mutator transaction binding the contract method 0x726bb465.
//
// Solidity: function batchAllowUsers(address[] _userAddresses) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactorSession) BatchAllowUsers(_userAddresses []common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.BatchAllowUsers(&_TransferRestrictModule.TransactOpts, _userAddresses)
}

// BatchDisallowUsers is a paid mutator transaction binding the contract method 0x32bb1f16.
//
// Solidity: function batchDisallowUsers(address[] _userAddresses) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactor) BatchDisallowUsers(opts *bind.TransactOpts, _userAddresses []common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.contract.Transact(opts, "batchDisallowUsers", _userAddresses)
}

// BatchDisallowUsers is a paid mutator transaction binding the contract method 0x32bb1f16.
//
// Solidity: function batchDisallowUsers(address[] _userAddresses) returns()
func (_TransferRestrictModule *TransferRestrictModuleSession) BatchDisallowUsers(_userAddresses []common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.BatchDisallowUsers(&_TransferRestrictModule.TransactOpts, _userAddresses)
}

// BatchDisallowUsers is a paid mutator transaction binding the contract method 0x32bb1f16.
//
// Solidity: function batchDisallowUsers(address[] _userAddresses) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactorSession) BatchDisallowUsers(_userAddresses []common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.BatchDisallowUsers(&_TransferRestrictModule.TransactOpts, _userAddresses)
}

// BindCompliance is a paid mutator transaction binding the contract method 0x4a932544.
//
// Solidity: function bindCompliance(address _compliance) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactor) BindCompliance(opts *bind.TransactOpts, _compliance common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.contract.Transact(opts, "bindCompliance", _compliance)
}

// BindCompliance is a paid mutator transaction binding the contract method 0x4a932544.
//
// Solidity: function bindCompliance(address _compliance) returns()
func (_TransferRestrictModule *TransferRestrictModuleSession) BindCompliance(_compliance common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.BindCompliance(&_TransferRestrictModule.TransactOpts, _compliance)
}

// BindCompliance is a paid mutator transaction binding the contract method 0x4a932544.
//
// Solidity: function bindCompliance(address _compliance) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactorSession) BindCompliance(_compliance common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.BindCompliance(&_TransferRestrictModule.TransactOpts, _compliance)
}

// DisallowUser is a paid mutator transaction binding the contract method 0x58591fb7.
//
// Solidity: function disallowUser(address _userAddress) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactor) DisallowUser(opts *bind.TransactOpts, _userAddress common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.contract.Transact(opts, "disallowUser", _userAddress)
}

// DisallowUser is a paid mutator transaction binding the contract method 0x58591fb7.
//
// Solidity: function disallowUser(address _userAddress) returns()
func (_TransferRestrictModule *TransferRestrictModuleSession) DisallowUser(_userAddress common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.DisallowUser(&_TransferRestrictModule.TransactOpts, _userAddress)
}

// DisallowUser is a paid mutator transaction binding the contract method 0x58591fb7.
//
// Solidity: function disallowUser(address _userAddress) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactorSession) DisallowUser(_userAddress common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.DisallowUser(&_TransferRestrictModule.TransactOpts, _userAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferRestrictModule.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_TransferRestrictModule *TransferRestrictModuleSession) Initialize() (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.Initialize(&_TransferRestrictModule.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactorSession) Initialize() (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.Initialize(&_TransferRestrictModule.TransactOpts)
}

// ModuleBurnAction is a paid mutator transaction binding the contract method 0x372491a2.
//
// Solidity: function moduleBurnAction(address _from, uint256 _value) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactor) ModuleBurnAction(opts *bind.TransactOpts, _from common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TransferRestrictModule.contract.Transact(opts, "moduleBurnAction", _from, _value)
}

// ModuleBurnAction is a paid mutator transaction binding the contract method 0x372491a2.
//
// Solidity: function moduleBurnAction(address _from, uint256 _value) returns()
func (_TransferRestrictModule *TransferRestrictModuleSession) ModuleBurnAction(_from common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.ModuleBurnAction(&_TransferRestrictModule.TransactOpts, _from, _value)
}

// ModuleBurnAction is a paid mutator transaction binding the contract method 0x372491a2.
//
// Solidity: function moduleBurnAction(address _from, uint256 _value) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactorSession) ModuleBurnAction(_from common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.ModuleBurnAction(&_TransferRestrictModule.TransactOpts, _from, _value)
}

// ModuleMintAction is a paid mutator transaction binding the contract method 0xf104a8c9.
//
// Solidity: function moduleMintAction(address _to, uint256 _value) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactor) ModuleMintAction(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TransferRestrictModule.contract.Transact(opts, "moduleMintAction", _to, _value)
}

// ModuleMintAction is a paid mutator transaction binding the contract method 0xf104a8c9.
//
// Solidity: function moduleMintAction(address _to, uint256 _value) returns()
func (_TransferRestrictModule *TransferRestrictModuleSession) ModuleMintAction(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.ModuleMintAction(&_TransferRestrictModule.TransactOpts, _to, _value)
}

// ModuleMintAction is a paid mutator transaction binding the contract method 0xf104a8c9.
//
// Solidity: function moduleMintAction(address _to, uint256 _value) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactorSession) ModuleMintAction(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.ModuleMintAction(&_TransferRestrictModule.TransactOpts, _to, _value)
}

// ModuleTransferAction is a paid mutator transaction binding the contract method 0x2cb7e1ec.
//
// Solidity: function moduleTransferAction(address _from, address _to, uint256 _value) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactor) ModuleTransferAction(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TransferRestrictModule.contract.Transact(opts, "moduleTransferAction", _from, _to, _value)
}

// ModuleTransferAction is a paid mutator transaction binding the contract method 0x2cb7e1ec.
//
// Solidity: function moduleTransferAction(address _from, address _to, uint256 _value) returns()
func (_TransferRestrictModule *TransferRestrictModuleSession) ModuleTransferAction(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.ModuleTransferAction(&_TransferRestrictModule.TransactOpts, _from, _to, _value)
}

// ModuleTransferAction is a paid mutator transaction binding the contract method 0x2cb7e1ec.
//
// Solidity: function moduleTransferAction(address _from, address _to, uint256 _value) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactorSession) ModuleTransferAction(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.ModuleTransferAction(&_TransferRestrictModule.TransactOpts, _from, _to, _value)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferRestrictModule.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TransferRestrictModule *TransferRestrictModuleSession) RenounceOwnership() (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.RenounceOwnership(&_TransferRestrictModule.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.RenounceOwnership(&_TransferRestrictModule.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TransferRestrictModule *TransferRestrictModuleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.TransferOwnership(&_TransferRestrictModule.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.TransferOwnership(&_TransferRestrictModule.TransactOpts, newOwner)
}

// UnbindCompliance is a paid mutator transaction binding the contract method 0x0694a5fb.
//
// Solidity: function unbindCompliance(address _compliance) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactor) UnbindCompliance(opts *bind.TransactOpts, _compliance common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.contract.Transact(opts, "unbindCompliance", _compliance)
}

// UnbindCompliance is a paid mutator transaction binding the contract method 0x0694a5fb.
//
// Solidity: function unbindCompliance(address _compliance) returns()
func (_TransferRestrictModule *TransferRestrictModuleSession) UnbindCompliance(_compliance common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.UnbindCompliance(&_TransferRestrictModule.TransactOpts, _compliance)
}

// UnbindCompliance is a paid mutator transaction binding the contract method 0x0694a5fb.
//
// Solidity: function unbindCompliance(address _compliance) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactorSession) UnbindCompliance(_compliance common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.UnbindCompliance(&_TransferRestrictModule.TransactOpts, _compliance)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_TransferRestrictModule *TransferRestrictModuleSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.UpgradeTo(&_TransferRestrictModule.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.UpgradeTo(&_TransferRestrictModule.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _TransferRestrictModule.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_TransferRestrictModule *TransferRestrictModuleSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.UpgradeToAndCall(&_TransferRestrictModule.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_TransferRestrictModule *TransferRestrictModuleTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _TransferRestrictModule.Contract.UpgradeToAndCall(&_TransferRestrictModule.TransactOpts, newImplementation, data)
}

// TransferRestrictModuleAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the TransferRestrictModule contract.
type TransferRestrictModuleAdminChangedIterator struct {
	Event *TransferRestrictModuleAdminChanged // Event containing the contract specifics and raw log

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
func (it *TransferRestrictModuleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferRestrictModuleAdminChanged)
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
		it.Event = new(TransferRestrictModuleAdminChanged)
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
func (it *TransferRestrictModuleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferRestrictModuleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferRestrictModuleAdminChanged represents a AdminChanged event raised by the TransferRestrictModule contract.
type TransferRestrictModuleAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*TransferRestrictModuleAdminChangedIterator, error) {

	logs, sub, err := _TransferRestrictModule.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &TransferRestrictModuleAdminChangedIterator{contract: _TransferRestrictModule.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *TransferRestrictModuleAdminChanged) (event.Subscription, error) {

	logs, sub, err := _TransferRestrictModule.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferRestrictModuleAdminChanged)
				if err := _TransferRestrictModule.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) ParseAdminChanged(log types.Log) (*TransferRestrictModuleAdminChanged, error) {
	event := new(TransferRestrictModuleAdminChanged)
	if err := _TransferRestrictModule.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferRestrictModuleBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the TransferRestrictModule contract.
type TransferRestrictModuleBeaconUpgradedIterator struct {
	Event *TransferRestrictModuleBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *TransferRestrictModuleBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferRestrictModuleBeaconUpgraded)
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
		it.Event = new(TransferRestrictModuleBeaconUpgraded)
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
func (it *TransferRestrictModuleBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferRestrictModuleBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferRestrictModuleBeaconUpgraded represents a BeaconUpgraded event raised by the TransferRestrictModule contract.
type TransferRestrictModuleBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*TransferRestrictModuleBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _TransferRestrictModule.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &TransferRestrictModuleBeaconUpgradedIterator{contract: _TransferRestrictModule.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *TransferRestrictModuleBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _TransferRestrictModule.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferRestrictModuleBeaconUpgraded)
				if err := _TransferRestrictModule.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) ParseBeaconUpgraded(log types.Log) (*TransferRestrictModuleBeaconUpgraded, error) {
	event := new(TransferRestrictModuleBeaconUpgraded)
	if err := _TransferRestrictModule.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferRestrictModuleComplianceBoundIterator is returned from FilterComplianceBound and is used to iterate over the raw logs and unpacked data for ComplianceBound events raised by the TransferRestrictModule contract.
type TransferRestrictModuleComplianceBoundIterator struct {
	Event *TransferRestrictModuleComplianceBound // Event containing the contract specifics and raw log

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
func (it *TransferRestrictModuleComplianceBoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferRestrictModuleComplianceBound)
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
		it.Event = new(TransferRestrictModuleComplianceBound)
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
func (it *TransferRestrictModuleComplianceBoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferRestrictModuleComplianceBoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferRestrictModuleComplianceBound represents a ComplianceBound event raised by the TransferRestrictModule contract.
type TransferRestrictModuleComplianceBound struct {
	Compliance common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterComplianceBound is a free log retrieval operation binding the contract event 0x1f7b76c58fb697eb53c6c7c1becb96911516a136e24d7ced386b2355358b75a3.
//
// Solidity: event ComplianceBound(address indexed _compliance)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) FilterComplianceBound(opts *bind.FilterOpts, _compliance []common.Address) (*TransferRestrictModuleComplianceBoundIterator, error) {

	var _complianceRule []interface{}
	for _, _complianceItem := range _compliance {
		_complianceRule = append(_complianceRule, _complianceItem)
	}

	logs, sub, err := _TransferRestrictModule.contract.FilterLogs(opts, "ComplianceBound", _complianceRule)
	if err != nil {
		return nil, err
	}
	return &TransferRestrictModuleComplianceBoundIterator{contract: _TransferRestrictModule.contract, event: "ComplianceBound", logs: logs, sub: sub}, nil
}

// WatchComplianceBound is a free log subscription operation binding the contract event 0x1f7b76c58fb697eb53c6c7c1becb96911516a136e24d7ced386b2355358b75a3.
//
// Solidity: event ComplianceBound(address indexed _compliance)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) WatchComplianceBound(opts *bind.WatchOpts, sink chan<- *TransferRestrictModuleComplianceBound, _compliance []common.Address) (event.Subscription, error) {

	var _complianceRule []interface{}
	for _, _complianceItem := range _compliance {
		_complianceRule = append(_complianceRule, _complianceItem)
	}

	logs, sub, err := _TransferRestrictModule.contract.WatchLogs(opts, "ComplianceBound", _complianceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferRestrictModuleComplianceBound)
				if err := _TransferRestrictModule.contract.UnpackLog(event, "ComplianceBound", log); err != nil {
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

// ParseComplianceBound is a log parse operation binding the contract event 0x1f7b76c58fb697eb53c6c7c1becb96911516a136e24d7ced386b2355358b75a3.
//
// Solidity: event ComplianceBound(address indexed _compliance)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) ParseComplianceBound(log types.Log) (*TransferRestrictModuleComplianceBound, error) {
	event := new(TransferRestrictModuleComplianceBound)
	if err := _TransferRestrictModule.contract.UnpackLog(event, "ComplianceBound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferRestrictModuleComplianceUnboundIterator is returned from FilterComplianceUnbound and is used to iterate over the raw logs and unpacked data for ComplianceUnbound events raised by the TransferRestrictModule contract.
type TransferRestrictModuleComplianceUnboundIterator struct {
	Event *TransferRestrictModuleComplianceUnbound // Event containing the contract specifics and raw log

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
func (it *TransferRestrictModuleComplianceUnboundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferRestrictModuleComplianceUnbound)
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
		it.Event = new(TransferRestrictModuleComplianceUnbound)
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
func (it *TransferRestrictModuleComplianceUnboundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferRestrictModuleComplianceUnboundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferRestrictModuleComplianceUnbound represents a ComplianceUnbound event raised by the TransferRestrictModule contract.
type TransferRestrictModuleComplianceUnbound struct {
	Compliance common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterComplianceUnbound is a free log retrieval operation binding the contract event 0x408b49d9be1c914c52a0227e18a077e5a892dddf32a26cfa94a5d9708fad7718.
//
// Solidity: event ComplianceUnbound(address indexed _compliance)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) FilterComplianceUnbound(opts *bind.FilterOpts, _compliance []common.Address) (*TransferRestrictModuleComplianceUnboundIterator, error) {

	var _complianceRule []interface{}
	for _, _complianceItem := range _compliance {
		_complianceRule = append(_complianceRule, _complianceItem)
	}

	logs, sub, err := _TransferRestrictModule.contract.FilterLogs(opts, "ComplianceUnbound", _complianceRule)
	if err != nil {
		return nil, err
	}
	return &TransferRestrictModuleComplianceUnboundIterator{contract: _TransferRestrictModule.contract, event: "ComplianceUnbound", logs: logs, sub: sub}, nil
}

// WatchComplianceUnbound is a free log subscription operation binding the contract event 0x408b49d9be1c914c52a0227e18a077e5a892dddf32a26cfa94a5d9708fad7718.
//
// Solidity: event ComplianceUnbound(address indexed _compliance)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) WatchComplianceUnbound(opts *bind.WatchOpts, sink chan<- *TransferRestrictModuleComplianceUnbound, _compliance []common.Address) (event.Subscription, error) {

	var _complianceRule []interface{}
	for _, _complianceItem := range _compliance {
		_complianceRule = append(_complianceRule, _complianceItem)
	}

	logs, sub, err := _TransferRestrictModule.contract.WatchLogs(opts, "ComplianceUnbound", _complianceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferRestrictModuleComplianceUnbound)
				if err := _TransferRestrictModule.contract.UnpackLog(event, "ComplianceUnbound", log); err != nil {
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

// ParseComplianceUnbound is a log parse operation binding the contract event 0x408b49d9be1c914c52a0227e18a077e5a892dddf32a26cfa94a5d9708fad7718.
//
// Solidity: event ComplianceUnbound(address indexed _compliance)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) ParseComplianceUnbound(log types.Log) (*TransferRestrictModuleComplianceUnbound, error) {
	event := new(TransferRestrictModuleComplianceUnbound)
	if err := _TransferRestrictModule.contract.UnpackLog(event, "ComplianceUnbound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferRestrictModuleInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the TransferRestrictModule contract.
type TransferRestrictModuleInitializedIterator struct {
	Event *TransferRestrictModuleInitialized // Event containing the contract specifics and raw log

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
func (it *TransferRestrictModuleInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferRestrictModuleInitialized)
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
		it.Event = new(TransferRestrictModuleInitialized)
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
func (it *TransferRestrictModuleInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferRestrictModuleInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferRestrictModuleInitialized represents a Initialized event raised by the TransferRestrictModule contract.
type TransferRestrictModuleInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) FilterInitialized(opts *bind.FilterOpts) (*TransferRestrictModuleInitializedIterator, error) {

	logs, sub, err := _TransferRestrictModule.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &TransferRestrictModuleInitializedIterator{contract: _TransferRestrictModule.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *TransferRestrictModuleInitialized) (event.Subscription, error) {

	logs, sub, err := _TransferRestrictModule.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferRestrictModuleInitialized)
				if err := _TransferRestrictModule.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_TransferRestrictModule *TransferRestrictModuleFilterer) ParseInitialized(log types.Log) (*TransferRestrictModuleInitialized, error) {
	event := new(TransferRestrictModuleInitialized)
	if err := _TransferRestrictModule.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferRestrictModuleOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TransferRestrictModule contract.
type TransferRestrictModuleOwnershipTransferredIterator struct {
	Event *TransferRestrictModuleOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TransferRestrictModuleOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferRestrictModuleOwnershipTransferred)
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
		it.Event = new(TransferRestrictModuleOwnershipTransferred)
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
func (it *TransferRestrictModuleOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferRestrictModuleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferRestrictModuleOwnershipTransferred represents a OwnershipTransferred event raised by the TransferRestrictModule contract.
type TransferRestrictModuleOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TransferRestrictModuleOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TransferRestrictModule.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TransferRestrictModuleOwnershipTransferredIterator{contract: _TransferRestrictModule.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TransferRestrictModuleOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TransferRestrictModule.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferRestrictModuleOwnershipTransferred)
				if err := _TransferRestrictModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TransferRestrictModule *TransferRestrictModuleFilterer) ParseOwnershipTransferred(log types.Log) (*TransferRestrictModuleOwnershipTransferred, error) {
	event := new(TransferRestrictModuleOwnershipTransferred)
	if err := _TransferRestrictModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferRestrictModuleUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the TransferRestrictModule contract.
type TransferRestrictModuleUpgradedIterator struct {
	Event *TransferRestrictModuleUpgraded // Event containing the contract specifics and raw log

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
func (it *TransferRestrictModuleUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferRestrictModuleUpgraded)
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
		it.Event = new(TransferRestrictModuleUpgraded)
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
func (it *TransferRestrictModuleUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferRestrictModuleUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferRestrictModuleUpgraded represents a Upgraded event raised by the TransferRestrictModule contract.
type TransferRestrictModuleUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*TransferRestrictModuleUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _TransferRestrictModule.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &TransferRestrictModuleUpgradedIterator{contract: _TransferRestrictModule.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *TransferRestrictModuleUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _TransferRestrictModule.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferRestrictModuleUpgraded)
				if err := _TransferRestrictModule.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) ParseUpgraded(log types.Log) (*TransferRestrictModuleUpgraded, error) {
	event := new(TransferRestrictModuleUpgraded)
	if err := _TransferRestrictModule.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferRestrictModuleUserAllowedIterator is returned from FilterUserAllowed and is used to iterate over the raw logs and unpacked data for UserAllowed events raised by the TransferRestrictModule contract.
type TransferRestrictModuleUserAllowedIterator struct {
	Event *TransferRestrictModuleUserAllowed // Event containing the contract specifics and raw log

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
func (it *TransferRestrictModuleUserAllowedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferRestrictModuleUserAllowed)
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
		it.Event = new(TransferRestrictModuleUserAllowed)
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
func (it *TransferRestrictModuleUserAllowedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferRestrictModuleUserAllowedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferRestrictModuleUserAllowed represents a UserAllowed event raised by the TransferRestrictModule contract.
type TransferRestrictModuleUserAllowed struct {
	Compliance  common.Address
	UserAddress common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUserAllowed is a free log retrieval operation binding the contract event 0x767240d3c58c1058ed268ca225eb1cc93d85cf46b6e85d8d120b39921475b1c1.
//
// Solidity: event UserAllowed(address _compliance, address _userAddress)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) FilterUserAllowed(opts *bind.FilterOpts) (*TransferRestrictModuleUserAllowedIterator, error) {

	logs, sub, err := _TransferRestrictModule.contract.FilterLogs(opts, "UserAllowed")
	if err != nil {
		return nil, err
	}
	return &TransferRestrictModuleUserAllowedIterator{contract: _TransferRestrictModule.contract, event: "UserAllowed", logs: logs, sub: sub}, nil
}

// WatchUserAllowed is a free log subscription operation binding the contract event 0x767240d3c58c1058ed268ca225eb1cc93d85cf46b6e85d8d120b39921475b1c1.
//
// Solidity: event UserAllowed(address _compliance, address _userAddress)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) WatchUserAllowed(opts *bind.WatchOpts, sink chan<- *TransferRestrictModuleUserAllowed) (event.Subscription, error) {

	logs, sub, err := _TransferRestrictModule.contract.WatchLogs(opts, "UserAllowed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferRestrictModuleUserAllowed)
				if err := _TransferRestrictModule.contract.UnpackLog(event, "UserAllowed", log); err != nil {
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

// ParseUserAllowed is a log parse operation binding the contract event 0x767240d3c58c1058ed268ca225eb1cc93d85cf46b6e85d8d120b39921475b1c1.
//
// Solidity: event UserAllowed(address _compliance, address _userAddress)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) ParseUserAllowed(log types.Log) (*TransferRestrictModuleUserAllowed, error) {
	event := new(TransferRestrictModuleUserAllowed)
	if err := _TransferRestrictModule.contract.UnpackLog(event, "UserAllowed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferRestrictModuleUserDisallowedIterator is returned from FilterUserDisallowed and is used to iterate over the raw logs and unpacked data for UserDisallowed events raised by the TransferRestrictModule contract.
type TransferRestrictModuleUserDisallowedIterator struct {
	Event *TransferRestrictModuleUserDisallowed // Event containing the contract specifics and raw log

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
func (it *TransferRestrictModuleUserDisallowedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferRestrictModuleUserDisallowed)
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
		it.Event = new(TransferRestrictModuleUserDisallowed)
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
func (it *TransferRestrictModuleUserDisallowedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferRestrictModuleUserDisallowedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferRestrictModuleUserDisallowed represents a UserDisallowed event raised by the TransferRestrictModule contract.
type TransferRestrictModuleUserDisallowed struct {
	Compliance  common.Address
	UserAddress common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUserDisallowed is a free log retrieval operation binding the contract event 0xe9ab2574b9b111f06f42baaa274030c7228b1dad45b248037101cd11de348240.
//
// Solidity: event UserDisallowed(address _compliance, address _userAddress)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) FilterUserDisallowed(opts *bind.FilterOpts) (*TransferRestrictModuleUserDisallowedIterator, error) {

	logs, sub, err := _TransferRestrictModule.contract.FilterLogs(opts, "UserDisallowed")
	if err != nil {
		return nil, err
	}
	return &TransferRestrictModuleUserDisallowedIterator{contract: _TransferRestrictModule.contract, event: "UserDisallowed", logs: logs, sub: sub}, nil
}

// WatchUserDisallowed is a free log subscription operation binding the contract event 0xe9ab2574b9b111f06f42baaa274030c7228b1dad45b248037101cd11de348240.
//
// Solidity: event UserDisallowed(address _compliance, address _userAddress)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) WatchUserDisallowed(opts *bind.WatchOpts, sink chan<- *TransferRestrictModuleUserDisallowed) (event.Subscription, error) {

	logs, sub, err := _TransferRestrictModule.contract.WatchLogs(opts, "UserDisallowed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferRestrictModuleUserDisallowed)
				if err := _TransferRestrictModule.contract.UnpackLog(event, "UserDisallowed", log); err != nil {
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

// ParseUserDisallowed is a log parse operation binding the contract event 0xe9ab2574b9b111f06f42baaa274030c7228b1dad45b248037101cd11de348240.
//
// Solidity: event UserDisallowed(address _compliance, address _userAddress)
func (_TransferRestrictModule *TransferRestrictModuleFilterer) ParseUserDisallowed(log types.Log) (*TransferRestrictModuleUserDisallowed, error) {
	event := new(TransferRestrictModuleUserDisallowed)
	if err := _TransferRestrictModule.contract.UnpackLog(event, "UserDisallowed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
