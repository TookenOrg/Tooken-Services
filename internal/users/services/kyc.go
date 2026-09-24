package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/users/database"
	"github.com/TookenOrg/tooken-services/internal/utils"
)

var (
	ErrMessageKycAlreadyExists   = errors.New("a KYC verification with status 'submitted' already exists for this user")
	ErrMessageInvalidCountryCode = errors.New("invalid country code")
	ErrMessageInvalidFullName    = errors.New("invalid full name")
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
	kycVerification, err := database.GetKycVerificationByUserIDAndStatus(ctx, userId, database.StatusSubmitted)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}
	if kycVerification != nil {
		return ErrMessageKycAlreadyExists
	}

	// 3 - Insert in kyc_verification with status submitted
	_, _, err = database.InsertKycVerification(ctx, userId, request)
	if err != nil {
		return
	}

	return nil
}
