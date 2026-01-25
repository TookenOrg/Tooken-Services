package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/globals"
)

// TODO: add user isFavorite
// https://github.com/TookenOrg/Tooken-Services/issues/31
const baseRealEstateQuery = `
SELECT
    q.*,
    CASE
        WHEN q.total_shares IS NULL OR q.total_shares = 0 THEN 0
        ELSE ROUND(
            (q.tokens_sold::numeric / q.total_shares::numeric) * 100,
            2
        )
    END AS tokens_sold_percentage
FROM (
    SELECT 
        re.id, 
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
        conf.payment_frequency_type AS payment_frequency_type_id,
        payt.name AS payment_frequency_type,

        FLOOR(random() * 2001)::int AS tokens_sold

    FROM ass.real_estate re
    LEFT JOIN ass.real_estate_type ret 
        ON re.estate_type = ret.id
    LEFT JOIN ass.real_estate_specification res 
        ON re.id = res.real_estate_id
    LEFT JOIN ass.real_estate_shares_config conf 
        ON re.id = conf.real_estate_id
    LEFT JOIN ass.payment_frequency_type payt 
        ON conf.payment_frequency_type = payt.id
    %s
) q
`

func scanRealEstate(scanner interface {
	Scan(dest ...any) error
}) (server.RealEstate, error) {

	var e server.RealEstate
	e.Configuration = &server.RealEstateConfiguration{}
	e.Specificiation = &server.RealEstateSpecification{}
	e.Progression = &server.RealEstateProgression{}

	var (
		active       sql.NullBool
		createdAt    sql.NullTime
		contractAddr sql.NullString
		estateTypeId sql.NullInt64
		estateType   sql.NullString

		lotSize, surfaceArea, poolSize, terraceSize,
		pricePerShare, yield, tokensSoldPctg sql.NullFloat64

		bedroomNumber, bathroomNumber, builtYear,
		totalShares, paymentFrequency,
		paymentFrequencyTypeId, tokensSold sql.NullInt64

		paymentFrequencyName sql.NullString
	)

	err := scanner.Scan(
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
		&tokensSoldPctg,
	)
	if err != nil {
		return e, err
	}

	// Assignations (exactement ton code, une seule fois)
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
		v := int(estateTypeId.Int64)
		e.EstateTypeId = &v
	}
	if estateType.Valid {
		e.EstateType = &estateType.String
	}

	// Specification
	if lotSize.Valid {
		v := float32(lotSize.Float64)
		e.Specificiation.LotSize = &v
	}
	if surfaceArea.Valid {
		v := float32(surfaceArea.Float64)
		e.Specificiation.SurfaceArea = &v
	}
	if bedroomNumber.Valid {
		v := int(bedroomNumber.Int64)
		e.Specificiation.BedroomNumber = &v
	}
	if bathroomNumber.Valid {
		v := int(bathroomNumber.Int64)
		e.Specificiation.BathroomNumber = &v
	}
	if poolSize.Valid {
		v := float32(poolSize.Float64)
		e.Specificiation.PoolSize = &v
	}
	if terraceSize.Valid {
		v := float32(terraceSize.Float64)
		e.Specificiation.TerraceSize = &v
	}
	if builtYear.Valid {
		v := int(builtYear.Int64)
		e.Specificiation.BuildYear = &v
	}

	// Configuration
	if totalShares.Valid {
		v := int(totalShares.Int64)
		e.Configuration.TotalShares = &v
	}
	if pricePerShare.Valid {
		v := float32(pricePerShare.Float64)
		e.Configuration.PricePerShare = &v
	}
	if yield.Valid {
		v := float32(yield.Float64)
		e.Configuration.Yield = &v
	}
	if paymentFrequency.Valid {
		v := int(paymentFrequency.Int64)
		e.Configuration.PaymentFrequency = &v
	}
	if paymentFrequencyTypeId.Valid {
		v := int(paymentFrequencyTypeId.Int64)
		e.Configuration.PaymentFrequencyTypeId = &v
	}
	if paymentFrequencyName.Valid {
		e.Configuration.PaymentFrequencyType = &paymentFrequencyName.String
	}

	// Progression
	if tokensSold.Valid {
		v := int(tokensSold.Int64)
		e.Progression.TokensSold = &v
	}
	if tokensSoldPctg.Valid {
		v := float32(tokensSoldPctg.Float64)
		e.Progression.TokensSoldPctg = &v
	}

	return e, nil
}

