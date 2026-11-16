package config

import (
	"strings"
	"testing"
	"time"
)

// validTestEnv returns a complete map of valid required and optional environment variables.
func validTestEnv() map[string]string {
	return map[string]string{
		"VERSION":           "1.2.3",
		"SERVICE_NAME":      "gopay-core",
		"HTTP_PORT":         "8080",
		"JWT_SECRET_KEY":    "my_test_secret_key_32_characters",
		"DB_HOST":           "localhost",
		"DB_PORT":           "5432",
		"DB_NAME":           "gopay_db",
		"DB_USER":           "postgres",
		"DB_PASSWORD":       "secret",
		"DB_ENABLE_SSL_MODE": "false",
		"REDIS_HOST":        "127.0.0.1",
		"REDIS_PORT":        "6379",
	}
}

func TestLoadFromLookup_Success(t *testing.T) {
	env := validTestEnv()
	cfg, err := LoadFromLookup(func(k string) string { return env[k] })

	if err != nil {
		t.Fatalf("expected valid config to load, got error: %v", err)
	}

	if cfg.HTTP.Port != 8080 {
		t.Errorf("HTTP.Port = %d, want 8080", cfg.HTTP.Port)
	}
	if cfg.JWT.SecretKey != "my_test_secret_key_32_characters" {
		t.Errorf("JWT.SecretKey mismatch")
	}
	if cfg.DB.Host != "localhost" || cfg.DB.Port != 5432 {
		t.Errorf("DB config mismatch: %+v", cfg.DB)
	}
	if cfg.DB.MaxConns != 25 { // default fallback
		t.Errorf("DB.MaxConns = %d, want default 25", cfg.DB.MaxConns)
	}
	if cfg.Timeouts.Shutdown != 10*time.Second {
		t.Errorf("Timeouts.Shutdown = %v, want 10s", cfg.Timeouts.Shutdown)
	}
}

func TestLoadFromLookup_AllMissingReportedAtOnce(t *testing.T) {
	// Completely empty environment
	emptyEnv := map[string]string{}
	_, err := LoadFromLookup(func(k string) string { return emptyEnv[k] })

	if err == nil {
		t.Fatal("expected error on empty environment, got nil")
	}

	errStr := err.Error()

	// Verify that ALL required keys are reported in ONE single error
	requiredKeys := []string{
		"HTTP_PORT",
		"JWT_SECRET_KEY",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
		"DB_USER",
		"DB_PASSWORD",
	}

	for _, key := range requiredKeys {
		if !strings.Contains(errStr, key) {
			t.Errorf("error message missing key %q; full error: %s", key, errStr)
		}
	}
}

func TestLoadFromLookup_PartialMissingReported(t *testing.T) {
	env := validTestEnv()
	delete(env, "HTTP_PORT")
	delete(env, "DB_PASSWORD")

	_, err := LoadFromLookup(func(k string) string { return env[k] })
	if err == nil {
		t.Fatal("expected error when keys are missing, got nil")
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "HTTP_PORT") {
		t.Errorf("expected error to mention HTTP_PORT, got: %s", errStr)
	}
	if !strings.Contains(errStr, "DB_PASSWORD") {
		t.Errorf("expected error to mention DB_PASSWORD, got: %s", errStr)
	}
	if strings.Contains(errStr, "JWT_SECRET_KEY") {
		t.Errorf("error must not mention present key JWT_SECRET_KEY, got: %s", errStr)
	}
}

func TestLoadFromLookup_InvalidIntegerReported(t *testing.T) {
	env := validTestEnv()
	env["HTTP_PORT"] = "not_a_number"
	env["DB_PORT"] = "invalid_port"

	_, err := LoadFromLookup(func(k string) string { return env[k] })
	if err == nil {
		t.Fatal("expected error for invalid integer, got nil")
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "HTTP_PORT") || !strings.Contains(errStr, "must be an integer") {
		t.Errorf("expected invalid integer error for HTTP_PORT, got: %s", errStr)
	}
	if !strings.Contains(errStr, "DB_PORT") {
		t.Errorf("expected invalid integer error for DB_PORT, got: %s", errStr)
	}
}
