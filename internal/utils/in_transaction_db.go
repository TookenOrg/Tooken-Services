package utils

import (
	"context"
	"database/sql"

	"github.com/TookenOrg/tooken-services/internal/globals"
)

// InTransaction runs fn in a transaction and rolls back on any failure,
// including a panic: leaving a transaction open would hold locks on the asset
// until the connection is recycled.
func InTransaction(ctx context.Context, fn func(*sql.Tx) error) (err error) {
	tx, err := globals.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}
