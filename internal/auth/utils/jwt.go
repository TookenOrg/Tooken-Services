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

const (
	AccessTokenLifetime = time.Hour
	// RefreshTokenLifetime bounds the whole session: refreshing issues a new
	// access token but never extends the refresh token. Nothing is stored
	// server side, so a stolen refresh token cannot be revoked; a fixed end
	// at least keeps the thief from renewing it forever.
	RefreshTokenLifetime = 8 * time.Hour
)

// Both kinds of token are signed with the same key, so only this claim tells
// them apart. Each parser demands its own value, which also refuses tokens
// issued before the claim existed.
const (
	tokenUseAccess  = "access"
	tokenUseRefresh = "refresh"
)

// ErrWrongTokenUse: the token is genuine but meant for the other purpose.
var ErrWrongTokenUse = errors.New("token not valid for this use")

type CustomClaims struct {
	UserID   int    `json:"userId"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	TokenUse string `json:"tokenUse"`
	jwt.RegisteredClaims
}

// RefreshClaims carries the account only: the role and email are read again
// from the database on every refresh.
type RefreshClaims struct {
	UserID   int    `json:"userId"`
	TokenUse string `json:"tokenUse"`
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
	token, _, err := NewAccessToken(userID, email, role)
	return token, err
}

// NewAccessToken also returns the expiry written in the token, so a client
// can schedule its refresh without decoding it.
func NewAccessToken(userID int, email, role string) (token string, expiresAt time.Time, err error) {
	now := time.Now()
	claims := CustomClaims{
		UserID:   userID,
		Email:    email,
		Role:     role,
		TokenUse: tokenUseAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenLifetime)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(JwtKey)
	if err != nil {
		return "", time.Time{}, err
	}

	// NumericDate drops the sub-second part: report what the token says.
	return token, claims.ExpiresAt.Time, nil
}

func GenerateRefreshToken(userID int) (string, error) {
	now := time.Now()
	claims := RefreshClaims{
		UserID:   userID,
		TokenUse: tokenUseRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(RefreshTokenLifetime)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(JwtKey)
}

// ParseJWT validates an access token. A refresh token is refused.
func ParseJWT(tokenStr string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	if err := parse(tokenStr, claims); err != nil {
		return nil, err
	}

	if claims.TokenUse != tokenUseAccess {
		return nil, ErrWrongTokenUse
	}

	return claims, nil
}

// ParseRefreshToken validates a refresh token. An access token is refused.
func ParseRefreshToken(tokenStr string) (*RefreshClaims, error) {
	claims := &RefreshClaims{}
	if err := parse(tokenStr, claims); err != nil {
		return nil, err
	}

	if claims.TokenUse != tokenUseRefresh || claims.UserID <= 0 {
		return nil, ErrWrongTokenUse
	}

	return claims, nil
}

func parse(tokenStr string, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(tokenStr, claims,
		func(*jwt.Token) (interface{}, error) { return JwtKey, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		// Every token this package issues expires; one that does not was not
		// issued here.
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}
