package main

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/martini-contrib/binding"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"time"
)

var _ = Describe("User resource", func() {
	var mockController *gomock.Controller

	BeforeEach(func() {
		mockController = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockController.Finish()
	})

	Describe("data structure", func() {
		Describe("SetPassword", func() {
			It("should set the password hash, the salt and the number of hashing iterations", func() {
				user := User{
					PasswordIterations: 0,
					PasswordSalt:       []byte("salty"),
					PasswordHash:       []byte("pass"),
				}

				user.SetPassword("test")

				Expect(user.PasswordIterations).To(Equal(passwordIterations))

				Expect(user.PasswordSalt).ToNot(Equal([]byte("salty")))
				Expect(user.PasswordSalt).To(HaveLen(saltBytes))

				Expect(user.PasswordHash).ToNot(Equal([]byte("pass")))
				Expect(user.PasswordHash).To(HaveLen(sha256.Size))
			})
		})

		Describe("ComputePasswordHash", func() {
			user := User{
				PasswordIterations: 10000,
				PasswordSalt:       []byte("salty"),
			}

			It("should return a result of the expected size", func() {
				Expect(user.ComputePasswordHash("password")).To(HaveLen(sha256.Size))
			})

			It("should return different results for different passwords", func() {
				hash1 := user.ComputePasswordHash("password1")
				hash2 := user.ComputePasswordHash("password2")

				Expect(hash1).ToNot(Equal(hash2))
			})
		})

		It("can be serialised to JSON", func() {
			user := User{
				UserID:             20,
				Email:              "test@example.com",
				PasswordIterations: 1000,
				PasswordSalt:       []byte("salty"),
				PasswordHash:       []byte("pass"),
				IsAdmin:            true,
				Created:            time.Date(2015, 10, 3, 6, 7, 8, 0, time.UTC),
			}

			bytes, err := json.Marshal(user)
			Expect(err).To(BeNil())
			Expect(string(bytes)).To(MatchJSON(`{"id":20,"email":"test@example.com","created":"2015-10-03T06:07:08Z"}`))
		})
	})

	Describe("POST data structure", func() {
		It("can be deserialised from JSON", func() {
			jsonString := `{"email":"test@example.com","password":"password1"}`
			var user PostUser
			err := json.Unmarshal([]byte(jsonString), &user)

			expectedUser := PostUser{
				Email:    "test@example.com",
				Password: "password1",
			}

			Expect(err).To(BeNil())
			Expect(user).To(Equal(expectedUser))
		})

		Describe("validation", func() {
			It("succeeds if all required properties are set", func() {
				errors := TestValidation(PostUser{Email: "test@example.com", Password: "password1"})
				Expect(errors).To(BeEmpty())
			})

			It("fails if email property is not set", func() {
				errors := TestValidation(PostUser{Password: "password1"})
				Expect(errors).ToNot(BeEmpty())
				Expect(errors[0].FieldNames).To(ContainElement("email"))
				Expect(errors[0].Classification).To(Equal(binding.RequiredError))
			})

			It("fails if email property is not a valid email address", func() {
				errors := TestValidation(PostUser{Email: "test", Password: "password1"})
				Expect(errors).ToNot(BeEmpty())
				Expect(errors[0].FieldNames).To(ContainElement("email"))
				Expect(errors[0].Classification).To(Equal("InvalidValue"))
			})

			It("fails if password property is not set", func() {
				errors := TestValidation(PostUser{Email: "test@example.com"})
				Expect(errors).ToNot(BeEmpty())
				Expect(errors[0].FieldNames).To(ContainElement("password"))
				Expect(errors[0].Classification).To(Equal(binding.RequiredError))
			})
		})
	})

	Describe("POST request handler", func() {
		var render *MockRender
		var db *MockDatabase

		BeforeEach(func() {
			render = NewMockRender(mockController)
			db = NewMockDatabase(mockController)
		})

		It("saves the user to the database with a hashed password and returns the ID of the newly created user", func() {
			userId := 1019

			createCall := db.EXPECT().CreateUser(gomock.Any()).Do(func(user *User) error {
				Expect(user.UserID).To(Equal(0))
				Expect(user.Email).To(Equal("test@example.com"))
				Expect(user.PasswordIterations).ToNot(BeZero())
				Expect(len(user.PasswordSalt)).To(BeNumerically(">", 0))
				Expect(len(user.PasswordHash)).To(BeNumerically(">", 0))
				Expect(user.IsAdmin).To(Equal(false))
				Expect(user.Created).ToNot(BeTemporally("==", time.Time{}))

				user.UserID = userId

				return nil
			})

			jsonCall := render.EXPECT().JSON(http.StatusCreated, gomock.Any()).Do(func(status int, value interface{}) {
				Expect(value).To(HaveKeyWithValue("id", userId))
			})

			gomock.InOrder(
				db.EXPECT().BeginTransaction(),
				createCall,
				db.EXPECT().CommitTransaction(),
				jsonCall,
				db.EXPECT().RollbackUncommittedTransaction(),
			)

			postUser(render, PostUser{Email: "test@example.com", Password: "password123"}, db, nil)
		})
	})
})
