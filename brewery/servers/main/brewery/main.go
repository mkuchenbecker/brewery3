//+build !test

package main

//nolint

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/mkuchenbecker/brewery3/brewery/logger"

	"github.com/kelseyhightower/envconfig"
	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"github.com/mkuchenbecker/brewery3/brewery/servers/brewery"
	"github.com/mkuchenbecker/brewery3/brewery/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func startBrewery(port int, brewery *brewery.Brewery) {
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

type Settings struct {
	MashThermPort0  int `envconfig:"MASH_THERM_PORT_0" default:"9110"`
	HermsThermPort0 int `envconfig:"HERMS_THERM_PORT_0" default:"9111"`
	BoilThermPort0  int `envconfig:"BOIL_THERM_PORT_0" default:"9112"`

	ElementPort0 int `envconfig:"ELEMENT_PORT_0" default:"9100"`

	BreweryPort0 int `envconfig:"BREWERY_PORT_0" default:"9000"`
}

func makeTemperatureClient(port int) (model.ThermometerClient, *grpc.ClientConn) {
	utils.Print(fmt.Sprintf("Connecting to client: %d", port))
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := model.NewThermometerClient(conn)
	return client, conn
}

func makeSwitchClient(port int) (model.SwitchClient, *grpc.ClientConn) {
	utils.Print(fmt.Sprintf("Connecting to client: %d", port))
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := model.NewSwitchClient(conn)
	return client, conn
}

func getSettings(prefix string) *Settings {
	var s Settings
	err := envconfig.Process(prefix, &s)
	if err != nil {
		log.Fatal(context.Background(), err.Error())
	}
	return &s
}

func main() { // nolint: deadcode
	utils.Print("waiting for discovery")
	time.Sleep(time.Second)
	utils.Print("starting clients")

	settings := getSettings("")

	mash, mashConn := makeTemperatureClient(settings.MashThermPort0)
	defer mashConn.Close()
	herms, hermsConn := makeTemperatureClient(settings.HermsThermPort0)
	defer hermsConn.Close()
	boil, boilConn := makeTemperatureClient(settings.BoilThermPort0)
	defer boilConn.Close()
	element, elementConn := makeSwitchClient(settings.ElementPort0)
	defer elementConn.Close()

	brewery := brewery.Brewery{
		MashSensor:  mash,
		HermsSensor: herms,
		BoilSensor:  boil,
		Element:     element,
		Logger:      logger.NewFake(),
	}
	startBrewery(settings.BreweryPort0, &brewery)
}
