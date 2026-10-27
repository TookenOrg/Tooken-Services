package utils

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func sign(t *testing.T, claims jwt.Claims) string {
	t.Helper()
	return signHS256With(t, claims, jwtKey)
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
			if claims, err := ParseJWT(signWith(method, access, jwtKey)); err == nil {
				t.Errorf("access token accepted: %+v", claims)
			}
			if claims, err := ParseRefreshToken(signWith(method, refresh, jwtKey)); err == nil {
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

func adminClaims() CustomClaims {
	return CustomClaims{UserID: 1, Email: "admin@tooken.test", Role: RoleAdmin, TokenUse: tokenUseAccess,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))}}
}

func signHS256With(t *testing.T, claims jwt.Claims, key []byte) string {
	t.Helper()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

// The key used to be a constant in this file, so it is public: a token signed
// with it must not authenticate anyone.
func TestATokenSignedWithTheFormerBuiltInKeyIsRefused(t *testing.T) {
	forged := signHS256With(t, adminClaims(), []byte("secret_key_to_change"))
	if claims, err := ParseJWT(forged); err == nil {
		t.Fatalf("forged ADMIN token accepted: %+v", claims)
	}
}

func TestSetSecretRefusesAWeakSecret(t *testing.T) {
	previous := jwtKey
	t.Cleanup(func() { jwtKey = previous })

	for _, tc := range []struct {
		secret string
		want   error
	}{
		{"", ErrSecretMissing},
		{"secret_key_to_change", ErrSecretTooShort},
		{strings.Repeat("k", MinSecretBytes-1), ErrSecretTooShort},
		{strings.Repeat("k", MinSecretBytes-1) + "\n", ErrSecretPadded},
		{" " + strings.Repeat("k", MinSecretBytes), ErrSecretPadded},
		{strings.Repeat(" ", MinSecretBytes), ErrSecretPadded},
	} {
		if err := SetSecret(tc.secret); !errors.Is(err, tc.want) {
			t.Errorf("SetSecret(%q) = %v, want %v", tc.secret, err, tc.want)
		}
		if !bytes.Equal(jwtKey, previous) {
			t.Fatalf("SetSecret(%q) was refused but replaced the key", tc.secret)
		}
	}

	issuedBefore, err := GenerateJWT(42, "investor@tooken.test", RoleUser)
	if err != nil {
		t.Fatal(err)
	}
	if err := SetSecret(strings.Repeat("k", MinSecretBytes)); err != nil {
		t.Fatalf("SetSecret(%d bytes) = %v, want nil", MinSecretBytes, err)
	}
	token, err := GenerateJWT(42, "investor@tooken.test", RoleUser)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseJWT(token); err != nil {
		t.Errorf("token signed with the new secret refused: %v", err)
	}
	if _, err := ParseJWT(issuedBefore); err == nil {
		t.Error("token signed with the replaced secret still accepted")
	}
}

// The HMAC library signs and verifies with an empty key. Until a secret is
// installed nothing is issued and nothing is accepted, not even a token signed
// with that empty key.
func TestWithoutASecretNoTokenIsIssuedOrAccepted(t *testing.T) {
	access, err := GenerateJWT(42, "investor@tooken.test", RoleUser)
	if err != nil {
		t.Fatal(err)
	}
	refresh, err := GenerateRefreshToken(42)
	if err != nil {
		t.Fatal(err)
	}

	previous := jwtKey
	jwtKey = nil
	t.Cleanup(func() { jwtKey = previous })

	if _, _, err := NewAccessToken(42, "investor@tooken.test", RoleUser); !errors.Is(err, ErrSecretMissing) {
		t.Errorf("NewAccessToken without a secret = %v, want ErrSecretMissing", err)
	}
	if _, err := GenerateRefreshToken(42); !errors.Is(err, ErrSecretMissing) {
		t.Errorf("GenerateRefreshToken without a secret = %v, want ErrSecretMissing", err)
	}
	for name, token := range map[string]string{
		"issued with the previous secret": access,
		"signed with an empty key":        signHS256With(t, adminClaims(), []byte{}),
	} {
		if claims, err := ParseJWT(token); err == nil {
			t.Errorf("access token %s accepted: %+v", name, claims)
		}
	}
	if claims, err := ParseRefreshToken(refresh); err == nil {
		t.Errorf("refresh token accepted: %+v", claims)
	}
}