func GetActiveRealEstates(ctx context.Context) ([]server.RealEstate, error) {
	query := fmt.Sprintf(baseRealEstateQuery, "WHERE re.active = true")

	rows, err := globals.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var estates []server.RealEstate
	for rows.Next() {
		e, err := scanRealEstate(rows)
		if err != nil {
			return nil, err
		}
		estates = append(estates, e)
	}

	return estates, rows.Err()
}

func GetRealEstateById(ctx context.Context, id int) (server.RealEstate, error) {
	query := fmt.Sprintf(
		baseRealEstateQuery,
		"WHERE re.id = $1 AND re.active = true",
	)

	row := globals.DB.QueryRowContext(ctx, query, id)
	return scanRealEstate(row)
}

// func GetActiveRealEstates(ctx context.Context) (realEstates []server.RealEstate, err error) {
// 	query := `
// 				SELECT
// 			q.*,

// 			/* % tokens sold */
// 			CASE
// 				WHEN q.total_shares IS NULL OR q.total_shares = 0 THEN 0
// 				ELSE ROUND(
// 					(q.tokens_sold::numeric / q.total_shares::numeric) * 100,
// 					2
// 				)
// 			END AS tokens_sold_percentage

// 		FROM (
// 			SELECT
// 				re.id,
// 				re.title,
// 				re.description,
// 				re.imageurl,
// 				re.estate_type AS estate_type_id,
// 				ret.name AS estate_type,
// 				re.active,
// 				re.created_at,
// 				re.contract_address,

// 				res.lot_size,
// 				res.surface_area,
// 				res.bedroom_number,
// 				res.bathroom_number,
// 				res.pool_size,
// 				res.terrace_size,
// 				res.built_year,

// 				conf.total_shares,
// 				conf.price_per_share,
// 				conf.yield,
// 				conf.payment_frequency,
// 				conf.payment_frequency_type AS payment_frequency_type_id,
// 				payt.name AS payment_frequency_type,

// 				/* Random tokens sold between 0 and 2000 */
// 				FLOOR(random() * 2001)::int AS tokens_sold

// 			FROM ass.real_estate re
// 			LEFT JOIN ass.real_estate_type ret
// 				ON re.estate_type = ret.id
// 			LEFT JOIN ass.real_estate_specification res
// 				ON re.id = res.real_estate_id
// 			LEFT JOIN ass.real_estate_shares_config conf
// 				ON re.id = conf.real_estate_id
// 			LEFT JOIN ass.payment_frequency_type payt
// 				ON conf.payment_frequency_type = payt.id

// 			WHERE re.active = true
// 		) q;
// 	`

// 	rows, err := globals.DB.Query(query)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var e server.RealEstate
// 		e.Configuration = &server.RealEstateConfiguration{}
// 		e.Specificiation = &server.RealEstateSpecification{}
// 		e.Progression = &server.RealEstateProgression{}

// 		var active sql.NullBool
// 		var createdAt sql.NullTime
// 		var contractAddr sql.NullString
// 		var estateTypeId sql.NullInt64
// 		var estateType sql.NullString
// 		var lotSize, surfaceArea, poolSize, terraceSize, pricePerShare, yield, tokensSoldPctg sql.NullFloat64
// 		var bedroomNumber, bathroomNumber, builtYear, totalShares, paymentFrequency, paymentFrequencyTypeId, tokensSold sql.NullInt64
// 		var paymentFrequencyName sql.NullString

// 		err := rows.Scan(
// 			&e.Id,
// 			&e.Title,
// 			&e.Description,
// 			&e.Imageurl,
// 			&estateTypeId,
// 			&estateType,
// 			&active,
// 			&createdAt,
// 			&contractAddr,
// 			&lotSize,
// 			&surfaceArea,
// 			&bedroomNumber,
// 			&bathroomNumber,
// 			&poolSize,
// 			&terraceSize,
// 			&builtYear,
// 			&totalShares,
// 			&pricePerShare,
// 			&yield,
// 			&paymentFrequency,
// 			&paymentFrequencyTypeId,
// 			&paymentFrequencyName,
// 			&tokensSold,
// 			&tokensSoldPctg,
// 		)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		if active.Valid {
// 			e.Active = &active.Bool
// 		}
// 		if createdAt.Valid {
// 			e.CreatedAt = &createdAt.Time
// 		}
// 		if contractAddr.Valid {
// 			e.ContractAddress = &contractAddr.String
// 		}
// 		if estateTypeId.Valid {
// 			val := int(estateTypeId.Int64)
// 			e.EstateTypeId = &val
// 		}
// 		if estateType.Valid {
// 			e.EstateType = &estateType.String
// 		}

