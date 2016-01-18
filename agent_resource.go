package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"strconv"
)

type Agent struct {
	AgentID int       `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

func postAgent(r render.Render, req *http.Request, db Database) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not read request.")
		r.Error(http.StatusInternalServerError)
		return
	}

	var agent Agent

	if err := json.Unmarshal(body, &agent); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not unmarshal request.")
		r.Error(http.StatusBadRequest)
		return
	}

	if agent.Name == "" {
		r.Text(http.StatusBadRequest, "Must specify name.")
		return
	}

	agent.Created = time.Now()

	if err := db.BeginTransaction(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not begin database transaction.")
		r.Error(http.StatusInternalServerError)
		return
	}

	defer db.RollbackUncommittedTransaction()

	if err := db.CreateAgent(&agent); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not create new agent.")
		r.Error(http.StatusInternalServerError)
		return
	}

	if err := db.CommitTransaction(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not commit transaction.")
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(http.StatusCreated, map[string]interface{}{"id": agent.AgentID})
}

func getAllAgents(r render.Render, db Database) {
	agents, err := db.GetAllAgents()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not get all agents.")
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(http.StatusOK, agents)
}

func getAgent(r render.Render, params martini.Params, db Database) {
	if err := db.BeginTransaction(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not begin database transaction.")
		r.Error(http.StatusInternalServerError)
		return
	}

	defer db.RollbackUncommittedTransaction()

	agentID, ok := extractAgentID(params, r, db)

	if !ok {
		return
	}

	agent := struct {
		Agent
		Variables []Variable `json:"variables"`
	}{}

	var err error

	if agent.Agent, err = db.GetAgentByID(agentID); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not agent info.")
		r.Error(http.StatusInternalServerError)
		return
	}

	if agent.Variables, err = db.GetVariablesForAgent(agentID); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not get variables for agent.")
		r.Error(http.StatusInternalServerError)
		return
	}

	if err := db.CommitTransaction(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not commit transaction.")
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(http.StatusOK, agent)
}

func extractAgentID(params martini.Params, r render.Render, db Database) (int, bool) {
	rawAgentID := params["agent_id"]
	agentID, err := strconv.Atoi(rawAgentID)

	if err != nil {
		r.Text(http.StatusNotFound, "Invalid agent ID.")
		return 0, false
	}

	if exists, err := db.CheckAgentIDExists(agentID); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not check if agent exists.")
		r.Error(http.StatusInternalServerError)
		return 0, false
	} else if !exists {
		r.Text(http.StatusNotFound, "Agent does not exist.")
		return 0, false
	}

	return agentID, true
}
