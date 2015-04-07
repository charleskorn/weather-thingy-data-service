package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database", func() {
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

	Describe("getMigrationSource", func() {
		It("returns all of the migrations", func() {
			expectedMigrations, _ := ioutil.ReadDir("db/migrations")
			expectedMigrationFileNames := make([]string, len(expectedMigrations))
			for i, m := range expectedMigrations {
				expectedMigrationFileNames[i] = m.Name()
			}

			migrations, _ := getMigrationSource().FindMigrations()
			migrationNames := make([]string, len(migrations))
			for i, m := range migrations {
				migrationNames[i] = m.Id
			}

			// If this test fails, you probably need to run 'generate.sh'.
			Expect(migrationNames).To(Equal(expectedMigrationFileNames))
		})
	})

	Describe("runMigrations", func() {
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
				_, err := db.DB().Exec("CREATE TABLE temp (name VARCHAR(100));")
				Expect(err).To(BeNil())
				err = db.BeginTransaction()
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
				_, err := db.Transaction().Exec("INSERT INTO temp (name) VALUES ('test');")
				Expect(err).To(BeNil())
				var count int
				err = db.Transaction().QueryRow("SELECT COUNT(*) FROM temp;").Scan(&count)
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
				_, err := db.DB().Exec("CREATE TABLE temp (name VARCHAR(100));")
				Expect(err).To(BeNil())
				err = db.BeginTransaction()
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
				_, err := db.Transaction().Exec("INSERT INTO temp (name) VALUES ('test');")
				Expect(err).To(BeNil())
				var count int
				err = db.Transaction().QueryRow("SELECT COUNT(*) FROM temp;").Scan(&count)
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

	Context("when connected to a database with all migrations applied", func() {
		BeforeEach(func() {
			db.RunMigrations()
		})

		Describe("CreateAgent", func() {
			It("saves new agents to the database", func() {
				created := time.Now().Round(time.Millisecond)
				agent := &Agent{Name: "Test agent", Created: created}

				Expect(db.BeginTransaction()).To(BeNil())
				err := db.CreateAgent(agent)
				Expect(err).To(BeNil())
				Expect(agent.AgentID).ToNot(Equal(0))

				Expect(db.CommitTransaction()).To(BeNil())

				var actualName string
				var actualCreated time.Time
				row := db.DB().QueryRow("SELECT name, created FROM agents WHERE agent_id = $1", agent.AgentID)
				err = row.Scan(&actualName, &actualCreated)

				Expect(err).To(BeNil())
				Expect(actualName).To(Equal("Test agent"))
				Expect(actualCreated).To(BeTemporally("==", created))
			})
		})

		Describe("GetAllAgents", func() {
			It("returns an empty list if there are no agents in the database", func() {
				agents, err := db.GetAllAgents()

				Expect(err).To(BeNil())
				Expect(agents).To(BeEmpty())
			})

			It("gets all agents from the database", func() {
				_, err := db.DB().Exec("INSERT INTO agents (agent_id, name, created) VALUES (1, 'Test Agent 1', '2015-03-30 12:00:00+10:00');")
				Expect(err).To(BeNil())
				_, err = db.DB().Exec("INSERT INTO agents (agent_id, name, created) VALUES (2, 'Test Agent 2', '2015-02-17 08:00:00+12:00');")
				Expect(err).To(BeNil())

				agents, err := db.GetAllAgents()

				Expect(err).To(BeNil())
				Expect(agents).To(HaveLen(2))
				Expect(agents[0].AgentID).To(Equal(1))
				Expect(agents[0].Name).To(Equal("Test Agent 1"))
				Expect(agents[0].Created).To(BeTemporally("==", time.Date(2015, 3, 30, 2, 0, 0, 0, time.UTC)))
				Expect(agents[1].AgentID).To(Equal(2))
				Expect(agents[1].Name).To(Equal("Test Agent 2"))
				Expect(agents[1].Created).To(BeTemporally("==", time.Date(2015, 2, 16, 20, 0, 0, 0, time.UTC)))
			})
		})

		Describe("CreateVariable", func() {
			It("saves new variables to the database", func() {
				created := time.Now().Round(time.Millisecond)
				variable := &Variable{Name: "Test variable", Units: "metres (m)", Created: created}

				Expect(db.BeginTransaction()).To(BeNil())
				err := db.CreateVariable(variable)
				Expect(err).To(BeNil())
				Expect(variable.VariableID).ToNot(Equal(0))

				Expect(db.CommitTransaction()).To(BeNil())

				var actualName, actualUnits string
				var actualCreated time.Time
				row := db.DB().QueryRow("SELECT name, units, created FROM variables WHERE variable_id = $1", variable.VariableID)
				err = row.Scan(&actualName, &actualUnits, &actualCreated)

				Expect(err).To(BeNil())
				Expect(actualName).To(Equal("Test variable"))
				Expect(actualUnits).To(Equal("metres (m)"))
				Expect(actualCreated).To(BeTemporally("==", created))
			})
		})

		Describe("AddDataPoint", func() {
			agentID := 2
			variableID := 6

			BeforeEach(func() {
				_, err := db.DB().Exec("INSERT INTO agents (agent_id, name) VALUES ($1, $2)", agentID, "Agent 2 name")
				Expect(err).To(BeNil())
				_, err = db.DB().Exec("INSERT INTO variables (variable_id, name, units) VALUES ($1, $2, $3)", variableID, "Variable 6 name", "Variable 6 units")
				Expect(err).To(BeNil())
			})

			It("adds the data point to the database", func() {
				dataTime := time.Now().Round(time.Millisecond)
				dataPoint := &DataPoint{AgentID: agentID, VariableID: variableID, Time: dataTime, Value: 100.67}

				Expect(db.BeginTransaction()).To(BeNil())
				err := db.AddDataPoint(dataPoint)
				Expect(err).To(BeNil())

				Expect(db.CommitTransaction()).To(BeNil())

				var actualAgentID, actualVariableID int
				var actualValue float64
				var actualTime time.Time
				row := db.DB().QueryRow("SELECT agent_id, variable_id, time, value FROM data")
				err = row.Scan(&actualAgentID, &actualVariableID, &actualTime, &actualValue)

				Expect(err).To(BeNil())
				Expect(actualAgentID).To(Equal(agentID))
				Expect(actualVariableID).To(Equal(variableID))
				Expect(actualTime).To(BeTemporally("==", dataTime))
				Expect(actualValue).To(Equal(100.67))
			})
		})
	})
})

func getTestDataSourceName() string {
	if envDataSource := os.Getenv("WEATHER_THINGY_TEST_DATA_SOURCE"); envDataSource != "" {
		log.Println("Using data source from WEATHER_THINGY_TEST_DATA_SOURCE environment variable: " + envDataSource)
		return envDataSource
	} else {
		defaultDataSource := "postgres://tests@localhost/weatherthingytest?sslmode=disable"
		log.Println("Using default data source: " + defaultDataSource)
		return defaultDataSource
	}
}

func extractDatabaseName(dataSourceName string) (string, error) {
	if url, err := url.Parse(dataSourceName); err != nil {
		return "", err
	} else {
		name := url.Path

		if strings.HasPrefix(name, "/") {
			name = name[1:]
		}

		return name, nil
	}
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
