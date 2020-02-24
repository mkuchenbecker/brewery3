package handler

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"

	"github.com/mkuchenbecker/brewery3/brewery/gpio"
	mocks "github.com/mkuchenbecker/brewery3/brewery/gpio/mocks"
	"github.com/pkg/errors"
)

func TestControllerReadTemperature(t *testing.T) {
	t.Parallel()

	sensor := newFakeSensor()
	sensor.sensors["A"] = 50
	sensor.sensors["B"] = 40

	ctrl := NewController(sensor, nil)

	cel, err := ctrl.ReadTemperature(gpio.TemperatureAddress("A"))

	assert.NoError(t, err)
	assert.Equal(t, float64(50), cel)
}

func TestControllerPowerPinHigh(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPin := mocks.NewMockGPIOPin(mockCtrl)
	mockPin.EXPECT().Output().Times(1)
	mockPin.EXPECT().High().Times(1)

	mockPins := mocks.NewMockIGPIO(mockCtrl)
	mockPins.EXPECT().Open().Return(nil).Times(1)
	mockPins.EXPECT().Close().Return(nil).Times(1)
	mockPins.EXPECT().Pin(uint8(5)).Return(mockPin).Times(1)

	ctrl := NewController(nil, mockPins)
	assert.Nil(t, ctrl.PowerPin(5, true))
}

func TestControllerPowerPinLow(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPin := mocks.NewMockGPIOPin(mockCtrl)
	mockPin.EXPECT().Output().Times(1)
	mockPin.EXPECT().Low().Times(1)

	mockPins := mocks.NewMockIGPIO(mockCtrl)
	mockPins.EXPECT().Open().Return(nil).Times(1)
	mockPins.EXPECT().Close().Return(nil).Times(1)
	mockPins.EXPECT().Pin(uint8(5)).Return(mockPin).Times(1)

	ctrl := NewController(nil, mockPins)
	assert.Nil(t, ctrl.PowerPin(5, false))
}

func TestControllerPowerPinOpenError(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	err := fmt.Errorf("error")
	mockPins := mocks.NewMockIGPIO(mockCtrl)
	mockPins.EXPECT().Open().Return(err).Times(1)

	ctrl := NewController(nil, mockPins)
	assert.Equal(t, err, errors.Cause(ctrl.PowerPin(5, false)))
}
