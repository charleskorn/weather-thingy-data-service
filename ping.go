package main

import (
	"github.com/martini-contrib/render"
	"net/http"
)

func getPing(r render.Render) {
	r.Text(http.StatusOK, "pong")
}
