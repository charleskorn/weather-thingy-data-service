package main

import "database/sql"

type CreateAgentInfo struct {
	Calls           []Agent
	AgentIDToReturn int
}

type MockDatabase struct {
	CreateAgentInfo CreateAgentInfo
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

func (d *MockDatabase) CreateAgent(agent *Agent) error {
	d.CreateAgentInfo.Calls = append(d.CreateAgentInfo.Calls, *agent)
	agent.AgentID = d.CreateAgentInfo.AgentIDToReturn

	return nil
}
