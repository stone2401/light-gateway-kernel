package pcore

import (
	"golang.org/x/time/rate"
)

type Limiter struct {
	rate.Limiter
}

// NewLimiter 创建一个新的Limiter实例。
//
// 参数:
//
//	limit int - 限流器的速率限制，单位为每秒请求的数量。
//
// 返回值:
//
//	*Limiter - 初始化后的Limiter指针。
func NewLimiter(limit int) *Limiter {
	return &Limiter{
		Limiter: *rate.NewLimiter(rate.Limit(limit), limit),
	}
}

func (l *Limiter) ProxyHandler(ctx *Context) {
	if l.Allow() {
		ctx.Next()
		return
	}
	// 修改状态码，返回
	ctx.Response.WriteHeader(429)
	ctx.Response.Write([]byte("Too Many Requests"))
	ctx.Abort()
}
