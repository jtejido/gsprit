package activity

import (
	"fmt"
	"gsprit/problem"
	"math"
)

// PickupService represents a pickup activity in a route.
type PickupService struct {
	problem.BaseActivity
	pickup              problem.Service
	arrTime             float64
	depTime             float64
	theoreticalEarliest float64
	theoreticalLatest   float64
	index               int
}

// NewPickupService creates a new PickupService from a Pickup.
func NewPickupService(pickup problem.Pickup) *PickupService {
	return &PickupService{
		pickup:              pickup,
		theoreticalEarliest: 0,
		theoreticalLatest:   math.MaxFloat64, // Equivalent to Double.MAX_VALUE
	}
}

// Duplicate creates a deep copy of PickupService.
func (ps *PickupService) Duplicate() problem.TourActivity {
	return &PickupService{
		pickup:              ps.pickup,
		arrTime:             ps.arrTime,
		depTime:             ps.depTime,
		index:               ps.index,
		theoreticalEarliest: ps.theoreticalEarliest,
		theoreticalLatest:   ps.theoreticalLatest,
	}
}

// Name returns the name/type of the pickup activity.
func (ps *PickupService) Name() string {
	return ps.pickup.Type()
}

// Location returns the location of the pickup.
func (ps *PickupService) Location() *problem.Location {
	return ps.pickup.Location()
}

// TheoreticalEarliestOperationStartTime returns the earliest start time for the operation.
func (ps *PickupService) TheoreticalEarliestOperationStartTime() float64 {
	return ps.theoreticalEarliest
}

// TheoreticalLatestOperationStartTime returns the latest start time for the operation.
func (ps *PickupService) TheoreticalLatestOperationStartTime() float64 {
	return ps.theoreticalLatest
}

// SetTheoreticalEarliestOperationStartTime sets the earliest allowed start time.
func (ps *PickupService) SetTheoreticalEarliestOperationStartTime(earliest float64) {
	ps.theoreticalEarliest = earliest
}

// SetTheoreticalLatestOperationStartTime sets the latest allowed start time.
func (ps *PickupService) SetTheoreticalLatestOperationStartTime(latest float64) {
	ps.theoreticalLatest = latest
}

// OperationTime returns the duration of the pickup operation.
func (ps *PickupService) OperationTime() float64 {
	return ps.pickup.ServiceDuration()
}

// ArrTime returns the arrival time at the pickup location.
func (ps *PickupService) ArrTime() float64 {
	return ps.arrTime
}

// EndTime returns the departure time from the pickup location.
func (ps *PickupService) EndTime() float64 {
	return ps.depTime
}

// SetArrTime sets the arrival time.
func (ps *PickupService) SetArrTime(arrTime float64) {
	ps.arrTime = arrTime
}

// SetEndTime sets the departure time.
func (ps *PickupService) SetEndTime(endTime float64) {
	ps.depTime = endTime
}

// Job returns the associated service job.
func (ps *PickupService) Job() problem.Job {
	return ps.pickup
}

// Size returns the required capacity for the pickup.
func (ps *PickupService) Size() *problem.Capacity {
	return ps.pickup.Size()
}

// String returns a string representation of the PickupService.
func (ps *PickupService) String() string {
	return fmt.Sprintf("[type=%s][locationId=%s][size=%v][twStart=%s][twEnd=%s]",
		ps.Name(), ps.Location().Id(), ps.Size(),
		Round(ps.TheoreticalEarliestOperationStartTime()),
		Round(ps.TheoreticalLatestOperationStartTime()))
}
