package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func helloWorld(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello, world!")
}
