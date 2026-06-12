package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
// Centralizing config enables consistent validation at startup rather than
// discovering missing values at runtime.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Log      LogConfig
	Cache    CacheConfig
	PubSub   PubSubConfig
}

type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	RequestTimeout  time.Duration
	RateLimitRPS    int
	RateLimitBurst  int
	Environment     string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
	Issuer string
}

type LogConfig struct {
	Level string
}

type CacheConfig struct {
	AccountListTTL     time.Duration
	BeneficiaryListTTL time.Duration
	TransferModesTTL   time.Duration
}

type PubSubConfig struct {
	ProjectID string
	TopicID   string
}

// Load reads environment variables and returns a validated Config.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8080"),
			ReadTimeout:     getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout:    getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:     getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: getDurationEnv("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
			RequestTimeout:  getDurationEnv("REQUEST_TIMEOUT", 30*time.Second),
			RateLimitRPS:    getIntEnv("RATE_LIMIT_RPS", 100),
			RateLimitBurst:  getIntEnv("RATE_LIMIT_BURST", 200),
			Environment:     getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "bankuser"),
			Password:        getEnv("DB_PASSWORD", "bankpass"),
			Name:            getEnv("DB_NAME", "bankdb"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "redis://localhost:6379/0"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
			Expiry: getDurationEnv("JWT_EXPIRY", 24*time.Hour),
			Issuer: getEnv("JWT_ISSUER", "bank-server"),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
		Cache: CacheConfig{
			AccountListTTL:     getDurationEnv("CACHE_ACCOUNT_LIST_TTL", 5*time.Minute),
			BeneficiaryListTTL: getDurationEnv("CACHE_BENEFICIARY_LIST_TTL", 5*time.Minute),
			TransferModesTTL:   getDurationEnv("CACHE_TRANSFER_MODES_TTL", 1*time.Hour),
		},
		PubSub: PubSubConfig{
			ProjectID: getEnv("GCP_PROJECT_ID", "tempBankLogs"),
			TopicID:   getEnv("GCP_PUBSUB_TOPIC_ID", "bankServerLogs"),
		},
	}

	return cfg, cfg.Validate()
}

// DSN returns the PostgreSQL connection string.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

// Validate ensures required configuration is present and sane.
func (c *Config) Validate() error {
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if len(c.JWT.Secret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	return nil
}

// String returns a string representation of the config with sensitive fields masked.
func (c *Config) String() string {
	return fmt.Sprintf(
		"Config{Server: %+v, Database: %+v, Redis: %+v, JWT: {Secret: %s, Expiry: %v, Issuer: %s}, Log: %+v, Cache: %+v, PubSub: %+v}",
		c.Server,
		c.Database.masked(),
		c.Redis.masked(),
		maskSecret(c.JWT.Secret),
		c.JWT.Expiry,
		c.JWT.Issuer,
		c.Log,
		c.Cache,
		c.PubSub,
	)
}

func (d DatabaseConfig) masked() string {
	return fmt.Sprintf(
		"DatabaseConfig{Host: %s, Port: %s, User: %s, Password: %s, Name: %s, SSLMode: %s, MaxOpenConns: %d, MaxIdleConns: %d, ConnMaxLifetime: %v}",
		d.Host,
		d.Port,
		d.User,
		maskSecret(d.Password),
		d.Name,
		d.SSLMode,
		d.MaxOpenConns,
		d.MaxIdleConns,
		d.ConnMaxLifetime,
	)
}

func (r RedisConfig) masked() string {
	return fmt.Sprintf(
		"RedisConfig{URL: %s, Password: %s, DB: %d}",
		r.URL,
		maskSecret(r.Password),
		r.DB,
	)
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

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
