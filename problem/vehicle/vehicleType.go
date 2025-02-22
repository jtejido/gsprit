package vehicle

import (
	"fmt"
	"gsprit/problem"
	"math"
	"reflect"
)

type VehicleCostParams struct {
	fix                  float64
	perTransportTimeUnit float64
	perDistanceUnit      float64
	perWaitingTimeUnit   float64
	perServiceTimeUnit   float64
}

func NewVehicleCostParams(fix, perTransportTimeUnit, perDistanceUnit, perWaitingTimeUnit, perServiceTimeUnit float64) *VehicleCostParams {
	return &VehicleCostParams{
		fix:                  fix,
		perTransportTimeUnit: perTransportTimeUnit,
		perDistanceUnit:      perDistanceUnit,
		perWaitingTimeUnit:   perWaitingTimeUnit,
		perServiceTimeUnit:   perServiceTimeUnit,
	}
}

func (vcp *VehicleCostParams) Fix() float64 {
	return vcp.fix
}
func (vcp *VehicleCostParams) PerTransportTimeUnit() float64 {
	return vcp.perTransportTimeUnit
}
func (vcp *VehicleCostParams) PerDistanceUnit() float64 {
	return vcp.perDistanceUnit
}
func (vcp *VehicleCostParams) PerWaitingTimeUnit() float64 {
	return vcp.perWaitingTimeUnit
}
func (vcp *VehicleCostParams) PerServiceTimeUnit() float64 {
	return vcp.perServiceTimeUnit
}

func (vcp *VehicleCostParams) SetFix(v float64) {
	vcp.fix = v
}
func (vcp *VehicleCostParams) SetPerTransportTimeUnit(v float64) {
	vcp.perTransportTimeUnit = v
}
func (vcp *VehicleCostParams) SetPerDistanceUnit(v float64) {
	vcp.perDistanceUnit = v
}
func (vcp *VehicleCostParams) SetPerWaitingTimeUnit(v float64) {
	vcp.perWaitingTimeUnit = v
}
func (vcp *VehicleCostParams) SetPerServiceTimeUnit(v float64) {
	vcp.perServiceTimeUnit = v
}

func (vcp *VehicleCostParams) String() string {
	return fmt.Sprintf("[fixed=%.2f][perTime=%.2f][perDistance=%.2f][perWaitingTimeUnit=%.2f]", vcp.fix, vcp.perTransportTimeUnit, vcp.perDistanceUnit, vcp.perWaitingTimeUnit)
}

// VehicleTypeBuilder constructs a VehicleType
type VehicleTypeBuilder struct {
	id                 string
	maxVelocity        float64
	fixedCost          float64
	perDistance        float64
	perTime            float64
	perWaitingTime     float64
	perServiceTime     float64
	profile            string
	capacityBuilder    *problem.CapacityBuilder
	capacityDimensions *problem.Capacity
	dimensionAdded     bool
	userData           any
}

// NewVehicleTypeBuilder initializes a new VehicleTypeBuilder
func NewVehicleTypeBuilder(id string) *VehicleTypeBuilder {
	if id == "" {
		panic("Vehicle type ID must not be empty.")
	}

	return &VehicleTypeBuilder{
		id:              id,
		maxVelocity:     math.MaxFloat64,
		fixedCost:       0.0,
		perDistance:     1.0,
		perTime:         0.0,
		perWaitingTime:  0.0,
		perServiceTime:  0.0,
		profile:         "car",
		capacityBuilder: problem.NewCapacityBuilder(),
	}
}

// SetUserData sets custom user data
func (b *VehicleTypeBuilder) SetUserData(userData any) *VehicleTypeBuilder {
	b.userData = userData
	return b
}

// SetMaxVelocity sets the maximum speed (m/s)
func (b *VehicleTypeBuilder) SetMaxVelocity(inMeterPerSeconds float64) *VehicleTypeBuilder {
	if inMeterPerSeconds < 0.0 {
		panic("The velocity of a vehicle type cannot be smaller than zero.")
	}
	b.maxVelocity = inMeterPerSeconds
	return b
}

// SetFixedCost sets a fixed cost for this vehicle type
func (b *VehicleTypeBuilder) SetFixedCost(fixedCost float64) *VehicleTypeBuilder {
	if fixedCost < 0.0 {
		panic("Fixed costs must not be smaller than zero.")
	}
	b.fixedCost = fixedCost
	return b
}

// SetCostPerDistance sets cost per distance unit
func (b *VehicleTypeBuilder) SetCostPerDistance(perDistance float64) *VehicleTypeBuilder {
	if perDistance < 0.0 {
		panic("Cost per distance must not be smaller than zero.")
	}
	b.perDistance = perDistance
	return b
}

// SetCostPerTime sets cost per time unit (deprecated, use SetCostPerTransportTime instead)
func (b *VehicleTypeBuilder) SetCostPerTime(perTime float64) *VehicleTypeBuilder {
	if perTime < 0.0 {
		panic("Cost per time must not be smaller than zero.")
	}
	b.perTime = perTime
	return b
}

// SetCostPerTransportTime sets cost per time unit
func (b *VehicleTypeBuilder) SetCostPerTransportTime(perTime float64) *VehicleTypeBuilder {
	if perTime < 0.0 {
		panic("Cost per transport time must not be smaller than zero.")
	}
	b.perTime = perTime
	return b
}

