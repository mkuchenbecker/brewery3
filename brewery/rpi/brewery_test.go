package rpi

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	model "github.com/brewery3/brewery/model/gomodel"

	mocks "github.com/brewery3/brewery/model/gomock"
	"github.com/golang/mock/gomock"
)

func TestElementPowerLevelToggle(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSwitch := mocks.NewMockSwitchClient(mockCtrl)

	// Base test.
	mockSwitch.EXPECT().On(context.Background(), gomock.Any()).Return(&model.OnResponse{}, nil).Times(1)
	mockSwitch.EXPECT().Off(context.Background(), gomock.Any()).Return(&model.OffResponse{}, nil).Times(1)
	// Interval is 1 second, time on is 0.75s, quit after 1 second.
	// Result is that the switch should turn on, .75s later turn off, and on the next iteration quit.

	brewery := Brewery{element: mockSwitch}
	err := brewery.elementPowerLevelToggle(
		time.Duration(750)*time.Millisecond,
		time.Duration(1)*time.Second,
		time.Duration(1)*time.Second)
	assert.NoError(t, err)
}
func TestElementPowerLevelToggleMultipleLoops(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSwitch := mocks.NewMockSwitchClient(mockCtrl)
	// Tests we always turn off after every
	mockSwitch.EXPECT().On(context.Background(), gomock.Any()).Return(&model.OnResponse{}, nil).Times(10)
	mockSwitch.EXPECT().Off(context.Background(), gomock.Any()).Return(&model.OffResponse{}, nil).Times(10)
	// Interval is 1 second, time on is 0.75s, quit after 1.25 second.
	// Two iterations should occur.
	brewery := Brewery{element: mockSwitch}
	err := brewery.elementPowerLevelToggle(
		time.Duration(50)*time.Millisecond,
		time.Duration(1000)*time.Millisecond,
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
	mockSwitch.EXPECT().Off(context.Background(),
		gomock.Any()).Return(&model.OffResponse{}, nil).Times(1)
	brewery := Brewery{element: mockSwitch}
	err := brewery.ElementOn(1 * time.Millisecond)
	assert.Error(t, err)

	mockSwitch.EXPECT().On(context.Background(),
		gomock.Any()).Return(&model.OnResponse{}, nil).Times(1)
	// We always attempt to turn the coiloff 3 times.
	mockSwitch.EXPECT().Off(context.Background(),
		gomock.Any()).Return(&model.OffResponse{}, fmt.Errorf("unable to turn coil off")).Times(3)
	err = brewery.ElementOn(1 * time.Millisecond)
	assert.Error(t, err)
}

func TestElementPowerLevelToggleError(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSwitch := mocks.NewMockSwitchClient(mockCtrl)
	mockSwitch.EXPECT().On(context.Background(),
		gomock.Any()).Return(&model.OnResponse{}, fmt.Errorf("unable to turn coil on")).Times(1)
	mockSwitch.EXPECT().Off(context.Background(),
		gomock.Any()).Return(&model.OffResponse{}, nil).Times(1)
	brewery := Brewery{element: mockSwitch}
	err := brewery.elementPowerLevelToggle(9*time.Millisecond, 100*time.Millisecond, 10*time.Millisecond)
	assert.Error(t, err)

	mockSwitch.EXPECT().On(context.Background(),
		gomock.Any()).Return(&model.OnResponse{}, nil).Times(1)
	mockSwitch.EXPECT().Off(context.Background(),
		gomock.Any()).Return(&model.OffResponse{}, fmt.Errorf("unable to turn coil off")).Times(3)

	err = brewery.elementPowerLevelToggle(9*time.Millisecond, 100*time.Millisecond, 10*time.Millisecond)
	assert.Error(t, err)

}

func TestElementPowerLevel0(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSwitch := mocks.NewMockSwitchClient(mockCtrl)

	mockSwitch.EXPECT().Off(context.Background(), gomock.Any()).Return(&model.OffResponse{}, nil).Times(1)
	brewery := Brewery{element: mockSwitch}
	err := brewery.ElementPowerLevel(0, 1)
	assert.NoError(t, err)
}

func TestElementPowerLevel101(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSwitch := mocks.NewMockSwitchClient(mockCtrl)

	mockSwitch.EXPECT().Off(context.Background(), gomock.Any()).Return(&model.OffResponse{}, nil).Times(1)
	brewery := Brewery{element: mockSwitch}
	err := brewery.ElementPowerLevel(101, 1)
	assert.NoError(t, err)
}

func TestElementPowerLevel100(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSwitch := mocks.NewMockSwitchClient(mockCtrl)

	mockSwitch.EXPECT().Off(context.Background(), gomock.Any()).Return(&model.OffResponse{}, nil).Times(1)
	mockSwitch.EXPECT().On(context.Background(), gomock.Any()).Return(&model.OnResponse{}, nil).Times(1)
	brewery := Brewery{element: mockSwitch}
	err := brewery.ElementPowerLevel(100, 1)
	assert.NoError(t, err)
}

func boilScheme() *model.ControlScheme {
	cs := model.ControlScheme{}
	cs.Scheme = &model.ControlScheme_Boil_{}
	return &cs
}

func NewMockBrewery(mockCtrl *gomock.Controller) *Brewery {
	return &Brewery{
		hermsSensor: mocks.NewMockThermometerClient(mockCtrl),
		mashSensor:  mocks.NewMockThermometerClient(mockCtrl),
		boilSensor:  mocks.NewMockThermometerClient(mockCtrl),
		element:     mocks.NewMockSwitchClient(mockCtrl),
	}
}

func TestBreweryGetConstraintsAndMashThermOn(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := NewMockBrewery(mockCtrl)
	brewery.ReplaceConfig(&model.ControlScheme{Scheme: &model.ControlScheme_Mash_{
		Mash: &model.ControlScheme_Mash{
			BoilMinTemp:  50,
			BoilMaxTemp:  100,
			HermsMinTemp: 51,
			HermsMaxTemp: 60,
			MashMinTemp:  49.5,
			MashMaxTemp:  50.5,
		}}})

	brewery.boilSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(60)}, nil).Times(1)

	brewery.hermsSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(55)}, nil).Times(1)

	brewery.mashSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(50.25)}, nil).Times(1)

	constraints, err := brewery.getTempConstraints()
	assert.NoError(t, err)

	assert.Equal(t, []Constraint{{min: 50, max: 100, val: 60},
		{min: 51, max: 60, val: 55},
		{min: 49.5, max: 50.5, val: 50.25}}, constraints)

	brewery.boilSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(60)}, nil).Times(1)

	brewery.hermsSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(55)}, nil).Times(1)

	brewery.mashSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(50.25)}, nil).Times(1)

	on, err := brewery.mashThermOn()
	assert.NoError(t, err)
	assert.False(t, on)
}

