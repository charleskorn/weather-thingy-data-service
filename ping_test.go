package main

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ping", func() {
	Describe("GET request handler", func() {
		var makeRequest = func() *httptest.ResponseRecorder {
			request, _ := http.NewRequest("GET", "/blah", nil)
			response := httptest.NewRecorder()

			getPing(response, request, nil)

			return response
		}

		Describe("when the request is valid", func() {
			var response *httptest.ResponseRecorder
			var responseBody string

			BeforeEach(func() {
				response = makeRequest()
				responseBody = string(response.Body.Bytes())
			})

			It("returns HTTP 200 response", func() {
				Expect(response.Code).To(Equal(http.StatusOK))
			})

			It("responds with the value 'pong'", func() {
				Expect(responseBody).To(Equal("pong"))
			})

			It("returns an appropriate Content-Type header", func() {
				Expect(response.HeaderMap).To(HaveKeyWithValue("Content-Type", []string{"text/plain; charset=utf-8"}))
			})
		})
	})
})
