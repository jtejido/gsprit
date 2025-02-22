package job

import (
	"gsprit/problem"
)

type ActivityBuilder struct {
	activityType problem.ActivityType
	location     *problem.Location
	timeWindows  []problem.TimeWindow
	serviceTime  float64
}

func NewActivityBuilder(location *problem.Location, activityType problem.ActivityType) *ActivityBuilder {
	return &ActivityBuilder{
		location:     location,
		activityType: activityType,
	}
}

func (b *ActivityBuilder) SetTimeWindows(timeWindows []problem.TimeWindow) *ActivityBuilder {
	b.timeWindows = timeWindows
	return b
}

func (b *ActivityBuilder) SetServiceTime(serviceTime float64) *ActivityBuilder {
	b.serviceTime = serviceTime
	return b
}

func (b *ActivityBuilder) Build() *Activity {
	return NewActivityFromBuilder(b)
}

// Activity represents a job activity
type Activity struct {
	activityType problem.ActivityType
	location     *problem.Location
	timeWindows  []problem.TimeWindow
	serviceTime  float64
}

func NewActivityFromBuilder(builder *ActivityBuilder) *Activity {
	return &Activity{
		activityType: builder.activityType,
		location:     builder.location,
		timeWindows:  builder.timeWindows,
		serviceTime:  builder.serviceTime,
	}
}

// NewActivity creates a new activity
func NewActivity(activityType problem.ActivityType, location *problem.Location, timeWindows []problem.TimeWindow, serviceTime float64) *Activity {
	return &Activity{
		activityType: activityType,
		location:     location,
		timeWindows:  timeWindows,
		serviceTime:  serviceTime,
	}
}

// GetActivityType returns the type of the activity
func (a *Activity) ActivityType() problem.ActivityType {
	return a.activityType
}

func (a *Activity) SetActivityType(at problem.ActivityType) {
	a.activityType = at
}

// GetLocation returns the location of the activity
func (a *Activity) Location() *problem.Location {
	return a.location
}

func (a *Activity) SetLocation(l *problem.Location) {
	a.location = l
}

// GetTimeWindows returns the time windows of the activity
func (a *Activity) TimeWindows() []problem.TimeWindow {
	return a.timeWindows
}

func (a *Activity) SetTimeWindows(tw []problem.TimeWindow) {
	a.timeWindows = tw
}

// GetServiceTime returns the service time of the activity
func (a *Activity) ServiceTime() float64 {
	return a.serviceTime
}

func (a *Activity) SetServiceTime(serviceTime float64) {
	a.serviceTime = serviceTime
}
