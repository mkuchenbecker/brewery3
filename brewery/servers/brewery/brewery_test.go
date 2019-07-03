package brewery

import (
	"context"
	"fmt"
	"testing"
	"time"

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
	err := brewery.elementOn()
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

func NewMockBrewery(mockCtrl *gomock.Controller) *Brewery {
	return &Brewery{
		HermsSensor: mocks.NewMockThermometerClient(mockCtrl),
		MashSensor:  mocks.NewMockThermometerClient(mockCtrl),
		BoilSensor:  mocks.NewMockThermometerClient(mockCtrl),
		Element:     mocks.NewMockSwitchClient(mockCtrl),
	}
}
func TestBreweryGetConstraints(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := NewMockBrewery(mockCtrl)
	_, err := brewery.Control(context.Background(), &model.ControlRequest{Scheme: MashConfig})
	assert.NoError(t, err)

	brewery.BoilSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(60)}, nil).Times(1)

	brewery.HermsSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(55)}, nil).Times(1)

	brewery.MashSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(50.25)}, nil).Times(1)

	constraints, err := brewery.getTempConstraints()
	assert.NoError(t, err)
	assert.Equal(t, []constraint{{min: 50, max: 100, val: 60},
		{min: 51, max: 60, val: 55},
		{min: 49.5, max: 50.5, val: 50.25}}, constraints)
}
func TestBreweryMashThermOn(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := NewMockBrewery(mockCtrl)
	_, err := brewery.Control(context.Background(), &model.ControlRequest{Scheme: MashConfig})
	assert.NoError(t, err)

	brewery.BoilSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(60)}, nil).Times(1)

	brewery.HermsSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(55)}, nil).Times(1)

	brewery.MashSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
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
	_, err := brewery.Control(context.Background(), &model.ControlRequest{Scheme: MashConfig})
	assert.NoError(t, err)

	brewery.BoilSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(60)}, nil).Times(1)

	brewery.HermsSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(55)}, nil).Times(1)

	brewery.MashSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(50.25)}, nil).Times(1)

	brewery.Element.(*mocks.MockSwitchClient).EXPECT().Off(context.Background(),
		&model.OffRequest{}).Return(&model.OffResponse{}, nil).Times(1)

	err = brewery.run()
	assert.NoError(t, err)

	brewery.BoilSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(60)}, nil).Times(1)
	brewery.HermsSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(55)}, nil).Times(1)
	brewery.MashSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(45)}, nil).Times(1)
	brewery.Element.(*mocks.MockSwitchClient).EXPECT().On(context.Background(),
		&model.OnRequest{}).Return(&model.OnResponse{}, nil).Times(1)

	err = brewery.run()
	assert.NoError(t, err)
}

func TestBreweryRunPower(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := NewMockBrewery(mockCtrl)
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
	brewery := NewMockBrewery(mockCtrl)
	err := brewery.run()
	assert.Nil(t, err)
}

func TestBreweryGetConstraintsError(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := NewMockBrewery(mockCtrl)

	brewery.BoilSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(60)}, fmt.Errorf("error boil")).Times(1)

	_, err := brewery.getTempConstraints()
	assert.Equal(t, fmt.Errorf("error boil"), err)

	brewery.BoilSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(60)}, nil).Times(1)

	brewery.HermsSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(55)}, fmt.Errorf("error herms")).Times(1)

	_, err = brewery.getTempConstraints()
	assert.Equal(t, fmt.Errorf("error herms"), err)

	brewery.BoilSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(60)}, nil).Times(1)

	brewery.HermsSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(55)}, nil).Times(1)

	brewery.MashSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(50.25)}, fmt.Errorf("error mash")).Times(1)

	_, err = brewery.getTempConstraints()
	assert.Equal(t, fmt.Errorf("error mash"), err)

}

func TestConstraints(t *testing.T) {
	t.Parallel()

	// Everything good.
	assert.Equal(t, 0, checkTempConstraints([]constraint{
		{min: 50, max: 100, val: 60},
		{min: 51, max: 60, val: 55},
		{min: 49.5, max: 50.5, val: 50.25},
	}))

	assert.Equal(t, -1, checkTempConstraints([]constraint{
		{min: 50, max: 100, val: 45}, //Boil too low
		{min: 51, max: 60, val: 55},
		{min: 49.5, max: 50.5, val: 50.25},
	}))

	assert.Equal(t, 1, checkTempConstraints([]constraint{
		{min: 50, max: 100, val: 101}, //Boil too high
		{min: 51, max: 60, val: 55},
		{min: 49.5, max: 50.5, val: 50.25},
	}))

	assert.Equal(t, -1, checkTempConstraints([]constraint{
		{min: 50, max: 100, val: 60},
		{min: 51, max: 60, val: 40}, // Stream too low
		{min: 49.5, max: 50.5, val: 50.25},
	}))

	assert.Equal(t, 1, checkTempConstraints([]constraint{
		{min: 50, max: 100, val: 60},
		{min: 51, max: 60, val: 61}, // Stream too high
		{min: 49.5, max: 50.5, val: 50.25},
	}))

	assert.Equal(t, -1, checkTempConstraints([]constraint{
		{min: 50, max: 100, val: 60},
		{min: 51, max: 60, val: 55},
		{min: 49.5, max: 50.5, val: 40}, // Mash too low
	}))

	assert.Equal(t, 1, checkTempConstraints([]constraint{
		{min: 50, max: 100, val: 60},
		{min: 51, max: 60, val: 55},
		{min: 49.5, max: 50.5, val: 52}, // Mash too high
	}))

	assert.Equal(t, 1, checkTempConstraints([]constraint{
		{min: 50, max: 100, val: 110},
		{min: 51, max: 60, val: 50},     //Stream too low
		{min: 49.5, max: 50.5, val: 51}, // Mash too high
	}))

	assert.Equal(t, 1, checkTempConstraints([]constraint{
		{min: 50, max: 100, val: 60},
		{min: 51, max: 60, val: 65},     // Stream too high
		{min: 49.5, max: 50.5, val: 45}, // Mash too low
	}))

}
