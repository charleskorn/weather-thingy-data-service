package main

import (
	"encoding/base64"
	"errors"
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

	BeforeEach(func() {
		mockController = gomock.NewController(GinkgoT())
		render = NewMockRender(mockController)
		db = NewMockDatabase(mockController)

		var err error
		request, err = http.NewRequest("SOMETHING", "/", nil)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		mockController.Finish()
	})

	Context("when no authentication header is provided", func() {
		It("returns HTTP 401 and sets the WWW-Authenticate header", func() {
			responseHeaders := http.Header{}

			render.EXPECT().Text(http.StatusUnauthorized, gomock.Any())
			render.EXPECT().Header().Return(responseHeaders)

			withAuthenticatedUser(render, request, nil, nil)

			Expect(responseHeaders.Get("WWW-Authenticate")).To(Equal(`Basic realm="weather-thingy-data-service"`))
		})
	})

	Context("when an authentication header for a different authentication method is provided", func() {
		It("returns HTTP 401 and sets the WWW-Authenticate header", func() {
			request.Header.Set("Authorization", "SomeOtherAuthMethod something")

			responseHeaders := http.Header{}

			render.EXPECT().Text(http.StatusUnauthorized, gomock.Any())
			render.EXPECT().Header().Return(responseHeaders)

			withAuthenticatedUser(render, request, nil, nil)

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

			withAuthenticatedUser(render, request, db, nil)

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

			withAuthenticatedUser(render, request, db, nil)

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

			withAuthenticatedUser(render, request, db, context)

			userType := reflect.TypeOf(User{})
			userFromContext := context.Get(userType)
			Expect(userFromContext.Interface().(User)).To(Equal(user))
		})
	})
})