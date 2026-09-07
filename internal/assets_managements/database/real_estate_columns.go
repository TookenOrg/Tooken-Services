package database

import (
	"fmt"
	"strings"
)

// Defect #1 fixed: the SELECT column list and the Scan target list are no
// longer two parallel lists someone has to remember to keep in sync — they are
// the same list.
//
// Inserting a column in the middle can no longer shift the scan silently: the
// SQL expression and the Go field travel together.
type realEstateColumn struct {
	expr string
	dest func(*RealEstateDTO) any
}

// selectList renders the expressions in order, indented so they stay readable
// in query logs.
func selectList(cols []realEstateColumn) string {
	exprs := make([]string, len(cols))
	for i, c := range cols {
		exprs[i] = "    " + c.expr
	}
	return strings.Join(exprs, ",\n")
}

// scanTargets projects the pointers in the same order as selectList.
func scanTargets(dto *RealEstateDTO, cols []realEstateColumn) []any {
	dest := make([]any, len(cols))
	for i, c := range cols {
		dest[i] = c.dest(dto)
	}
	return dest
}

// realEstateCoreColumns covers the asset, its type, its status and its
// progression. Used by both the list and the detail queries.
var realEstateCoreColumns = []realEstateColumn{
	{"re.id", func(d *RealEstateDTO) any { return &d.Id }},
	{"re.title", func(d *RealEstateDTO) any { return &d.Title }},
	{"re.description", func(d *RealEstateDTO) any { return &d.Description }},
	{"re.imageurl", func(d *RealEstateDTO) any { return &d.Imageurl }},
	{"re.estate_type AS estate_type_id", func(d *RealEstateDTO) any { return &d.EstateTypeId }},
	{"ret.name AS estate_type", func(d *RealEstateDTO) any { return &d.EstateTypeName }},
	{"re.active", func(d *RealEstateDTO) any { return &d.Active }},
	{"re.status_id", func(d *RealEstateDTO) any { return &d.StatusId }},
	{"rst.code AS status_code", func(d *RealEstateDTO) any { return &d.StatusCode }},
	{"re.issuer_id", func(d *RealEstateDTO) any { return &d.IssuerId }},
	{"re.token_id", func(d *RealEstateDTO) any { return &d.TokenId }},
	{"re.contract_address", func(d *RealEstateDTO) any { return &d.ContractAddress }},
	{"re.created_at", func(d *RealEstateDTO) any { return &d.CreatedAt }},
	{"re.updated_at", func(d *RealEstateDTO) any { return &d.UpdatedAt }},
	{"re.published_at", func(d *RealEstateDTO) any { return &d.PublishedAt }},

	// Presence of the 1-1 tables is read on the join key, never guessed from a
	// business value that could legitimately be NULL.
	{"(res.real_estate_id IS NOT NULL) AS has_specification",
		func(d *RealEstateDTO) any { return &d.specificationPresent }},
	{"res.surface_area", func(d *RealEstateDTO) any { return &d.Specification.SurfaceArea }},
	{"res.lot_size", func(d *RealEstateDTO) any { return &d.Specification.LotSize }},
	{"res.pool_size", func(d *RealEstateDTO) any { return &d.Specification.PoolSize }},
	{"res.terrace_size", func(d *RealEstateDTO) any { return &d.Specification.TerraceSize }},
	{"res.bedroom_number", func(d *RealEstateDTO) any { return &d.Specification.BedroomNumber }},
	{"res.bathroom_number", func(d *RealEstateDTO) any { return &d.Specification.BathroomNumber }},
	{"res.built_year", func(d *RealEstateDTO) any { return &d.Specification.BuiltYear }},
	{"res.energy_class", func(d *RealEstateDTO) any { return &d.Specification.EnergyClass }},
	{"res.ges_class", func(d *RealEstateDTO) any { return &d.Specification.GesClass }},

	{"(addr.real_estate_id IS NOT NULL) AS has_location",
		func(d *RealEstateDTO) any { return &d.locationPresent }},
	// city and country_code are NOT NULL in the table, but the LEFT JOIN turns
	// them into NULL for an asset without an address. COALESCE keeps the scan
	// targets as plain strings; the empty value is never read since the
	// Location sub-struct is dropped when has_location is false.
	{"COALESCE(addr.city, '') AS city", func(d *RealEstateDTO) any { return &d.Location.City }},
	{"COALESCE(addr.country_code, '') AS country_code",
		func(d *RealEstateDTO) any { return &d.Location.CountryCode }},

	{"(conf.real_estate_id IS NOT NULL) AS has_shares_config",
		func(d *RealEstateDTO) any { return &d.sharesConfigPresent }},
	{"conf.compartment_ref", func(d *RealEstateDTO) any { return &d.SharesConfig.CompartmentRef }},
	{"conf.total_shares", func(d *RealEstateDTO) any { return &d.SharesConfig.TotalShares }},
	{"conf.price_per_share", func(d *RealEstateDTO) any { return &d.SharesConfig.PricePerShare }},
	{"conf.total_valuation", func(d *RealEstateDTO) any { return &d.SharesConfig.TotalValuation }},
	{"conf.currency_code", func(d *RealEstateDTO) any { return &d.SharesConfig.CurrencyCode }},
	{"conf.yield", func(d *RealEstateDTO) any { return &d.SharesConfig.Yield }},
	{"conf.payment_frequency", func(d *RealEstateDTO) any { return &d.SharesConfig.PaymentFrequency }},
	{"conf.payment_frequency_type AS payment_frequency_type_id",
		func(d *RealEstateDTO) any { return &d.SharesConfig.PaymentFrequencyTypeId }},
	{"payt.name AS payment_frequency_type",
		func(d *RealEstateDTO) any { return &d.SharesConfig.PaymentFrequencyTypeName }},
	{"conf.min_investment", func(d *RealEstateDTO) any { return &d.SharesConfig.MinInvestment }},
	{"conf.max_investment", func(d *RealEstateDTO) any { return &d.SharesConfig.MaxInvestment }},
	{"conf.entry_fee_rate", func(d *RealEstateDTO) any { return &d.SharesConfig.EntryFeeRate }},
	{"conf.mgmt_fee_rate", func(d *RealEstateDTO) any { return &d.SharesConfig.MgmtFeeRate }},
	{"conf.exit_fee_rate", func(d *RealEstateDTO) any { return &d.SharesConfig.ExitFeeRate }},

	{"COALESCE(sold.reserved, 0)::bigint AS tokens_sold",
		func(d *RealEstateDTO) any { return &d.Progression.TokensSold }},
	{`CASE
        WHEN conf.total_shares IS NULL OR conf.total_shares = 0 THEN 0
        ELSE ROUND((COALESCE(sold.reserved, 0)::numeric / conf.total_shares) * 100, 2)
    END AS tokens_sold_pctg`,
		func(d *RealEstateDTO) any { return &d.Progression.TokensSoldPctg }},
}

