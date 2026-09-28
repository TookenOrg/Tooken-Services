package services

import (
	"context"
	"database/sql"
	"time"

	"errors"
	"strings"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/users/database"
	"github.com/TookenOrg/tooken-services/internal/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
)

var (
	ErrKycVerificationNotFound      = errors.New("KYC verification not found")
	ErrMessageKycAlreadyExists      = errors.New("a KYC verification with status 'submitted' already exists for this user")
	ErrMessageKycStillValid         = errors.New("a valid KYC verification already exists for this user")
	ErrMessageInvalidCountryCode    = errors.New("invalid country code")
	ErrMessageInvalidFullName       = errors.New("invalid full name")
	ErrMessageInvalidExpiresAt      = errors.New("invalid expires_at")
	ErrMessageInvalidStatus         = errors.New("invalid status for the KYC verification")
	ErrMessageInvalidReason         = errors.New("invalid reason for rejecting the KYC verification")
	ErrMessageInvalidWalletAddress  = errors.New("invalid wallet address for the user")
	ErrMessageIncoherentUserState   = errors.New("incoherent user state: missing wallet but has on-chain identity")
	ErrMessageFailedToVerifyOnChain = errors.New("failed to verify KYC on-chain")
)

func (s *Service) PostKycVerifications(ctx context.Context, userId int, request *server.KycVerificationRequest) (err error) {

	// 0 - Trim and uppercase the declared country code and declared full name
	request.DeclaredFullName = strings.TrimSpace(request.DeclaredFullName)
	if request.DeclaredFullName == "" {
		return ErrMessageInvalidFullName
	}

	request.DeclaredCountryCode = strings.TrimSpace(strings.ToUpper(request.DeclaredCountryCode))
	if request.DeclaredCountryCode == "" || len(request.DeclaredCountryCode) != 2 || !utils.IsAlphaCountryCode(request.DeclaredCountryCode) {
		return ErrMessageInvalidCountryCode
	}

	// 1 - Get user to check if they exist
	_, err = database.GetUserByID(ctx, userId)
	if err != nil {
		return
	}

	// 2 - Refuse a new filing while the previous one still stands
	latestKyc, latestErr := database.GetLatestKycVerificationByUserID(ctx, userId)
	if latestErr != nil && !errors.Is(latestErr, sql.ErrNoRows) {
		return latestErr
	}
	var latest *server.KycVerification
	if latestErr == nil {
		latest = &latestKyc
	}
	if err = canAcceptNewKycFiling(latest, time.Now()); err != nil {
		return err
	}

	// 3 - Insert in kyc_verification with status submitted
	_, _, err = database.InsertKycVerification(ctx, userId, request)
	if err != nil {
		return
	}

	return nil
}

// canAcceptNewKycFiling decides whether a user may file a new KYC verification,
// given their latest one. latest is nil when the user has never filed.
//
// Filing again must stay possible after a rejection, a revocation or an expiry:
// a single refusal must not lock an investor out for good. What must not happen
// is a filing landing on top of an entitlement that still holds. The projection
// trigger of migration 000025 reads the latest submission, so a fresh
// 'submitted' row drags usr.users back to 'pending' and clears kyc_expires_at —
// while the chain still verifies that investor. Milestone 3 gates every order on
// kyc_status = 'verified', so the investor would be silently refused.
func canAcceptNewKycFiling(latest *server.KycVerification, now time.Time) error {
	if latest == nil {
		return nil
	}

	switch latest.Status {
	case database.StatusSubmitted:
		return ErrMessageKycAlreadyExists
	case database.StatusApproved:
		// A missing expiry is an anomaly, never a licence to file again:
		// approving always sets one.
		if latest.ExpiresAt == nil || latest.ExpiresAt.After(now) {
			return ErrMessageKycStillValid
		}
	}

	return nil
}

func (s *Service) GetListKycVerifications(ctx context.Context, userId *int, status *string) (kycVerifications []server.KycVerification, err error) {
	return database.GetKycVerificationByUserIDAndStatus(ctx, userId, status)
}
func (s *Service) GetKycVerificationById(ctx context.Context, verificationId int) (kycVerification server.KycVerification, err error) {
	kycVerification, err = database.GetKycVerificationByID(ctx, verificationId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return kycVerification, ErrKycVerificationNotFound
		}
		return kycVerification, err
	}
	return kycVerification, nil
}

