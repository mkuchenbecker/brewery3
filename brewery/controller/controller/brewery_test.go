package brewery

import (
	"context"
	"fmt"
	"testing"

	"github.com/mkuchenbecker/brewery3/brewery/logger"

	"github.com/stretchr/testify/assert"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"

	"github.com/golang/mock/gomock"
	mocks "github.com/mkuchenbecker/brewery3/brewery/model/gomock"
)

var MashConfig = &model.ControlScheme{Scheme: &model.ControlScheme_Mash_{
	Mash: &model.ControlScheme_Mash{
		BoilMinTemp:  50,
		BoilMaxTemp:  100,
		HermsMinTemp: 51,
		HermsMaxTemp: 60,
		MashMinTemp:  49.5,
		MashMaxTemp:  50.5,
	}}}

type mockBrewery struct {
	*Brewery
}

func newMockBrewery(mockCtrl *gomock.Controller) *mockBrewery {
	mb := mockBrewery{Brewery: &Brewery{
		HermsSensor: mocks.NewMockThermometerClient(mockCtrl),
		MashSensor:  mocks.NewMockThermometerClient(mockCtrl),
		BoilSensor:  mocks.NewMockThermometerClient(mockCtrl),
		Element:     mocks.NewMockSwitchClient(mockCtrl),
		Logger:      logger.NewFake(),
	}}
	return &mb
}

func (b *mockBrewery) SetMash(temp float64) {
	b.MashSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: temp}, nil)
}
func (b *mockBrewery) SetHerms(temp float64) {
	b.HermsSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: temp}, nil)

}
func (b *mockBrewery) SetBoil(temp float64) {
	b.BoilSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: temp}, nil)
}

func TestElementOnError(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSwitch := mocks.NewMockSwitchClient(mockCtrl)
	mockSwitch.EXPECT().On(context.Background(),
		gomock.Any()).Return(&model.OnResponse{}, fmt.Errorf("unable to turn coil on")).Times(1)
	brewery := Brewery{Element: mockSwitch}
	err := brewery.elementOn(context.Background())
	assert.Error(t, err)

}

func TestBreweryRunMash(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := newMockBrewery(mockCtrl)
	_, err := brewery.Control(context.Background(), &model.ControlRequest{Scheme: MashConfig})
	assert.NoError(t, err)

	brewery.SetBoil(60)
	brewery.SetHerms(55)
	brewery.SetMash(50.25)
	brewery.Element.(*mocks.MockSwitchClient).EXPECT().On(context.Background(),
		&model.OnRequest{}).Return(&model.OnResponse{}, nil).Times(1)

	brewery.Element.(*mocks.MockSwitchClient).EXPECT().Off(context.Background(),
		&model.OffRequest{}).Return(&model.OffResponse{}, nil).Times(1)

	err = brewery.run()
	assert.NoError(t, err)

	brewery.SetBoil(60)
	brewery.SetHerms(55)
	brewery.SetMash(45)

	err = brewery.run()
	assert.NoError(t, err)
}

func TestBreweryRunBoil(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := newMockBrewery(mockCtrl)
	_, err := brewery.Control(context.Background(), &model.ControlRequest{Scheme: &model.ControlScheme{Scheme: &model.ControlScheme_Boil_{}}})
	assert.NoError(t, err)

	brewery.Element.(*mocks.MockSwitchClient).EXPECT().On(context.Background(),
		&model.OnRequest{}).Return(&model.OnResponse{}, nil).Times(1)

	err = brewery.run()
	assert.NoError(t, err)
}

func TestBreweryRunNoScheme(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	brewery := newMockBrewery(mockCtrl)
	err := brewery.run()
	assert.Nil(t, err)
}

func TestBreweryMashThermOn(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := newMockBrewery(mockCtrl)
	brewery.Scheme = MashConfig

	// Everything good.
	brewery.SetBoil(60)
	brewery.SetHerms(55)
	brewery.SetMash(50.25)
	on, err := brewery.mashThermOn(context.Background())
	assert.NoError(t, err)
	assert.False(t, on)

	brewery.SetBoil(45) // Boil low
	brewery.SetHerms(55)
	brewery.SetMash(50.25)
	on, err = brewery.mashThermOn(context.Background())
	assert.NoError(t, err)
	assert.True(t, on)

	brewery.SetBoil(60)
	brewery.SetHerms(45) // Herms low
	brewery.SetMash(50.25)
	on, err = brewery.mashThermOn(context.Background())
	assert.NoError(t, err)
	assert.True(t, on)

	brewery.SetBoil(60)
	brewery.SetHerms(55)
	brewery.SetMash(30) // Mash low
	on, err = brewery.mashThermOn(context.Background())
	assert.NoError(t, err)
	assert.True(t, on)

	// When the mash is too high, always turn off.
	brewery.SetBoil(45) // Boil low
	brewery.SetHerms(55)
	brewery.SetMash(70) // Mash High
	on, err = brewery.mashThermOn(context.Background())
	assert.NoError(t, err)
	assert.False(t, on)

	// Even though the mash is too low, the herms is too high so turn off.
	// This protects us from overshooting the temperature.
	brewery.SetBoil(60)
	brewery.SetHerms(100)
	brewery.SetMash(30)
	on, err = brewery.mashThermOn(context.Background())
	assert.NoError(t, err)
	assert.False(t, on)

}
