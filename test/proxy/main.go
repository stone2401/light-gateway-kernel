package main

import (
	"net/http"

	loadbalance "github.com/stone2401/light-gateway-kernel/pcore/load_balance"
	"github.com/stone2401/light-gateway-kernel/pkg/sdk"
)

func main() {
	b := loadbalance.NewRandomBalance()
	b.AddNode("http://localhost:2401", 1)
	proxy := sdk.NewSingleHostReverseProxy(b)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
	http.ListenAndServe(":2400", mux)
}
