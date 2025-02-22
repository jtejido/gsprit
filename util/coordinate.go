package util

import "fmt"

// Coordinate represents a 2D point with X and Y coordinates.
type Coordinate struct {
	X float64
	Y float64
}

// NewCoordinate creates a new instance of Coordinate.
func NewCoordinate(x, y float64) *Coordinate {
	return &Coordinate{X: x, Y: y}
}

// String returns a string representation of the Coordinate.
func (c Coordinate) String() string {
	return fmt.Sprintf("[x=%.6f][y=%.6f]", c.X, c.Y)
}

// Equals checks if two coordinates are equal.
func (c Coordinate) Equals(other any) bool {
	if v, ok := other.(*Coordinate); ok {
		return c.X == v.X && c.Y == v.Y
	}

	return false

}
