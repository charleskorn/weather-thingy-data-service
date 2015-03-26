package main

import (
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

const (
	TESTING_ADDRESS = ":8081"
)

func urlFor(path string) string {
	return "http://localhost" + TESTING_ADDRESS + path
}

func TestDataService(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("HTTP endpoints", func() {
		g.BeforeEach(func() {
			go startServer(TESTING_ADDRESS)
		})

		g.AfterEach(func() {
			stopServer()
		})

		g.Describe("/", func() {
			g.It("returns the welcome message", func() {
				resp, err := http.Get(urlFor("/"))
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(string(body)).To(Equal("Hello, world!"))
			})
		})
	})
}
