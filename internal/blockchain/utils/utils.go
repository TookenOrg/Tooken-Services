package utils

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/shopspring/decimal"
)

// maxExactDigits is the number of significant decimal digits a float64 is
// guaranteed to carry without alteration. Beyond that, the value received is
// not necessarily the value sent, and no check performed here can recover it.
const maxExactDigits = 15

// maxExactInteger is 2^53, above which consecutive integers are no longer all
// representable as a float64: 10000000000000001 arrives as 10000000000000000.
const maxExactInteger = float64(1 << 53)

// ConvertFloatToWei turns a human-readable amount into its on-chain integer
// representation.
//
// It refuses rather than approximates. The result is written to an immutable
// ledger, so returning a value that merely resembles the one requested is worse
// than returning an error: a rejected request can be retried, a wrong mint
// cannot be undone. Every case below used to pass silently.
func ConvertFloatToWei(f float64, decimalNumber int64) (wei *big.Int, err error) {
	// decimal.NewFromFloat panics on NaN and Inf. Mint and burn run in a
	// goroutine, so that panic would take the whole process down.
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return nil, errors.New("amount must be a finite number")
	}

	if f < 0 {
		return nil, errors.New("amount must not be negative")
	}

	// A negative scale would divide instead of multiplying, and could still
	// land on an integer: 100 with -2 decimals used to return 1.
	if decimalNumber < 0 {
		return nil, fmt.Errorf("decimals must not be negative, got %d", decimalNumber)
	}

	if f >= maxExactInteger {
		return nil, fmt.Errorf(
			"amount %v is too large to be exact: above %.0f a float64 cannot represent every unit",
			f, maxExactInteger)
	}

	dec := decimal.NewFromFloat(f)
	if dec.NumDigits() > maxExactDigits {
		return nil, fmt.Errorf(
			"amount %v carries more than %d significant digits and cannot be trusted to the last one",
			f, maxExactDigits)
	}

	scale := decimal.NewFromInt(10).Pow(decimal.NewFromInt(decimalNumber))
	res := dec.Mul(scale)

	// BigInt truncates towards zero. Without this check, 1.005 on a token with
	// 2 decimals silently became 1.00, and 0.9 on a token with 0 decimals
	// became a confirmed transaction that issued nothing at all.
	if !res.IsInteger() {
		return nil, fmt.Errorf(
			"amount %v is finer than the %d decimals of the token and would be truncated to %s",
			f, decimalNumber, res.Truncate(0).Div(scale).String())
	}

	return res.BigInt(), nil
}
