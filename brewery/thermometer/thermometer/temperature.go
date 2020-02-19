package sensors

import (
	"context"
	"sync"
	"time"

	"github.com/mkuchenbecker/brewery3/brewery/utils"

	"github.com/mkuchenbecker/brewery3/brewery/gpio"
	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
)

// ThermometerServer implements thermometer service interface.
type ThermometerServer struct {
	controller         gpio.Controller
	address            gpio.TemperatureAddress
	mux                sync.Mutex // Ensures multiple reads are not simultaneous.
	currentTemp        float64
	err                error
	logIntervalSeconds time.Duration
}

// NewThermometerServer creates a new Thermometer Server.
func NewThermometerServer(controller gpio.Controller, address gpio.TemperatureAddress) (*ThermometerServer, error) {
	s := ThermometerServer{controller: controller, address: address, currentTemp: 0, err: nil, logIntervalSeconds: 5}
	err := s.update()
	go s.backgroundUpdate(utils.UpdateInterval)
	return &s, err
}

func (s *ThermometerServer) backgroundUpdate(interval time.Duration) {
	for {
		var lastLogTime time.Time
		currentTime := time.Now()
		if lastLogTime.Add(s.logIntervalSeconds).Before(currentTime) {
			lastLogTime = currentTime
		}
		err := s.update()
		if err != nil {
			utils.LogError(nil, err, "temperature read error")
		}
		timer := time.NewTimer(interval)
		<-timer.C
	}
}

func (s *ThermometerServer) update() (err error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	var temp float64
	temp, s.err = s.controller.ReadTemperature(s.address)
	if err != nil {
		return err
	}
	s.currentTemp = temp
	return s.err
}

// Get implements Thermometer.Get function.
func (s *ThermometerServer) Get(ctx context.Context, req *model.GetRequest) (*model.GetResponse, error) {
	return &model.GetResponse{Temperature: s.currentTemp}, s.err
}
