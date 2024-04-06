package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
)

func main() {
	r := gin.Default()

	// // 创建一个熔断器
	// cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
	// 	Name:    "example",
	// 	Timeout: 5 * time.Second,
	// 	ReadyToTrip: func(counts gobreaker.Counts) bool {
	// 		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
	// 		return counts.Requests >= 3 && failureRatio >= 0.6
	// 	},
	// 	OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
	// 		fmt.Printf("Circuit breaker state changed to: %s\n", to)
	// 	},
	// })

	// // 使用熔断器的中间件
	// r.Use(func(c *gin.Context) {
	// 	state := cb.State()
	// 	if state == gobreaker.StateOpen {
	// 		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "circuit breaker is open"})
	// 		c.Abort()
	// 		return
	// 	}

	// 	cb.Execute(func() (interface{}, error) {
	// 		c.Next() // 继续处理请求
	// 		if c.Writer.Status() >= http.StatusInternalServerError {
	// 			return nil, fmt.Errorf("Request failed with status code %d", c.Writer.Status())
	// 		}
	// 		return c.Writer.Status(), nil
	// 	})
	// })
	// 创建熔断器
	// cb := pcore.NewFuseEntry("example", 5*time.Second, 3, 0.6)
	// r.Use(cb.FuseHandler)
	num := atomic.Int64{}
	num.Store(0)
	// 示例路由
	r.GET("/ping", func(c *gin.Context) {
		c.Abort()
		fmt.Println("请求进来了！！！", num.Load())
		num.Add(1)
		if num.Load() >= 5 {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "模拟请求错误"})
	})
	serevr := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	http2.ConfigureServer(serevr, &http2.Server{})
	serevr.ListenAndServeTLS("../core_test/default.pem", "../core_test/default.key")
}
