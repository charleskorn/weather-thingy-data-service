package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"strconv"
)

type Agent struct {
	AgentID int       `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

func postAgent(w http.ResponseWriter, r *http.Request, _ httprouter.Params, db Database) bool {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("Could not read request: ", err)
		http.Error(w, "Could not read request.", http.StatusInternalServerError)
		return false
	}

	var agent Agent

	if err := json.Unmarshal(body, &agent); err != nil {
		log.Println("Could not unmarshal request: ", err)
		http.Error(w, "Could not parse request body.", http.StatusBadRequest)
		return false
	}

	if agent.Name == "" {
		http.Error(w, "Must specify name.", http.StatusBadRequest)
		return false
	}

	agent.Created = time.Now()

	if err := db.CreateAgent(&agent); err != nil {
		log.Println("Could not create new agent: ", err)
		http.Error(w, "Could not create new agent.", http.StatusInternalServerError)
		return false
	}

	response, err := json.Marshal(map[string]interface{}{"id": agent.AgentID})

	if err != nil {
		log.Println("Could not generate response: ", err)
		http.Error(w, "Could not generate response.", http.StatusInternalServerError)
		return false
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
	return true
}

func getAllAgents(w http.ResponseWriter, r *http.Request, _ httprouter.Params, db Database) {
	agents, err := db.GetAllAgents()

	if err != nil {
		log.Println("Could not get all agents: ", err)
		http.Error(w, "Could not get all agents.", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(agents)

	if err != nil {
		log.Println("Could not generate response: ", err)
		http.Error(w, "Could not generate response.", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func getAgent(w http.ResponseWriter, r *http.Request, params httprouter.Params, db Database) bool {
	agentID, ok := extractAgentID(params, w, db)

	if !ok {
		return false
	}

	agent := struct {
		Agent
		Variables []Variable `json:"variables"`
	}{}

	var err error

	if agent.Agent, err = db.GetAgentByID(agentID); err != nil {
		log.Println("Could not agent info: ", err)
		http.Error(w, "Could not get agent details.", http.StatusInternalServerError)
		return false
	}

	if agent.Variables, err = db.GetVariablesForAgent(agentID); err != nil {
		log.Println("Could not get variables for agent: ", err)
		http.Error(w, "Could not get some agent details.", http.StatusInternalServerError)
		return false
	}

	response, err := json.Marshal(agent)

	if err != nil {
		log.Println("Could not generate response: ", err)
		http.Error(w, "Could not generate response.", http.StatusInternalServerError)
		return false
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return false
}

func extractAgentID(params httprouter.Params, w http.ResponseWriter, db Database) (int, bool) {
	rawAgentID := params.ByName("agent_id")
	agentID, err := strconv.Atoi(rawAgentID)

	if err != nil {
		http.Error(w, "Invalid agent ID.", http.StatusBadRequest)
		return 0, false
	}

	if exists, err := db.CheckAgentIDExists(agentID); err != nil {
		log.Println("Could not check if agent exists: ", err)
		http.Error(w, "Could not check if agent exists.", http.StatusInternalServerError)
		return 0, false
	} else if !exists {
		http.Error(w, "Agent does not exist.", http.StatusNotFound)
		return 0, false
	}

	return agentID, true
}
