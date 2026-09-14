package utils

import (
	"crypto/ecdsa"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func getECDSAKeys() (privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) {
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatalf("Error getting private key: %s", err.Error())
	}
	publicKey = &privateKey.PublicKey
	return
}

func GetEthFrom() common.Address {
	_, publicKeyPtr := getECDSAKeys()
	if publicKeyPtr == nil {
		log.Fatalf("Error getting ECDSA keys")
	}
	publicKey := *publicKeyPtr
	return crypto.PubkeyToAddress(publicKey)
}
