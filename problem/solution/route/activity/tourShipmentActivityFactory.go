package activity

import "gsprit/problem"

type TourShipmentActivityFactory interface {
	CreatePickup(shipment problem.Shipment) problem.AbstractActivity
	CreateDelivery(shipment problem.Shipment) problem.AbstractActivity
}
