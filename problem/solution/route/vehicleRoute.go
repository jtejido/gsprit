package route

import (
	"fmt"
	"gsprit/problem"
	"gsprit/problem/driver"
	"gsprit/problem/job"
	"gsprit/problem/solution/route/activity"
	"gsprit/problem/vehicle"
	"math"
)

func EmptyRoute() *VehicleRoute {
	return NewVehicleRouteBuilder(vehicle.NewNoVehicle(), driver.NewNoDriver()).Build()
}

type JobActivityFactory interface {
	CreateActivities(job problem.Job) []problem.AbstractActivity
}

type defaultJobActivityFactory struct {
	serviceActivityFactory  *activity.DefaultTourActivityFactory
	shipmentActivityFactory *activity.DefaultShipmentActivityFactory
}

func newDefaultJobActivityFactory() *defaultJobActivityFactory {
	return &defaultJobActivityFactory{
		//serviceActivityFactory:  serviceFactory,
		shipmentActivityFactory: new(activity.DefaultShipmentActivityFactory),
	}
}

func (f *defaultJobActivityFactory) CreateActivities(job problem.Job) []problem.AbstractActivity {
	var acts []problem.AbstractActivity
	if job.JobType().IsBreak() {
		acts = append(acts, activity.NewBreakActivity(job.(problem.Break)))
	} else if job.JobType().IsService() {
		acts = append(acts, f.serviceActivityFactory.CreateActivity(job.(problem.Service)))
	} else if job.JobType().IsShipment() {
		acts = append(acts, f.shipmentActivityFactory.CreatePickup(job.(problem.Shipment)))
		acts = append(acts, f.shipmentActivityFactory.CreateDelivery(job.(problem.Shipment)))
	}

	return acts
}

// VehicleRouteBuilder is a builder for constructing VehicleRoute instances.
type VehicleRouteBuilder struct {
	vehicle            problem.Vehicle
	driver             problem.Driver
	start              *activity.Start
	end                *activity.End
	tourActivities     *activity.TourActivities
	openShipments      map[*job.Shipment]bool
	openActivities     map[*job.Shipment]problem.TourActivity
	jobActivityFactory JobActivityFactory
}

// NewVehicleRouteBuilder creates a new VehicleRouteBuilder.
func NewVehicleRouteBuilder(vehicle problem.Vehicle, driver problem.Driver) *VehicleRouteBuilder {
	if vehicle == nil || driver == nil {
		panic("null arguments not accepted. Use vehicle.CreateNoVehicle() and driver.NewNoDriver()")
	}

	start := activity.NewStart(vehicle.StartLocation(), vehicle.EarliestDeparture(), math.MaxFloat64)
	start.SetEndTime(vehicle.EarliestDeparture())

	end := activity.NewEnd(vehicle.EndLocation(), 0.0, vehicle.LatestArrival())

	return &VehicleRouteBuilder{
		vehicle:            vehicle,
		driver:             driver,
		start:              start,
		end:                end,
		tourActivities:     activity.NewTourActivities(),
		openShipments:      make(map[*job.Shipment]bool),
		openActivities:     make(map[*job.Shipment]problem.TourActivity),
		jobActivityFactory: newDefaultJobActivityFactory(),
	}
}

// SetJobActivityFactory sets a custom job activity factory.
func (b *VehicleRouteBuilder) SetJobActivityFactory(factory JobActivityFactory) *VehicleRouteBuilder {
	b.jobActivityFactory = factory
	return b
}

// SetDepartureTime sets the departure time of the vehicle route.
func (b *VehicleRouteBuilder) SetDepartureTime(departureTime float64) *VehicleRouteBuilder {
	if departureTime < b.start.EndTime() {
		panic("departureTime < vehicle.EarliestDepartureTime. This must not be.")
	}
	b.start.SetEndTime(departureTime)
	return b
}

// AddService adds a service activity to the route.
func (b *VehicleRouteBuilder) AddService(service problem.Service) *VehicleRouteBuilder {
	return b.AddServiceWithTimeWindow(service, service.TimeWindow())
}

