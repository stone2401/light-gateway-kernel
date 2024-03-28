package pcore

import (
	"context"
	"errors"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stone2401/light-gateway-kernel/pkg/sdk"
)

type Engine struct {
	app    *gin.Engine
	proxy  *httputil.ReverseProxy
	server *http.Server
}

type Handler func(*http.Request) (code int, err error)

func NewEngine(b sdk.Balance, handler ...Handler) *Engine {
	gin.SetMode(gin.ReleaseMode)
	ginHandler := []gin.HandlerFunc{}
	for _, h := range handler {
		ginHandler = append(ginHandler, func(c *gin.Context) {
			code, err := h(c.Request)
			if err != nil {
				c.String(code, err.Error())
				c.Abort()
			}
			c.Next()
		})
	}
	app := gin.New()
	app.Use(ginHandler...)
	return &Engine{
		app:   app,
		proxy: sdk.NewSingleHostReverseProxy(b),
	}
}

func (e *Engine) Register(path string, headler ...Handler) error {
	// 判断是不是 /结尾
	if path[len(path)-1] != '/' {
		path = path + "/"
	}
	// 添加 *action
	path = path + "*action"
	var err error
	defer func() {
		if err := recover(); err != nil {
			err = errors.New(err.(string))
		}
	}()
	ginHandler := []gin.HandlerFunc{}
	for _, h := range headler {
		ginHandler = append(ginHandler, func(c *gin.Context) {
			code, err := h(c.Request)
			if err != nil {
				c.String(code, err.Error())
				c.Abort()
			}
			c.Next()
		})
	}
	ginHandler = append(ginHandler, func(ctx *gin.Context) {
		ctx.Request.URL.Path = ctx.Param("action")
		e.proxy.ServeHTTP(ctx.Writer, ctx.Request)
	})
	e.app.Any(path, ginHandler...)
	return err
}

func (e *Engine) Start(addr string) error {
	if e.server != nil {
		return errors.New("engine is running")
	}
	if addr == "" {
		return errors.New("address is empty")
	}
	e.server = &http.Server{
		Addr:    addr,
		Handler: e.app,
	}
	return e.server.ListenAndServe()
}

func (e *Engine) SyncStart(addr string) error {
	if e.server != nil {
		return errors.New("engine is running")
	}
	if addr == "" {
		return errors.New("address is empty")
	}
	e.server = &http.Server{
		Addr:    addr,
		Handler: e.app,
	}
	go e.server.ListenAndServe()
	return nil
}

func (e *Engine) Stop() {
	ctx := context.Background()
	go func() {
		select {
		case <-ctx.Done():
		case <-time.After(15 * time.Second):
			e.server.Shutdown(ctx)
		}
	}()
	e.server.Shutdown(ctx)
}
