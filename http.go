package main

import (
	"github.com/Sirupsen/logrus"
	"net"
	"net/http"
	"time"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
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
	m.Use(Log())
	m.Use(gzip.All())
	m.Use(martini.Recovery())
	m.Use(method.Override())
	m.Use(render.Renderer())

	r := martini.NewRouter()
	r.Group("/v1", func(g martini.Router) {
		g.Get("/ping", getPing)

		r.Group("", func(g martini.Router) {
			g.Get("/agents", getAllAgents)
			g.Post("/agents", binding.Bind(Agent{}), postAgent)
			g.Get("/agents/:agent_id", getAgent)
			g.Get("/agents/:agent_id/data", getData)
			g.Post("/agents/:agent_id/data", binding.Bind(PostDataPoints{}), postDataPoints)
			g.Post("/variables", binding.Bind(Variable{}), postVariable)
		}, withDatabaseConnection)
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
			logrus.WithError(err).Error("Error occurred while listening for requests.")
		}
	}
}

func withDatabaseConnection(config Config, context martini.Context, r render.Render, log *logrus.Entry) {
	db, err := connectToDatabase(config.DataSourceName)

	if err != nil {
		log.WithError(err).Error("Could not connect to database.")
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
