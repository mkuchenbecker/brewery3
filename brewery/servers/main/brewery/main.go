//+build !test

package main

//nolint

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"gopkg.in/yaml.v2"

	"github.com/kelseyhightower/envconfig"
	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"github.com/mkuchenbecker/brewery3/brewery/servers/brewery"
	"github.com/mkuchenbecker/brewery3/brewery/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func startBrewery(port int32, brewery *brewery.Brewery) {
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
	Config string `envconfig:"CONFIG_FILENAME" default:"config.yaml"`
}

func makeTemperatureClient(address string, port int32) (model.ThermometerClient, *grpc.ClientConn) {
	utils.Print(fmt.Sprintf("Connecting to client: %d", port))
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", address, port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := model.NewThermometerClient(conn)
	return client, conn
}

func makeSwitchClient(address string, port int32) (model.SwitchClient, *grpc.ClientConn) {
	utils.Print(fmt.Sprintf("Connecting to client: %d", port))
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", address, port), grpc.WithInsecure())
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

	settings := getSettings("")
	data, err := ioutil.ReadFile(settings.Config)
	if err != nil {
		panic(err)
	}
	brewSettings := model.BrewerySettings{}
	yaml.Unmarshal(data, &brewSettings)

	mainBrewrey := brewery.Brewery{
		Temperatures: make(map[brewery.SensorName]model.ThermometerClient),
	}
	for _, therm := range brewSettings.Thermometers {
		thermomerter, connection := makeTemperatureClient(therm.Addr, therm.Port)
		defer func() {
			err := connection.Close()
			if err != nil {
				utils.LogError(&utils.DefualtPrinter{}, err, "")
			}
		}()
		mainBrewrey.Temperatures[brewery.SensorName(therm.Name)] = thermomerter
	}

	heater, connection := makeSwitchClient(brewSettings.Heater.Addr, brewSettings.Heater.Port)
	defer func() {
		err := connection.Close()
		if err != nil {
			utils.LogError(&utils.DefualtPrinter{}, err, "")
		}
	}()
	mainBrewrey.Element = heater

	startBrewery(brewSettings.Config.Port, &mainBrewrey)
}
