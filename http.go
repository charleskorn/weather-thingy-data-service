package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/graceful"
)

const SHUTDOWN_TIMEOUT = 2 * time.Second

var server *graceful.Server

type RouteWithDatabase func(http.ResponseWriter, *http.Request, httprouter.Params, Database)
type RouteWithDatabaseTransaction func(http.ResponseWriter, *http.Request, httprouter.Params, Database) bool

func startServer(config Config) {
	router := httprouter.New()
	router.POST("/v1/agents", withDatabaseTransaction(postAgent, config))
	router.GET("/v1/agents", withDatabase(getAllAgents, config))
	router.POST("/v1/variables", withDatabaseTransaction(postVariable, config))

	server = &graceful.Server{
		Timeout: SHUTDOWN_TIMEOUT,
		Server:  &http.Server{Addr: config.ServerAddress, Handler: LoggingRouter{Handler: router}},
	}

	if err := server.ListenAndServe(); err != nil {
		if opErr, ok := err.(*net.OpError); !ok || (ok && opErr.Op != "accept") {
			log.Fatal(err)
		}
	}
}

func withDatabase(handler RouteWithDatabase, config Config) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		db, err := connectToDatabase(config.DataSourceName)

		if err != nil {
			log.Print("Could not connect to database: ", err)
			SimpleError(w, http.StatusInternalServerError)
			return
		}

		defer db.Close()

		handler(w, r, p, db)
	}
}

func withDatabaseTransaction(handler RouteWithDatabaseTransaction, config Config) httprouter.Handle {
	wrapped := func(w http.ResponseWriter, r *http.Request, p httprouter.Params, db Database) {
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

	return withDatabase(wrapped, config)
}

func stopServer() {
	server.Stop(SHUTDOWN_TIMEOUT)
}

func SimpleError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
