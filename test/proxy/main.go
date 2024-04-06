package main

import (
	"net/http"

	"github.com/stone2401/light-gateway-kernel/pkg/sdk"
)

func main() {
	b := sdk.NewRandomBalance()
	b.AddNode("http://localhost:2401", 1)
	proxy := sdk.NewSingleHostReverseProxy(b)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
	http.ListenAndServe(":2400", mux)
}
