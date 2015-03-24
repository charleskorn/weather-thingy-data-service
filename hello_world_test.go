package main

import (
	"net/http/httptest"
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

func TestHelloWorld(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Hello World request handler", func() {
		g.It("returns the welcome message", func() {
			resp := httptest.NewRecorder()
			helloWorld(resp, nil, nil)
			body, _ := ioutil.ReadAll(resp.Body)
			Expect(string(body)).To(Equal("Hello, world!"))
		})
	})
}
