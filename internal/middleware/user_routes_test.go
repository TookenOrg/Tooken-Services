package middleware

import (
	"strings"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/getkin/kin-openapi/openapi3"
)

// An operation that names a user in its path acts on that user's own data.
// Without a security block the middleware never parses a token, the handler
// cannot compare the path to the caller, and anybody acts on anybody: the
// favorites were shipped that way. The check runs over the whole spec, so a
// personal route added later is held to the same rule without being listed.
func TestRoutesNamingAUserRequireAToken(t *testing.T) {
	const spec = "../../api/openapi.yaml"

	routes, err := LoadProtectedRoutesFromOpenAPI(spec)
	if err != nil {
		t.Fatalf("LoadProtectedRoutesFromOpenAPI: %v", err)
	}
	doc, err := openapi3.NewLoader().LoadFromFile(spec)
	if err != nil {
		t.Fatalf("load %s: %v", spec, err)
	}

	checked := 0
	for path, item := range doc.Paths.Map() {
		if !strings.Contains(path, "{userId}") {
			continue
		}
		for method := range item.Operations() {
			checked++
			full := globals.BaseURL + toGinPath(path)
			mode, ok := routes[full][strings.ToUpper(method)]
			if !ok {
				t.Errorf("%s %s is public: it acts on the user of its path and must declare bearerAuth", method, full)
				continue
			}
			if mode != AuthRequired {
				t.Errorf("%s %s is %v, want AuthRequired: an anonymous caller has no user to compare the path to", method, full, mode)
			}
		}
	}

	// The favorites are the routes this rule was written for: if they stopped
	// naming a user, the loop above would pass without checking anything.
	if checked < 2 {
		t.Fatalf("%d operation(s) name a user in their path, want at least the two favorite ones", checked)
	}
}
