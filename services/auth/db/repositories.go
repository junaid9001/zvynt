package db

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

var (
	ErrDuplicateKey        = errors.New("duplicate value")
	ErrInternalServerError = errors.New("internal server error")
	ErrRecordNotFound      = errors.New("record not found")
)

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user User) (*User, error) {

	err := s.db.Create(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrDuplicateKey
		}

		log.Print(err)
		return nil, ErrInternalServerError
	}
	return &user, nil
}

func (s *Store) GetUserByEmail(email string) (*User, error) {
	var user User
	err := s.db.Where("email=?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, ErrRecordNotFound

		}

		log.Print(err)
		return nil, ErrInternalServerError

	}

	return &user, nil
}

func (s *Store) GetUserByID(id string) (*User, error) {

	var user User

	UUID := uuid.MustParse(id)

	err := s.db.First(&user, UUID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}

		log.Print(err)
		return nil, ErrInternalServerError

	}

	return &user, nil
}
