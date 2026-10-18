package services

import (
	"fmt"
	"math/big"
	"strings"

	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
)

// sharesToBaseUnits converts a number of shares into the token base units the
// capping modules compare balances against (B3): "50" with 18 decimals is
// 50 × 10^18.
//
// Only a strictly positive whole number is accepted. A share cannot be split,
// and a cap of zero blocks every mint (B2), so "0", "-3", "1.5", "1e3" and ""
// are all errors rather than values to round.
func sharesToBaseUnits(shares string, decimals int) (*big.Int, error) {
	if decimals < 0 {
		return nil, fmt.Errorf("decimals must not be negative, got %d", decimals)
	}
	// big.Int.SetString alone would also accept a sign and "_" separators.
	if shares == "" || strings.TrimLeft(shares, "0123456789") != "" {
		return nil, fmt.Errorf("shares must be a whole number, got %q", shares)
	}
	n, _ := new(big.Int).SetString(shares, 10)
	if n.Sign() == 0 {
		return nil, fmt.Errorf("shares must be greater than zero, got %q", shares)
	}

	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	return n.Mul(n, scale), nil
}

// encodeCapSettings returns the two calls the factory replays on the capping
// modules while it deploys the token: setMaxBalance(perInvestor) on
// MaxBalanceModule, then setSupplyLimit(supply) on SupplyLimitModule.
//
// Both values are in base units. The order is the order of the modules in
// TokenDetails (B1), and these setters only accept calls coming from the
// token's compliance (B4), which is why they cannot be sent afterwards.
func encodeCapSettings(perInvestor, supply *big.Int) (maxBalance, supplyLimit []byte, err error) {
	maxBalanceABI, err := contracts.MaxBalanceModuleMetaData.GetAbi()
	if err != nil {
		return nil, nil, err
	}
	if maxBalance, err = maxBalanceABI.Pack("setMaxBalance", perInvestor); err != nil {
		return nil, nil, fmt.Errorf("encode setMaxBalance: %w", err)
	}

	supplyLimitABI, err := contracts.SupplyLimitModuleMetaData.GetAbi()
	if err != nil {
		return nil, nil, err
	}
	if supplyLimit, err = supplyLimitABI.Pack("setSupplyLimit", supply); err != nil {
		return nil, nil, fmt.Errorf("encode setSupplyLimit: %w", err)
	}

	return maxBalance, supplyLimit, nil
}
