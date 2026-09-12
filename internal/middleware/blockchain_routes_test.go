package middleware

import (
	"testing"

	"github.com/TookenOrg/tooken-services/internal/globals"
)

// Route protection lives in the OpenAPI document, not in the code: an operation
// without a security block is served to anyone, and nothing in the build, the
// vet pass or the test suite notices. The blockchain operations were shipped
// that way - deploying the factory, creating a token, minting and burning were
// all reachable without a token, on endpoints that spend the platform private
// key. Adding a route is the moment the omission happens again, so the list
// below is the reminder: a new write path under /trex or /contract has to be
// added here, and that only passes once the spec declares bearerAuth.
func TestBlockchainWriteRoutesRequireAToken(t *testing.T) {
	routes, err := LoadProtectedRoutesFromOpenAPI("../../api/openapi.yaml")
	if err != nil {
		t.Fatalf("LoadProtectedRoutesFromOpenAPI: %v", err)
	}

	writeRoutes := []struct{ method, path string }{
		{"POST", "/trex/init/implementations"},
		{"POST", "/trex/init/contract/factory/identity"},
		{"POST", "/trex/init/authority/configure"},
		{"POST", "/trex/init/contract/factory/trexFactory"},
		{"POST", "/trex/init/suite/deploy"},
		{"POST", "/contract/identity"},
		{"POST", "/contract/identity/claims/add"},
		{"POST", "/contract/token"},
		{"POST", "/contract/token/mint"},
		{"POST", "/contract/token/burn"},
	}

	for _, r := range writeRoutes {
		full := globals.BaseURL + r.path
		mode, ok := routes[full][r.method]
		if !ok {
			t.Errorf("%s %s is public: it deploys or moves tokens and must declare bearerAuth", r.method, full)
			continue
		}
		if mode != AuthRequired {
			t.Errorf("%s %s is %v, want AuthRequired: optional auth still serves anonymous callers", r.method, full, mode)
		}
	}
}

// The counterpart: read-only blockchain endpoints stay open. Locking them would
// be the easy over-correction, and it would break the public catalogue, which
// has to show a token to a visitor who has not signed up yet.
func TestBlockchainReadRoutesStayPublic(t *testing.T) {
	routes, err := LoadProtectedRoutesFromOpenAPI("../../api/openapi.yaml")
	if err != nil {
		t.Fatalf("LoadProtectedRoutesFromOpenAPI: %v", err)
	}

	readRoutes := []struct{ method, path string }{
		{"GET", "/trex/suite"},
		{"GET", "/contract/token/:tokenAddress/infos"},
	}

	for _, r := range readRoutes {
		full := globals.BaseURL + r.path
		if _, ok := routes[full][r.method]; ok {
			t.Errorf("%s %s now requires a token: it only reads on-chain state", r.method, full)
		}
	}
}
