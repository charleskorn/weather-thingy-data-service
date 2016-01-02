package main

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"time"
)

type LoggingRouter struct {
	Handler http.Handler
}

type LoggingResponseWriter struct {
	Status       int
	ResponseSize int
	RealWriter   http.ResponseWriter
}

func (r LoggingRouter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.WithFields(log.Fields{
		"method":        request.Method,
		"url":           request.URL.String(),
		"remoteAddress": request.RemoteAddr,
		"requestLength": request.ContentLength,
	}).Info("Request received.")

	interceptor := LoggingResponseWriter{Status: http.StatusOK, RealWriter: writer}
	startTime := time.Now()

	defer func() {
		elapsed := time.Now().Sub(startTime) / time.Millisecond

		log.WithFields(log.Fields{
			"method":              request.Method,
			"url":                 request.URL.String(),
			"remoteAddress":       request.RemoteAddr,
			"requestLength":       request.ContentLength,
			"responseStatus":      interceptor.Status,
			"responseLength":      interceptor.ResponseSize,
			"millisecondsElapsed": elapsed,
		}).Info("Request complete.")
	}()

	r.Handler.ServeHTTP(&interceptor, request)
}

func (w *LoggingResponseWriter) Header() http.Header {
	return w.RealWriter.Header()
}

func (w *LoggingResponseWriter) Write(b []byte) (int, error) {
	w.ResponseSize += len(b)
	return w.RealWriter.Write(b)
}

func (w *LoggingResponseWriter) WriteHeader(status int) {
	w.Status = status
	w.RealWriter.WriteHeader(status)
}
