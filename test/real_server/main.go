package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

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
	w.WriteHeader(500)
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
	r := NewRealServer(":8080")
	r.Run()
	r2 := NewRealServer(":8081")
	r2.Run()
	// 结束监听
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
