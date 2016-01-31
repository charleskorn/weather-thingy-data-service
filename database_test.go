package main

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PostgresDatabase", func() {
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

			// If this test fails, you probably need to run 'make generate'.
			Expect(migrationNames).To(Equal(expectedMigrationFileNames))
		})
	})
})
