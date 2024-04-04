package main

import (
	"fmt"
	"time"

	"github.com/stone2401/light-gateway-kernel/pcore"
	"github.com/stone2401/light-gateway-kernel/pkg/sdk"
)

// 主函数：创建一个带有随机余额的SDK实例，配置节点和限流器，并启动代理服务器。
func main() {
    // 创建一个随机余额的SDK实例
    b := sdk.NewRandomBalance()
    // 向SDK实例添加一个节点
    b.AddNode("http://127.0.0.1:8080", 1)

    // 创建一个每秒最多处理10个请求的限流器
    limiter := pcore.NewLimiter(10)
    // 创建一个熔断器，超过5秒内处理5个请求失败后，接下来5秒内将拒绝处理请求
    fuse := pcore.NewFuseEntry("test", 5*time.Second, 5, 0.5)
    // 使用限流器和熔断器创建代理引擎
    proxy := pcore.NewEngine(b, limiter.ProxyHandler, fuse.FuseHandler)
    // 注册代理路由，匹配所有以"/{name}"开始的请求
    proxy.Register("/{name}")
    // 启动代理服务器并同步监听8083端口，返回启动状态
    fmt.Println(proxy.SyncStart(":8083"))
}
