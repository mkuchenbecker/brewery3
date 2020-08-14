package detectable

import (
	"math"

	protos "github.com/golang001/deltav/model/gomodel"
)

const C = float64(1000) // Units / second
var InverseCSquared = math.Pow(1/C, 2)

func DefaultInverseC() float64 {
	return math.Pow(float64(1)/C, 2)
}

type Property struct {
	Intensity    float32
	PropertyType PropertyType
}

type PropertyType int

var (
	GammaRadiation PropertyType = 1
	EMRadiation    PropertyType = 2
	RFRadiation    PropertyType = 2
)

type DetectRequest struct {
	Pos    protos.Position
	Range  float64
	Filter []PropertyType
}

type DetectResponse struct {
	detected []Detectable
}

type PositionSystem interface {
	Detect(DetectRequest) DetectResponse
	Register(Detectable)
}

type Detectable interface {
	GetPosition() protos.Position
	GetProperty(PropertyType) (*Property, error)
	GetID() string
	Compare(DetectRequest) bool
}
