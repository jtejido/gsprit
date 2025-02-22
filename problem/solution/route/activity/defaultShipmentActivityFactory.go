package activity

import (
	"gsprit/problem"
)

var _ TourShipmentActivityFactory = (*DefaultShipmentActivityFactory)(nil)

type DefaultShipmentActivityFactory struct {
}

func (f *DefaultShipmentActivityFactory) CreatePickup(shipment problem.Shipment) problem.AbstractActivity {
	return NewPickupShipment(shipment)
}

func (f *DefaultShipmentActivityFactory) CreateDelivery(shipment problem.Shipment) problem.AbstractActivity {
	return NewDeliverShipment(shipment)
}
