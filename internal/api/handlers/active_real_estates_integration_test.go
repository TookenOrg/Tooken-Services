//go:build integration

// HTTP coverage of GET /assets/real-estates/active, the public catalogue that
// staff also use to see what is not published yet.
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	_ "github.com/lib/pq"
)

const activeRealEstatesPath = "/assets/real-estates/active"

func TestListActiveRealEstates(t *testing.T) {
	h := newOrderHarness(t)
	fixtures := newAssetFixtures(t, h)
	issuer := fixtures.activeIssuer(t, "Catalogue")
	published := fixtures.asset(t, issuer, true)
	draft := fixtures.asset(t, issuer, false)
	deleted := fixtures.deletedAsset(t, issuer)

	user := newTrackedUser(t, h, authUtils.RoleUser)
	manager := newTrackedUser(t, h, authUtils.RoleManager)
	admin := newTrackedUser(t, h, authUtils.RoleAdmin)

	// The shared database holds assets from other suites, so the listing is
	// checked on the three identifiers seeded here, never on its length.
	for _, c := range []struct {
		caller     string
		token      string
		seesDrafts bool
	}{
		{"an anonymous visitor", "", false},
		{"an investor", user.token, false},
		{"a manager", manager.token, true},
		{"an admin", admin.token, true},
	} {
		t.Run(c.caller+" gets the catalogue of his role", func(t *testing.T) {
			code, body := callAPI(h, http.MethodGet, activeRealEstatesPath, c.token, "")
			expectCode(t, code, body, http.StatusOK)
			var listed []server.RealEstateSummary
			if err := json.Unmarshal([]byte(body), &listed); err != nil {
				t.Fatalf("decoding the listing: %v", err)
			}
			seen := map[int]bool{}
			for _, e := range listed {
				seen[e.Id] = true
			}

			if !seen[published] {
				t.Errorf("published asset %d missing", published)
			}
			if seen[draft] != c.seesDrafts {
				t.Errorf("draft asset %d listed = %v, want %v", draft, seen[draft], c.seesDrafts)
			}
			if seen[deleted] {
				t.Errorf("deleted asset %d listed", deleted)
			}
		})
	}

	// An empty catalogue would tell visitors there is nothing to invest in;
	// a failure must say so instead.
	t.Run("a database failure answers 500 without leaking the cause", func(t *testing.T) {
		withClosedDatabase(t, func() {
			code, body := callAPI(h, http.MethodGet, activeRealEstatesPath, "", "")
			expectCode(t, code, body, http.StatusInternalServerError)
			expectAPIMessage(t, body, "Unable to retrieve real estates.")
			if strings.Contains(body, "database is closed") {
				t.Fatalf("the cause leaked to the caller: %s", body)
			}
		})
	})
}
