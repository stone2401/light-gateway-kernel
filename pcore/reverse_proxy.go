package pcore

import (
	"context"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/stone2401/light-gateway-kernel/pkg/balance"
	"github.com/stone2401/light-gateway-kernel/pkg/zlog"
)

type GatwayReverseProxy struct {
	ctx         context.Context
	cancel      context.CancelFunc
	addr        string // proxy server port
	reversProxy *httputil.ReverseProxy
	server      http.Server
}

func NewGatwayReverseProxy(addr string, balancer balance.Balance) *GatwayReverseProxy {
	reversProxy := NewSingleHostReverseProxy(balancer)
	ctx, cancel := context.WithCancel(context.Background())
	return &GatwayReverseProxy{
		ctx:         ctx,
		cancel:      cancel,
		addr:        addr,
		reversProxy: reversProxy,
		server: http.Server{
			Addr:    addr,
			Handler: reversProxy,
		},
	}
}

func (g *GatwayReverseProxy) Start() error {
	zlog.Zlog().Info("start proxy server")
	go g.server.ListenAndServe()
	return nil
}

func (g *GatwayReverseProxy) Stop() error {
	go func() {
		// 等待 15 seconds 执行 cancel
		select {
		case <-time.After(15 * time.Second):
			g.cancel()
		case <-g.ctx.Done():
		}
	}()
	err := g.server.Shutdown(g.ctx)
	if err != nil {
		zlog.Zlog().Error(err.Error())
	}
	g.cancel()
	return nil
}

func NewSingleHostReverseProxy(balancer balance.Balance) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		token := req.Header.Get("X-Forwarded-For")
		addr, err := balancer.GetNode(token)
		if err != nil {
			if errors.Is(err, balance.ErrorNotFoundNode) {
				// TODO: log
				zlog.Zlog().Error("not found node")
				return
			}
		}
		target, err := url.Parse(addr)
		if err != nil {
			return
		}
		rewriteRequestURL(req, target)
	}
	return &httputil.ReverseProxy{Director: director}
}

func rewriteRequestURL(req *http.Request, target *url.URL) {
	targetQuery := target.RawQuery
	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
	if targetQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = targetQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}
