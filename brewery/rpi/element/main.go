package element

import (
	"fmt"
	"log"
	"net"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"github.com/mkuchenbecker/brewery3/brewery/rpi/gpio"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartHeater(port int, pin uint8) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serve := grpc.NewServer()
	model.RegisterSwitchServer(serve, &HeaterServer{pin: pin, ctrl: gpio.GetDefaultController()})
	// Register reflection service on gRPC server.
	reflection.Register(serve)
	if err := serve.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
