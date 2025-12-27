package utils

import (
	"math/big"

	"github.com/shopspring/decimal"
)

func ConvertFloatToWei(f float64, decimalNumber int64) (wei *big.Int, err error) {
	dec := decimal.NewFromFloat(f)
	if dec.IsNegative() {
		return
	}

	scale := decimal.NewFromInt(10).Pow(decimal.NewFromInt(decimalNumber))
	res := dec.Mul(scale)

	return res.BigInt(), nil
}
