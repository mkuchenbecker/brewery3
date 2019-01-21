package sensors

import (
	"fmt"
	"log"
	"net"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"github.com/mkuchenbecker/brewery3/brewery/rpi/gpio"
	"github.com/mkuchenbecker/brewery3/brewery/utils"
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
	server := &ThermometerServer{address: addr,
		ctrl: gpio.GetDefaultController()}
	model.RegisterThermometerServer(serve, server)

	utils.Print("Serving Traffic")

	reflection.Register(serve)
	if err := serve.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
