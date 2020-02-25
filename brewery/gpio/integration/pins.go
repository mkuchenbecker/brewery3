package integration

import (
	"github.com/mkuchenbecker/brewery3/brewery/gpio"
	rpio "github.com/stianeikeland/go-rpio"
)

type DefaultPins struct {
}

func (d *DefaultPins) Open() error {
	return rpio.Open()
}

func (d *DefaultPins) Close() error {
	return rpio.Close()
}

func (d *DefaultPins) Pin(pin uint8) gpio.Pin {
	return rpio.Pin(pin)
}
