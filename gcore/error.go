package gcore

import (
	"errors"
	"gengine/common"
	"github.com/gin-gonic/gin"
)

const (
	HTTP_STATUS_CODE = 400
)

var codes = map[int]string{
	HTTP_STATUS_CODE: "HTTP_STATUS_CODE",
}

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				status := 400 //default status==400
				if value, exists := ctx.Get(codes[HTTP_STATUS_CODE]); exists {
					if v, ok := value.(int); ok {
						status = v
					}
				}
				ctx.AbortWithStatusJSON(status, gin.H{"error": e})
			}
		}()
		ctx.Next()
	}
}

func (this *Filters) Throw(err error, context *Context) {
	if errors.Is(err, errors.New("")) {
		context.Set(common.CtxHttpStatusCode, 400)
	}
	panic(err)
}

func Error(err error, msg ...string) {
	if err == nil {
		return
	} else {
		errMsg := err.Error()
		if len(msg) > 0 {
			errMsg = msg[0]
		}
		panic(errMsg)
	}
}
