package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type DataPoint struct {
	AgentID    int
	VariableID int
	Time       time.Time
	Value      float64
}

type PostDataPoints struct {
	Time time.Time
	Data []PostDataPoint
}

type PostDataPoint struct {
	Variable string
	Value    float64
}

type GetDataResult struct {
	Data []GetDataResultVariable `json:"data"`
}

type GetDataResultVariable struct {
	VariableID int                `json:"id"`
	Name       string             `json:"name"`
	Units      string             `json:"units"`
	Points     map[string]float64 `json:"points"`
}

func postDataPoints(w http.ResponseWriter, r *http.Request, params httprouter.Params, db Database) bool {
	agentID, ok := extractAgentID(params, w, db)

	if !ok {
		return false
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("Could not read request: ", err)
		http.Error(w, "Could not read request.", http.StatusInternalServerError)
		return false
	}

	var data PostDataPoints

	if err := json.Unmarshal(body, &data); err != nil {
		log.Println("Could not unmarshal request: ", err)
		http.Error(w, "Could not parse request body.", http.StatusBadRequest)
		return false
	}

	if !validatePostRequest(data, w) {
		return false
	}

	for _, point := range data.Data {
		variableID, err := db.GetVariableIDForName(point.Variable)

		if err != nil {
			if variableID == -1 {
				http.Error(w, fmt.Sprintf("Could not find variable with name '%v'.", point.Variable), http.StatusBadRequest)
			} else {
				http.Error(w, "Could not find variable.", http.StatusInternalServerError)
			}

			return false
		}

		if err := db.AddDataPoint(DataPoint{AgentID: agentID, VariableID: variableID, Value: point.Value, Time: data.Time}); err != nil {
			log.Println("Could not save data: ", err)
			http.Error(w, "Could not save data.", http.StatusInternalServerError)
			return false
		}
	}

	w.WriteHeader(http.StatusCreated)
	return true
}

func validatePostRequest(data PostDataPoints, w http.ResponseWriter) bool {
	if data.Time.Equal(time.Time{}) {
		http.Error(w, "Must specify time value.", http.StatusBadRequest)
		return false
	}

	if len(data.Data) == 0 {
		http.Error(w, "Must include at least one data point.", http.StatusBadRequest)
		return false
	}

	for _, point := range data.Data {
		if point.Variable == "" {
			http.Error(w, "Must include variable name.", http.StatusBadRequest)
			return false
		}
	}

	return true
}

func getData(w http.ResponseWriter, r *http.Request, params httprouter.Params, db Database) bool {
	agentID, ok := extractAgentID(params, w, db)

	if !ok {
		return false
	}

	variables, fromTime, toTime, ok := extractGetParameters(w, r)

	if !ok {
		return false
	}

	result := GetDataResult{}

	for _, variableID := range variables {
		variable, err := db.GetVariableByID(variableID)

		if err != nil {
			log.Println("Could not get variable info:", err)
			http.Error(w, "Could not get variable info.", http.StatusInternalServerError)
			return false
		}

		variableResult := GetDataResultVariable{VariableID: variableID, Name: variable.Name, Units: variable.Units}
		variableResult.Points, err = db.GetData(agentID, variableID, fromTime, toTime)

		if err != nil {
			log.Println("Could not retrieve data:", err)
			http.Error(w, "Could not retrieve data.", http.StatusInternalServerError)
			return false
		}

		result.Data = append(result.Data, variableResult)
	}

	response, err := json.Marshal(result)

	if err != nil {
		log.Println("Could not generate response: ", err)
		http.Error(w, "Could not generate response.", http.StatusInternalServerError)
		return false
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

	return true
}

func extractGetParameters(w http.ResponseWriter, r *http.Request) ([]int, time.Time, time.Time, bool) {
	if r.URL.Query().Get("variable") == "" {
		http.Error(w, "Must specify variable with 'variable' URL parameter.", http.StatusBadRequest)
		return nil, time.Time{}, time.Time{}, false
	}

	if r.URL.Query().Get("date_from") == "" {
		http.Error(w, "Must specify to date with 'date_from' URL parameter.", http.StatusBadRequest)
		return nil, time.Time{}, time.Time{}, false
	}

	if r.URL.Query().Get("date_to") == "" {
		http.Error(w, "Must specify from date with 'date_to' URL parameter.", http.StatusBadRequest)
		return nil, time.Time{}, time.Time{}, false
	}

	fromDate, err := time.Parse(time.RFC3339, r.URL.Query().Get("date_from"))

	if err != nil {
		http.Error(w, "Cannot parse from date value.", http.StatusBadRequest)
		return nil, time.Time{}, time.Time{}, false
	}

	toDate, err := time.Parse(time.RFC3339, r.URL.Query().Get("date_to"))

	if err != nil {
		http.Error(w, "Cannot parse to date value.", http.StatusBadRequest)
		return nil, time.Time{}, time.Time{}, false
	}

	if fromDate.After(toDate) {
		http.Error(w, "From date is after to date.", http.StatusBadRequest)
		return nil, time.Time{}, time.Time{}, false
	}

	variables := []int{}

	for _, rawID := range r.URL.Query()["variable"] {
		id, err := strconv.Atoi(rawID)

		if err != nil {
			http.Error(w, fmt.Sprintf("Variable '%v' is not an integer.", id), http.StatusBadRequest)
			return nil, time.Time{}, time.Time{}, false
		}

		variables = append(variables, id)
	}

	return variables, fromDate, toDate, true
}
