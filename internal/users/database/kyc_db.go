package database

import (
	"context"
	"fmt"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
)

const (
	StatusSubmitted = "submitted"
)

func InsertKycVerification(ctx context.Context, userId int, request *server.KycVerificationRequest) (id int, createdAt time.Time, err error) {
	query := `
		INSERT INTO usr.kyc_verification (
			user_id,
			status,
			declared_full_name,
			declared_country_code,
			submitted_at
		)
		VALUES (
			$1,  
			$2, 
			$3,  
			$4,
			$5
		)
		RETURNING id, created_at;
    `
	err = globals.DB.QueryRowContext(ctx, query, userId, StatusSubmitted, request.DeclaredFullName, request.DeclaredCountryCode, time.Now()).Scan(&id, &createdAt)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to insert kyc verification: %w", err)
	}

	logger.LogInfo("KYC verification inserted. ID=%d", id)
	return id, createdAt, nil
}

func GetKycVerificationByUserIDAndStatus(ctx context.Context, userId *int, status *string) (kycVerification []server.KycVerification, err error) {
	// A nil filter becomes NULL and is ignored.
	query := `
		SELECT k.id, k.user_id, k.status, k.declared_full_name, k.declared_country_code, k.created_at, k.updated_at, u.email
		FROM usr.kyc_verification k
		JOIN usr.users u ON u.id = k.user_id
		WHERE ($1::integer IS NULL OR k.user_id = $1)
		  AND ($2::text IS NULL OR k.status = $2);
	`
	kycVerification = []server.KycVerification{}
	rows, err := globals.DB.QueryContext(ctx, query, userId, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get kyc verification: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var kv server.KycVerification
		if err := rows.Scan(
			&kv.Id,
			&kv.UserId,
			&kv.Status,
			&kv.DeclaredFullName,
			&kv.DeclaredCountryCode,
			&kv.CreatedAt,
			&kv.UpdatedAt,
			&kv.Email,
		); err != nil {
			return nil, fmt.Errorf("failed to scan kyc verification: %w", err)
		}
		kycVerification = append(kycVerification, kv)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate kyc verification rows: %w", err)
	}

	return kycVerification, nil
}
