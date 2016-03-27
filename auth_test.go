package main

import (
	"encoding/base64"
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"reflect"
)

var _ = Describe("Authentication", func() {
	var mockController *gomock.Controller
	var render *MockRender
	var db *MockDatabase
	var request *http.Request
	var logger *logrus.Entry

	BeforeEach(func() {
		mockController = gomock.NewController(GinkgoT())
		render = NewMockRender(mockController)
		db = NewMockDatabase(mockController)
		logger = logrus.NewEntry(logrus.StandardLogger())

		var err error
		request, err = http.NewRequest("SOMETHING", "/", nil)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		mockController.Finish()
	})

	Describe("withAuthenticatedUser", func() {
		Context("when no authentication header is provided", func() {
			It("returns HTTP 401 and sets the WWW-Authenticate header", func() {
				responseHeaders := http.Header{}

				render.EXPECT().Text(http.StatusUnauthorized, gomock.Any())
				render.EXPECT().Header().Return(responseHeaders)

				withAuthenticatedUser(render, request, nil, logger, nil)

				Expect(responseHeaders.Get("WWW-Authenticate")).To(Equal(`Basic realm="weather-thingy-data-service"`))
			})
		})

		Context("when an authentication header for a different authentication method is provided", func() {
			It("returns HTTP 401 and sets the WWW-Authenticate header", func() {
				request.Header.Set("Authorization", "SomeOtherAuthMethod something")

				responseHeaders := http.Header{}

				render.EXPECT().Text(http.StatusUnauthorized, gomock.Any())
				render.EXPECT().Header().Return(responseHeaders)

				withAuthenticatedUser(render, request, nil, logger, nil)

				Expect(responseHeaders.Get("WWW-Authenticate")).To(Equal(`Basic realm="weather-thingy-data-service"`))
			})
		})

		Context("when an authentication header with a user that does not exist is provided", func() {
			It("returns HTTP 401 and sets the WWW-Authenticate header", func() {
				encodedUsernameAndPassword := base64.StdEncoding.EncodeToString([]byte("user@test.com:password123"))
				request.Header.Set("Authorization", "Basic "+encodedUsernameAndPassword)

				responseHeaders := http.Header{}

				render.EXPECT().Text(http.StatusUnauthorized, gomock.Any())
				render.EXPECT().Header().Return(responseHeaders)
				db.EXPECT().GetUserByEmail("user@test.com").Return(User{}, errors.New("The user doesn't exist"))

				withAuthenticatedUser(render, request, db, logger, nil)

				Expect(responseHeaders.Get("WWW-Authenticate")).To(Equal(`Basic realm="weather-thingy-data-service"`))
			})
		})

		Context("when an authentication header with a username and password that do not match is provided", func() {
			It("returns HTTP 401 and sets the WWW-Authenticate header", func() {
				encodedUsernameAndPassword := base64.StdEncoding.EncodeToString([]byte("user@test.com:password123"))
				request.Header.Set("Authorization", "Basic "+encodedUsernameAndPassword)

				responseHeaders := http.Header{}
				user := User{Email: "user@test.com"}
				user.SetPassword("differentpassword")

				render.EXPECT().Text(http.StatusUnauthorized, gomock.Any())
				render.EXPECT().Header().Return(responseHeaders)
				db.EXPECT().GetUserByEmail("user@test.com").Return(user, nil)

				withAuthenticatedUser(render, request, db, logger, nil)

				Expect(responseHeaders.Get("WWW-Authenticate")).To(Equal(`Basic realm="weather-thingy-data-service"`))
			})
		})

		Context("when an authentication header with a username and password that do match is provided", func() {
			It("does not render a response and sets the user in the request context", func() {
				encodedUsernameAndPassword := base64.StdEncoding.EncodeToString([]byte("user@test.com:password123"))
				request.Header.Set("Authorization", "Basic "+encodedUsernameAndPassword)

				user := User{Email: "user@test.com"}
				user.SetPassword("password123")

				context := NewTestContext()

				db.EXPECT().GetUserByEmail("user@test.com").Return(user, nil)

				withAuthenticatedUser(render, request, db, logger, context)

				userType := reflect.TypeOf(User{})
				userFromContext := context.Get(userType)
				Expect(userFromContext.Interface().(User)).To(Equal(user))
			})
		})
	})

	Context("requireAdminUser", func() {
		Context("when the user is an administrator", func() {
			It("does not render a response", func() {
				user := User{
					IsAdmin: true,
				}

				requireAdminUser(user, render, logger)
			})
		})

		Context("when the user is not an administrator", func() {
			It("returns a HTTP 403 response", func() {
				user := User{
					IsAdmin: false,
				}

				render.EXPECT().Error(http.StatusForbidden)

				requireAdminUser(user, render, logger)
			})
		})
	})

	Context("withAuthenticatedAgent", func() {
		Context("when no authorisation header is provided", func() {
			It("returns HTTP 401 and sets the WWW-Authenticate header", func() {
				params := martini.Params{}
				responseHeaders := http.Header{}

				render.EXPECT().Text(http.StatusUnauthorized, gomock.Any())
				render.EXPECT().Header().Return(responseHeaders)

				withAuthenticatedAgent(render, request, params, db, logger, nil)

				Expect(responseHeaders.Get("WWW-Authenticate")).To(Equal(`weather-thingy-agent-token`))
			})
		})

		Context("when an authentication header for a different authentication method is provided", func() {
			It("returns HTTP 401 and sets the WWW-Authenticate header", func() {
				params := martini.Params{}
				request.Header.Set("Authorization", "SomeOtherAuthMethod something")
				responseHeaders := http.Header{}

				render.EXPECT().Text(http.StatusUnauthorized, gomock.Any())
				render.EXPECT().Header().Return(responseHeaders)

				withAuthenticatedAgent(render, request, params, db, logger, nil)

				Expect(responseHeaders.Get("WWW-Authenticate")).To(Equal(`weather-thingy-agent-token`))
			})
		})

		Context("when the agent ID is not an integer", func() {
			It("returns HTTP 401 and sets the WWW-Authenticate header", func() {
				params := martini.Params{"agent_id": "abc123"}
				request.Header.Set("Authorization", "weather-thingy-agent-token something")
				responseHeaders := http.Header{}

				gomock.InOrder(
					db.EXPECT().BeginTransaction(),
					render.EXPECT().Header().Return(responseHeaders),
					render.EXPECT().Text(http.StatusUnauthorized, gomock.Any()),
					db.EXPECT().RollbackUncommittedTransaction(),
				)

				withAuthenticatedAgent(render, request, params, db, logger, nil)

				Expect(responseHeaders.Get("WWW-Authenticate")).To(Equal(`weather-thingy-agent-token`))
			})
		})

		Context("when the agent does not exist", func() {
			It("returns HTTP 401 and sets the WWW-Authenticate header", func() {
				params := martini.Params{"agent_id": "123"}
				request.Header.Set("Authorization", "weather-thingy-agent-token something")
				responseHeaders := http.Header{}

				gomock.InOrder(
					db.EXPECT().BeginTransaction(),
					db.EXPECT().CheckAgentIDExists(123).Return(false, nil),
					render.EXPECT().Header().Return(responseHeaders),
					render.EXPECT().Text(http.StatusUnauthorized, gomock.Any()),
					db.EXPECT().RollbackUncommittedTransaction(),
				)

				withAuthenticatedAgent(render, request, params, db, logger, nil)

				Expect(responseHeaders.Get("WWW-Authenticate")).To(Equal(`weather-thingy-agent-token`))
			})
		})

		Context("when the token provided does not match the agent's token", func() {
			It("returns HTTP 401 and sets the WWW-Authenticate header", func() {
				params := martini.Params{"agent_id": "123"}
				request.Header.Set("Authorization", "weather-thingy-agent-token something")
				responseHeaders := http.Header{}

				agent := Agent{}
				agent.SetToken("somethingelse")

				gomock.InOrder(
					db.EXPECT().BeginTransaction(),
					db.EXPECT().CheckAgentIDExists(123).Return(true, nil),
					db.EXPECT().GetAgentByID(123).Return(agent, nil),
					render.EXPECT().Header().Return(responseHeaders),
					render.EXPECT().Text(http.StatusUnauthorized, gomock.Any()),
					db.EXPECT().RollbackUncommittedTransaction(),
				)

				withAuthenticatedAgent(render, request, params, db, logger, nil)

				Expect(responseHeaders.Get("WWW-Authenticate")).To(Equal(`weather-thingy-agent-token`))
			})
		})

		Context("when the token provided does match the agent's token", func() {
			It("does not render a response and sets the agent in the request context", func() {
				params := martini.Params{"agent_id": "123"}
				request.Header.Set("Authorization", "weather-thingy-agent-token thetoken")
				context := NewTestContext()

				agent := Agent{}
				agent.SetToken("thetoken")

				gomock.InOrder(
					db.EXPECT().BeginTransaction(),
					db.EXPECT().CheckAgentIDExists(123).Return(true, nil),
					db.EXPECT().GetAgentByID(123).Return(agent, nil),
					db.EXPECT().RollbackUncommittedTransaction(),
				)

				withAuthenticatedAgent(render, request, params, db, logger, context)

				agentType := reflect.TypeOf(Agent{})
				agentFromContext := context.Get(agentType)
				Expect(agentFromContext.Interface().(Agent)).To(Equal(agent))
			})
		})
	})
})
