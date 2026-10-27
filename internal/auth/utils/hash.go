package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// MaxPasswordBytes is the most bcrypt hashes: it refuses longer input. The
// limit is in bytes, not characters, so a non-ASCII password reaches it sooner.
const MaxPasswordBytes = 72

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
