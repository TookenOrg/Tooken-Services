package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/auth/utils"
	usersDb "github.com/TookenOrg/tooken-services/internal/users/database"
)

// ErrInvalidRefreshToken: the refresh token is forged, expired, not a refresh
// token, or names an account that no longer exists. Either way the user has
// to sign in again.
var ErrInvalidRefreshToken = errors.New("invalid refresh token")

// Refresh issues a new access token for the session the refresh token belongs
// to. The account is read again, so the new token carries the current role,
// and its lifetime is capped at the refresh token's own expiry so the session
// cannot outlive its eight-hour ceiling.
func (s *Service) Refresh(ctx context.Context, refreshToken string) (session server.JwtToken, err error) {
	claims, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		return server.JwtToken{}, fmt.Errorf("%w: %w", ErrInvalidRefreshToken, err)
	}

	var maxExpiresAt time.Time
	if claims.ExpiresAt != nil {
		maxExpiresAt = claims.ExpiresAt.Time
	}

	session, err = issueAccessToken(ctx, claims.UserID, refreshToken, maxExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return server.JwtToken{}, fmt.Errorf("%w: account %d no longer exists", ErrInvalidRefreshToken, claims.UserID)
	}
	if errors.Is(err, utils.ErrSessionExpired) {
		return server.JwtToken{}, fmt.Errorf("%w: %w", ErrInvalidRefreshToken, err)
	}

	return session, err
}

// openSession starts the session of a user who just proved who they are.
func openSession(ctx context.Context, userID int) (session server.JwtToken, err error) {
	refreshToken, err := utils.GenerateRefreshToken(userID)
	if err != nil {
		return server.JwtToken{}, fmt.Errorf("cannot sign the refresh token: %w", err)
	}

	return issueAccessToken(ctx, userID, refreshToken, time.Time{})
}

// issueAccessToken describes the account the way GET /users/me does, so the
// client gets the same user wherever it reads it.
func issueAccessToken(ctx context.Context, userID int, refreshToken string, maxExpiresAt time.Time) (session server.JwtToken, err error) {
	user, err := usersDb.GetUserByID(ctx, userID)
	if err != nil {
		return server.JwtToken{}, fmt.Errorf("cannot read account %d: %w", userID, err)
	}
	projectSessionUserKycStatus(user, time.Now())

	var role string
	if user.Role != nil {
		role = string(*user.Role)
	}

	token, expiresAt, err := utils.NewAccessTokenUntil(user.Id, string(user.Email), role, maxExpiresAt)
	if err != nil {
		return server.JwtToken{}, fmt.Errorf("cannot sign the access token: %w", err)
	}

	return server.JwtToken{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         *user,
	}, nil
}

func projectSessionUserKycStatus(user *server.User, now time.Time) {
	if user == nil || user.KycStatus == nil {
		return
	}
	if user.KycExpiresAt != nil && user.KycExpiresAt.Before(now) {
		expired := server.UserKycStatusExpired
		user.KycStatus = &expired
	}
}
