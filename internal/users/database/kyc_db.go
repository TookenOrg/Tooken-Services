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

type KycVerificationDTO struct {
	Id                  int
	UserId              int
	Status              string
	DeclaredFullName    string
	DeclaredCountryCode string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

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

func GetKycVerificationByUserIDAndStatus(ctx context.Context, userId int, status string) (kycVerification *KycVerificationDTO, err error) {
	query := `
		SELECT id, user_id, status, declared_full_name, declared_country_code, created_at, updated_at
		FROM usr.kyc_verification
		WHERE user_id = $1 AND status = $2
		LIMIT 1;
	`
	kycVerification = &KycVerificationDTO{}
	err = globals.DB.QueryRowContext(ctx, query, userId, status).Scan(
		&kycVerification.Id,
		&kycVerification.UserId,
		&kycVerification.Status,
		&kycVerification.DeclaredFullName,
		&kycVerification.DeclaredCountryCode,
		&kycVerification.CreatedAt,
		&kycVerification.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get kyc verification: %w", err)
	}
	return kycVerification, nil
}
