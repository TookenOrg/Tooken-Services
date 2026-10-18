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

// MaxBalanceModuleMetaData contains all meta data concerning the MaxBalanceModule contract.
var MaxBalanceModuleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_id\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_balance\",\"type\":\"uint256[]\"}],\"name\":\"InvalidPresetValues\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"MaxBalanceExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"OnlyComplianceOwnerCanCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"TokenAlreadyBound\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"ComplianceBound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"ComplianceUnbound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_id\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_balance\",\"type\":\"uint256\"}],\"name\":\"IDBalancePreSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_maxBalance\",\"type\":\"uint256\"}],\"name\":\"MaxBalanceSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"PresetCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_id\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_balance\",\"type\":\"uint256[]\"}],\"name\":\"batchPreSetModuleState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"bindCompliance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"canComplianceBind\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_identity\",\"type\":\"address\"}],\"name\":\"getIDBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"isComplianceBound\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isPlugAndPlay\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"moduleBurnAction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"moduleCheck\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"moduleMintAction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"moduleTransferAction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_id\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_balance\",\"type\":\"uint256\"}],\"name\":\"preSetModuleState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"presetCompleted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_max\",\"type\":\"uint256\"}],\"name\":\"setMaxBalance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_compliance\",\"type\":\"address\"}],\"name\":\"unbindCompliance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523060805234801561001457600080fd5b5060805161252761004c60003960008181610843015281816108c801528181610c1d01528181610ca20152610d8c01526125276000f3fe6080604052600436106101755760003560e01c80638129fc1c116100cb578063cc203f751161007f578063f104a8c911610059578063f104a8c914610431578063f2fde38b14610451578063f7c4722c1461047157600080fd5b8063cc203f75146103dd578063dd50b74c146103fd578063e6f5e8071461041d57600080fd5b80638f27c0eb116100b05780638f27c0eb1461037d5780639d51d9b71461039d578063bcc21053146103bd57600080fd5b80638129fc1c146103405780638da5cb5b1461035557600080fd5b8063372491a21161012d5780634f1ef286116101075780634f1ef286146102f557806352d1902d14610308578063715018a61461032b57600080fd5b8063372491a21461025d5780634a9325441461027d5780634cf4d2951461029d57600080fd5b806306fdde031161015e57806306fdde03146101d15780632cb7e1ec1461021d5780633659cfe61461023d57600080fd5b8063013b7ce41461017a5780630694a5fb146101af575b600080fd5b34801561018657600080fd5b5061019a610195366004612016565b610491565b60405190151581526020015b60405180910390f35b3480156101bb57600080fd5b506101cf6101ca366004612069565b610522565b005b3480156101dd57600080fd5b50604080518082018252601081527f4d617842616c616e63654d6f64756c6500000000000000000000000000000000602082015290516101a691906120aa565b34801561022957600080fd5b506101cf6102383660046120dd565b6106cc565b34801561024957600080fd5b506101cf610258366004612069565b610839565b34801561026957600080fd5b506101cf61027836600461211e565b6109b4565b34801561028957600080fd5b506101cf610298366004612069565b610a84565b3480156102a957600080fd5b5061019a6102b8366004612069565b6001600160a01b031660009081527ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae00602052604090205460ff1690565b6101cf610303366004612160565b610c13565b34801561031457600080fd5b5061031d610d7f565b6040519081526020016101a6565b34801561033757600080fd5b506101cf610e44565b34801561034c57600080fd5b506101cf610e58565b34801561036157600080fd5b506033546040516001600160a01b0390911681526020016101a6565b34801561038957600080fd5b506101cf610398366004612270565b610f77565b3480156103a957600080fd5b506101cf6103b83660046122f3565b611185565b3480156103c957600080fd5b5061019a6103d8366004612069565b611244565b3480156103e957600080fd5b5061031d6103f836600461230c565b611355565b34801561040957600080fd5b506101cf610418366004612069565b611382565b34801561042957600080fd5b50600061019a565b34801561043d57600080fd5b506101cf61044c36600461211e565b61146c565b34801561045d57600080fd5b506101cf61046c366004612069565b611591565b34801561047d57600080fd5b506101cf61048c3660046120dd565b61161e565b6001600160a01b038116600090815260ca60205260408120548311156104b95750600061051a565b60006104c5838661175d565b6001600160a01b03808516600090815260ca602090815260408083205460cb83528184209486168452939091529020549192509061050490869061235b565b111561051457600091505061051a565b60019150505b949350505050565b3360009081527ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae00602081905260409091205460ff166105a85760405162461bcd60e51b815260206004820152601e60248201527f6f6e6c7920626f756e6420636f6d706c69616e63652063616e2063616c6c000060448201526064015b60405180910390fd5b7ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae006001600160a01b03831661061f5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f206164647265737300604482015260640161059f565b336001600160a01b038416146106815760405162461bcd60e51b815260206004820152602160248201527f6f6e6c7920636f6d706c69616e636520636f6e74726163742063616e2063616c6044820152601b60fa1b606482015260840161059f565b6001600160a01b038316600081815260208390526040808220805460ff19169055517f408b49d9be1c914c52a0227e18a077e5a892dddf32a26cfa94a5d9708fad77189190a2505050565b3360009081527ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae00602081905260409091205460ff1661074d5760405162461bcd60e51b815260206004820152601e60248201527f6f6e6c7920626f756e6420636f6d706c69616e63652063616e2063616c6c0000604482015260640161059f565b6000610759338661175d565b90506000610767338661175d565b33600090815260cb602090815260408083206001600160a01b038516845290915281208054929350869290919061079f90849061235b565b909155505033600090815260cb602090815260408083206001600160a01b0386168452909152812080548692906107d790849061236e565b909155505033600090815260ca602090815260408083205460cb83528184206001600160a01b038616855290925290912054111561083157604051630d9ada5960e21b81523360048201526024810185905260440161059f565b505050505050565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036108c65760405162461bcd60e51b815260206004820152602c60248201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060448201526b19195b1959d85d1958d85b1b60a21b606482015260840161059f565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166109217f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b6001600160a01b03161461098c5760405162461bcd60e51b815260206004820152602c60248201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060448201526b6163746976652070726f787960a01b606482015260840161059f565b61099581611907565b604080516000808252602082019092526109b19183919061190f565b50565b3360009081527ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae00602081905260409091205460ff16610a355760405162461bcd60e51b815260206004820152601e60248201527f6f6e6c7920626f756e6420636f6d706c69616e63652063616e2063616c6c0000604482015260640161059f565b6000610a41338561175d565b33600090815260cb602090815260408083206001600160a01b0385168452909152812080549293508592909190610a7990849061236e565b909155505050505050565b7ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae006001600160a01b038216610afb5760405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420617267756d656e74202d207a65726f206164647265737300604482015260640161059f565b6001600160a01b03821660009081526020829052604090205460ff1615610b645760405162461bcd60e51b815260206004820152601860248201527f636f6d706c69616e636520616c726561647920626f756e640000000000000000604482015260640161059f565b336001600160a01b03831614610bc65760405162461bcd60e51b815260206004820152602160248201527f6f6e6c7920636f6d706c69616e636520636f6e74726163742063616e2063616c6044820152601b60fa1b606482015260840161059f565b6001600160a01b038216600081815260208390526040808220805460ff19166001179055517f1f7b76c58fb697eb53c6c7c1becb96911516a136e24d7ced386b2355358b75a39190a25050565b6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000163003610ca05760405162461bcd60e51b815260206004820152602c60248201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060448201526b19195b1959d85d1958d85b1b60a21b606482015260840161059f565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610cfb7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b6001600160a01b031614610d665760405162461bcd60e51b815260206004820152602c60248201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060448201526b6163746976652070726f787960a01b606482015260840161059f565b610d6f82611907565b610d7b8282600161190f565b5050565b6000306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614610e1f5760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c0000000000000000606482015260840161059f565b507f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc90565b610e4c611aaf565b610e566000611b09565b565b600054610100900460ff1615808015610e785750600054600160ff909116105b80610e925750303b158015610e92575060005460ff166001145b610f045760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a6564000000000000000000000000000000000000606482015260840161059f565b6000805460ff191660011790558015610f27576000805461ff0019166101001790555b610f2f611b68565b80156109b1576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a150565b821580610f845750828114155b15610fc55784848484846040517f6ab857fa00000000000000000000000000000000000000000000000000000000815260040161059f959493929190612381565b336001600160a01b0316856001600160a01b0316638da5cb5b6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561100d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110319190612427565b6001600160a01b0316146110635760405163e8649dab60e01b81526001600160a01b038616600482015260240161059f565b60405163a446d49f60e01b81523060048201526001600160a01b0386169063a446d49f90602401602060405180830381865afa1580156110a7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110cb9190612444565b156110f45760405163346a0c0160e21b81526001600160a01b038616600482015260240161059f565b60005b8381101561115a576111488686868481811061111557611115612466565b905060200201602081019061112a9190612069565b85858581811061113c5761113c612466565b90506020020135611be3565b806111528161247c565b9150506110f7565b5050506001600160a01b03909216600090815260c960205260409020805460ff191660011790555050565b3360009081527ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae00602081905260409091205460ff166112065760405162461bcd60e51b815260206004820152601e60248201527f6f6e6c7920626f756e6420636f6d706c69616e63652063616e2063616c6c0000604482015260640161059f565b33600081815260ca6020526040808220859055518492917f8bec153ee727d17f3f63aa27cbb6bd64eb978b88821b30002a4b5517aeb10cec91a35050565b6001600160a01b038116600090815260c9602052604081205460ff161561126d57506001919050565b6000826001600160a01b0316636a3edf286040518163ffffffff1660e01b8152600401602060405180830381865afa1580156112ad573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112d19190612427565b90506000816001600160a01b03166318160ddd6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611313573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113379190612495565b90508060000361134b575060019392505050565b5060009392505050565b6001600160a01b03808316600090815260cb60209081526040808320938516835292905220545b92915050565b336001600160a01b0316816001600160a01b0316638da5cb5b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156113ca573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113ee9190612427565b6001600160a01b0316146114205760405163e8649dab60e01b81526001600160a01b038216600482015260240161059f565b6001600160a01b038116600081815260c96020526040808220805460ff19166001179055517f40bfb91138c8514efd65a280ed9f436317a8e93584e42ef52b3b863e68bd500f9190a250565b3360009081527ff6cc97de1266c180cd39f3b311632644143ce7873d2927755382ad4b39e8ae00602081905260409091205460ff166114ed5760405162461bcd60e51b815260206004820152601e60248201527f6f6e6c7920626f756e6420636f6d706c69616e63652063616e2063616c6c0000604482015260640161059f565b60006114f9338561175d565b33600090815260cb602090815260408083206001600160a01b038516845290915281208054929350859290919061153190849061235b565b909155505033600090815260ca602090815260408083205460cb83528184206001600160a01b038616855290925290912054111561158b57604051630d9ada5960e21b81523360048201526024810184905260440161059f565b50505050565b611599611aaf565b6001600160a01b0381166116155760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f6464726573730000000000000000000000000000000000000000000000000000606482015260840161059f565b6109b181611b09565b336001600160a01b0316836001600160a01b0316638da5cb5b6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611666573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061168a9190612427565b6001600160a01b0316146116bc5760405163e8649dab60e01b81526001600160a01b038416600482015260240161059f565b60405163a446d49f60e01b81523060048201526001600160a01b0384169063a446d49f90602401602060405180830381865afa158015611700573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117249190612444565b1561174d5760405163346a0c0160e21b81526001600160a01b038416600482015260240161059f565b611758838383611be3565b505050565b600080836001600160a01b0316636a3edf286040518163ffffffff1660e01b8152600401602060405180830381865afa15801561179e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117c29190612427565b6001600160a01b031663134e18f46040518163ffffffff1660e01b8152600401602060405180830381865afa1580156117ff573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118239190612427565b6040517ff0eb5e540000000000000000000000000000000000000000000000000000000081526001600160a01b038581166004830152919091169063f0eb5e5490602401602060405180830381865afa158015611884573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118a89190612427565b90506001600160a01b0381166119005760405162461bcd60e51b815260206004820152601260248201527f6964656e74697479206e6f7420666f756e640000000000000000000000000000604482015260640161059f565b9392505050565b6109b1611aaf565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff16156119425761175883611c44565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa92505050801561199c575060408051601f3d908101601f1916820190925261199991810190612495565b60015b611a0e5760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201527f6f6e206973206e6f742055555053000000000000000000000000000000000000606482015260840161059f565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8114611aa35760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f7860448201527f6961626c65555549440000000000000000000000000000000000000000000000606482015260840161059f565b50611758838383611d0f565b6033546001600160a01b03163314610e565760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015260640161059f565b603380546001600160a01b0383811673ffffffffffffffffffffffffffffffffffffffff19831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff16611bd35760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b606482015260840161059f565b611bdb611d34565b610e56611da7565b6001600160a01b03838116600081815260cb602090815260408083209487168084529482529182902085905590518481527f0a777417e1c0a4b24c5f556b0c3dd522c37b99b2c3ca4bf9fea2fa1293c660c9910160405180910390a3505050565b6001600160a01b0381163b611cc15760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201527f6f74206120636f6e747261637400000000000000000000000000000000000000606482015260840161059f565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0392909216919091179055565b611d1883611e12565b600082511180611d255750805b156117585761158b8383611e52565b600054610100900460ff16611d9f5760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b606482015260840161059f565b610e56611e77565b600054610100900460ff16610e565760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b606482015260840161059f565b611e1b81611c44565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b606061190083836040518060600160405280602781526020016124cb60279139611eeb565b600054610100900460ff16611ee25760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b606482015260840161059f565b610e5633611b09565b6060600080856001600160a01b031685604051611f0891906124ae565b600060405180830381855af49150503d8060008114611f43576040519150601f19603f3d011682016040523d82523d6000602084013e611f48565b606091505b5091509150611f5986838387611f63565b9695505050505050565b60608315611fd2578251600003611fcb576001600160a01b0385163b611fcb5760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640161059f565b508161051a565b61051a8383815115611fe75781518083602001fd5b8060405162461bcd60e51b815260040161059f91906120aa565b6001600160a01b03811681146109b157600080fd5b6000806000806080858703121561202c57600080fd5b843561203781612001565b9350602085013561204781612001565b925060408501359150606085013561205e81612001565b939692955090935050565b60006020828403121561207b57600080fd5b813561190081612001565b60005b838110156120a1578181015183820152602001612089565b50506000910152565b60208152600082518060208401526120c9816040850160208701612086565b601f01601f19169190910160400192915050565b6000806000606084860312156120f257600080fd5b83356120fd81612001565b9250602084013561210d81612001565b929592945050506040919091013590565b6000806040838503121561213157600080fd5b823561213c81612001565b946020939093013593505050565b634e487b7160e01b600052604160045260246000fd5b6000806040838503121561217357600080fd5b823561217e81612001565b9150602083013567ffffffffffffffff8082111561219b57600080fd5b818501915085601f8301126121af57600080fd5b8135818111156121c1576121c161214a565b604051601f8201601f19908116603f011681019083821181831017156121e9576121e961214a565b8160405282815288602084870101111561220257600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b60008083601f84011261223657600080fd5b50813567ffffffffffffffff81111561224e57600080fd5b6020830191508360208260051b850101111561226957600080fd5b9250929050565b60008060008060006060868803121561228857600080fd5b853561229381612001565b9450602086013567ffffffffffffffff808211156122b057600080fd5b6122bc89838a01612224565b909650945060408801359150808211156122d557600080fd5b506122e288828901612224565b969995985093965092949392505050565b60006020828403121561230557600080fd5b5035919050565b6000806040838503121561231f57600080fd5b823561232a81612001565b9150602083013561233a81612001565b809150509250929050565b634e487b7160e01b600052601160045260246000fd5b8082018082111561137c5761137c612345565b8181038181111561137c5761137c612345565b6001600160a01b03868116825260606020808401829052908301869052600091879160808501845b898110156123d05784356123bc81612001565b8416825293820193908201906001016123a9565b5085810360408701528681527f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff87111561240957600080fd5b8660051b935083888383013790920190910198975050505050505050565b60006020828403121561243957600080fd5b815161190081612001565b60006020828403121561245657600080fd5b8151801515811461190057600080fd5b634e487b7160e01b600052603260045260246000fd5b60006001820161248e5761248e612345565b5060010190565b6000602082840312156124a757600080fd5b5051919050565b600082516124c0818460208701612086565b919091019291505056fe416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a2646970667358221220b96c2b82660a7e91093fdbf6985024cbd9ef98b54cc1fe27f034443c57fff1e464736f6c63430008110033",
}

