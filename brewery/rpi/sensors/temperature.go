package sensors

import (
	"context"

	model "github.com/brewery3/brewery/model/gomodel"
	gpio "github.com/brewery3/brewery/rpi/gpio/igpio"
)

// HeaterServer implements switch.
type ThermometerServer struct {
	ctrl    gpio.Controller
	address gpio.TemperatureAddress
}

func (s *ThermometerServer) Get(ctx context.Context, req *model.GetRequest) (*model.GetResponse, error) {
	temp, err := s.ctrl.ReadTemperature(s.address)
	return &model.GetResponse{Temperature: temp}, err
}
