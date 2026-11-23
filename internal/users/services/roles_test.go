package services

import (
	"errors"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/api/server"
)

func TestUpdateUserRoleRejectsInvalidRolesBeforeDatabase(t *testing.T) {
	for _, role := range []server.UserRole{"", "user", "OWNER", " ADMIN", "MANAGER "} {
		t.Run(string(role), func(t *testing.T) {
			_, err := (&Service{}).UpdateUserRole(t.Context(), 1, 2, role)
			if !errors.Is(err, ErrInvalidRole) {
				t.Fatalf("error = %v, want invalid role", err)
			}
		})
	}
}
