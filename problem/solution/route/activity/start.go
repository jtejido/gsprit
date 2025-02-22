package activity

import (
	"fmt"
	"gsprit/problem"
)

type Start struct {
	problem.BaseActivity
	location                          *problem.Location
	theoreticalEarliestOperationStart float64
	theoreticalLatestOperationStart   float64
	endTime                           float64
	arrTime                           float64
	capacity                          *problem.Capacity
}

func NewStart(location *problem.Location, theoreticalStart, theoreticalEnd float64) *Start {
	res := &Start{
		location:                          location,
		theoreticalEarliestOperationStart: theoreticalStart,
		theoreticalLatestOperationStart:   theoreticalEnd,
		endTime:                           theoreticalStart,
		capacity:                          problem.NewCapacity(make([]int, 1)),
	}
	res.SetIndex(-1)
	return res
}

func (s *Start) TheoreticalEarliestOperationStartTime() float64 {
	return s.theoreticalEarliestOperationStart
}

func (s *Start) TheoreticalLatestOperationStartTime() float64 {
	return s.theoreticalLatestOperationStart
}

func (s *Start) SetTheoreticalEarliestOperationStartTime(time float64) {
	s.theoreticalEarliestOperationStart = time
}

func (s *Start) SetTheoreticalLatestOperationStartTime(time float64) {
	s.theoreticalLatestOperationStart = time
}

func (s *Start) Location() *problem.Location {
	return s.location
}

func (s *Start) SetLocation(l *problem.Location) {
	s.location = l
}

func (s *Start) OperationTime() float64 {
	return 0.0
}

func (s *Start) Name() string {
	return "start"
}

func (s *Start) ArrTime() float64 {
	return s.arrTime
}

func (s *Start) EndTime() float64 {
	return s.endTime
}

func (s *Start) SetArrTime(arrTime float64) {
	s.arrTime = arrTime
}

func (s *Start) SetEndTime(endTime float64) {
	s.endTime = endTime
}

func (s *Start) Duplicate() problem.TourActivity {
	res := &Start{
		BaseActivity:                      s.BaseActivity,
		location:                          s.location,
		theoreticalEarliestOperationStart: s.theoreticalEarliestOperationStart,
		theoreticalLatestOperationStart:   s.theoreticalLatestOperationStart,
		endTime:                           s.endTime,
		arrTime:                           s.arrTime,
	}
	res.SetIndex(-1)
	return res
}

func (s *Start) Size() *problem.Capacity {
	return s.capacity
}

func (s *Start) String() string {
	return fmt.Sprintf("[type=%s][location=%s][twStart=%.2f][twEnd=%.2f]",
		s.Name(), s.location.String(), s.theoreticalEarliestOperationStart, s.theoreticalLatestOperationStart)
}

func (s *Start) Copy() *Start {
	res := &Start{
		location:                          s.Location(),
		theoreticalEarliestOperationStart: s.TheoreticalEarliestOperationStartTime(),
		theoreticalLatestOperationStart:   s.TheoreticalLatestOperationStartTime(),
		endTime:                           s.EndTime(),
	}
	res.SetIndex(-1)
	return res
}
