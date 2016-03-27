package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"net/http"
	"strconv"
	"time"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

type DataPoint struct {
	AgentID    int
	VariableID int
	Time       time.Time
	Value      float64
}

type PostDataPoints struct {
	Time time.Time `json:"time" binding:"required"`
	Data []PostDataPoint
}

type PostDataPoint struct {
	Variable string `json:"variable" binding:"required"`
	Value    float64
}

type GetDataResult struct {
	Data []GetDataResultVariable `json:"data"`
}

type GetDataResultVariable struct {
	VariableID           int                `json:"id"`
	Name                 string             `json:"name"`
	Units                string             `json:"units"`
	DisplayDecimalPlaces int                `json:"displayDecimalPlaces"`
	Points               map[string]float64 `json:"points"`
}

func postDataPoints(render render.Render, data PostDataPoints, agent Agent, db Database, log *logrus.Entry) {
	if err := db.BeginTransaction(); err != nil {
		log.WithError(err).Error("Could not begin database transaction.")
		render.Error(http.StatusInternalServerError)
		return
	}

	defer db.RollbackUncommittedTransaction()

	for _, point := range data.Data {
		variableID, err := db.GetVariableIDForName(point.Variable)

		if err != nil {
			if variableID == -1 {
				render.Text(http.StatusBadRequest, fmt.Sprintf("Could not find variable with name '%v'.", point.Variable))
			} else {
				render.Error(http.StatusInternalServerError)
			}

			return
		}

		if err := db.AddDataPoint(DataPoint{AgentID: agent.AgentID, VariableID: variableID, Value: point.Value, Time: data.Time}); err != nil {
			log.WithError(err).Error("Could not save data.")
			render.Error(http.StatusInternalServerError)
			return
		}
	}

	if err := db.CommitTransaction(); err != nil {
		log.WithError(err).Error("Could not commit transaction.")
		render.Error(http.StatusInternalServerError)
		return
	}

	render.Status(http.StatusCreated)
}

func (data PostDataPoints) Validate(errors binding.Errors, _ *http.Request) binding.Errors {
	if len(data.Data) == 0 {
		errors = append(errors, binding.Error{
			FieldNames:     []string{"data"},
			Classification: binding.RequiredError,
			Message:        "Must include at least one data point.",
		})
	}

	// HACK: Until https://github.com/martini-contrib/binding/issues/40 is fixed, we have to validate that the variable name is present ourselves.
	for _, point := range data.Data {
		if point.Variable == "" {
			errors = append(errors, binding.Error{
				FieldNames:     []string{"variable"},
				Classification: binding.RequiredError,
				Message:        "Must provide a variable name for each data point.",
			})
		}
	}

	return errors
}

func getData(render render.Render, req *http.Request, params martini.Params, db Database, user User, log *logrus.Entry) {
	if err := db.BeginTransaction(); err != nil {
		log.WithError(err).Error("Could not begin database transaction.")
		render.Error(http.StatusInternalServerError)
		return
	}

	defer db.RollbackUncommittedTransaction()

	agentID, ok := extractAgentID(params, render, db, log)

	if !ok {
		return
	}

	agent, err := db.GetAgentByID(agentID)

	if err != nil {
		log.WithError(err).Error("Could not get agent.")
		render.Error(http.StatusInternalServerError)
		return
	}

	if agent.OwnerUserID != user.UserID {
		log.WithError(err).Error("User does not own this agent.")
		render.Error(http.StatusForbidden)
		return
	}

	variables, fromTime, toTime, ok := extractGetParameters(render, req)

	if !ok {
		return
	}

	result := GetDataResult{}

	for _, variableID := range variables {
		variable, err := db.GetVariableByID(variableID)

		if err != nil {
			log.WithError(err).Error("Could not get variable info.")
			render.Error(http.StatusInternalServerError)
			return
		}

		variableResult := GetDataResultVariable{VariableID: variableID, Name: variable.Name, Units: variable.Units, DisplayDecimalPlaces: variable.DisplayDecimalPlaces}
		variableResult.Points, err = db.GetData(agentID, variableID, fromTime, toTime)

		if err != nil {
			log.WithError(err).Error("Could not retrieve data.")
			render.Error(http.StatusInternalServerError)
			return
		}

		result.Data = append(result.Data, variableResult)
	}

	if err := db.CommitTransaction(); err != nil {
		log.WithError(err).Error("Could not commit transaction.")
		render.Error(http.StatusInternalServerError)
		return
	}

	render.JSON(http.StatusOK, result)
}

func extractGetParameters(render render.Render, req *http.Request) ([]int, time.Time, time.Time, bool) {
	if req.URL.Query().Get("variable") == "" {
		render.Text(http.StatusBadRequest, "Must specify variable with 'variable' URL parameter.")
		return nil, time.Time{}, time.Time{}, false
	}

	if req.URL.Query().Get("date_from") == "" {
		render.Text(http.StatusBadRequest, "Must specify to date with 'date_from' URL parameter.")
		return nil, time.Time{}, time.Time{}, false
	}

	if req.URL.Query().Get("date_to") == "" {
		render.Text(http.StatusBadRequest, "Must specify from date with 'date_to' URL parameter.")
		return nil, time.Time{}, time.Time{}, false
	}

	fromDate, err := time.Parse(time.RFC3339, req.URL.Query().Get("date_from"))

	if err != nil {
		render.Text(http.StatusBadRequest, "Cannot parse from date value.")
		return nil, time.Time{}, time.Time{}, false
	}

	toDate, err := time.Parse(time.RFC3339, req.URL.Query().Get("date_to"))

	if err != nil {
		render.Text(http.StatusBadRequest, "Cannot parse to date value.")
		return nil, time.Time{}, time.Time{}, false
	}

	if fromDate.After(toDate) {
		render.Text(http.StatusBadRequest, "From date is after to date.")
		return nil, time.Time{}, time.Time{}, false
	}

	variables := []int{}

	for _, rawID := range req.URL.Query()["variable"] {
		id, err := strconv.Atoi(rawID)

		if err != nil {
			render.Text(http.StatusBadRequest, fmt.Sprintf("Variable '%v' is not an integer.", id))
			return nil, time.Time{}, time.Time{}, false
		}

		variables = append(variables, id)
	}

	return variables, fromDate, toDate, true
}
