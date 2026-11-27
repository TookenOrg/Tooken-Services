package utils

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// Precomputed cost-14 bcrypt hash for "Sup3r-secret!" so unit tests stay fast
// while exercising NeedsRehash on a real legacy hash.
const legacyCost14Hash = "$2a$14$vNtfHEd3PJ2x.WP0d0oizuOeKY5GoHffm9OUIB1SFhk23V/fPq.Ha"

func TestHashPasswordUsesPasswordCost(t *testing.T) {
	const password = "Sup3r-secret!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}

	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		t.Fatalf("bcrypt.Cost: %v", err)
	}
	if cost != PasswordCost {
		t.Errorf("cost = %d, want %d", cost, PasswordCost)
	}
	if !CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash rejected the password that produced the hash")
	}
}

func TestNeedsRehash(t *testing.T) {
	cost12Hash, err := HashPassword("Sup3r-secret!")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	cost10Bytes, err := bcrypt.GenerateFromPassword([]byte("Sup3r-secret!"), 10)
	if err != nil {
		t.Fatalf("GenerateFromPassword cost 10: %v", err)
	}

	tests := []struct {
		name string
		hash string
		want bool
	}{
		{name: "current cost", hash: cost12Hash, want: false},
		{name: "higher legacy cost", hash: legacyCost14Hash, want: true},
		{name: "lower cost", hash: string(cost10Bytes), want: true},
		{name: "malformed hash", hash: "not-a-bcrypt-hash", want: false},
		{name: "empty hash", hash: "", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := NeedsRehash(tc.hash); got != tc.want {
				t.Errorf("NeedsRehash(%q) = %v, want %v", tc.hash, got, tc.want)
			}
		})
	}
}
