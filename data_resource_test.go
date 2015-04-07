package main

import (
	"encoding/json"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Data resource", func() {
	Describe("POST data structure", func() {
		It("can be deserialised from JSON", func() {
			jsonString := `{"time":"2015-03-26T14:35:00Z","data":[` +
				`{"variable":"temperature","value":10.675},` +
				`{"variable":"humidity","value":90}` +
				`]}`
			var postData PostDataPoints
			err := json.Unmarshal([]byte(jsonString), &postData)

			expectedPostData := PostDataPoints{
				Time: time.Date(2015, 3, 26, 14, 35, 0, 0, time.UTC),
				Data: []PostDataPoint{
					PostDataPoint{Variable: "temperature", Value: 10.675},
					PostDataPoint{Variable: "humidity", Value: 90},
				},
			}
			Expect(err).To(BeNil())
			Expect(postData).To(Equal(expectedPostData))
		})
	})
})
