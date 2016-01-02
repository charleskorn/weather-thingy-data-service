package main

import "io"
import "net/http"
import "github.com/julienschmidt/httprouter"

func getPing(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "pong")
}
