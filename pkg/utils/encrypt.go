package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
