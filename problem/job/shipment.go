package job

import (
	"fmt"
	"gsprit/problem"
	"gsprit/problem/solution/route/activity"
	"math"
)

var _ problem.Shipment = (*Shipment)(nil)

type ShipmentBuilder struct {
	id                                             string
	pickupServiceTime                              float64
	deliveryServiceTime                            float64
	capacityBuilder                                *problem.CapacityBuilder
	capacity                                       *problem.Capacity
	skillBuilder                                   *problem.SkillsBuilder
	skills                                         *problem.Skills
	name                                           string
	pickupLocation                                 *problem.Location
	deliveryLocation                               *problem.Location
	deliveryTimeWindows                            activity.TimeWindows
	pickupTimeWindows                              activity.TimeWindows
	priority                                       int
	userData                                       any
	maxTimeInVehicle                               float64
	pickup                                         problem.Activity
	delivery                                       problem.Activity
	deliveryTimeWindowAdded, pickupTimeWindowAdded bool
}

func NewShipmentBuilder(id string) *ShipmentBuilder {
	if id == "" {
		panic("ID must not be empty.")
	}
	tw1, _ := activity.NewTimeWindow(0., math.MaxFloat64)
	tw2, _ := activity.NewTimeWindow(0., math.MaxFloat64)
	ptw := activity.NewTimeWindows()
	ptw.Add(tw1)

	dtw := activity.NewTimeWindows()
	dtw.Add(tw2)
	return &ShipmentBuilder{
		id:                  id,
		capacityBuilder:     problem.NewCapacityBuilder(),
		skillBuilder:        problem.NewSkillsBuilder(),
		name:                "no-name",
		pickupTimeWindows:   ptw,
		deliveryTimeWindows: dtw,
		priority:            2,
		maxTimeInVehicle:    math.MaxFloat64,
	}
}

func (b *ShipmentBuilder) SetUserData(userData any) *ShipmentBuilder {
	b.userData = userData
	return b
}

func (b *ShipmentBuilder) SetPickupLocation(location *problem.Location) *ShipmentBuilder {
	b.pickupLocation = location
	return b
}

func (b *ShipmentBuilder) SetPickupServiceTime(serviceTime float64) *ShipmentBuilder {
	if serviceTime < 0.0 {
		panic("The service time of a shipment must not be < 0.0.")
	}
	b.pickupServiceTime = serviceTime
	return b
}

func (b *ShipmentBuilder) SetPickupTimeWindow(timeWindow problem.TimeWindow) *ShipmentBuilder {
	if timeWindow == nil {
		panic("The delivery time window must not be null.")
	}
	b.pickupTimeWindows = activity.NewTimeWindows()
	b.pickupTimeWindows.Add(timeWindow)
	return b
}

func (b *ShipmentBuilder) SetDeliveryLocation(location *problem.Location) *ShipmentBuilder {
	b.deliveryLocation = location
	return b
}

func (b *ShipmentBuilder) SetDeliveryServiceTime(serviceTime float64) *ShipmentBuilder {
	if serviceTime < 0.0 {
		panic("The service time of a delivery must not be < 0.0.")
	}
	b.deliveryServiceTime = serviceTime
	return b
}

func (b *ShipmentBuilder) SetDeliveryTimeWindow(timeWindow problem.TimeWindow) *ShipmentBuilder {
	if timeWindow == nil {
		panic("The delivery time window must not be null.")
	}
	b.deliveryTimeWindows = activity.NewTimeWindows()
	b.deliveryTimeWindows.Add(timeWindow)
	return b
}

func (b *ShipmentBuilder) AddSizeDimension(index, value int) *ShipmentBuilder {
	if value < 0 {
		panic(fmt.Sprintf("The capacity value must not be negative, but is %d.", value))
	}
	b.capacityBuilder.AddDimension(index, value)
	return b
}

func (b *ShipmentBuilder) AddAllSizeDimensions(capacity problem.Capacity) *ShipmentBuilder {
	for i := 0; i < capacity.NuOfDimensions(); i++ {
		b.AddSizeDimension(i, capacity.Get(i))
	}
	return b
}

