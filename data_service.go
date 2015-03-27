package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"flag"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/graceful"
	"os"
)

type Config struct {
	ServerAddress  string
	DataSourceName string
}

var server *graceful.Server

func startServer(config Config) {
	router := httprouter.New()
	router.GET("/", helloWorld)

	server = &graceful.Server{
		Timeout: 10 * time.Second,
		Server:  &http.Server{Addr: config.ServerAddress, Handler: router},
	}

	if err := server.ListenAndServe(); err != nil {
		if opErr, ok := err.(*net.OpError); !ok || (ok && opErr.Op != "accept") {
			log.Fatal(err)
		}
	}
}

func stopServer() {
	server.Stop(10 * time.Second)
}

func readOptions() Config {
	var args Config

	flagSet := flag.NewFlagSet("weather-thingy-data-service", flag.ExitOnError)
	flagSet.StringVar(&args.ServerAddress, "address", ":8080", "The port (and optional address) the server should listen on.")
	flagSet.StringVar(&args.DataSourceName, "dataSource", "postgres://weatherthingy@localhost/weatherthingy?sslmode=disable", "The data source URL to use.")
	flagSet.Parse(os.Args[1:])

	return args
}

func main() {
	config := readOptions()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting up...")

	log.Println("Connecting to database...")
	db, err := connectToDatabase(config.DataSourceName)

	if err != nil {
		log.Fatal("Could not connect to database: ", err)
	}

	log.Println("Checking for pending migrations...")

	if n, err := db.RunMigrations(); err != nil {
		log.Fatal("Could not apply migrations to database: ", err)
	} else {
		log.Printf("Applied %d migrations.", n)
	}

	log.Println("Starting server...")
	startServer(config)

	log.Println("Shutting down...")
	db.Close()

	log.Println("Shut down normally.")
}
