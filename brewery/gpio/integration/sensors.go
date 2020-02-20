package integration

import (
	"fmt"

	"github.com/mkuchenbecker/brewery3/brewery/gpio"

	temperature "github.com/yryz/ds18b20"
)

type defaulttemperature struct {
}

func (d *defaulttemperature) Sensors() ([]string, error) {
	return temperature.Sensors()
}

func (d *defaulttemperature) Temperature(sensor gpio.TemperatureAddress) (float64, error) {
	return temperature.Temperature(string(sensor))
}

// VerifyTemperatureAddress returns an error if the sensor is not found.
func VerifyTemperatureAddress(address gpio.TemperatureAddress) error {
	sensors, err := temperature.Sensors()
	if err != nil {
		return err
	}

	for _, sensor := range sensors {
		if sensor == string(address) {
			return nil
		}
	}

	return fmt.Errorf("sensor not found %s", address)
}
