package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	JWTSecret   []byte
	TokenExpiry time.Duration
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	tokenExpiry, err := time.ParseDuration(envOrDefault("TOKEN_TTL", "8h"))
	if err != nil || tokenExpiry <= 0 {
		return Config{}, fmt.Errorf("TOKEN_TTL must be a positive Go duration")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" && envOrDefault("APP_ENV", "development") == "development" {
		secret = "notix-development-only-secret-not-for-production"
	}
	if len(secret) < 32 {
		return Config{}, fmt.Errorf("JWT_SECRET must contain at least 32 characters; set it before running outside development")
	}

	port := envOrDefault("PORT", "8080")
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return Config{}, fmt.Errorf("PORT must be an integer from 1 to 65535")
	}

	return Config{
		Port:        port,
		DBHost:      envOrDefault("DB_HOST", "localhost"),
		DBPort:      envOrDefault("DB_PORT", "5434"),
		DBUser:      envOrDefault("DB_USER", "notix"),
		DBPassword:  envOrDefault("DB_PASSWORD", "notix_dev_password"),
		DBName:      envOrDefault("DB_NAME", "notix"),
		DBSSLMode:   envOrDefault("DB_SSLMODE", "disable"),
		JWTSecret:   []byte(secret),
		TokenExpiry: tokenExpiry,
	}, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
