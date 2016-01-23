package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"net/http"
	"time"
)

func Log() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		log.WithFields(log.Fields{
			"method":              req.Method,
			"url":                 req.URL.String(),
			"remoteAddress":       req.RemoteAddr,
			"requestLength":       req.ContentLength,
		}).Info("Request processing started.")

		start := time.Now()

		rw := res.(martini.ResponseWriter)
		c.Next()

		elapsed := time.Since(start).Seconds() * 1000

		log.WithFields(log.Fields{
			"method":              req.Method,
			"url":                 req.URL.String(),
			"remoteAddress":       req.RemoteAddr,
			"requestLength":       req.ContentLength,
			"responseStatus":      rw.Status(),
			"responseSize":        rw.Size(),
			"millisecondsElapsed": elapsed,
		}).Info("Request processing complete.")
	}
}
