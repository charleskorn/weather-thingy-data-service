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
	testDataSourceName := "postgres://tests@localhost/weatherthingytest?sslmode=disable"

	var db Database

	BeforeEach(func() {
		if envDataSource := os.Getenv("WEATHER_THINGY_TEST_DATA_SOURCE"); envDataSource != "" {
			log.Println("Using data source from WEATHER_THINGY_TEST_DATA_SOURCE: " + envDataSource)
			testDataSourceName = envDataSource
		} else {
			log.Println("Using default data source: " + testDataSourceName)
		}

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
			db, err := connectToDatabase(testDataSourceName)

			Expect(err).To(BeNil())
			Expect(db).ToNot(BeNil())
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

	Context("when connected to a database with all migrations applied", func() {
		BeforeEach(func() {
			db.RunMigrations()
		})

		It("saves new agents to the database", func() {
			created := time.Now().Round(time.Millisecond)
			agent := &Agent{Name: "Test agent", Created: created}

			err := db.CreateAgent(agent)
			Expect(err).To(BeNil())
			Expect(agent.AgentID).ToNot(Equal(0))

			var actualName string
			var actualCreated time.Time
			row := db.DB().QueryRow("SELECT name, created FROM agents WHERE agent_id = $1", agent.AgentID)
			err = row.Scan(&actualName, &actualCreated)

			Expect(err).To(BeNil())
			Expect(actualName).To(Equal("Test agent"))
			Expect(actualCreated).To(BeTemporally("==", created))
		})
	})
})

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
