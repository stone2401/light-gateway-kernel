package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	urlStr := "http://127.0.0.1:8080/base"
	uri, err := url.Parse(urlStr)
	if err != nil {
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(uri)
	proxy.Director = func(r *http.Request) {
		r.URL.Scheme = uri.Scheme
		r.URL.Path = uri.Path
		r.URL.Host = uri.Host
		fmt.Println("请求打进来了！！！")
	}
	proxy.ModifyResponse = func(r *http.Response) error {
		fmt.Println("响应返回了！")
		return nil
	}
	http.ListenAndServe(":8083", proxy)
}
