package element

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"

	"github.com/golang/mock/gomock"
	mocks "github.com/mkuchenbecker/brewery3/brewery/rpi/gpio/mocks"
)

func TestHeaterOn(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockController := mocks.NewMockController(mockCtrl)
	mockController.EXPECT().PowerPin(5, true).Return(nil).Times(1)

	server := HeaterServer{ctrl: mockController, pin: 5}

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
	mockController.EXPECT().PowerPin(5, true).Return(expErr).Times(1)

	server := HeaterServer{ctrl: mockController, pin: 5}

	_, err := server.On(context.Background(), &model.OnRequest{})
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
}

func TestHeaterOff(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockController := mocks.NewMockController(mockCtrl)
	mockController.EXPECT().PowerPin(5, false).Return(nil).Times(1)

	server := HeaterServer{ctrl: mockController, pin: 5}

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
	mockController.EXPECT().PowerPin(5, false).Return(expErr).Times(1)

	server := HeaterServer{ctrl: mockController, pin: 5}

	_, err := server.Off(context.Background(), &model.OffRequest{})
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
}
