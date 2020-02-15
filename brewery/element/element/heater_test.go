package element

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"

	"github.com/golang/mock/gomock"
	mocks "github.com/mkuchenbecker/brewery3/brewery/gpio/mocks"
)

func TestHeaterOn(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockController := mocks.NewMockController(mockCtrl)
	mockController.EXPECT().PowerPin(uint8(5), true).Return(nil).Times(1)

	server := NewHeaterServer(mockController, 5)

	res, err := server.On(context.Background(), &model.OnRequest{})
	assert.NoError(t, err)
	assert.Equal(t, &model.OnResponse{}, res)
}

func TestHeaterOnError(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expErr := fmt.Errorf("error")

	mockController := mocks.NewMockController(mockCtrl)
	mockController.EXPECT().PowerPin(uint8(5), true).Return(expErr).Times(1)

	server := NewHeaterServer(mockController, 5)

	_, err := server.On(context.Background(), &model.OnRequest{})
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
}

func TestHeaterOff(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockController := mocks.NewMockController(mockCtrl)
	mockController.EXPECT().PowerPin(uint8(5), false).Return(nil).Times(1)

	server := NewHeaterServer(mockController, 5)

	res, err := server.Off(context.Background(), &model.OffRequest{})
	assert.NoError(t, err)
	assert.Equal(t, &model.OffResponse{}, res)
}

func TestHeaterOffError(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expErr := fmt.Errorf("error")

	mockController := mocks.NewMockController(mockCtrl)
	mockController.EXPECT().PowerPin(uint8(5), false).Return(expErr).Times(1)

	server := NewHeaterServer(mockController, 5)

	_, err := server.Off(context.Background(), &model.OffRequest{})
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
}

func TestHeaterToggle(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockController := mocks.NewMockController(mockCtrl)
	mockController.EXPECT().PowerPin(uint8(5), true).Return(nil).Times(1)
	mockController.EXPECT().PowerPin(uint8(5), false).Return(nil).Times(1)

	server := NewHeaterServer(mockController, 5)

	res, err := server.ToggleOn(context.Background(), &model.ToggleOnRequest{IntervalMs: 1}) // Should fire off after 1ms.
	assert.NoError(t, err)
	assert.Equal(t, &model.ToggleOnResponse{}, res)

	// Wait for the Off to fire.
	time.Sleep(100 * time.Millisecond)
}
