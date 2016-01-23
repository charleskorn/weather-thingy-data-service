package main

import (
	"encoding/json"
	"time"

	"net/http"

	"github.com/golang/mock/gomock"
	"github.com/martini-contrib/binding"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

		Describe("validation", func() {
			It("succeeds if all properties are set to valid values", func() {
				errors := TestValidation(Variable{Name: "Distance", Units: "metres (m)", DisplayDecimalPlaces: 2})
				Expect(errors).To(BeEmpty())
			})

			It("fails if name property is not set", func() {
				errors := TestValidation(Variable{Units: "metres (m)", DisplayDecimalPlaces: 2})
				Expect(errors).ToNot(BeEmpty())
				Expect(errors[0].FieldNames).To(ContainElement("name"))
				Expect(errors[0].Classification).To(Equal(binding.RequiredError))
			})

			It("fails if units property is not set", func() {
				errors := TestValidation(Variable{Name: "Distance", DisplayDecimalPlaces: 2})
				Expect(errors).ToNot(BeEmpty())
				Expect(errors[0].FieldNames).To(ContainElement("units"))
				Expect(errors[0].Classification).To(Equal(binding.RequiredError))
			})

			It("fails if display decimal places property is zero", func() {
				errors := TestValidation(Variable{Name: "Distance", Units: "metres (m)", DisplayDecimalPlaces: 0})
				Expect(errors).ToNot(BeEmpty())
				Expect(errors[0].FieldNames).To(ContainElement("displayDecimalPlaces"))
				Expect(errors[0].Classification).To(Equal(binding.RequiredError))
			})

			It("fails if display decimal places property is negative", func() {
				errors := TestValidation(Variable{Name: "Distance", Units: "metres (m)", DisplayDecimalPlaces: -1})
				Expect(errors).ToNot(BeEmpty())
				Expect(errors[0].FieldNames).To(ContainElement("displayDecimalPlaces"))
				Expect(errors[0].Classification).To(Equal("OutOfRangeError"))
			})
		})
	})

	Describe("POST request handler", func() {
		var db *MockDatabase
		var render *MockRender

		BeforeEach(func() {
			db = NewMockDatabase(mockController)
			render = NewMockRender(mockController)
		})

		It("saves the variable to the database and returns HTTP 201 response", func() {
			createVariableCall := db.EXPECT().CreateVariable(gomock.Any()).Do(func(variable *Variable) error {
				Expect(variable.Name).To(Equal("New variable name"))
				Expect(variable.Units).To(Equal("metres (m)"))
				Expect(variable.VariableID).To(Equal(0))
				Expect(variable.DisplayDecimalPlaces).To(Equal(2))
				Expect(variable.Created).ToNot(BeTemporally("==", time.Time{}))

				variable.VariableID = 1019

				return nil
			})

			gomock.InOrder(
				db.EXPECT().BeginTransaction(),
				createVariableCall,
				db.EXPECT().CommitTransaction(),
				render.EXPECT().JSON(http.StatusCreated, map[string]interface{}{"id": 1019}),
				db.EXPECT().RollbackUncommittedTransaction(),
			)

			variable := Variable{Name: "New variable name", Units: "metres (m)", DisplayDecimalPlaces: 2}
			postVariable(render, variable, db, nil)
		})
	})
})
