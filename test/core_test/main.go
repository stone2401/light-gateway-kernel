package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/stone2401/light-gateway-kernel/pcore"
	"github.com/stone2401/light-gateway-kernel/pkg/sdk"
)

// 主函数：创建一个带有随机余额的SDK实例，配置节点和限流器，并启动代理服务器。
func main() {
	// 启动pprof，在端口6060上监听以进行性能分析
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	// 创建一个随机余额的SDK实例
	b := sdk.NewRandomBalance()
	// 向SDK实例添加一个节点
	b.AddNode("http://127.0.0.1:8080", 1)

	// 创建一个每秒最多处理10个请求的限流器
	limiter := pcore.NewLimiter(10000)
	// 创建一个熔断器，超过5秒内处理5个请求失败后，接下来5秒内将拒绝处理请求
	fuse := pcore.NewFuseEntry("test", 5*time.Second, 5, 0.5)

	// 创建一个计数器，每10秒打印一次计数和时间``
	counter := pcore.NewCounter(10)

	// 使用限流器和熔断器创建代理引擎
	proxy := pcore.NewEngine(b, counter.CounterHandler, limiter.ProxyHandler, fuse.FuseHandler)
	// 注册代理路由，匹配所有以"/{name}"开始的请求
	proxy.Register("/base")

	go func() {
		// 循环监听计数器增加的事件，并打印计数和时间
		for {
			entry := <-counter.Gain()
			fmt.Println(entry.Count, entry.Time.Format("2006-01-02 15:04:05"))
		}
	}()
	// 启动代理服务器并同步监听8083端口，返回启动状态
	fmt.Println(proxy.Start(":8083"))
}
