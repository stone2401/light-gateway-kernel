package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	client := http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 5,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 5 * time.Second,
			}).DialContext,
			// TLSClientConfig: &tls.Config{
			// 	InsecureSkipVerify: true},
		},
	}
	go requestFunc(client, "go 1")
	wg.Add(1)
	go requestFunc(client, "go 2")
	wg.Wait()
}

func requestFunc(client http.Client, name string) {
	req2, err := http.NewRequest("GET", "http://127.0.0.1:8083/base", nil)
	if err != nil {
		fmt.Println(err)
	}
	now2 := time.Now()
	for i := range 100000 {
		if i%1000 == 0 {
			fmt.Println(i)
			fmt.Println(time.Since(now2))
			now2 = time.Now()
		}
		resp, err := client.Do(req2)
		if err != nil {
			fmt.Println(err, name)
			break
		}
		if resp.StatusCode != 200 {
			data, _ := io.ReadAll(resp.Body)
			fmt.Println(resp.StatusCode, string(data), name)
			time.Sleep(1 * time.Second)
		}
		resp.Body.Close()
	}
	wg.Done()
}
