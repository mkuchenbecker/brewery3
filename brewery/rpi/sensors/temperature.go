package sensors

import (
	"context"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	gpio "github.com/mkuchenbecker/brewery3/brewery/rpi/gpio/igpio"
)

// HeaterServer implements switch.
type ThermometerServer struct {
	Controller gpio.Controller
	Address    gpio.TemperatureAddress
}

func (s *ThermometerServer) Get(ctx context.Context, req *model.GetRequest) (*model.GetResponse, error) {
	temp, err := s.Controller.ReadTemperature(s.Address)
	return &model.GetResponse{Temperature: temp}, err
}
