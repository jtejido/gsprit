package activity

import (
	"fmt"
	"gsprit/problem"
	"math"
)

// BreakActivity represents a scheduled break during a route.
type BreakActivity struct {
	problem.BaseActivity
	index    int
	arrTime  float64
	endTime  float64
	location *problem.Location
	duration float64
	earliest float64
	latest   float64
	breakJob problem.Break
}

// NewBreakActivity creates a new BreakActivity instance.
func NewBreakActivity(breakJob problem.Break) *BreakActivity {
	return &BreakActivity{
		breakJob: breakJob,
		duration: breakJob.ServiceDuration(),
		earliest: 0,
		latest:   math.MaxFloat64,
	}
}

// ArrTime returns the arrival time of the activity.
func (b *BreakActivity) ArrTime() float64 {
	return b.arrTime
}

// SetArrTime sets the arrival time of the activity.
func (b *BreakActivity) SetArrTime(arrTime float64) {
	b.arrTime = arrTime
}

// EndTime returns the end time of the activity.
func (b *BreakActivity) EndTime() float64 {
	return b.endTime
}

// SetEndTime sets the end time of the activity.
func (b *BreakActivity) SetEndTime(endTime float64) {
	b.endTime = endTime
}

// TheoreticalEarliestOperationStartTime returns the earliest possible start time.
func (b *BreakActivity) TheoreticalEarliestOperationStartTime() float64 {
	return b.earliest
}

// SetTheoreticalEarliestOperationStartTime sets the earliest start time.
func (b *BreakActivity) SetTheoreticalEarliestOperationStartTime(earliest float64) {
	b.earliest = earliest
}

// TheoreticalLatestOperationStartTime returns the latest possible start time.
func (b *BreakActivity) TheoreticalLatestOperationStartTime() float64 {
	return b.latest
}

// SetTheoreticalLatestOperationStartTime sets the latest start time.
func (b *BreakActivity) SetTheoreticalLatestOperationStartTime(latest float64) {
	b.latest = latest
}

// OperationTime returns the duration of the break.
func (b *BreakActivity) OperationTime() float64 {
	return b.duration
}

// SetOperationTime sets the duration of the break.
func (b *BreakActivity) SetOperationTime(duration float64) {
	b.duration = duration
}

// Location returns the location of the break activity.
func (b *BreakActivity) Location() *problem.Location {
	return b.location
}

// SetLocation sets the location of the break activity.
func (b *BreakActivity) SetLocation(breakLocation *problem.Location) {
	b.location = breakLocation
}

// Job returns the associated break job.
func (b *BreakActivity) Job() problem.Job {
	return b.breakJob
}

// Name returns the name/type of the break activity.
func (b *BreakActivity) Name() string {
	return b.breakJob.Type()
}

// Duplicate creates a new instance with the same attributes.
func (b *BreakActivity) Duplicate() problem.TourActivity {
	return &BreakActivity{
		breakJob: b.breakJob,
		arrTime:  b.ArrTime(),
		endTime:  b.EndTime(),
		location: b.Location(),
		index:    b.Index(),
		earliest: b.TheoreticalEarliestOperationStartTime(),
		latest:   b.TheoreticalLatestOperationStartTime(),
		duration: b.OperationTime(),
	}
}

// Size returns the size/capacity associated with the break.
func (b *BreakActivity) Size() *problem.Capacity {
	return b.breakJob.Size()
}

// Index returns the activity index.
func (b *BreakActivity) Index() int {
	return b.index
}

// SetIndex sets the activity index.
func (b *BreakActivity) SetIndex(index int) {
	b.index = index
}

// String provides a string representation of the break activity.
func (b *BreakActivity) String() string {
	return fmt.Sprintf("[type=%s][location=%v][size=%v][twStart=%s][twEnd=%s]",
		b.Name(), b.Location(), b.Size(),
		Round(b.TheoreticalEarliestOperationStartTime()),
		Round(b.TheoreticalLatestOperationStartTime()))
}
