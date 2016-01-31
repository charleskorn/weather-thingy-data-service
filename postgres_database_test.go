package main

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	"net/url"
	"os"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PostgresDatabase", func() {
	var testDataSourceName string
	var db Database

	BeforeEach(func() {
		testDataSourceName = getTestDataSourceName()
		removeTestDatabase(testDataSourceName, true)

		var err error
		db, err = connectToDatabase(testDataSourceName)

		if err != nil {
			Fail("Cannot connect to test database: " + err.Error())
		}
	})

	AfterEach(func() {
		db.Close()
		removeTestDatabase(testDataSourceName, false)
	})

	Describe("connectToDatabase", func() {
		It("connects to the database", func() {
			testConnection, err := connectToDatabase(testDataSourceName)

			Expect(err).To(BeNil())
			defer testConnection.Close()

			Expect(testConnection).ToNot(BeNil())
		})
	})

	Describe("RunMigrations", func() {
		It("applies all of the migrations", func() {
			migrations, _ := getMigrationSource().FindMigrations()
			expectedMigrationCount := len(migrations)

			_, err := db.RunMigrations()
			Expect(err).To(BeNil())

			var actualMigrationCount int
			err = db.DB().QueryRow("SELECT COUNT(*) FROM gorp_migrations;").Scan(&actualMigrationCount)
			Expect(err).To(BeNil())
			Expect(actualMigrationCount).To(Equal(expectedMigrationCount))
		})
	})

	Describe("BeginTransaction", func() {
		Context("when there is no active transaction", func() {
			It("does not return an error", func() {
				Expect(db.BeginTransaction()).To(BeNil())
			})

			It("sets Transaction", func() {
				db.BeginTransaction()
				Expect(db.Transaction()).ToNot(BeNil())
			})

			AfterEach(func() {
				db.RollbackTransaction()
			})
		})

		Context("when there is already an active transaction", func() {
			It("returns an error", func() {
				db.BeginTransaction()
				Expect(db.BeginTransaction()).ToNot(BeNil())
			})

			AfterEach(func() {
				db.RollbackTransaction()
			})
		})
	})

	Describe("CommitTransaction", func() {
		Context("when there is no active transaction", func() {
			It("returns an error", func() {
				Expect(db.CommitTransaction()).ToNot(BeNil())
			})
		})

		Context("when there is an active transaction", func() {
			BeforeEach(func() {
				ExpectSucceeded(db.DB().Exec("CREATE TABLE temp (name VARCHAR(100));"))
				err := db.BeginTransaction()
				Expect(err).To(BeNil())
			})

			It("does not return an error", func() {
				Expect(db.CommitTransaction()).To(BeNil())
			})

			It("sets Transaction to nil", func() {
				db.CommitTransaction()
				Expect(db.Transaction()).To(BeNil())
			})

			It("applies changes made to the database", func() {
				ExpectSucceeded(db.Transaction().Exec("INSERT INTO temp (name) VALUES ('test');"))
				var count int
				err := db.Transaction().QueryRow("SELECT COUNT(*) FROM temp;").Scan(&count)
				Expect(err).To(BeNil())
				Expect(count).To(Equal(1))

				err = db.CommitTransaction()
				Expect(err).To(BeNil())

				err = db.DB().QueryRow("SELECT COUNT(*) FROM temp;").Scan(&count)
				Expect(err).To(BeNil())
				Expect(count).To(Equal(1))
			})
		})
	})

	Describe("RollbackTransaction", func() {
		Context("when there is no active transaction", func() {
			It("returns an error", func() {
				Expect(db.RollbackTransaction()).ToNot(BeNil())
			})
		})

		Context("when there is an active transaction", func() {
			BeforeEach(func() {
				ExpectSucceeded(db.DB().Exec("CREATE TABLE temp (name VARCHAR(100));"))
				err := db.BeginTransaction()
				Expect(err).To(BeNil())
			})

			It("does not return an error", func() {
				Expect(db.RollbackTransaction()).To(BeNil())
			})

			It("sets Transaction to nil", func() {
				db.RollbackTransaction()
				Expect(db.Transaction()).To(BeNil())
			})

			It("reverts changes made to the database", func() {
				ExpectSucceeded(db.Transaction().Exec("INSERT INTO temp (name) VALUES ('test');"))
				var count int
				err := db.Transaction().QueryRow("SELECT COUNT(*) FROM temp;").Scan(&count)
				Expect(err).To(BeNil())
				Expect(count).To(Equal(1))

				err = db.RollbackTransaction()
				Expect(err).To(BeNil())

				err = db.DB().QueryRow("SELECT COUNT(*) FROM temp;").Scan(&count)
				Expect(err).To(BeNil())
				Expect(count).To(Equal(0))
			})
		})
	})

	Describe("RollbackUncommittedTransaction", func() {
		Context("when there is no active transaction", func() {
			It("does not return an error", func() {
				Expect(db.RollbackUncommittedTransaction()).To(BeNil())
			})
		})

		Context("when there is an active transaction", func() {
			BeforeEach(func() {
				Expect(db.BeginTransaction()).To(BeNil())
			})

			It("does not return an error", func() {
				Expect(db.RollbackUncommittedTransaction()).To(BeNil())
			})

			It("sets Transaction to nil", func() {
				db.RollbackUncommittedTransaction()
				Expect(db.Transaction()).To(BeNil())
			})
		})
	})

	Context("when connected to a database with all migrations applied", func() {
		BeforeEach(func() {
			db.RunMigrations()

			CreateTestData(db)
		})

		Describe("CreateAgent", func() {
			It("saves new agents to the database", func() {
				created := time.Now().Round(time.Millisecond)
				agent := &Agent{Name: "Test agent", OwnerUserID: 3001, Created: created}

				Expect(db.BeginTransaction()).To(BeNil())
				err := db.CreateAgent(agent)
				Expect(err).To(BeNil())
				Expect(agent.AgentID).ToNot(Equal(0))

				Expect(db.CommitTransaction()).To(BeNil())

				var actualName string
				var actualCreated time.Time
				var actualOwnerUserId int
				row := db.DB().QueryRow("SELECT name, owner_user_id, created FROM agents WHERE agent_id = $1", agent.AgentID)
				err = row.Scan(&actualName, &actualOwnerUserId, &actualCreated)

				Expect(err).To(BeNil())
				Expect(actualName).To(Equal("Test agent"))
				Expect(actualCreated).To(BeTemporally("==", created))
				Expect(actualOwnerUserId).To(Equal(3001))
			})
		})

		Describe("GetAllAgents", func() {
			It("returns an empty list if there are no agents in the database", func() {
				ExpectSucceeded(db.DB().Exec("DELETE FROM data;"))
				ExpectSucceeded(db.DB().Exec("DELETE FROM agents;"))
				agents, err := db.GetAllAgents()

				Expect(err).To(BeNil())
				Expect(agents).To(BeEmpty())
			})

			It("gets all agents from the database", func() {
				agents, err := db.GetAllAgents()

				Expect(err).To(BeNil())
				Expect(agents).To(HaveLen(2))
				Expect(agents[0].AgentID).To(Equal(1001))
				Expect(agents[0].Name).To(Equal("First agent"))
				Expect(agents[0].Created).To(BeTemporally("==", time.Date(2015, 3, 30, 2, 0, 0, 0, time.UTC)))
				Expect(agents[1].AgentID).To(Equal(1002))
				Expect(agents[1].Name).To(Equal("Second agent"))
				Expect(agents[1].Created).To(BeTemporally("==", time.Date(2015, 2, 16, 20, 0, 0, 0, time.UTC)))
			})
		})

		Describe("CheckAgentIDExists", func() {
			BeforeEach(func() {
				err := db.BeginTransaction()
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				db.RollbackTransaction()
			})

			It("returns true if the agent exists", func() {
				exists, err := db.CheckAgentIDExists(1002)
				Expect(err).To(BeNil())
				Expect(exists).To(BeTrue())
			})

			It("returns false if the agent does not exist", func() {
				exists, err := db.CheckAgentIDExists(9001)
				Expect(err).To(BeNil())
				Expect(exists).To(BeFalse())
			})
		})

		Describe("CreateVariable", func() {
			It("saves new variables to the database", func() {
				created := time.Now().Round(time.Millisecond)
				variable := &Variable{Name: "Test variable", Units: "metres (m)", DisplayDecimalPlaces: 2, Created: created}

				Expect(db.BeginTransaction()).To(BeNil())
				err := db.CreateVariable(variable)
				Expect(err).To(BeNil())
				Expect(variable.VariableID).ToNot(Equal(0))

				Expect(db.CommitTransaction()).To(BeNil())

				var actualName, actualUnits string
				var actualCreated time.Time
				var actualDisplayDecimalPlaces int
				row := db.DB().QueryRow("SELECT name, units, display_decimal_places, created "+
					"FROM variables WHERE variable_id = $1", variable.VariableID)
				err = row.Scan(&actualName, &actualUnits, &actualDisplayDecimalPlaces, &actualCreated)

				Expect(err).To(BeNil())
				Expect(actualName).To(Equal("Test variable"))
				Expect(actualUnits).To(Equal("metres (m)"))
				Expect(actualDisplayDecimalPlaces).To(Equal(2))
				Expect(actualCreated).To(BeTemporally("==", created))
			})
		})

		Describe("AddDataPoint", func() {
			It("adds the data point to the database", func() {
				dataTime := time.Now().Round(time.Millisecond)
				dataPoint := DataPoint{AgentID: 1002, VariableID: 2002, Time: dataTime, Value: 100.67}

				Expect(db.BeginTransaction()).To(BeNil())
				err := db.AddDataPoint(dataPoint)
				Expect(err).To(BeNil())

				Expect(db.CommitTransaction()).To(BeNil())

				var actualAgentID, actualVariableID int
				var actualValue float64
				var actualTime time.Time
				row := db.DB().QueryRow("SELECT agent_id, variable_id, time, value FROM data WHERE agent_id = 1002 AND variable_id = 2002;")
				err = row.Scan(&actualAgentID, &actualVariableID, &actualTime, &actualValue)

				Expect(err).To(BeNil())
				Expect(actualAgentID).To(Equal(1002))
				Expect(actualVariableID).To(Equal(2002))
				Expect(actualTime).To(BeTemporally("==", dataTime))
				Expect(actualValue).To(Equal(100.67))
			})
		})

		Describe("GetVariableIDForName", func() {
			BeforeEach(func() {
				err := db.BeginTransaction()
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				db.RollbackTransaction()
			})

			It("returns the variable ID if the variable exists", func() {
				id, err := db.GetVariableIDForName("distance")
				Expect(err).To(BeNil())
				Expect(id).To(Equal(2001))
			})

			It("returns -1 if the variable does not exist", func() {
				id, err := db.GetVariableIDForName("temperature")
				Expect(err).ToNot(BeNil())
				Expect(id).To(Equal(-1))
			})
		})

		Describe("GetData", func() {
			BeforeEach(func() {
				err := db.BeginTransaction()
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				db.RollbackTransaction()
			})

			It("returns the data matching the criteria given", func() {
				data, err := db.GetData(1001, 2002, time.Date(2015, 4, 7, 15, 0, 30, 0, time.UTC), time.Date(2015, 4, 7, 15, 2, 30, 0, time.UTC))
				Expect(err).To(BeNil())
				Expect(data).To(HaveLen(2))
				Expect(data).To(HaveKeyWithValue("2015-04-07T15:01:00Z", float64(103)))
				Expect(data).To(HaveKeyWithValue("2015-04-07T15:02:00Z", float64(104)))
			})
		})

		Describe("GetVariableByID", func() {
			BeforeEach(func() {
				err := db.BeginTransaction()
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				db.RollbackTransaction()
			})

			It("returns the variable if it exists", func() {
				variable, err := db.GetVariableByID(2001)
				Expect(err).To(BeNil())
				Expect(variable.VariableID).To(Equal(2001))
				Expect(variable.Name).To(Equal("distance"))
				Expect(variable.Units).To(Equal("metres"))
				Expect(variable.DisplayDecimalPlaces).To(Equal(2))
				Expect(variable.Created).To(BeTemporally("==", time.Date(2015, 4, 7, 15, 0, 0, 0, time.UTC)))
			})

			It("fails if the variable does not exist", func() {
				variable, err := db.GetVariableByID(9002)
				Expect(err).ToNot(BeNil())
				Expect(variable).To(Equal(Variable{}))
			})
		})

		Describe("GetVariablesForAgent", func() {
			BeforeEach(func() {
				err := db.BeginTransaction()
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				db.RollbackTransaction()
			})

			It("returns the details of every variable associated with the agent", func() {
				variables, err := db.GetVariablesForAgent(1002)
				Expect(err).To(BeNil())
				Expect(variables).To(HaveLen(1))
				Expect(variables[0].VariableID).To(Equal(2001))
				Expect(variables[0].Name).To(Equal("distance"))
				Expect(variables[0].Units).To(Equal("metres"))
				Expect(variables[0].DisplayDecimalPlaces).To(Equal(2))
				Expect(variables[0].Created).To(BeTemporally("==", time.Date(2015, 4, 7, 15, 0, 0, 0, time.UTC)))
			})
		})

		Describe("GetAgentByID", func() {
			BeforeEach(func() {
				ExpectSucceeded(db.DB().Exec("INSERT INTO agents (agent_id, name, owner_user_id, created) VALUES (1, 'Test Agent 1', 3001, '2015-03-30 12:00:00+0:00');"))

				err := db.BeginTransaction()
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				db.RollbackTransaction()
			})

			It("returns the agent if it exists", func() {
				agent, err := db.GetAgentByID(1)
				Expect(err).To(BeNil())
				Expect(agent.AgentID).To(Equal(1))
				Expect(agent.Name).To(Equal("Test Agent 1"))
				Expect(agent.OwnerUserID).To(Equal(3001))
				Expect(agent.Created).To(BeTemporally("==", time.Date(2015, 3, 30, 12, 0, 0, 0, time.UTC)))
			})

			It("fails if the variable does not exist", func() {
				agent, err := db.GetAgentByID(2)
				Expect(err).ToNot(BeNil())
				Expect(agent).To(Equal(Agent{}))
			})
		})

		Describe("CreateUser", func() {
			It("saves new users to the database", func() {
				created := time.Now().Round(time.Millisecond)
				user := &User{
					UserID:             0,
					Email:              "test@example.com",
					PasswordIterations: 1000,
					PasswordSalt:       []byte("salty"),
					PasswordHash:       []byte("pass"),
					IsAdmin:            true,
					Created:            created,
				}

				Expect(db.BeginTransaction()).To(BeNil())
				err := db.CreateUser(user)
				Expect(err).To(BeNil())
				Expect(user.UserID).ToNot(Equal(0))

				Expect(db.CommitTransaction()).To(BeNil())

				var actualEmail string
				var actualPasswordIterations int
				var actualPasswordSalt []byte
				var actualPasswordHash []byte
				var actualIsAdmin bool
				var actualCreated time.Time
				row := db.DB().QueryRow("SELECT email, password_iterations, password_salt, password_hash, is_admin, created FROM users WHERE user_id = $1", user.UserID)
				err = row.Scan(&actualEmail, &actualPasswordIterations, &actualPasswordSalt, &actualPasswordHash, &actualIsAdmin, &actualCreated)

				Expect(err).To(BeNil())
				Expect(actualEmail).To(Equal("test@example.com"))
				Expect(actualPasswordIterations).To(Equal(1000))
				Expect(actualPasswordSalt).To(Equal([]byte("salty")))
				Expect(actualPasswordHash).To(Equal([]byte("pass")))
				Expect(actualIsAdmin).To(Equal(true))
				Expect(actualCreated).To(BeTemporally("==", created))
			})
		})
	})
})

