package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"fmt"
	"github.com/julienschmidt/httprouter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Data resource", func() {
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
						VariableID: 10,
						Name:       "distance",
						Units:      "metres",
						Points: map[string]float64{
							"2015-03-26T14:35:00Z": 15.3,
							"2015-03-26T14:40:00Z": 15.0},
					},
				},
			}

			expectedJson := `{"data":[` +
				`{"id":10,"name":"distance","units":"metres","points":{` +
				`"2015-03-26T14:35:00Z":15.3,"2015-03-26T14:40:00Z":15` +
				`}}]}`
			json, err := json.Marshal(data)
			Expect(err).To(BeNil())
			Expect(string(json)).To(Equal(expectedJson))
		})
	})

	Describe("POST request handler", func() {
		var makeRequest = func(body string, agentID string, db Database) (*httptest.ResponseRecorder, bool) {
			request, _ := http.NewRequest("POST", "/blah", strings.NewReader(body))
			response := httptest.NewRecorder()
			params := httprouter.Params{httprouter.Param{Key: "agent_id", Value: agentID}}

			returnValue := postDataPoints(response, request, params, db)

			return response, returnValue
		}

		var db *MockDatabase
		existentAgentID := "10"

		BeforeEach(func() {
			db = &MockDatabase{}
			db.CheckAgentIDExistsInfo.AgentIDs = []int{10}
			db.GetVariableIDForNameInfo.Variables = map[string]int{"temperature": 12}
		})

		Describe("when the request is valid", func() {
			var response *httptest.ResponseRecorder
			var returnValue bool

			BeforeEach(func() {
				response, returnValue = makeRequest(`{"time":"2015-05-06T10:15:30Z","data":[{"variable":"temperature","value":10.5}]}`, existentAgentID, db)
			})

			It("returns HTTP 201 response", func() {
				Expect(response.Code).To(Equal(http.StatusCreated))
			})

			It("saves the data point to the database", func() {
				Expect(db.AddDataPointInfo.Calls).To(HaveLen(1))
				dataPoint := db.AddDataPointInfo.Calls[0]
				Expect(dataPoint.AgentID).To(Equal(10))
				Expect(dataPoint.VariableID).To(Equal(12))
				Expect(dataPoint.Value).To(Equal(10.5))
				Expect(dataPoint.Time).To(BeTemporally("==", time.Date(2015, 5, 6, 10, 15, 30, 0, time.UTC)))
			})

			It("returns true to commit the transaction", func() {
				Expect(returnValue).To(BeTrue())
			})
		})

		Describe("when the request is invalid", func() {
			TheRequestFailsWithCode := func(request string, agentID string, responseCode int) {
				var response *httptest.ResponseRecorder
				var returnValue bool

				BeforeEach(func() {
					response, returnValue = makeRequest(request, agentID, db)
				})

				It(fmt.Sprintf("returns HTTP %v response", responseCode), func() {
					Expect(response.Code).To(Equal(responseCode))
				})

				It("does not save the variable to the database", func() {
					Expect(len(db.CreateVariableInfo.Calls)).To(Equal(0))
				})

				It("returns false to rollback the transaction", func() {
					Expect(returnValue).To(BeFalse())
				})
			}

			TheRequestFails := func(request string, agentID string) {
				TheRequestFailsWithCode(request, agentID, http.StatusBadRequest)
			}

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
				TheRequestFails(`{"time":"2015-01-02T03:04:05Z","data":[{"value":10]}`, existentAgentID)
			})

			Describe("because the value field is empty", func() {
				TheRequestFails(`{"time":"2015-01-02T03:04:05Z","data":[{"variable":"temperature","value":""}]}`, existentAgentID)
			})

			Describe("because the value field is not a number", func() {
				TheRequestFails(`{"time":"2015-01-02T03:04:05Z","data":[{"variable":"temperature","value":"abc"}]}`, existentAgentID)
			})

			Describe("because the agent ID is not an integer", func() {
				TheRequestFails(`{"time":"2015-01-02T03:04:05Z","data":[{"variable":"temperature","value":10}]}`, "abc")
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
		var makeRequest = func(query string, agentID string, db Database) (*httptest.ResponseRecorder, bool) {
			request, _ := http.NewRequest("GET", "/blah?"+query, strings.NewReader(""))
			response := httptest.NewRecorder()
			params := httprouter.Params{httprouter.Param{Key: "agent_id", Value: agentID}}

			returnValue := getData(response, request, params, db)

			return response, returnValue
		}

		var db *MockDatabase

		BeforeEach(func() {
			db = &MockDatabase{}
			db.CheckAgentIDExistsInfo.AgentIDs = []int{1}
			db.GetVariableByIDInfo.Variables = map[int]Variable{
				123: Variable{Name: "temperature", Units: "°C"},
				321: Variable{Name: "humidity", Units: "%"},
			}
			db.GetDataInfo.Values = AgentValues{
				1: VariableValues{
					123: ValueSet{
						"2015-03-27T06:00:00Z": 100,
						"2015-03-27T09:00:00Z": 105,
					},
					321: ValueSet{
						"2015-03-27T08:00:00Z": 100.5,
						"2015-03-27T12:00:00Z": 80.9,
					},
				},
			}
		})

		Context("when the request is valid", func() {
			It("returns the data", func() {
				resp, returnValue := makeRequest("variable=123&variable=321&date_from=2015-03-27T05:00:00Z&date_to=2015-03-28T23:50:45Z", "1", db)

				Expect(resp.Code).To(Equal(http.StatusOK))
				Expect(string(resp.Body.Bytes())).To(Equal(`{"data":[` +
					`{"id":123,"name":"temperature","units":"°C","points":{"2015-03-27T06:00:00Z":100,"2015-03-27T09:00:00Z":105}},` +
					`{"id":321,"name":"humidity","units":"%","points":{"2015-03-27T08:00:00Z":100.5,"2015-03-27T12:00:00Z":80.9}}` +
					`]}`))
				Expect(returnValue).To(BeTrue())
				Expect(resp.HeaderMap).To(HaveKeyWithValue("Content-Type", []string{"application/json; charset=utf-8"}))
			})
		})

		Context("when the request is invalid", func() {
			TheRequestFailsWithCode := func(query string, agentID string, responseCode int) {
				var response *httptest.ResponseRecorder
				var returnValue bool

				BeforeEach(func() {
					response, returnValue = makeRequest(query, agentID, db)
				})

				It(fmt.Sprintf("returns HTTP %v response", responseCode), func() {
					Expect(response.Code).To(Equal(responseCode))
				})

				It("does not save the variable to the database", func() {
					Expect(len(db.CreateVariableInfo.Calls)).To(Equal(0))
				})

				It("returns false to rollback the transaction", func() {
					Expect(returnValue).To(BeFalse())
				})
			}

			TheRequestFails := func(query string, agentID string) {
				TheRequestFailsWithCode(query, agentID, http.StatusBadRequest)
			}

			Context("because the agent ID does not exist", func() {
				TheRequestFailsWithCode("variable=123&variable=321&date_from=2015-03-27T05:00:00Z&date_to=2015-03-28T23:50:45Z", "909090", http.StatusNotFound)
			})

			Context("because the agent ID is not an integer", func() {
				TheRequestFails("variable=123&variable=321&date_from=2015-03-27T05:00:00Z&date_to=2015-03-28T23:50:45Z", "abc")
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
