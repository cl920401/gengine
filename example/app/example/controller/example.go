package controller

import (
	"gengine/example/app/example/router"
	"gengine/gcore"
)

type Example struct {

}

func NewExampleController() *Example {
	return &Example{}
}

func (this *Example) Example(ctx *gcore.Context) string {
	return "test"
}

func (this *Example) Name() string {
	return "Example"
}

func (this *Example) RouterMap() gcore.RoutesInfo {
	return router.ExampleRouters
}