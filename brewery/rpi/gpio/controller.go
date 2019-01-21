package gpio

import (
	"github.com/mkuchenbecker/brewery3/brewery/rpi/gpio/igpio"
	"github.com/mkuchenbecker/brewery3/brewery/utils"
)

type GPIOController struct {
	Sensors  igpio.SensorArray
	gpioPins igpio.IGpio
}

func (gpio *GPIOController) PowerPin(pinNum uint8, on bool) (err error) {
	err = gpio.gpioPins.Open()
	if err != nil {
		return err
	}
	defer utils.DeferErrReturn(gpio.gpioPins.Close, &err)
	pin := gpio.gpioPins.Pin(pinNum)
	if on {
		pin.High()
	} else {
		pin.Low()
	}
	return nil
}

func (gp *GPIOController) ReadTemperature(sensor igpio.TemperatureAddress) (float64, error) {
	temp, err := gp.Sensors.Temperature(igpio.Sensor(sensor))
	return float64(temp), err
}

func GetDefaultController() *GPIOController {
	return &GPIOController{Sensors: &DefaultSensorArray{}, gpioPins: &DefaultPins{}}
}
