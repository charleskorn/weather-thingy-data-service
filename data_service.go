package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/graceful"
)

var server *graceful.Server

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

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting up...")

	log.Println("Connecting to database...")
	db, err := connectToDatabase("postgres://user:pass@localhost/weatherthingy?sslmode=disable")

	if err != nil {
		log.Fatal("Could not connect to database.", err)
	}

	if err := db.runMigrations(); err != nil {
		log.Fatal("Could not apply migrations to database.", err)
	}

	log.Println("Starting server...")
	startServer()

	log.Println("Shut down normally.")
}
