package element

import (
	"context"

	"github.com/mkuchenbecker/brewery3/brewery/utils"

	"github.com/mkuchenbecker/brewery3/brewery/gpio"
	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
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
