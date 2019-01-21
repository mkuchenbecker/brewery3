package rpi

import (
	"fmt"
	"log"
	"net"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartBrewery(port int, brewery *Brewery) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serve := grpc.NewServer()
	model.RegisterBreweryServer(serve, brewery)
	// Register reflection service on gRPC server.
	reflection.Register(serve)
	if err := serve.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
