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

var server *graceful.Server
var dataSourceName string

func startServer() {
	router := httprouter.New()
	router.GET("/", helloWorld)

	server = &graceful.Server{
		Timeout: 10 * time.Second,
		Server:  &http.Server{Addr: ":8080", Handler: router},
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

func readOptions() {
	flagSet := flag.NewFlagSet("weather-thingy-data-service", flag.ExitOnError)
	flagSet.StringVar(&dataSourceName, "dataSource", "postgres://weatherthingy@localhost/weatherthingy?sslmode=disable", "The data source URL to use.")
	flagSet.Parse(os.Args[1:])
}

func main() {
	readOptions()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting up...")

	log.Println("Connecting to database...")
	db, err := connectToDatabase(dataSourceName)

	if err != nil {
		log.Fatal("Could not connect to database: ", err)
	}

	log.Println("Checking for pending migrations...")

	if n, err := db.runMigrations(); err != nil {
		log.Fatal("Could not apply migrations to database: ", err)
	} else {
		log.Printf("Applied %d migrations.", n)
	}

	log.Println("Starting server...")
	startServer()

	log.Println("Shut down normally.")
}
