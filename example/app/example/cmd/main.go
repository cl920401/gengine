package main

import (
	"gengine/gcore"
	"gengine/example/app/example/controller"
	"gengine/filter"
)

func main() {
	gcore.Init().
		Attach(filter.NewDefault()).
		Group("/example", controller.NewExampleController()).
		Launch()
	return
}