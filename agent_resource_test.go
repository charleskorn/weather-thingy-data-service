package main

import (
	"testing"

	"encoding/json"
	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"time"
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
	})
}
