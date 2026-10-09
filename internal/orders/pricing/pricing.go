package pricing

import "github.com/shopspring/decimal"

// PricedOrder is what the investor owes, frozen at order time.
// Every amount is already rounded to 2 decimals, exactly as the database CHECKs expect.
type PricedOrder struct {
	UnitPrice   decimal.Decimal // copied as is from shares_config (8 decimals)
	GrossAmount decimal.Decimal // round(quantity × unitPrice, 2)
	FeeAmount   decimal.Decimal // sum of Fees[i].Amount; zero when there is no fee
	AmountDue   decimal.Decimal // GrossAmount + FeeAmount
	Fees        []PricedFee     // empty when there is no fee (U4), never nil
}

// PricedFee is one line of iss.issuance_order_fees.
type PricedFee struct {
	Code       string          // "ENTRY"
	Rate       decimal.Decimal // percentage, e.g. 2.0000 for 2 %
	BaseAmount decimal.Decimal // what the rate applies to: GrossAmount for ENTRY
	Amount     decimal.Decimal // round(BaseAmount × Rate / 100, 2)
}

// PriceOrder freezes what the investor owes. Rounding must match the
// database CHECKs (round half away from zero, 2 decimals).
func PriceOrder(quantity int64, unitPrice decimal.Decimal, entryFeeRate *decimal.Decimal) PricedOrder {

	pricedOrder := PricedOrder{
		UnitPrice: unitPrice,
		Fees:      []PricedFee{},
	}

	total := unitPrice.Mul(decimal.NewFromInt(quantity)).Round(2)
	pricedOrder.GrossAmount = total
	amountDue := total
	if entryFeeRate != nil && entryFeeRate.GreaterThan(decimal.Zero) {
		entryFee := total.Mul(*entryFeeRate).Div(decimal.NewFromInt(100)).Round(2)
		if entryFee.GreaterThan(decimal.Zero) {
			amountDue = total.Add(entryFee)
			pricedOrder.FeeAmount = entryFee
			pricedOrder.Fees = append(pricedOrder.Fees, PricedFee{
				Code:       "ENTRY",
				Rate:       *entryFeeRate,
				BaseAmount: total,
				Amount:     entryFee,
			})
		}
	}

	pricedOrder.AmountDue = amountDue

	return pricedOrder
}
