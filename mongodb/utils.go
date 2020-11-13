package mongodb

import (
	"golang.org/x/crypto/bcrypt"
)

func hashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func comparePasswords(hash, pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
}