// MaxBalanceModuleABI is the input ABI used to generate the binding from.
// Deprecated: Use MaxBalanceModuleMetaData.ABI instead.
var MaxBalanceModuleABI = MaxBalanceModuleMetaData.ABI

// MaxBalanceModuleBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MaxBalanceModuleMetaData.Bin instead.
var MaxBalanceModuleBin = MaxBalanceModuleMetaData.Bin

// DeployMaxBalanceModule deploys a new Ethereum contract, binding an instance of MaxBalanceModule to it.
func DeployMaxBalanceModule(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MaxBalanceModule, error) {
	parsed, err := MaxBalanceModuleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MaxBalanceModuleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MaxBalanceModule{MaxBalanceModuleCaller: MaxBalanceModuleCaller{contract: contract}, MaxBalanceModuleTransactor: MaxBalanceModuleTransactor{contract: contract}, MaxBalanceModuleFilterer: MaxBalanceModuleFilterer{contract: contract}}, nil
}

// MaxBalanceModule is an auto generated Go binding around an Ethereum contract.
type MaxBalanceModule struct {
	MaxBalanceModuleCaller     // Read-only binding to the contract
	MaxBalanceModuleTransactor // Write-only binding to the contract
	MaxBalanceModuleFilterer   // Log filterer for contract events
}

