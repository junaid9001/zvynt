package shared

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   string `json:"user_id"`
	UserRole string `json:"role"`
	jwt.RegisteredClaims
}
