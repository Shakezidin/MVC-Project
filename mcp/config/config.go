package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BankAPIBaseURL string
	JWTToken       string
	Log            LogConfig
}

type LogConfig struct {
	Level string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("No .env file found")
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

		Log: LogConfig{
			Level: getEnv(
				"LOG_LEVEL",
				"info",
			)},
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
