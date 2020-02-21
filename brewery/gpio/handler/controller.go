package handler

import (
	"github.com/mkuchenbecker/brewery3/brewery/gpio"
	"github.com/mkuchenbecker/brewery3/brewery/utils"
	"github.com/pkg/errors"
)

// GPIOController is an implementation of the Controller interface.
type gpioController struct {
	Sensors  gpio.GPIOTemperature
	gpioPins gpio.IGPIO
}

// NewController constructs a new Controller.
func NewController(sensors gpio.GPIOTemperature, gpioPins gpio.IGPIO) gpio.Controller {
	return &gpioController{Sensors: sensors, gpioPins: gpioPins}
}

func (ctrl *gpioController) PowerPin(pinNum uint8, on bool) (err error) {
	err = ctrl.gpioPins.Open()
	if err != nil {
		return errors.Wrap(err, "unable to open gpio")
	}
	defer utils.DeferErrReturn(ctrl.gpioPins.Close, &err)
	pin := ctrl.gpioPins.Pin(pinNum)
	pin.Output()
	if on {
		pin.High()
	} else {
		pin.Low()
	}
	return nil
}

func (ctrl *gpioController) ReadTemperature(sensor gpio.TemperatureAddress) (float64, error) {
	return ctrl.Sensors.Temperature(sensor)
}
