package cli

import (
	"context"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	mock "github.com/mkuchenbecker/brewery3/brewery/model/gomock"
	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"github.com/stretchr/testify/assert"
)

func TestParseTemp(t *testing.T) {
	t.Parallel()
	var temp float64
	var err error

	temp, err = parseTemp("100")
	assert.Equal(t, float64(100), temp)
	assert.NoError(t, err)

	temp, err = parseTemp("0")
	assert.Equal(t, float64(0), temp)
	assert.NoError(t, err)

	temp, err = parseTemp("32f")
	assert.Equal(t, float64(0), temp)
	assert.NoError(t, err)

	_, err = parseTemp("abc")
	assert.Error(t, err)
}

func TestGetMashRequest(t *testing.T) {
	t.Parallel()
	req := getMashRequest(50)
	mash := req.Scheme.GetMash()

	assert.Equal(t, &model.ControlScheme_Mash{
		MashMaxTemp:  50.5,
		MashMinTemp:  50,
		HermsMaxTemp: 65,
		HermsMinTemp: 50,
		BoilMinTemp:  50,
		BoilMaxTemp:  100,
	}, mash)
}

func TestGetPowerRequest(t *testing.T) {
	t.Parallel()
	req := getPowerRequst(50)
	pow := req.Scheme.GetPower()

	assert.Equal(t, &model.ControlScheme_Power{
		PowerLevel: 50,
	}, pow)
}

func TestGetBoilRequest(t *testing.T) {
	t.Parallel()
	req := getBoilRequest()
	pow := req.Scheme.GetPower()

	assert.Equal(t, &model.ControlScheme_Power{
		PowerLevel: 75,
	}, pow)
}

func TestCLIPower(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := mock.NewMockBreweryClient(mockCtrl)

	req := getPowerRequst(55)
	brewery.EXPECT().Control(context.Background(),
		req).Return(&model.ControlResponse{}, nil).Times(1)

	args := os.Args[0:1]
	args = append(args, "-power=55")
	assert.NoError(t, Run(brewery, args))
}

func TestCLIPowerErr(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := mock.NewMockBreweryClient(mockCtrl)

	args := os.Args[0:1]
	args = append(args, "-power=101")
	assert.Error(t, Run(brewery, args))
}

func TestCLIMash(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := mock.NewMockBreweryClient(mockCtrl)

	req := getMashRequest(65)
	brewery.EXPECT().Control(context.Background(),
		req).Return(&model.ControlResponse{}, nil).Times(1)

	args := os.Args[0:1]
	args = append(args, "-mash=65")
	assert.NoError(t, Run(brewery, args))
}

func TestCLIMashErr(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := mock.NewMockBreweryClient(mockCtrl)

	args := os.Args[0:1]
	args = append(args, "-mash=110")
	assert.Error(t, Run(brewery, args))
}

func TestCLIBoil(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := mock.NewMockBreweryClient(mockCtrl)

	req := getBoilRequest()
	brewery.EXPECT().Control(context.Background(),
		req).Return(&model.ControlResponse{}, nil).Times(1)

	args := os.Args[0:1]
	args = append(args, "-boil")
	assert.NoError(t, Run(brewery, args))
}
