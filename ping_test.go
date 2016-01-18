package main

import (
	"net/http"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Ping", func() {
	var mockController *gomock.Controller

	BeforeEach(func() {
		mockController = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockController.Finish()
	})

	Describe("GET request handler", func() {
		var render *MockRender

		BeforeEach(func() {
			render = NewMockRender(mockController)
		})

		It("writes a HTTP 200 response with the value 'pong'", func() {
			render.EXPECT().Text(http.StatusOK, "pong")

			getPing(render)
		})
	})
})
