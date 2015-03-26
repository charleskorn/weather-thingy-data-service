package main

import (
	"io/ioutil"
	"os"
	"testing"

	"database/sql"
	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"log"
	"net/url"
	"strings"
)

func TestDatabase(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	testDataSourceName := "postgres://tests@localhost/weatherthingytest?sslmode=disable"

	g.Describe("Database", func() {
		var db *Database

		g.Before(func() {
			if envDataSource := os.Getenv("WEATHER_THINGY_TEST_DATA_SOURCE"); envDataSource != "" {
				testDataSourceName = envDataSource
			}
		})

		g.BeforeEach(func() {
			removeTestDatabase(testDataSourceName, true)

			var err error
			db, err = connectToDatabase(testDataSourceName)

			if err != nil {
				log.Fatal("Cannot connect to test database: ", err)
			}
		})

		g.AfterEach(func() {
			db.DatabaseHandle.Close()
			removeTestDatabase(testDataSourceName, false)
		})

		g.Describe("connectToDatabase", func() {
			g.It("connects to the database", func() {
				db, err := connectToDatabase(testDataSourceName)

				Expect(err).To(BeNil())
				Expect(db).ToNot(BeNil())
				Expect(db.DatabaseHandle).ToNot(BeNil())
			})
		})

		g.Describe("getMigrationSource", func() {
			g.It("returns all of the migrations", func() {
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

		g.Describe("runMigrations", func() {
			g.It("applies all of the migrations", func() {
				migrations, _ := getMigrationSource().FindMigrations()
				expectedMigrationCount := len(migrations)

				err := db.runMigrations()
				Expect(err).To(BeNil())

				var actualMigrationCount int
				err = db.DatabaseHandle.QueryRow("SELECT COUNT(*) FROM gorp_migrations;").Scan(&actualMigrationCount)
				Expect(err).To(BeNil())
				Expect(actualMigrationCount).To(Equal(expectedMigrationCount))
			})
		})
	})
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
		log.Fatal("Could not extract database name from data source string '"+dataSourceName+"': ", err)
	}

	url, _ := url.Parse(dataSourceName)
	url.Path = "/postgres"

	db, err := sql.Open("postgres", url.String())

	if err != nil {
		log.Fatal("Could not connect to 'postgres' database to remove test database: ", err)
	}

	defer db.Close()

	if _, err = db.Exec("DROP DATABASE IF EXISTS " + testDatabaseName + ";"); err != nil {
		log.Fatal("Could not drop existing test database: ", err)
	}

	if recreate {
		if _, err = db.Exec("CREATE DATABASE " + testDatabaseName + ";"); err != nil {
			log.Fatal("Could not create test database: ", err)
		}
	}
}