func (b *ShipmentBuilder) AddRequiredSkill(skill string) *ShipmentBuilder {
	b.skillBuilder.AddSkill(skill)
	return b
}

func (b *ShipmentBuilder) AddAllRequiredSkills(skills []string) *ShipmentBuilder {
	b.skillBuilder.AddAllSkills(skills)
	return b
}

func (b *ShipmentBuilder) SetName(name string) *ShipmentBuilder {
	b.name = name
	return b
}

func (b *ShipmentBuilder) AddDeliveryTimeWindow(timeWindow problem.TimeWindow) *ShipmentBuilder {
	if timeWindow == nil {
		panic("The time window must not be null.")
	}
	if !b.deliveryTimeWindowAdded {
		b.deliveryTimeWindows = activity.NewTimeWindows()
		b.deliveryTimeWindowAdded = true
	}
	b.deliveryTimeWindows.Add(timeWindow)
	return b
}

func (b *ShipmentBuilder) AddAllDeliveryTimeWindows(timeWindows []problem.TimeWindow) *ShipmentBuilder {
	for _, tw := range timeWindows {
		b.AddDeliveryTimeWindow(tw)
	}
	return b
}

func (b *ShipmentBuilder) AddPickupTimeWindow(timeWindow problem.TimeWindow) *ShipmentBuilder {
	if timeWindow == nil {
		panic("The time window must not be null.")
	}
	if !b.pickupTimeWindowAdded {
		b.pickupTimeWindows = activity.NewTimeWindows()
		b.pickupTimeWindowAdded = true
	}
	b.pickupTimeWindows.Add(timeWindow)
	return b
}

func (b *ShipmentBuilder) AddAllPickupTimeWindows(timeWindows []problem.TimeWindow) *ShipmentBuilder {
	for _, tw := range timeWindows {
		b.AddPickupTimeWindow(tw)
	}
	return b
}

func (b *ShipmentBuilder) SetPriority(priority int) *ShipmentBuilder {
	if priority < 1 || priority > 10 {
		panic("The priority value is not valid. Only 1 (very high) to 10 (very low) are allowed.")
	}
	b.priority = priority
	return b
}

func (b *ShipmentBuilder) SetMaxTimeInVehicle(maxTimeInVehicle float64) *ShipmentBuilder {
	if maxTimeInVehicle < 0 {
		panic("The maximum time in vehicle must be positive.")
	}
	b.maxTimeInVehicle = maxTimeInVehicle
	return b
}

func (b *ShipmentBuilder) Build() *Shipment {
	if b.pickupLocation == nil {
		panic("The pickup location is missing.")
	}
	if b.deliveryLocation == nil {
		panic("The delivery location is missing.")
	}
	b.capacity = b.capacityBuilder.Build()
	b.skills = b.skillBuilder.Build()
	b.pickup = NewActivityBuilder(b.pickupLocation, problem.ActivityTypePickup).
		SetServiceTime(b.pickupServiceTime).
		SetTimeWindows(b.pickupTimeWindows.TimeWindows()).
		Build()
	b.delivery = NewActivityBuilder(b.deliveryLocation, problem.ActivityTypeDelivery).
		SetServiceTime(b.deliveryServiceTime).
		SetTimeWindows(b.deliveryTimeWindows.TimeWindows()).
		Build()

	return newShipmentFromBuilder(b)
}

// Shipment represents a job that includes a pickup and delivery.
type Shipment struct {
	problem.BaseJob
	id                  string
	pickupServiceTime   float64
	deliveryServiceTime float64
	capacity            *problem.Capacity
	skills              *problem.Skills
	name                string
	pickupLocation      *problem.Location
	deliveryLocation    *problem.Location
	deliveryTimeWindows activity.TimeWindows
	pickupTimeWindows   activity.TimeWindows
	priority            int
	maxTimeInVehicle    float64
	activities          []problem.Activity
}

