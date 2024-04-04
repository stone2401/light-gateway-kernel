package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.Use(func(h http.Handler) http.Handler {
		if h == nil {
			return nil
		}
		return h
	})
	// 匹配根路径 /
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, this is root path!")
	})

	// 匹配 /api 路径
	r.PathPrefix("/api").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, this is /api path!")
	})

	// 匹配 /api/v1/ping 路径
	r.HandleFunc("/api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, this is /api/v1/ping path!")
	})
	// 创建一个HTTP服务器并指定路由
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	// 启动一个goroutine，10秒后注册 /hello 路由
	go func() {
		time.Sleep(10 * time.Second)
		r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, this is /hello path!")
		})
	}()
	server.ListenAndServe()
}
