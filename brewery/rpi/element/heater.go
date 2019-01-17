package element

import (
	"context"
	"fmt"
	"time"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	gpio "github.com/mkuchenbecker/brewery3/brewery/rpi/gpio/igpio"
)

// HeaterServer implements switch.
type HeaterServer struct {
	ctrl gpio.Controller
	pin  int
}

func (s *HeaterServer) On(ctx context.Context, req *model.OnRequest) (*model.OnResponse, error) {
	fmt.Printf("On: %s - %s\n", "req.ID", time.Now().String())
	err := s.ctrl.PowerPin(s.pin, true)
	return &model.OnResponse{}, err
}

func (s *HeaterServer) Off(ctx context.Context, req *model.OffRequest) (*model.OffResponse, error) {
	fmt.Printf("Off: %s - %s\n", "req.ID", time.Now().String())
	err := s.ctrl.PowerPin(s.pin, false)
	return &model.OffResponse{}, err
}
