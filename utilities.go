package main

import "net/http"

func SimpleError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
