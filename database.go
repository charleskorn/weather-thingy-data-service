package main

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/rubenv/sql-migrate"
)

type Database interface {
	RunMigrations() (int, error)
	Close()
	DB() *sql.DB

	CreateAgent(agent *Agent) error
}

type PostgresDatabase struct {
	DatabaseHandle *sql.DB
}

func getMigrationSource() migrate.MigrationSource {
	return &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      "db/migrations",
	}
}

func connectToDatabase(dataSourceName string) (Database, error) {
	db, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		return nil, err
	}

	return &PostgresDatabase{DatabaseHandle: db}, nil
}

func (d *PostgresDatabase) RunMigrations() (int, error) {
	migrationSource := getMigrationSource()

	n, err := migrate.Exec(d.DatabaseHandle, "postgres", migrationSource, migrate.Up)

	if err != nil {
		return 0, err
	}

	return n, nil
}

func (d *PostgresDatabase) Close() {
	d.DatabaseHandle.Close()
}

func (d *PostgresDatabase) DB() *sql.DB {
	return d.DatabaseHandle
}

func (d *PostgresDatabase) CreateAgent(agent *Agent) error {
	row := d.DB().QueryRow("INSERT INTO agents (name, created) VALUES ($1, $2) RETURNING agent_id", agent.Name, agent.Created)
	return row.Scan(&agent.AgentID)
}
