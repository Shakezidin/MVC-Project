package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BankAPIBaseURL string
	JWTToken       string
	LogLevel       string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		BankAPIBaseURL: getEnv(
			"BANK_API_BASE_URL",
			"http://localhost:8080",
		),

		JWTToken: getEnv(
			"JWT_TOKEN",
			"",
		),

		LogLevel: getEnv(
			"LOG_LEVEL",
			"info",
		),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
