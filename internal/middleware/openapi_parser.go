package middleware

import (
	"strings"

	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/getkin/kin-openapi/openapi3"
)

func LoadProtectedRoutesFromOpenAPI(openapiPath string) (map[string]map[string]bool, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(openapiPath)
	if err != nil {
		return nil, logger.LogError("error loading OpenAPI: %s", err)
	}

	if err := doc.Validate(loader.Context); err != nil {
		return nil, logger.LogError("invalid OpenAPI: %s", err)
	}

	protectedRoutes := make(map[string]map[string]bool)

	for path, pathItem := range doc.Paths.Map() {
		for method, operation := range pathItem.Operations() {
			if operation.Security != nil && len(*operation.Security) > 0 {
				for _, secReq := range *operation.Security {
					if _, hasBearerAuth := secReq["bearerAuth"]; hasBearerAuth {
						fullPath := globals.BaseURL + path
						if protectedRoutes[fullPath] == nil {
							protectedRoutes[fullPath] = make(map[string]bool)
						}
						protectedRoutes[fullPath][strings.ToUpper(method)] = true

						logger.LogInfo("✓Protected route detected: %s %s", strings.ToUpper(method), fullPath)
					}
				}
			}
		}
	}

	return protectedRoutes, nil
}
