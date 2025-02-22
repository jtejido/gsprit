package activity

import (
	"fmt"
	"gsprit/problem"
	"math"
)

// PickupShipment represents an activity where a shipment is picked up.
type PickupShipment struct {
	problem.BaseActivity
	shipment problem.Shipment
	endTime  float64
	arrTime  float64
	earliest float64
	latest   float64
	index    int
}

// NewPickupShipment creates a new PickupShipment instance.
func NewPickupShipment(shipment problem.Shipment) *PickupShipment {
	return &PickupShipment{
		shipment: shipment,
		earliest: 0,
		latest:   math.MaxFloat64, // Double.MAX_VALUE equivalent
	}
}

// Job returns the associated job.
func (ps *PickupShipment) Job() problem.Job {
	return ps.shipment
}

// SetEarliestOperationStartTime sets the earliest operation start time.
func (ps *PickupShipment) SetTheoreticalEarliestOperationStartTime(earliest float64) {
	ps.earliest = earliest
}

// SetLatestOperationStartTime sets the latest operation start time.
func (ps *PickupShipment) SetTheoreticalLatestOperationStartTime(latest float64) {
	ps.latest = latest
}

// Name returns the name of the activity.
func (ps *PickupShipment) Name() string {
	return "pickupShipment"
}

// Location returns the pickup location of the shipment.
func (ps *PickupShipment) Location() *problem.Location {
	return ps.shipment.PickupLocation()
}

// EarliestOperationStartTime returns the earliest operation start time.
func (ps *PickupShipment) TheoreticalEarliestOperationStartTime() float64 {
	return ps.earliest
}

func (ps *PickupShipment) TheoreticalLatestOperationStartTime() float64 {
	return ps.latest
}

// EarliestOperationStartTime returns the earliest operation start time.
func (ps *PickupShipment) ArrTime() float64 {
	return ps.arrTime
}

func (ps *PickupShipment) EndTime() float64 {
	return ps.endTime
}

func (ps *PickupShipment) OperationTime() float64 {
	return ps.shipment.PickupServiceTime()
}

func (ps *PickupShipment) SetArrTime(arrTime float64) {
	ps.arrTime = arrTime
}

func (ps *PickupShipment) SetEndTime(endTime float64) {
	ps.endTime = endTime
}

func (ps *PickupShipment) Size() *problem.Capacity {
	return ps.shipment.Size()
}

func (ps *PickupShipment) Duplicate() problem.TourActivity {
	return &PickupShipment{
		shipment: ps.shipment,
		arrTime:  ps.arrTime,
		endTime:  ps.endTime,
		index:    ps.index,
		earliest: ps.earliest,
		latest:   ps.latest,
	}
}

func (ps *PickupShipment) String() string {
	return fmt.Sprintf("[type=%s][locationId=%s][size=%s][twStart=%s][twEnd=%s]", ps.Name(), ps.Location().Id(), ps.Size().String(), Round(ps.TheoreticalEarliestOperationStartTime()), Round(ps.TheoreticalLatestOperationStartTime()))
}
