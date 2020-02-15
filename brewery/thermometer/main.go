//+build !test

package main

import (
	"fmt"
	"log"
	"net"

	"github.com/mkuchenbecker/brewery3/brewery/gpio/integration"
	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"github.com/mkuchenbecker/brewery3/brewery/servers/element"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() { // nolint: deadcode
	port := 9100
	pin := uint8(1)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serve := grpc.NewServer()
	heater := element.NewHeaterServer(integration.NewDefaultController(), pin)
	model.RegisterSwitchServer(serve, heater)
	// Register reflection service on gRPC server.
	reflection.Register(serve)
	if err := serve.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
