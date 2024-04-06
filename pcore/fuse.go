package pcore

import (
	"time"

	"github.com/sony/gobreaker"
)

type FuseEntry struct {
	breaker gobreaker.CircuitBreaker
}

// NewFuseEntry 创建一个熔断器条目。
//
// 参数：
//
//	name: 熔断器的名称。
//	timeout: 熔断器打开后的恢复超时时间。
//	requestNum: 触发熔断的请求数量阈值。
//	failureRatio: 触发熔断的失败率阈值。
//
// 返回值：
//
//	返回一个初始化好的 *FuseEntry 指针。
func NewFuseEntry(name string, timeout time.Duration, requestNum int, failureRatio float64) *FuseEntry {
	// 使用 go-circuit-breaker 创建一个熔断器，并配置名称、超时时间和熔断判断逻辑。
	return &FuseEntry{
		breaker: *gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    name,
			Timeout: timeout,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				// 判断是否达到熔断条件：请求数达到阈值且失败率超过设定值。
				return counts.Requests >= uint32(requestNum) && float64(counts.TotalFailures)/float64(counts.Requests) >= failureRatio
			},
		}),
	}
}

// FuseHandler 熔断器中间件。
func (f *FuseEntry) FuseHandler(ctx *Context) {
	if f.breaker.State() == gobreaker.StateOpen {
		// 如果熔断器处于关闭状态，则直接返回错误。
		ctx.Response.WriteHeader(500)
		ctx.Response.Write([]byte(gobreaker.ErrOpenState.Error()))
		ctx.Abort()
		return
	}
	f.breaker.Execute(func() (interface{}, error) {
		// 执行被保护的代码，并返回结果。
		ctx.Next()
		// 判断　响应状态码是否为大于 500
		if ctx.Response.Code >= 500 {
			ctx.Abort()
			return nil, gobreaker.ErrOpenState
		}
		return nil, nil
	})
}
