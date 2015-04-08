package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Agent resource", func() {
	Describe("data structure", func() {
		It("can be serialised to JSON", func() {
			agent := Agent{AgentID: 1039, Name: "Cool agent", Created: time.Date(2015, 3, 26, 14, 35, 0, 0, time.UTC)}

			bytes, err := json.Marshal(agent)
			Expect(err).To(BeNil())
			Expect(string(bytes)).To(MatchJSON(`{"id":1039,"name":"Cool agent","created":"2015-03-26T14:35:00Z"}`))
		})

		It("can be deserialised from JSON", func() {
			jsonString := `{"id":1039,"name":"Cool agent","created":"2015-03-26T14:35:00Z"}`
			var agent Agent
			err := json.Unmarshal([]byte(jsonString), &agent)

			expectedAgent := Agent{AgentID: 1039, Name: "Cool agent", Created: time.Date(2015, 3, 26, 14, 35, 0, 0, time.UTC)}
			Expect(err).To(BeNil())
			Expect(agent).To(Equal(expectedAgent))
		})
	})

	Describe("POST request handler", func() {
		var makeRequest = func(body string, db Database) (*httptest.ResponseRecorder, bool) {
			request, _ := http.NewRequest("POST", "/blah", strings.NewReader(body))
			response := httptest.NewRecorder()

			returnValue := postAgent(response, request, nil, db)

			return response, returnValue
		}

		var db *MockDatabase

		BeforeEach(func() {
			db = &MockDatabase{}
		})

		Describe("when the request is valid", func() {
			var response *httptest.ResponseRecorder
			var responseBody string
			var returnValue bool

			BeforeEach(func() {
				db.CreateAgentInfo.AgentIDToReturn = 1019
				response, returnValue = makeRequest(`{"name":"New agent name"}`, db)
				responseBody = string(response.Body.Bytes())
			})

			It("returns HTTP 201 response", func() {
				Expect(response.Code).To(Equal(http.StatusCreated))
			})

			It("saves the agent to the database", func() {
				Expect(db.CreateAgentInfo.Calls).To(HaveLen(1))
				agent := db.CreateAgentInfo.Calls[0]
				Expect(agent.Name).To(Equal("New agent name"))
				Expect(agent.AgentID).To(Equal(0))
				Expect(agent.Created).ToNot(BeTemporally("==", time.Time{}))
			})

			It("returns the newly created agent's ID", func() {
				Expect(responseBody).To(Equal(`{"id":1019}`))
			})

			It("returns an appropriate Content-Type header", func() {
				Expect(response.HeaderMap).To(HaveKeyWithValue("Content-Type", []string{"application/json; charset=utf-8"}))
			})

			It("returns true to commit the transaction", func() {
				Expect(returnValue).To(BeTrue())
			})
		})

		Describe("when the request is invalid", func() {
			TheRequestFails := func(request string) {
				var response *httptest.ResponseRecorder
				var returnValue bool

				BeforeEach(func() {
					response, returnValue = makeRequest(request, db)
				})

				It("returns HTTP 400 response", func() {
					Expect(response.Code).To(Equal(http.StatusBadRequest))
				})

				It("does not save the agent to the database", func() {
					Expect(len(db.CreateAgentInfo.Calls)).To(Equal(0))
				})

				It("returns false to rollback the transaction", func() {
					Expect(returnValue).To(BeFalse())
				})
			}

			Describe("because there is no name field", func() {
				TheRequestFails(`{}`)
			})

			Describe("because the name field is empty", func() {
				TheRequestFails(`{"name":""}`)
			})
		})
	})

	Describe("GET all request handler", func() {
		var makeRequest = func(db Database) *httptest.ResponseRecorder {
			request, _ := http.NewRequest("GET", "/blah", strings.NewReader(""))
			response := httptest.NewRecorder()

			getAllAgents(response, request, nil, db)

			return response
		}

		var db *MockDatabase

		BeforeEach(func() {
			db = &MockDatabase{}
			db.GetAllAgentsInfo.AgentsToReturn = []Agent{
				Agent{AgentID: 1234, Name: "The name", Created: time.Date(2015, 3, 27, 8, 0, 0, 0, time.UTC)},
			}
		})

		It("returns a list of all agents", func() {
			resp := makeRequest(db)

			Expect(resp.Code).To(Equal(http.StatusOK))
			Expect(resp.HeaderMap).To(HaveKeyWithValue("Content-Type", []string{"application/json; charset=utf-8"}))
			Expect(string(resp.Body.Bytes())).To(Equal(`[{"id":1234,"name":"The name","created":"2015-03-27T08:00:00Z"}]`))
		})
	})
})