// MaxBalanceModuleCaller is an auto generated read-only Go binding around an Ethereum contract.
type MaxBalanceModuleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MaxBalanceModuleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MaxBalanceModuleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MaxBalanceModuleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MaxBalanceModuleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MaxBalanceModuleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MaxBalanceModuleSession struct {
	Contract     *MaxBalanceModule // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MaxBalanceModuleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MaxBalanceModuleCallerSession struct {
	Contract *MaxBalanceModuleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// MaxBalanceModuleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MaxBalanceModuleTransactorSession struct {
	Contract     *MaxBalanceModuleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// MaxBalanceModuleRaw is an auto generated low-level Go binding around an Ethereum contract.
type MaxBalanceModuleRaw struct {
	Contract *MaxBalanceModule // Generic contract binding to access the raw methods on
}

// MaxBalanceModuleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MaxBalanceModuleCallerRaw struct {
	Contract *MaxBalanceModuleCaller // Generic read-only contract binding to access the raw methods on
}

// MaxBalanceModuleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MaxBalanceModuleTransactorRaw struct {
	Contract *MaxBalanceModuleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMaxBalanceModule creates a new instance of MaxBalanceModule, bound to a specific deployed contract.
func NewMaxBalanceModule(address common.Address, backend bind.ContractBackend) (*MaxBalanceModule, error) {
	contract, err := bindMaxBalanceModule(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MaxBalanceModule{MaxBalanceModuleCaller: MaxBalanceModuleCaller{contract: contract}, MaxBalanceModuleTransactor: MaxBalanceModuleTransactor{contract: contract}, MaxBalanceModuleFilterer: MaxBalanceModuleFilterer{contract: contract}}, nil
}

// NewMaxBalanceModuleCaller creates a new read-only instance of MaxBalanceModule, bound to a specific deployed contract.
func NewMaxBalanceModuleCaller(address common.Address, caller bind.ContractCaller) (*MaxBalanceModuleCaller, error) {
	contract, err := bindMaxBalanceModule(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MaxBalanceModuleCaller{contract: contract}, nil
}

// NewMaxBalanceModuleTransactor creates a new write-only instance of MaxBalanceModule, bound to a specific deployed contract.
func NewMaxBalanceModuleTransactor(address common.Address, transactor bind.ContractTransactor) (*MaxBalanceModuleTransactor, error) {
	contract, err := bindMaxBalanceModule(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MaxBalanceModuleTransactor{contract: contract}, nil
}

// NewMaxBalanceModuleFilterer creates a new log filterer instance of MaxBalanceModule, bound to a specific deployed contract.
func NewMaxBalanceModuleFilterer(address common.Address, filterer bind.ContractFilterer) (*MaxBalanceModuleFilterer, error) {
	contract, err := bindMaxBalanceModule(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MaxBalanceModuleFilterer{contract: contract}, nil
}

// bindMaxBalanceModule binds a generic wrapper to an already deployed contract.
func bindMaxBalanceModule(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MaxBalanceModuleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MaxBalanceModule *MaxBalanceModuleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MaxBalanceModule.Contract.MaxBalanceModuleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MaxBalanceModule *MaxBalanceModuleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.MaxBalanceModuleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MaxBalanceModule *MaxBalanceModuleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.MaxBalanceModuleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MaxBalanceModule *MaxBalanceModuleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MaxBalanceModule.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MaxBalanceModule *MaxBalanceModuleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MaxBalanceModule *MaxBalanceModuleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.contract.Transact(opts, method, params...)
}

// CanComplianceBind is a free data retrieval call binding the contract method 0xbcc21053.
//
// Solidity: function canComplianceBind(address _compliance) view returns(bool)
func (_MaxBalanceModule *MaxBalanceModuleCaller) CanComplianceBind(opts *bind.CallOpts, _compliance common.Address) (bool, error) {
	var out []interface{}
	err := _MaxBalanceModule.contract.Call(opts, &out, "canComplianceBind", _compliance)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanComplianceBind is a free data retrieval call binding the contract method 0xbcc21053.
//
// Solidity: function canComplianceBind(address _compliance) view returns(bool)
func (_MaxBalanceModule *MaxBalanceModuleSession) CanComplianceBind(_compliance common.Address) (bool, error) {
	return _MaxBalanceModule.Contract.CanComplianceBind(&_MaxBalanceModule.CallOpts, _compliance)
}

// CanComplianceBind is a free data retrieval call binding the contract method 0xbcc21053.
//
// Solidity: function canComplianceBind(address _compliance) view returns(bool)
func (_MaxBalanceModule *MaxBalanceModuleCallerSession) CanComplianceBind(_compliance common.Address) (bool, error) {
	return _MaxBalanceModule.Contract.CanComplianceBind(&_MaxBalanceModule.CallOpts, _compliance)
}

// GetIDBalance is a free data retrieval call binding the contract method 0xcc203f75.
//
// Solidity: function getIDBalance(address _compliance, address _identity) view returns(uint256)
func (_MaxBalanceModule *MaxBalanceModuleCaller) GetIDBalance(opts *bind.CallOpts, _compliance common.Address, _identity common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MaxBalanceModule.contract.Call(opts, &out, "getIDBalance", _compliance, _identity)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetIDBalance is a free data retrieval call binding the contract method 0xcc203f75.
//
// Solidity: function getIDBalance(address _compliance, address _identity) view returns(uint256)
func (_MaxBalanceModule *MaxBalanceModuleSession) GetIDBalance(_compliance common.Address, _identity common.Address) (*big.Int, error) {
	return _MaxBalanceModule.Contract.GetIDBalance(&_MaxBalanceModule.CallOpts, _compliance, _identity)
}

// GetIDBalance is a free data retrieval call binding the contract method 0xcc203f75.
//
// Solidity: function getIDBalance(address _compliance, address _identity) view returns(uint256)
func (_MaxBalanceModule *MaxBalanceModuleCallerSession) GetIDBalance(_compliance common.Address, _identity common.Address) (*big.Int, error) {
	return _MaxBalanceModule.Contract.GetIDBalance(&_MaxBalanceModule.CallOpts, _compliance, _identity)
}

// IsComplianceBound is a free data retrieval call binding the contract method 0x4cf4d295.
//
// Solidity: function isComplianceBound(address _compliance) view returns(bool)
func (_MaxBalanceModule *MaxBalanceModuleCaller) IsComplianceBound(opts *bind.CallOpts, _compliance common.Address) (bool, error) {
	var out []interface{}
	err := _MaxBalanceModule.contract.Call(opts, &out, "isComplianceBound", _compliance)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsComplianceBound is a free data retrieval call binding the contract method 0x4cf4d295.
//
// Solidity: function isComplianceBound(address _compliance) view returns(bool)
func (_MaxBalanceModule *MaxBalanceModuleSession) IsComplianceBound(_compliance common.Address) (bool, error) {
	return _MaxBalanceModule.Contract.IsComplianceBound(&_MaxBalanceModule.CallOpts, _compliance)
}

// IsComplianceBound is a free data retrieval call binding the contract method 0x4cf4d295.
//
// Solidity: function isComplianceBound(address _compliance) view returns(bool)
func (_MaxBalanceModule *MaxBalanceModuleCallerSession) IsComplianceBound(_compliance common.Address) (bool, error) {
	return _MaxBalanceModule.Contract.IsComplianceBound(&_MaxBalanceModule.CallOpts, _compliance)
}

// IsPlugAndPlay is a free data retrieval call binding the contract method 0xe6f5e807.
//
// Solidity: function isPlugAndPlay() pure returns(bool)
func (_MaxBalanceModule *MaxBalanceModuleCaller) IsPlugAndPlay(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MaxBalanceModule.contract.Call(opts, &out, "isPlugAndPlay")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPlugAndPlay is a free data retrieval call binding the contract method 0xe6f5e807.
//
// Solidity: function isPlugAndPlay() pure returns(bool)
func (_MaxBalanceModule *MaxBalanceModuleSession) IsPlugAndPlay() (bool, error) {
	return _MaxBalanceModule.Contract.IsPlugAndPlay(&_MaxBalanceModule.CallOpts)
}

// IsPlugAndPlay is a free data retrieval call binding the contract method 0xe6f5e807.
//
// Solidity: function isPlugAndPlay() pure returns(bool)
func (_MaxBalanceModule *MaxBalanceModuleCallerSession) IsPlugAndPlay() (bool, error) {
	return _MaxBalanceModule.Contract.IsPlugAndPlay(&_MaxBalanceModule.CallOpts)
}

// ModuleCheck is a free data retrieval call binding the contract method 0x013b7ce4.
//
// Solidity: function moduleCheck(address , address _to, uint256 _value, address _compliance) view returns(bool)
func (_MaxBalanceModule *MaxBalanceModuleCaller) ModuleCheck(opts *bind.CallOpts, arg0 common.Address, _to common.Address, _value *big.Int, _compliance common.Address) (bool, error) {
	var out []interface{}
	err := _MaxBalanceModule.contract.Call(opts, &out, "moduleCheck", arg0, _to, _value, _compliance)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ModuleCheck is a free data retrieval call binding the contract method 0x013b7ce4.
//
// Solidity: function moduleCheck(address , address _to, uint256 _value, address _compliance) view returns(bool)
func (_MaxBalanceModule *MaxBalanceModuleSession) ModuleCheck(arg0 common.Address, _to common.Address, _value *big.Int, _compliance common.Address) (bool, error) {
	return _MaxBalanceModule.Contract.ModuleCheck(&_MaxBalanceModule.CallOpts, arg0, _to, _value, _compliance)
}

// ModuleCheck is a free data retrieval call binding the contract method 0x013b7ce4.
//
// Solidity: function moduleCheck(address , address _to, uint256 _value, address _compliance) view returns(bool)
func (_MaxBalanceModule *MaxBalanceModuleCallerSession) ModuleCheck(arg0 common.Address, _to common.Address, _value *big.Int, _compliance common.Address) (bool, error) {
	return _MaxBalanceModule.Contract.ModuleCheck(&_MaxBalanceModule.CallOpts, arg0, _to, _value, _compliance)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string _name)
func (_MaxBalanceModule *MaxBalanceModuleCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MaxBalanceModule.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string _name)
func (_MaxBalanceModule *MaxBalanceModuleSession) Name() (string, error) {
	return _MaxBalanceModule.Contract.Name(&_MaxBalanceModule.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string _name)
func (_MaxBalanceModule *MaxBalanceModuleCallerSession) Name() (string, error) {
	return _MaxBalanceModule.Contract.Name(&_MaxBalanceModule.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MaxBalanceModule *MaxBalanceModuleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MaxBalanceModule.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MaxBalanceModule *MaxBalanceModuleSession) Owner() (common.Address, error) {
	return _MaxBalanceModule.Contract.Owner(&_MaxBalanceModule.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MaxBalanceModule *MaxBalanceModuleCallerSession) Owner() (common.Address, error) {
	return _MaxBalanceModule.Contract.Owner(&_MaxBalanceModule.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_MaxBalanceModule *MaxBalanceModuleCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _MaxBalanceModule.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_MaxBalanceModule *MaxBalanceModuleSession) ProxiableUUID() ([32]byte, error) {
	return _MaxBalanceModule.Contract.ProxiableUUID(&_MaxBalanceModule.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_MaxBalanceModule *MaxBalanceModuleCallerSession) ProxiableUUID() ([32]byte, error) {
	return _MaxBalanceModule.Contract.ProxiableUUID(&_MaxBalanceModule.CallOpts)
}

// BatchPreSetModuleState is a paid mutator transaction binding the contract method 0x8f27c0eb.
//
// Solidity: function batchPreSetModuleState(address _compliance, address[] _id, uint256[] _balance) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactor) BatchPreSetModuleState(opts *bind.TransactOpts, _compliance common.Address, _id []common.Address, _balance []*big.Int) (*types.Transaction, error) {
	return _MaxBalanceModule.contract.Transact(opts, "batchPreSetModuleState", _compliance, _id, _balance)
}

// BatchPreSetModuleState is a paid mutator transaction binding the contract method 0x8f27c0eb.
//
// Solidity: function batchPreSetModuleState(address _compliance, address[] _id, uint256[] _balance) returns()
func (_MaxBalanceModule *MaxBalanceModuleSession) BatchPreSetModuleState(_compliance common.Address, _id []common.Address, _balance []*big.Int) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.BatchPreSetModuleState(&_MaxBalanceModule.TransactOpts, _compliance, _id, _balance)
}

// BatchPreSetModuleState is a paid mutator transaction binding the contract method 0x8f27c0eb.
//
// Solidity: function batchPreSetModuleState(address _compliance, address[] _id, uint256[] _balance) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactorSession) BatchPreSetModuleState(_compliance common.Address, _id []common.Address, _balance []*big.Int) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.BatchPreSetModuleState(&_MaxBalanceModule.TransactOpts, _compliance, _id, _balance)
}

// BindCompliance is a paid mutator transaction binding the contract method 0x4a932544.
//
// Solidity: function bindCompliance(address _compliance) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactor) BindCompliance(opts *bind.TransactOpts, _compliance common.Address) (*types.Transaction, error) {
	return _MaxBalanceModule.contract.Transact(opts, "bindCompliance", _compliance)
}

// BindCompliance is a paid mutator transaction binding the contract method 0x4a932544.
//
// Solidity: function bindCompliance(address _compliance) returns()
func (_MaxBalanceModule *MaxBalanceModuleSession) BindCompliance(_compliance common.Address) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.BindCompliance(&_MaxBalanceModule.TransactOpts, _compliance)
}

// BindCompliance is a paid mutator transaction binding the contract method 0x4a932544.
//
// Solidity: function bindCompliance(address _compliance) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactorSession) BindCompliance(_compliance common.Address) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.BindCompliance(&_MaxBalanceModule.TransactOpts, _compliance)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MaxBalanceModule.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_MaxBalanceModule *MaxBalanceModuleSession) Initialize() (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.Initialize(&_MaxBalanceModule.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactorSession) Initialize() (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.Initialize(&_MaxBalanceModule.TransactOpts)
}

// ModuleBurnAction is a paid mutator transaction binding the contract method 0x372491a2.
//
// Solidity: function moduleBurnAction(address _from, uint256 _value) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactor) ModuleBurnAction(opts *bind.TransactOpts, _from common.Address, _value *big.Int) (*types.Transaction, error) {
	return _MaxBalanceModule.contract.Transact(opts, "moduleBurnAction", _from, _value)
}

// ModuleBurnAction is a paid mutator transaction binding the contract method 0x372491a2.
//
// Solidity: function moduleBurnAction(address _from, uint256 _value) returns()
func (_MaxBalanceModule *MaxBalanceModuleSession) ModuleBurnAction(_from common.Address, _value *big.Int) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.ModuleBurnAction(&_MaxBalanceModule.TransactOpts, _from, _value)
}

// ModuleBurnAction is a paid mutator transaction binding the contract method 0x372491a2.
//
// Solidity: function moduleBurnAction(address _from, uint256 _value) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactorSession) ModuleBurnAction(_from common.Address, _value *big.Int) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.ModuleBurnAction(&_MaxBalanceModule.TransactOpts, _from, _value)
}

// ModuleMintAction is a paid mutator transaction binding the contract method 0xf104a8c9.
//
// Solidity: function moduleMintAction(address _to, uint256 _value) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactor) ModuleMintAction(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _MaxBalanceModule.contract.Transact(opts, "moduleMintAction", _to, _value)
}

// ModuleMintAction is a paid mutator transaction binding the contract method 0xf104a8c9.
//
// Solidity: function moduleMintAction(address _to, uint256 _value) returns()
func (_MaxBalanceModule *MaxBalanceModuleSession) ModuleMintAction(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.ModuleMintAction(&_MaxBalanceModule.TransactOpts, _to, _value)
}

// ModuleMintAction is a paid mutator transaction binding the contract method 0xf104a8c9.
//
// Solidity: function moduleMintAction(address _to, uint256 _value) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactorSession) ModuleMintAction(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.ModuleMintAction(&_MaxBalanceModule.TransactOpts, _to, _value)
}

// ModuleTransferAction is a paid mutator transaction binding the contract method 0x2cb7e1ec.
//
// Solidity: function moduleTransferAction(address _from, address _to, uint256 _value) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactor) ModuleTransferAction(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _MaxBalanceModule.contract.Transact(opts, "moduleTransferAction", _from, _to, _value)
}

// ModuleTransferAction is a paid mutator transaction binding the contract method 0x2cb7e1ec.
//
// Solidity: function moduleTransferAction(address _from, address _to, uint256 _value) returns()
func (_MaxBalanceModule *MaxBalanceModuleSession) ModuleTransferAction(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.ModuleTransferAction(&_MaxBalanceModule.TransactOpts, _from, _to, _value)
}

// ModuleTransferAction is a paid mutator transaction binding the contract method 0x2cb7e1ec.
//
// Solidity: function moduleTransferAction(address _from, address _to, uint256 _value) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactorSession) ModuleTransferAction(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.ModuleTransferAction(&_MaxBalanceModule.TransactOpts, _from, _to, _value)
}

// PreSetModuleState is a paid mutator transaction binding the contract method 0xf7c4722c.
//
// Solidity: function preSetModuleState(address _compliance, address _id, uint256 _balance) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactor) PreSetModuleState(opts *bind.TransactOpts, _compliance common.Address, _id common.Address, _balance *big.Int) (*types.Transaction, error) {
	return _MaxBalanceModule.contract.Transact(opts, "preSetModuleState", _compliance, _id, _balance)
}

