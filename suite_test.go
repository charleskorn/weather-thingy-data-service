package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"testing"
	"time"

	"database/sql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

func TestDataService(t *testing.T) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(GinkgoWriter)
	RegisterFailHandler(Fail)
	RunSpecs(t, "weather-thingy Data Service")
}

type BeParsableAndEqualToMatcher struct {
	CompareTo time.Time
}

func BeParsableAndEqualTo(compareTo time.Time) types.GomegaMatcher {
	return &BeParsableAndEqualToMatcher{
		CompareTo: compareTo,
	}
}

func (matcher *BeParsableAndEqualToMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to be", matcher.CompareTo)
}

func (matcher *BeParsableAndEqualToMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to be %s", matcher.CompareTo)
}

func (matcher *BeParsableAndEqualToMatcher) Match(actual interface{}) (bool, error) {
	s, ok := actual.(string)

	if !ok {
		return false, fmt.Errorf("'%#v' is not a string", actual)
	}

	t, err := time.Parse(time.RFC3339, s)

	if err != nil {
		return false, fmt.Errorf("Could not parse value '%s' as a RFC3339 date/time value", s)
	}

	return t.Equal(matcher.CompareTo), nil
}

func ExpectSucceeded(_ sql.Result, err error) {
	Expect(err).To(BeNil())
}
