package passhash

import (
	"golang.org/x/crypto/bcrypt"
)

// HashString returns a hashed string and an error
func HashString(password string) (string, error) {
	key, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(key), nil
}

func CompareHashAndPassword(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
