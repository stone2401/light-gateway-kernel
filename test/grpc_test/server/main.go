package main

import (
	"context"
	"fmt"
	"net"

	"github.com/stone2401/light-gateway-kernel/test/grpc_test/echo"
	"google.golang.org/grpc"
)

type EchoServer struct {
	echo.EchoServer
}

func (e *EchoServer) UnaryEcho(_ context.Context, message *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{
		Message: message.Message,
	}, nil
}

func (e *EchoServer) ServerStreamingEcho(_ *echo.EchoRequest, stream echo.Echo_ServerStreamingEchoServer) error {
	for i := 0; i < 10; i++ {
		stream.Send(&echo.EchoResponse{
			Message: "hello",
		})
	}
	return nil
}

func (e *EchoServer) ClientStreamingEcho(stream echo.Echo_ClientStreamingEchoServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			return err
		}
		fmt.Printf("ClientStreamingEcho in.Message: %v\n", in.Message)
	}
}

func (e *EchoServer) BidirectionalStreamingEcho(stream echo.Echo_BidirectionalStreamingEchoServer) error {
	for i := 0; i < 10; i++ {
		in, err := stream.Recv()
		if err != nil {
			return err
		}
		fmt.Printf("BidirectionalStreamingEcho in.Message: %v\n", in.Message)
		err = stream.Send(&echo.EchoResponse{
			Message: "hello",
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":2500")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	echo.RegisterEchoServer(s, &EchoServer{})
	s.Serve(lis)

}
