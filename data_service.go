package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/graceful"
)

var server *graceful.Server

func helloWorld(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello, world!")
}

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
	startServer()
}
