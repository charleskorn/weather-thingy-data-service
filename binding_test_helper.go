package main

import (
	"github.com/martini-contrib/binding"
	"reflect"

	. "github.com/onsi/gomega"
	"net/http"
	"strings"
)

func TestValidation(obj interface{}) binding.Errors {
	context := NewTestContext()
	context.Map(context)

	request, _ := http.NewRequest("TEST", "/blah", strings.NewReader(""))
	context.Map(request)

	validatorFunc := binding.Validate(obj)
	_, err := context.Invoke(validatorFunc)
	Expect(err).ShouldNot(HaveOccurred())

	errorsType := reflect.TypeOf((binding.Errors)(nil))
	errors := context.Get(errorsType)

	return errors.Interface().(binding.Errors)
}
