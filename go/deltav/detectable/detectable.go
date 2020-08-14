package detectable

import (
	"fmt"
	"math"

	"github.com/golang001/deltav/common/counters"
	protos "github.com/golang001/deltav/model/gomodel"
)

// TraceDetectable is a stuct point-in-time history of an object.
type TraceDetectable struct {
	Pos        protos.Position
	Properties map[PropertyType]*Property
	ID         string
}

func NewTraceDetectable(x float64, y float64, z float64, t int64, id string) *TraceDetectable {
	return &TraceDetectable{
		Pos:        protos.Position{Position: &protos.Vector3{X: x, Y: y, Z: z}, T: t},
		ID:         id,
		Properties: make(map[PropertyType]*Property)}
}

func (td *TraceDetectable) AddProperty(p Property) *TraceDetectable {
	td.Properties[p.PropertyType] = &p
	return td
}

func (td *TraceDetectable) GetPosition() protos.Position {
	return td.Pos
}

func (td *TraceDetectable) GetProperty(propType PropertyType) (*Property, error) {
	prop, ok := td.Properties[propType]
	if !ok {
		return nil, fmt.Errorf("property not found")
	}
	return prop, nil
}

func (td *TraceDetectable) GetID() string {
	return td.ID
}

func (td *TraceDetectable) Compare(req DetectRequest) bool {
	if req.Pos.T < td.Pos.T {
		return false
	}
	deltaPos := td.Pos.Subtract(&req.Pos)
	sqMag := deltaPos.Position.MagnitudeSquared()
	lightTravelTimeSquared := sqMag * InverseCSquared
	deltaTimeSquared := deltaPos.T * deltaPos.T
	if math.Abs(lightTravelTimeSquared-float64(deltaTimeSquared)) >= 1 {
		return false
	}

	if req.Range*req.Range < sqMag {
		return false
	}

	if len(req.Filter) == 0 {
		return true
	}

	for _, pt := range req.Filter {
		for key := range td.Properties {
			if key == pt {
				return true
			}
		}
	}

	return false
}

type DetectableHistory struct {
	hist             map[int64]Detectable
	maxMagnitudeSqrd float64
}

func NewDetectableHistory() *DetectableHistory {
	return &DetectableHistory{hist: make(map[int64]Detectable), maxMagnitudeSqrd: 0}
}

func (hist *DetectableHistory) Insert(d Detectable) {
	hist.hist[d.GetPosition().T] = d
	pos := d.GetPosition()
	var mag float64 = pos.Position.MagnitudeSquared()
	if hist.maxMagnitudeSqrd < mag {
		hist.maxMagnitudeSqrd = mag
	}
}

func (hist *DetectableHistory) Prune(minTime int64) {
	for k := range hist.hist {
		if k < minTime {
			delete(hist.hist, k)
		}
	}
}

type DetectableDatabase struct {
	//Todo, make threadsafe
	db    map[string](*DetectableHistory)
	count *counters.Counter
}

func (db *DetectableDatabase) Size() int64 {
	total := int64(0)
	for _, hist := range db.db {
		total += int64(len(hist.hist))
	}
	return total
}

func NewDetectableDatabase() *DetectableDatabase {
	return &DetectableDatabase{
		db:    make(map[string]*DetectableHistory),
		count: counters.New(),
	}
}

func (db *DetectableDatabase) Register(obj Detectable) {
	hist, ok := db.db[obj.GetID()]
	if !ok {
		hist = NewDetectableHistory()
	}

	hist.Insert(obj)
	db.db[obj.GetID()] = hist
}

func detectThreaded(req DetectRequest, hist *DetectableHistory, detected chan *Detectable) {
	for i := req.Pos.T; i >= 0; i-- {
		d, ok := hist.hist[i]
		if !ok {
			continue
		}
		if d.Compare(req) {
			detected <- &d
			break
		}
	}
	detected <- nil
}

func (db *DetectableDatabase) Prune(currTime int64) {
	// Finds the max distance from the origin from a trackable object
	// and how long light takes to travel to it. Then, it doubles that
	// time and removes all elements further than that away from the
	// origin.
	maxMagSqrd := float64(0)
	for _, hist := range db.db {
		if hist.maxMagnitudeSqrd > maxMagSqrd {
			maxMagSqrd = hist.maxMagnitudeSqrd
		}
	}

	fmt.Printf("maxMag %f", maxMagSqrd)
	minTime := currTime - int64(math.Sqrt(maxMagSqrd)*2*math.Sqrt(InverseCSquared))
	fmt.Printf("pruning before time %d", minTime)
	done := make([](chan bool), len(db.db))
	i := 0
	for _, hist := range db.db {
		done[i] = make(chan bool)
		go func(hist *DetectableHistory, minTime int64, d chan bool) {
			hist.Prune(minTime)
			d <- true
		}(hist, minTime, done[i])
		i++
	}
	for _, d := range done {
		<-d
	}
}

func (db *DetectableDatabase) Detect(req DetectRequest) DetectResponse {
	resp := DetectResponse{detected: []Detectable{}}
	detectedChannels := make([](chan *Detectable), len(db.db))
	counter := 0
	for key, hist := range db.db {
		db.count.Inc(fmt.Sprintf("detecting.%s", key))
		detectedChannels[counter] = make(chan *Detectable)
		go detectThreaded(req, hist, detectedChannels[counter])
		counter++
	}
	for _, ch := range detectedChannels {
		if det := <-ch; det != nil {
			resp.detected = append(resp.detected, *(det))
		}
	}
	return resp
}
