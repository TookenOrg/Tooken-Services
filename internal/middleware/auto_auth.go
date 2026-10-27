package middleware

import (
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"

	"github.com/gin-gonic/gin"
)

// protectedRoutes holds the routes read from the spec by the last successful
// InitAuth. It is nil until then, and the middleware refuses every request
// while it is.
var protectedRoutes atomic.Pointer[map[string]map[string]AuthMode]

// InitAuth reads from the OpenAPI document which operations require a token.
//
// A failed load keeps the routes already loaded. Assigning its result anyway
// would leave the map nil, and a nil map protects nothing: every route would
// be served without a token, silently.
func InitAuth(openapiPath string) error {
	routes, err := LoadProtectedRoutesFromOpenAPI(openapiPath)
	if err != nil {
		return err
	}
	protectedRoutes.Store(&routes)

	logger.LogInfo("✓ %d protected routes loaded from OpenApi files", len(routes))
	return nil
}

func AutoAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// FullPath is the route pattern (/real-estates/:id), not the concrete
		// URL: it is what can be compared to the OpenAPI paths. It is empty
		// when no route matched, in which case there is nothing to protect.
		path := c.FullPath()
		method := c.Request.Method

		routes := protectedRoutes.Load()
		if routes == nil {
			// Without the list nothing tells a protected route from a public
			// one, and serving them all would hand out every protected one.
			logger.LogError("❌ Request refused: InitAuth was never called, no route is known to be protected")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Authentication is not configured.",
			})
			return
		}

		methods, exists := (*routes)[path]
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
				"message": "Authorization header required",
			})
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid authorization format. Use: Bearer <token>",
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
				"message": "Invalid or expired token",
			})
			return
		}

		// Save claims in gCtx
		c.Set("user_claims", claims)

		c.Next()
	}
}
