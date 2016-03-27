package main

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
	"strings"
)

const authenticationRealm string = "weather-thingy-data-service"
const tokenAuthenticationScheme string = "weather-thingy-agent-token"

func withAuthenticatedUser(render render.Render, req *http.Request, db Database, log *logrus.Entry, c martini.Context) {
	authorizationHeader := req.Header.Get("Authorization")
	prefix := "Basic "

	if !strings.HasPrefix(authorizationHeader, prefix) {
		log.Error("Authentication failed because there was no Authorization header or it was not for HTTP basic authentication.")
		respondWithUserAuthenticationFailed(render, "You must authenticate with a HTTP basic authentication header to access this resource.")
		return
	}

	encoded := strings.TrimPrefix(authorizationHeader, prefix)
	decoded, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		log.WithError(err).Error("Could not decode base64-encoded part of Authorization header.")
		respondWithUserAuthenticationFailed(render, "You must authenticate with a HTTP basic authentication header to access this resource.")
		return
	}

	parts := strings.Split(string(decoded), ":")

	if len(parts) != 2 {
		log.Error("Decoded part of Authorization header is invalid.")
		respondWithUserAuthenticationFailed(render, "You must authenticate with a HTTP basic authentication header to access this resource.")
		return
	}

	email := parts[0]
	password := parts[1]

	user, err := db.GetUserByEmail(email)

	if err != nil || subtle.ConstantTimeCompare(user.ComputePasswordHash(password), user.PasswordHash) != 1 {
		log.Error("Authentication failed because the email address or password do not match any known user.")
		respondWithUserAuthenticationFailed(render, "Email address or password do not match any known user.")
		return
	}

	c.Map(user)
}

func respondWithUserAuthenticationFailed(render render.Render, message string) {
	render.Header().Set("WWW-Authenticate", `Basic realm="`+authenticationRealm+`"`)
	render.Text(http.StatusUnauthorized, message)
}

func requireAdminUser(user User, render render.Render, log *logrus.Entry) {
	if !user.IsAdmin {
		render.Error(http.StatusForbidden)
	}
}

func withAuthenticatedAgent(render render.Render, req *http.Request, params martini.Params, db Database, log *logrus.Entry, c martini.Context) {
	authorizationHeader := req.Header.Get("Authorization")
	prefix := tokenAuthenticationScheme + " "

	if !strings.HasPrefix(authorizationHeader, prefix) {
		respondWithAgentAuthenticationFailed(render, fmt.Sprintf("You must authenticate with a HTTP Authorization: %s header to access this resource.", tokenAuthenticationScheme))
		return
	}

	if err := db.BeginTransaction(); err != nil {
		log.WithError(err).Error("Could not begin database transaction.")
		render.Error(http.StatusInternalServerError)
		return
	}

	defer db.RollbackUncommittedTransaction()

	agentID, err := strconv.Atoi(params["agent_id"])

	if err != nil {
		log.WithError(err).Error("Agent ID is invalid.")
		respondWithAgentAuthenticationFailed(render, "Agent ID or token are invalid or incorrect.")
		return
	}

	exists := false

	if exists, err = db.CheckAgentIDExists(agentID); err != nil {
		log.WithError(err).Error("Could not check if agent exists.")
		render.Error(http.StatusInternalServerError)
		return
	} else if !exists {
		log.Error("Authentication failed because the agent does not exist.")
		respondWithAgentAuthenticationFailed(render, "Agent ID or token are invalid or incorrect.")
		return
	}

	agent, err := db.GetAgentByID(agentID)

	if err != nil {
		log.WithError(err).Error("Retrieving agent failed.")
		render.Error(http.StatusInternalServerError)
		return
	}

	token := strings.TrimPrefix(authorizationHeader, prefix)

	if subtle.ConstantTimeCompare(agent.ComputeTokenHash(token), agent.TokenHash) != 1 {
		log.Error("Authentication failed because the token does not match the agent ID given.")
		respondWithAgentAuthenticationFailed(render, "Agent ID or token are invalid or incorrect.")
		return
	}

	c.Map(agent)
}

func respondWithAgentAuthenticationFailed(render render.Render, message string) {
	render.Header().Set("WWW-Authenticate", tokenAuthenticationScheme)
	render.Text(http.StatusUnauthorized, message)
}
