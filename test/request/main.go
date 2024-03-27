package main

import (
	"fmt"
	"net/http"
	"sync"
)

var wg sync.WaitGroup

func main() {
	wg.Add(3)
	go func() {
		for range 100 {
			resp, err := http.Get("http://127.0.0.1:8083/base1/hello")
			if err != nil {
				fmt.Println(err)
				break
			}
			if resp.StatusCode != 200 {
				fmt.Println(resp.StatusCode)
				break
			}
		}
		wg.Done()
	}()
	go func() {
		for range 100 {
			resp, err := http.Get("http://127.0.0.1:8083/base1/hello")
			if err != nil {
				fmt.Println(err)
				break
			}
			if resp.StatusCode != 200 {
				fmt.Println(resp.StatusCode)
				break
			}
		}
		wg.Done()
	}()
	go func() {
		for range 100 {
			resp, err := http.Get("http://127.0.0.1:8083/base1/hello")
			if err != nil {
				fmt.Println(err)
				break
			}
			if resp.StatusCode != 200 {
				fmt.Println(resp.StatusCode)
				break
			}
		}
		wg.Done()
	}()
	wg.Wait()
}
