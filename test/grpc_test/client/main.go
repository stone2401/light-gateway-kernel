package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/stone2401/light-gateway-kernel/test/grpc_test/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:2500", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	defer conn.Close()
	ech := echo.NewEchoClient(conn)
	fmt.Println("ech.UnaryEcho")
	fmt.Println(ech.UnaryEcho(context.Background(), &echo.EchoRequest{Message: "hello"}))
	{
		fmt.Println("ech.ClientStreamingEcho")
		stream, err := ech.ClientStreamingEcho(context.Background())
		if err != nil {
			return
		}

		for i := 0; i < 10; i++ {
			stream.Send(&echo.EchoRequest{Message: "hello"})
		}
		resp, _ := stream.CloseAndRecv()
		fmt.Println(resp)
	}
	time.Sleep(time.Second)
	{
		fmt.Println("ech.ServerStreamingEcho")
		stream, err := ech.ServerStreamingEcho(context.Background(), &echo.EchoRequest{Message: "hello"})
		if err != nil {
			return
		}
		for {
			resp, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				return
			}
			fmt.Println(resp)
		}
	}
	time.Sleep(time.Second)
	{
		fmt.Println("ech.BidirectionalStreamingEcho")
		stream, err := ech.BidirectionalStreamingEcho(context.Background())
		if err != nil {
			return
		}
		for i := 0; i < 10; i++ {
			stream.Send(&echo.EchoRequest{Message: "hello"})
			fmt.Println(stream.Recv())
		}
		stream.CloseSend()
	}
}
