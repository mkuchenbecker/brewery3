package gpio

import (
	"github.com/mkuchenbecker/brewery3/brewery/rpi/gpio/igpio"
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

func (d *DefaultPins) Pin(pin uint8) igpio.PinGpio {
	return rpio.Pin(pin)
}
