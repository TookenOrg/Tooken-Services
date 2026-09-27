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
	"github.com/samber/lo"
)

var (
	ErrKycVerificationNotFound      = errors.New("KYC verification not found")
	ErrMessageKycAlreadyExists      = errors.New("a KYC verification with status 'submitted' already exists for this user")
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

	// 2 - Get Kyc verification to check if one already exists for the user with status "submitted"
	kycVerifications, err := database.GetKycVerificationByUserIDAndStatus(ctx, &userId, lo.ToPtr(database.StatusSubmitted))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}
	if len(kycVerifications) > 0 {
		return ErrMessageKycAlreadyExists
	}

	// 3 - Insert in kyc_verification with status submitted
	_, _, err = database.InsertKycVerification(ctx, userId, request)
	if err != nil {
		return
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
