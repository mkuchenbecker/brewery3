package integration

import (
	"github.com/mkuchenbecker/brewery3/brewery/gpio"
	rpio "github.com/stianeikeland/go-rpio"
)

type defaultPins struct {
}

func (d *defaultPins) Open() error {
	return rpio.Open()
}

func (d *defaultPins) Close() error {
	return rpio.Close()
}

func (d *defaultPins) Pin(pin uint8) gpio.GPIOPin {
	return rpio.Pin(pin)
}