// AddServiceWithTimeWindow adds a service with a specified time window.
func (b *VehicleRouteBuilder) AddServiceWithTimeWindow(service problem.Service, tw problem.TimeWindow) *VehicleRouteBuilder {
	if service == nil {
		panic("service must not be nil")
	}
	acts := b.jobActivityFactory.CreateActivities(service)
	act := acts[0]
	act.SetTheoreticalEarliestOperationStartTime(tw.Start())
	act.SetTheoreticalLatestOperationStartTime(tw.End())
	b.tourActivities.AddActivityToEnd(act)
	return b
}

// AddBreak adds a break activity.
func (b *VehicleRouteBuilder) AddBreak(breakJob problem.Break, tw activity.TimeWindow, location *problem.Location) *VehicleRouteBuilder {
	if breakJob == nil {
		panic("break must not be nil")
	}
	acts := b.jobActivityFactory.CreateActivities(breakJob)
	breakAct := acts[0].(*activity.BreakActivity)
	breakAct.SetTheoreticalEarliestOperationStartTime(tw.Start())
	breakAct.SetTheoreticalLatestOperationStartTime(tw.End())
	breakAct.SetLocation(location)
	b.tourActivities.AddActivityToEnd(breakAct)
	return b
}

// AddPickup adds a pickup to the route.
func (b *VehicleRouteBuilder) AddPickup(pickup problem.Pickup) *VehicleRouteBuilder {
	return b.AddService(pickup)
}

// AddPickupWithTimeWindow adds a pickup with a specific time window.
func (b *VehicleRouteBuilder) AddPickupWithTimeWindow(pickup problem.Pickup, tw problem.TimeWindow) *VehicleRouteBuilder {
	return b.AddServiceWithTimeWindow(pickup, tw)
}

// AddDelivery adds a delivery to the route.
func (b *VehicleRouteBuilder) AddDelivery(delivery problem.Delivery) *VehicleRouteBuilder {
	return b.AddService(delivery)
}

// AddDeliveryWithTimeWindow adds a delivery with a specific time window.
func (b *VehicleRouteBuilder) AddDeliveryWithTimeWindow(delivery problem.Delivery, tw problem.TimeWindow) *VehicleRouteBuilder {
	return b.AddServiceWithTimeWindow(delivery, tw)
}

// AddPickupForShipment adds a pickup activity for a shipment.
func (b *VehicleRouteBuilder) AddPickupForShipment(shipment *job.Shipment) *VehicleRouteBuilder {
	return b.AddPickupForShipmentWithTimeWindow(shipment, shipment.PickupTimeWindow())
}

// AddPickupForShipmentWithTimeWindow adds a shipment pickup with a time window.
func (b *VehicleRouteBuilder) AddPickupForShipmentWithTimeWindow(shipment *job.Shipment, tw problem.TimeWindow) *VehicleRouteBuilder {
	if b.openShipments[shipment] {
		panic("shipment has already been added. Cannot add it twice.")
	}

	acts := b.jobActivityFactory.CreateActivities(shipment)
	act := acts[0]
	act.SetTheoreticalEarliestOperationStartTime(tw.Start())
	act.SetTheoreticalLatestOperationStartTime(tw.End())

	b.tourActivities.AddActivityToEnd(act)
	b.openShipments[shipment] = true
	b.openActivities[shipment] = acts[1] // Store the delivery activity for later.
	return b
}

// AddDeliveryForShipment adds a delivery activity for a shipment.
func (b *VehicleRouteBuilder) AddDeliveryForShipment(shipment *job.Shipment) *VehicleRouteBuilder {
	return b.AddDeliveryForShipmentWithTimeWindow(shipment, shipment.DeliveryTimeWindow())
}

// AddDeliveryForShipmentWithTimeWindow adds a shipment delivery with a time window.
func (b *VehicleRouteBuilder) AddDeliveryForShipmentWithTimeWindow(shipment *job.Shipment, tw problem.TimeWindow) *VehicleRouteBuilder {
	if !b.openShipments[shipment] {
		panic(fmt.Sprintf("cannot deliver shipment. Shipment %v needs to be picked up first.", shipment))
	}

	act := b.openActivities[shipment]
	act.SetTheoreticalEarliestOperationStartTime(tw.Start())
	act.SetTheoreticalLatestOperationStartTime(tw.End())
	b.tourActivities.AddActivityToEnd(act)

	delete(b.openShipments, shipment)
	delete(b.openActivities, shipment)
	return b
}