func getTestDataSourceName() string {
	envDataSource := os.Getenv("WEATHER_THINGY_TEST_DATA_SOURCE")

	if envDataSource != "" {
		log.Println("Using data source from WEATHER_THINGY_TEST_DATA_SOURCE environment variable: " + envDataSource)
		return envDataSource
	}

	defaultDataSource := "postgres://tests@localhost/weatherthingytest?sslmode=disable"
	log.Println("Using default data source: " + defaultDataSource)
	return defaultDataSource
}

func extractDatabaseName(dataSourceName string) (string, error) {
	url, err := url.Parse(dataSourceName)

	if err != nil {
		return "", err
	}

	name := url.Path

	if strings.HasPrefix(name, "/") {
		name = name[1:]
	}

	return name, nil
}

func removeTestDatabase(dataSourceName string, recreate bool) {
	testDatabaseName, err := extractDatabaseName(dataSourceName)

	if err != nil {
		Fail("Could not extract database name from data source string '" + dataSourceName + "': " + err.Error())
	}

	url, _ := url.Parse(dataSourceName)
	url.Path = "/postgres"

	db, err := sql.Open("postgres", url.String())

	if err != nil {
		Fail("Could not connect to 'postgres' database to remove test database: " + err.Error())
	}

	defer db.Close()

	if _, err = db.Exec("DROP DATABASE IF EXISTS " + testDatabaseName + ";"); err != nil {
		Fail("Could not drop existing test database: " + err.Error())
	}

	if recreate {
		if _, err = db.Exec("CREATE DATABASE " + testDatabaseName + ";"); err != nil {
			Fail("Could not create test database: " + err.Error())
		}
	}
}

