package job

import (
	"fmt"
	"gsprit/problem"
	"gsprit/problem/solution/route/activity"
	"math"
)

var _ problem.Pickup = (*Pickup)(nil)

type PickupBuilder struct {
	ServiceBuilder[*Pickup]
}

func NewPickupBuilder(id string) *PickupBuilder {
	if id == "" {
		panic("Service ID must not be empty.")
	}
	tws := activity.NewTimeWindows()
	tw, _ := activity.NewTimeWindow(0.0, math.MaxFloat64)
	tws.Add(tw)

	b := &PickupBuilder{
		ServiceBuilder: ServiceBuilder[*Pickup]{
			id:               id,
			serviceType:      "service",
			capacityBuilder:  problem.NewCapacityBuilder(),
			skillsBuilder:    problem.NewSkillsBuilder(),
			timeWindows:      tws,
			name:             "no-name",
			priority:         2,
			maxTimeInVehicle: math.MaxFloat64,
		},
	}

	return b
}

func (b *PickupBuilder) SetMaxTimeInVehicle(maxTimeInVehicle float64) *PickupBuilder {
	panic("maxTimeInVehicle is not yet supported for Pickups and Services (only for Deliveries and Shipments)")
}

// SetType sets the service type (e.g., Service, Pickup, Delivery, Break).
func (b *PickupBuilder) SetType(serviceType string) *PickupBuilder {
	b.ServiceBuilder.SetType(serviceType)
	return b
}

// SetLocation sets the location of the service.
func (b *PickupBuilder) SetLocation(location *problem.Location) *PickupBuilder {
	b.ServiceBuilder.SetLocation(location)
	return b
}

// SetServiceTime sets the service time duration.
func (b *PickupBuilder) SetServiceTime(serviceTime float64) *PickupBuilder {
	b.ServiceBuilder.SetServiceTime(serviceTime)
	return b
}

// SetUserData sets user-specific data.
func (b *PickupBuilder) SetUserData(userData any) *PickupBuilder {
	b.ServiceBuilder.SetUserData(userData)
	return b
}

// AddSizeDimension adds a dimension to the service's capacity.
func (b *PickupBuilder) AddSizeDimension(dimensionIndex, dimensionValue int) *PickupBuilder {
	b.ServiceBuilder.AddSizeDimension(dimensionIndex, dimensionValue)
	return b
}

// SetTimeWindow sets the primary time window for the service.
func (b *PickupBuilder) SetTimeWindow(tw problem.TimeWindow) *PickupBuilder {
	b.ServiceBuilder.SetTimeWindow(tw)
	return b
}

// AddTimeWindow adds a time window to the service.
func (b *PickupBuilder) AddTimeWindow(tw problem.TimeWindow) *PickupBuilder {
	b.ServiceBuilder.AddTimeWindow(tw)
	return b
}

// AddTimeWindowByRange adds a time window with a start and end time.
func (b *PickupBuilder) AddTimeWindowByRange(earliest, latest float64) *PickupBuilder {
	b.ServiceBuilder.AddTimeWindowByRange(earliest, latest)
	return b
}

// AddAllTimeWindows adds multiple time windows.
func (b *PickupBuilder) AddAllTimeWindows(timeWindows []problem.TimeWindow) *PickupBuilder {
	b.ServiceBuilder.AddAllTimeWindows(timeWindows)
	return b
}

// AddRequiredSkill adds a required skill to the service.
func (b *PickupBuilder) AddRequiredSkill(skill string) *PickupBuilder {
	b.ServiceBuilder.AddRequiredSkill(skill)
	return b
}

// AddAllRequiredSkills adds multiple required skills.
func (b *PickupBuilder) AddAllRequiredSkills(skills []string) *PickupBuilder {
	b.ServiceBuilder.AddAllRequiredSkills(skills)
	return b
}

// AddAllRequiredSkillsFromSkills adds skills from a Skills object.
func (b *PickupBuilder) AddAllRequiredSkillsFromSkills(skills *problem.Skills) *PickupBuilder {
	b.ServiceBuilder.AddAllRequiredSkillsFromSkills(skills)
	return b
}

// AddAllSizeDimensions adds all dimensions from a capacity object.
func (b *PickupBuilder) AddAllSizeDimensions(capacity *problem.Capacity) *PickupBuilder {
	b.ServiceBuilder.AddAllSizeDimensions(capacity)
	return b
}

// SetPriority sets the priority of the service.
func (b *PickupBuilder) SetPriority(priority int) *PickupBuilder {
	b.ServiceBuilder.SetPriority(priority)
	return b
}

func (b *PickupBuilder) Build() *Pickup {
	if b.location == nil {
		panic("location is missing")
	}
	b.SetType("pickup")
	b.capacity = b.capacityBuilder.Build()
	b.skills = b.skillsBuilder.Build()
	b.activity = NewActivityBuilder(b.location, problem.ActivityTypePickup).
		SetServiceTime(b.serviceTime).
		SetTimeWindows(b.timeWindows.TimeWindows()).
		Build()

	return newPickupFromBuilder(b)
}

type Pickup struct {
	Service
}

func newPickupFromBuilder(b *PickupBuilder) *Pickup {
	res := &Pickup{
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
	}
	res.SetUserData(b.userData)
	return res
}

func NewPickup(id string, location *problem.Location) (*Pickup, error) {
	if location == nil {
		return nil, fmt.Errorf("location is missing")
	}
	twi := activity.NewTimeWindows()
	tw, _ := activity.NewTimeWindow(0., math.MaxFloat64)
	twi.Add(tw)

	return &Pickup{
		Service: Service{
			id:               id,
			name:             "no-name",
			t:                "pickup",
			size:             problem.NewCapacity(make([]int, 1)),
			skills:           problem.NewSkills(),
			timeWindows:      twi,
			maxTimeInVehicle: math.MaxFloat64,
			priority:         2,
			activities:       []problem.Activity{NewActivity(problem.ActivityTypePickup, nil, twi.TimeWindows(), 0.)},
		},
	}, nil
}

func (p *Pickup) JobType() problem.JobType {
	return problem.JobTypePickupService
}
