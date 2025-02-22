package vehicle

import (
	"fmt"
	"gsprit/problem"
	"math"
)

// VehicleBuilder constructs a VehicleImpl
type VehicleBuilder struct {
	id            string
	earliestStart float64
	latestArrival float64
	returnToDepot bool
	vehicleType   problem.VehicleType
	skillBuilder  *problem.SkillsBuilder
	skills        *problem.Skills
	startLocation *problem.Location
	endLocation   *problem.Location
	vehicleBreak  problem.Break
	userData      any
}

// NewVehicleBuilder initializes a new VehicleBuilder
func NewVehicleBuilder(id string) *VehicleBuilder {
	if id == "" {
		panic("Vehicle ID must not be empty.")
	}

	return &VehicleBuilder{
		id:            id,
		earliestStart: 0.0,
		latestArrival: math.MaxFloat64,
		returnToDepot: true,
		vehicleType:   NewVehicleTypeBuilder("default").Build(),
		skillBuilder:  problem.NewSkillsBuilder(),
	}
}

// SetType sets the vehicle type
func (b *VehicleBuilder) SetType(vehicleType problem.VehicleType) *VehicleBuilder {
	if vehicleType == nil {
		panic("Vehicle type must not be nil.")
	}
	b.vehicleType = vehicleType
	return b
}

// SetUserData sets custom user data
func (b *VehicleBuilder) SetUserData(userData any) *VehicleBuilder {
	b.userData = userData
	return b
}

// SetReturnToDepot sets whether the vehicle must return to depot
func (b *VehicleBuilder) SetReturnToDepot(returnToDepot bool) *VehicleBuilder {
	b.returnToDepot = returnToDepot
	return b
}

// SetStartLocation sets the start location of the vehicle
func (b *VehicleBuilder) SetStartLocation(startLocation *problem.Location) *VehicleBuilder {
	if startLocation == nil {
		panic(fmt.Sprintf("Start location of vehicle %s must not be nil.", b.id))
	}
	b.startLocation = startLocation
	return b
}

// SetEndLocation sets the end location of the vehicle
func (b *VehicleBuilder) SetEndLocation(endLocation *problem.Location) *VehicleBuilder {
	b.endLocation = endLocation
	return b
}

// SetEarliestStart sets the earliest departure time
func (b *VehicleBuilder) SetEarliestStart(earliestStart float64) *VehicleBuilder {
	if earliestStart < 0 {
		panic(fmt.Sprintf("The earliest start time of vehicle %s must not be negative.", b.id))
	}
	b.earliestStart = earliestStart
	return b
}

// SetLatestArrival sets the latest arrival time
func (b *VehicleBuilder) SetLatestArrival(latestArrival float64) *VehicleBuilder {
	if latestArrival < 0 {
		panic(fmt.Sprintf("The latest arrival time of vehicle %s must not be negative.", b.id))
	}
	b.latestArrival = latestArrival
	return b
}

// AddSkill adds a single skill to the vehicle
func (b *VehicleBuilder) AddSkill(skill string) *VehicleBuilder {
	if skill == "" {
		panic(fmt.Sprintf("Skill of vehicle %s must not be empty.", b.id))
	}
	b.skillBuilder.AddSkill(skill)
	return b
}

// AddAllSkills adds multiple skills
func (b *VehicleBuilder) AddAllSkills(skills []string) *VehicleBuilder {
	if skills == nil {
		panic(fmt.Sprintf("Skills of vehicle %s must not be nil.", b.id))
	}
	b.skillBuilder.AddAllSkills(skills)
	return b
}

// AddSkillsFromObject adds skills from an existing Skills object
func (b *VehicleBuilder) AddSkillsFromObject(skills *problem.Skills) *VehicleBuilder {
	if skills != nil {
		b.skillBuilder.AddAllSkills(skills.Values())
	}
	return b
}

// SetBreak sets a break for the vehicle
func (b *VehicleBuilder) SetBreak(vehicleBreak problem.Break) *VehicleBuilder {
	b.vehicleBreak = vehicleBreak
	return b
}

// Build constructs and returns a VehicleImpl
func (b *VehicleBuilder) Build() *Vehicle {
	if b.latestArrival < b.earliestStart {
		panic(fmt.Sprintf("The latest arrival time of vehicle %s must not be smaller than its start time.", b.id))
	}

	if b.startLocation != nil && b.endLocation != nil {
		if b.startLocation.Id() != b.endLocation.Id() && !b.returnToDepot {
			panic(fmt.Sprintf("You specified both the end location and that the vehicle %s does not need to return to its end location. This must not be.", b.id))
		}
	}

	// If end location is not specified, default to start location
	if b.startLocation != nil && b.endLocation == nil {
		b.endLocation = b.startLocation
	}

	if b.startLocation == nil && b.endLocation == nil {
		panic(fmt.Sprintf("Every vehicle requires a start location, but vehicle %s does not have one.", b.id))
	}

	// Finalize skills
	b.skills = b.skillBuilder.Build()

	return newVehicleFromBuilder(b)
}

