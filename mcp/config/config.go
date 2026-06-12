package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	BankAPIBaseURL string
	JWTToken       string
	Log            LogConfig
	PubSub         PubSubConfig
}

type LogConfig struct {
	Level string
}

type PubSubConfig struct {
	ProjectID string
	TopicID   string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		BankAPIBaseURL: getEnv("BANK_API_BASE_URL", "http://localhost:8080"),
		JWTToken:       getEnv("JWT_TOKEN", ""),
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
		PubSub: PubSubConfig{
			ProjectID: getEnv("GCP_PROJECT_ID", "tempBankLogs"),
			TopicID:   getEnv("GCP_PUBSUB_TOPIC_ID", "MCPServerLogs"),
		},
	}
}

func (c *Config) String() string {
	return "Config{BankAPIBaseURL: " + c.BankAPIBaseURL + ", JWTToken: " + maskSecret(c.JWTToken) + ", Log: " + c.Log.Level + ", PubSub: " + c.PubSub.ProjectID + "/" + c.PubSub.TopicID + "}"
}

func maskSecret(s string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) <= 4 {
		return strings.Repeat("*", len(s))
	}
	return s[:2] + strings.Repeat("*", len(s)-4) + s[len(s)-2:]
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
