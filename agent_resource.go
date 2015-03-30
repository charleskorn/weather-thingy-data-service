package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Agent struct {
	AgentID int       `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

func postAgent(w http.ResponseWriter, r *http.Request, _ httprouter.Params, db Database) bool {
	var agent Agent
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("Could not read request: ", err)
		http.Error(w, "Could not read request.", http.StatusInternalServerError)
		return false
	}

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
	return
}
