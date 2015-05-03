package main

import (
	"database/sql"
	"time"

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
	GetAllAgents() ([]Agent, error)
	CreateVariable(variable *Variable) error
	AddDataPoint(dataPoint DataPoint) error
	CheckAgentIDExists(agentID int) (bool, error)
	GetVariableIDForName(name string) (int, error)
	GetData(agentID int, variableID int, fromDate time.Time, toDate time.Time) (map[string]float64, error)
	GetVariableByID(variableID int) (Variable, error)
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
