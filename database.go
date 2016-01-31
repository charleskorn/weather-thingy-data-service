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
	RollbackUncommittedTransaction() error

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
	GetVariablesForAgent(agentID int) ([]Variable, error)
	GetAgentByID(agentID int) (Agent, error)
	CreateUser(user *User) error
}

func getMigrationSource() migrate.MigrationSource {
	return &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      "db/migrations",
	}
}
