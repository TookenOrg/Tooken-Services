package models

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type SignatureResult struct {
	Wallet          string
	IdentityToClaim string
	Topic           *big.Int
	DataBytes       []byte
	DataHash        common.Hash
	EthSignedHash   common.Hash
	Signature       string // Hex string representation of the signature
	SignatureBytes  []byte // Raw bytes of the signature
	V               uint8
	R               common.Hash
	S               common.Hash
	RecoveredAddr   common.Address
	Key             common.Hash
}
