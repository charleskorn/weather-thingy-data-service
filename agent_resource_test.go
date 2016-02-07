package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-martini/martini"
	"github.com/golang/mock/gomock"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Agent resource", func() {
	var mockController *gomock.Controller

	BeforeEach(func() {
		mockController = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockController.Finish()
	})

	Describe("data structure", func() {
		It("can be serialised to JSON", func() {
			agent := Agent{AgentID: 1039, Name: "Cool agent", OwnerUserID: 2456, Created: time.Date(2015, 3, 26, 14, 35, 0, 0, time.UTC)}

			bytes, err := json.Marshal(agent)
			Expect(err).To(BeNil())
			Expect(string(bytes)).To(MatchJSON(`{"id":1039,"name":"Cool agent","ownerUserId":2456,"created":"2015-03-26T14:35:00Z"}`))
		})

		It("can be deserialised from JSON", func() {
			jsonString := `{"id":1039,"name":"Cool agent","ownerUserId":2456,"created":"2015-03-26T14:35:00Z"}`
			var agent Agent
			err := json.Unmarshal([]byte(jsonString), &agent)

			expectedAgent := Agent{AgentID: 1039, Name: "Cool agent", OwnerUserID: 2456, Created: time.Date(2015, 3, 26, 14, 35, 0, 0, time.UTC)}
			Expect(err).To(BeNil())
			Expect(agent).To(Equal(expectedAgent))
		})

		Describe("validation", func() {
			It("succeeds if all required properties are set", func() {
				errors := TestValidation(Agent{Name: "Test Agent"})
				Expect(errors).To(BeEmpty())
			})

			It("fails if name property is not set", func() {
				errors := TestValidation(Agent{})
				Expect(errors).ToNot(BeEmpty())
				Expect(errors[0].FieldNames).To(ContainElement("name"))
				Expect(errors[0].Classification).To(Equal(binding.RequiredError))
			})
		})
	})

	Describe("POST request handler", func() {
		var render *MockRender
		var db *MockDatabase

		BeforeEach(func() {
			render = NewMockRender(mockController)
			db = NewMockDatabase(mockController)
		})

		It("saves the agent to the database and returns the ID of the newly created agent", func() {
			agentId := 1019
			user := User{UserID: 2349}

			createCall := db.EXPECT().CreateAgent(gomock.Any()).Do(func(agent *Agent) error {
				Expect(agent.Name).To(Equal("New agent name"))
				Expect(agent.AgentID).To(Equal(0))
				Expect(agent.Created).ToNot(BeTemporally("==", time.Time{}))
				Expect(agent.OwnerUserID).To(Equal(user.UserID))

				agent.AgentID = agentId

				return nil
			})

			jsonCall := render.EXPECT().JSON(http.StatusCreated, gomock.Any()).Do(func(status int, value interface{}) {
				Expect(value).To(HaveKeyWithValue("id", agentId))
			})

			gomock.InOrder(
				db.EXPECT().BeginTransaction(),
				createCall,
				db.EXPECT().CommitTransaction(),
				jsonCall,
				db.EXPECT().RollbackUncommittedTransaction(),
			)

			postAgent(render, Agent{Name: "New agent name"}, db, user, nil)
		})
	})

	Describe("GET all request handler", func() {
		var render *MockRender
		var db *MockDatabase

		var agents []Agent = []Agent{
			Agent{AgentID: 1234, Name: "The name", Created: time.Date(2015, 3, 27, 8, 0, 0, 0, time.UTC)},
		}

		BeforeEach(func() {
			render = NewMockRender(mockController)
			db = NewMockDatabase(mockController)
		})

		It("returns a list of all agents", func() {
			db.EXPECT().GetAllAgents().Return(agents, nil)
			render.EXPECT().JSON(http.StatusOK, agents)

			getAllAgents(render, db, nil)
		})
	})

	Describe("GET agent request handler", func() {
		var makeRequest = func(render render.Render, db Database, agentID string) {
			params := martini.Params{
				"agent_id": agentID,
			}

			getAgent(render, params, db, nil)
		}

		var db *MockDatabase
		var render *MockRender

		BeforeEach(func() {
			db = NewMockDatabase(mockController)
			render = NewMockRender(mockController)
		})

		Context("when the request is valid", func() {
			It("returns HTTP 200 response with the details of the agent", func() {
				getAgentCall := db.EXPECT().GetAgentByID(1234).Return(
					Agent{AgentID: 1234, Name: "The name", OwnerUserID: 5678, Created: time.Date(2015, 3, 27, 8, 0, 0, 0, time.UTC)},
					nil)

				getVariablesCall := db.EXPECT().GetVariablesForAgent(1234).Return(
					[]Variable{
						Variable{VariableID: 2001, Name: "distance", Units: "metres", DisplayDecimalPlaces: 1, Created: time.Date(2015, 3, 20, 18, 0, 0, 0, time.UTC)},
					},
					nil)

				jsonCall := render.EXPECT().JSON(http.StatusOK, gomock.Any()).Do(func(status int, value interface{}) {
					bytes, err := json.Marshal(value)
					Expect(err).To(BeNil())

					json := string(bytes)
					Expect(json).To(MatchJSON(`{` +
						`"id":1234,` +
						`"ownerUserId":5678,` +
						`"name":"The name",` +
						`"created":"2015-03-27T08:00:00Z",` +
						`"variables":[{"id":2001,"name":"distance","units":"metres","displayDecimalPlaces":1,"created":"2015-03-20T18:00:00Z"}]` +
						`}`))
				})

				gomock.InOrder(
					db.EXPECT().BeginTransaction(),
					db.EXPECT().CheckAgentIDExists(1234).Return(true, nil),
					getAgentCall,
					getVariablesCall,
					db.EXPECT().CommitTransaction(),
					jsonCall,
					db.EXPECT().RollbackUncommittedTransaction(),
				)

				makeRequest(render, db, "1234")
			})
		})

		Context("when the request is invalid", func() {
			Context("because the agent does not exist", func() {
				It("returns HTTP 404 response", func() {
					gomock.InOrder(
						db.EXPECT().BeginTransaction(),
						db.EXPECT().CheckAgentIDExists(5).Return(false, nil),
						render.EXPECT().Text(http.StatusNotFound, gomock.Any()),
						db.EXPECT().RollbackUncommittedTransaction(),
					)

					makeRequest(render, db, "5")
				})
			})

			Context("because the agent ID is not an integer", func() {
				It("returns HTTP 404 response", func() {
					gomock.InOrder(
						db.EXPECT().BeginTransaction(),
						render.EXPECT().Text(http.StatusNotFound, gomock.Any()),
						db.EXPECT().RollbackUncommittedTransaction(),
					)

					makeRequest(render, db, "abc")
				})
			})
		})
	})
})
