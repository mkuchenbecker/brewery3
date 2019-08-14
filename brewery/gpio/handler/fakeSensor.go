package handler

import (
	"fmt"
	"sync"

	"github.com/mkuchenbecker/brewery3/brewery/gpio"
)

type fakeSensor struct {
	sensors map[string]float64
	err     error
	mux     sync.RWMutex
}

func newFakeSensor() *fakeSensor {
	return &fakeSensor{
		sensors: make(map[string]float64),
	}
}

func (s *fakeSensor) Sensors() ([]string, error) {
	return []string{}, fmt.Errorf("unimplemented")
}

func (s *fakeSensor) Temperature(sensor gpio.TemperatureAddress) (float64, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	val, ok := s.sensors[string(sensor)]
	if !ok {
		return 0, fmt.Errorf("not found")
	}
	return val, s.err
}
