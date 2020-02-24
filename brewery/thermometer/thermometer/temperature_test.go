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

	mockTemp := mocks.NewMockTemperature(mockCtrl)
	mockTemp.EXPECT().Temperature(gpio.TemperatureAddress("address123")).Return(float64(50), nil).Times(1)

	server := ThermometerServer{temperature: mockTemp, address: gpio.TemperatureAddress("address123")}
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

	mockTemp := mocks.NewMockTemperature(mockCtrl)
	mockTemp.EXPECT().Temperature(gpio.TemperatureAddress("address123")).Return(float64(0), fmt.Errorf("temperatureError")).Times(1)

	server := ThermometerServer{temperature: mockTemp, address: gpio.TemperatureAddress("address123")}
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

	mockTemp := mocks.NewMockTemperature(mockCtrl)
	// The expected number of calls is 2, once in the background and once in the foreground.
	mockTemp.EXPECT().Temperature(gpio.TemperatureAddress("address123")).Return(float64(50), nil).Times(1)
	mockTemp.EXPECT().Temperature(gpio.TemperatureAddress("address123")).Return(float64(60), nil).Times(1)

	server, err := NewThermometerServer(mockTemp, gpio.TemperatureAddress("address123"), 2.5)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond) // Sleep so the background process has time to fire.

	res, err := server.Get(context.Background(), &model.GetRequest{})
	assert.NoError(t, err)
	assert.Equal(t, &model.GetResponse{Temperature: 62.5}, res) // Equal to the second call temperature and not the first.
}
