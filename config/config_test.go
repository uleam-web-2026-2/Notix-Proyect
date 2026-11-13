package config

import (
	"strings"
	"testing"
)

func TestLoadRequiresDatabasePassword(t *testing.T) {
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("JWT_SECRET", "test-jwt-secret-with-at-least-32-characters")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "DB_PASSWORD must be set") {
		t.Fatalf("Load() error = %v, want DB_PASSWORD required error", err)
	}
}

func TestLoadRequiresJWTSecret(t *testing.T) {
	t.Setenv("DB_PASSWORD", "configured-password")
	t.Setenv("JWT_SECRET", "")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "JWT_SECRET must contain at least 32 characters") {
		t.Fatalf("Load() error = %v, want JWT_SECRET length error", err)
	}
}

func TestLoadUsesEnvironmentSecrets(t *testing.T) {
	const (
		password = "configured-password"
		secret   = "test-jwt-secret-with-at-least-32-characters"
	)
	t.Setenv("DB_PASSWORD", password)
	t.Setenv("JWT_SECRET", secret)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.DBPassword != password {
		t.Errorf("DBPassword = %q, want configured environment value", cfg.DBPassword)
	}
	if string(cfg.JWTSecret) != secret {
		t.Errorf("JWTSecret = %q, want configured environment value", cfg.JWTSecret)
	}
}
