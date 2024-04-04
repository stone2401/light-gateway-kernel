# GatwayReverseProxy 实现

## 1. httputil.ReversProxy
go 官方提供的反向代理实现了HTTPServer方法，此方法在整个http库中的做用都是处理请求．

我们可以通过重写　`Director` , `ModifyResponse` , `ErrorHandler`　方法修改代理的行为

`Director`：请求转发前回调方法，主要用与请求拦截，转发地址修改等

`ModifyResponse`：请求转发结束后回调，此方法内，我们可以修改返回内容，或者更具请求返回值返回err进行统计

`ErrorHandler`：接收　ModifyResponse　等请求中的错误回调，可以用与错误处理

`Transport`：中可以设置超时时间等内容

## 2. GatwayReverseProxy　
在　`GatwayReverseProxy`　中我们重写了　`Director`　方法，并且支持传入一个　`负载均衡器`

负载均衡器独立与整个代理过程，仅仅提供获取节点方法即可，自行组织整个节点添加以及服务发现等功能

`server`：http.Server　将用于端口监听，服务关闭工作（`Shutdown`方法）