func (s *Service) ApproveKycVerification(ctx context.Context, verificationId int, expiresAt time.Time, decidedBy int) (kycVerification server.KycVerification, err error) {

	// 0 - Validate the expiresAt parameter
	if expiresAt.IsZero() {
		return kycVerification, ErrMessageInvalidExpiresAt
	}

	// 1 - Check if the expiresAt parameter is in the future
	if expiresAt.Before(time.Now()) {
		return kycVerification, ErrMessageInvalidExpiresAt
	}

	// 2 - Check current status of the KYC verification to ensure it can be approved
	currentKycVerification, err := database.GetKycVerificationByID(ctx, verificationId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return kycVerification, ErrKycVerificationNotFound
		}
		return kycVerification, err
	}
	if currentKycVerification.Status != database.StatusSubmitted {
		return kycVerification, ErrMessageInvalidStatus
	}

	// 3 - Approve the KYC verification
	kycVerification, err = database.ApproveKycVerification(ctx, verificationId, expiresAt, decidedBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return kycVerification, ErrKycVerificationNotFound
		}
		return kycVerification, err
	}
	return kycVerification, nil
}

func (s *Service) RejectKycVerification(ctx context.Context, verificationId int, reason string, decidedBy int) (kycVerification server.KycVerification, err error) {

	reason = strings.TrimSpace(reason)
	// 0 - Validate the reason parameter
	if reason == "" {
		return kycVerification, ErrMessageInvalidReason
	}

	// 1 - Check current status of the KYC verification to ensure it can be rejected
	currentKycVerification, err := database.GetKycVerificationByID(ctx, verificationId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return kycVerification, ErrKycVerificationNotFound
		}
		return kycVerification, err
	}
	if currentKycVerification.Status != database.StatusSubmitted {
		return kycVerification, ErrMessageInvalidStatus
	}

	kycVerification, err = database.RejectKycVerification(ctx, verificationId, reason, decidedBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return kycVerification, ErrKycVerificationNotFound
		}
		return kycVerification, err
	}
	return kycVerification, nil
}

func (s *Service) RevokeKycVerification(ctx context.Context, verificationId int, reason string, revokedBy int) (kycVerification server.KycVerification, err error) {

	reason = strings.TrimSpace(reason)
	// 0 - Validate the reason parameter
	if reason == "" {
		return kycVerification, ErrMessageInvalidReason
	}

	// 1 - Check current status of the KYC verification to ensure it can be revoked
	currentKycVerification, err := database.GetKycVerificationByID(ctx, verificationId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return kycVerification, ErrKycVerificationNotFound
		}
		return kycVerification, err
	}
	if currentKycVerification.Status != database.StatusApproved {
		return kycVerification, ErrMessageInvalidStatus
	}

	kycVerification, err = database.RevokeKycVerification(ctx, verificationId, reason, revokedBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return kycVerification, ErrKycVerificationNotFound
		}
		return kycVerification, err
	}
	return kycVerification, nil
}

func (s *Service) SyncApprovedKycVerificationOnChain(ctx context.Context, verificationId int) (response server.KycOnChainSyncResponse, err error) {

	// 0 - Check current status of the KYC verification to ensure it can be synced on-chain
	currentKycVerification, err := database.GetKycVerificationByID(ctx, verificationId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response, ErrKycVerificationNotFound
		}
		return response, err
	}
	if currentKycVerification.Status != database.StatusApproved {
		return response, ErrMessageInvalidStatus
	}
	if !currentKycVerification.ExpiresAt.After(time.Now()) {
		return response, ErrMessageInvalidExpiresAt
	}
	if strings.TrimSpace(currentKycVerification.DeclaredCountryCode) == "" {
		return response, ErrMessageInvalidCountryCode
	}

	response, isVerifiedOnChain, err := s.getKYCOnChainSync().SyncApprovedKYCOnChain(ctx, currentKycVerification)
	if err != nil {
		return response, err
	}
	if !isVerifiedOnChain {
		return response, ErrMessageFailedToVerifyOnChain
	}

	if currentKycVerification.KycStatus != nil &&
		*currentKycVerification.KycStatus == database.StatusVerified &&
		currentKycVerification.KycVerifiedAt != nil {
		response.KycStatus = database.StatusVerified
		response.KycVerifiedAt = *currentKycVerification.KycVerifiedAt
		return response, nil
	}

	logger.LogInfo("KYC verification ID=%d has been successfully verified on-chain", currentKycVerification.Id)
	verifiedAt, err := database.MarkUserKycVerified(ctx, currentKycVerification.UserId)
	if err != nil {
		return response, err
	}
	response.KycStatus = database.StatusVerified
	response.KycVerifiedAt = verifiedAt

	return response, nil
}
