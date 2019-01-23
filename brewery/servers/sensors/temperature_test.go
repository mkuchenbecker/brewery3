package sensors

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mkuchenbecker/brewery3/brewery/gpio"
	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"

	"github.com/golang/mock/gomock"
	mocks "github.com/mkuchenbecker/brewery3/brewery/gpio/mocks"
)

func TestReadTemperature(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockController := mocks.NewMockController(mockCtrl)
	mockController.EXPECT().ReadTemperature(gpio.TemperatureAddress("address123")).Return(float64(50), nil).Times(1)

	server := ThermometerServer{controller: mockController, address: gpio.TemperatureAddress("address123")}
	err := server.update()
	assert.NoError(t, err)

	res, err := server.Get(context.Background(), &model.GetRequest{})
	assert.NoError(t, err)
	assert.Equal(t, &model.GetResponse{Temperature: 50}, res)
}

func TestReadTemperatureError(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockController := mocks.NewMockController(mockCtrl)
	mockController.EXPECT().ReadTemperature(gpio.TemperatureAddress("address123")).Return(float64(0), fmt.Errorf("temperatureError")).Times(1)

	server := ThermometerServer{controller: mockController, address: gpio.TemperatureAddress("address123")}
	err := server.update()
	assert.Error(t, err)

	_, err = server.Get(context.Background(), &model.GetRequest{})
	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("temperatureError"), err)
}

func TestTemperatureConstructor(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockController := mocks.NewMockController(mockCtrl)
	// The expected number of calls is 2, once in the background and once in the foreground.
	mockController.EXPECT().ReadTemperature(gpio.TemperatureAddress("address123")).Return(float64(50), nil).Times(1)
	mockController.EXPECT().ReadTemperature(gpio.TemperatureAddress("address123")).Return(float64(60), nil).Times(1)

	server, err := NewThermometerServer(mockController, gpio.TemperatureAddress("address123"))
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond) // Sleep so the background process has time to fire.

	res, err := server.Get(context.Background(), &model.GetRequest{})
	assert.NoError(t, err)
	assert.Equal(t, &model.GetResponse{Temperature: 60}, res) // Equal to the second call temperature and not the first.
}
