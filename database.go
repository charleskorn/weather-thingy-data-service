package main

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/rubenv/sql-migrate"
)

type Database struct {
	DatabaseHandle *sql.DB
}

func getMigrationSource() migrate.MigrationSource {
	return &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      "db/migrations",
	}
}

func connectToDatabase(dataSourceName string) (*Database, error) {
	db, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		return nil, err
	}

	return &Database{DatabaseHandle: db}, nil
}

func (d *Database) runMigrations() (int, error) {
	migrationSource := getMigrationSource()

	n, err := migrate.Exec(d.DatabaseHandle, "postgres", migrationSource, migrate.Up)

	if err != nil {
		return 0, err
	}

	return n, nil
}
