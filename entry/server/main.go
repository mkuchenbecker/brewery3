package main

//nolint

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/mkuchenbecker/brewery3/brewery/gpio"
	"github.com/mkuchenbecker/brewery3/brewery/gpio/integration"
	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"github.com/mkuchenbecker/brewery3/brewery/servers"
	"github.com/mkuchenbecker/brewery3/brewery/servers/element"
	"github.com/mkuchenbecker/brewery3/brewery/servers/sensors"
	"github.com/mkuchenbecker/brewery3/brewery/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func startBrewery(port int, brewery *servers.Brewery) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serve := grpc.NewServer()
	model.RegisterBreweryServer(serve, brewery)
	go brewery.StartRunLoop()
	// Register reflection service on gRPC server.
	reflection.Register(serve)
	if err := serve.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func startThermometer(port int, address string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	serve := grpc.NewServer()
	addr := gpio.TemperatureAddress(address)
	err = integration.VerifyTemperatureAddress(addr)
	if err != nil {
		log.Fatalf("failed to read address: %v\n", err)
	}
	server, err := sensors.NewThermometerServer(integration.NewDefaultController(), addr)
	if err != nil {
		log.Fatalf("failed to make thermometer server: %v\n", err)
	}

	model.RegisterThermometerServer(serve, server)

	utils.Print("Serving Traffic")

	reflection.Register(serve)
	if err := serve.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}

func startHeater(port int, pin uint8) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serve := grpc.NewServer()
	model.RegisterSwitchServer(serve, element.NewHeaterServer(integration.NewDefaultController(), pin))
	// Register reflection service on gRPC server.
	reflection.Register(serve)
	if err := serve.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

const (
	mashAddr  = "28-0315712c08ff"
	mashPort  = 8110
	hermsAddr = "28-0315715039ff"
	hermsPort = 8111
	boilAddr  = "28-031571188aff"
	boilPort  = 8112

	elementPin  = 11
	elementPort = 8120
)

func makeTemperatureClient(port int, address string) (model.ThermometerClient, *grpc.ClientConn) {
	utils.Print(fmt.Sprintf("Connecting to client: %d", port))
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := model.NewThermometerClient(conn)
	res, err := client.Get(context.Background(), &model.GetRequest{})
	if err != nil {
		panic(err)
	}
	utils.Print(fmt.Sprintf("temp: %f", res.Temperature))
	return client, conn
}

func makeSwitchClient(port int, pin uint8) (model.SwitchClient, *grpc.ClientConn) {
	utils.Print(fmt.Sprintf("Connecting to client: %d", port))
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := model.NewSwitchClient(conn)
	_, err = client.Off(context.Background(), &model.OffRequest{})
	if err != nil {
		panic(err)
	}
	return client, conn
}

func main() {
	utils.Print("starting backends")
	go startThermometer(mashPort, mashAddr)
	go startThermometer(hermsPort, hermsAddr)
	go startThermometer(boilPort, boilAddr)
	go startHeater(elementPort, elementPin)

	utils.Print("waiting for discovery")
	time.Sleep(time.Second)
	utils.Print("starting clients")

	mash, mashConn := makeTemperatureClient(mashPort, mashAddr)
	defer mashConn.Close()
	herms, hermsConn := makeTemperatureClient(hermsPort, hermsAddr)
	defer hermsConn.Close()
	boil, boilConn := makeTemperatureClient(boilPort, boilAddr)
	defer boilConn.Close()

	element, elementConn := makeSwitchClient(elementPort, elementPin)
	defer elementConn.Close()

	brewery := servers.Brewery{
		MashSensor:  mash,
		HermsSensor: herms,
		BoilSensor:  boil,
		Element:     element,
	}
	startBrewery(8100, &brewery)
}