// realEstateBaseFrom holds the joins required by the columns above.
//
// The issuer, the token, the media and the full address are absent on purpose:
// they are loaded separately by the detail query.
const realEstateBaseFrom = `
FROM ass.real_estate re
LEFT JOIN ass.real_estate_type ret
    ON re.estate_type = ret.id
LEFT JOIN ass.real_estate_status rst
    ON re.status_id = rst.id
LEFT JOIN ass.real_estate_specification res
    ON re.id = res.real_estate_id
LEFT JOIN ass.real_estate_shares_config conf
    ON re.id = conf.real_estate_id
LEFT JOIN ass.payment_frequency_type payt
    ON conf.payment_frequency_type = payt.id

-- The address table is joined, but only city and country_code are ever
-- selected here. The street, the postal code and the coordinates pinpoint a
-- tokenized asset: they stay out of any query that can serve a public list.
LEFT JOIN ass.real_estate_address addr
    ON re.id = addr.real_estate_id

-- Shares already taken. The aggregate is computed ONCE and then joined: a
-- correlated subquery would be re-executed per asset (N+1) — measured at 196 ms
-- against 9 ms on 500 assets / 50,000 orders. The same form serves the detail
-- query: PostgreSQL pushes the re.id = $1 predicate into the aggregate.
--
-- What counts as a reserved share is not encoded here but carried by
-- iss.issuance_order_statuses.counts_as_reserved: adding a status to the flow
-- requires no change to this query.
LEFT JOIN (
    SELECT o.asset_id, SUM(o.quantity) AS reserved
    FROM iss.issuance_orders o
    JOIN iss.issuance_order_statuses s
        ON s.id = o.status_id
    WHERE s.counts_as_reserved
    GROUP BY o.asset_id
) sold
    ON sold.asset_id = re.id
`

// buildRealEstateQuery assembles the SELECT from the column table.
//
// `clauses` carries everything that follows the FROM: the WHERE filter and,
// for any query that can return several rows, an explicit ORDER BY.
func buildRealEstateQuery(cols []realEstateColumn, clauses string) string {
	return fmt.Sprintf("SELECT\n%s\n%s%s", selectList(cols), realEstateBaseFrom, clauses)
}
