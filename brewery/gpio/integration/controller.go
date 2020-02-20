package integration

import (
	"github.com/mkuchenbecker/brewery3/brewery/gpio"
	"github.com/mkuchenbecker/brewery3/brewery/gpio/handler"
)

// NewDefaultController creates a new default controller that interfaces with the
// Raspberry Pi Hardware.
func NewDefaultController() gpio.Controller {
	return handler.NewController(&defaulttemperature{}, &defaultPins{})
}
