package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/martini-contrib/render"
	"io/ioutil"
	"net/http"
	"time"
)

type Variable struct {
	VariableID           int       `json:"id"`
	Name                 string    `json:"name"`
	Units                string    `json:"units"`
	DisplayDecimalPlaces int       `json:"displayDecimalPlaces"`
	Created              time.Time `json:"created"`
}

func postVariable(render render.Render, r *http.Request, db Database) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not read request.")
		render.Error(http.StatusInternalServerError)
		return
	}

	var variable Variable

	if err := json.Unmarshal(body, &variable); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not unmarshal request.")
		render.Text(http.StatusBadRequest, "Could not parse request.")
		return
	}

	if variable.Name == "" {
		render.Text(http.StatusBadRequest, "Must specify name.")
		return
	}

	if variable.Units == "" {
		render.Text(http.StatusBadRequest, "Must specify units.")
		return
	}

	if variable.DisplayDecimalPlaces < 0 {
		render.Text(http.StatusBadRequest, "displayDecimalPlaces must be non-negative.")
		return
	}

	if err := db.BeginTransaction(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not begin database transaction.")
		render.Error(http.StatusInternalServerError)
		return
	}

	defer db.RollbackUncommittedTransaction()

	variable.Created = time.Now()

	if err := db.CreateVariable(&variable); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not create new variable.")
		render.Error(http.StatusInternalServerError)
		return
	}

	if err := db.CommitTransaction(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not commit transaction.")
		render.Error(http.StatusInternalServerError)
		return
	}

	render.JSON(http.StatusCreated, map[string]interface{}{"id": variable.VariableID})
}
