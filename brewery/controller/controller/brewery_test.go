package brewery

import (
	"context"
	"fmt"
	"testing"
	"time"

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

func TestPowerToggle(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSwitch := mocks.NewMockSwitchClient(mockCtrl)

	mockSwitch.EXPECT().ToggleOn(context.Background(),
		&model.ToggleOnRequest{IntervalMs: 75}).Return(&model.ToggleOnResponse{}, nil).Times(1)
	mockSwitch.EXPECT().Off(context.Background(),
		gomock.Any()).Return(&model.OffResponse{}, nil).Times(1) // Cleanup action for toggle.

	brewery := Brewery{Element: mockSwitch}
	err := brewery.powerToggle(
		75, // milliseconds
		time.Duration(100)*time.Millisecond,
		time.Duration(100)*time.Millisecond)
	assert.NoError(t, err)
}

func TestElementPowerLevelToggleMultipleLoops(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSwitch := mocks.NewMockSwitchClient(mockCtrl)
	mockSwitch.EXPECT().ToggleOn(context.Background(),
		&model.ToggleOnRequest{IntervalMs: 5}).Return(&model.ToggleOnResponse{}, nil).Times(10)
	mockSwitch.EXPECT().Off(context.Background(),
		gomock.Any()).Return(&model.OffResponse{}, nil).Times(1)

	brewery := Brewery{Element: mockSwitch}
	err := brewery.powerToggle(
		5,
		time.Duration(10)*time.Millisecond,
		time.Duration(100)*time.Millisecond)
	assert.NoError(t, err)
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

func TestElementPowerLevelToggleError(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSwitch := mocks.NewMockSwitchClient(mockCtrl)

	mockSwitch.EXPECT().ToggleOn(context.Background(),
		&model.ToggleOnRequest{IntervalMs: 9}).Return(&model.ToggleOnResponse{},
		fmt.Errorf("unable to turn coil on")).Times(1)

	mockSwitch.EXPECT().Off(context.Background(),
		gomock.Any()).Return(&model.OffResponse{}, nil).Times(1)
	brewery := Brewery{Element: mockSwitch}
	err := brewery.powerToggle(9, 10*time.Millisecond, 10*time.Millisecond)
	assert.Error(t, err)

	mockSwitch.EXPECT().ToggleOn(context.Background(),
		&model.ToggleOnRequest{IntervalMs: 5}).Return(&model.ToggleOnResponse{},
		nil).Times(1)
	// Happens in multiples of 3 due to retries.
	mockSwitch.EXPECT().Off(context.Background(),
		gomock.Any()).Return(&model.OffResponse{}, fmt.Errorf("unable to turn coil off")).Times(3)

	err = brewery.powerToggle(5, 10*time.Millisecond, 10*time.Millisecond)
	assert.Error(t, err)
}

func TestElementPowerLevel0(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSwitch := mocks.NewMockSwitchClient(mockCtrl)

	mockSwitch.EXPECT().Off(context.Background(), gomock.Any()).Return(&model.OffResponse{}, nil).Times(1)
	brewery := Brewery{Element: mockSwitch}
	err := brewery.powerLevel(0)
	assert.NoError(t, err)
}

func TestElementPowerLevel101(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSwitch := mocks.NewMockSwitchClient(mockCtrl)

	mockSwitch.EXPECT().Off(context.Background(), gomock.Any()).Return(&model.OffResponse{}, nil).Times(1)
	brewery := Brewery{Element: mockSwitch}
	err := brewery.powerLevel(101)
	assert.NoError(t, err)
}

func TestElementPowerLevel100(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSwitch := mocks.NewMockSwitchClient(mockCtrl)

	mockSwitch.EXPECT().On(context.Background(), gomock.Any()).Return(&model.OnResponse{}, nil).Times(1)
	brewery := Brewery{Element: mockSwitch}
	err := brewery.powerLevel(100)
	assert.NoError(t, err)
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

func TestBreweryRunPower(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := newMockBrewery(mockCtrl)
	_, err := brewery.Control(context.Background(), &model.ControlRequest{Scheme: &model.ControlScheme{Scheme: &model.ControlScheme_Power_{
		Power: &model.ControlScheme_Power{PowerLevel: 100}}}})
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
