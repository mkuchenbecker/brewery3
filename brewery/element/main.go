package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/mkuchenbecker/brewery3/brewery/utils"

	"github.com/mkuchenbecker/brewery3/brewery/element/element"
	"github.com/mkuchenbecker/brewery3/brewery/gpio/integration"
	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() { // nolint: deadcode
	strPort := os.Getenv("PORT")
	strPin := os.Getenv("GPIO_PIN")

	port, err := strconv.ParseInt(strPort, 10, 32)
	if err != nil {
		log.Fatalf("Invalid port is not 32 bit int: %s", strPort)
	}

	pin, err := strconv.ParseInt(strPin, 10, 8)
	if err != nil {
		log.Fatalf("invalid pin : %s", strPin)
	}

	utils.Printf("Starting heater on port: %d; pin: %d", port, pin)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serve := grpc.NewServer()
	heater := element.NewHeaterServer(integration.NewDefaultController(), uint8(pin))
	model.RegisterSwitchServer(serve, heater)
	// Register reflection service on gRPC server.
	reflection.Register(serve)
	if err := serve.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
