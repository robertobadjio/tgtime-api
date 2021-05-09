package service

import (
	"golang.org/x/crypto/bcrypt"
)

func CheckAuth(passwordHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))

	return err == nil
}

