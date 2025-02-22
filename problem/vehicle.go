package problem

type Vehicle interface {
	HasID
	HasIndex
	EarliestDeparture() float64
	LatestArrival() float64
	Type() VehicleType
	IsReturnToDepot() bool
	StartLocation() *Location
	EndLocation() *Location
	VehicleTypeIdentifier() VehicleTypeKey
	Skills() *Skills
	UserData() any
	Break() Break
	String() string
}

type AbstractVehicle interface {
	Vehicle
	SetIndex(index int)
}

type BaseVehicle struct {
	index             int
	vehicleIdentifier VehicleTypeKey
	userData          any
}

type BaseTypeKey struct {
	index int
}

func (a *BaseTypeKey) Index() int {
	return a.index
}

func (a *BaseTypeKey) SetIndex(index int) {
	a.index = index
}

type AbstractVehicleTypeKey interface {
	HasIndex
	SetIndex(index int)
}

func (av *BaseVehicle) UserData() any {
	return av.userData
}

func (av *BaseVehicle) SetUserData(userData any) {
	av.userData = userData
}

func (av *BaseVehicle) Index() int {
	return av.index
}

func (av *BaseVehicle) SetIndex(index int) {
	av.index = index
}

func (av *BaseVehicle) VehicleTypeIdentifier() VehicleTypeKey {
	return av.vehicleIdentifier
}

func (av *BaseVehicle) SetVehicleIdentifier(vehicleTypeIdentifier VehicleTypeKey) {
	av.vehicleIdentifier = vehicleTypeIdentifier
}

type VehicleType interface {
	TypeId() string
	CapacityDimensions() *Capacity
	MaxVelocity() float64
	VehicleCostParams() VehicleCostParams
	Profile() string
	UserData() any
	String() string
	Equals(any) bool
}

type VehicleCostParams interface {
	Fix() float64
	PerTransportTimeUnit() float64
	PerDistanceUnit() float64
	PerWaitingTimeUnit() float64
	PerServiceTimeUnit() float64
	SetFix(v float64)
	SetPerTransportTimeUnit(v float64)
	SetPerDistanceUnit(v float64)
	SetPerWaitingTimeUnit(v float64)
	SetPerServiceTimeUnit(v float64)
	String() string
}

type VehicleTypeKey interface {
	AbstractVehicleTypeKey
	Type() string
	SetType(typeID string)

	StartLocation() string
	SetStartLocation(startLocation string)

	EndLocation() string
	SetEndLocation(endLocation string)

	EarliestStart() float64
	SetEarliestStart(earliestStart float64)

	LatestEnd() float64
	SetLatestEnd(latestEnd float64)

	Skills() *Skills
	SetSkills(skills *Skills)

	ReturnToDepot() bool
	SetReturnToDepot(returnToDepot bool)

	Equals(other VehicleTypeKey) bool
	String() string
}
