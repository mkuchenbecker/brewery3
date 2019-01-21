package gpio

import (
	"fmt"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"

	"github.com/mkuchenbecker/brewery3/brewery/rpi/gpio/igpio"
	mocks "github.com/mkuchenbecker/brewery3/brewery/rpi/gpio/mocks"
)

type fakeSensor struct {
	sensors map[igpio.Sensor]igpio.Celsius
	err     error
	mux     sync.RWMutex
}

func newFakeSensor() *fakeSensor {
	return &fakeSensor{
		sensors: make(map[igpio.Sensor]igpio.Celsius),
	}
}

func (s *fakeSensor) Sensors() ([]igpio.Sensor, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	retSlice := make([]igpio.Sensor, 0)
	for k := range s.sensors {
		retSlice = append(retSlice, k)
	}
	return retSlice, s.err

}

func (s *fakeSensor) Temperature(sensor igpio.Sensor) (igpio.Celsius, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	val, ok := s.sensors[sensor]
	if !ok {
		return igpio.Celsius(0), fmt.Errorf("not found")
	}
	return val, s.err
}

func TestControllerReadTemperature(t *testing.T) {
	t.Parallel()

	sensor := newFakeSensor()
	sensor.sensors[igpio.Sensor("A")] = 50
	sensor.sensors[igpio.Sensor("B")] = 40

	ctrl := GPIOController{Sensors: sensor}

	cel, err := ctrl.ReadTemperature(igpio.TemperatureAddress("A"))

	assert.NoError(t, err)
	assert.Equal(t, float64(50), cel)
}

func TestControllerPowerPinHigh(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPin := mocks.NewMockPinGpio(mockCtrl)
	mockPin.EXPECT().High().Times(1)

	mockPins := mocks.NewMockIGpio(mockCtrl)
	mockPins.EXPECT().Open().Return(nil).Times(1)
	mockPins.EXPECT().Close().Return(nil).Times(1)
	mockPins.EXPECT().Pin(uint8(5)).Return(mockPin).Times(1)

	ctrl := GPIOController{gpioPins: mockPins}
	assert.Nil(t, ctrl.PowerPin(5, true))
}

func TestControllerPowerPinLow(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPin := mocks.NewMockPinGpio(mockCtrl)
	mockPin.EXPECT().Low().Times(1)

	mockPins := mocks.NewMockIGpio(mockCtrl)
	mockPins.EXPECT().Open().Return(nil).Times(1)
	mockPins.EXPECT().Close().Return(nil).Times(1)
	mockPins.EXPECT().Pin(uint8(5)).Return(mockPin).Times(1)

	ctrl := GPIOController{gpioPins: mockPins}
	assert.Nil(t, ctrl.PowerPin(5, false))
}

func TestControllerPowerPinOpenError(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	err := fmt.Errorf("error")
	mockPins := mocks.NewMockIGpio(mockCtrl)
	mockPins.EXPECT().Open().Return(err).Times(1)

	ctrl := GPIOController{gpioPins: mockPins}
	assert.Equal(t, err, ctrl.PowerPin(5, false))
}

func TestControllerPowerPinCloseError(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPin := mocks.NewMockPinGpio(mockCtrl)
	mockPin.EXPECT().Low().Times(1)

	err := fmt.Errorf("error")
	mockPins := mocks.NewMockIGpio(mockCtrl)
	mockPins.EXPECT().Open().Return(nil).Times(1)
	mockPins.EXPECT().Close().Return(err).Times(1)
	mockPins.EXPECT().Pin(uint8(5)).Return(mockPin).Times(1)

	ctrl := GPIOController{gpioPins: mockPins}
	assert.Equal(t, err, ctrl.PowerPin(5, false))
}
