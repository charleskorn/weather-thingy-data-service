package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"time"
)

type Variable struct {
	VariableID           int       `json:"id"`
	Name                 string    `json:"name" binding:"required"`
	Units                string    `json:"units" binding:"required"`
	DisplayDecimalPlaces int       `json:"displayDecimalPlaces" binding:"required"`
	Created              time.Time `json:"created"`
}

func postVariable(render render.Render, variable Variable, db Database, log *logrus.Entry) {
	if err := db.BeginTransaction(); err != nil {
		log.WithError(err).Error("Could not begin database transaction.")
		render.Error(http.StatusInternalServerError)
		return
	}

	defer db.RollbackUncommittedTransaction()

	variable.Created = time.Now()

	if err := db.CreateVariable(&variable); err != nil {
		log.WithError(err).Error("Could not create new variable.")
		render.Error(http.StatusInternalServerError)
		return
	}

	if err := db.CommitTransaction(); err != nil {
		log.WithError(err).Error("Could not commit transaction.")
		render.Error(http.StatusInternalServerError)
		return
	}

	render.JSON(http.StatusCreated, map[string]interface{}{"id": variable.VariableID})
}

func (variable Variable) Validate(errors binding.Errors, _ *http.Request) binding.Errors {
	if variable.DisplayDecimalPlaces < 0 {
		errors = append(errors, binding.Error{
			FieldNames:     []string{"displayDecimalPlaces"},
			Classification: "OutOfRangeError",
			Message:        "displayDecimalPlaces must be positive.",
		})
	}

	return errors
}