func TestBreweryRunMash(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := NewMockBrewery(mockCtrl)
	brewery.ReplaceConfig(&model.ControlScheme{Scheme: &model.ControlScheme_Mash_{
		Mash: &model.ControlScheme_Mash{
			BoilMinTemp:  50,
			BoilMaxTemp:  100,
			HermsMinTemp: 51,
			HermsMaxTemp: 60,
			MashMinTemp:  49.5,
			MashMaxTemp:  50.5,
		}}})

	brewery.boilSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(60)}, nil).Times(1)

	brewery.hermsSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(55)}, nil).Times(1)

	brewery.mashSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(50.25)}, nil).Times(1)

	brewery.element.(*mocks.MockSwitchClient).EXPECT().Off(context.Background(),
		&model.OffRequest{}).Return(&model.OffResponse{}, nil).Times(1)

	err := brewery.Run(1)
	assert.NoError(t, err)

	brewery.boilSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(60)}, nil).Times(1)
	brewery.hermsSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(55)}, nil).Times(1)
	brewery.mashSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(45)}, nil).Times(1)
	brewery.element.(*mocks.MockSwitchClient).EXPECT().On(context.Background(),
		&model.OnRequest{}).Return(&model.OnResponse{}, nil).Times(1)
	brewery.element.(*mocks.MockSwitchClient).EXPECT().Off(context.Background(),
		&model.OffRequest{}).Return(&model.OffResponse{}, nil).Times(1)

	err = brewery.Run(1)
	assert.NoError(t, err)
}

