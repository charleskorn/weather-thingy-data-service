package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"fmt"

	"github.com/go-martini/martini"
	"github.com/golang/mock/gomock"
	"github.com/martini-contrib/render"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Data resource", func() {
	var mockController *gomock.Controller

	BeforeEach(func() {
		mockController = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockController.Finish()
	})

	Describe("POST data structure", func() {
		It("can be deserialised from JSON", func() {
			jsonString := `{"time":"2015-03-26T14:35:00Z","data":[` +
				`{"variable":"temperature","value":10.675},` +
				`{"variable":"humidity","value":90}` +
				`]}`
			var postData PostDataPoints
			err := json.Unmarshal([]byte(jsonString), &postData)

			expectedPostData := PostDataPoints{
				Time: time.Date(2015, 3, 26, 14, 35, 0, 0, time.UTC),
				Data: []PostDataPoint{
					PostDataPoint{Variable: "temperature", Value: 10.675},
					PostDataPoint{Variable: "humidity", Value: 90},
				},
			}
			Expect(err).To(BeNil())
			Expect(postData).To(Equal(expectedPostData))
		})
	})

	Describe("GET data structure", func() {
		It("can be serialised to JSON", func() {
			data := GetDataResult{
				Data: []GetDataResultVariable{
					GetDataResultVariable{
						VariableID:           10,
						Name:                 "distance",
						Units:                "metres",
						DisplayDecimalPlaces: 30,
						Points: map[string]float64{
							"2015-03-26T14:35:00Z": 15.3,
							"2015-03-26T14:40:00Z": 15.0},
					},
				},
			}

			expectedJSON := `{"data":[{` +
				`"id":10,` +
				`"name":"distance",` +
				`"units":"metres",` +
				`"displayDecimalPlaces":30,` +
				`"points":{"2015-03-26T14:35:00Z":15.3,"2015-03-26T14:40:00Z":15}` +
				`}]}`
			json, err := json.Marshal(data)
			Expect(err).To(BeNil())
			Expect(string(json)).To(MatchJSON(expectedJSON))
		})
	})

	Describe("POST request handler", func() {
		var makeRequest = func(body string, agentID string, render render.Render, db Database) {
			request, _ := http.NewRequest("POST", "/blah", strings.NewReader(body))
			params := martini.Params{"agent_id": agentID}

			postDataPoints(render, request, params, db)
		}

		var db *MockDatabase
		var render *MockRender
		existentAgentID := "10"

		BeforeEach(func() {
			db = NewMockDatabase(mockController)
			render = NewMockRender(mockController)
		})

		Describe("when the request is valid", func() {
			It("saves the data point to the database and returns HTTP 201 response", func() {
				createCall := db.EXPECT().AddDataPoint(DataPoint{
					AgentID:    10,
					VariableID: 12,
					Value:      10.5,
					Time:       time.Date(2015, 5, 6, 10, 15, 30, 0, time.UTC),
				})

				gomock.InOrder(
					db.EXPECT().BeginTransaction(),
					db.EXPECT().CheckAgentIDExists(10).Return(true, nil),
					db.EXPECT().GetVariableIDForName("temperature").Return(12, nil),
					createCall,
					db.EXPECT().CommitTransaction(),
					render.EXPECT().Status(http.StatusCreated),
					db.EXPECT().RollbackUncommittedTransaction(),
				)

				makeRequest(`{"time":"2015-05-06T10:15:30Z","data":[{"variable":"temperature","value":10.5}]}`, existentAgentID, render, db)
			})
		})

		Describe("when the request is invalid", func() {
			TheRequestFailsWithCode := func(request string, agentID string, responseCode int) {
				It(fmt.Sprintf("does not save the variable to the database and returns HTTP %v response", responseCode), func() {
					render.EXPECT().Text(responseCode, gomock.Any())

					makeRequest(request, agentID, render, db)
				})
			}

			TheRequestFails := func(request string, agentID string) {
				TheRequestFailsWithCode(request, agentID, http.StatusBadRequest)
			}

			BeforeEach(func() {
				gomock.InOrder(
					db.EXPECT().BeginTransaction(),
					db.EXPECT().CheckAgentIDExists(909090).Return(false, nil).AnyTimes(),
					db.EXPECT().CheckAgentIDExists(10).Return(true, nil).AnyTimes(),
					db.EXPECT().GetVariableIDForName("nothing").Return(-1, errors.New("Doesn't exisit")).AnyTimes(),
					db.EXPECT().RollbackUncommittedTransaction(),
				)
			})

			Describe("because there are no fields", func() {
				TheRequestFails(`{}`, existentAgentID)
			})

			Describe("because the time field is empty", func() {
				TheRequestFails(`{"time":"","data":[{"variable":"temperature","value":10}]}`, existentAgentID)
			})

			Describe("because the time field is missing", func() {
				TheRequestFails(`{"data":[{"variable":"temperature","value":10}]}`, existentAgentID)
			})

			Describe("because the time field is in an invalid format", func() {
				TheRequestFails(`{"time":"blah","data":[{"variable":"temperature","value":10}]}`, existentAgentID)
			})

			Describe("because the data field is empty", func() {
				TheRequestFails(`{"time":"2015-01-02T03:04:05Z","data":[]}`, existentAgentID)
			})

			Describe("because the variable field is empty", func() {
				TheRequestFails(`{"time":"2015-01-02T03:04:05Z","data":[{"variable":"","value":10}]}`, existentAgentID)
			})

			Describe("because the variable field is missing", func() {
				TheRequestFails(`{"time":"2015-01-02T03:04:05Z","data":[{"value":10}]}`, existentAgentID)
			})

			Describe("because the value field is empty", func() {
				TheRequestFails(`{"time":"2015-01-02T03:04:05Z","data":[{"variable":"temperature","value":""}]}`, existentAgentID)
			})

			Describe("because the value field is not a number", func() {
				TheRequestFails(`{"time":"2015-01-02T03:04:05Z","data":[{"variable":"temperature","value":"abc"}]}`, existentAgentID)
			})

			Describe("because the agent ID is not an integer", func() {
				TheRequestFailsWithCode(`{"time":"2015-01-02T03:04:05Z","data":[{"variable":"temperature","value":10}]}`, "abc", http.StatusNotFound)
			})

			Describe("because the agent ID does not match any known agent", func() {
				TheRequestFailsWithCode(`{"time":"2015-01-02T03:04:05Z","data":[{"variable":"temperature","value":10}]}`, "909090", http.StatusNotFound)
			})

			Describe("because the variable name does not match any known variable", func() {
				TheRequestFails(`{"time":"2015-01-02T03:04:05Z","data":[{"variable":"nothing","value":10}]}`, existentAgentID)
			})
		})
	})

	Describe("GET request handler", func() {
		var makeRequest = func(query string, agentID string, render render.Render, db Database) {
			request, _ := http.NewRequest("GET", "/blah?"+query, strings.NewReader(""))
			params := martini.Params{"agent_id": agentID}

			getData(render, request, params, db)
		}

		var db *MockDatabase
		var render *MockRender

		BeforeEach(func() {
			db = NewMockDatabase(mockController)
			render = NewMockRender(mockController)
		})

		Context("when the request is valid", func() {
			variable123 := Variable{Name: "temperature", Units: "°C", DisplayDecimalPlaces: 1}
			variable321 := Variable{Name: "humidity", Units: "%", DisplayDecimalPlaces: 2}

			variable123Data := map[string]float64{
				"2015-03-27T06:00:00Z": 100,
				"2015-03-27T09:00:00Z": 105,
			}

			variable321Data := map[string]float64{
				"2015-03-27T08:00:00Z": 100.5,
				"2015-03-27T12:00:00Z": 80.9,
			}

			It("returns the data", func() {
				jsonCall := render.EXPECT().JSON(http.StatusOK, gomock.Any()).Do(func(status int, value interface{}) {
					bytes, err := json.Marshal(value)
					Expect(err).To(BeNil())

					json := string(bytes)
					Expect(json).To(MatchJSON(`{"data":[` +
						`{"id":123,"name":"temperature","units":"°C","displayDecimalPlaces":1,"points":{"2015-03-27T06:00:00Z":100,"2015-03-27T09:00:00Z":105}},` +
						`{"id":321,"name":"humidity","units":"%","displayDecimalPlaces":2,"points":{"2015-03-27T08:00:00Z":100.5,"2015-03-27T12:00:00Z":80.9}}` +
						`]}`))
				})

				fromDate := time.Date(2015, 3, 27, 5, 0, 0, 0, time.UTC)
				toDate := time.Date(2015, 3, 28, 23, 50, 45, 0, time.UTC)

				gomock.InOrder(
					db.EXPECT().BeginTransaction(),
					db.EXPECT().CheckAgentIDExists(1).Return(true, nil),
					db.EXPECT().GetVariableByID(123).Return(variable123, nil),
					db.EXPECT().GetData(1, 123, fromDate, toDate).Return(variable123Data, nil),
					db.EXPECT().GetVariableByID(321).Return(variable321, nil),
					db.EXPECT().GetData(1, 321, fromDate, toDate).Return(variable321Data, nil),
					db.EXPECT().CommitTransaction(),
					jsonCall,
					db.EXPECT().RollbackUncommittedTransaction(),
				)

				makeRequest("variable=123&variable=321&date_from=2015-03-27T05:00:00Z&date_to=2015-03-28T23:50:45Z", "1", render, db)
			})
		})

		Context("when the request is invalid", func() {
			TheRequestFailsWithCode := func(query string, agentID string, responseCode int) {
				It(fmt.Sprintf("returns HTTP %v response", responseCode), func() {
					render.EXPECT().Text(responseCode, gomock.Any())
					makeRequest(query, agentID, render, db)
				})
			}

			TheRequestFails := func(query string, agentID string) {
				TheRequestFailsWithCode(query, agentID, http.StatusBadRequest)
			}

			BeforeEach(func() {
				gomock.InOrder(
					db.EXPECT().BeginTransaction(),
					db.EXPECT().CheckAgentIDExists(909090).Return(false, nil).AnyTimes(),
					db.EXPECT().CheckAgentIDExists(1).Return(true, nil).AnyTimes(),
					db.EXPECT().RollbackUncommittedTransaction(),
				)
			})

			Context("because the agent ID does not exist", func() {
				TheRequestFailsWithCode("variable=123&variable=321&date_from=2015-03-27T05:00:00Z&date_to=2015-03-28T23:50:45Z", "909090", http.StatusNotFound)
			})

			Context("because the agent ID is not an integer", func() {
				TheRequestFailsWithCode("variable=123&variable=321&date_from=2015-03-27T05:00:00Z&date_to=2015-03-28T23:50:45Z", "abc", http.StatusNotFound)
			})

			Context("because no variables are specified", func() {
				TheRequestFails("date_from=2015-03-27T05:00:00Z&date_to=2015-03-28T23:50:45Z", "1")
			})

			Context("because a variable is not an integer", func() {
				TheRequestFails("variable=123&variable=abc&date_from=2015-03-27T05:00:00Z&date_to=2015-03-28T23:50:45Z", "1")
			})

			Context("because the from date is missing", func() {
				TheRequestFails("variable=123&variable=321&date_to=2015-03-28T23:50:45Z", "1")
			})

			Context("because the from date is in an invalid format", func() {
				TheRequestFails("variable=123&variable=321&date_from=2015-03-27&date_to=2015-03-28T23:50:45Z", "1")
			})

			Context("because the to date is missing", func() {
				TheRequestFails("variable=123&variable=321&date_from=2015-03-27T05:00:00Z", "1")
			})

			Context("because the to date is in an invalid format", func() {
				TheRequestFails("variable=123&variable=321&date_from=2015-03-27T05:00:00Z&date_to=2015-03-28", "1")
			})

			Context("because the to date is before the from date", func() {
				TheRequestFails("variable=123&variable=321&date_from=2015-03-27T05:00:00Z&date_to=2015-03-26T23:50:45Z", "1")
			})
		})
	})
})
