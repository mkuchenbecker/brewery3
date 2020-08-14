package detectable

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	protos "github.com/golang001/deltav/model/gomodel"
	"github.com/stretchr/testify/assert"
)

func getLargeRandomPS(gameLengthMinutes int64, trackableEntities int) *DetectableDatabase {
	objs := []*TraceDetectable{}
	for i := 0; i < trackableEntities; i++ {
		objs = append(objs,
			NewTraceDetectable(0, 0, 0, 0, fmt.Sprintf("%d", i)))
	}

	db := NewDetectableDatabase()

	for i := int64(0); i < gameLengthMinutes*60; i++ {
		for _, obj := range objs {
			obj.Pos.Position.X += (rand.Float64())
			obj.Pos.Position.Y += (rand.Float64())
			obj.Pos.Position.Z += (rand.Float64())
			obj.Pos.T = i
			db.Register(obj)
		}
	}
	return db
}

func TestDetectableLookback(t *testing.T) {
	InverseCSquared = 1 //Overwrite to make math easier.
	var pos PositionSystem = NewDetectableDatabase()
	pos.Register(NewTraceDetectable(0, 0, 0, 0, "A"))
	pos.Register(NewTraceDetectable(0, 0, 0, 1, "B"))
	pos.Register(NewTraceDetectable(0, 0, 0, 2, "C"))
	pos.Register(NewTraceDetectable(0, 0, 0, 3, "D"))
	pos.Register(NewTraceDetectable(0, 0, 0, 4, "E"))

	req := DetectRequest{Pos: protos.Position{Position: &protos.Vector3{X: 1, Y: 0, Z: 0}, T: 3},
		Range: 10, Filter: []PropertyType{}}

	resp := pos.Detect(req)

	assert.Equal(t, 1, len(resp.detected))
	assert.Equal(t, "C", resp.detected[0].GetID())
}

func TestDetectableRange(t *testing.T) {
	InverseCSquared = 1 //Overwrite to make math easier.
	var pos PositionSystem = NewDetectableDatabase()
	prop := Property{Intensity: 10, PropertyType: GammaRadiation}
	pos.Register(NewTraceDetectable(10, 0, 0, 0, "A").AddProperty(prop))
	pos.Register(NewTraceDetectable(100, 0, 0, 0, "B").AddProperty(prop))
	pos.Register(NewTraceDetectable(10, 0, 0, 0, "C")) //Is not gamma.

	// A is found but c is not because we're filtering for gamma.
	req := DetectRequest{Pos: protos.Position{Position: &protos.Vector3{X: 0, Y: 0, Z: 0}, T: 10},
		Range: 10, Filter: []PropertyType{GammaRadiation}}
	resp := pos.Detect(req)
	assert.Equal(t, 1, len(resp.detected))
	assert.Equal(t, "A", resp.detected[0].GetID())

	//Won't see B because its too far away
	req = DetectRequest{Pos: protos.Position{Position: &protos.Vector3{X: 0, Y: 0, Z: 0}, T: 100},
		Range: 10, Filter: []PropertyType{GammaRadiation}}

	resp = pos.Detect(req)
	assert.Equal(t, 0, len(resp.detected))

	//Will see B with higher range.
	req = DetectRequest{Pos: protos.Position{Position: &protos.Vector3{X: 0, Y: 0, Z: 0}, T: 100},
		Range: 1000, Filter: []PropertyType{GammaRadiation}}

	resp = pos.Detect(req)

	assert.Equal(t, 1, len(resp.detected))
	assert.Equal(t, "B", resp.detected[0].GetID())
}

