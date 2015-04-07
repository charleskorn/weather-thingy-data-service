package main

import "time"

type DataPoint struct {
	AgentID    int
	VariableID int
	Time       time.Time
	Value      float64
}

type PostDataPoints struct {
	Time time.Time
	Data []PostDataPoint
}

type PostDataPoint struct {
	Variable string
	Value    float64
}
