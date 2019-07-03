//+build !test

package sensors

// import (
// 	"fmt"
// 	"log"
// 	"net"

// 	"github.com/mkuchenbecker/brewery3/brewery/gpio"
// 	"github.com/mkuchenbecker/brewery3/brewery/gpio/integration"
// 	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
// 	"github.com/mkuchenbecker/brewery3/brewery/utils"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/reflection"
// )

// func startThermometer(port int, address string) { // nolint: deadcode
// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v\n", err)
// 	}
// 	serve := grpc.NewServer()
// 	addr := gpio.TemperatureAddress(address)
// 	err = integration.VerifyTemperatureAddress(addr)
// 	if err != nil {
// 		log.Fatalf("failed to read address: %v\n", err)
// 	}
// 	server, err := NewThermometerServer(integration.NewDefaultController(), addr)
// 	if err != nil {
// 		log.Fatalf("failed to make thermometer server: %v\n", err)
// 	}

// 	model.RegisterThermometerServer(serve, server)

// 	utils.Print("Serving Traffic")

// 	reflection.Register(serve)
// 	if err := serve.Serve(lis); err != nil {
// 		log.Fatalf("failed to serve: %v\n", err)
// 	}
// }
