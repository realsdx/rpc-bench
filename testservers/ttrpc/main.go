// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/realsdx/ttrpc-bench/testproto/pbttrpc"

	"github.com/containerd/ttrpc"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	// pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	var port = flag.Int("port", 50051, "The server port")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	s, err := ttrpc.NewServer()
	if err != nil {
		log.Fatalf("failed to satrt ttrpc: %v", err)
	}
	defer s.Close()

	pb.RegisterGreeterService(s, &server{})
	log.Printf("TTRPC Server listening at %v", lis.Addr())
	if err := s.Serve(context.Background(), lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
