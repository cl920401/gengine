package gcore

import (
	"sync"
)

type Middleware interface {
	OnRequest(*Context) error
	OnResponse(*Context, interface{}) (interface{}, error)
}

func getFilters() *Filters {
	filters_once.Do(func() {
		filters = &Filters{}
	})
	return filters
}

var filters *Filters
var filters_once sync.Once

type Filters struct {
	filters []Middleware
}

func NewFilters() *Filters {
	return &Filters{}
}


func (this *Filters) AddFilters(f ...Middleware) {
	if f != nil && len(f) > 0 {
		this.filters = append(this.filters, f...)
	}
}


func (this *Filters) before(ctx *Context) {
	for _, f := range this.filters {
		err := f.OnRequest(ctx)
		if err != nil {
			this.Throw(err, ctx)
		}
	}
}
func (this *Filters) after(ctx *Context, ret interface{}) interface{} {
	var result = ret
	for _, f := range this.filters {
		r, err := f.OnResponse(ctx, result)
		if err != nil {
			this.Throw(err, ctx)
		}
		result = r
	}
	return result
}

func (this *Filters) Filter(responder Responder, ctx *Context) interface{} {
	this.before(ctx)
	var ret interface{}
	innerNode := getInnerRouter().getRoute(ctx.Request.Method, ctx.Request.URL.Path)
	var innerFairingHandler *Filters
	if innerNode.fullPath != "" && innerNode.handlers != nil { //create inner fairinghandler for route-level middlerware.  hook like
		if fs, ok := innerNode.handlers.([]Middleware); ok {
			innerFairingHandler = NewFilters()
			innerFairingHandler.AddFilters(fs...)
		}
	}
	// exec route-level middleware
	if innerFairingHandler != nil {
		innerFairingHandler.before(ctx)
	}
	if s1, ok := responder.(StringResponder); ok {
		ret = s1(ctx)
	}
	if s2, ok := responder.(JsonResponder); ok {
		ret = s2(ctx)
	}
	if s3, ok := responder.(SqlResponder); ok {
		ret = s3(ctx)
	}
	if s4, ok := responder.(SqlQueryResponder); ok {
		ret = s4(ctx)
	}
	// exec route-level middleware
	if innerFairingHandler != nil {
		ret = innerFairingHandler.after(ctx, ret)
	}
	return getFilters().after(ctx, ret)
}
