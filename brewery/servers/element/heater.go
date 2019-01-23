package element

import (
	"context"
	"time"

	"github.com/mkuchenbecker/brewery3/brewery/utils"

	"github.com/mkuchenbecker/brewery3/brewery/gpio"
	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
)

// HeaterServer implements switch.
type HeaterServer struct {
	controller gpio.Controller
	Pin        uint8
	offError   error
}

// NewHeaterServer constructs a HeaterServer from the supplied parameters.
// ctrl is the GPIOController and pin is the GPIO pin for the server to use.
func NewHeaterServer(ctrl gpio.Controller, pin uint8) *HeaterServer {
	return &HeaterServer{controller: ctrl, Pin: pin}
}

// On handles the Switch.On function.
func (s *HeaterServer) On(ctx context.Context, req *model.OnRequest) (*model.OnResponse, error) {
	utils.Print("Heater On")
	err := s.controller.PowerPin(s.Pin, true)
	return &model.OnResponse{}, err
}

// Off handles the Switch.Off function.
func (s *HeaterServer) Off(ctx context.Context, req *model.OffRequest) (*model.OffResponse, error) {
	utils.Print("Heater Off")
	err := s.controller.PowerPin(s.Pin, false)
	return &model.OffResponse{}, err
}

// ToggleOn handles the Switch.ToggleOn function. It turns the switch on and then off
// after a period of time.
func (s *HeaterServer) ToggleOn(ctx context.Context, req *model.ToggleOnRequest) (*model.ToggleOnResponse, error) {
	utils.Print("Heater On")

	go utils.BackgroundErrReturn(func() error {
		timer := time.NewTimer(time.Duration(req.IntervalMs) * time.Millisecond)
		<-timer.C
		defer utils.Print("Heater Off")
		s.offError = s.controller.PowerPin(s.Pin, false)
		return s.offError
	})
	err := s.controller.PowerPin(s.Pin, true)
	if err != nil {
		return &model.ToggleOnResponse{}, err
	}

	return &model.ToggleOnResponse{}, err
}
