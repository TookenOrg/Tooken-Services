package services

import (
	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/assets_managements/database"
	"github.com/shopspring/decimal"
)

// Boundary between the database schema and the API contract.
//
// This is the only file that knows about both. The database layer no longer
// depends on the code generated from openapi.yaml: changing the contract no
// longer forces a change in the SQL, and the other way around.
//
// ⚠️ Accepted temporary debt: server.* still exposes amounts as *float32.
// decimalToFloat32 is the single place where the exact precision of the NUMERIC
// values is lost. Switching the OpenAPI amounts to `type: string` will remove
// that conversion — there will be one function to replace.

func toServerRealEstate(dto database.RealEstateDTO) server.RealEstate {
	e := server.RealEstate{
		Id:              dto.Id,
		Title:           dto.Title,
		Description:     dto.Description,
		Imageurl:        dto.Imageurl,
		EstateTypeId:    dto.EstateTypeId,
		EstateType:      dto.EstateTypeName,
		Active:          dto.Active,
		CreatedAt:       dto.CreatedAt,
		ContractAddress: dto.ContractAddress,
		Progression: &server.RealEstateProgression{
			TokensSold:     intPtr(int(dto.Progression.TokensSold)),
			TokensSoldPctg: decimalToFloat32(dto.Progression.TokensSoldPctg),
		},
	}

	if s := dto.Specification; s != nil {
		e.Specificiation = &server.RealEstateSpecification{
			SurfaceArea:    nullDecimalToFloat32(s.SurfaceArea),
			LotSize:        nullDecimalToFloat32(s.LotSize),
			PoolSize:       nullDecimalToFloat32(s.PoolSize),
			TerraceSize:    nullDecimalToFloat32(s.TerraceSize),
			BedroomNumber:  s.BedroomNumber,
			BathroomNumber: s.BathroomNumber,
			BuildYear:      s.BuiltYear,
		}
	}

	if c := dto.SharesConfig; c != nil {
		e.Configuration = &server.RealEstateConfiguration{
			TotalShares:            nullDecimalToInt(c.TotalShares),
			PricePerShare:          nullDecimalToFloat32(c.PricePerShare),
			Yield:                  nullDecimalToFloat32(c.Yield),
			PaymentFrequency:       c.PaymentFrequency,
			PaymentFrequencyTypeId: c.PaymentFrequencyTypeId,
			PaymentFrequencyType:   c.PaymentFrequencyTypeName,
		}
	}

	return e
}

func toServerRealEstates(dtos []database.RealEstateDTO) []server.RealEstate {
	if dtos == nil {
		return nil
	}

	estates := make([]server.RealEstate, 0, len(dtos))
	for _, dto := range dtos {
		estates = append(estates, toServerRealEstate(dto))
	}

	return estates
}

func intPtr(v int) *int { return &v }

func decimalToFloat32(d decimal.Decimal) *float32 {
	v, _ := d.Float64()
	f := float32(v)
	return &f
}

func nullDecimalToFloat32(d decimal.NullDecimal) *float32 {
	if !d.Valid {
		return nil
	}
	return decimalToFloat32(d.Decimal)
}

// nullDecimalToInt: total_shares is a NUMERIC(20,0), so an exact integer.
// IntPart truncates beyond int64, which stays far above any realistic share
// count.
func nullDecimalToInt(d decimal.NullDecimal) *int {
	if !d.Valid {
		return nil
	}
	return intPtr(int(d.Decimal.IntPart()))
}
