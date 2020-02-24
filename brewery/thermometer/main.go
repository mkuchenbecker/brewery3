//+build !test

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/mkuchenbecker/brewery3/brewery/gpio"
	"github.com/mkuchenbecker/brewery3/brewery/gpio/integration"
	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	sensors "github.com/mkuchenbecker/brewery3/brewery/thermometer/thermometer"
	"github.com/mkuchenbecker/brewery3/brewery/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func startThermometer(port int, address string, adjustment float64) { // nolint: deadcode
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	utils.Print("Thermometer Starting")
	if err != nil {
		utils.Printf("failed to listen: %v\n", err)
		log.Fatalf("failed to listen: %v\n", err)
	}
	serve := grpc.NewServer()
	addr := gpio.TemperatureAddress(address)
	err = integration.VerifyTemperatureAddress(addr)
	if err != nil {
		utils.LogError(&utils.DefualtPrinter{}, err, "failed to read address")
	}
	server, err := sensors.NewThermometerServer(integration.DefaultTemperature{}, addr, adjustment)
	if err != nil {
		utils.Printf("failed to temperature")
		utils.LogError(&utils.DefualtPrinter{}, err, "failed to temperature")
	}

	model.RegisterThermometerServer(serve, server)

	utils.Print("Serving Traffic")

	reflection.Register(serve)
	if err := serve.Serve(lis); err != nil {
		utils.Printf("failed to serve: %v\n", err)
		log.Fatalf("failed to serve: %v\n", err)
	}
}

func main() {
	strPort := os.Getenv("THERM_PORT")
	address := os.Getenv("THERM_ADDR")
	strAdjustment := os.Getenv("CALIBRATION_ADJUSTMENT")

	utils.Printf("Port:%s; address: %s; adjustment: %s", strPort, address, strAdjustment)

	port, err := strconv.ParseInt(strPort, 10, 32)
	if err != nil {
		log.Fatalf("Invalid port is not 32 bit int: %s", strPort)
	}

	adjustment, err := strconv.ParseFloat(strAdjustment, 64)
	if err != nil {
		log.Fatalf("invalid calibration adjustment: %s", strAdjustment)
	}

	startThermometer(int(port), address, adjustment)
}
