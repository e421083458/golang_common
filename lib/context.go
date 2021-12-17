package lib

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
)

type Trace struct {
	TraceId     string
	SpanId      string
	Caller      string
	SrcMethod   string
	HintCode    int64
	HintContent string
}

type TraceContext struct {
	Trace
	CSpanId string
}

func SetGinTraceContext(c *gin.Context, trace *TraceContext) error {
	if trace == nil || c == nil {
		return errors.New("context is nil")
	}
	c.Set("trace", trace)
	return nil
}

func SetTraceContext(ctx context.Context, trace *TraceContext) context.Context {
	if trace == nil {
		return ctx
	}
	return context.WithValue(ctx, "trace", trace)
}

func GetTraceContext(ctx context.Context) *TraceContext {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		traceIntraceContext, exists := ginCtx.Get("trace")
		if !exists {
			return NewTrace()
		}
		traceContext, ok := traceIntraceContext.(*TraceContext)
		if ok {
			return traceContext
		}
		return NewTrace()
	}

	if contextInterface, ok := ctx.(context.Context); ok {
		traceContext, ok := contextInterface.Value("trace").(*TraceContext)
		if ok {
			return traceContext
		}
		return NewTrace()

	}
	return NewTrace()
}
