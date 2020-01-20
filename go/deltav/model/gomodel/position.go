package deltav_model

/*
This file is a value-added addition helper functions on the position/vector struct.
TODO: Decouple wire interface with go working model.
*/
import "math"

// MagnitudeSquared takes a Vector3 and gets the magnitude squared.
// Magnitude will not be written because Square root is expensive and
// one should be cognisant of expensive operations so this is fast by
// default.
func (p *Vector3) MagnitudeSquared() float64 {
	return math.Pow(p.X, 2) + math.Pow(p.Y, 2) + math.Pow(p.Z, 2)
}

// Subtract gets the difference between two vectors.
func (p *Vector3) Subtract(n *Vector3) *Vector3 {
	return &Vector3{X: p.X - n.X, Y: p.Y - n.Y, Z: p.Z - n.Z}
}

// Subtract gets the difference between two positions.
func (p *Position) Subtract(n *Position) *Position {
	return &Position{Position: p.Position.Subtract(n.Position), T: p.T - n.T}
}
