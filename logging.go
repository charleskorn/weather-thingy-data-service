package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/twinj/uuid"
	"net/http"
	"time"
)

func Log() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		requestId := uuid.NewV4().String()

		logger := logrus.WithFields(logrus.Fields{
			"method":        req.Method,
			"url":           req.URL.String(),
			"remoteAddress": req.RemoteAddr,
			"requestLength": req.ContentLength,
			"requestId":     requestId,
		})

		logger.Info("Request processing started.")
		start := time.Now()

		rw := res.(martini.ResponseWriter)

		c.Map(logger)
		c.Next()

		elapsed := time.Since(start).Seconds() * 1000

		logger.WithFields(logrus.Fields{
			"responseStatus":      rw.Status(),
			"responseSize":        rw.Size(),
			"millisecondsElapsed": elapsed,
		}).Info("Request processing complete.")
	}
}
