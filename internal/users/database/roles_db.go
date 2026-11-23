package database

import (
	"context"
	"errors"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/globals"
)

var (
	ErrRoleForbidden = errors.New("current administrator role required")
	ErrRoleConflict  = errors.New("cannot demote yourself or the last administrator")
)

func UpdateUserRole(ctx context.Context, actorID, userID int, role server.UserRole) (*server.User, error) {
	tx, err := globals.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Bound wait time and lock only administrator rows plus the actor and
	// target accounts in primary-key order, so concurrent role changes
	// serialize without blocking ordinary sign-ups or profile updates.
	if _, err := tx.ExecContext(ctx, `SET LOCAL lock_timeout = '5s'`); err != nil {
		return nil, err
	}
	lockedRows, err := tx.QueryContext(ctx, `
		SELECT id
		FROM usr.users
		WHERE role = 'ADMIN' OR id IN ($1, $2)
		ORDER BY id
		FOR UPDATE`, actorID, userID)
	if err != nil {
		return nil, err
	}
	if err := lockedRows.Close(); err != nil {
		return nil, err
	}
	var isAdmin bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS (
		SELECT 1 FROM usr.users WHERE id = $1 AND role = 'ADMIN'
	)`, actorID).Scan(&isAdmin); err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, ErrRoleForbidden
	}
	var currentRole string
	if err := tx.QueryRowContext(ctx, `SELECT role FROM usr.users WHERE id = $1`, userID).Scan(&currentRole); err != nil {
		return nil, err
	}
	if currentRole == "ADMIN" && role != "ADMIN" {
		var adminCount int
		if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM usr.users WHERE role = 'ADMIN'`).Scan(&adminCount); err != nil {
			return nil, err
		}
		if actorID == userID || adminCount <= 1 {
			return nil, ErrRoleConflict
		}
	}
	user, err := scanUser(tx.QueryRowContext(ctx, `
		UPDATE usr.users SET role = $2, updated_at = now()
		WHERE id = $1
		RETURNING id, full_name, email, created_at, role, kyc_status,
		          kyc_expires_at, kyc_verified_at, country_code, updated_at, last_connexion`,
		userID, role))
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return user, nil
}
