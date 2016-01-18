package main

import (
	"encoding/json"
	"time"

	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/golang/mock/gomock"
	"github.com/martini-contrib/render"
)

var _ = Describe("Variables resource", func() {
	var mockController *gomock.Controller

	BeforeEach(func() {
		mockController = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockController.Finish()
	})

	Describe("data structure", func() {
		It("can be serialised to JSON", func() {
			variable := Variable{
				VariableID:           1039,
				Name:                 "Distance to floor",
				Units:                "metres (m)",
				DisplayDecimalPlaces: 1,
				Created:              time.Date(2015, 3, 26, 14, 35, 0, 0, time.UTC),
			}

			bytes, err := json.Marshal(variable)
			Expect(err).To(BeNil())
			Expect(string(bytes)).To(MatchJSON(`{"id":1039,"name":"Distance to floor","units":"metres (m)","displayDecimalPlaces":1,"created":"2015-03-26T14:35:00Z"}`))
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
		var makeRequest = func(body string, render render.Render, db Database) {
			request, _ := http.NewRequest("POST", "/blah", strings.NewReader(body))

			postVariable(render, request, db)
		}

		var db *MockDatabase
		var render *MockRender

		BeforeEach(func() {
			db = NewMockDatabase(mockController)
			render = NewMockRender(mockController)
		})

		Describe("when the request is valid", func() {
			It("saves the variable to the database and returns HTTP 201 response", func() {
				db.EXPECT().CreateVariable(gomock.Any()).Do(func(variable *Variable) error {
					Expect(variable.Name).To(Equal("New variable name"))
					Expect(variable.Units).To(Equal("metres (m)"))
					Expect(variable.VariableID).To(Equal(0))
					Expect(variable.DisplayDecimalPlaces).To(Equal(2))
					Expect(variable.Created).ToNot(BeTemporally("==", time.Time{}))

					variable.VariableID = 1019

					return nil
				})

				render.EXPECT().JSON(http.StatusCreated, map[string]interface{}{"id": 1019})

				makeRequest(`{"name":"New variable name","units":"metres (m)","displayDecimalPlaces":2}`, render, db)
			})
		})

		Describe("when the request is invalid", func() {
			TheRequestFails := func(request string) {
				BeforeEach(func() {
					makeRequest(request, render, db)
				})

				It("returns HTTP 400 response", func() {
					render.EXPECT().Error(http.StatusBadRequest)
				})
			}

			Describe("because there are no fields", func() {
				TheRequestFails(`{}`)
			})

			Describe("because the name field is empty", func() {
				TheRequestFails(`{"name":"","units":"something","displayDecimalPlaces":1}`)
			})

			Describe("because the name field is missing", func() {
				TheRequestFails(`{"units":"something","displayDecimalPlaces":1}`)
			})

			Describe("because the units field is empty", func() {
				TheRequestFails(`{"name":"something","units":"","displayDecimalPlaces":1}`)
			})

			Describe("because the units field is missing", func() {
				TheRequestFails(`{"name":"something","displayDecimalPlaces":1}`)
			})

			Describe("because the display decimal places field is not an integer", func() {
				TheRequestFails(`{"name":"something","units":"something","displayDecimalPlaces":"abc"}`)
			})

			Describe("because the display decimal places field is negative", func() {
				TheRequestFails(`{"name":"something","units":"something","displayDecimalPlaces":-1}`)
			})
		})
	})
})
