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

	Describe("/v1/agents", func() {
		Context("POST", func() {
			It("saves the agent to the database and returns the agent ID", func() {
				resp, err := http.Post(urlFor("/v1/agents"), "application/json", strings.NewReader(`{"name":"New agent name"}`))

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
				Expect(resp.Header).To(HaveKeyWithValue("Content-Type", []string{"application/json; charset=utf-8"}))

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

		Context("GET", func() {
			It("returns all agents", func() {
				db, err := connectToDatabase(testDataSourceName)
				Expect(err).To(BeNil())
				defer db.Close()

				_, err = db.DB().Exec("INSERT INTO agents (agent_id, name, created) VALUES (1, 'Test Agent 1', '2015-03-30 12:00:00+10:00');")
				Expect(err).To(BeNil())
				_, err = db.DB().Exec("INSERT INTO agents (agent_id, name, created) VALUES (2, 'Test Agent 2', '2015-02-17 08:00:00+12:00');")
				Expect(err).To(BeNil())

				resp, err := http.Get(urlFor("/v1/agents"))

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Header).To(HaveKeyWithValue("Content-Type", []string{"application/json; charset=utf-8"}))

				responseBytes, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())

				var response []map[string]interface{}
				err = json.Unmarshal(responseBytes, &response)

				Expect(err).To(BeNil())
				Expect(response).To(HaveLen(2))
				Expect(response[0]).To(HaveKeyWithValue("id", float64(1)))
				Expect(response[0]).To(HaveKeyWithValue("name", "Test Agent 1"))
				Expect(response[0]).To(HaveKeyWithValue("created", BeParsableAndEqualTo(time.Date(2015, 3, 30, 2, 0, 0, 0, time.UTC))))
				Expect(response[1]).To(HaveKeyWithValue("id", float64(2)))
				Expect(response[1]).To(HaveKeyWithValue("name", "Test Agent 2"))
				Expect(response[1]).To(HaveKeyWithValue("created", BeParsableAndEqualTo(time.Date(2015, 2, 16, 20, 0, 0, 0, time.UTC))))
			})
		})
	})

	Describe("/v1/agents/:agent_id/data", func() {
		Context("POST", func() {
			It("saves the data to the database", func() {
				db, err := connectToDatabase(testDataSourceName)
				Expect(err).To(BeNil())
				defer db.Close()

				_, err = db.DB().Exec("INSERT INTO agents (agent_id, name) VALUES (1004, 'Test Agent 1');")
				Expect(err).To(BeNil())
				_, err = db.DB().Exec("INSERT INTO variables (variable_id, name, units) VALUES (1005, 'distance', 'metres');")
				Expect(err).To(BeNil())

				resp, err := http.Post(urlFor("/v1/agents/1004/data"), "application/json", strings.NewReader(`{"time":"2015-05-06T10:15:30Z","data":[{"variable":"distance","value":10.5}]}`))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				responseBytes, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(string(responseBytes)).To(Equal(""))

				var agentID, variableID int
				var value float64
				var actualTime time.Time
				err = db.DB().QueryRow("SELECT agent_id, variable_id, value, time FROM data;").Scan(&agentID, &variableID, &value, &actualTime)
				Expect(err).To(BeNil())
				Expect(agentID).To(Equal(1004))
				Expect(variableID).To(Equal(1005))
				Expect(value).To(Equal(10.5))
				Expect(actualTime).To(BeTemporally("==", time.Date(2015, 5, 6, 10, 15, 30, 0, time.UTC)))
			})
		})
	})

	Describe("/v1/variables", func() {
		Context("POST", func() {
			It("saves the variable to the database and returns the variable ID", func() {
				resp, err := http.Post(urlFor("/v1/variables"), "application/json", strings.NewReader(`{"name":"New variable name","units":"seconds (s)"}`))

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
				Expect(resp.Header).To(HaveKeyWithValue("Content-Type", []string{"application/json; charset=utf-8"}))

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
				var name, units string
				var created time.Time
				err = db.DB().QueryRow("SELECT name, units, created FROM variables WHERE variable_id = $1;", id).Scan(&name, &units, &created)
				Expect(err).To(BeNil())
				Expect(name).To(Equal("New variable name"))
				Expect(units).To(Equal("seconds (s)"))
				Expect(created).To(BeTemporally("~", time.Now(), 100*time.Millisecond))
			})
		})
	})
})
