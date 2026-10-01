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
	StatusApproved  = "approved"
	StatusVerified  = "verified"
	StatusRejected  = "rejected"
	StatusRevoked   = "revoked"
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
		SELECT k.id,
				k.user_id,
				k.status,
				k.declared_full_name,
				k.declared_country_code,
				k.created_at,
				k.updated_at,
				u.email,
				k.submitted_at,
				k.decided_at,
				k.decided_by,
				k.rejection_reason,
				k.expires_at,
				k.revoked_at,
				k.revoked_by,
				k.revocation_reason
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
			&kv.SubmittedAt,
			&kv.DecidedAt,
			&kv.DecidedBy,
			&kv.RejectionReason,
			&kv.ExpiresAt,
			&kv.RevokedAt,
			&kv.RevokedBy,
			&kv.RevocationReason,
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

// GetLatestKycVerificationByUserID returns the most recent filing of a user, or
// sql.ErrNoRows when that user has never filed.
//
// The ordering deliberately mirrors usr.kyc_verification_sync_user() (migration
// 000025). The guard that accepts a new filing and the trigger that projects
// usr.users must read the same row; were they to disagree, the API would allow
// a filing the projection then interprets differently.
func GetLatestKycVerificationByUserID(ctx context.Context, userId int) (kycVerification server.KycVerification, err error) {
	query := `
		SELECT 	k.id,
				k.user_id,
				k.status,
				k.declared_full_name,
				k.declared_country_code,
				k.created_at,
				k.updated_at,
				u.email,
				k.submitted_at,
				k.decided_at,
				k.decided_by,
				k.rejection_reason,
				k.expires_at,
				k.revoked_at,
				k.revoked_by,
				k.revocation_reason
		FROM usr.kyc_verification k
		JOIN usr.users u ON u.id = k.user_id
		WHERE k.user_id = $1
		ORDER BY k.submitted_at DESC, k.id DESC
		LIMIT 1;
	`
	err = globals.DB.QueryRowContext(ctx, query, userId).Scan(
		&kycVerification.Id,
		&kycVerification.UserId,
		&kycVerification.Status,
		&kycVerification.DeclaredFullName,
		&kycVerification.DeclaredCountryCode,
		&kycVerification.CreatedAt,
		&kycVerification.UpdatedAt,
		&kycVerification.Email,
		&kycVerification.SubmittedAt,
		&kycVerification.DecidedAt,
		&kycVerification.DecidedBy,
		&kycVerification.RejectionReason,
		&kycVerification.ExpiresAt,
		&kycVerification.RevokedAt,
		&kycVerification.RevokedBy,
		&kycVerification.RevocationReason,
	)
	if err != nil {
		return kycVerification, fmt.Errorf("failed to get latest kyc verification: %w", err)
	}
	return kycVerification, nil
}

func GetKycVerificationByID(ctx context.Context, verificationId int) (kycVerification server.KycVerification, err error) {
	query := `
		SELECT 	k.id,
				k.user_id,
				k.status,
				k.declared_full_name,
				k.declared_country_code,
				k.created_at,
				k.updated_at,
				u.email,
				k.submitted_at,
				k.decided_at,
				k.decided_by,
				k.rejection_reason,
				k.expires_at,
				k.revoked_at,
				k.revoked_by,
				k.revocation_reason,
				u.kyc_status,
				u.kyc_verified_at
		FROM usr.kyc_verification k
		JOIN usr.users u ON u.id = k.user_id
		WHERE k.id = $1;
	`
	err = globals.DB.QueryRowContext(ctx, query, verificationId).Scan(
		&kycVerification.Id,
		&kycVerification.UserId,
		&kycVerification.Status,
		&kycVerification.DeclaredFullName,
		&kycVerification.DeclaredCountryCode,
		&kycVerification.CreatedAt,
		&kycVerification.UpdatedAt,
		&kycVerification.Email,
		&kycVerification.SubmittedAt,
		&kycVerification.DecidedAt,
		&kycVerification.DecidedBy,
		&kycVerification.RejectionReason,
		&kycVerification.ExpiresAt,
		&kycVerification.RevokedAt,
		&kycVerification.RevokedBy,
		&kycVerification.RevocationReason,
		&kycVerification.KycStatus,
		&kycVerification.KycVerifiedAt,
	)
	if err != nil {
		return kycVerification, fmt.Errorf("failed to get kyc verification by ID: %w", err)
	}
	return kycVerification, nil
}

