package main

import (
	"crypto/subtle"
	"encoding/base64"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strings"
)

const authenticationRealm string = "weather-thingy-data-service"

func withAuthenticatedUser(render render.Render, req *http.Request, db Database, c martini.Context) {
	authorizationHeader := req.Header.Get("Authorization")
	prefix := "Basic "

	if !strings.HasPrefix(authorizationHeader, prefix) {
		respondWithAuthenticationFailed(render, "You must authenticate with a HTTP basic authentication header to access this resource.")
		return
	}

	encoded := strings.TrimPrefix(authorizationHeader, prefix)
	decoded, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		respondWithAuthenticationFailed(render, "You must authenticate with a HTTP basic authentication header to access this resource.")
		return
	}

	parts := strings.Split(string(decoded), ":")

	if len(parts) != 2 {
		respondWithAuthenticationFailed(render, "You must authenticate with a HTTP basic authentication header to access this resource.")
		return
	}

	email := parts[0]
	password := parts[1]

	user, err := db.GetUserByEmail(email)

	if err != nil || subtle.ConstantTimeCompare(user.ComputePasswordHash(password), user.PasswordHash) != 1 {
		respondWithAuthenticationFailed(render, "Email address or password do not match any known user.")
		return
	}

	c.Map(user)
}

func respondWithAuthenticationFailed(render render.Render, message string) {
	render.Header().Set("WWW-Authenticate", `Basic realm="`+authenticationRealm+`"`)
	render.Text(http.StatusUnauthorized, message)
}
