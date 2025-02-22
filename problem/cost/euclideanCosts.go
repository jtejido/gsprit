package cost

import (
	"gsprit/problem"
	"gsprit/util"
	"math"
)

type EuclideanCosts struct {
	AbstractForwardVehicleRoutingTransportCosts
	Speed        float64
	DetourFactor float64
}

func NewEuclideanCosts() *EuclideanCosts {
	res := &EuclideanCosts{
		Speed:        1.0,
		DetourFactor: 1.0,
	}

	res.Spi = res

	return res
}

func (e *EuclideanCosts) TransportCost(from, to *problem.Location, time float64, driver problem.Driver, vehicle problem.Vehicle) float64 {
	distance := e.calculateDistance(from, to)
	if vehicle != nil && vehicle.Type() != nil {
		return distance * vehicle.Type().VehicleCostParams().PerDistanceUnit()
	}
	return distance
}

func (e *EuclideanCosts) TransportTime(from, to *problem.Location, time float64, driver problem.Driver, vehicle problem.Vehicle) float64 {
	return e.calculateDistance(from, to) / e.Speed
}

func (e *EuclideanCosts) Distance(from, to *problem.Location, departureTime float64, vehicle problem.Vehicle) float64 {
	return e.calculateDistance(from, to)
}

func (e *EuclideanCosts) calculateDistance(from, to *problem.Location) float64 {
	return e.calculateDistanceCoords(from.Coordinate(), to.Coordinate())
}

func (e *EuclideanCosts) calculateDistanceCoords(from, to *util.Coordinate) float64 {
	if from == nil || to == nil {
		panic("Cannot calculate Euclidean distance. Coordinates are missing.")
	}
	return EuclideanDistance(from, to) * e.DetourFactor
}

func EuclideanDistance(from, to *util.Coordinate) float64 {
	// Simple Euclidean distance formula √((x2 - x1)^2 + (y2 - y1)^2)
	dx := to.X - from.X
	dy := to.Y - from.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (e *EuclideanCosts) String() string {
	return "[name=crowFlyCosts]"
}
