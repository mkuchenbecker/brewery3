package servers

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

func TestElementPowerLevelToggle(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSwitch := mocks.NewMockSwitchClient(mockCtrl)

	// Base test.
	mockSwitch.EXPECT().On(context.Background(), gomock.Any()).Return(&model.OnResponse{}, nil).Times(1)
	//Off always gets called at the end of the function in addition to each toggle.
	mockSwitch.EXPECT().Off(context.Background(), gomock.Any()).Return(&model.OffResponse{}, nil).Times(2)
	// Interval is 1 second, time on is 0.75s, quit after 1 second.
	// Result is that the switch should turn on, .75s later turn off, and on the next iteration quit.

	brewery := Brewery{Element: mockSwitch}
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
	mockSwitch.EXPECT().Off(context.Background(), gomock.Any()).Return(&model.OffResponse{}, nil).Times(11)
	// Interval is 1 second, time on is 0.75s, quit after 1.25 second.
	// Two iterations should occur.
	brewery := Brewery{Element: mockSwitch}
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
	brewery := Brewery{Element: mockSwitch}
	err := brewery.ElementOn()
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
	brewery := Brewery{Element: mockSwitch}
	err := brewery.elementPowerLevelToggle(9*time.Millisecond, 100*time.Millisecond, 10*time.Millisecond)
	assert.Error(t, err)

	mockSwitch.EXPECT().On(context.Background(),
		gomock.Any()).Return(&model.OnResponse{}, nil).Times(1)
	// Happens in multiples of 3 due to retries.
	mockSwitch.EXPECT().Off(context.Background(),
		gomock.Any()).Return(&model.OffResponse{}, fmt.Errorf("unable to turn coil off")).Times(6)

	err = brewery.elementPowerLevelToggle(9*time.Millisecond, 100*time.Millisecond, 10*time.Millisecond)
	assert.Error(t, err)
}

func TestElementPowerLevel0(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSwitch := mocks.NewMockSwitchClient(mockCtrl)

	mockSwitch.EXPECT().Off(context.Background(), gomock.Any()).Return(&model.OffResponse{}, nil).Times(1)
	brewery := Brewery{Element: mockSwitch}
	err := brewery.ElementPowerLevel(0)
	assert.NoError(t, err)
}

func TestElementPowerLevel101(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSwitch := mocks.NewMockSwitchClient(mockCtrl)

	mockSwitch.EXPECT().Off(context.Background(), gomock.Any()).Return(&model.OffResponse{}, nil).Times(1)
	brewery := Brewery{Element: mockSwitch}
	err := brewery.ElementPowerLevel(101)
	assert.NoError(t, err)
}

func TestElementPowerLevel100(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSwitch := mocks.NewMockSwitchClient(mockCtrl)

	mockSwitch.EXPECT().On(context.Background(), gomock.Any()).Return(&model.OnResponse{}, nil).Times(1)
	brewery := Brewery{Element: mockSwitch}
	err := brewery.ElementPowerLevel(100)
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

func TestGetConstraintsBackground(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := NewMockBrewery(mockCtrl)
	_, err := brewery.Control(context.Background(), &model.ControlRequest{Scheme: MashConfig})
	assert.NoError(t, err)
	//Since the tempsread is true the RPCs should still occur but the current temperatures are used.
	brewery.boilTemp = 10
	brewery.mashTemp = 10
	brewery.hermsTemp = 10
	brewery.tempsRead = true

	brewery.BoilSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(60)}, nil).Times(1)

	brewery.HermsSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(55)}, nil).Times(1)

	brewery.MashSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(50.25)}, nil).Times(1)

	constraints, err := brewery.getTempConstraints()
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond) // Pause so the background RPCs can finish.
	// The expected values are all 10 because the temperatures are read in the background.
	assert.Equal(t, []constraint{{min: 50, max: 100, val: 10},
		{min: 51, max: 60, val: 10},
		{min: 49.5, max: 50.5, val: 10}}, constraints)
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

	err = brewery.Run()
	assert.NoError(t, err)

	brewery.tempsRead = false // Reset so we send GRPCs to get new temperatures.
	brewery.BoilSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(60)}, nil).Times(1)
	brewery.HermsSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(55)}, nil).Times(1)
	brewery.MashSensor.(*mocks.MockThermometerClient).EXPECT().Get(context.Background(),
		&model.GetRequest{}).Return(&model.GetResponse{Temperature: float64(45)}, nil).Times(1)
	brewery.Element.(*mocks.MockSwitchClient).EXPECT().On(context.Background(),
		&model.OnRequest{}).Return(&model.OnResponse{}, nil).Times(1)

	err = brewery.Run()
	assert.NoError(t, err)
}

func TestBreweryRunBoil(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := NewMockBrewery(mockCtrl)

	_, err := brewery.Control(context.Background(), &model.ControlRequest{Scheme: &model.ControlScheme{Scheme: &model.ControlScheme_Boil_{
		Boil: &model.ControlScheme_Boil{}}}})
	assert.NoError(t, err)

	brewery.Element.(*mocks.MockSwitchClient).EXPECT().On(context.Background(),
		&model.OnRequest{}).Return(&model.OnResponse{}, nil).Times(1)
	err = brewery.Run()
	assert.NoError(t, err)
}

func TestBreweryPower(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	brewery := NewMockBrewery(mockCtrl)
	_, err := brewery.Control(context.Background(), &model.ControlRequest{Scheme: &model.ControlScheme{Scheme: &model.ControlScheme_Power_{
		Power: &model.ControlScheme_Power{PowerLevel: 100}}}})
	assert.NoError(t, err)

	brewery.Element.(*mocks.MockSwitchClient).EXPECT().On(context.Background(),
		&model.OnRequest{}).Return(&model.OnResponse{}, nil).Times(1)

	err = brewery.Run()
	assert.NoError(t, err)
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
