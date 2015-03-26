package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Agent struct {
	AgentID int       `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

func postAgent(w http.ResponseWriter, r *http.Request, _ httprouter.Params, db Database) {
	var agent Agent
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Could not read request: ", err)
		http.Error(w, "Could not read request.", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &agent); err != nil {
		log.Printf("Could not unmarshal request: ", err)
		http.Error(w, "Could not parse request body.", http.StatusBadRequest)
		return
	}

	if agent.Name == "" {
		http.Error(w, "Must specify name.", http.StatusBadRequest)
		return
	}

	if err := db.CreateAgent(&agent); err != nil {
		log.Printf("Could not create new agent: ", err)
		http.Error(w, "Could not create new agent.", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(map[string]interface{}{"id": agent.AgentID})

	if err != nil {
		log.Printf("Could not generate response: ", err)
		http.Error(w, "Could not generate response.", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
