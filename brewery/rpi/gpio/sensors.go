package gpio

import (
	"fmt"

	"github.com/mkuchenbecker/brewery3/brewery/rpi/gpio/igpio"
	temperature "github.com/yryz/ds18b20"
)

type DefaultSensorArray struct {
}

func (d *DefaultSensorArray) Sensors() ([]igpio.Sensor, error) {
	strSli, err := temperature.Sensors()
	retSli := make([]igpio.Sensor, 0)
	if err != nil {
		return retSli, err
	}
	for _, s := range strSli {
		retSli = append(retSli, igpio.Sensor(s))
	}
	return retSli, nil
}

func (d *DefaultSensorArray) Temperature(sensor igpio.Sensor) (igpio.Celsius, error) {
	temp, err := temperature.Temperature(string(sensor))
	return igpio.Celsius(temp), err
}

func NewTemperatureAddress(address string, sensorArray igpio.SensorArray) (igpio.TemperatureAddress, error) {
	sensors, err := sensorArray.Sensors()
	if err != nil {
		return "", err
	}

	for _, sensor := range sensors {
		if sensor == igpio.Sensor(address) {
			return igpio.TemperatureAddress(address), nil
		}
	}

	return "", fmt.Errorf("sensor not found %s", address)
}
