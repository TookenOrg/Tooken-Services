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
// The DTO models the whole v3 schema; this file decides what actually leaves
// the server. Two deliberate omissions:
//
//   - the full address (street, postal code, coordinates). Both real estate
//     endpoints are public today, so only the coarse location is exposed. The
//     rest waits for the endpoint to be able to tell a KYC-verified caller from
//     an anonymous one.
//   - the fee grid, the investment bounds and the compartment reference, which
//     only mean something once the subscription funnel exists.
//
// Amounts cross the boundary as exact decimal strings: decimal.Decimal.String()
// never rounds, where a float32 turned 199.99 into 199.99001.

func toServerRealEstateSummary(dto database.RealEstateDTO) server.RealEstateSummary {
	e := server.RealEstateSummary{
		Id:              dto.Id,
		Title:           dto.Title,
		Description:     dto.Description,
		Imageurl:        dto.Imageurl,
		EstateTypeId:    dto.EstateTypeId,
		EstateType:      dto.EstateTypeName,
		Status:          strPtr(dto.StatusCode),
		StatusId:        intPtr(dto.StatusId),
		Active:          dto.Active,
		CreatedAt:       dto.CreatedAt,
		ContractAddress: dto.ContractAddress,
		Progression: &server.RealEstateProgression{
			TokensSold:     intPtr(int(dto.Progression.TokensSold)),
			TokensSoldPctg: decimalToFloat32(dto.Progression.TokensSoldPctg),
		},
	}

	if l := dto.Location; l != nil {
		e.Location = &server.RealEstateLocation{
			City:        strPtr(l.City),
			CountryCode: strPtr(l.CountryCode),
		}
	}

	if s := dto.Specification; s != nil {
		e.Specification = &server.RealEstateSpecification{
			SurfaceArea:    nullDecimalToFloat32(s.SurfaceArea),
			LotSize:        nullDecimalToFloat32(s.LotSize),
			PoolSize:       nullDecimalToFloat32(s.PoolSize),
			TerraceSize:    nullDecimalToFloat32(s.TerraceSize),
			BedroomNumber:  s.BedroomNumber,
			BathroomNumber: s.BathroomNumber,
			BuildYear:      s.BuiltYear,
			EnergyClass:    s.EnergyClass,
			GesClass:       s.GesClass,
		}
	}

	if c := dto.SharesConfig; c != nil {
		e.Configuration = &server.RealEstateConfiguration{
			TotalShares:            nullDecimalToString(c.TotalShares),
			PricePerShare:          nullDecimalToString(c.PricePerShare),
			TotalValuation:         nullDecimalToString(c.TotalValuation),
			CurrencyCode:           c.CurrencyCode,
			Yield:                  nullDecimalToFloat32(c.Yield),
			PaymentFrequency:       c.PaymentFrequency,
			PaymentFrequencyTypeId: c.PaymentFrequencyTypeId,
			PaymentFrequencyType:   c.PaymentFrequencyTypeName,
		}
	}

	return e
}

func toServerRealEstates(dtos []database.RealEstateDTO) []server.RealEstateSummary {
	if dtos == nil {
		return nil
	}

	estates := make([]server.RealEstateSummary, 0, len(dtos))
	for _, dto := range dtos {
		estates = append(estates, toServerRealEstateSummary(dto))
	}

	return estates
}

func toServerRealEstate(dto database.RealEstateDTO) server.RealEstate {
	s := toServerRealEstateSummary(dto)

	e := server.RealEstate{
		Id:              s.Id,
		Title:           s.Title,
		Description:     s.Description,
		Imageurl:        s.Imageurl,
		EstateTypeId:    s.EstateTypeId,
		EstateType:      s.EstateType,
		Status:          s.Status,
		StatusId:        s.StatusId,
		Active:          s.Active,
		CreatedAt:       s.CreatedAt,
		ContractAddress: s.ContractAddress,
		Location:        s.Location,
		Specification:   s.Specification,
		Configuration:   s.Configuration,
		Progression:     s.Progression,
		UpdatedAt:       &dto.UpdatedAt,
		PublishedAt:     dto.PublishedAt,
	}

	if t := dto.Token; t != nil {
		e.Token = &server.RealEstateToken{
			Id:        intPtr(t.Id),
			Symbol:    strPtr(t.Symbol),
			TokenName: strPtr(t.TokenName),
			Address:   strPtr(t.Address),
			NbDecimal: intPtr(t.NbDecimal),
		}
	}

	if i := dto.Issuer; i != nil {
		e.Issuer = &server.RealEstateIssuer{
			Id:          intPtr(i.Id),
			Name:        strPtr(i.Name),
			LegalForm:   strPtr(i.LegalForm),
			CountryCode: strPtr(i.CountryCode),
			Status:      strPtr(i.StatusCode),
		}
	}

	if len(dto.Media) > 0 {
		media := make([]server.RealEstateMedia, 0, len(dto.Media))
		for _, m := range dto.Media {
			media = append(media, server.RealEstateMedia{
				Id:        intPtr(m.Id),
				Url:       strPtr(m.Url),
				AltText:   m.AltText,
				MediaType: strPtr(m.MediaType),
				Position:  intPtr(m.Position),
				IsCover:   boolPtr(m.IsCover),
			})
		}
		e.Media = &media
	}

	return e
}

func intPtr(v int) *int       { return &v }
func strPtr(v string) *string { return &v }
func boolPtr(v bool) *bool    { return &v }

// nullDecimalToString emits the exact stored value. Formatting — thousands
// separators, fixed decimal places — is the client's job; the API's job is not
// to lose a cent on the way.
func nullDecimalToString(d decimal.NullDecimal) *string {
	if !d.Valid {
		return nil
	}
	return strPtr(d.Decimal.String())
}

// decimalToFloat32 is kept only for the quantities where a float is harmless:
// a percentage and physical measurements. No monetary amount goes through it.
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
