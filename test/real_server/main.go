package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/stone2401/light-gateway-kernel/pcore"
	"github.com/stone2401/light-gateway-kernel/pkg/sdk"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/websocket"
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
	w.WriteHeader(500)
	w.Write([]byte("error"))
}

func (r *RealServer) IndexHtml(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("http request: %v %s\n", req.URL.Path, r.addr)
	w.WriteHeader(200)
	f, err := os.ReadFile("index.html")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("index.html not found"))
		return
	}
	w.Write(f)
}

func (r *RealServer) WebSocketsHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("http request: %v %s\n", req.URL.Path, r.addr)
	conn := websocket.Handler(func(c *websocket.Conn) {
		defer c.Close()
		w, err := c.NewFrameWriter(websocket.TextFrame)
		if err != nil {
			return
		}
		for {
			r, err := c.NewFrameReader()
			if err != nil {
				return
			}
			msg, err := io.ReadAll(r)
			if err != nil {
				return
			}
			cmd := exec.Command("zsh", "-c", string(msg))
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				return
			}
			cmd.Stderr = cmd.Stdout
			if err = cmd.Start(); err != nil {
				return
			}
			io.Copy(w, stdout)
		}
	})
	conn.ServeHTTP(w, req)
}

func (r *RealServer) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/base", r.Hello)
	mux.HandleFunc("/base/error", r.ErrorHandle)
	mux.HandleFunc("/index.html", r.IndexHtml)
	mux.HandleFunc("/api/ws", r.WebSocketsHandler)
	go func() {
		http.ListenAndServe(r.addr, mux)
	}()
}

func main() {
	r := NewRealServer(":8082")
	r.Run()
	r2 := NewRealServer(":8081")
	r2.Run()
	client, err := clientv3.New(clientv3.Config{Endpoints: []string{"127.0.0.1:2379"}})
	if err != nil {
		panic(err)
	}
	has := sdk.NewMurmurHasher()
	nodes := []*pcore.NodeInfo{
		{
			Ip:     "http://127.0.0.1:8082",
			Weight: 1,
		},
		{
			Ip:     "http://127.0.0.1:8081",
			Weight: 1,
		},
	}
	for _, node := range nodes {
		b, _ := json.Marshal(node)
		enc := has.Encrypt(string(b))
		client.Put(context.Background(), "8080/base"+strconv.Itoa(int(enc)), string(b))
	}
	// 结束监听
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// http.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
	// 	// w.Write([]byte("pong"))
	// 	http.NotFound(w, r)
	// })
	// go func() {
	// 	defer func() {
	// 		if p := recover(); p != nil {
	// 			fmt.Println(p)
	// 		}
	// 	}()
	// 	time.Sleep(10 * time.Second)
	// 	http.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {})
	// }()
	// err := http.ListenAndServe(":8090", nil)
	// if err != nil {
	// 	panic(err)
	// }
}
