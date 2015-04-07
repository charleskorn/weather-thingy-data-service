package main

import "time"

type DataPoint struct {
	AgentID    int
	VariableID int
	Time       time.Time
	Value      float64
}
