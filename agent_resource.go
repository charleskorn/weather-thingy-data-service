package main

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"strconv"
)

type Agent struct {
	AgentID     int       `json:"id"`
	OwnerUserID int       `json:"ownerUserId"`
	Name        string    `json:"name" binding:"required"`
	Created     time.Time `json:"created"`
}

func postAgent(r render.Render, agent Agent, db Database, log *logrus.Entry) {
	agent.Created = time.Now()
	agent.OwnerUserID = -1

	if err := db.BeginTransaction(); err != nil {
		log.WithError(err).Error("Could not begin database transaction.")
		r.Error(http.StatusInternalServerError)
		return
	}

	defer db.RollbackUncommittedTransaction()

	if err := db.CreateAgent(&agent); err != nil {
		log.WithError(err).Error("Could not create new agent.")
		r.Error(http.StatusInternalServerError)
		return
	}

	if err := db.CommitTransaction(); err != nil {
		log.WithError(err).Error("Could not commit transaction.")
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(http.StatusCreated, map[string]interface{}{"id": agent.AgentID})
}

func getAllAgents(r render.Render, db Database, log *logrus.Entry) {
	agents, err := db.GetAllAgents()

	if err != nil {
		log.WithError(err).Error("Could not get all agents.")
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(http.StatusOK, agents)
}

func getAgent(r render.Render, params martini.Params, db Database, log *logrus.Entry) {
	if err := db.BeginTransaction(); err != nil {
		log.WithError(err).Error("Could not begin database transaction.")
		r.Error(http.StatusInternalServerError)
		return
	}

	defer db.RollbackUncommittedTransaction()

	agentID, ok := extractAgentID(params, r, db, log)

	if !ok {
		return
	}

	agent := struct {
		Agent
		Variables []Variable `json:"variables"`
	}{}

	var err error

	if agent.Agent, err = db.GetAgentByID(agentID); err != nil {
		log.WithError(err).Error("Could not agent info.")
		r.Error(http.StatusInternalServerError)
		return
	}

	if agent.Variables, err = db.GetVariablesForAgent(agentID); err != nil {
		log.WithError(err).Error("Could not get variables for agent.")
		r.Error(http.StatusInternalServerError)
		return
	}

	if err := db.CommitTransaction(); err != nil {
		log.WithError(err).Error("Could not commit transaction.")
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(http.StatusOK, agent)
}

func extractAgentID(params martini.Params, r render.Render, db Database, log *logrus.Entry) (int, bool) {
	rawAgentID := params["agent_id"]
	agentID, err := strconv.Atoi(rawAgentID)

	if err != nil {
		r.Text(http.StatusNotFound, "Invalid agent ID.")
		return 0, false
	}

	if exists, err := db.CheckAgentIDExists(agentID); err != nil {
		log.WithError(err).Error("Could not check if agent exists.")
		r.Error(http.StatusInternalServerError)
		return 0, false
	} else if !exists {
		r.Text(http.StatusNotFound, "Agent does not exist.")
		return 0, false
	}

	return agentID, true
}
