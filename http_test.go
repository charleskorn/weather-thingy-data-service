package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"encoding/base64"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"io"
)

const (
	TestingAddress            = ":8081"
	StatusUnprocessableEntity = 422
)

func urlFor(path string) string {
	return "http://localhost" + TestingAddress + path
}

func haveJSONContentType() types.GomegaMatcher {
	return HaveKeyWithValue("Content-Type", Or(Equal([]string{"application/json; charset=UTF-8"}), Equal([]string{"application/json; charset=utf-8"})))
}

var _ = Describe("HTTP endpoints", func() {
	var testDataSourceName string
	var db Database
	var testUser User
	const testUserPassword string = "TestPassword123"

	BeforeEach(func() {
		testDataSourceName = getTestDataSourceName()
		removeTestDatabase(testDataSourceName, true)

		var err error
		db, err = connectToDatabase(testDataSourceName)
		Expect(err).To(BeNil())
		_, err = db.RunMigrations()
		Expect(err).To(BeNil())

		testUser = User{
			Email:   "validuser@testing.com",
			IsAdmin: false,
			Created: time.Now(),
		}

		testUser.SetPassword(testUserPassword)

		Expect(db.BeginTransaction()).To(Succeed())
		Expect(db.CreateUser(&testUser)).To(Succeed())
		Expect(db.CommitTransaction()).To(Succeed())

		go startServer(Config{ServerAddress: TestingAddress, DataSourceName: testDataSourceName})
	})

	AfterEach(func() {
		db.Close()
		stopServer()
		removeTestDatabase(testDataSourceName, false)
	})

	postWithAuthentication := func(url string, contentType string, body io.Reader) *http.Response {
		request, err := http.NewRequest("POST", url, body)
		Expect(err).To(BeNil())

		request.Header.Set("Content-Type", contentType)
		request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(testUser.Email+":"+testUserPassword)))

		resp, err := http.DefaultClient.Do(request)
		Expect(err).To(BeNil())

		return resp
	}

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
			It("saves the agent to the database and returns the agent ID", func() {
				resp := postWithAuthentication(urlFor("/v1/agents"), "application/json", strings.NewReader(`{"name":"New agent name"}`))

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
				var ownerUserId int
				var created time.Time

				err = db.DB().QueryRow("SELECT name, owner_user_id, created FROM agents WHERE agent_id = $1;", id).Scan(&name, &ownerUserId, &created)
				Expect(err).To(BeNil())
				Expect(name).To(Equal("New agent name"))
				Expect(ownerUserId).To(Equal(testUser.UserID))
				Expect(created).To(BeTemporally("~", time.Now(), 100*time.Millisecond))
			})
		})

		Context("GET", func() {
			It("returns all agents", func() {
				ExpectSucceeded(db.DB().Exec("INSERT INTO users (user_id, email, password_iterations, password_salt, password_hash, is_admin) VALUES ($1, $2, $3, $4, $5, $6)", 3001, "blah@blah.com", 0, []byte{}, []byte{}, false))
				ExpectSucceeded(db.DB().Exec("INSERT INTO agents (agent_id, name, owner_user_id, created) VALUES ($1, $2, $3, $4);", 1, "Test Agent 1", 3001, "2015-03-30 12:00:00+10:00"))
				ExpectSucceeded(db.DB().Exec("INSERT INTO agents (agent_id, name, owner_user_id, created) VALUES ($1, $2, $3, $4);", 2, "Test Agent 2", 3001, "2015-02-17 08:00:00+12:00"))

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
				ExpectSucceeded(db.DB().Exec("INSERT INTO users (user_id, email, password_iterations, password_salt, password_hash, is_admin) VALUES ($1, $2, $3, $4, $5, $6)", 3001, "blah@blah.com", 0, []byte{}, []byte{}, false))
				ExpectSucceeded(db.DB().Exec("INSERT INTO agents (agent_id, name, owner_user_id, created) VALUES ($1, $2, $3, $4)", 1001, "First agent", 3001, "2015-04-05T03:00:00Z"))
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
			BeforeEach(func() {
				ExpectSucceeded(db.DB().Exec("INSERT INTO users (user_id, email, password_iterations, password_salt, password_hash, is_admin) VALUES ($1, $2, $3, $4, $5, $6)", 3001, "blah@blah.com", 0, []byte{}, []byte{}, false))
				ExpectSucceeded(db.DB().Exec("INSERT INTO agents (agent_id, name, owner_user_id) VALUES ($1, $2, $3);", 1004, "Test Agent 1", 3001))
				ExpectSucceeded(db.DB().Exec("INSERT INTO variables (variable_id, name, units, display_decimal_places) VALUES ($1, $2, $3, $4);", 1005, "distance", "metres", 1))
			})

			It("saves the data to the database", func() {
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
				ExpectSucceeded(db.DB().Exec("INSERT INTO users (user_id, email, password_iterations, password_salt, password_hash, is_admin) VALUES ($1, $2, $3, $4, $5, $6)", 3001, "blah@blah.com", 0, []byte{}, []byte{}, false))
				ExpectSucceeded(db.DB().Exec("INSERT INTO agents (agent_id, name, owner_user_id) VALUES ($1, $2, $3);", 1004, "Test Agent 1", 3001))
				ExpectSucceeded(db.DB().Exec("INSERT INTO variables (variable_id, name, units, display_decimal_places) VALUES ($1, $2, $3, $4);", 1005, "distance", "metres", 1))
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

	Describe("/v1/users", func() {
		Context("POST", func() {
			It("saves the user to the database and returns the variable ID", func() {
				resp, err := http.Post(urlFor("/v1/users"), "application/json", strings.NewReader(`{"email":"test@testing.com","password":"test123"}`))

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
				var email string
				var isAdmin bool
				var created time.Time
				row := db.DB().QueryRow("SELECT email, is_admin, created FROM users WHERE user_id = $1;", id)
				err = row.Scan(&email, &isAdmin, &created)
				Expect(err).To(BeNil())
				Expect(email).To(Equal("test@testing.com"))
				Expect(isAdmin).To(Equal(false))
				Expect(created).To(BeTemporally("~", time.Now(), 1000*time.Millisecond))
			})
		})
	})
})
