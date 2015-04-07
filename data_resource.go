package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
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

		db.AddDataPoint(DataPoint{AgentID: agentID, VariableID: variableID, Value: point.Value, Time: data.Time})
	}

	w.WriteHeader(http.StatusCreated)
	return true
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
