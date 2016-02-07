package main

import (
	"github.com/martini-contrib/binding"
	"reflect"

	. "github.com/onsi/gomega"
	"net/http"
	"strings"
)

func TestValidation(json string, objectType interface{}) binding.Errors {
	context := NewTestContext()
	context.Map(context)

	request, _ := http.NewRequest("TEST", "/blah", strings.NewReader(json))
	request.Header.Set("Content-Type", "application/json")
	context.Map(request)

	validatorFunc := binding.Bind(objectType)
	_, err := context.Invoke(validatorFunc)
	Expect(err).ShouldNot(HaveOccurred())

	errorsType := reflect.TypeOf((binding.Errors)(nil))
	errors := context.Get(errorsType)

	return errors.Interface().(binding.Errors)
}
