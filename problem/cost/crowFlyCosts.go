package cost

import (
	"gsprit/problem"
	"gsprit/util"
	"math"
)

// CrowFlyCosts represents Euclidean distance-based travel costs.
type CrowFlyCosts struct {
	EuclideanCosts
	locations func(id string) *util.Coordinate
}

// NewCrowFlyCosts initializes a new CrowFlyCosts with given locations.
func NewCrowFlyCosts(locations func(id string) *util.Coordinate) *CrowFlyCosts {
	res := &CrowFlyCosts{locations: locations}
	res.Spi = res
	return res
}

// CalculateDistance computes the Euclidean distance between two locations.
func (c *CrowFlyCosts) CalculateDistance(fromLocation, toLocation *problem.Location) float64 {
	var from, to *util.Coordinate

	if fromLocation.Coordinate() != nil && toLocation.Coordinate() != nil {
		from = fromLocation.Coordinate()
		to = toLocation.Coordinate()
	} else if c.locations != nil {
		from = c.locations(fromLocation.Id())
		to = c.locations(toLocation.Id())
	}

	return calculateDistance(from, to)
}

// calculateDistance computes the Euclidean distance between two coordinates.
func calculateDistance(from, to *util.Coordinate) float64 {
	if from == nil || to == nil {
		return 0.0
	}
	dx := from.X - to.X
	dy := from.Y - to.Y
	return math.Sqrt(dx*dx + dy*dy)
}
