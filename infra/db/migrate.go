package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

// MigrateUp applies all pending database migrations from the specified directory.
func MigrateUp(db *sqlx.DB, dir string) (int, error) {
	migrationSrc := &migrate.FileMigrationSource{
		Dir: dir,
	}
	return migrate.Exec(db.DB, "postgres", migrationSrc, migrate.Up)
}

// MigrateDown rolls back the specified number of migrations from the specified directory.
func MigrateDown(db *sqlx.DB, dir string, steps int) (int, error) {
	migrationSrc := &migrate.FileMigrationSource{
		Dir: dir,
	}
	return migrate.ExecMax(db.DB, "postgres", migrationSrc, migrate.Down, steps)
}