// Build constructs the VehicleRoute instance.
func (b *VehicleRouteBuilder) Build() *VehicleRoute {
	if len(b.openShipments) > 0 {
		panic("there are still shipments that have not been delivered yet.")
	}

	if !b.vehicle.IsReturnToDepot() && !b.tourActivities.IsEmpty() {
		// Set the end location to the last activity’s location if vehicle does not return to depot.
		lastAct := b.tourActivities.Activities()[len(b.tourActivities.Activities())-1]
		b.end.SetLocation(lastAct.Location())
	}

	return NewVehicleRoute(b)
}

type VehicleRoute struct {
	tourActivities problem.TourActivities
	vehicle        problem.Vehicle
	driver         problem.Driver
	start          *activity.Start
	end            *activity.End
	index          int
}

func NewVehicleRoute(builder *VehicleRouteBuilder) *VehicleRoute {
	return &VehicleRoute{
		tourActivities: builder.tourActivities,
		vehicle:        builder.vehicle,
		driver:         builder.driver,
		start:          builder.start,
		end:            builder.end,
		index:          -1,
	}
}

func (vr *VehicleRoute) Index() int {
	return vr.index
}

func (vr *VehicleRoute) SetIndex(index int) {
	vr.index = index
}

func (vr *VehicleRoute) Activities() []problem.TourActivity {
	return vr.tourActivities.Activities()
}

func (vr *VehicleRoute) TourActivities() problem.TourActivities {
	return vr.tourActivities
}

func (vr *VehicleRoute) Vehicle() problem.Vehicle {
	return vr.vehicle
}

func (vr *VehicleRoute) Driver() problem.Driver {
	return vr.driver
}

func (vr *VehicleRoute) SetVehicleAndDepartureTime(vehicle problem.Vehicle, vehicleDepTime float64) {
	vr.vehicle = vehicle
	vr.setStartAndEnd(vehicle, vehicleDepTime)
}

func (vr *VehicleRoute) setStartAndEnd(vehicle problem.Vehicle, vehicleDepTime float64) {
	if vr.start == nil && vr.end == nil {
		vr.start = activity.NewStart(vehicle.StartLocation(), vehicle.EarliestDeparture(), vehicle.LatestArrival())
		vr.end = activity.NewEnd(vehicle.EndLocation(), vehicle.EarliestDeparture(), vehicle.LatestArrival())
	}
	vr.start.SetEndTime(max(vehicleDepTime, vehicle.EarliestDeparture()))
	vr.start.SetLocation(vehicle.StartLocation())
	vr.end.SetLocation(vehicle.EndLocation())
}

func (vr *VehicleRoute) DepartureTime() (float64, error) {
	if vr.start == nil {
		return 0., fmt.Errorf("cannot get departureTime without having a vehicle on this route")
	}
	return vr.start.EndTime(), nil
}

func (vr *VehicleRoute) IsEmpty() bool {
	return vr.tourActivities.IsEmpty()
}

func (vr *VehicleRoute) Start() *activity.Start {
	return vr.start
}

func (vr *VehicleRoute) End() *activity.End {
	return vr.end
}

func (vr *VehicleRoute) String() string {
	return fmt.Sprintf("[start=%v][end=%v][departureTime=%.2f][vehicle=%v][driver=%v][nuOfActs=%d]",
		vr.start, vr.end, vr.start.EndTime(), vr.vehicle, vr.driver, len(vr.tourActivities.Activities()))
}

func (vr *VehicleRoute) Copy() *VehicleRoute {
	return &VehicleRoute{
		start:          vr.start.Copy(),
		end:            vr.end.Copy(),
		tourActivities: vr.tourActivities.Copy(),
		vehicle:        vr.vehicle,
		driver:         vr.driver,
	}
}
