package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func helloWorld(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello, world!")
}

func main() {
	router := httprouter.New()
	router.GET("/", helloWorld)

	log.Fatal(http.ListenAndServe(":8080", router))
}
