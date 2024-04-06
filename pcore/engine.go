package pcore

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/stone2401/light-gateway-kernel/pkg/sdk"
)

// Engine http 代理引擎
// app *gin.Engine  gin引擎,用于注册路由
// proxy *http util.ReverseProxy 反向代理
// server *http.Server 服务，用于监听
type Engine struct {
	proxy    *httputil.ReverseProxy
	server   *http.Server
	mux      *http.ServeMux
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
	Code int
}

// WriteHeader 写入响应头
func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.Code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Hijack 劫持连接
func (w *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("ResponseWriter does not support Hijack")
	}
	return hijacker.Hijack()
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
		mux:      http.NewServeMux(),
		Handlers: handler,
	}
}

// Register 注册路由
// path string 路由
// header gin.HandlerFunc 处理函数
func (e *Engine) Register(path string, b sdk.Balance, h ...Handler) error {
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
	proxy := sdk.NewSingleHostReverseProxy(b)
	if b == nil {
		proxy = e.proxy
	}

	e.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{
			Request:  r,
			Response: &ResponseWriter{ResponseWriter: w},
			handlers: append(e.Handlers, append(h, func(ctx *Context) {
				proxy.ServeHTTP(ctx.Response, ctx.Request)
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

func (e *Engine) initServer(addr string) error {
	if e.server != nil {
		return errors.New("engine is running")
	}
	if addr == "" {
		return errors.New("address is empty")
	}
	e.server = &http.Server{
		Addr:    addr,
		Handler: e.mux,
	}
	return nil
}

// Start 启动服务
func (e *Engine) Start(addr string) error {
	if err := e.initServer(addr); err != nil {
		return err
	}
	return e.server.ListenAndServe()
}

func (e *Engine) StartTls(addr string, certFile, keyFile string) error {
	if err := e.initServer(addr); err != nil {
		return err
	}
	return e.server.ListenAndServeTLS(certFile, keyFile)
}

// AsyncStart 同步启动服务
func (e *Engine) AsyncStart(addr string) error {
	if err := e.initServer(addr); err != nil {
		return err
	}
	go func() {
		_ = e.server.ListenAndServe()
	}()
	return nil
}

// AsyncStartTls 同步启动服务
func (e *Engine) AsyncStartTls(addr string, certFile, keyFile string) error {
	if err := e.initServer(addr); err != nil {
		return err
	}
	go func() {
		_ = e.server.ListenAndServeTLS(certFile, keyFile)
	}()
	return nil
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
	err := e.server.Shutdown(ctx)
	if err != nil {
		return
	}
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
