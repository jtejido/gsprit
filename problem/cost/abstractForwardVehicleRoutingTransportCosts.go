package cost

import (
	"gsprit/problem"
)

type AbstractForwardVehicleRoutingTransportCosts struct {
	Spi VehicleRoutingTransportCosts
}

func (a *AbstractForwardVehicleRoutingTransportCosts) BackwardTransportTime(from, to *problem.Location, arrivalTime float64, driver problem.Driver, vehicle problem.Vehicle) float64 {
	return a.Spi.TransportTime(from, to, arrivalTime, driver, vehicle)
}

func (a *AbstractForwardVehicleRoutingTransportCosts) BackwardTransportCost(from, to *problem.Location, arrivalTime float64, driver problem.Driver, vehicle problem.Vehicle) float64 {
	return a.Spi.TransportCost(from, to, arrivalTime, driver, vehicle)
}
