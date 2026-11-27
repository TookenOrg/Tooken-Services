package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// MaxPasswordBytes is the most bcrypt hashes: it refuses longer input. The
// limit is in bytes, not characters, so a non-ASCII password reaches it sooner.
const MaxPasswordBytes = 72

// PasswordCost is the bcrypt cost used for every new hash.
const PasswordCost = 12

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// NeedsRehash reports whether hash was produced with a cost other than PasswordCost.
// A malformed hash returns false.
func NeedsRehash(hash string) bool {
	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		return false
	}
	return cost != PasswordCost
}
