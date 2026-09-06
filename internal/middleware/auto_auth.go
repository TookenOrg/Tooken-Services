package middleware

import (
	"net/http"
	"strings"

	"github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Global var for protected routes (with bearer)
var protectedRoutes map[string]map[string]AuthMode

func InitAuth(openapiPath string) error {
	var err error
	protectedRoutes, err = LoadProtectedRoutesFromOpenAPI(openapiPath)
	if err != nil {
		return err
	}

	logger.LogInfo("✓ %d protected routes loaded from OpenApi files", len(protectedRoutes))
	return nil
}

func AutoAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// FullPath is the route pattern (/real-estates/:id), not the concrete
		// URL: it is what can be compared to the OpenAPI paths. It is empty
		// when no route matched, in which case there is nothing to protect.
		path := c.FullPath()
		method := c.Request.Method

		methods, exists := protectedRoutes[path]
		if !exists {
			c.Next()
			return
		}

		mode, declared := methods[method]
		if !declared {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			// An optional-auth route serves anonymous callers; the handler
			// sees no claims and answers with the public view.
			if mode == AuthOptional {
				c.Next()
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format. Use: ******",
			})
			return
		}

		// A malformed or expired token is rejected even where auth is optional:
		// silently downgrading to anonymous would answer 200 with a truncated
		// payload, and the caller would never learn its session expired.
		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}

		// Save claims in gCtx
		c.Set("user_claims", claims)

		c.Next()
	}
}
