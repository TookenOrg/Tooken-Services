package database

import (
	"context"

	"github.com/TookenOrg/tooken-services/internal/globals"
)

func CheckEmailExists(ctx context.Context, email string) (exists bool, err error) {
	query := `
        SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`

	err = globals.DB.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return
	}

	return
}
