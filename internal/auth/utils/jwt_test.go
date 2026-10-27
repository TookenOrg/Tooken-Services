package utils

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func sign(t *testing.T, claims jwt.Claims) string {
	t.Helper()
	return signHS256With(t, claims, JwtKey)
}

// Refresh tokens used to carry nothing but an expiry. Signed with the same key,
// one presented as a Bearer token passed the signature check and authenticated
// a caller with no user id and no role, for eight hours.
func TestParseJWTRefusesARefreshToken(t *testing.T) {
	legacy := sign(t, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	if claims, err := ParseJWT(legacy); err == nil {
		t.Fatalf("a refresh token without a use was accepted as an access token: %+v", claims)
	}
}

func TestTheTwoKindsOfTokenAreNotInterchangeable(t *testing.T) {
	access, expiresAt, err := NewAccessToken(42, "investor@tooken.test", RoleUser)
	if err != nil {
		t.Fatal(err)
	}
	refresh, err := GenerateRefreshToken(42)
	if err != nil {
		t.Fatal(err)
	}

	claims, err := ParseJWT(access)
	if err != nil {
		t.Fatalf("access token refused: %v", err)
	}
	if claims.UserID != 42 || claims.Role != RoleUser || !claims.ExpiresAt.Time.Equal(expiresAt) {
		t.Errorf("access claims: got %+v, want user 42, role USER, expiry %s", claims, expiresAt)
	}

	refreshClaims, err := ParseRefreshToken(refresh)
	if err != nil {
		t.Fatalf("refresh token refused: %v", err)
	}
	if refreshClaims.UserID != 42 {
		t.Errorf("refresh user: got %d, want 42", refreshClaims.UserID)
	}

	if _, err := ParseJWT(refresh); !errors.Is(err, ErrWrongTokenUse) {
		t.Errorf("refresh token as access token: got %v, want ErrWrongTokenUse", err)
	}
	if _, err := ParseRefreshToken(access); !errors.Is(err, ErrWrongTokenUse) {
		t.Errorf("access token as refresh token: got %v, want ErrWrongTokenUse", err)
	}
}

// HMAC keys verify any HS* algorithm unless the parser pins one: the server
// signs HS256 only, so anything else is not one of its tokens.
func TestParsersAcceptHS256Only(t *testing.T) {
	inAnHour := jwt.NewNumericDate(time.Now().Add(time.Hour))
	access := CustomClaims{UserID: 42, Role: RoleAdmin, TokenUse: tokenUseAccess,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: inAnHour}}
	refresh := RefreshClaims{UserID: 42, TokenUse: tokenUseRefresh,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: inAnHour}}

	signWith := func(method jwt.SigningMethod, claims jwt.Claims, key any) string {
		token, err := jwt.NewWithClaims(method, claims).SignedString(key)
		if err != nil {
			t.Fatal(err)
		}
		return token
	}

	for _, method := range []jwt.SigningMethod{jwt.SigningMethodHS384, jwt.SigningMethodHS512} {
		t.Run(method.Alg(), func(t *testing.T) {
			if claims, err := ParseJWT(signWith(method, access, JwtKey)); err == nil {
				t.Errorf("access token accepted: %+v", claims)
			}
			if claims, err := ParseRefreshToken(signWith(method, refresh, JwtKey)); err == nil {
				t.Errorf("refresh token accepted: %+v", claims)
			}
		})
	}
	t.Run("none", func(t *testing.T) {
		unsigned := signWith(jwt.SigningMethodNone, access, jwt.UnsafeAllowNoneSignatureType)
		if claims, err := ParseJWT(unsigned); err == nil {
			t.Errorf("unsigned token accepted: %+v", claims)
		}
	})
}

func TestParseRefreshTokenRefusals(t *testing.T) {
	inAnHour := jwt.NewNumericDate(time.Now().Add(time.Hour))

	tests := map[string]string{
		"expired": sign(t, RefreshClaims{UserID: 42, TokenUse: tokenUseRefresh,
			RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute))}}),
		"no expiry": sign(t, RefreshClaims{UserID: 42, TokenUse: tokenUseRefresh}),
		"no account": sign(t, RefreshClaims{TokenUse: tokenUseRefresh,
			RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: inAnHour}}),
		"issued before the use claim": sign(t, jwt.RegisteredClaims{ExpiresAt: inAnHour}),
		"another key": signHS256With(t, RefreshClaims{UserID: 42, TokenUse: tokenUseRefresh,
			RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: inAnHour}}, []byte("not-the-key")),
		"garbage": "not.a.token",
	}

	for name, token := range tests {
		t.Run(name, func(t *testing.T) {
			if claims, err := ParseRefreshToken(token); err == nil {
				t.Errorf("accepted: %+v", claims)
			}
		})
	}
}

func signHS256With(t *testing.T, claims jwt.Claims, key []byte) string {
	t.Helper()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return token
}
