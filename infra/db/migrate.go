package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

func MigrateDB(db *sqlx.DB, dir string) error {
	migrationSrc := &migrate.FileMigrationSource{
		Dir: dir,
	}

	_, err := migrate.Exec(db.DB, "postgres", migrationSrc, migrate.Up)
	if err != nil {
		return err
	}

	fmt.Println("Migration completed successfully")
	return nil
}