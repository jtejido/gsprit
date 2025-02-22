package job

import (
	"fmt"
	"gsprit/problem"
	"gsprit/problem/solution/route/activity"
	"math"
)

var _ problem.Delivery = (*Delivery)(nil)

type Delivery struct {
	Service
}

type DeliveryBuilder struct {
	ServiceBuilder[*Delivery]
}

func NewDeliveryBuilder(id string) *DeliveryBuilder {
	if id == "" {
		panic("Service ID must not be empty.")
	}
	tws := activity.NewTimeWindows()
	tw, _ := activity.NewTimeWindow(0.0, math.MaxFloat64)
	tws.Add(tw)

	b := &DeliveryBuilder{
		ServiceBuilder: ServiceBuilder[*Delivery]{
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

func (b *DeliveryBuilder) SetMaxTimeInVehicle(maxTimeInVehicle float64) *DeliveryBuilder {
	if maxTimeInVehicle < 0 {
		panic("maxTimeInVehicle should not be negative")
	}
	b.maxTimeInVehicle = maxTimeInVehicle
	return b
}

// SetType sets the service type (e.g., Service, Pickup, Delivery, Break).
func (b *DeliveryBuilder) SetType(serviceType string) *DeliveryBuilder {
	b.ServiceBuilder.SetType(serviceType)
	return b
}

// SetLocation sets the location of the service.
func (b *DeliveryBuilder) SetLocation(location *problem.Location) *DeliveryBuilder {
	b.ServiceBuilder.SetLocation(location)
	return b
}

// SetServiceTime sets the service time duration.
func (b *DeliveryBuilder) SetServiceTime(serviceTime float64) *DeliveryBuilder {
	b.ServiceBuilder.SetServiceTime(serviceTime)
	return b
}

// SetUserData sets user-specific data.
func (b *DeliveryBuilder) SetUserData(userData any) *DeliveryBuilder {
	b.ServiceBuilder.SetUserData(userData)
	return b
}

// AddSizeDimension adds a dimension to the service's capacity.
func (b *DeliveryBuilder) AddSizeDimension(dimensionIndex, dimensionValue int) *DeliveryBuilder {
	b.ServiceBuilder.AddSizeDimension(dimensionIndex, dimensionValue)
	return b
}

// SetTimeWindow sets the primary time window for the service.
func (b *DeliveryBuilder) SetTimeWindow(tw problem.TimeWindow) *DeliveryBuilder {
	b.ServiceBuilder.SetTimeWindow(tw)
	return b
}

// AddTimeWindow adds a time window to the service.
func (b *DeliveryBuilder) AddTimeWindow(tw problem.TimeWindow) *DeliveryBuilder {
	b.ServiceBuilder.AddTimeWindow(tw)
	return b
}

// AddTimeWindowByRange adds a time window with a start and end time.
func (b *DeliveryBuilder) AddTimeWindowByRange(earliest, latest float64) *DeliveryBuilder {
	b.ServiceBuilder.AddTimeWindowByRange(earliest, latest)
	return b
}

// AddAllTimeWindows adds multiple time windows.
func (b *DeliveryBuilder) AddAllTimeWindows(timeWindows []problem.TimeWindow) *DeliveryBuilder {
	b.ServiceBuilder.AddAllTimeWindows(timeWindows)
	return b
}

// AddRequiredSkill adds a required skill to the service.
func (b *DeliveryBuilder) AddRequiredSkill(skill string) *DeliveryBuilder {
	b.ServiceBuilder.AddRequiredSkill(skill)
	return b
}

// AddAllRequiredSkills adds multiple required skills.
func (b *DeliveryBuilder) AddAllRequiredSkills(skills []string) *DeliveryBuilder {
	return b
}

// AddAllRequiredSkillsFromSkills adds skills from a Skills object.
func (b *DeliveryBuilder) AddAllRequiredSkillsFromSkills(skills *problem.Skills) *DeliveryBuilder {
	b.ServiceBuilder.AddAllRequiredSkillsFromSkills(skills)
	return b
}

// AddAllSizeDimensions adds all dimensions from a capacity object.
func (b *DeliveryBuilder) AddAllSizeDimensions(capacity *problem.Capacity) *DeliveryBuilder {
	b.ServiceBuilder.AddAllSizeDimensions(capacity)
	return b
}

// SetPriority sets the priority of the service.
func (b *DeliveryBuilder) SetPriority(priority int) *DeliveryBuilder {
	b.ServiceBuilder.SetPriority(priority)
	return b
}

func (b *DeliveryBuilder) Build() *Delivery {
	if b.location == nil {
		panic("location is missing")
	}
	b.SetType("delivery")
	b.capacity = b.capacityBuilder.Build()
	b.skills = b.skillsBuilder.Build()
	b.activity = NewActivityBuilder(b.location, problem.ActivityTypeDelivery).
		SetServiceTime(b.serviceTime).
		SetTimeWindows(b.timeWindows.TimeWindows()).
		Build()

	return newDeliveryFromBuilder(b)
}

func newDeliveryFromBuilder(b *DeliveryBuilder) *Delivery {
	res := &Delivery{
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

func NewDelivery(id string, location *problem.Location) (*Delivery, error) {
	if location == nil {
		return nil, fmt.Errorf("location is missing")
	}
	twi := activity.NewTimeWindows()
	tw, _ := activity.NewTimeWindow(0., math.MaxFloat64)
	twi.Add(tw)

	return &Delivery{
		Service: Service{
			id:               id,
			name:             "no-name",
			t:                "delivery",
			size:             problem.NewCapacity(make([]int, 1)),
			skills:           problem.NewSkills(),
			timeWindows:      twi,
			maxTimeInVehicle: math.MaxFloat64,
			priority:         2,
			activities:       []problem.Activity{NewActivity(problem.ActivityTypeDelivery, nil, twi.TimeWindows(), 0.)},
		},
	}, nil
}

func (p *Delivery) JobType() problem.JobType {
	return problem.JobTypeDeliveryService
}
