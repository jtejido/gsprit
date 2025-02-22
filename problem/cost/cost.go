package cost

import "gsprit/problem"

type ForwardTransportTime interface {
	TransportTime(from, to *problem.Location, departureTime float64, driver problem.Driver, vehicle problem.Vehicle) float64
}

type BackwardTransportTime interface {
	BackwardTransportTime(from, to *problem.Location, arrivalTime float64, driver problem.Driver, vehicle problem.Vehicle) float64
}

type TransportTime interface {
	ForwardTransportTime
	BackwardTransportTime
}

type ForwardTransportCost interface {
	TransportCost(from, to *problem.Location, departureTime float64, driver problem.Driver, vehicle problem.Vehicle) float64
}

type BackwardTransportCost interface {
	BackwardTransportCost(from, to *problem.Location, arrivalTime float64, driver problem.Driver, vehicle problem.Vehicle) float64
}

type TransportCost interface {
	ForwardTransportCost
	BackwardTransportCost
}

type TransportDistance interface {
	Distance(from, to *problem.Location, departureTime float64, vehicle problem.Vehicle) float64
}

type VehicleRoutingTransportCosts interface {
	TransportTime
	TransportCost
	TransportDistance
	String() string
}

type VehicleRoutingActivityCosts interface {
	ActivityCost(tourAct problem.TourActivity, arrivalTime float64, drv problem.Driver, veh problem.Vehicle) float64
	ActivityDuration(tourAct problem.TourActivity, arrivalTime float64, drv problem.Driver, veh problem.Vehicle) float64
}

var Time = struct {
	TourEnd   float64
	TourStart float64
	Undefined float64
}{
	TourEnd:   -2.0,
	TourStart: -1.0,
	Undefined: -3.0,
}

type Parameter interface {
	PenaltyForMissedTimeWindow() float64
}
