package main

import (
	"github.com/stone2401/light-gateway-kernel/pcore"
	"github.com/stone2401/light-gateway-kernel/pcore/load_balance"
)

func main() {
	b := load_balance.NewRandomBalance()
	b.AddNode("127.0.0.1:6379", 1)
	b.AddNode("127.0.0.1:6380", 1)
	b.AddNode("127.0.0.1:6381", 1)
	proxy := pcore.NewTcpEngine(b, 0, 0, 0)
	proxy.ListenAndServe("tcp", ":2400")
}
