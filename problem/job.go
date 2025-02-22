package problem

type ActivityType int

const (
	ActivityTypePickup ActivityType = iota
	ActivityTypeDelivery
	ActivityTypeService
	ActivityTypeBreak
)

type Activity interface {
	ActivityType() ActivityType
	SetActivityType(at ActivityType)
	Location() *Location
	SetLocation(l *Location)
	TimeWindows() []TimeWindow
	SetTimeWindows(tw []TimeWindow)
	ServiceTime() float64
	SetServiceTime(serviceTime float64)
}

type JobType int

const (
	JobTypeShipment JobType = iota
	JobTypeService
	JobTypePickupService
	JobTypeDeliveryService
	JobTypeBreakService
)

func (jt JobType) String() string {
	switch jt {
	case JobTypeShipment:
		return "SHIPMENT"
	case JobTypeService:
		return "SERVICE"
	case JobTypePickupService:
		return "PICKUP_SERVICE"
	case JobTypeDeliveryService:
		return "DELIVERY_SERVICE"
	case JobTypeBreakService:
		return "BREAK_SERVICE"
	default:
		return "UNKNOWN"
	}
}

func (jt JobType) IsShipment() bool {
	return jt == JobTypeShipment
}

func (jt JobType) IsService() bool {
	return !jt.IsShipment()
}

func (jt JobType) IsPickup() bool {
	return jt == JobTypePickupService
}

func (jt JobType) IsDelivery() bool {
	return jt == JobTypeDeliveryService
}

func (jt JobType) IsBreak() bool {
	return jt == JobTypeBreakService
}

type Job interface {
	HasID
	HasIndex
	Size() *Capacity
	RequiredSkills() *Skills
	Name() string
	Priority() int
	MaxTimeInVehicle() float64
	Activities() []Activity
	JobType() JobType
	String() string
}

type BaseJob struct {
	index    int
	userData any
}

func (j *BaseJob) Index() int {
	return j.index
}

func (j *BaseJob) SetIndex(index int) {
	j.index = index
}

func (j *BaseJob) UserData() any {
	return j.userData
}

func (j *BaseJob) SetUserData(d any) {
	j.userData = d
}

type AbstractJob interface {
	Job
	SetIndex(index int)
	UserData() any
	SetUserData(d any)
}

type Service interface {
	AbstractJob
	TimeWindows() []TimeWindow
	ServiceDuration() float64
	Location() *Location
	TimeWindow() TimeWindow
	Type() string
}

type Break interface {
	Service
}

type Pickup interface {
	Service
}

type Delivery interface {
	Service
}

type Shipment interface {
	AbstractJob
	PickupLocation() *Location
	PickupServiceTime() float64
	DeliveryLocation() *Location
	DeliveryServiceTime() float64
	DeliveryTimeWindow() TimeWindow
	PickupTimeWindow() TimeWindow
}
