package activity

import (
	"fmt"
	"gsprit/problem"
	"math"
)

// DeliverShipment represents a delivery activity for a shipment.
type DeliverShipment struct {
	shipment problem.Shipment
	endTime  float64
	arrTime  float64
	capacity *problem.Capacity
	earliest float64
	latest   float64
	index    int
}

// NewDeliverShipment creates a new delivery activity for a shipment.
func NewDeliverShipment(shipment problem.Shipment) *DeliverShipment {
	return &DeliverShipment{
		shipment: shipment,
		capacity: problem.Invert(shipment.Size()),
		earliest: 0,
		latest:   math.MaxFloat64,
	}
}

// Duplicate creates a copy of the DeliverShipment activity.
func (ds *DeliverShipment) Duplicate() problem.TourActivity {
	return &DeliverShipment{
		shipment: ds.shipment,
		arrTime:  ds.arrTime,
		endTime:  ds.endTime,
		index:    ds.index,
		earliest: ds.earliest,
		latest:   ds.latest,
	}
}

// Job returns the associated shipment.
func (ds *DeliverShipment) Job() problem.Job {
	return ds.shipment
}

// SetEarliestOperationStartTime sets the earliest operation start time.
func (ds *DeliverShipment) SetTheoreticalEarliestOperationStartTime(earliest float64) {
	ds.earliest = earliest
}

// SetLatestOperationStartTime sets the latest operation start time.
func (ds *DeliverShipment) SetTheoreticalLatestOperationStartTime(latest float64) {
	ds.latest = latest
}

// Name returns the name of the activity.
func (ds *DeliverShipment) Name() string {
	return "deliverShipment"
}

// Location returns the delivery location.
func (ds *DeliverShipment) Location() *problem.Location {
	return ds.shipment.DeliveryLocation()
}

// EarliestOperationStartTime returns the earliest allowed operation start time.
func (ds *DeliverShipment) TheoreticalEarliestOperationStartTime() float64 {
	return ds.earliest
}

// LatestOperationStartTime returns the latest allowed operation start time.
func (ds *DeliverShipment) TheoreticalLatestOperationStartTime() float64 {
	return ds.latest
}

// OperationTime returns the delivery service time.
func (ds *DeliverShipment) OperationTime() float64 {
	return ds.shipment.DeliveryServiceTime()
}

// ArrTime returns the arrival time.
func (ds *DeliverShipment) ArrTime() float64 {
	return ds.arrTime
}

// EndTime returns the end time.
func (ds *DeliverShipment) EndTime() float64 {
	return ds.endTime
}

// SetArrTime sets the arrival time.
func (ds *DeliverShipment) SetArrTime(arrTime float64) {
	ds.arrTime = arrTime
}

// SetEndTime sets the end time.
func (ds *DeliverShipment) SetEndTime(endTime float64) {
	ds.endTime = endTime
}

// Size returns the capacity associated with the delivery.
func (ds *DeliverShipment) Size() *problem.Capacity {
	return ds.capacity
}

// Index returns the activity index.
func (ds *DeliverShipment) Index() int {
	return ds.index
}

// SetIndex sets the activity index.
func (ds *DeliverShipment) SetIndex(index int) {
	ds.index = index
}

// String returns a string representation of the DeliverShipment activity.
func (ds *DeliverShipment) String() string {
	return fmt.Sprintf("[type=%s][locationId=%s][size=%v][twStart=%s][twEnd=%s]",
		ds.Name(), ds.Location().Id(), ds.Size(),
		Round(ds.TheoreticalEarliestOperationStartTime()), Round(ds.TheoreticalLatestOperationStartTime()))
}