// PreSetModuleState is a paid mutator transaction binding the contract method 0xf7c4722c.
//
// Solidity: function preSetModuleState(address _compliance, address _id, uint256 _balance) returns()
func (_MaxBalanceModule *MaxBalanceModuleSession) PreSetModuleState(_compliance common.Address, _id common.Address, _balance *big.Int) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.PreSetModuleState(&_MaxBalanceModule.TransactOpts, _compliance, _id, _balance)
}

// PreSetModuleState is a paid mutator transaction binding the contract method 0xf7c4722c.
//
// Solidity: function preSetModuleState(address _compliance, address _id, uint256 _balance) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactorSession) PreSetModuleState(_compliance common.Address, _id common.Address, _balance *big.Int) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.PreSetModuleState(&_MaxBalanceModule.TransactOpts, _compliance, _id, _balance)
}

// PresetCompleted is a paid mutator transaction binding the contract method 0xdd50b74c.
//
// Solidity: function presetCompleted(address _compliance) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactor) PresetCompleted(opts *bind.TransactOpts, _compliance common.Address) (*types.Transaction, error) {
	return _MaxBalanceModule.contract.Transact(opts, "presetCompleted", _compliance)
}

// PresetCompleted is a paid mutator transaction binding the contract method 0xdd50b74c.
//
// Solidity: function presetCompleted(address _compliance) returns()
func (_MaxBalanceModule *MaxBalanceModuleSession) PresetCompleted(_compliance common.Address) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.PresetCompleted(&_MaxBalanceModule.TransactOpts, _compliance)
}

