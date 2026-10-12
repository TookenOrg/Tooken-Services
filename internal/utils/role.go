package utils

import (
	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/TookenOrg/tooken-services/internal/middleware"
	"github.com/gin-gonic/gin"
)

// IsStaff reports whether the caller manages the catalogue.
//
// The two real estate read endpoints are declared with an optional security
// block, so claims are present only when a token was supplied: no token means
// an anonymous visitor, and the public view.
func IsStaff(gCtx *gin.Context) bool {
	claims, ok := middleware.GetUserClaims(gCtx)
	if !ok {
		return false
	}

	return claims.HasRole(authUtils.RoleManager, authUtils.RoleAdmin)
}