func newShipmentFromBuilder(builder *ShipmentBuilder) *Shipment {
	activities := make([]problem.Activity, 0)
	activities = append(activities, builder.pickup)
	activities = append(activities, builder.delivery)
	res := &Shipment{
		id:                  builder.id,
		pickupServiceTime:   builder.pickupServiceTime,
		deliveryServiceTime: builder.deliveryServiceTime,
		capacity:            builder.capacity,
		skills:              builder.skills,
		name:                builder.name,
		pickupLocation:      builder.pickupLocation,
		deliveryLocation:    builder.deliveryLocation,
		deliveryTimeWindows: builder.deliveryTimeWindows,
		pickupTimeWindows:   builder.pickupTimeWindows,
		priority:            builder.priority,
		maxTimeInVehicle:    builder.maxTimeInVehicle,
		activities:          activities,
	}
	res.SetUserData(builder.userData)
	return res
}

func NewShipment(id string, pickupLocation_, deliveryLocation_ *problem.Location) (*Shipment, error) {
	if pickupLocation_ == nil {
		return nil, fmt.Errorf("the pickup location is missing")
	}
	if deliveryLocation_ == nil {
		return nil, fmt.Errorf("the delivery location is missing")
	}
	tw1, _ := activity.NewTimeWindow(0., math.MaxFloat64)
	tw2, _ := activity.NewTimeWindow(0., math.MaxFloat64)
	ptw := activity.NewTimeWindows()
	ptw.Add(tw1)

	dtw := activity.NewTimeWindows()
	dtw.Add(tw2)

	return &Shipment{
		id:                  id,
		pickupServiceTime:   0.,
		deliveryServiceTime: 0.,
		capacity:            problem.NewCapacity(make([]int, 1)),
		skills:              problem.NewSkills(),
		name:                "no-name",
		pickupLocation:      pickupLocation_,
		deliveryLocation:    deliveryLocation_,
		pickupTimeWindows:   ptw,
		deliveryTimeWindows: dtw,
		priority:            2,
		maxTimeInVehicle:    math.MaxFloat64,
	}, nil
}

// Methods for Shipment
func (s *Shipment) Id() string {
	return s.id
}

func (s *Shipment) PickupLocation() *problem.Location {
	return s.pickupLocation
}

func (s *Shipment) PickupServiceTime() float64 {
	return s.pickupServiceTime
}

func (s *Shipment) DeliveryLocation() *problem.Location {
	return s.deliveryLocation
}

func (s *Shipment) DeliveryServiceTime() float64 {
	return s.deliveryServiceTime
}

func (s *Shipment) DeliveryTimeWindow() problem.TimeWindow {
	return s.deliveryTimeWindows.TimeWindows()[0]
}

func (s *Shipment) PickupTimeWindow() problem.TimeWindow {
	return s.pickupTimeWindows.TimeWindows()[0]
}

func (s *Shipment) Size() *problem.Capacity {
	return s.capacity
}

func (s *Shipment) RequiredSkills() *problem.Skills {
	return s.skills
}

func (s *Shipment) Name() string {
	return s.name
}

func (s *Shipment) Priority() int {
	return s.priority
}

func (s *Shipment) MaxTimeInVehicle() float64 {
	return s.maxTimeInVehicle
}

func (s *Shipment) String() string {
	return fmt.Sprintf("[id=%s][name=%s][pickupLocation=%v][deliveryLocation=%v][capacity=%v][pickupServiceTime=%.2f][deliveryServiceTime=%.2f][pickupTimeWindows=%v][deliveryTimeWindows=%v]",
		s.id, s.name, s.pickupLocation, s.deliveryLocation, s.capacity, s.pickupServiceTime, s.deliveryServiceTime, s.pickupTimeWindows, s.deliveryTimeWindows)
}

func (s *Shipment) Activities() []problem.Activity {
	return s.activities
}

func (b *Shipment) JobType() problem.JobType {
	return problem.JobTypeShipment
}
