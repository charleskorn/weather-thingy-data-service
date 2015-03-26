package main

import (
	"time"
)

type Agent struct {
	AgentID int       `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}
