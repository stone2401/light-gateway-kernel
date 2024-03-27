package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/stone2401/light-gateway-kernel/pcore"
	"github.com/stone2401/light-gateway-kernel/pkg/balance"
)

func main() {
	b := balance.NewRobinBalance()
	
	b.AddNode("http://127.0.0.1:8080", 1)
	b.AddNode("http://127.0.0.1:8081", 1)
	pcore.NewGatwayReverseProxy(":8083", b).Start()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
