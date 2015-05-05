package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/rubenv/sql-migrate"
)

type PostgresDatabase struct {
	DatabaseHandle     *sql.DB
	CurrentTransaction *sql.Tx
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

func (d *PostgresDatabase) GetAllAgents() ([]Agent, error) {
	rows, err := d.DB().Query("SELECT agent_id, name, created FROM agents;")

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	agents := []Agent{}

	for rows.Next() {
		agent := Agent{}

		if err := rows.Scan(&agent.AgentID, &agent.Name, &agent.Created); err != nil {
			return nil, err
		}

		agents = append(agents, agent)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return agents, nil
}

func (d *PostgresDatabase) CreateVariable(variable *Variable) error {
	if err := d.ensureTransaction(); err != nil {
		return err
	}

	row := d.CurrentTransaction.QueryRow("INSERT INTO variables (name, units, display_decimal_places, created) "+
		"VALUES ($1, $2, $3, $4) RETURNING variable_id", variable.Name, variable.Units, variable.DisplayDecimalPlaces, variable.Created)
	return row.Scan(&variable.VariableID)
}

func (d *PostgresDatabase) AddDataPoint(dataPoint DataPoint) error {
	if err := d.ensureTransaction(); err != nil {
		return err
	}

	_, err := d.CurrentTransaction.Exec("INSERT INTO data (agent_id, variable_id, time, value) VALUES ($1, $2, $3, $4);", dataPoint.AgentID, dataPoint.VariableID, dataPoint.Time, dataPoint.Value)
	return err
}

func (d *PostgresDatabase) CheckAgentIDExists(agentID int) (bool, error) {
	if err := d.ensureTransaction(); err != nil {
		return false, err
	}

	row := d.CurrentTransaction.QueryRow("SELECT COUNT(*) FROM agents WHERE agent_id = $1;", agentID)
	count := 0

	if err := row.Scan(&count); err != nil {
		return false, err
	}

	return (count > 0), nil
}

func (d *PostgresDatabase) GetVariableIDForName(name string) (int, error) {
	if err := d.ensureTransaction(); err != nil {
		return 0, err
	}

	rows, err := d.CurrentTransaction.Query("SELECT variable_id FROM variables WHERE name = $1;", name)

	if err != nil {
		return 0, err
	}

	defer rows.Close()

	if !rows.Next() {
		return -1, fmt.Errorf("Cannot find variable with name '%s'.", name)
	}

	var variableID int
	if err := rows.Scan(&variableID); err != nil {
		return 0, err
	}

	return variableID, nil
}

func (d *PostgresDatabase) GetData(agentID int, variableID int, fromDate time.Time, toDate time.Time) (map[string]float64, error) {
	rows, err := d.DB().Query("SELECT value, time FROM data WHERE agent_id = $1 AND variable_id = $2 AND time >= $3 AND time <= $4;",
		agentID, variableID, fromDate, toDate)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	m := map[string]float64{}

	for rows.Next() {
		var value float64
		var t time.Time

		if err := rows.Scan(&value, &t); err != nil {
			return nil, err
		}

		m[t.In(time.UTC).Format(time.RFC3339)] = value
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return m, nil
}

func (d *PostgresDatabase) GetVariableByID(variableID int) (Variable, error) {
	if err := d.ensureTransaction(); err != nil {
		return Variable{}, err
	}

	variable := Variable{}
	row := d.CurrentTransaction.QueryRow("SELECT variable_id, name, units, display_decimal_places, created FROM variables WHERE variable_id = $1;", variableID)

	if err := row.Scan(&variable.VariableID, &variable.Name, &variable.Units, &variable.DisplayDecimalPlaces, &variable.Created); err != nil {
		return Variable{}, err
	}

	return variable, nil
}

func (d *PostgresDatabase) GetVariablesForAgent(agentID int) ([]Variable, error) {
	if err := d.ensureTransaction(); err != nil {
		return []Variable{}, err
	}

	rows, err := d.CurrentTransaction.Query("SELECT variable_id, name, units, display_decimal_places, created FROM variables "+
		"WHERE variable_id IN (SELECT DISTINCT variable_id FROM data WHERE agent_id = $1);",
		agentID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	variables := []Variable{}

	for rows.Next() {
		variable := Variable{}

		if err := rows.Scan(&variable.VariableID, &variable.Name, &variable.Units,
			&variable.DisplayDecimalPlaces, &variable.Created); err != nil {
			return nil, err
		}

		variables = append(variables, variable)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return variables, nil
}

func (d *PostgresDatabase) GetAgentByID(agentID int) (Agent, error) {
	if err := d.ensureTransaction(); err != nil {
		return Agent{}, err
	}

	agent := Agent{}
	row := d.CurrentTransaction.QueryRow("SELECT agent_id, name, created FROM agents WHERE agent_id = $1;", agentID)

	if err := row.Scan(&agent.AgentID, &agent.Name, &agent.Created); err != nil {
		return Agent{}, err
	}

	return agent, nil
}

func (d *PostgresDatabase) ensureTransaction() error {
	if d.CurrentTransaction == nil {
		return errors.New("An active transaction is required to call this method.")
	}

	return nil
}
