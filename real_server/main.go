package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type RealServer struct {
	addr string
}

func NewRealServer(addr string) *RealServer {
	return &RealServer{
		addr: addr,
	}
}

func (r *RealServer) Hello(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("http request: %v, addr: %v\n", req.URL.Path, r.addr)
	w.Write([]byte("hello" + r.addr))
}

func (r *RealServer) ErrorHandle(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("http request: %v\n", req.URL.Path)
	w.Write([]byte("error"))
}

func (r *RealServer) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.Hello)
	mux.HandleFunc("/base/error", r.ErrorHandle)
	go func() {
		http.ListenAndServe(r.addr, mux)
	}()
}

func main() {
	r := NewRealServer(":8080")
	r.Run()
	r2 := NewRealServer(":8081")
	r2.Run()
	// 结束监听
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
