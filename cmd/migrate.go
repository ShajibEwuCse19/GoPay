package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/ShajibEwuCse19/gopay-core/config"
	"github.com/ShajibEwuCse19/gopay-core/infra/db"
	"github.com/ShajibEwuCse19/gopay-core/infra/logger"
)

// Migrate routes migration sub-commands ('up', 'down') with an optional directory override.
func Migrate(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: go run . migrate [up|down] [migrations_dir]")
		os.Exit(1)
	}

	action := args[0]
	if action != "up" && action != "down" {
		fmt.Fprintf(os.Stderr, "Unknown migrate command: %q\nUsage: go run . migrate [up|down] [migrations_dir]\n", action)
		os.Exit(1)
	}

	// 1. Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config error:\n%v\n", err)
		os.Exit(1)
	}

	// 2. Initialize logger
	logger.Init(cfg.ServiceName, cfg.LogLevel)

	// 3. Connect to Database
	dbCon, err := db.NewConnection(cfg)
	if err != nil {
		slog.Error("Database connection failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer dbCon.Close()

	// 4. Resolve migrations directory: CLI argument > Config/Env (DB_MIGRATIONS_PATH) > default "embedded"
	migrationsDir := cfg.DB.MigrationsPath
	if len(args) > 1 && strings.TrimSpace(args[1]) != "" {
		migrationsDir = strings.TrimSpace(args[1])
	}

	sourceLabel := migrationsDir
	if sourceLabel == "" || sourceLabel == "embedded" {
		sourceLabel = "embedded://migrations"
	}

	// 5. Execute Migration based on action
	switch action {
	case "up":
		n, err := db.MigrateUp(dbCon, migrationsDir)
		if err != nil {
			slog.Error("Failed to apply migrations", slog.String("source", sourceLabel), slog.String("error", err.Error()))
			os.Exit(1)
		}
		slog.Info("Migrations applied successfully", slog.String("source", sourceLabel), slog.Int("applied_count", n))

	case "down":
		// Safely rollback 1 step
		n, err := db.MigrateDown(dbCon, migrationsDir, 1)
		if err != nil {
			slog.Error("Failed to rollback migration", slog.String("source", sourceLabel), slog.String("error", err.Error()))
			os.Exit(1)
		}
		slog.Info("Migration rolled back successfully", slog.String("source", sourceLabel), slog.Int("rolled_back_count", n))
	}
}
