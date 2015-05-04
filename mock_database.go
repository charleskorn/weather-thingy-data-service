package main

import (
	"database/sql"
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

type ValueSet map[string]float64
type VariableValues map[int]ValueSet
type AgentValues map[int]VariableValues

type GetDataInfo struct {
	Values AgentValues
}

type GetVariableByIDInfo struct {
	Variables map[int]Variable
}

type GetVariablesForAgentInfo struct {
	Variables []Variable
}

type GetAgentByIDInfo struct {
	Agents map[int]Agent
}

type MockDatabase struct {
	CreateAgentInfo          CreateAgentInfo
	GetAllAgentsInfo         GetAllAgentsInfo
	CreateVariableInfo       CreateVariableInfo
	AddDataPointInfo         AddDataPointInfo
	CheckAgentIDExistsInfo   CheckAgentIDExistsInfo
	GetVariableIDForNameInfo GetVariableIDForNameInfo
	GetDataInfo              GetDataInfo
	GetVariableByIDInfo      GetVariableByIDInfo
	GetVariablesForAgentInfo GetVariablesForAgentInfo
	GetAgentByIDInfo         GetAgentByIDInfo
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
		return -1, fmt.Errorf("Cannot find variable with name '%s'.", name)
	}

	return id, nil
}

func (d *MockDatabase) GetData(agentID int, variableID int, fromDate time.Time, toDate time.Time) (map[string]float64, error) {
	if agent, ok := d.GetDataInfo.Values[agentID]; !ok {
		return nil, fmt.Errorf("No mocked data for agent with ID %v.", agentID)
	} else {
		if values, ok := agent[variableID]; !ok {
			return ValueSet{}, nil
		} else {
			return values, nil
		}
	}
}

func (d *MockDatabase) GetVariableByID(variableID int) (Variable, error) {
	if variable, ok := d.GetVariableByIDInfo.Variables[variableID]; !ok {
		return Variable{}, fmt.Errorf("No mocked variable with ID %v.", variableID)
	} else {
		return variable, nil
	}
}

func (d *MockDatabase) GetVariablesForAgent(agentID int) ([]Variable, error) {
	return d.GetVariablesForAgentInfo.Variables, nil
}

func (d *MockDatabase) GetAgentByID(agentID int) (Agent, error) {
	if agent, ok := d.GetAgentByIDInfo.Agents[agentID]; !ok {
		return Agent{}, fmt.Errorf("No mocked agent with ID %v.", agentID)
	} else {
		return agent, nil
	}
}