func ApproveKycVerification(ctx context.Context, verificationId int, expiresAt time.Time, decidedBy int) (kycVerification server.KycVerification, err error) {
	query := `
		UPDATE usr.kyc_verification
		SET status = 'approved',
			expires_at = $2,
			decided_at = NOW(),
			decided_by = $3
		WHERE id = $1
		RETURNING id,
				  user_id,
				  status,
				  declared_full_name,
				  declared_country_code,
				  created_at,
				  updated_at,
				  (SELECT email FROM usr.users WHERE id = usr.kyc_verification.user_id),
				  submitted_at,
				  decided_at,
				  decided_by,
				  rejection_reason,
				  expires_at,
				  revoked_at,
				  revoked_by,
				  revocation_reason;
	`
	err = globals.DB.QueryRowContext(ctx, query, verificationId, expiresAt, decidedBy).Scan(
		&kycVerification.Id,
		&kycVerification.UserId,
		&kycVerification.Status,
		&kycVerification.DeclaredFullName,
		&kycVerification.DeclaredCountryCode,
		&kycVerification.CreatedAt,
		&kycVerification.UpdatedAt,
		&kycVerification.Email,
		&kycVerification.SubmittedAt,
		&kycVerification.DecidedAt,
		&kycVerification.DecidedBy,
		&kycVerification.RejectionReason,
		&kycVerification.ExpiresAt,
		&kycVerification.RevokedAt,
		&kycVerification.RevokedBy,
		&kycVerification.RevocationReason,
	)
	if err != nil {
		return kycVerification, fmt.Errorf("failed to approve kyc verification: %w", err)
	}
	return kycVerification, nil
}

func RejectKycVerification(ctx context.Context, verificationId int, reason string, decidedBy int) (kycVerification server.KycVerification, err error) {
	query := `
		UPDATE usr.kyc_verification
		SET status = 'rejected',
			rejection_reason = $2,
			decided_at = NOW(),
			decided_by = $3
		WHERE id = $1
		RETURNING id,
				  user_id,
				  status,
				  declared_full_name,
				  declared_country_code,
				  created_at,
				  updated_at,
				  (SELECT email FROM usr.users WHERE id = usr.kyc_verification.user_id),
				  submitted_at,
				  decided_at,
				  decided_by,
				  rejection_reason,
				  expires_at,
				  revoked_at,
				  revoked_by,
				  revocation_reason;
	`
	err = globals.DB.QueryRowContext(ctx, query, verificationId, reason, decidedBy).Scan(
		&kycVerification.Id,
		&kycVerification.UserId,
		&kycVerification.Status,
		&kycVerification.DeclaredFullName,
		&kycVerification.DeclaredCountryCode,
		&kycVerification.CreatedAt,
		&kycVerification.UpdatedAt,
		&kycVerification.Email,
		&kycVerification.SubmittedAt,
		&kycVerification.DecidedAt,
		&kycVerification.DecidedBy,
		&kycVerification.RejectionReason,
		&kycVerification.ExpiresAt,
		&kycVerification.RevokedAt,
		&kycVerification.RevokedBy,
		&kycVerification.RevocationReason,
	)
	if err != nil {
		return kycVerification, fmt.Errorf("failed to reject kyc verification: %w", err)
	}
	return kycVerification, nil
}

func RevokeKycVerification(ctx context.Context, verificationId int, reason string, revokedBy int) (kycVerification server.KycVerification, err error) {
	query := `
		UPDATE usr.kyc_verification
		SET status = 'revoked',
			revocation_reason = $2,
			revoked_at = NOW(),
			revoked_by = $3
		WHERE id = $1
		RETURNING id,
				  user_id,
				  status,
				  declared_full_name,
				  declared_country_code,
				  created_at,
				  updated_at,
				  (SELECT email FROM usr.users WHERE id = usr.kyc_verification.user_id),
				  submitted_at,
				  decided_at,
				  decided_by,
				  rejection_reason,
				  expires_at,
				  revoked_at,
				  revoked_by,
				  revocation_reason;
	`
	err = globals.DB.QueryRowContext(ctx, query, verificationId, reason, revokedBy).Scan(
		&kycVerification.Id,
		&kycVerification.UserId,
		&kycVerification.Status,
		&kycVerification.DeclaredFullName,
		&kycVerification.DeclaredCountryCode,
		&kycVerification.CreatedAt,
		&kycVerification.UpdatedAt,
		&kycVerification.Email,
		&kycVerification.SubmittedAt,
		&kycVerification.DecidedAt,
		&kycVerification.DecidedBy,
		&kycVerification.RejectionReason,
		&kycVerification.ExpiresAt,
		&kycVerification.RevokedAt,
		&kycVerification.RevokedBy,
		&kycVerification.RevocationReason,
	)
	if err != nil {
		return kycVerification, fmt.Errorf("failed to revoke kyc verification: %w", err)
	}
	return kycVerification, nil
}
