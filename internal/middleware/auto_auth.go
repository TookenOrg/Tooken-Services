package middleware

import (
	"net/http"
	"strings"

	"github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Global var for protected routes (with bearer)
var protectedRoutes map[string]map[string]bool

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
		path := c.Request.URL.Path
		method := c.Request.Method

		if methods, exists := protectedRoutes[path]; exists {
			if methods[method] {
				// This route must be protected
				authHeader := c.GetHeader("Authorization")

				if authHeader == "" {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"error": "Authorization header required",
					})
					return
				}

				if !strings.HasPrefix(authHeader, "Bearer ") {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"error": "Invalid authorization format. Use: Bearer <token>",
					})
					return
				}

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
			}
		}

		c.Next()
	}
}
