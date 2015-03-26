package main

import (
	"io/ioutil"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hello World request handler", func() {
	It("returns the welcome message", func() {
		resp := httptest.NewRecorder()
		helloWorld(resp, nil, nil)
		body, _ := ioutil.ReadAll(resp.Body)
		Expect(string(body)).To(Equal("Hello, world!"))
	})
})