// 		// Specification
// 		if lotSize.Valid {
// 			val := float32(lotSize.Float64)
// 			e.Specificiation.LotSize = &val
// 		}
// 		if surfaceArea.Valid {
// 			val := float32(surfaceArea.Float64)
// 			e.Specificiation.SurfaceArea = &val
// 		}
// 		if bedroomNumber.Valid {
// 			val := int(bedroomNumber.Int64)
// 			e.Specificiation.BedroomNumber = &val
// 		}
// 		if bathroomNumber.Valid {
// 			val := int(bathroomNumber.Int64)
// 			e.Specificiation.BathroomNumber = &val
// 		}
// 		if poolSize.Valid {
// 			val := float32(poolSize.Float64)
// 			e.Specificiation.PoolSize = &val
// 		}
// 		if terraceSize.Valid {
// 			val := float32(terraceSize.Float64)
// 			e.Specificiation.TerraceSize = &val
// 		}
// 		if builtYear.Valid {
// 			val := int(builtYear.Int64)
// 			e.Specificiation.BuildYear = &val
// 		}

// 		// Configuration
// 		if totalShares.Valid {
// 			val := int(totalShares.Int64)
// 			e.Configuration.TotalShares = &val
// 		}
// 		if pricePerShare.Valid {
// 			val := float32(pricePerShare.Float64)
// 			e.Configuration.PricePerShare = &val
// 		}
// 		if yield.Valid {
// 			val := float32(yield.Float64)
// 			e.Configuration.Yield = &val
// 		}
// 		if paymentFrequency.Valid {
// 			val := int(paymentFrequency.Int64)
// 			e.Configuration.PaymentFrequency = &val
// 		}
// 		if paymentFrequencyTypeId.Valid {
// 			val := int(paymentFrequencyTypeId.Int64)
// 			e.Configuration.PaymentFrequencyTypeId = &val
// 		}
// 		if paymentFrequencyName.Valid {
// 			e.Configuration.PaymentFrequencyType = &paymentFrequencyName.String
// 		}

// 		// Progression
// 		if tokensSold.Valid {
// 			val := int(tokensSold.Int64)
// 			e.Progression.TokensSold = &val
// 		}

// 		if tokensSoldPctg.Valid {
// 			val := float32(tokensSoldPctg.Float64)
// 			e.Progression.TokensSoldPctg = &val

// 		}

// 		realEstates = append(realEstates, e)
// 	}

// 	if err := rows.Err(); err != nil {
// 		log.Fatal(err)
// 	}
// 	return
// }

// func GetRealEstateById(ctx context.Context, id int) (realEstate server.RealEstate, err error) {
// 	query := `
// 	SELECT
// 		q.*,

// 		/* % tokens sold */
// 		CASE
// 			WHEN q.total_shares IS NULL OR q.total_shares = 0 THEN 0
// 			ELSE ROUND(
// 				(q.tokens_sold::numeric / q.total_shares::numeric) * 100,
// 				2
// 			)
// 		END AS tokens_sold_percentage

// 	FROM (
// 		SELECT
// 			re.id,
// 			re.title,
// 			re.description,
// 			re.imageurl,
// 			re.estate_type AS estate_type_id,
// 			ret.name AS estate_type,
// 			re.active,
// 			re.created_at,
// 			re.contract_address,

// 			res.lot_size,
// 			res.surface_area,
// 			res.bedroom_number,
// 			res.bathroom_number,
// 			res.pool_size,
// 			res.terrace_size,
// 			res.built_year,

// 			conf.total_shares,
// 			conf.price_per_share,
// 			conf.yield,
// 			conf.payment_frequency,
// 			conf.payment_frequency_type AS payment_frequency_type_id,
// 			payt.name AS payment_frequency_type,

// 			/* Random tokens sold between 0 and 2000 */
// 			FLOOR(random() * 2001)::int AS tokens_sold

// 		FROM ass.real_estate re
// 		LEFT JOIN ass.real_estate_type ret
// 			ON re.estate_type = ret.id
// 		LEFT JOIN ass.real_estate_specification res
// 			ON re.id = res.real_estate_id
// 		LEFT JOIN ass.real_estate_shares_config conf
// 			ON re.id = conf.real_estate_id
// 		LEFT JOIN ass.payment_frequency_type payt
// 			ON conf.payment_frequency_type = payt.id

// 		WHERE re.id = $1 AND re.active = true
// 	) q;
// 	`

// 	// On utilise QueryRow avec le paramètre id
// 	row := globals.DB.QueryRowContext(ctx, query, id)

// 	// Initialisation des sous-structs
// 	realEstate.Configuration = &server.RealEstateConfiguration{}
// 	realEstate.Specificiation = &server.RealEstateSpecification{}
// 	realEstate.Progression = &server.RealEstateProgression{}

