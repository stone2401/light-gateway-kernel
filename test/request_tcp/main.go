package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	// 连接tcp
	conn, err := net.Dial("tcp", ":2400")
	if err != nil {
		return
	}
	for {
		conn.Write([]byte("hello tcp\r\n"))
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(buf[:n]))
		time.Sleep(1 * time.Second)
	}
}
