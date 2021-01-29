package gcore

import (
	"fmt"
	"gengine/common"
	"gengine/injector"
	"github.com/gin-gonic/gin"
	"reflect"
)

type Bean interface {
	Name() string
}

type Core struct {
	engine *gin.Engine
	exprData     map[string]interface{}
	exprMethod     map[string]interface{}
	currentGroup string // temp-var for group string
}

func Init() *Core {
	g := &Core{
		engine: gin.Default(),
	}
	g.engine.Use(ErrorHandler()) //强迫加载的异常处理中间件
	return g
}

func (this *Core) Beans(beans ...Bean) *Core {
	gin.Recovery()
	for _, bean := range beans {
		this.exprData[bean.Name()] = bean
		injector.BeanFactory.Set(bean)
	}
	return this
}

func (this *Core) Handle(httpMethod, relativePath string, handler interface{}) *Core {
	if h := Convert(handler); h != nil {
		getInnerRouter().addRoute(httpMethod, common.GetPath(this.currentGroup, relativePath), h) // for future
		this.group.Handle(httpMethod, relativePath, h)
	}
	return this
}

func (this *Core) Group(group string, controllers ...IController) *Core {
	g := this.engine.Group(group)
	for _, controller := range controllers {
		this.Beans(controller)
		for k, router := range controller.RouterMap() {
			if _, ok := reflect.TypeOf(controller).MethodByName(k); ok {

			}
		}
	}
	return this
}


// 注册中间件
func (this *Core) Attach(fs ...Middleware) *Core {
	for _, f := range fs {
		injector.BeanFactory.Set(f)
	}
	getFilters().AddFilters(fs...)
	return this
}

func (this *Core) Launch() {
	var port int32 = 8080
	if err := this.engine.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic(fmt.Errorf("server run failed : [%w]", err))
	}
}