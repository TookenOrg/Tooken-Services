package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/users/database"
)

var (
	ErrInvalidRole   = errors.New("role must be USER, MANAGER or ADMIN")
	ErrRoleForbidden = database.ErrRoleForbidden
	ErrRoleConflict  = database.ErrRoleConflict
)

func validUserRole(role server.UserRole) bool {
	return role == "USER" || role == "MANAGER" || role == "ADMIN"
}

func (s *Service) UpdateUserRole(ctx context.Context, actorID, userID int, role server.UserRole) (*server.User, error) {
	if !validUserRole(role) {
		return nil, ErrInvalidRole
	}
	user, err := database.UpdateUserRole(ctx, actorID, userID, role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	user.KycStatus = defineKycStatus(*user)
	return user, nil
}
