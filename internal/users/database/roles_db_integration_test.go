//go:build integration

package database

import (
	"errors"
	"sync"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/internal/testenv"
	"github.com/TookenOrg/tooken-services/pkg/logger"
)

func TestConcurrentAdministratorDemotions(t *testing.T) {
	db := testenv.ScratchDatabase(t, "tooken_it_user_roles")
	oldDB := globals.DB
	globals.DB = db
	t.Cleanup(func() { globals.DB = oldDB })
	logger.Init(false)
	var administrators [2]int
	for i, email := range []string{"first@example.test", "second@example.test"} {
		if err := db.QueryRow(`INSERT INTO usr.users (full_name, email, password, role)
			VALUES ('Administrator', $1, 'x', 'ADMIN') RETURNING id`, email).Scan(&administrators[i]); err != nil {
			t.Fatal(err)
		}
	}

	start := make(chan struct{})
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for i := range administrators {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_, err := UpdateUserRole(t.Context(), administrators[i], administrators[1-i], "USER")
			results <- err
		}(i)
	}
	close(start)
	wg.Wait()
	close(results)
	successes, refusals := 0, 0
	for err := range results {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, ErrRoleForbidden):
			refusals++
		default:
			t.Fatalf("unexpected demotion error: %v", err)
		}
	}
	if successes != 1 || refusals != 1 {
		t.Fatalf("successes=%d refusals=%d, want one of each", successes, refusals)
	}
	var adminCount int
	if err := db.QueryRow(`SELECT count(*) FROM usr.users WHERE role = 'ADMIN'`).Scan(&adminCount); err != nil {
		t.Fatal(err)
	}
	if adminCount != 1 {
		t.Fatalf("administrators remaining = %d, want 1", adminCount)
	}
}
