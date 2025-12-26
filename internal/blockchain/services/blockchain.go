package services

import (
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
)

func SetupEthClient() (*ethclient.Client, error) {
	rpcURL := os.Getenv("WS_RPC_URL")
	return ethclient.Dial(rpcURL)
}
