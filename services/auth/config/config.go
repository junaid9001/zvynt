package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APP_PORT string

	DB_URL     string
	JWT_SECRET string
}

func LoadConfig() *Config {
	godotenv.Load()
	appPort := os.Getenv("APP_PORT")
	jwtSecret := os.Getenv("JWT_SECRET")
	dbUrl := os.Getenv("AUTH_DB_URL")

	return &Config{
		APP_PORT:   appPort,
		JWT_SECRET: jwtSecret,
		DB_URL:     dbUrl,
	}

}
