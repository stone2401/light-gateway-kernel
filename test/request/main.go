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
	req, err := http.NewRequest("GET", "http://127.0.0.1:8083/base", nil)
	if err != nil {
		fmt.Println(err)
	}
	wg.Add(1)
	now := time.Now()
	go func() {
		for i := range 100000 {
			if i%1000 == 0 {
				fmt.Println(i)
				fmt.Println(time.Since(now))
				now = time.Now()
			}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err, "go 11")
				break
			}
			if resp.StatusCode != 200 {
				// 读取ｂｏｄｙ
				data, _ := io.ReadAll(resp.Body)
				fmt.Println(resp.StatusCode, string(data), "go 12")
				time.Sleep(1 * time.Second)
			}
			resp.Body.Close()
		}
		wg.Done()
	}()
	req2, err := http.NewRequest("GET", "http://127.0.0.1:8083/base", nil)
	if err != nil {
		fmt.Println(err)
	}
	wg.Add(1)
	now2 := time.Now()
	go func() {
		for i := range 100000 {
			if i%1000 == 0 {
				fmt.Println(i)
				fmt.Println(time.Since(now2))
				now2 = time.Now()
			}
			resp, err := client.Do(req2)
			if err != nil {
				fmt.Println(err, "go 2")
				break
			}
			if resp.StatusCode != 200 {
				// 读取ｂｏｄｙ
				data, _ := io.ReadAll(resp.Body)
				fmt.Println(resp.StatusCode, string(data), "go 2")
				time.Sleep(1 * time.Second)
			}
			resp.Body.Close()
		}
		wg.Done()
	}()
	wg.Wait()
}
