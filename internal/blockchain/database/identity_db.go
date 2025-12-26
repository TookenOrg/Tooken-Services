package database

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

func GetIdentityAddrByUserId(ctx context.Context, userId int) (identityAddr *common.Address, err error) {
	return &common.MaxAddress, nil
}