// 	// Variables Null* pour le scan
// 	var active sql.NullBool
// 	var createdAt sql.NullTime
// 	var contractAddr sql.NullString
// 	var estateTypeId sql.NullInt64
// 	var estateType sql.NullString
// 	var lotSize, surfaceArea, poolSize, terraceSize, pricePerShare, yield, tokensSoldPctg sql.NullFloat64
// 	var bedroomNumber, bathroomNumber, builtYear, totalShares, paymentFrequency, paymentFrequencyTypeId, tokensSold sql.NullInt64
// 	var paymentFrequencyName sql.NullString

// 	err = row.Scan(
// 		&realEstate.Id,
// 		&realEstate.Title,
// 		&realEstate.Description,
// 		&realEstate.Imageurl,
// 		&estateTypeId,
// 		&estateType,
// 		&active,
// 		&createdAt,
// 		&contractAddr,
// 		&lotSize,
// 		&surfaceArea,
// 		&bedroomNumber,
// 		&bathroomNumber,
// 		&poolSize,
// 		&terraceSize,
// 		&builtYear,
// 		&totalShares,
// 		&pricePerShare,
// 		&yield,
// 		&paymentFrequency,
// 		&paymentFrequencyTypeId,
// 		&paymentFrequencyName,
// 		&tokensSold,
// 		&tokensSoldPctg,
// 	)
// 	if err != nil {
// 		logger.LogError("%s", err.Error())
// 		return
// 	}

// 	// Assignation sécurisée
// 	if active.Valid {
// 		realEstate.Active = &active.Bool
// 	}
// 	if createdAt.Valid {
// 		realEstate.CreatedAt = &createdAt.Time
// 	}
// 	if contractAddr.Valid {
// 		realEstate.ContractAddress = &contractAddr.String
// 	}
// 	if estateTypeId.Valid {
// 		val := int(estateTypeId.Int64)
// 		realEstate.EstateTypeId = &val
// 	}
// 	if estateType.Valid {
// 		realEstate.EstateType = &estateType.String
// 	}

// 	// Specification
// 	if lotSize.Valid {
// 		val := float32(lotSize.Float64)
// 		realEstate.Specificiation.LotSize = &val
// 	}
// 	if surfaceArea.Valid {
// 		val := float32(surfaceArea.Float64)
// 		realEstate.Specificiation.SurfaceArea = &val
// 	}
// 	if bedroomNumber.Valid {
// 		val := int(bedroomNumber.Int64)
// 		realEstate.Specificiation.BedroomNumber = &val
// 	}
// 	if bathroomNumber.Valid {
// 		val := int(bathroomNumber.Int64)
// 		realEstate.Specificiation.BathroomNumber = &val
// 	}
// 	if poolSize.Valid {
// 		val := float32(poolSize.Float64)
// 		realEstate.Specificiation.PoolSize = &val
// 	}
// 	if terraceSize.Valid {
// 		val := float32(terraceSize.Float64)
// 		realEstate.Specificiation.TerraceSize = &val
// 	}
// 	if builtYear.Valid {
// 		val := int(builtYear.Int64)
// 		realEstate.Specificiation.BuildYear = &val
// 	}

// 	// Configuration
// 	if totalShares.Valid {
// 		val := int(totalShares.Int64)
// 		realEstate.Configuration.TotalShares = &val
// 	}
// 	if pricePerShare.Valid {
// 		val := float32(pricePerShare.Float64)
// 		realEstate.Configuration.PricePerShare = &val
// 	}
// 	if yield.Valid {
// 		val := float32(yield.Float64)
// 		realEstate.Configuration.Yield = &val
// 	}
// 	if paymentFrequency.Valid {
// 		val := int(paymentFrequency.Int64)
// 		realEstate.Configuration.PaymentFrequency = &val
// 	}
// 	if paymentFrequencyTypeId.Valid {
// 		val := int(paymentFrequencyTypeId.Int64)
// 		realEstate.Configuration.PaymentFrequencyTypeId = &val
// 	}
// 	if paymentFrequencyName.Valid {
// 		realEstate.Configuration.PaymentFrequencyType = &paymentFrequencyName.String
// 	}

// 	// Progression
// 	if tokensSold.Valid {
// 		val := int(tokensSold.Int64)
// 		realEstate.Progression.TokensSold = &val
// 	}
// 	if tokensSoldPctg.Valid {
// 		val := float32(tokensSoldPctg.Float64)
// 		realEstate.Progression.TokensSoldPctg = &val
// 	}

// 	return
// }
