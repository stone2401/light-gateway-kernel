package sdk

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/stone2401/light-gateway-kernel/pkg/zlog"
)

// 代理服务
// ctx context.Context 上下文
// cancel context.CancelFunc 取消上下文
// addr string 代理服务地址
// reversProxy *httputil.ReverseProxy 反向代理核心
// server http.Server 服务，用于监听代理服务
type GatwayReverseProxy struct {
	ctx         context.Context
	cancel      context.CancelFunc
	addr        string // proxy server port
	reversProxy *httputil.ReverseProxy
	server      http.Server
}

// 代理拦截
// Director 代理发送前的拦截
// ModifyResponse 代理接收后的拦截
// ErrorHandler 代理错误的拦截
type Interceptor interface {
	// 代理发送前的拦截
	Director(*http.Request)
	// 代理接收后的拦截
	ModifyResponse(*http.Response) error
	// 代理错误的拦截
	ErrorHandler(http.ResponseWriter, *http.Request, error)
}

// 初始化 reverse proxy
// addr string 代理服务地址
// balancer Balance 负载均衡器
func NewGatwayReverseProxy(addr string, balancer Balance) *GatwayReverseProxy {
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

// 启动代理服务
func (g *GatwayReverseProxy) Start() error {
	zlog.Zlog().Info("start proxy server")
	go g.server.ListenAndServe()
	return nil
}

// 启动代理服务
func (g *GatwayReverseProxy) AsyncStart() error {
	zlog.Zlog().Info("start proxy server")
	go g.server.ListenAndServe()
	return nil
}

// 停止代理服务
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

// 初始化　reverse proxy
// balancer Balance 负载均衡器
func NewSingleHostReverseProxy(balance Balance) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		// TODO: 负载均衡
		token := req.Header.Get("X-Forwarded-For")
		addr, err := balance.GetNode(token)
		if err != nil {
			if errors.Is(err, ErrorNotFoundNode) {
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
	return &httputil.ReverseProxy{Director: director, Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			// RootCAs: func() *x509.CertPool {
			// 	pool := x509.NewCertPool()
			// 	file, _ := os.ReadFile("./ca/server.crt")
			// 	fmt.Println(string(file))
			// 	pool.AppendCertsFromPEM(file)
			// 	return pool
			// }(),
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}}
}

func rewriteRequestURL(req *http.Request, target *url.URL) {
	targetQuery := target.RawQuery
	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
	req.Host = target.Host
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
