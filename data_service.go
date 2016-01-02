package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"os"
)

type Config struct {
	ServerAddress  string
	DataSourceName string
}

func readOptions() Config {
	var args Config

	flagSet := flag.NewFlagSet("weather-thingy-data-service", flag.ExitOnError)
	flagSet.StringVar(&args.ServerAddress, "address", ":8080", "The port (and optional address) the server should listen on.")
	flagSet.StringVar(&args.DataSourceName, "dataSource", "postgres://weatherthingy@localhost/weatherthingy?sslmode=disable", "The data source URL to use.")
	flagSet.Parse(os.Args[1:])

	return args
}

func runMigrations(config Config) {
	log.Println("Connecting to database...")
	db, err := connectToDatabase(config.DataSourceName)

	if err != nil {
		log.Fatal("Could not connect to database: ", err)
	}

	defer db.Close()

	log.Println("Checking for pending migrations...")

	if n, err := db.RunMigrations(); err != nil {
		log.Fatal("Could not apply migrations to database: ", err)
	} else {
		log.Printf("Applied %d migrations.", n)
	}
}

func main() {
	config := readOptions()

	log.SetFormatter(&log.JSONFormatter{})
	log.Println("Starting up...")

	runMigrations(config)

	log.Printf("Starting server on %s...", config.ServerAddress)
	startServer(config)

	log.Println("Shut down normally.")
}
