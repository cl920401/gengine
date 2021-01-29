package filter

import (
	"fmt"
	. "gengine/common"
	"gengine/gcore"
	"github.com/hashicorp/go-uuid"
	"time"
)

type Default struct{}

func NewDefault() *Default {
	return &Default{}
}

func setRequestID(ctx *gcore.Context) {
	rid := ctx.GetHeader(HeaderRequestID)
	if rid == "" {
		rid, _ = uuid.GenerateUUID()
	}
	ctx.Set(CtxRequestID, rid)
}

func getRequestID(ctx *gcore.Context) string {
	if rid, exists := ctx.Get(CtxRequestID); exists {
		if str, ok := rid.(string); ok {
			return str
		}
	}
	rid, _ := uuid.GenerateUUID()
	return rid
}

func (this *Default) OnRequest(ctx *gcore.Context) error {
	setRequestID(ctx)
	return nil
}

func (this *Default) OnResponse(ctx *gcore.Context, ret interface{}) (interface{}, error) {
	if e, exists := ctx.Get(CtxError); exists {
		s := fmt.Sprintf("<%v>[%s]request error:%v", time.Now(), getRequestID(ctx), e)
		fmt.Println(s)
		return ret, e.(error)
	}
	return ret, nil
}

