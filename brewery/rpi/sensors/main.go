package sensors

import (
	"fmt"
	"log"
	"net"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"github.com/mkuchenbecker/brewery3/brewery/rpi/gpio"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartThermometer(port int, address string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	serve := grpc.NewServer()
	addr, err := gpio.NewTemperatureAddress(address, &gpio.DefaultSensorArray{})
	if err != nil {
		log.Fatalf("failed to read address: %v\n", err)
	}
	model.RegisterThermometerServer(serve, &ThermometerServer{address: addr})
	// Register reflection service on gRPC server.

	reflection.Register(serve)
	if err := serve.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
