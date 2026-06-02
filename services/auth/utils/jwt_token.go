package utils

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/junaid9001/zvynt/pkg/shared"
)

func GenerateJwt(jwtSecret, userID, userRole string) (string, error) {

	claims := shared.Claims{
		UserID:   userID,
		UserRole: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "zvynt",
		},
	}

	tokenStr := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	token, err := tokenStr.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Println("failed to sign jwt token ", err)
		return "", err
	}

	return token, nil

}
