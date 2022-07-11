package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("failed to generate hash password: %w", err)
	}
	return string(bytes), nil
}

func CheckPassword(password, hashPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}