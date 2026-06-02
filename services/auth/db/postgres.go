package db

import (
	"log"

	"github.com/junaid9001/zvynt/auth/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) *gorm.DB {

	db, err := gorm.Open(postgres.Open(cfg.DB_URL), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db connection success")

	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
