package config

import (
	"os"
	"testing"
	"time"
)

func TestDurationEnv(t *testing.T) {
	const key = "TEST_DURATION_VAR"
	defer os.Unsetenv(key)

	t.Run("returns parsed duration when valid", func(t *testing.T) {
		os.Setenv(key, "30s")
		got := durationEnv(key, 5*time.Second)
		if got != 30*time.Second {
			t.Errorf("got %v, want 30s", got)
		}
	})

	t.Run("returns default when env is empty", func(t *testing.T) {
		os.Unsetenv(key)
		got := durationEnv(key, 10*time.Second)
		if got != 10*time.Second {
			t.Errorf("got %v, want 10s", got)
		}
	})

	t.Run("returns default when env is invalid string", func(t *testing.T) {
		os.Setenv(key, "invalid_time")
		got := durationEnv(key, 10*time.Second)
		if got != 10*time.Second {
			t.Errorf("got %v, want 10s", got)
		}
	})

	t.Run("returns default when duration is zero or negative", func(t *testing.T) {
		os.Setenv(key, "0s")
		if got := durationEnv(key, 10*time.Second); got != 10*time.Second {
			t.Errorf("got %v, want 10s", got)
		}

		os.Setenv(key, "-5s")
		if got := durationEnv(key, 10*time.Second); got != 10*time.Second {
			t.Errorf("got %v, want 10s", got)
		}
	})
}

func TestIntEnv(t *testing.T) {
	const key = "TEST_INT_VAR"
	defer os.Unsetenv(key)

	t.Run("returns parsed int when valid", func(t *testing.T) {
		os.Setenv(key, "8080")
		got := intEnv(key, 3000)
		if got != 8080 {
			t.Errorf("got %d, want 8080", got)
		}
	})

	t.Run("returns default when unset", func(t *testing.T) {
		os.Unsetenv(key)
		got := intEnv(key, 5432)
		if got != 5432 {
			t.Errorf("got %d, want 5432", got)
		}
	})

	t.Run("returns default when invalid integer", func(t *testing.T) {
		os.Setenv(key, "not_an_int")
		got := intEnv(key, 6379)
		if got != 6379 {
			t.Errorf("got %d, want 6379", got)
		}
	})
}

func TestGetEnv(t *testing.T) {
	const key = "TEST_STRING_VAR"
	defer os.Unsetenv(key)

	os.Setenv(key, "custom_val")
	if got := getEnv(key, "fallback"); got != "custom_val" {
		t.Errorf("got %s, want custom_val", got)
	}

	os.Unsetenv(key)
	if got := getEnv(key, "fallback"); got != "fallback" {
		t.Errorf("got %s, want fallback", got)
	}
}

func TestConfigSections(t *testing.T) {
	cfg := &Config{
		Version:     "1.0.0",
		ServiceName: "GoPay",
		HTTP: HTTPConfig{
			Port:         8080,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		DB: DBConfig{
			Host:     "localhost",
			Port:     5432,
			Name:     "gopay",
			MaxConns: 25,
			MinConns: 5,
		},
		JWT: JWTConfig{
			SecretKey:       "super_secret",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 7 * 24 * time.Hour,
		},
		Redis: RedisConfig{
			Host: "localhost",
			Port: 6379,
		},
		Timeouts: TimeoutConfig{
			CreateCustomer: 5 * time.Second,
			Shutdown:       10 * time.Second,
		},
	}

	if cfg.HTTP.Port != 8080 {
		t.Errorf("HTTP.Port = %d, want 8080", cfg.HTTP.Port)
	}
	if cfg.DB.Port != 5432 {
		t.Errorf("DB.Port = %d, want 5432", cfg.DB.Port)
	}
	if cfg.JWT.AccessTokenTTL != 15*time.Minute {
		t.Errorf("JWT.AccessTokenTTL = %v, want 15m", cfg.JWT.AccessTokenTTL)
	}
	if cfg.Redis.Port != 6379 {
		t.Errorf("Redis.Port = %d, want 6379", cfg.Redis.Port)
	}
	if cfg.Timeouts.Shutdown != 10*time.Second {
		t.Errorf("Timeouts.Shutdown = %v, want 10s", cfg.Timeouts.Shutdown)
	}
}