// PresetCompleted is a paid mutator transaction binding the contract method 0xdd50b74c.
//
// Solidity: function presetCompleted(address _compliance) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactorSession) PresetCompleted(_compliance common.Address) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.PresetCompleted(&_MaxBalanceModule.TransactOpts, _compliance)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MaxBalanceModule.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MaxBalanceModule *MaxBalanceModuleSession) RenounceOwnership() (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.RenounceOwnership(&_MaxBalanceModule.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.RenounceOwnership(&_MaxBalanceModule.TransactOpts)
}

// SetMaxBalance is a paid mutator transaction binding the contract method 0x9d51d9b7.
//
// Solidity: function setMaxBalance(uint256 _max) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactor) SetMaxBalance(opts *bind.TransactOpts, _max *big.Int) (*types.Transaction, error) {
	return _MaxBalanceModule.contract.Transact(opts, "setMaxBalance", _max)
}

// SetMaxBalance is a paid mutator transaction binding the contract method 0x9d51d9b7.
//
// Solidity: function setMaxBalance(uint256 _max) returns()
func (_MaxBalanceModule *MaxBalanceModuleSession) SetMaxBalance(_max *big.Int) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.SetMaxBalance(&_MaxBalanceModule.TransactOpts, _max)
}

// SetMaxBalance is a paid mutator transaction binding the contract method 0x9d51d9b7.
//
// Solidity: function setMaxBalance(uint256 _max) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactorSession) SetMaxBalance(_max *big.Int) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.SetMaxBalance(&_MaxBalanceModule.TransactOpts, _max)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _MaxBalanceModule.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MaxBalanceModule *MaxBalanceModuleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.TransferOwnership(&_MaxBalanceModule.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.TransferOwnership(&_MaxBalanceModule.TransactOpts, newOwner)
}

// UnbindCompliance is a paid mutator transaction binding the contract method 0x0694a5fb.
//
// Solidity: function unbindCompliance(address _compliance) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactor) UnbindCompliance(opts *bind.TransactOpts, _compliance common.Address) (*types.Transaction, error) {
	return _MaxBalanceModule.contract.Transact(opts, "unbindCompliance", _compliance)
}

