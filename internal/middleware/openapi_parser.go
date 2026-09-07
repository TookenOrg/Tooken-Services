package middleware

import (
	"strings"

	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/getkin/kin-openapi/openapi3"
)

// toGinPath rewrites an OpenAPI path template into the form gin exposes through
// c.FullPath(): /real-estates/{id} becomes /real-estates/:id.
//
// Without this the two never match. The map was keyed on the template while the
// middleware compared it to the concrete request path, so every protected route
// carrying a parameter was silently reachable without a token.
func toGinPath(path string) string {
	var b strings.Builder

	for {
		open := strings.IndexByte(path, '{')
		if open < 0 {
			b.WriteString(path)
			return b.String()
		}

		close := strings.IndexByte(path[open:], '}')
		if close < 0 {
			b.WriteString(path)
			return b.String()
		}
		close += open

		b.WriteString(path[:open])
		b.WriteByte(':')
		b.WriteString(path[open+1 : close])
		path = path[close+1:]
	}
}

// AuthMode says how a route treats the Authorization header.
type AuthMode int

const (
	// AuthRequired: no valid token, no access.
	AuthRequired AuthMode = iota
	// AuthOptional: the route serves anonymous callers, but parses the token
	// when one is presented so the handler can enrich its answer. Declared in
	// openapi.yaml by listing an empty requirement alongside the scheme:
	//
	//     security:
	//       - {}
	//       - bearerAuth: []
	AuthOptional
)

func LoadProtectedRoutesFromOpenAPI(openapiPath string) (map[string]map[string]AuthMode, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(openapiPath)
	if err != nil {
		return nil, logger.LogError("error loading OpenAPI: %s", err)
	}

	if err := doc.Validate(loader.Context); err != nil {
		return nil, logger.LogError("invalid OpenAPI: %s", err)
	}

	protectedRoutes := make(map[string]map[string]AuthMode)

	for path, pathItem := range doc.Paths.Map() {
		for method, operation := range pathItem.Operations() {
			if operation.Security == nil || len(*operation.Security) == 0 {
				continue
			}

			hasBearerAuth := false
			allowsAnonymous := false

			for _, secReq := range *operation.Security {
				if _, ok := secReq["bearerAuth"]; ok {
					hasBearerAuth = true
				}
				// An empty requirement is how OpenAPI spells "authentication
				// is one option among others", i.e. anonymous is accepted.
				if len(secReq) == 0 {
					allowsAnonymous = true
				}
			}

			if !hasBearerAuth {
				continue
			}

			mode := AuthRequired
			label := "Protected"
			if allowsAnonymous {
				mode = AuthOptional
				label = "Optional-auth"
			}

			fullPath := globals.BaseURL + toGinPath(path)
			if protectedRoutes[fullPath] == nil {
				protectedRoutes[fullPath] = make(map[string]AuthMode)
			}
			protectedRoutes[fullPath][strings.ToUpper(method)] = mode

			logger.LogInfo("✓%s route detected: %s %s", label, strings.ToUpper(method), fullPath)
		}
	}

	return protectedRoutes, nil
}
