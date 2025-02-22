package activity

import (
	"fmt"
	"gsprit/problem"
	"math"
)

// DeliverService represents a delivery activity
type DeliverService struct {
	problem.BaseActivity
	delivery            problem.Delivery
	capacity            *problem.Capacity
	arrTime             float64
	endTime             float64
	theoreticalEarliest float64
	theoreticalLatest   float64
	index               int
}

// NewDeliverService creates a new delivery activity
func NewDeliverService(delivery problem.Delivery) *DeliverService {
	return &DeliverService{
		delivery:            delivery,
		capacity:            problem.Invert(delivery.Size()),
		theoreticalEarliest: 0,
		theoreticalLatest:   math.MaxFloat64,
	}
}

// Name returns the name of the activity
func (ds *DeliverService) Name() string {
	return ds.delivery.Type()
}

// Location returns the delivery location
func (ds *DeliverService) Location() *problem.Location {
	return ds.delivery.Location()
}

// SetTheoreticalEarliestOperationStartTime sets the earliest operation start time
func (ds *DeliverService) SetTheoreticalEarliestOperationStartTime(earliest float64) {
	ds.theoreticalEarliest = earliest
}

// SetTheoreticalLatestOperationStartTime sets the latest operation start time
func (ds *DeliverService) SetTheoreticalLatestOperationStartTime(latest float64) {
	ds.theoreticalLatest = latest
}

// TheoreticalEarliestOperationStartTime returns the earliest operation start time
func (ds *DeliverService) TheoreticalEarliestOperationStartTime() float64 {
	return ds.theoreticalEarliest
}

// TheoreticalLatestOperationStartTime returns the latest operation start time
func (ds *DeliverService) TheoreticalLatestOperationStartTime() float64 {
	return ds.theoreticalLatest
}

// OperationTime returns the operation time
func (ds *DeliverService) OperationTime() float64 {
	return ds.delivery.ServiceDuration()
}

// ArrTime returns the arrival time
func (ds *DeliverService) ArrTime() float64 {
	return ds.arrTime
}

// EndTime returns the end time
func (ds *DeliverService) EndTime() float64 {
	return ds.endTime
}

// SetArrTime sets the arrival time
func (ds *DeliverService) SetArrTime(arrTime float64) {
	ds.arrTime = arrTime
}

// SetEndTime sets the end time
func (ds *DeliverService) SetEndTime(endTime float64) {
	ds.endTime = endTime
}

// Duplicate creates a copy of the activity
func (ds *DeliverService) Duplicate() problem.TourActivity {
	return &DeliverService{
		delivery:            ds.delivery,
		capacity:            ds.capacity,
		arrTime:             ds.arrTime,
		endTime:             ds.endTime,
		index:               ds.index,
		theoreticalEarliest: ds.theoreticalEarliest,
		theoreticalLatest:   ds.theoreticalLatest,
	}
}

// Job returns the associated delivery job
func (ds *DeliverService) Job() problem.Job {
	return ds.delivery
}

// Size returns the capacity required for the delivery
func (ds *DeliverService) Size() *problem.Capacity {
	return ds.capacity
}

// String returns a string representation of the delivery service
func (ds *DeliverService) String() string {
	return fmt.Sprintf("[type=%s][locationId=%s][size=%v][twStart=%s][twEnd=%s]",
		ds.Name(), ds.Location().Id(), ds.Size(),
		Round(ds.TheoreticalEarliestOperationStartTime()),
		Round(ds.TheoreticalLatestOperationStartTime()))
}
