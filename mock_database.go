package main

import "database/sql"

type CreateAgentInfo struct {
	Calls           []Agent
	AgentIDToReturn int
}

type GetAllAgentsInfo struct {
	AgentsToReturn []Agent
}

type MockDatabase struct {
	CreateAgentInfo  CreateAgentInfo
	GetAllAgentsInfo GetAllAgentsInfo
}

func (d *MockDatabase) RunMigrations() (int, error) {
	panic("Cannot call RunMigrations() on a MockDatabase")
}

func (d *MockDatabase) Close() {
	panic("Cannot call Close() on a MockDatabase")
}

func (d *MockDatabase) DB() *sql.DB {
	panic("Cannot call DB() on a MockDatabase")
}

func (d *MockDatabase) Transaction() *sql.Tx {
	panic("Cannot call Transaction() on a MockDatabase")
}

func (d *MockDatabase) BeginTransaction() error {
	panic("Cannot call BeginTransaction() on a MockDatabase")
}

func (d *MockDatabase) CommitTransaction() error {
	panic("Cannot call CommitTransaction() on a MockDatabase")
}

func (d *MockDatabase) RollbackTransaction() error {
	panic("Cannot call RollbackTransaction() on a MockDatabase")
}

func (d *MockDatabase) CreateAgent(agent *Agent) error {
	d.CreateAgentInfo.Calls = append(d.CreateAgentInfo.Calls, *agent)
	agent.AgentID = d.CreateAgentInfo.AgentIDToReturn

	return nil
}

func (d *MockDatabase) GetAllAgents() ([]Agent, error) {
	return d.GetAllAgentsInfo.AgentsToReturn, nil
}
