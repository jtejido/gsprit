package cost

import (
	"gsprit/problem"
	"math"
)

type WaitingTimeCosts struct {
}

func (c *WaitingTimeCosts) ActivityCost(tourAct problem.TourActivity, arrivalTime float64, driver problem.Driver, vehicle problem.Vehicle) float64 {
	if vehicle != nil {
		waiting := vehicle.Type().VehicleCostParams().PerWaitingTimeUnit() * math.Max(0., tourAct.TheoreticalEarliestOperationStartTime()-arrivalTime)
		servicing := vehicle.Type().VehicleCostParams().PerServiceTimeUnit() * c.ActivityDuration(tourAct, arrivalTime, driver, vehicle)
		return waiting + servicing
	}
	return 0.
}

func (c *WaitingTimeCosts) ActivityDuration(tourAct problem.TourActivity, arrivalTime float64, driver problem.Driver, vehicle problem.Vehicle) float64 {
	return tourAct.OperationTime()
}