func TestDetectableMultiItem(t *testing.T) {
	InverseCSquared = 1 //Overwrite to make math easier.
	var pos PositionSystem = NewDetectableDatabase()

	for i := int64(0); i < 100; i++ {
		r := rand.Float64() * 1000
		pos.Register(NewTraceDetectable(r, r, r, i, "A"))
		pos.Register(NewTraceDetectable(r, r, r, i, "B"))
		pos.Register(NewTraceDetectable(r, r, r, i, "C"))
		pos.Register(NewTraceDetectable(r, r, r, i, "D"))
		pos.Register(NewTraceDetectable(r, r, r, i, "E"))
		pos.Register(NewTraceDetectable(r, r, r, i, "F"))
	}

	pos.Register(NewTraceDetectable(2, 0, 0, 10, "A"))
	pos.Register(NewTraceDetectable(0, 2, 0, 10, "B"))
	pos.Register(NewTraceDetectable(0, 0, 2, 10, "C"))
	pos.Register(NewTraceDetectable(-2, 0, 0, 10, "D"))
	pos.Register(NewTraceDetectable(0, -2, 0, 10, "E"))
	pos.Register(NewTraceDetectable(0, 0, -2, 10, "F"))

	assert.Equal(t, 6, len(pos.(*DetectableDatabase).db))
	assert.Equal(t, 100, len(pos.(*DetectableDatabase).db["A"].hist))
	assert.Equal(t, 100, len(pos.(*DetectableDatabase).db["F"].hist))

	req := DetectRequest{Pos: protos.Position{Position: &protos.Vector3{X: 0, Y: 0, Z: 0}, T: 12},
		Range: 10, Filter: []PropertyType{}}

	resp := pos.Detect(req)

	assert.Equal(t, 6, len(resp.detected))
	countMap := make(map[string]int)
	for _, val := range resp.detected {
		// Each key should be unique A-F
		_, ok := countMap[val.GetID()]
		assert.False(t, ok)
		countMap[val.GetID()]++
	}
}

func TestGetLargeRandomPs(t *testing.T) {
	pos := getLargeRandomPS(60, 1000)
	assert.Equal(t, 1000, len(pos.db))
	for _, hist := range pos.db {
		assert.Equal(t, 60*60, len(hist.hist))
	}
}

func TestDetectPerformance(t *testing.T) {
	InverseCSquared = 0.001
	pos := getLargeRandomPS(60, 1000)

	beforeTime := time.Now()
	req := DetectRequest{Pos: protos.Position{Position: &protos.Vector3{X: 0, Y: 0, Z: 0}, T: 60},
		Range: 10, Filter: []PropertyType{}}

	pos.Detect(req)
	afterTime := time.Now()

	assert.WithinDuration(t, beforeTime, afterTime, 10*time.Millisecond)
}

func TestPrune(t *testing.T) {
	InverseCSquared = 0.01
	pos := NewDetectableDatabase()

	pos.Register(NewTraceDetectable(10, 0, 0, 1, "A"))
	pos.Register(NewTraceDetectable(20, 0, 0, 2, "A"))
	pos.Register(NewTraceDetectable(30, 0, 0, 3, "A"))
	pos.Register(NewTraceDetectable(40, 0, 0, 4, "A"))
	pos.Register(NewTraceDetectable(50, 0, 0, 5, "A"))
	pos.Register(NewTraceDetectable(60, 0, 0, 6, "A"))
	pos.Register(NewTraceDetectable(70, 0, 0, 7, "A"))
	pos.Register(NewTraceDetectable(80, 0, 0, 8, "A"))
	pos.Register(NewTraceDetectable(90, 0, 0, 9, "A"))
	// 10 seconds to origin, 20 seconds retained. Its the furtherst object form the
	// origin so we double the time it takes light to reach the origin and set that as the retention limit.
	pos.Register(NewTraceDetectable(100, 0, 0, 10, "A"))

	pos.Prune(10)
	assert.Equal(t, int64(10), pos.Size()) // Nothing should get pruned.

	for i := int64(11); i <= 30; i++ {
		pos.Register(NewTraceDetectable(100, 0, 0, i, "A"))

	}
	assert.Equal(t, int64(30), pos.Size())
	pos.Prune(20)
	assert.Equal(t, int64(30), pos.Size())
	pos.Prune(30)
	assert.Equal(t, int64(21), pos.Size()) //Inly the last 20 seconds are retained.
}

// func TestPrunePerfromance(t *testing.T) {
// 	InverseCSquared = DefaultInverseC()
// 	fmt.Printf("%f\n", InverseCSquared)
// 	pos := getLargeRandomPS(60, 1000)

// 	assert.Equal(t, int64(60*60*1000), pos.Size())

// 	// Pruning a big db is 100 ms
// 	beforeTime := time.Now()
// 	pos.Prune(60 * 60)
// 	afterTime := time.Now()

// 	assert.WithinDuration(t, beforeTime, afterTime, 100*time.Millisecond)
// 	assert.Equal(t, int64(6*1000), pos.Size())

// 	//Pruining a "pruned" db is fast.
// 	beforeTime = time.Now()
// 	pos.Prune(60 * 60)
// 	afterTime = time.Now()

// 	assert.WithinDuration(t, beforeTime, afterTime, 20*time.Millisecond)
// 	assert.Equal(t, int64(6*1000), pos.Size())
// }