func TestBreweryRunBoil(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := NewMockBrewery(mockCtrl)
	brewery.ReplaceConfig(&model.ControlScheme{Scheme: &model.ControlScheme_Boil_{
		Boil: &model.ControlScheme_Boil{}}})

	brewery.element.(*mocks.MockSwitchClient).EXPECT().On(context.Background(),
		&model.OnRequest{}).Return(&model.OnResponse{}, nil).Times(1)
	brewery.element.(*mocks.MockSwitchClient).EXPECT().Off(context.Background(),
		&model.OffRequest{}).Return(&model.OffResponse{}, nil).Times(1)

	err := brewery.Run(1)
	assert.NoError(t, err)
}

func TestBreweryPower(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := NewMockBrewery(mockCtrl)
	brewery.ReplaceConfig(&model.ControlScheme{Scheme: &model.ControlScheme_Power_{
		Power: &model.ControlScheme_Power{PowerLevel: 100}}})

	brewery.element.(*mocks.MockSwitchClient).EXPECT().On(context.Background(),
		&model.OnRequest{}).Return(&model.OnResponse{}, nil).Times(1)
	brewery.element.(*mocks.MockSwitchClient).EXPECT().Off(context.Background(),
		&model.OffRequest{}).Return(&model.OffResponse{}, nil).Times(1)

	err := brewery.Run(1)
	assert.NoError(t, err)
}
func TestBreweryGetConstraintsError(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := NewMockBrewery(mockCtrl)

	brewery.boilSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(60)}, fmt.Errorf("error boil")).Times(1)

	_, err := brewery.getTempConstraints()
	assert.Equal(t, fmt.Errorf("error boil"), err)

	brewery.boilSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(60)}, nil).Times(1)

	brewery.hermsSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(55)}, fmt.Errorf("error herms")).Times(1)

	_, err = brewery.getTempConstraints()
	assert.Equal(t, fmt.Errorf("error herms"), err)

	brewery.boilSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(60)}, nil).Times(1)

	brewery.hermsSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(55)}, nil).Times(1)

	brewery.mashSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(50.25)}, fmt.Errorf("error mash")).Times(1)

	_, err = brewery.getTempConstraints()
	assert.Equal(t, fmt.Errorf("error mash"), err)

}

func TestConstraints(t *testing.T) {
	t.Parallel()

	// Everything good.
	assert.Equal(t, 0, checkTempConstraints([]Constraint{
		{min: 50, max: 100, val: 60},
		{min: 51, max: 60, val: 55},
		{min: 49.5, max: 50.5, val: 50.25},
	}))

	assert.Equal(t, -1, checkTempConstraints([]Constraint{
		{min: 50, max: 100, val: 45}, //Boil too low
		{min: 51, max: 60, val: 55},
		{min: 49.5, max: 50.5, val: 50.25},
	}))

	assert.Equal(t, 1, checkTempConstraints([]Constraint{
		{min: 50, max: 100, val: 101}, //Boil too high
		{min: 51, max: 60, val: 55},
		{min: 49.5, max: 50.5, val: 50.25},
	}))

	assert.Equal(t, -1, checkTempConstraints([]Constraint{
		{min: 50, max: 100, val: 60},
		{min: 51, max: 60, val: 40}, // Stream too low
		{min: 49.5, max: 50.5, val: 50.25},
	}))

	assert.Equal(t, 1, checkTempConstraints([]Constraint{
		{min: 50, max: 100, val: 60},
		{min: 51, max: 60, val: 61}, // Stream too high
		{min: 49.5, max: 50.5, val: 50.25},
	}))

	assert.Equal(t, -1, checkTempConstraints([]Constraint{
		{min: 50, max: 100, val: 60},
		{min: 51, max: 60, val: 55},
		{min: 49.5, max: 50.5, val: 40}, // Mash too low
	}))

	assert.Equal(t, 1, checkTempConstraints([]Constraint{
		{min: 50, max: 100, val: 60},
		{min: 51, max: 60, val: 55},
		{min: 49.5, max: 50.5, val: 52}, // Mash too high
	}))

	assert.Equal(t, 1, checkTempConstraints([]Constraint{
		{min: 50, max: 100, val: 110},
		{min: 51, max: 60, val: 50},     //Stream too low
		{min: 49.5, max: 50.5, val: 51}, // Mash too high
	}))

	assert.Equal(t, 1, checkTempConstraints([]Constraint{
		{min: 50, max: 100, val: 60},
		{min: 51, max: 60, val: 65},     // Stream too high
		{min: 49.5, max: 50.5, val: 45}, // Mash too low
	}))

}
