package job

import (
	"fmt"
	"gsprit/problem"
	"gsprit/problem/solution/route/activity"
	"math"
)

// ServiceBuilder is a generic builder for constructing various service-related job types.
type ServiceBuilder[T problem.AbstractJob] struct {
	id               string
	location         *problem.Location
	serviceType      string
	serviceTime      float64
	capacityBuilder  *problem.CapacityBuilder
	capacity         *problem.Capacity
	skillsBuilder    *problem.SkillsBuilder
	skills           *problem.Skills
	name             string
	timeWindows      activity.TimeWindows
	twAdded          bool
	priority         int
	userData         any
	maxTimeInVehicle float64
	activity         problem.Activity
}

// NewServiceBuilder creates a new instance of ServiceBuilder with a specified ID.
func NewServiceBuilder[T problem.Service](id string) *ServiceBuilder[T] {
	if id == "" {
		panic("Service ID must not be empty.")
	}
	tws := activity.NewTimeWindows()
	tw, _ := activity.NewTimeWindow(0.0, math.MaxFloat64)
	tws.Add(tw)

	return &ServiceBuilder[T]{
		id:               id,
		serviceType:      "service",
		capacityBuilder:  problem.NewCapacityBuilder(),
		skillsBuilder:    problem.NewSkillsBuilder(),
		timeWindows:      tws,
		name:             "no-name",
		priority:         2,
		maxTimeInVehicle: math.MaxFloat64,
	}
}

// SetType sets the service type (e.g., Service, Pickup, Delivery, Break).
func (b *ServiceBuilder[T]) SetType(serviceType string) *ServiceBuilder[T] {
	b.serviceType = serviceType
	return b
}

// SetLocation sets the location of the service.
func (b *ServiceBuilder[T]) SetLocation(location *problem.Location) *ServiceBuilder[T] {
	b.location = location
	return b
}

// SetServiceTime sets the service time duration.
func (b *ServiceBuilder[T]) SetServiceTime(serviceTime float64) *ServiceBuilder[T] {
	if serviceTime < 0 {
		panic("The service time of a service must be greater than or equal to zero.")
	}
	b.serviceTime = serviceTime
	return b
}

// SetUserData sets user-specific data.
func (b *ServiceBuilder[T]) SetUserData(userData any) *ServiceBuilder[T] {
	b.userData = userData
	return b
}

// AddSizeDimension adds a dimension to the service's capacity.
func (b *ServiceBuilder[T]) AddSizeDimension(dimensionIndex, dimensionValue int) *ServiceBuilder[T] {
	if dimensionValue < 0 {
		panic("The capacity value must not be negative.")
	}
	b.capacityBuilder.AddDimension(dimensionIndex, dimensionValue)
	return b
}

// SetTimeWindow sets the primary time window for the service.
func (b *ServiceBuilder[T]) SetTimeWindow(tw problem.TimeWindow) *ServiceBuilder[T] {
	b.timeWindows = activity.NewTimeWindows()
	b.timeWindows.Add(tw)
	return b
}

// AddTimeWindow adds a time window to the service.
func (b *ServiceBuilder[T]) AddTimeWindow(tw problem.TimeWindow) *ServiceBuilder[T] {
	if !b.twAdded {
		b.timeWindows = activity.NewTimeWindows()
		b.twAdded = true
	}
	b.timeWindows.Add(tw)
	return b
}

// AddTimeWindowByRange adds a time window with a start and end time.
func (b *ServiceBuilder[T]) AddTimeWindowByRange(earliest, latest float64) *ServiceBuilder[T] {
	tw, err := activity.NewTimeWindow(earliest, latest)
	if err != nil {
		panic(err)
	}
	return b.AddTimeWindow(tw)
}

// AddAllTimeWindows adds multiple time windows.
func (b *ServiceBuilder[T]) AddAllTimeWindows(timeWindows []problem.TimeWindow) *ServiceBuilder[T] {
	for _, tw := range timeWindows {
		b.AddTimeWindow(tw)
	}
	return b
}

// AddRequiredSkill adds a required skill to the service.
func (b *ServiceBuilder[T]) AddRequiredSkill(skill string) *ServiceBuilder[T] {
	b.skillsBuilder.AddSkill(skill)
	return b
}

// AddAllRequiredSkills adds multiple required skills.
func (b *ServiceBuilder[T]) AddAllRequiredSkills(skills []string) *ServiceBuilder[T] {
	b.skillsBuilder.AddAllSkills(skills)
	return b
}

