package activity

import (
	"fmt"
	"gsprit/problem"
)

// ServiceActivity represents an activity associated with a service.
type ServiceActivity struct {
	problem.BaseActivity
	index               int
	arrTime             float64
	endTime             float64
	theoreticalEarliest float64
	theoreticalLatest   float64
	service             problem.Service
}

// NewServiceActivity creates a new ServiceActivity instance.
func NewServiceActivity(service problem.Service) *ServiceActivity {
	return &ServiceActivity{
		service: service,
	}
}

// ArrTime returns the arrival time of the activity.
func (s *ServiceActivity) ArrTime() float64 {
	return s.arrTime
}

// SetArrTime sets the arrival time of the activity.
func (s *ServiceActivity) SetArrTime(arrTime float64) {
	s.arrTime = arrTime
}

// EndTime returns the end time of the activity.
func (s *ServiceActivity) EndTime() float64 {
	return s.endTime
}

// SetEndTime sets the end time of the activity.
func (s *ServiceActivity) SetEndTime(endTime float64) {
	s.endTime = endTime
}

// TheoreticalEarliestOperationStartTime returns the earliest possible start time.
func (s *ServiceActivity) TheoreticalEarliestOperationStartTime() float64 {
	return s.theoreticalEarliest
}

// TheoreticalLatestOperationStartTime returns the latest possible start time.
func (s *ServiceActivity) TheoreticalLatestOperationStartTime() float64 {
	return s.theoreticalLatest
}

// SetTheoreticalEarliestOperationStartTime sets the earliest start time.
func (s *ServiceActivity) SetTheoreticalEarliestOperationStartTime(earliest float64) {
	s.theoreticalEarliest = earliest
}

// SetTheoreticalLatestOperationStartTime sets the latest start time.
func (s *ServiceActivity) SetTheoreticalLatestOperationStartTime(latest float64) {
	s.theoreticalLatest = latest
}

// OperationTime returns the service duration.
func (s *ServiceActivity) OperationTime() float64 {
	return s.service.ServiceDuration()
}

// Location returns the location of the service activity.
func (s *ServiceActivity) Location() *problem.Location {
	return s.service.Location()
}

// Job returns the associated service job.
func (s *ServiceActivity) Job() problem.Job {
	return s.service
}

// Name returns the name/type of the service.
func (s *ServiceActivity) Name() string {
	return s.service.Type()
}

// Duplicate creates a new instance with the same attributes.
func (s *ServiceActivity) Duplicate() problem.TourActivity {
	return &ServiceActivity{
		service:             s.service,
		arrTime:             s.ArrTime(),
		endTime:             s.EndTime(),
		index:               s.Index(),
		theoreticalEarliest: s.TheoreticalEarliestOperationStartTime(),
		theoreticalLatest:   s.TheoreticalLatestOperationStartTime(),
	}
}

// Size returns the size/capacity associated with the service.
func (s *ServiceActivity) Size() *problem.Capacity {
	return s.service.Size()
}

// Index returns the activity index.
func (s *ServiceActivity) Index() int {
	return s.index
}

// SetIndex sets the activity index.
func (s *ServiceActivity) SetIndex(index int) {
	s.index = index
}

// String provides a string representation of the service activity.
func (s *ServiceActivity) String() string {
	return fmt.Sprintf("[type=%s][location=%v][size=%v][twStart=%s][twEnd=%s]",
		s.Name(), s.Location(), s.Size(),
		Round(s.TheoreticalEarliestOperationStartTime()),
		Round(s.TheoreticalLatestOperationStartTime()))
}
