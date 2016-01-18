package main

import (
	log "github.com/Sirupsen/logrus"
	"net"
	"net/http"
	"time"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/method"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/strict"
	"github.com/stretchr/graceful"
)

const ShutdownTimeout = 2 * time.Second

var server *graceful.Server


func startServer(config Config) {
	m := martini.New()
	m.Use(gzip.All())
	m.Use(martini.Recovery())
	m.Use(method.Override())
	m.Use(render.Renderer())

	r := martini.NewRouter()
	r.Group("/v1", func(g martini.Router) {
		g.Get("/ping", getPing)
		g.Get("/agents", withDatabaseConnection, getAllAgents)
		g.Post("/agents", withDatabaseConnection, postAgent)
		g.Get("/agents/:agent_id", withDatabaseConnection, getAgent)
		g.Get("/agents/:agent_id/data", withDatabaseConnection, getData)
		g.Post("/agents/:agent_id/data", withDatabaseConnection, postDataPoints)
		g.Post("/variables", withDatabaseConnection, postVariable)
	})

	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)

	r.NotFound(strict.MethodNotAllowed, strict.NotFound)

	m.Map(config)

	server = &graceful.Server{
		Timeout: ShutdownTimeout,
		Server:  &http.Server{Addr: config.ServerAddress, Handler: m},
	}

	if err := server.ListenAndServe(); err != nil {
		if opErr, ok := err.(*net.OpError); !ok || (ok && opErr.Op != "accept") {
			log.WithFields(log.Fields{"error": err}).Error("Error occurred while listening for requests.")
		}
	}
}

func withDatabaseConnection(config Config, context martini.Context, r render.Render) {
	db, err := connectToDatabase(config.DataSourceName)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Could not connect to database.")
		r.Error(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	context.Map(db)
	context.Next()
}

func stopServer() {
	server.Stop(ShutdownTimeout)
}
