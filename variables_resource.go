package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Variable struct {
	VariableID           int       `json:"id"`
	Name                 string    `json:"name"`
	Units                string    `json:"units"`
	DisplayDecimalPlaces int       `json:"displayDecimalPlaces"`
	Created              time.Time `json:"created"`
}

func postVariable(w http.ResponseWriter, r *http.Request, _ httprouter.Params, db Database) bool {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not read request.")
		http.Error(w, "Could not read request.", http.StatusInternalServerError)
		return false
	}

	var variable Variable

	if err := json.Unmarshal(body, &variable); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not unmarshal request.")
		http.Error(w, "Could not parse request body.", http.StatusBadRequest)
		return false
	}

	if variable.Name == "" {
		http.Error(w, "Must specify name.", http.StatusBadRequest)
		return false
	}

	if variable.Units == "" {
		http.Error(w, "Must specify units.", http.StatusBadRequest)
		return false
	}

	if variable.DisplayDecimalPlaces < 0 {
		http.Error(w, "displayDecimalPlaces must be non-negative.", http.StatusBadRequest)
		return false
	}

	variable.Created = time.Now()

	if err := db.CreateVariable(&variable); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not create new variable.")
		http.Error(w, "Could not create new variable.", http.StatusInternalServerError)
		return false
	}

	response, err := json.Marshal(map[string]interface{}{"id": variable.VariableID})

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not generate response.")
		http.Error(w, "Could not generate response.", http.StatusInternalServerError)
		return false
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
	return true
}
