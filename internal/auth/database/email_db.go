package database

import (
	"context"
	"fmt"

	"github.com/TookenOrg/tooken-services/internal/globals"
)

func CheckEmailExists(ctx context.Context, email string) (exists bool, err error) {
	// The table is schema-qualified like every other query of the package. The
	// unqualified name used to resolve only when the connection carried a
	// search_path containing usr, and would silently target a different table
	// otherwise.
	query := `
        SELECT EXISTS(SELECT 1 FROM usr.users WHERE lower(email) = lower($1))`

	err = globals.DB.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email: %w", err)
	}

	return
}
