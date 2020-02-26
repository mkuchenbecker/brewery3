package element

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/mkuchenbecker/brewery3/brewery/gpio/integration"
	"github.com/mkuchenbecker/brewery3/brewery/utils"

	"github.com/mkuchenbecker/brewery3/brewery/gpio"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// HeaterServer implements switch.
type HeaterServer struct {
	elementPin gpio.Pin
}

// NewHeaterServer constructs a HeaterServer from the supplied parameters.
// ctrl is the GPIOController and pin is the GPIO pin for the server to use.
func NewHeaterServer(pin gpio.Pin) *HeaterServer {
	return &HeaterServer{elementPin: pin}
}

// On handles the Switch.On function.
func (s *HeaterServer) On(ctx context.Context, req *model.OnRequest) (*model.OnResponse, error) {
	utils.Print("Heater On")
	s.elementPin.High()
	return &model.OnResponse{}, nil
}

// Off handles the Switch.Off function.
func (s *HeaterServer) Off(ctx context.Context, req *model.OffRequest) (*model.OffResponse, error) {
	utils.Print("Heater Off")
	s.elementPin.Low()
	return &model.OffResponse{}, nil
}

func StartElement(pinNum int64, port int64) {
	utils.Printf("Starting heater on port: %d; pin: %d", port, pinNum)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serve := grpc.NewServer()

	utils.Print("setting up pins")
	pins := integration.SysfsPins{}
	if err = pins.Open(); err != nil {
		log.Fatalf("failed to open gpio: %v", err)
	}
	defer pins.Close()

	utils.Print("getting pin")
	pin := pins.Pin(uint(pinNum))
	pin.Output()

	heater := NewHeaterServer(pin)
	model.RegisterSwitchServer(serve, heater)
	// Register reflection service on gRPC server.
	reflection.Register(serve)
	if err := serve.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
