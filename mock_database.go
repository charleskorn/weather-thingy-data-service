package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type CreateAgentInfo struct {
	Calls           []Agent
	AgentIDToReturn int
}

type GetAllAgentsInfo struct {
	AgentsToReturn []Agent
}

type CreateVariableInfo struct {
	Calls              []Variable
	VariableIDToReturn int
}

type AddDataPointInfo struct {
	Calls []DataPoint
}

type CheckAgentIDExistsInfo struct {
	AgentIDs []int
}

type GetVariableIDForNameInfo struct {
	Variables map[string]int
}

type MockDatabase struct {
	CreateAgentInfo          CreateAgentInfo
	GetAllAgentsInfo         GetAllAgentsInfo
	CreateVariableInfo       CreateVariableInfo
	AddDataPointInfo         AddDataPointInfo
	CheckAgentIDExistsInfo   CheckAgentIDExistsInfo
	GetVariableIDForNameInfo GetVariableIDForNameInfo
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

func (d *MockDatabase) CreateVariable(variable *Variable) error {
	d.CreateVariableInfo.Calls = append(d.CreateVariableInfo.Calls, *variable)
	variable.VariableID = d.CreateVariableInfo.VariableIDToReturn

	return nil
}

func (d *MockDatabase) AddDataPoint(dataPoint DataPoint) error {
	d.AddDataPointInfo.Calls = append(d.AddDataPointInfo.Calls, dataPoint)
	return nil
}

func (d *MockDatabase) CheckAgentIDExists(agentID int) (bool, error) {
	for _, id := range d.CheckAgentIDExistsInfo.AgentIDs {
		if id == agentID {
			return true, nil
		}
	}

	return false, nil
}

func (d *MockDatabase) GetVariableIDForName(name string) (int, error) {
	id, ok := d.GetVariableIDForNameInfo.Variables[name]

	if !ok {
		return -1, errors.New(fmt.Sprintf("Cannot find variable with name '%s'.", name))
	}

	return id, nil
}

func (d *MockDatabase) GetData(agentID int, variableID int, fromDate time.Time, toDate time.Time) (map[string]float64, error) {
	panic("Cannot call GetData() on a MockDatabase")
}
