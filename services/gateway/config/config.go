package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APP_PORT   string
	REDIS_ADDR string
	JWT_SECRET string
	AUTH_ADDR  string
}

func LoadConfig() *Config {
	godotenv.Load()
	appPort := os.Getenv("APP_PORT")
	redisAddr := os.Getenv("REDIS_ADDR")
	jwtSecret := os.Getenv("JWT_SECRET")
	authAddr := os.Getenv("AUTH_ADDR")

	return &Config{
		APP_PORT:   appPort,
		REDIS_ADDR: redisAddr,
		JWT_SECRET: jwtSecret,
		AUTH_ADDR:  authAddr,
	}

}
