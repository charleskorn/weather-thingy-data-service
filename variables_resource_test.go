package main

import (
	"encoding/json"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strings"
)

var _ = Describe("Variables resource", func() {
	Describe("data structure", func() {
		It("can be serialised to JSON", func() {
			variable := Variable{VariableID: 1039, Name: "Distance to floor", Units: "metres (m)", Created: time.Date(2015, 3, 26, 14, 35, 0, 0, time.UTC)}

			bytes, err := json.Marshal(variable)
			Expect(err).To(BeNil())
			Expect(string(bytes)).To(MatchJSON(`{"id":1039,"name":"Distance to floor","units":"metres (m)","created":"2015-03-26T14:35:00Z"}`))
		})

		It("can be deserialised from JSON", func() {
			jsonString := `{"id":1039,"name":"Distance to floor","units":"metres (m)","created":"2015-03-26T14:35:00Z"}`
			var variable Variable
			err := json.Unmarshal([]byte(jsonString), &variable)

			expectedVariable := Variable{VariableID: 1039, Name: "Distance to floor", Units: "metres (m)", Created: time.Date(2015, 3, 26, 14, 35, 0, 0, time.UTC)}
			Expect(err).To(BeNil())
			Expect(variable).To(Equal(expectedVariable))
		})
	})

	Describe("POST request handler", func() {
		var makeRequest = func(body string, db Database) (*httptest.ResponseRecorder, bool) {
			request, _ := http.NewRequest("POST", "/blah", strings.NewReader(body))
			response := httptest.NewRecorder()

			returnValue := postVariable(response, request, nil, db)

			return response, returnValue
		}

		var db *MockDatabase

		BeforeEach(func() {
			db = &MockDatabase{}
		})

		Describe("when the request is valid", func() {
			var response *httptest.ResponseRecorder
			var responseBody string
			var returnValue bool

			BeforeEach(func() {
				db.CreateVariableInfo.VariableIDToReturn = 1019
				response, returnValue = makeRequest(`{"name":"New variable name","units":"metres (m)"}`, db)
				responseBody = string(response.Body.Bytes())
			})

			It("returns HTTP 201 response", func() {
				Expect(response.Code).To(Equal(http.StatusCreated))
			})

			It("saves the variable to the database", func() {
				Expect(db.CreateVariableInfo.Calls).To(HaveLen(1))
				variable := db.CreateVariableInfo.Calls[0]
				Expect(variable.Name).To(Equal("New variable name"))
				Expect(variable.Units).To(Equal("metres (m)"))
				Expect(variable.VariableID).To(Equal(0))
				Expect(variable.Created).ToNot(BeTemporally("==", time.Time{}))
			})

			It("returns the newly created variable's ID", func() {
				Expect(responseBody).To(Equal(`{"id":1019}`))
			})

			It("returns an appropriate Content-Type header", func() {
				Expect(response.HeaderMap).To(HaveKeyWithValue("Content-Type", []string{"application/json; charset=utf-8"}))
			})

			It("returns true to commit the transaction", func() {
				Expect(returnValue).To(BeTrue())
			})
		})

		Describe("when the request is invalid", func() {
			Describe("because there are no fields", func() {
				var response *httptest.ResponseRecorder
				var returnValue bool

				BeforeEach(func() {
					response, returnValue = makeRequest(`{}`, db)
				})

				It("returns HTTP 400 response", func() {
					Expect(response.Code).To(Equal(http.StatusBadRequest))
				})

				It("does not save the variable to the database", func() {
					Expect(len(db.CreateVariableInfo.Calls)).To(Equal(0))
				})

				It("returns false to rollback the transaction", func() {
					Expect(returnValue).To(BeFalse())
				})
			})

			Describe("because the name field is empty", func() {
				var response *httptest.ResponseRecorder
				var returnValue bool

				BeforeEach(func() {
					response, returnValue = makeRequest(`{"name":"","units":"something"}`, db)
				})

				It("returns HTTP 400 response", func() {
					Expect(response.Code).To(Equal(http.StatusBadRequest))
				})

				It("does not save the variable to the database", func() {
					Expect(len(db.CreateVariableInfo.Calls)).To(Equal(0))
				})

				It("returns false to rollback the transaction", func() {
					Expect(returnValue).To(BeFalse())
				})
			})

			Describe("because the name field is missing", func() {
				var response *httptest.ResponseRecorder
				var returnValue bool

				BeforeEach(func() {
					response, returnValue = makeRequest(`{"units":"something"}`, db)
				})

				It("returns HTTP 400 response", func() {
					Expect(response.Code).To(Equal(http.StatusBadRequest))
				})

				It("does not save the variable to the database", func() {
					Expect(len(db.CreateVariableInfo.Calls)).To(Equal(0))
				})

				It("returns false to rollback the transaction", func() {
					Expect(returnValue).To(BeFalse())
				})
			})

			Describe("because the units field is empty", func() {
				var response *httptest.ResponseRecorder
				var returnValue bool

				BeforeEach(func() {
					response, returnValue = makeRequest(`{"name":"something","units":""}`, db)
				})

				It("returns HTTP 400 response", func() {
					Expect(response.Code).To(Equal(http.StatusBadRequest))
				})

				It("does not save the variable to the database", func() {
					Expect(len(db.CreateVariableInfo.Calls)).To(Equal(0))
				})

				It("returns false to rollback the transaction", func() {
					Expect(returnValue).To(BeFalse())
				})
			})

			Describe("because the units field is missing", func() {
				var response *httptest.ResponseRecorder
				var returnValue bool

				BeforeEach(func() {
					response, returnValue = makeRequest(`{"name":"something"}`, db)
				})

				It("returns HTTP 400 response", func() {
					Expect(response.Code).To(Equal(http.StatusBadRequest))
				})

				It("does not save the variable to the database", func() {
					Expect(len(db.CreateVariableInfo.Calls)).To(Equal(0))
				})

				It("returns false to rollback the transaction", func() {
					Expect(returnValue).To(BeFalse())
				})
			})
		})
	})
})
