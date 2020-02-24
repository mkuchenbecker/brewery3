package element

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"

	"github.com/golang/mock/gomock"
	mocks "github.com/mkuchenbecker/brewery3/brewery/gpio/mocks"
)

func TestHeaterOn(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pin := mocks.NewMockPin(mockCtrl)
	pin.EXPECT().High().Times(1)

	server := NewHeaterServer(pin)

	res, err := server.On(context.Background(), &model.OnRequest{})
	assert.NoError(t, err)
	assert.Equal(t, &model.OnResponse{}, res)
}

func TestHeaterOff(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pin := mocks.NewMockPin(mockCtrl)
	pin.EXPECT().Low().Times(1)

	server := NewHeaterServer(pin)

	res, err := server.Off(context.Background(), &model.OffRequest{})
	assert.NoError(t, err)
	assert.Equal(t, &model.OffResponse{}, res)
}
