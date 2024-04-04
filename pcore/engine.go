package pcore

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gorilla/mux"
	"github.com/stone2401/light-gateway-kernel/pkg/sdk"
)

// Engine http 代理引擎
// app *gin.Engine  gin引擎,用于注册路由
// proxy *http util.ReverseProxy 反向代理
// server *http.Server 服务，用于监听
type Engine struct {
	proxy    *httputil.ReverseProxy
	server   *http.Server
	router   *mux.Router
	Handlers []Handler
}

type Context struct {
	Request  *http.Request
	Response *ResponseWriter
	index    int8
	handlers []Handler
	context.Context
}

type Handler func(ctx *Context)

// ResponseWriter 响应写入器
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader 写入响应头
func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// NewEngine 创建一个新的Engine实例。
//
// 参数:
// b - sdk.Balance类型，负载均衡器。
// handler - 可变参数，类型为Handler，中间价限流等功能。
//
// 返回值:
// 返回一个指向Engine结构体的指针。
func NewEngine(b sdk.Balance, handler ...Handler) *Engine {
	// 创建一个新的Engine实例，初始化proxy和router。
	return &Engine{
		proxy:    sdk.NewSingleHostReverseProxy(b),
		router:   mux.NewRouter(),
		Handlers: handler,
	}
}

// Register 注册路由
// path string 路由
// header gin.HandlerFunc 处理函数
func (e *Engine) Register(path string, h ...Handler) error {
	// 添加 *action
	var err error
	defer func() {
		if err := recover(); err != nil {
			var ok = true
			err, ok = err.(error)
			if !ok {
				err = errors.New("unknow panic")
			}

		}
	}()
	e.router.PathPrefix(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := &ResponseWriter{ResponseWriter: w}
		ctx := &Context{
			Request:  r,
			Response: resp,
			index:    -1,
			handlers: append(e.Handlers, append(h, func(ctx *Context) {
				e.proxy.ServeHTTP(resp, ctx.Request)
			})...),
			Context: context.Background(),
		}
		ctx.Next()
	})
	return err
}

// Use 添加中间件
func (e *Engine) Use(h ...Handler) {
	e.Handlers = append(e.Handlers, h...)
}

// Start 启动服务
func (e *Engine) Start(addr string) error {
	go e.SyncStart(addr)
	return nil
}

// SyncStart 同步启动服务
func (e *Engine) SyncStart(addr string) error {
	fmt.Println("start server")
	if e.server != nil {
		return errors.New("engine is running")
	}
	if addr == "" {
		return errors.New("address is empty")
	}
	e.server = &http.Server{
		Addr:    addr,
		Handler: e.router,
	}
	return e.server.ListenAndServe()
}

func (e *Engine) Stop() {
	ctx := context.Background()
	go func() {
		select {
		case <-ctx.Done():
		case <-time.After(15 * time.Second):
			ctx.Done()
		}
	}()
	e.server.Shutdown(ctx)
}

func (ctx *Context) Next() {
	ctx.index++
	for ctx.index < int8(len(ctx.handlers)) {
		ctx.handlers[ctx.index](ctx)
		ctx.index++
	}
}

func (ctx *Context) Abort() {
	ctx.index = int8(len(ctx.handlers))
}
