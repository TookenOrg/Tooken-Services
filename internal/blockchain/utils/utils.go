package utils

import (
	"errors"
	"math/big"

	"github.com/shopspring/decimal"
)

func ConvertFloatToWei(f float64, decimalNumber int64) (wei *big.Int, err error) {
	dec := decimal.NewFromFloat(f)
	if dec.IsNegative() {
		return nil, errors.New("amount must not be negative")
	}

	scale := decimal.NewFromInt(10).Pow(decimal.NewFromInt(decimalNumber))
	res := dec.Mul(scale)

	return res.BigInt(), nil
}
