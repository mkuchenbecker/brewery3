package main

import (
	"context"
	"fmt"
	"time"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"github.com/mkuchenbecker/brewery3/brewery/rpi"
	"github.com/mkuchenbecker/brewery3/brewery/rpi/element"
	"github.com/mkuchenbecker/brewery3/brewery/rpi/sensors"
	"github.com/mkuchenbecker/brewery3/brewery/utils"
	"google.golang.org/grpc"
)

const (
	MashAddr   = "28-0315712c08ff"
	HermsAddr  = "28-0315715039ff"
	BoilAddr   = "28-031571188aff"
	ElementPin = 11
)

func MakeTemperatureClient(port int, address string) (model.ThermometerClient, *grpc.ClientConn) {
	utils.Print(fmt.Sprintf("Starting temperature server on port: %d", port))
	go sensors.StartThermometer(port, address)
	utils.Print(fmt.Sprintf("Waiting for discovery on port: %d", port))
	time.Sleep(2 * time.Second)
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

func MakeSwitchClient(port int, pin uint8) (model.SwitchClient, *grpc.ClientConn) {
	utils.Print(fmt.Sprintf("Starting switch server on port: %d", port))
	go element.StartHeater(port, pin)
	utils.Print(fmt.Sprintf("Waiting for discovery on port: %d", port))
	time.Sleep(5 * time.Second)
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
	mash, mashConn := MakeTemperatureClient(8090, MashAddr)
	defer mashConn.Close()
	herms, hermsConn := MakeTemperatureClient(8091, HermsAddr)
	defer hermsConn.Close()
	boil, boilConn := MakeTemperatureClient(8092, BoilAddr)
	defer boilConn.Close()

	element, elementConn := MakeSwitchClient(8110, ElementPin)
	defer elementConn.Close()

	brewery := rpi.Brewery{
		MashSensor:  mash,
		HermsSensor: herms,
		BoilSensor:  boil,
		Element:     element,
	}
	rpi.StartBrewery(8100, &brewery)
}
