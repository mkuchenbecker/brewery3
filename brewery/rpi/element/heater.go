package element

import (
	"context"

	"github.com/mkuchenbecker/brewery3/brewery/utils"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	gpio "github.com/mkuchenbecker/brewery3/brewery/rpi/gpio/igpio"
)

// HeaterServer implements switch.
type HeaterServer struct {
	ctrl gpio.Controller
	pin  uint8
}

func (s *HeaterServer) On(ctx context.Context, req *model.OnRequest) (*model.OnResponse, error) {
	utils.Print("Heater On")
	err := s.ctrl.PowerPin(s.pin, true)
	return &model.OnResponse{}, err
}

func (s *HeaterServer) Off(ctx context.Context, req *model.OffRequest) (*model.OffResponse, error) {
	utils.Print("Heater Off")
	err := s.ctrl.PowerPin(s.pin, false)
	return &model.OffResponse{}, err
}
