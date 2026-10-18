package database

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/lib/pq"
)

// IsContractRoleNameTaken is what lets a token creation that lost the race for a
// shared module adopt the winner's instead of failing (B7). It must recognise the
// violation through the wrapping the deployment adds, and nothing else.
func TestIsContractRoleNameTaken(t *testing.T) {
	taken := &pq.Error{Code: "23505", Constraint: contractRoleNameConstraint}

	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"the violation itself", taken, true},
		{"wrapped twice, as the deployment returns it", fmt.Errorf("deploy X: %w", fmt.Errorf("insert role: %w", taken)), true},
		{"another unique constraint", &pq.Error{Code: "23505", Constraint: "contract_role_tx_hash_key"}, false},
		{"another error on the same constraint", &pq.Error{Code: "23514", Constraint: contractRoleNameConstraint}, false},
		{"a plain error", fmt.Errorf("duplicate key value violates unique constraint"), false},
		{"no row", sql.ErrNoRows, false},
		{"nil", nil, false},
	}
	for _, c := range cases {
		if got := IsContractRoleNameTaken(c.err); got != c.want {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}
