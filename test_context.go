package main

import (
	"github.com/codegangsta/inject"
	"github.com/go-martini/martini"
)

type TestContext struct {
	inject.Injector
}

func (TestContext) Next() {
	panic("Can't call Next on a TestContext.")
}

func (TestContext) Written() bool {
	panic("Can't call Next on a TestContext.")
}

func NewTestContext() martini.Context {
	return TestContext{inject.New()}
}
