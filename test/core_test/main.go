package main

import (
	"errors"
	"net/http"

	"github.com/stone2401/light-gateway-kernel/pcore"
	"github.com/stone2401/light-gateway-kernel/pkg/sdk"
)

func main() {
	b := sdk.NewRandomBalance()
	b.AddNode("http://127.0.0.1:8080/base", 1)
	b.AddNode("http://127.0.0.1:8081/base", 1)

	proxy := pcore.NewEngine(b, pcore.NewRateLimiter(1))
	proxy.Register("/base", func(r *http.Request) (code int, err error) {
		return http.StatusTooManyRequests, errors.New("error")
	})
	proxy.Register("/base1", func(r *http.Request) (code int, err error) {
		return 200, nil
	})
	proxy.Start(":8083")
}
