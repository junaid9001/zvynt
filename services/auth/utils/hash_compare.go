package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to hash password ", err)
		return []byte{}, err
	}

	return hashed, nil
}

func CompareHashAndPassword(hashedPassword, Password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(Password))

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
