package main

import (
	"crypto/subtle"
	"encoding/base64"
	"github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strings"
)

const authenticationRealm string = "weather-thingy-data-service"

func withAuthenticatedUser(render render.Render, req *http.Request, db Database, log *logrus.Entry, c martini.Context) {
	authorizationHeader := req.Header.Get("Authorization")
	prefix := "Basic "

	if !strings.HasPrefix(authorizationHeader, prefix) {
		log.Error("Authentication failed because there was no Authorization header or it was not for HTTP basic authentication.")
		respondWithAuthenticationFailed(render, "You must authenticate with a HTTP basic authentication header to access this resource.")
		return
	}

	encoded := strings.TrimPrefix(authorizationHeader, prefix)
	decoded, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		log.WithError(err).Error("Could not decode base64-encoded part of Authorization header.")
		respondWithAuthenticationFailed(render, "You must authenticate with a HTTP basic authentication header to access this resource.")
		return
	}

	parts := strings.Split(string(decoded), ":")

	if len(parts) != 2 {
		log.Error("Decoded part of Authorization header is invalid.")
		respondWithAuthenticationFailed(render, "You must authenticate with a HTTP basic authentication header to access this resource.")
		return
	}

	email := parts[0]
	password := parts[1]

	user, err := db.GetUserByEmail(email)

	if err != nil || subtle.ConstantTimeCompare(user.ComputePasswordHash(password), user.PasswordHash) != 1 {
		log.Error("Authentication failed because the email address or password do not match any known user.")
		respondWithAuthenticationFailed(render, "Email address or password do not match any known user.")
		return
	}

	c.Map(user)
}

func respondWithAuthenticationFailed(render render.Render, message string) {
	render.Header().Set("WWW-Authenticate", `Basic realm="`+authenticationRealm+`"`)
	render.Text(http.StatusUnauthorized, message)
}
