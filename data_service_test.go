package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

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
	var testDataSourceName string

	BeforeEach(func() {
		testDataSourceName = getTestDataSourceName()
		removeTestDatabase(testDataSourceName, true)
		db, err := connectToDatabase(testDataSourceName)
		Expect(err).To(BeNil())
		_, err = db.RunMigrations()
		Expect(err).To(BeNil())
		db.Close()

		go startServer(Config{ServerAddress: TESTING_ADDRESS, DataSourceName: testDataSourceName})
	})

	AfterEach(func() {
		stopServer()
		removeTestDatabase(testDataSourceName, false)
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

	Describe("/v1/agents", func() {
		Context("POST", func() {
			It("saves the agent to the database and returns the agent ID", func() {
				resp, err := http.Post(urlFor("/v1/agents"), "application/json", strings.NewReader(`{"name":"New agent name"}`))

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
				Expect(resp.Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))

				responseBytes, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())

				var response map[string]interface{}
				err = json.Unmarshal(responseBytes, &response)

				Expect(err).To(BeNil())
				Expect(response).To(HaveKey("id"))

				db, err := connectToDatabase(testDataSourceName)
				Expect(err).To(BeNil())
				defer db.Close()

				id := int(response["id"].(float64))
				var name string
				var created time.Time
				err = db.DB().QueryRow("SELECT name, created FROM agents WHERE agent_id = $1;", id).Scan(&name, &created)
				Expect(err).To(BeNil())
				Expect(name).To(Equal("New agent name"))
				Expect(created).To(BeTemporally("~", time.Now(), 100*time.Millisecond))
			})
		})
	})
})
