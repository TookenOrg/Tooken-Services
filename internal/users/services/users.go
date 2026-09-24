package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/users/database"
	"github.com/samber/lo"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

func (s *Service) GetUserByID(ctx context.Context, userID int) (user *server.User, err error) {
	user, err = database.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrUserNotFound
		}
		return
	}
	if user == nil {
		err = ErrUserNotFound
		return
	}
	user.KycStatus = defineKycStatus(*user)
	return
}

func defineKycStatus(user server.User) (status *server.UserKycStatus) {
	if user.KycStatus == nil {
		return nil
	}

	isExpired := user.KycExpiresAt != nil && user.KycExpiresAt.Before(time.Now())
	if isExpired {
		status = lo.ToPtr(server.UserKycStatus(server.Expired))
	} else {
		status = user.KycStatus
	}
	return
}
