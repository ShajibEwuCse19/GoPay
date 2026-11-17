package db

import (
	"github.com/ShajibEwuCse19/gopay-core/migrations"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

// getMigrationSource returns an embedded migration source by default,
// or a disk-based FileMigrationSource if a physical directory path is provided.
func getMigrationSource(dir string) migrate.MigrationSource {
	if dir == "" || dir == "embedded" {
		return &migrate.EmbedFileSystemMigrationSource{
			FileSystem: migrations.FS,
			Root:       "/",
		}
	}

	return &migrate.FileMigrationSource{
		Dir: dir,
	}
}

// MigrateUp applies all pending database migrations (embedded or from specified directory).
func MigrateUp(db *sqlx.DB, dir string) (int, error) {
	source := getMigrationSource(dir)
	return migrate.Exec(db.DB, "postgres", source, migrate.Up)
}

// MigrateDown rolls back the specified number of migrations (embedded or from specified directory).
func MigrateDown(db *sqlx.DB, dir string, steps int) (int, error) {
	source := getMigrationSource(dir)
	return migrate.ExecMax(db.DB, "postgres", source, migrate.Down, steps)
}
