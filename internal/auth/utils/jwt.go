package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte("secret_key_to_change")

// Role values, mirroring the CHECK constraint on usr.users.role added by
// migration 000009.
const (
	RoleUser    = "USER"
	RoleManager = "MANAGER"
	RoleAdmin   = "ADMIN"
)

type CustomClaims struct {
	UserID int    `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// HasRole reports whether the caller holds one of the given roles.
//
// The role is carried by the token rather than read from the database on every
// request, so a role revoked mid-session only takes effect once the token
// expires. With a one hour lifetime that window is accepted.
func (c *CustomClaims) HasRole(roles ...string) bool {
	for _, role := range roles {
		if c.Role == role {
			return true
		}
	}

	return false
}

func GenerateJWT(userID int, email, role string) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

func GenerateRefreshToken() (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 8)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

func ParseJWT(tokenStr string) (claims *CustomClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
