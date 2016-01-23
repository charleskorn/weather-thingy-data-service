package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

const (
	TestingAddress = ":8081"
)

func urlFor(path string) string {
	return "http://localhost" + TestingAddress + path
}

func haveJSONContentType() types.GomegaMatcher {
	return HaveKeyWithValue("Content-Type", Or(Equal([]string{"application/json; charset=UTF-8"}), Equal([]string{"application/json; charset=utf-8"})))
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

		go startServer(Config{ServerAddress: TestingAddress, DataSourceName: testDataSourceName})
	})

	AfterEach(func() {
		stopServer()
		removeTestDatabase(testDataSourceName, false)
	})

	Describe("/v1/ping", func() {
		Context("GET", func() {
			It("responds with 'pong'", func() {
				resp, err := http.Get(urlFor("/v1/ping"))

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Header).To(HaveKeyWithValue("Content-Type", []string{"text/plain; charset=UTF-8"}))

				responseBytes, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())

				response := string(responseBytes)
				Expect(response).To(Equal("pong"))
			})
		})
	})

	Describe("/v1/agents", func() {
		Context("POST", func() {
			var db Database

			BeforeEach(func() {
				var err error
				db, err = connectToDatabase(testDataSourceName)
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				db.Close()
			})

			Context("when the agent is valid", func() {
				It("saves the agent to the database and returns the agent ID", func() {
					resp, err := http.Post(urlFor("/v1/agents"), "application/json", strings.NewReader(`{"name":"New agent name"}`))

					Expect(err).To(BeNil())
					Expect(resp.StatusCode).To(Equal(http.StatusCreated))
					Expect(resp.Header).To(haveJSONContentType())

					responseBytes, err := ioutil.ReadAll(resp.Body)
					Expect(err).To(BeNil())

					var response map[string]interface{}
					err = json.Unmarshal(responseBytes, &response)

					Expect(err).To(BeNil())
					Expect(response).To(HaveKey("id"))

					id := int(response["id"].(float64))
					var name string
					var created time.Time
					err = db.DB().QueryRow("SELECT name, created FROM agents WHERE agent_id = $1;", id).Scan(&name, &created)
					Expect(err).To(BeNil())
					Expect(name).To(Equal("New agent name"))
					Expect(created).To(BeTemporally("~", time.Now(), 100*time.Millisecond))
				})
			})

			Context("when the agent is invalid", func() {
				It("does not save the agent to the database and returns a HTTP 422 response", func() {
					var count int
					err := db.DB().QueryRow("SELECT COUNT(*) FROM agents;").Scan(&count)
					Expect(err).To(BeNil())
					Expect(count).To(Equal(0))

					resp, err := http.Post(urlFor("/v1/agents"), "application/json", strings.NewReader(`{"name":""}`))

					Expect(err).To(BeNil())
					Expect(resp.StatusCode).To(Equal(422))
					Expect(resp.Header).To(haveJSONContentType())

					err = db.DB().QueryRow("SELECT COUNT(*) FROM agents;").Scan(&count)
					Expect(err).To(BeNil())
					Expect(count).To(Equal(0))
				})
			})

		})

		Context("GET", func() {
			It("returns all agents", func() {
				db, err := connectToDatabase(testDataSourceName)
				Expect(err).To(BeNil())
				defer db.Close()

				ExpectSucceeded(db.DB().Exec("INSERT INTO agents (agent_id, name, created) VALUES (1, 'Test Agent 1', '2015-03-30 12:00:00+10:00');"))
				ExpectSucceeded(db.DB().Exec("INSERT INTO agents (agent_id, name, created) VALUES (2, 'Test Agent 2', '2015-02-17 08:00:00+12:00');"))

				resp, err := http.Get(urlFor("/v1/agents"))

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Header).To(haveJSONContentType())

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

	Describe("/v1/agents/:agent_id", func() {
		Context("GET", func() {
			It("returns all details of the agent", func() {
				db, err := connectToDatabase(testDataSourceName)
				Expect(err).To(BeNil())
				defer db.Close()

				ExpectSucceeded(db.DB().Exec("INSERT INTO agents (agent_id, name, created) VALUES ($1, $2, $3)", 1001, "First agent", "2015-04-05T03:00:00Z"))
				ExpectSucceeded(db.DB().Exec("INSERT INTO variables (variable_id, name, units, display_decimal_places, created) VALUES ($1, $2, $3, $4, $5)", 2001, "distance", "metres", 1, "2015-04-07T15:00:00Z"))
				ExpectSucceeded(db.DB().Exec("INSERT INTO variables (variable_id, name, units, display_decimal_places, created) VALUES ($1, $2, $3, $4, $5)", 2002, "humidity", "%", 1, "2015-04-07T15:00:00Z"))
				ExpectSucceeded(db.DB().Exec("INSERT INTO data (agent_id, variable_id, value, time) VALUES ($1, $2, $3, $4)", 1001, 2001, 100, "2015-04-07T15:00:00Z"))

				resp, err := http.Get(urlFor("/v1/agents/1001"))

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Header).To(haveJSONContentType())

				responseBytes, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())

				var response map[string]interface{}
				err = json.Unmarshal(responseBytes, &response)

				Expect(err).To(BeNil())
				Expect(response).To(HaveKeyWithValue("id", float64(1001)))
				Expect(response).To(HaveKeyWithValue("name", "First agent"))
				Expect(response).To(HaveKeyWithValue("created", BeParsableAndEqualTo(time.Date(2015, 4, 5, 3, 0, 0, 0, time.UTC))))
				Expect(response).To(HaveKey("variables"))

				variables := response["variables"].([]interface{})
				Expect(variables).To(HaveLen(1))

				variable := variables[0]

				Expect(variable).To(HaveKeyWithValue("id", float64(2001)))
				Expect(variable).To(HaveKeyWithValue("name", "distance"))
				Expect(variable).To(HaveKeyWithValue("units", "metres"))
				Expect(variable).To(HaveKeyWithValue("displayDecimalPlaces", float64(1)))
				Expect(variable).To(HaveKeyWithValue("created", BeParsableAndEqualTo(time.Date(2015, 4, 7, 15, 0, 0, 0, time.UTC))))
			})
		})
	})

	Describe("/v1/agents/:agent_id/data", func() {
		Context("POST", func() {
			It("saves the data to the database", func() {
				db, err := connectToDatabase(testDataSourceName)
				Expect(err).To(BeNil())
				defer db.Close()

				ExpectSucceeded(db.DB().Exec("INSERT INTO agents (agent_id, name) VALUES (1004, 'Test Agent 1');"))
				ExpectSucceeded(db.DB().Exec("INSERT INTO variables (variable_id, name, units, display_decimal_places) VALUES (1005, 'distance', 'metres', 1);"))

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

		Context("GET", func() {
			It("retrieves the data from the database", func() {
				db, err := connectToDatabase(testDataSourceName)
				Expect(err).To(BeNil())
				defer db.Close()

				ExpectSucceeded(db.DB().Exec("INSERT INTO agents (agent_id, name) VALUES (1004, 'Test Agent 1');"))
				ExpectSucceeded(db.DB().Exec("INSERT INTO variables (variable_id, name, units, display_decimal_places) VALUES (1005, 'distance', 'metres', 1);"))
				ExpectSucceeded(db.DB().Exec("INSERT INTO data (agent_id, variable_id, value, time) VALUES ($1, $2, $3, $4)", 1004, 1005, 103, "2015-04-07T15:00:00Z"))
				ExpectSucceeded(db.DB().Exec("INSERT INTO data (agent_id, variable_id, value, time) VALUES ($1, $2, $3, $4)", 1004, 1005, 104, "2015-04-07T15:01:00Z"))
				ExpectSucceeded(db.DB().Exec("INSERT INTO data (agent_id, variable_id, value, time) VALUES ($1, $2, $3, $4)", 1004, 1005, 105, "2015-04-07T15:02:00Z"))
				ExpectSucceeded(db.DB().Exec("INSERT INTO data (agent_id, variable_id, value, time) VALUES ($1, $2, $3, $4)", 1004, 1005, 106, "2015-04-07T15:03:00Z"))

				resp, err := http.Get(urlFor("/v1/agents/1004/data?variable=1005&date_from=2015-04-07T15:00:30Z&date_to=2015-04-07T15:02:30Z"))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Header).To(haveJSONContentType())

				responseBytes, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())

				responseString := string(responseBytes)
				Expect(responseString).To(MatchJSON(`{"data":[` +
					`{"id":1005,"name":"distance","units":"metres","displayDecimalPlaces":1,"points":{"2015-04-07T15:01:00Z":104,"2015-04-07T15:02:00Z":105}}` +
					`]}`))
			})
		})
	})

	Describe("/v1/variables", func() {
		Context("POST", func() {
			It("saves the variable to the database and returns the variable ID", func() {
				resp, err := http.Post(urlFor("/v1/variables"), "application/json", strings.NewReader(`{"name":"New variable name","units":"seconds (s)","displayDecimalPlaces":2}`))

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
				Expect(resp.Header).To(haveJSONContentType())

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
				var displayDecimalPlaces int
				row := db.DB().QueryRow("SELECT name, units, display_decimal_places, created FROM variables WHERE variable_id = $1;", id)
				err = row.Scan(&name, &units, &displayDecimalPlaces, &created)
				Expect(err).To(BeNil())
				Expect(name).To(Equal("New variable name"))
				Expect(units).To(Equal("seconds (s)"))
				Expect(displayDecimalPlaces).To(Equal(2))
				Expect(created).To(BeTemporally("~", time.Now(), 100*time.Millisecond))
			})
		})
	})
})
