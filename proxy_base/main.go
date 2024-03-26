package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	proxyAddr = "http://127.0.0.1:8081"
)

func Proxy(w http.ResponseWriter, r *http.Request) {
	uri, err := url.Parse(proxyAddr)
	if err != nil {
		return
	}
	r.URL.Path = uri.Path
	r.URL.Host = uri.Host
	r.URL.Scheme = uri.Scheme

	transport := http.DefaultTransport
	resp, err := transport.RoundTrip(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	for key, headle := range resp.Header {
		for _, head := range headle {
			w.Header().Add(key, head)
		}
	}
	io.Copy(w, resp.Body)
}

func main() {
	http.HandleFunc("/proxy", Proxy)
	http.ListenAndServe(":8083", nil)
}
