package igpio

//go:generate mockgen -destination=../mocks/mock_gpio.go github.com/mkuchenbecker/brewery3/brewery/rpi/gpio/igpio IGpio,Controller,PinGpio

type TemperatureAddress string

type Controller interface {
	PowerPin(pin int, on bool) error
	ReadTemperature(address TemperatureAddress) (float64, error)
}

type IGpio interface {
	Open() error
	Close() error //more needed, check retuen types
	Pin(uint8) PinGpio
}

type PinGpio interface {
	High()
	Low()
}

type SensorArray interface {
	Sensors() ([]Sensor, error)
	Temperature(Sensor) (Celsius, error)
}

type Sensor string
type Celsius float64
