package db

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name     string    `gorm:"type:varchar(100);not null"`
	Email    string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	Password string    `gorm:"not null"`
	Role     string    `gorm:"not null;default:user"`
}
