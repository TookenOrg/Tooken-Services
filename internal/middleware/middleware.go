package middleware

import (
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
