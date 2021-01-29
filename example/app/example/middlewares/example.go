package middlewares

import "gengine/gcore"

type Example struct {

}

func NewOutputError() *Example {
	return &Example{}
}

func (this *Example) OnRequest(ctx *gcore.Context) error {
	return nil
}

func (this *Example) OnResponse(ctx *gcore.Context, result interface{}) (interface{}, error) {
	return result, nil
}