package router

import (
	"gengine/gcore"
)

var ExampleRouters = gcore.RoutesInfo{
	"Example": {
		Method:  "GET",
		Path:    "/example",
		Note:    "测试接口",
	},
}