package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/graceful"
)

type Config struct {
	ServerAddress  string
	DataSourceName string
}

const SHUTDOWN_TIMEOUT = 2 * time.Second

var server *graceful.Server

type RouteWithDatabase func(http.ResponseWriter, *http.Request, httprouter.Params, Database) bool

func startServer(config Config) {
	router := httprouter.New()
	router.GET("/", helloWorld)
	router.POST("/v1/agents", wrapRouteInDatabaseTransaction(postAgent, config))

	server = &graceful.Server{
		Timeout: SHUTDOWN_TIMEOUT,
		Server:  &http.Server{Addr: config.ServerAddress, Handler: router},
	}

	if err := server.ListenAndServe(); err != nil {
		if opErr, ok := err.(*net.OpError); !ok || (ok && opErr.Op != "accept") {
			log.Fatal(err)
		}
	}
}

func wrapRouteInDatabaseTransaction(handler RouteWithDatabase, config Config) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		db, err := connectToDatabase(config.DataSourceName)

		if err != nil {
			log.Print("Could not connect to database: ", err)
			SimpleError(w, http.StatusInternalServerError)
			return
		}

		defer db.Close()

		if err := db.BeginTransaction(); err != nil {
			log.Print("Could not begin transaction: ", err)
			SimpleError(w, http.StatusInternalServerError)
			return
		}

		response := httptest.NewRecorder()
		commitTransaction := handler(response, r, p, db)

		if commitTransaction {
			if err := db.CommitTransaction(); err != nil {
				log.Print("Could not commit transaction: ", err)
				SimpleError(w, http.StatusInternalServerError)
				return
			}
		} else {
			if err := db.RollbackTransaction(); err != nil {
				log.Print("Could not rollback transaction: ", err)
				SimpleError(w, http.StatusInternalServerError)
				return
			}
		}

		for k, v := range response.Header() {
			w.Header()[k] = v
		}

		w.WriteHeader(response.Code)
		w.Write(response.Body.Bytes())
	}
}

func stopServer() {
	server.Stop(SHUTDOWN_TIMEOUT)
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

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting up...")

	runMigrations(config)

	log.Println("Starting server...")
	startServer(config)

	log.Println("Shut down normally.")
}
