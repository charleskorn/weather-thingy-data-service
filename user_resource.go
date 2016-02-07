package main

import (
	"crypto/rand"
	"crypto/sha256"
	"github.com/Sirupsen/logrus"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/text/unicode/norm"
	"io"
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

const passwordIterations = 100000
const saltBytes int = 32

func (user *User) SetPassword(password string) error {
	salt := make([]byte, saltBytes)

	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return err
	}

	user.PasswordIterations = passwordIterations
	user.PasswordSalt = salt
	user.PasswordHash = user.ComputePasswordHash(password)

	return nil
}

func (user *User) ComputePasswordHash(password string) []byte {
	passwordBytes := norm.NFC.Bytes([]byte(password))

	return pbkdf2.Key(passwordBytes, user.PasswordSalt, user.PasswordIterations, sha256.Size, sha256.New)
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
