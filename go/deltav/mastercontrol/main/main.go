package main

import (
	"fmt"
	"log"
	"net"
	"time"

	master "github.com/golang001/deltav/mastercontrol"
	protos "github.com/golang001/deltav/model/gomodel"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct {
	mc *master.MasterControl
}

func (s *server) Register(ctx context.Context, req *protos.RegisterRequest) (*protos.RegisterResponse, error) {
	fmt.Printf("Register: %s - %s\n", "req.ID", time.Now().String())
	return &protos.RegisterResponse{Effect: "success"}, nil
}

func (s *server) Detect(ctx context.Context, req *protos.DetectRequest) (*protos.DetectResponse, error) {
	return &protos.DetectResponse{}, nil
}

func (s *server) Get(ctx context.Context, req *protos.GetRequest) (*protos.GetResponse, error) {
	return &protos.GetResponse{}, nil
}

func (s *server) Initialize(ctx context.Context, req *protos.InitializeRequest) (*protos.InitializeResponse, error) {
	return &protos.InitializeResponse{}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serve := grpc.NewServer()
	protos.RegisterWorldModelServer(serve, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(serve)
	if err := serve.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