// SetCostPerWaitingTime sets cost per waiting time
func (b *VehicleTypeBuilder) SetCostPerWaitingTime(perWaitingTime float64) *VehicleTypeBuilder {
	if perWaitingTime < 0.0 {
		panic("Cost per waiting time must not be smaller than zero.")
	}
	b.perWaitingTime = perWaitingTime
	return b
}

// SetCostPerServiceTime sets cost per service time
func (b *VehicleTypeBuilder) SetCostPerServiceTime(perServiceTime float64) *VehicleTypeBuilder {
	b.perServiceTime = perServiceTime
	return b
}

// AddCapacityDimension adds a new capacity dimension
func (b *VehicleTypeBuilder) AddCapacityDimension(dimIndex, dimVal int) *VehicleTypeBuilder {
	if dimVal < 0 {
		panic("The capacity value must not be negative.")
	}
	if b.capacityDimensions != nil {
		panic("Either use AddCapacityDimension() or SetCapacityDimensions(), but not both.")
	}
	b.dimensionAdded = true
	b.capacityBuilder.AddDimension(dimIndex, dimVal)
	return b
}

// SetCapacityDimensions sets predefined capacity dimensions
func (b *VehicleTypeBuilder) SetCapacityDimensions(capacity *problem.Capacity) *VehicleTypeBuilder {
	if b.dimensionAdded {
		panic("Either use AddCapacityDimension() or SetCapacityDimensions(), but not both.")
	}
	b.capacityDimensions = capacity
	return b
}

// SetProfile sets the profile (e.g., "car", "truck", etc.)
func (b *VehicleTypeBuilder) SetProfile(profile string) *VehicleTypeBuilder {
	b.profile = profile
	return b
}

// Build constructs and returns a VehicleTypeImpl
func (b *VehicleTypeBuilder) Build() *VehicleType {
	if b.capacityDimensions == nil {
		b.capacityDimensions = b.capacityBuilder.Build()
	}
	return newVehicleTypeFromBuilder(b)
}

type VehicleType struct {
	typeId             string
	profile            string
	vehicleCostParams  problem.VehicleCostParams
	capacityDimensions *problem.Capacity
	maxVelocity        float64
	userData           any
}

func newVehicleTypeFromBuilder(builder *VehicleTypeBuilder) *VehicleType {
	return &VehicleType{
		typeId:             builder.id,
		profile:            builder.profile,
		maxVelocity:        builder.maxVelocity,
		vehicleCostParams:  NewVehicleCostParams(builder.fixedCost, builder.perTime, builder.perDistance, builder.perWaitingTime, builder.perServiceTime),
		capacityDimensions: builder.capacityDimensions,
		userData:           builder.userData,
	}
}

func NewDefaultVehicleType(id string) *VehicleType {
	return &VehicleType{
		typeId:             id,
		profile:            "car",
		maxVelocity:        math.MaxFloat64,
		vehicleCostParams:  NewVehicleCostParams(0., 0., 1., 0., 0.),
		capacityDimensions: problem.NewCapacity(make([]int, 1)),
	}
}
func NewVehicleType(id, profile string, maxVelocity float64, fixedCost, perTime, perDistance, perWaitingTime, perServiceTime float64, capacity *problem.Capacity, userData any) (*VehicleType, error) {
	if maxVelocity < 0.0 {
		return nil, fmt.Errorf("the velocity of a vehicle (type) cannot be smaller than zero")
	}
	if fixedCost < 0.0 {
		return nil, fmt.Errorf("fixed costs must not be smaller than zero")
	}
	if perDistance < 0.0 {
		return nil, fmt.Errorf("cost per distance must not be smaller than zero")
	}
	if perTime < 0.0 {
		return nil, fmt.Errorf("cost per time must not be smaller than zero")
	}

	if capacity == nil {
		capacity = problem.NewCapacity(make([]int, 1))
	}

	return &VehicleType{
		typeId:             id,
		profile:            profile,
		maxVelocity:        maxVelocity,
		vehicleCostParams:  NewVehicleCostParams(fixedCost, perTime, perDistance, perWaitingTime, perServiceTime),
		capacityDimensions: capacity,
		userData:           userData,
	}, nil
}

func (v *VehicleType) TypeId() string {
	return v.typeId
}

func (v *VehicleType) UserData() any {
	return v.userData
}

func (v *VehicleType) VehicleCostParams() problem.VehicleCostParams {
	return v.vehicleCostParams
}

func (v *VehicleType) MaxVelocity() float64 {
	return v.maxVelocity
}

func (v *VehicleType) CapacityDimensions() *problem.Capacity {
	return v.capacityDimensions
}

func (v *VehicleType) Profile() string {
	return v.profile
}

func (v *VehicleType) String() string {
	return fmt.Sprintf("[typeId=%s][capacity=%v][costs=%v]", v.typeId, v.capacityDimensions, v.vehicleCostParams)
}

func (v *VehicleType) Equals(o any) bool {
	other, ok := o.(*VehicleType)
	if !ok {
		return false
	}
	return reflect.DeepEqual(v, other)
}
