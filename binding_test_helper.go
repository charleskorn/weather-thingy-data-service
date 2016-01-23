package main

import (
	"github.com/codegangsta/inject"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"reflect"

	. "github.com/onsi/gomega"
	"net/http"
	"strings"
)

type TestContext struct {
	inject.Injector
}

func TestValidation(obj interface{}) binding.Errors {
	var context martini.Context = TestContext{inject.New()}
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

func (TestContext) Next() {
	panic("Can't call Next on a TestContext.")
}

func (TestContext) Written() bool {
	panic("Can't call Next on a TestContext.")
}