type Vehicle struct {
	problem.BaseVehicle
	id                string
	t                 problem.VehicleType
	earliestDeparture float64
	latestArrival     float64
	returnToDepot     bool
	skills            *problem.Skills
	startLocation     *problem.Location
	endLocation       *problem.Location
	vehicleBreak      problem.Break
}

func newVehicleFromBuilder(builder *VehicleBuilder) *Vehicle {
	res := &Vehicle{
		id:                builder.id,
		t:                 builder.vehicleType,
		earliestDeparture: builder.earliestStart,
		latestArrival:     builder.latestArrival,
		returnToDepot:     builder.returnToDepot,
		skills:            builder.skills,
		endLocation:       builder.endLocation,
		startLocation:     builder.startLocation,
		vehicleBreak:      builder.vehicleBreak,
	}
	res.SetUserData(builder.userData)
	res.SetVehicleIdentifier(NewVehicleTypeKey(res.t.TypeId(), res.startLocation.Id(), res.endLocation.Id(), res.earliestDeparture, res.latestArrival, res.skills, res.returnToDepot))
	return res
}

func NewVehicle(
	id string,
	vehicleType problem.VehicleType,
	returnToDepot bool,
	startLocation *problem.Location,
	endLocation *problem.Location,
	earliestStart float64,
	latestArrival float64,
	skills *problem.Skills,
	vehicleBreak problem.Break,
	userData any,
) (*Vehicle, error) {

	if vehicleType == nil {
		vehicleType = NewDefaultVehicleType("default")
	}

	if startLocation == nil {
		return nil, fmt.Errorf("start location must not be nil")
	}
	if earliestStart < 0 {
		return nil, fmt.Errorf("earliest start time must not be negative")
	}
	if latestArrival < 0 {
		return nil, fmt.Errorf("latest arrival time must not be negative")
	}
	if endLocation == nil && returnToDepot {
		endLocation = startLocation
	}
	if latestArrival == 0 {
		latestArrival = math.MaxFloat64
	}

	if skills == nil {
		skills = problem.NewSkills()
	}

	res := &Vehicle{
		id:                id,
		t:                 vehicleType,
		earliestDeparture: earliestStart,
		latestArrival:     latestArrival,
		returnToDepot:     returnToDepot,
		skills:            skills,
		startLocation:     startLocation,
		endLocation:       endLocation,
		vehicleBreak:      vehicleBreak,
	}
	res.SetUserData(userData)
	res.SetVehicleIdentifier(NewVehicleTypeKey(vehicleType.TypeId(), startLocation.Id(), endLocation.Id(), earliestStart, latestArrival, skills, returnToDepot))

	return res, nil
}

func (v *Vehicle) Id() string                       { return v.id }
func (v *Vehicle) Type() problem.VehicleType        { return v.t }
func (v *Vehicle) EarliestDeparture() float64       { return v.earliestDeparture }
func (v *Vehicle) LatestArrival() float64           { return v.latestArrival }
func (v *Vehicle) IsReturnToDepot() bool            { return v.returnToDepot }
func (v *Vehicle) Skills() *problem.Skills          { return v.skills }
func (v *Vehicle) StartLocation() *problem.Location { return v.startLocation }
func (v *Vehicle) EndLocation() *problem.Location   { return v.endLocation }
func (v *Vehicle) Break() problem.Break             { return v.vehicleBreak }

func (v *Vehicle) String() string {
	return fmt.Sprintf("[id=%s][type=%v][startLocation=%v][endLocation=%v][isReturnToDepot=%v][skills=%v]",
		v.id, v.t, v.startLocation, v.endLocation, v.returnToDepot, v.skills)
}

type NoVehicle struct {
	Vehicle
}

func NewNoVehicle() *NoVehicle {
	return &NoVehicle{
		Vehicle: Vehicle{
			id:                "noVehicle",
			t:                 NewDefaultVehicleType("noType"),
			earliestDeparture: 0,
			latestArrival:     0,
			skills:            nil,
			returnToDepot:     false,
			startLocation:     nil,
			endLocation:       nil,
			vehicleBreak:      nil,
		},
	}
}
