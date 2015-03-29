package main

import (
	"log"
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
	interceptor := LoggingResponseWriter{Status: http.StatusOK, RealWriter: writer}
	startTime := time.Now()

	defer func() {
		elapsed := time.Now().Sub(startTime) / time.Millisecond

		log.Printf("%s %s from %s with %d bytes, response code %d with %d bytes in %d ms",
			request.Method, request.URL.String(), request.RemoteAddr, request.ContentLength,
			interceptor.Status, interceptor.ResponseSize, elapsed)
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