func ExpectSucceeded(_ sql.Result, err error) {
	Expect(err).To(BeNil())
}

func CreateTestData(db Database) {
	ExpectSucceeded(db.DB().Exec("INSERT INTO users (user_id, email, password_iterations, password_salt, password_hash, is_admin, created) VALUES ($1, $2, $3, $4, $5, $6, $7)", 3001, "blah@blah.com", 0, []byte{}, []byte{}, false, "2015-03-30 11:58:00+10:00"))
	ExpectSucceeded(db.DB().Exec("INSERT INTO agents (agent_id, name, owner_user_id, created) VALUES ($1, $2, $3, $4)", 1001, "First agent", 3001, "2015-03-30 12:00:00+10:00"))
	ExpectSucceeded(db.DB().Exec("INSERT INTO agents (agent_id, name, owner_user_id, created) VALUES ($1, $2, $3, $4)", 1002, "Second agent", 3001, "2015-02-17 08:00:00+12:00"))
	ExpectSucceeded(db.DB().Exec("INSERT INTO variables (variable_id, name, units, display_decimal_places, created) VALUES ($1, $2, $3, $4, $5)", 2001, "distance", "metres", 2, "2015-04-07T15:00:00Z"))
	ExpectSucceeded(db.DB().Exec("INSERT INTO variables (variable_id, name, units, display_decimal_places, created) VALUES ($1, $2, $3, $4, $5)", 2002, "humidity", "%", 2, "2015-04-07T15:00:00Z"))
	ExpectSucceeded(db.DB().Exec("INSERT INTO data (agent_id, variable_id, value, time) VALUES ($1, $2, $3, $4)", 1001, 2001, 100, "2015-04-07T15:00:00Z"))
	ExpectSucceeded(db.DB().Exec("INSERT INTO data (agent_id, variable_id, value, time) VALUES ($1, $2, $3, $4)", 1001, 2002, 101, "2015-04-07T15:00:00Z"))
	ExpectSucceeded(db.DB().Exec("INSERT INTO data (agent_id, variable_id, value, time) VALUES ($1, $2, $3, $4)", 1001, 2002, 103, "2015-04-07T15:01:00Z"))
	ExpectSucceeded(db.DB().Exec("INSERT INTO data (agent_id, variable_id, value, time) VALUES ($1, $2, $3, $4)", 1001, 2002, 104, "2015-04-07T15:02:00Z"))
	ExpectSucceeded(db.DB().Exec("INSERT INTO data (agent_id, variable_id, value, time) VALUES ($1, $2, $3, $4)", 1001, 2002, 105, "2015-04-07T15:03:00Z"))
	ExpectSucceeded(db.DB().Exec("INSERT INTO data (agent_id, variable_id, value, time) VALUES ($1, $2, $3, $4)", 1002, 2001, 102, "2015-04-07T15:00:00Z"))
}