// UnbindCompliance is a paid mutator transaction binding the contract method 0x0694a5fb.
//
// Solidity: function unbindCompliance(address _compliance) returns()
func (_MaxBalanceModule *MaxBalanceModuleSession) UnbindCompliance(_compliance common.Address) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.UnbindCompliance(&_MaxBalanceModule.TransactOpts, _compliance)
}

// UnbindCompliance is a paid mutator transaction binding the contract method 0x0694a5fb.
//
// Solidity: function unbindCompliance(address _compliance) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactorSession) UnbindCompliance(_compliance common.Address) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.UnbindCompliance(&_MaxBalanceModule.TransactOpts, _compliance)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _MaxBalanceModule.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_MaxBalanceModule *MaxBalanceModuleSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.UpgradeTo(&_MaxBalanceModule.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.UpgradeTo(&_MaxBalanceModule.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _MaxBalanceModule.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_MaxBalanceModule *MaxBalanceModuleSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.UpgradeToAndCall(&_MaxBalanceModule.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_MaxBalanceModule *MaxBalanceModuleTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _MaxBalanceModule.Contract.UpgradeToAndCall(&_MaxBalanceModule.TransactOpts, newImplementation, data)
}

// MaxBalanceModuleAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the MaxBalanceModule contract.
type MaxBalanceModuleAdminChangedIterator struct {
	Event *MaxBalanceModuleAdminChanged // Event containing the contract specifics and raw log

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
func (it *MaxBalanceModuleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MaxBalanceModuleAdminChanged)
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
		it.Event = new(MaxBalanceModuleAdminChanged)
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
func (it *MaxBalanceModuleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MaxBalanceModuleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MaxBalanceModuleAdminChanged represents a AdminChanged event raised by the MaxBalanceModule contract.
type MaxBalanceModuleAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*MaxBalanceModuleAdminChangedIterator, error) {

	logs, sub, err := _MaxBalanceModule.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &MaxBalanceModuleAdminChangedIterator{contract: _MaxBalanceModule.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *MaxBalanceModuleAdminChanged) (event.Subscription, error) {

	logs, sub, err := _MaxBalanceModule.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MaxBalanceModuleAdminChanged)
				if err := _MaxBalanceModule.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_MaxBalanceModule *MaxBalanceModuleFilterer) ParseAdminChanged(log types.Log) (*MaxBalanceModuleAdminChanged, error) {
	event := new(MaxBalanceModuleAdminChanged)
	if err := _MaxBalanceModule.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MaxBalanceModuleBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the MaxBalanceModule contract.
type MaxBalanceModuleBeaconUpgradedIterator struct {
	Event *MaxBalanceModuleBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *MaxBalanceModuleBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MaxBalanceModuleBeaconUpgraded)
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
		it.Event = new(MaxBalanceModuleBeaconUpgraded)
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
func (it *MaxBalanceModuleBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MaxBalanceModuleBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MaxBalanceModuleBeaconUpgraded represents a BeaconUpgraded event raised by the MaxBalanceModule contract.
type MaxBalanceModuleBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*MaxBalanceModuleBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _MaxBalanceModule.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &MaxBalanceModuleBeaconUpgradedIterator{contract: _MaxBalanceModule.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *MaxBalanceModuleBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _MaxBalanceModule.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MaxBalanceModuleBeaconUpgraded)
				if err := _MaxBalanceModule.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_MaxBalanceModule *MaxBalanceModuleFilterer) ParseBeaconUpgraded(log types.Log) (*MaxBalanceModuleBeaconUpgraded, error) {
	event := new(MaxBalanceModuleBeaconUpgraded)
	if err := _MaxBalanceModule.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MaxBalanceModuleComplianceBoundIterator is returned from FilterComplianceBound and is used to iterate over the raw logs and unpacked data for ComplianceBound events raised by the MaxBalanceModule contract.
type MaxBalanceModuleComplianceBoundIterator struct {
	Event *MaxBalanceModuleComplianceBound // Event containing the contract specifics and raw log

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
func (it *MaxBalanceModuleComplianceBoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MaxBalanceModuleComplianceBound)
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
		it.Event = new(MaxBalanceModuleComplianceBound)
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
func (it *MaxBalanceModuleComplianceBoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MaxBalanceModuleComplianceBoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MaxBalanceModuleComplianceBound represents a ComplianceBound event raised by the MaxBalanceModule contract.
type MaxBalanceModuleComplianceBound struct {
	Compliance common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterComplianceBound is a free log retrieval operation binding the contract event 0x1f7b76c58fb697eb53c6c7c1becb96911516a136e24d7ced386b2355358b75a3.
//
// Solidity: event ComplianceBound(address indexed _compliance)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) FilterComplianceBound(opts *bind.FilterOpts, _compliance []common.Address) (*MaxBalanceModuleComplianceBoundIterator, error) {

	var _complianceRule []interface{}
	for _, _complianceItem := range _compliance {
		_complianceRule = append(_complianceRule, _complianceItem)
	}

	logs, sub, err := _MaxBalanceModule.contract.FilterLogs(opts, "ComplianceBound", _complianceRule)
	if err != nil {
		return nil, err
	}
	return &MaxBalanceModuleComplianceBoundIterator{contract: _MaxBalanceModule.contract, event: "ComplianceBound", logs: logs, sub: sub}, nil
}

// WatchComplianceBound is a free log subscription operation binding the contract event 0x1f7b76c58fb697eb53c6c7c1becb96911516a136e24d7ced386b2355358b75a3.
//
// Solidity: event ComplianceBound(address indexed _compliance)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) WatchComplianceBound(opts *bind.WatchOpts, sink chan<- *MaxBalanceModuleComplianceBound, _compliance []common.Address) (event.Subscription, error) {

	var _complianceRule []interface{}
	for _, _complianceItem := range _compliance {
		_complianceRule = append(_complianceRule, _complianceItem)
	}

	logs, sub, err := _MaxBalanceModule.contract.WatchLogs(opts, "ComplianceBound", _complianceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MaxBalanceModuleComplianceBound)
				if err := _MaxBalanceModule.contract.UnpackLog(event, "ComplianceBound", log); err != nil {
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
func (_MaxBalanceModule *MaxBalanceModuleFilterer) ParseComplianceBound(log types.Log) (*MaxBalanceModuleComplianceBound, error) {
	event := new(MaxBalanceModuleComplianceBound)
	if err := _MaxBalanceModule.contract.UnpackLog(event, "ComplianceBound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MaxBalanceModuleComplianceUnboundIterator is returned from FilterComplianceUnbound and is used to iterate over the raw logs and unpacked data for ComplianceUnbound events raised by the MaxBalanceModule contract.
type MaxBalanceModuleComplianceUnboundIterator struct {
	Event *MaxBalanceModuleComplianceUnbound // Event containing the contract specifics and raw log

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
func (it *MaxBalanceModuleComplianceUnboundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MaxBalanceModuleComplianceUnbound)
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
		it.Event = new(MaxBalanceModuleComplianceUnbound)
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
func (it *MaxBalanceModuleComplianceUnboundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MaxBalanceModuleComplianceUnboundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MaxBalanceModuleComplianceUnbound represents a ComplianceUnbound event raised by the MaxBalanceModule contract.
type MaxBalanceModuleComplianceUnbound struct {
	Compliance common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterComplianceUnbound is a free log retrieval operation binding the contract event 0x408b49d9be1c914c52a0227e18a077e5a892dddf32a26cfa94a5d9708fad7718.
//
// Solidity: event ComplianceUnbound(address indexed _compliance)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) FilterComplianceUnbound(opts *bind.FilterOpts, _compliance []common.Address) (*MaxBalanceModuleComplianceUnboundIterator, error) {

	var _complianceRule []interface{}
	for _, _complianceItem := range _compliance {
		_complianceRule = append(_complianceRule, _complianceItem)
	}

	logs, sub, err := _MaxBalanceModule.contract.FilterLogs(opts, "ComplianceUnbound", _complianceRule)
	if err != nil {
		return nil, err
	}
	return &MaxBalanceModuleComplianceUnboundIterator{contract: _MaxBalanceModule.contract, event: "ComplianceUnbound", logs: logs, sub: sub}, nil
}

// WatchComplianceUnbound is a free log subscription operation binding the contract event 0x408b49d9be1c914c52a0227e18a077e5a892dddf32a26cfa94a5d9708fad7718.
//
// Solidity: event ComplianceUnbound(address indexed _compliance)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) WatchComplianceUnbound(opts *bind.WatchOpts, sink chan<- *MaxBalanceModuleComplianceUnbound, _compliance []common.Address) (event.Subscription, error) {

	var _complianceRule []interface{}
	for _, _complianceItem := range _compliance {
		_complianceRule = append(_complianceRule, _complianceItem)
	}

	logs, sub, err := _MaxBalanceModule.contract.WatchLogs(opts, "ComplianceUnbound", _complianceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MaxBalanceModuleComplianceUnbound)
				if err := _MaxBalanceModule.contract.UnpackLog(event, "ComplianceUnbound", log); err != nil {
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
func (_MaxBalanceModule *MaxBalanceModuleFilterer) ParseComplianceUnbound(log types.Log) (*MaxBalanceModuleComplianceUnbound, error) {
	event := new(MaxBalanceModuleComplianceUnbound)
	if err := _MaxBalanceModule.contract.UnpackLog(event, "ComplianceUnbound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MaxBalanceModuleIDBalancePreSetIterator is returned from FilterIDBalancePreSet and is used to iterate over the raw logs and unpacked data for IDBalancePreSet events raised by the MaxBalanceModule contract.
type MaxBalanceModuleIDBalancePreSetIterator struct {
	Event *MaxBalanceModuleIDBalancePreSet // Event containing the contract specifics and raw log

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
func (it *MaxBalanceModuleIDBalancePreSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MaxBalanceModuleIDBalancePreSet)
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
		it.Event = new(MaxBalanceModuleIDBalancePreSet)
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
func (it *MaxBalanceModuleIDBalancePreSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MaxBalanceModuleIDBalancePreSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MaxBalanceModuleIDBalancePreSet represents a IDBalancePreSet event raised by the MaxBalanceModule contract.
type MaxBalanceModuleIDBalancePreSet struct {
	Compliance common.Address
	Id         common.Address
	Balance    *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterIDBalancePreSet is a free log retrieval operation binding the contract event 0x0a777417e1c0a4b24c5f556b0c3dd522c37b99b2c3ca4bf9fea2fa1293c660c9.
//
// Solidity: event IDBalancePreSet(address indexed _compliance, address indexed _id, uint256 _balance)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) FilterIDBalancePreSet(opts *bind.FilterOpts, _compliance []common.Address, _id []common.Address) (*MaxBalanceModuleIDBalancePreSetIterator, error) {

	var _complianceRule []interface{}
	for _, _complianceItem := range _compliance {
		_complianceRule = append(_complianceRule, _complianceItem)
	}
	var _idRule []interface{}
	for _, _idItem := range _id {
		_idRule = append(_idRule, _idItem)
	}

	logs, sub, err := _MaxBalanceModule.contract.FilterLogs(opts, "IDBalancePreSet", _complianceRule, _idRule)
	if err != nil {
		return nil, err
	}
	return &MaxBalanceModuleIDBalancePreSetIterator{contract: _MaxBalanceModule.contract, event: "IDBalancePreSet", logs: logs, sub: sub}, nil
}

// WatchIDBalancePreSet is a free log subscription operation binding the contract event 0x0a777417e1c0a4b24c5f556b0c3dd522c37b99b2c3ca4bf9fea2fa1293c660c9.
//
// Solidity: event IDBalancePreSet(address indexed _compliance, address indexed _id, uint256 _balance)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) WatchIDBalancePreSet(opts *bind.WatchOpts, sink chan<- *MaxBalanceModuleIDBalancePreSet, _compliance []common.Address, _id []common.Address) (event.Subscription, error) {

	var _complianceRule []interface{}
	for _, _complianceItem := range _compliance {
		_complianceRule = append(_complianceRule, _complianceItem)
	}
	var _idRule []interface{}
	for _, _idItem := range _id {
		_idRule = append(_idRule, _idItem)
	}

	logs, sub, err := _MaxBalanceModule.contract.WatchLogs(opts, "IDBalancePreSet", _complianceRule, _idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MaxBalanceModuleIDBalancePreSet)
				if err := _MaxBalanceModule.contract.UnpackLog(event, "IDBalancePreSet", log); err != nil {
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

// ParseIDBalancePreSet is a log parse operation binding the contract event 0x0a777417e1c0a4b24c5f556b0c3dd522c37b99b2c3ca4bf9fea2fa1293c660c9.
//
// Solidity: event IDBalancePreSet(address indexed _compliance, address indexed _id, uint256 _balance)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) ParseIDBalancePreSet(log types.Log) (*MaxBalanceModuleIDBalancePreSet, error) {
	event := new(MaxBalanceModuleIDBalancePreSet)
	if err := _MaxBalanceModule.contract.UnpackLog(event, "IDBalancePreSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MaxBalanceModuleInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the MaxBalanceModule contract.
type MaxBalanceModuleInitializedIterator struct {
	Event *MaxBalanceModuleInitialized // Event containing the contract specifics and raw log

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
func (it *MaxBalanceModuleInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MaxBalanceModuleInitialized)
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
		it.Event = new(MaxBalanceModuleInitialized)
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
func (it *MaxBalanceModuleInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MaxBalanceModuleInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MaxBalanceModuleInitialized represents a Initialized event raised by the MaxBalanceModule contract.
type MaxBalanceModuleInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) FilterInitialized(opts *bind.FilterOpts) (*MaxBalanceModuleInitializedIterator, error) {

	logs, sub, err := _MaxBalanceModule.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &MaxBalanceModuleInitializedIterator{contract: _MaxBalanceModule.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *MaxBalanceModuleInitialized) (event.Subscription, error) {

	logs, sub, err := _MaxBalanceModule.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MaxBalanceModuleInitialized)
				if err := _MaxBalanceModule.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_MaxBalanceModule *MaxBalanceModuleFilterer) ParseInitialized(log types.Log) (*MaxBalanceModuleInitialized, error) {
	event := new(MaxBalanceModuleInitialized)
	if err := _MaxBalanceModule.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MaxBalanceModuleMaxBalanceSetIterator is returned from FilterMaxBalanceSet and is used to iterate over the raw logs and unpacked data for MaxBalanceSet events raised by the MaxBalanceModule contract.
type MaxBalanceModuleMaxBalanceSetIterator struct {
	Event *MaxBalanceModuleMaxBalanceSet // Event containing the contract specifics and raw log

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
func (it *MaxBalanceModuleMaxBalanceSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MaxBalanceModuleMaxBalanceSet)
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
		it.Event = new(MaxBalanceModuleMaxBalanceSet)
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
func (it *MaxBalanceModuleMaxBalanceSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MaxBalanceModuleMaxBalanceSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MaxBalanceModuleMaxBalanceSet represents a MaxBalanceSet event raised by the MaxBalanceModule contract.
type MaxBalanceModuleMaxBalanceSet struct {
	Compliance common.Address
	MaxBalance *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterMaxBalanceSet is a free log retrieval operation binding the contract event 0x8bec153ee727d17f3f63aa27cbb6bd64eb978b88821b30002a4b5517aeb10cec.
//
// Solidity: event MaxBalanceSet(address indexed _compliance, uint256 indexed _maxBalance)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) FilterMaxBalanceSet(opts *bind.FilterOpts, _compliance []common.Address, _maxBalance []*big.Int) (*MaxBalanceModuleMaxBalanceSetIterator, error) {

	var _complianceRule []interface{}
	for _, _complianceItem := range _compliance {
		_complianceRule = append(_complianceRule, _complianceItem)
	}
	var _maxBalanceRule []interface{}
	for _, _maxBalanceItem := range _maxBalance {
		_maxBalanceRule = append(_maxBalanceRule, _maxBalanceItem)
	}

	logs, sub, err := _MaxBalanceModule.contract.FilterLogs(opts, "MaxBalanceSet", _complianceRule, _maxBalanceRule)
	if err != nil {
		return nil, err
	}
	return &MaxBalanceModuleMaxBalanceSetIterator{contract: _MaxBalanceModule.contract, event: "MaxBalanceSet", logs: logs, sub: sub}, nil
}

// WatchMaxBalanceSet is a free log subscription operation binding the contract event 0x8bec153ee727d17f3f63aa27cbb6bd64eb978b88821b30002a4b5517aeb10cec.
//
// Solidity: event MaxBalanceSet(address indexed _compliance, uint256 indexed _maxBalance)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) WatchMaxBalanceSet(opts *bind.WatchOpts, sink chan<- *MaxBalanceModuleMaxBalanceSet, _compliance []common.Address, _maxBalance []*big.Int) (event.Subscription, error) {

	var _complianceRule []interface{}
	for _, _complianceItem := range _compliance {
		_complianceRule = append(_complianceRule, _complianceItem)
	}
	var _maxBalanceRule []interface{}
	for _, _maxBalanceItem := range _maxBalance {
		_maxBalanceRule = append(_maxBalanceRule, _maxBalanceItem)
	}

	logs, sub, err := _MaxBalanceModule.contract.WatchLogs(opts, "MaxBalanceSet", _complianceRule, _maxBalanceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MaxBalanceModuleMaxBalanceSet)
				if err := _MaxBalanceModule.contract.UnpackLog(event, "MaxBalanceSet", log); err != nil {
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

// ParseMaxBalanceSet is a log parse operation binding the contract event 0x8bec153ee727d17f3f63aa27cbb6bd64eb978b88821b30002a4b5517aeb10cec.
//
// Solidity: event MaxBalanceSet(address indexed _compliance, uint256 indexed _maxBalance)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) ParseMaxBalanceSet(log types.Log) (*MaxBalanceModuleMaxBalanceSet, error) {
	event := new(MaxBalanceModuleMaxBalanceSet)
	if err := _MaxBalanceModule.contract.UnpackLog(event, "MaxBalanceSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MaxBalanceModuleOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the MaxBalanceModule contract.
type MaxBalanceModuleOwnershipTransferredIterator struct {
	Event *MaxBalanceModuleOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *MaxBalanceModuleOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MaxBalanceModuleOwnershipTransferred)
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
		it.Event = new(MaxBalanceModuleOwnershipTransferred)
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
func (it *MaxBalanceModuleOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MaxBalanceModuleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MaxBalanceModuleOwnershipTransferred represents a OwnershipTransferred event raised by the MaxBalanceModule contract.
type MaxBalanceModuleOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MaxBalanceModuleOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MaxBalanceModule.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MaxBalanceModuleOwnershipTransferredIterator{contract: _MaxBalanceModule.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MaxBalanceModuleOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MaxBalanceModule.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MaxBalanceModuleOwnershipTransferred)
				if err := _MaxBalanceModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_MaxBalanceModule *MaxBalanceModuleFilterer) ParseOwnershipTransferred(log types.Log) (*MaxBalanceModuleOwnershipTransferred, error) {
	event := new(MaxBalanceModuleOwnershipTransferred)
	if err := _MaxBalanceModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MaxBalanceModulePresetCompletedIterator is returned from FilterPresetCompleted and is used to iterate over the raw logs and unpacked data for PresetCompleted events raised by the MaxBalanceModule contract.
type MaxBalanceModulePresetCompletedIterator struct {
	Event *MaxBalanceModulePresetCompleted // Event containing the contract specifics and raw log

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
func (it *MaxBalanceModulePresetCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MaxBalanceModulePresetCompleted)
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
		it.Event = new(MaxBalanceModulePresetCompleted)
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
func (it *MaxBalanceModulePresetCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MaxBalanceModulePresetCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MaxBalanceModulePresetCompleted represents a PresetCompleted event raised by the MaxBalanceModule contract.
type MaxBalanceModulePresetCompleted struct {
	Compliance common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPresetCompleted is a free log retrieval operation binding the contract event 0x40bfb91138c8514efd65a280ed9f436317a8e93584e42ef52b3b863e68bd500f.
//
// Solidity: event PresetCompleted(address indexed _compliance)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) FilterPresetCompleted(opts *bind.FilterOpts, _compliance []common.Address) (*MaxBalanceModulePresetCompletedIterator, error) {

	var _complianceRule []interface{}
	for _, _complianceItem := range _compliance {
		_complianceRule = append(_complianceRule, _complianceItem)
	}

	logs, sub, err := _MaxBalanceModule.contract.FilterLogs(opts, "PresetCompleted", _complianceRule)
	if err != nil {
		return nil, err
	}
	return &MaxBalanceModulePresetCompletedIterator{contract: _MaxBalanceModule.contract, event: "PresetCompleted", logs: logs, sub: sub}, nil
}

// WatchPresetCompleted is a free log subscription operation binding the contract event 0x40bfb91138c8514efd65a280ed9f436317a8e93584e42ef52b3b863e68bd500f.
//
// Solidity: event PresetCompleted(address indexed _compliance)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) WatchPresetCompleted(opts *bind.WatchOpts, sink chan<- *MaxBalanceModulePresetCompleted, _compliance []common.Address) (event.Subscription, error) {

	var _complianceRule []interface{}
	for _, _complianceItem := range _compliance {
		_complianceRule = append(_complianceRule, _complianceItem)
	}

	logs, sub, err := _MaxBalanceModule.contract.WatchLogs(opts, "PresetCompleted", _complianceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MaxBalanceModulePresetCompleted)
				if err := _MaxBalanceModule.contract.UnpackLog(event, "PresetCompleted", log); err != nil {
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

// ParsePresetCompleted is a log parse operation binding the contract event 0x40bfb91138c8514efd65a280ed9f436317a8e93584e42ef52b3b863e68bd500f.
//
// Solidity: event PresetCompleted(address indexed _compliance)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) ParsePresetCompleted(log types.Log) (*MaxBalanceModulePresetCompleted, error) {
	event := new(MaxBalanceModulePresetCompleted)
	if err := _MaxBalanceModule.contract.UnpackLog(event, "PresetCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MaxBalanceModuleUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the MaxBalanceModule contract.
type MaxBalanceModuleUpgradedIterator struct {
	Event *MaxBalanceModuleUpgraded // Event containing the contract specifics and raw log

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
func (it *MaxBalanceModuleUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MaxBalanceModuleUpgraded)
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
		it.Event = new(MaxBalanceModuleUpgraded)
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
func (it *MaxBalanceModuleUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MaxBalanceModuleUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MaxBalanceModuleUpgraded represents a Upgraded event raised by the MaxBalanceModule contract.
type MaxBalanceModuleUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*MaxBalanceModuleUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _MaxBalanceModule.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &MaxBalanceModuleUpgradedIterator{contract: _MaxBalanceModule.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_MaxBalanceModule *MaxBalanceModuleFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *MaxBalanceModuleUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _MaxBalanceModule.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MaxBalanceModuleUpgraded)
				if err := _MaxBalanceModule.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_MaxBalanceModule *MaxBalanceModuleFilterer) ParseUpgraded(log types.Log) (*MaxBalanceModuleUpgraded, error) {
	event := new(MaxBalanceModuleUpgraded)
	if err := _MaxBalanceModule.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
