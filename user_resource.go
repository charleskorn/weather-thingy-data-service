package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"strings"
	"time"
)

type User struct {
	UserID             int       `json:"id"`
	Email              string    `json:"email"`
	PasswordIterations int       `json:"-"`
	PasswordSalt       []byte    `json:"-"`
	PasswordHash       []byte    `json:"-"`
	IsAdmin            bool      `json:"-"`
	Created            time.Time `json:"created"`
}

type PostUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (user *User) SetPassword(password string) error {
	var err error

	if user.PasswordSalt, err = generateHashingSalt(); err != nil {
		return err
	}

	user.PasswordIterations = hashIterations
	user.PasswordHash = user.ComputePasswordHash(password)

	return nil
}

func (user *User) ComputePasswordHash(password string) []byte {
	return computePasswordHash(password, user.PasswordSalt, user.PasswordIterations)
}

func (user PostUser) Validate(errors binding.Errors, _ *http.Request) binding.Errors {
	// TODO: proper email validation
	if user.Email != "" && !strings.Contains(user.Email, "@") {
		errors = append(errors, binding.Error{
			FieldNames:     []string{"email"},
			Classification: "InvalidValue",
			Message:        "Email address is not valid.",
		})
	}

	return errors
}

func postUser(r render.Render, postedUser PostUser, db Database, log *logrus.Entry) {
	newUser := User{
		Email:   postedUser.Email,
		IsAdmin: false,
		Created: time.Now(),
	}

	newUser.SetPassword(postedUser.Password)

	if err := db.BeginTransaction(); err != nil {
		log.WithError(err).Error("Could not begin database transaction.")
		r.Error(http.StatusInternalServerError)
		return
	}

	defer db.RollbackUncommittedTransaction()

	if err := db.CreateUser(&newUser); err != nil {
		log.WithError(err).Error("Could not create new user.")
		r.Error(http.StatusInternalServerError)
		return
	}

	if err := db.CommitTransaction(); err != nil {
		log.WithError(err).Error("Could not commit transaction.")
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(http.StatusCreated, map[string]interface{}{"id": newUser.UserID})
}
