package database

import (
	"time"

	"github.com/shopspring/decimal"
)

// The DTOs mirror the v3 schema (see migrations 000001 to 000010), not the API
// contract. Translation to server.* lives in the services layer: that is what
// keeps this layer independent from the code generated out of openapi.yaml.
//
// Conventions:
//   - NOT NULL column          -> value field
//   - nullable column          -> pointer (database/sql knows how to fill a **T)
//   - NUMERIC amount / measure -> decimal.Decimal, never a float
//     (nullable: decimal.NullDecimal, since decimal.Decimal refuses to scan NULL)
//   - missing 1-1 table        -> nil sub-struct
//
// A field reached through a LEFT JOIN is nullable *because of the join*, even
// when the column is NOT NULL in its own table: row presence is carried by the
// sub-struct, never guessed from a value.

// RealEstateDTO is the full projection of a real estate asset.
type RealEstateDTO struct {
	Id              int
	Title           *string
	Description     *string
	Imageurl        *string
	EstateTypeId    *int
	EstateTypeName  *string
	Active          *bool
	StatusId        int
	StatusCode      string
	IssuerId        *int
	TokenId         *int
	ContractAddress *string
	CreatedAt       *time.Time
	UpdatedAt       time.Time
	PublishedAt     *time.Time

	// Coarse location, read by both the list and the detail. Separate from
	// Address on purpose: this is the part that is safe to expose publicly.
	Location *RealEstateLocationDTO

	// 1-1 subsets. nil means the row does not exist in the database, or was not
	// requested by the query (see the list, which does not load the full address).
	Specification *RealEstateSpecificationDTO
	SharesConfig  *RealEstateSharesConfigDTO
	Address       *RealEstateAddressDTO
	Issuer        *RealEstateIssuerDTO
	Token         *RealEstateTokenDTO

	// Always present: an asset with no order has a progression of zero, not a
	// missing progression.
	Progression RealEstateProgressionDTO

	// 1-N, loaded separately (detail only).
	Media []RealEstateMediaDTO

	// Presence markers filled by the scan. Unexported: they are an
	// implementation detail, callers read Specification == nil.
	specificationPresent bool
	sharesConfigPresent  bool
	locationPresent      bool
}

// RealEstateLocationDTO holds the only two address columns a public endpoint
// may read. The street, the postal code and the coordinates are never selected
// by the list query, so no amount of mapper carelessness can leak them.
type RealEstateLocationDTO struct {
	City        string
	CountryCode string
}

// RealEstateSpecificationDTO maps ass.real_estate_specification (1-1).
type RealEstateSpecificationDTO struct {
	SurfaceArea    decimal.NullDecimal
	LotSize        decimal.NullDecimal
	PoolSize       decimal.NullDecimal
	TerraceSize    decimal.NullDecimal
	BedroomNumber  *int
	BathroomNumber *int
	BuiltYear      *int
	EnergyClass    *string
	GesClass       *string
}

// RealEstateSharesConfigDTO maps ass.real_estate_shares_config (1-1).
//
// TotalValuation is GENERATED ALWAYS in the database: it cannot diverge from
// TotalShares × PricePerShare, and is never written from Go.
type RealEstateSharesConfigDTO struct {
	CompartmentRef           *string
	TotalShares              decimal.NullDecimal
	PricePerShare            decimal.NullDecimal
	TotalValuation           decimal.NullDecimal
	CurrencyCode             *string
	Yield                    decimal.NullDecimal
	PaymentFrequency         *int
	PaymentFrequencyTypeId   *int
	PaymentFrequencyTypeName *string
	MinInvestment            decimal.NullDecimal
	MaxInvestment            decimal.NullDecimal
	EntryFeeRate             decimal.NullDecimal
	MgmtFeeRate              decimal.NullDecimal
	ExitFeeRate              decimal.NullDecimal
}

// RealEstateAddressDTO maps ass.real_estate_address (1-1).
//
// The most sensitive data in the model: loaded by the detail endpoint, never by
// the list. Latitude and Longitude are NUMERIC(9,6) and are always either both
// set or both null (constraint real_estate_address_coordinates_ck).
type RealEstateAddressDTO struct {
	Street           string
	StreetComplement *string
	PostalCode       string
	City             string
	Region           *string
	CountryCode      string
	Latitude         decimal.NullDecimal
	Longitude        decimal.NullDecimal
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// RealEstateIssuerDTO maps ass.issuer, the legal issuing vehicle.
type RealEstateIssuerDTO struct {
	Id                 int
	Name               string
	LegalForm          string
	RegistrationNumber *string
	CountryCode        string
	LeiCode            *string
	StatusId           int
	StatusCode         string
}

// RealEstateTokenDTO maps blk.token, the ERC-3643 contract of the asset.
//
// nil means "not tokenized yet": an explicit business state since migration
// 000008, where contract_address could only say "empty string".
type RealEstateTokenDTO struct {
	Id                    int
	Symbol                string
	TokenName             string
	Address               string
	NbDecimal             int
	ModularComplianceAddr *string
	CreatedAt             *time.Time
}

// RealEstateMediaDTO maps ass.real_estate_media (1-N), ordered by Position.
type RealEstateMediaDTO struct {
	Id        int
	Url       string
	AltText   *string
	MediaType string
	Position  int
	IsCover   bool
	CreatedAt time.Time
}

// RealEstateProgressionDTO holds the reserved shares, computed in SQL.
//
// What counts as "reserved" is carried by
// iss.issuance_order_statuses.counts_as_reserved (migration 000010), not by this
// code: adding a status to the flow requires no change here.
type RealEstateProgressionDTO struct {
	TokensSold     int64
	TokensSoldPctg decimal.Decimal
}
