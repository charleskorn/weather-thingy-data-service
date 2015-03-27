package main

import (
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	TESTING_ADDRESS = ":8081"
)

func urlFor(path string) string {
	return "http://localhost" + TESTING_ADDRESS + path
}

var _ = Describe("HTTP endpoints", func() {
	BeforeEach(func() {
		go startServer(Config{ServerAddress: TESTING_ADDRESS})
	})

	AfterEach(func() {
		stopServer()
	})

	Describe("/", func() {
		It("returns the welcome message", func() {
			resp, err := http.Get(urlFor("/"))
			Expect(err).To(BeNil())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("Hello, world!"))
		})
	})
})
