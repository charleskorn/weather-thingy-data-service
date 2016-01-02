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
	log.Info("Connecting to database...")
	db, err := connectToDatabase(config.DataSourceName)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Could not connect to database.")
	}

	defer db.Close()

	log.Info("Checking for pending migrations...")

	if n, err := db.RunMigrations(); err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Could not apply migrations to database.")
	} else if n == 0 {
		log.Info("Database is already up to date, no migrations applied.")
	} else {
		log.WithFields(log.Fields{"migrationCount": n}).Info("Applied one or more migrations.", n)
	}
}

func main() {
	config := readOptions()

	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Starting up...")

	runMigrations(config)

	log.WithFields(log.Fields{"serverAddress": config.ServerAddress}).Infof("Starting server on %s...", config.ServerAddress)
	startServer(config)

	log.Info("Shut down normally.")
}
