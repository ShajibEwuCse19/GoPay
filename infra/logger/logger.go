package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Init initializes and sets the global slog JSON logger with service name and caller source.
func Init(serviceName, levelStr string) *slog.Logger {
	var level slog.Level
	switch strings.ToLower(strings.TrimSpace(levelStr)) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true, 
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)

	logger := slog.New(handler).With(
		slog.String("service", serviceName),
	)

	slog.SetDefault(logger)

	return logger
}