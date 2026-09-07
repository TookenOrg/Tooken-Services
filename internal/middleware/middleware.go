package middleware

import (
	"net/http"

	"github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/gin-gonic/gin"
)

func GetUserClaims(c *gin.Context) (*utils.CustomClaims, bool) {
	value, exists := c.Get("user_claims")
	if !exists {
		return nil, false
	}

	claims, ok := value.(*utils.CustomClaims)
	return claims, ok
}

// RequireRole aborts the request unless the caller holds one of the given
// roles, and reports whether the handler may proceed.
//
// It answers 401 when no valid token was presented and 403 when the token is
// valid but the role is not allowed: telling the two apart is what lets a
// client know whether logging in again would help.
//
// The route must also carry a security block in openapi.yaml, otherwise
// AutoAuthMiddleware never parses the token and there are no claims to check.
func RequireRole(c *gin.Context, roles ...string) bool {
	claims, ok := GetUserClaims(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Authentication required.",
		})
		return false
	}

	if !claims.HasRole(roles...) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Your account is not allowed to perform this operation.",
		})
		return false
	}

	return true
}
