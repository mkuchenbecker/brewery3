package element

import (
	"context"

	"github.com/mkuchenbecker/brewery3/brewery/utils"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"github.com/mkuchenbecker/brewery3/brewery/rpi/gpio/igpio"
	gpio "github.com/mkuchenbecker/brewery3/brewery/rpi/gpio/igpio"
)

// HeaterServer implements switch.
type HeaterServer struct {
	Controller gpio.Controller
	Pin        uint8
}

func NewHeaterServer(ctrl igpio.Controller, pin uint8) *HeaterServer {
	return &HeaterServer{Controller: ctrl, Pin: pin}
}

func (s *HeaterServer) On(ctx context.Context, req *model.OnRequest) (*model.OnResponse, error) {
	utils.Print("Heater On")
	err := s.Controller.PowerPin(s.Pin, true)
	return &model.OnResponse{}, err
}

func (s *HeaterServer) Off(ctx context.Context, req *model.OffRequest) (*model.OffResponse, error) {
	utils.Print("Heater Off")
	err := s.Controller.PowerPin(s.Pin, false)
	return &model.OffResponse{}, err
}
