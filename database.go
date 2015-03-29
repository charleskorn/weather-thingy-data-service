package main

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
	"github.com/rubenv/sql-migrate"
)

type Database interface {
	RunMigrations() (int, error)
	Close()
	BeginTransaction() error
	CommitTransaction() error
	RollbackTransaction() error

	DB() *sql.DB
	Transaction() *sql.Tx

	CreateAgent(agent *Agent) error
}

type PostgresDatabase struct {
	DatabaseHandle     *sql.DB
	CurrentTransaction *sql.Tx
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

func (d *PostgresDatabase) Transaction() *sql.Tx {
	return d.CurrentTransaction
}

func (d *PostgresDatabase) BeginTransaction() error {
	if d.CurrentTransaction != nil {
		return errors.New("Cannot call BeginTransaction when there is already a transaction in progress.")
	}

	tx, err := d.DatabaseHandle.Begin()

	if err != nil {
		return err
	}

	d.CurrentTransaction = tx
	return nil
}

func (d *PostgresDatabase) CommitTransaction() error {
	if d.CurrentTransaction == nil {
		return errors.New("Cannot call CommitTransaction when there is no transaction in progress.")
	}

	if err := d.CurrentTransaction.Commit(); err != nil {
		return err
	}

	d.CurrentTransaction = nil
	return nil
}

func (d *PostgresDatabase) RollbackTransaction() error {
	if d.CurrentTransaction == nil {
		return errors.New("Cannot call RollbackTransaction when there is no transaction in progress.")
	}

	if err := d.CurrentTransaction.Rollback(); err != nil {
		return err
	}

	d.CurrentTransaction = nil
	return nil
}

func (d *PostgresDatabase) CreateAgent(agent *Agent) error {
	if err := d.ensureTransaction(); err != nil {
		return err
	}

	row := d.CurrentTransaction.QueryRow("INSERT INTO agents (name, created) VALUES ($1, $2) RETURNING agent_id", agent.Name, agent.Created)
	return row.Scan(&agent.AgentID)
}

func (d *PostgresDatabase) ensureTransaction() error {
	if d.CurrentTransaction == nil {
		return errors.New("An active transaction is required to call this method.")
	}

	return nil
}
