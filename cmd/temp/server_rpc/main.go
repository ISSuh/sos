package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ISSuh/sos/internal/rpc"
	"google.golang.org/grpc"
)

type server struct {
	rpc.UnimplementedGreeterServer
}

func (s *server) SayHello(c context.Context, res *rpc.HelloRequest) (*rpc.HelloReply, error) {
	fmt.Printf("[TEST] name : %s\n", res.Name)
	return &rpc.HelloReply{
		Message: "test",
	}, nil
}

func main() {
	address := "0.0.0.0:33669"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := &server{}
	s := grpc.NewServer()
	rpc.RegisterGreeterServer(s, server)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
