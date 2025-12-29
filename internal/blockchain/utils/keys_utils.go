package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"io"
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

func EncryptAESGCM(plaintext string) (string, error) {

	// TODO: https://github.com/TookenOrg/tooken-services/issues/12
	key := []byte("12345678901234567890123456789012")
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
