package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestAgentResource(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Agent resource", func() {
		g.Describe("data structure", func() {
			g.It("can be serialised to JSON", func() {
				agent := Agent{AgentID: 1039, Name: "Cool agent", Created: time.Date(2015, 3, 26, 14, 35, 0, 0, time.UTC)}

				bytes, err := json.Marshal(agent)
				Expect(err).To(BeNil())
				Expect(string(bytes)).To(Equal(`{"id":1039,"name":"Cool agent","created":"2015-03-26T14:35:00Z"}`))
			})

			g.It("can be deserialised from JSON", func() {
				jsonString := `{"id":1039,"name":"Cool agent","created":"2015-03-26T14:35:00Z"}`
				var agent Agent
				err := json.Unmarshal([]byte(jsonString), &agent)

				expectedAgent := Agent{AgentID: 1039, Name: "Cool agent", Created: time.Date(2015, 3, 26, 14, 35, 0, 0, time.UTC)}
				Expect(err).To(BeNil())
				Expect(agent).To(Equal(expectedAgent))
			})
		})

		g.Describe("POST request handler", func() {
			var makeRequest = func(body string, db Database) *httptest.ResponseRecorder {
				request, _ := http.NewRequest("POST", "/blah", strings.NewReader(body))
				response := httptest.NewRecorder()

				postAgent(response, request, nil, db)

				return response
			}

			var db *MockDatabase

			g.BeforeEach(func() {
				db = &MockDatabase{}
			})

			g.Describe("when the request is valid", func() {
				var response *httptest.ResponseRecorder
				var responseBody string

				g.BeforeEach(func() {
					db.CreateAgentInfo.AgentIDToReturn = 1019
					response = makeRequest(`{"name":"New agent name"}`, db)
					responseBody = string(response.Body.Bytes())
				})

				g.It("returns HTTP 201 response", func() {
					Expect(response.Code).To(Equal(http.StatusCreated))
				})

				g.It("saves the agent to the database", func() {
					Expect(db.CreateAgentInfo.Calls).To(Equal([]Agent{Agent{Name: "New agent name"}}))
				})

				g.It("returns the newly created agent's ID", func() {
					Expect(responseBody).To(Equal(`{"id":1019}`))
				})

				g.It("returns an appropriate Content-Type header", func() {
					Expect(response.HeaderMap["Content-Type"]).To(Equal([]string{"application/json"}))
				})
			})

			g.Describe("when the request is invalid", func() {
				g.Describe("because there is no name field", func() {
					var response *httptest.ResponseRecorder

					g.BeforeEach(func() {
						response = makeRequest(`{}`, db)
					})

					g.It("returns HTTP 400 response", func() {
						Expect(response.Code).To(Equal(http.StatusBadRequest))
					})

					g.It("does not save the agent to the database", func() {
						Expect(len(db.CreateAgentInfo.Calls)).To(Equal(0))
					})
				})

				g.Describe("because the name field is empty", func() {
					var response *httptest.ResponseRecorder

					g.BeforeEach(func() {
						response = makeRequest(`{"name":""}`, db)
					})

					g.It("returns HTTP 400 response", func() {
						Expect(response.Code).To(Equal(http.StatusBadRequest))
					})

					g.It("does not save the agent to the database", func() {
						Expect(len(db.CreateAgentInfo.Calls)).To(Equal(0))
					})
				})
			})
		})
	})
}
