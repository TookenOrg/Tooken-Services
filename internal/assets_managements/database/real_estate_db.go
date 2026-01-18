package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/globals"
)

func GetActiveRealEstates(ctx context.Context) (realEstates []server.RealEstate, err error) {
	query := `
		SELECT re.id, 
			re.title, 
			re.description, 
			re.imageurl, 
			re.estate_type AS estate_type_id,
			ret.name AS estate_type, 
			re.active, 
			re.created_at, 
			re.contract_address,
			res.lot_size, 
			res.surface_area, 
			res.bedroom_number, 
			res.bathroom_number, 
			res.pool_size, 
			res.terrace_size,
			res.built_year,
			conf.total_shares, 
			conf.price_per_share, 
			conf.yield, 
			conf.payment_frequency, 
			conf.payment_frequency_type as payment_frequency_type_id,
			payt.name as payment_frequency_type,
			0 as token_sold
		FROM ass.real_estate re
		LEFT JOIN ass.real_estate_type ret ON re.estate_type = ret.id
		LEFT JOIN ass.real_estate_specification res ON re.id = res.real_estate_id
		LEFT JOIN ass.real_estate_shares_config conf ON re.id = conf.real_estate_id
		LEFT JOIN ass.payment_frequency_type payt ON conf.payment_frequency_type = payt.id
		WHERE re.active = true;
	`

	rows, err := globals.DB.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var e server.RealEstate
		e.Configuration = &server.RealEstateConfiguration{}
		e.Specificiation = &server.RealEstateSpecification{}
		e.Progression = &server.RealEstateProgression{}

		var active sql.NullBool
		var createdAt sql.NullTime
		var contractAddr sql.NullString
		var estateTypeId sql.NullInt64
		var estateType sql.NullString
		var lotSize, surfaceArea, poolSize, terraceSize, pricePerShare, yield sql.NullFloat64
		var bedroomNumber, bathroomNumber, builtYear, totalShares, paymentFrequency, paymentFrequencyTypeId, tokensSold sql.NullInt64
		var paymentFrequencyName sql.NullString

		err := rows.Scan(
			&e.Id,
			&e.Title,
			&e.Description,
			&e.Imageurl,
			&estateTypeId,
			&estateType,
			&active,
			&createdAt,
			&contractAddr,
			&lotSize,
			&surfaceArea,
			&bedroomNumber,
			&bathroomNumber,
			&poolSize,
			&terraceSize,
			&builtYear,
			&totalShares,
			&pricePerShare,
			&yield,
			&paymentFrequency,
			&paymentFrequencyTypeId,
			&paymentFrequencyName,
			&tokensSold,
		)
		if err != nil {
			log.Fatal(err)
		}

		if active.Valid {
			e.Active = &active.Bool
		}
		if createdAt.Valid {
			e.CreatedAt = &createdAt.Time
		}
		if contractAddr.Valid {
			e.ContractAddress = &contractAddr.String
		}
		if estateTypeId.Valid {
			val := int(estateTypeId.Int64)
			e.EstateTypeId = &val
		}
		if estateType.Valid {
			e.EstateType = &estateType.String
		}

		// Specification
		if lotSize.Valid {
			val := float32(lotSize.Float64)
			e.Specificiation.LotSize = &val
		}
		if surfaceArea.Valid {
			val := float32(surfaceArea.Float64)
			e.Specificiation.SurfaceArea = &val
		}
		if bedroomNumber.Valid {
			val := int(bedroomNumber.Int64)
			e.Specificiation.BedroomNumber = &val
		}
		if bathroomNumber.Valid {
			val := int(bathroomNumber.Int64)
			e.Specificiation.BathroomNumber = &val
		}
		if poolSize.Valid {
			val := float32(poolSize.Float64)
			e.Specificiation.PoolSize = &val
		}
		if terraceSize.Valid {
			val := float32(terraceSize.Float64)
			e.Specificiation.TerraceSize = &val
		}
		if builtYear.Valid {
			val := int(builtYear.Int64)
			e.Specificiation.BuildYear = &val
		}

		// Configuration
		if totalShares.Valid {
			val := int(totalShares.Int64)
			e.Configuration.TotalShares = &val
		}
		if pricePerShare.Valid {
			val := float32(pricePerShare.Float64)
			e.Configuration.PricePerShare = &val
		}
		if yield.Valid {
			val := float32(yield.Float64)
			e.Configuration.Yield = &val
		}
		if paymentFrequency.Valid {
			val := int(paymentFrequency.Int64)
			e.Configuration.PaymentFrequency = &val
		}
		if paymentFrequencyTypeId.Valid {
			val := int(paymentFrequencyTypeId.Int64)
			e.Configuration.PaymentFrequencyTypeId = &val
		}
		if paymentFrequencyName.Valid {
			e.Configuration.PaymentFrequencyType = &paymentFrequencyName.String
		}

		// Progression
		if tokensSold.Valid {
			val := int(tokensSold.Int64)
			e.Progression.SharesSold = &val
		}

		realEstates = append(realEstates, e)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return
}
