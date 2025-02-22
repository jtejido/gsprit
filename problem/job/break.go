package job

import (
	"gsprit/problem"
	"gsprit/problem/solution/route/activity"
	"math"
)

var _ problem.Break = (*Break)(nil)

type Break struct {
	Service
	variableLocation bool
}

type BreakBuilder struct {
	ServiceBuilder[*Break]
	variableLocation bool
}

func NewBreakBuilder(id string) *BreakBuilder {
	if id == "" {
		panic("Service ID must not be empty.")
	}
	tws := activity.NewTimeWindows()
	tw, _ := activity.NewTimeWindow(0.0, math.MaxFloat64)
	tws.Add(tw)

	b := &BreakBuilder{
		ServiceBuilder: ServiceBuilder[*Break]{
			id:               id,
			serviceType:      "service",
			capacityBuilder:  problem.NewCapacityBuilder(),
			skillsBuilder:    problem.NewSkillsBuilder(),
			timeWindows:      tws,
			name:             "no-name",
			priority:         2,
			maxTimeInVehicle: math.MaxFloat64,
		},
		variableLocation: true,
	}

	return b
}

func (b *BreakBuilder) SetMaxTimeInVehicle(maxTimeInVehicle float64) *BreakBuilder {
	if maxTimeInVehicle < 0 {
		panic("maxTimeInVehicle should not be negative")
	}
	b.maxTimeInVehicle = maxTimeInVehicle
	return b
}

// SetType sets the service type (e.g., Service, Pickup, Delivery, Break).
func (b *BreakBuilder) SetType(serviceType string) *BreakBuilder {
	b.ServiceBuilder.SetType(serviceType)
	return b
}

// SetLocation sets the location of the service.
func (b *BreakBuilder) SetLocation(location *problem.Location) *BreakBuilder {
	b.ServiceBuilder.SetLocation(location)
	return b
}

// SetServiceTime sets the service time duration.
func (b *BreakBuilder) SetServiceTime(serviceTime float64) *BreakBuilder {
	b.ServiceBuilder.SetServiceTime(serviceTime)
	return b
}

// SetUserData sets user-specific data.
func (b *BreakBuilder) SetUserData(userData any) *BreakBuilder {
	b.ServiceBuilder.SetUserData(userData)
	return b
}

// AddSizeDimension adds a dimension to the service's capacity.
func (b *BreakBuilder) AddSizeDimension(dimensionIndex, dimensionValue int) *BreakBuilder {
	b.ServiceBuilder.AddSizeDimension(dimensionIndex, dimensionValue)
	return b
}

// SetTimeWindow sets the primary time window for the service.
func (b *BreakBuilder) SetTimeWindow(tw problem.TimeWindow) *BreakBuilder {
	b.ServiceBuilder.SetTimeWindow(tw)
	return b
}

// AddTimeWindow adds a time window to the service.
func (b *BreakBuilder) AddTimeWindow(tw problem.TimeWindow) *BreakBuilder {
	b.ServiceBuilder.AddTimeWindow(tw)
	return b
}

// AddTimeWindowByRange adds a time window with a start and end time.
func (b *BreakBuilder) AddTimeWindowByRange(earliest, latest float64) *BreakBuilder {
	b.ServiceBuilder.AddTimeWindowByRange(earliest, latest)
	return b
}

// AddAllTimeWindows adds multiple time windows.
func (b *BreakBuilder) AddAllTimeWindows(timeWindows []problem.TimeWindow) *BreakBuilder {
	b.ServiceBuilder.AddAllTimeWindows(timeWindows)
	return b
}

// AddRequiredSkill adds a required skill to the service.
func (b *BreakBuilder) AddRequiredSkill(skill string) *BreakBuilder {
	b.ServiceBuilder.AddRequiredSkill(skill)
	return b
}

// AddAllRequiredSkills adds multiple required skills.
func (b *BreakBuilder) AddAllRequiredSkills(skills []string) *BreakBuilder {
	return b
}

// AddAllRequiredSkillsFromSkills adds skills from a Skills object.
func (b *BreakBuilder) AddAllRequiredSkillsFromSkills(skills *problem.Skills) *BreakBuilder {
	b.ServiceBuilder.AddAllRequiredSkillsFromSkills(skills)
	return b
}

// AddAllSizeDimensions adds all dimensions from a capacity object.
func (b *BreakBuilder) AddAllSizeDimensions(capacity *problem.Capacity) *BreakBuilder {
	b.ServiceBuilder.AddAllSizeDimensions(capacity)
	return b
}

// SetPriority sets the priority of the service.
func (b *BreakBuilder) SetPriority(priority int) *BreakBuilder {
	b.ServiceBuilder.SetPriority(priority)
	return b
}

func (b *BreakBuilder) Build() *Break {
	if b.location == nil {
		b.variableLocation = false
	}
	b.SetType("break")
	b.capacity = b.capacityBuilder.Build()
	b.skills = b.skillsBuilder.Build()
	b.activity = NewActivityBuilder(b.location, problem.ActivityTypeBreak).
		SetServiceTime(b.serviceTime).
		SetTimeWindows(b.timeWindows.TimeWindows()).
		Build()

	return newBreakFromBuilder(b)
}

func newBreakFromBuilder(b *BreakBuilder) *Break {
	res := &Break{
		Service: Service{
			id:               b.id,
			serviceTime:      b.serviceTime,
			t:                b.serviceType,
			size:             b.capacity,
			skills:           b.skills,
			name:             b.name,
			location:         b.location,
			timeWindows:      b.timeWindows,
			priority:         b.priority,
			maxTimeInVehicle: b.maxTimeInVehicle,
			activities:       []problem.Activity{b.activity},
		},
		variableLocation: b.variableLocation,
	}
	res.SetUserData(b.userData)
	return res
}

func (b *Break) HasVariableLocation() bool {
	return b.variableLocation
}

func (b *Break) JobType() problem.JobType {
	return problem.JobTypeBreakService
}