// AddAllRequiredSkillsFromSkills adds skills from a Skills object.
func (b *ServiceBuilder[T]) AddAllRequiredSkillsFromSkills(skills *problem.Skills) *ServiceBuilder[T] {
	b.skillsBuilder.AddAllSkills(skills.Values())
	return b
}

// AddAllSizeDimensions adds all dimensions from a capacity object.
func (b *ServiceBuilder[T]) AddAllSizeDimensions(capacity *problem.Capacity) *ServiceBuilder[T] {
	for i := 0; i < capacity.NuOfDimensions(); i++ {
		b.AddSizeDimension(i, capacity.Get(i))
	}
	return b
}

// SetPriority sets the priority of the service.
func (b *ServiceBuilder[T]) SetPriority(priority int) *ServiceBuilder[T] {
	if priority < 1 || priority > 10 {
		panic("The priority value is not valid. Only 1 (very high) to 10 (very low) are allowed.")
	}
	b.priority = priority
	return b
}

// SetMaxTimeInVehicle sets the maximum allowed time in a vehicle (only for Deliveries & Shipments).
func (b *ServiceBuilder[T]) SetMaxTimeInVehicle(maxTimeInVehicle float64) *ServiceBuilder[T] {
	if _, ok := any(*new(T)).(problem.Delivery); !ok {
		if _, ok := any(*new(T)).(problem.Shipment); !ok {
			panic("The maximum time in vehicle is only supported for Deliveries and Shipments.")
		}
	}
	if maxTimeInVehicle < 0 {
		panic("maxTimeInVehicle should be positive.")
	}
	b.maxTimeInVehicle = maxTimeInVehicle
	return b
}

// Build constructs the Service instance.
func (b *ServiceBuilder[T]) Build() *Service {
	b.SetType("service")
	b.capacity = b.capacityBuilder.Build()
	b.skills = b.skillsBuilder.Build()

	b.activity = NewActivityBuilder(b.location, problem.ActivityTypeService).
		SetServiceTime(b.serviceTime).
		SetTimeWindows(b.timeWindows.TimeWindows()).
		Build()

	return NewServiceFromBuilder(b)
}

type Service struct {
	problem.BaseJob
	id               string
	t                string
	serviceTime      float64
	size             *problem.Capacity
	skills           *problem.Skills
	name             string
	location         *problem.Location
	timeWindows      activity.TimeWindows
	priority         int
	maxTimeInVehicle float64
	activities       []problem.Activity
}

func NewServiceFromBuilder[T problem.AbstractJob](b *ServiceBuilder[T]) *Service {
	service := new(Service)
	service.SetUserData(b.userData)
	service.id = b.id
	service.serviceTime = b.serviceTime
	service.t = b.serviceType
	service.size = b.capacity
	service.skills = b.skills
	service.name = b.name
	service.location = b.location
	service.timeWindows = b.timeWindows
	service.priority = b.priority
	service.maxTimeInVehicle = b.maxTimeInVehicle

	service.activities = append(service.activities, b.activity)
	return service
}

func (s *Service) TimeWindows() []problem.TimeWindow {
	return s.timeWindows.TimeWindows()
}

func (s *Service) JobType() problem.JobType {
	return problem.JobTypeService
}

func (s *Service) Location() *problem.Location {
	return s.location
}

func (s *Service) ServiceDuration() float64 {
	return s.serviceTime
}

func (s *Service) TimeWindow() problem.TimeWindow {
	return s.timeWindows.TimeWindows()[0]
}

func (s *Service) Type() string {
	return s.t
}

func (s *Service) Id() string {
	return s.id
}

func (s *Service) String() string {
	return fmt.Sprintf("[id=%s][name=%s][type=%s][location=%v][capacity=%v][serviceTime=%.2f][timeWindows=%v]",
		s.id, s.name, s.t, s.location, s.size, s.serviceTime, s.timeWindows)
}

func (s *Service) Size() *problem.Capacity {
	return s.size
}

func (s *Service) RequiredSkills() *problem.Skills {
	return s.skills
}

func (s *Service) Name() string {
	return s.name
}

func (s *Service) Priority() int {
	return s.priority
}

func (s *Service) MaxTimeInVehicle() float64 {
	return s.maxTimeInVehicle
}

func (s *Service) Activities() []problem.Activity {
	return s.activities
}
