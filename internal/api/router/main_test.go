package router

import (
	"os"
	"testing"

	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
)

// The server installs its secret from JWT_SECRET at startup; tests install one
// of their own so that tokens can be issued and parsed.
func TestMain(m *testing.M) {
	if err := authUtils.SetSecret("test-only-secret-at-least-32-bytes"); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
