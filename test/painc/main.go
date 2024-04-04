package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	conn := atomic.Int64{}
	conn.Store(100)
	fmt.Println(conn.Load())
}
