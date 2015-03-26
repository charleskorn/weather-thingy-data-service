package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"log"
	"testing"
)

func TestDataService(t *testing.T) {
	log.SetOutput(GinkgoWriter)
	RegisterFailHandler(Fail)
	RunSpecs(t, "weather-thingy Data Service")
}
