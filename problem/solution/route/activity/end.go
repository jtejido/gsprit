package activity

import (
	"fmt"
	"gsprit/problem"
)

type End struct {
	problem.BaseActivity
	location                          *problem.Location
	theoreticalEarliestOperationStart float64
	theoreticalLatestOperationStart   float64
	endTime                           float64
	arrTime                           float64
	capacity                          *problem.Capacity
}

func NewEnd(location *problem.Location, theoreticalStart, theoreticalEnd float64) *End {
	res := &End{
		location:                          location,
		theoreticalEarliestOperationStart: theoreticalStart,
		theoreticalLatestOperationStart:   theoreticalEnd,
		endTime:                           theoreticalStart,
		capacity:                          problem.NewCapacity(make([]int, 1)),
	}
	res.SetIndex(-2)
	return res
}

func (e *End) TheoreticalEarliestOperationStartTime() float64 {
	return e.theoreticalEarliestOperationStart
}

func (e *End) TheoreticalLatestOperationStartTime() float64 {
	return e.theoreticalLatestOperationStart
}

func (e *End) SetTheoreticalEarliestOperationStartTime(time float64) {
	e.theoreticalEarliestOperationStart = time
}

func (e *End) SetTheoreticalLatestOperationStartTime(time float64) {
	e.theoreticalLatestOperationStart = time
}

func (e *End) Location() *problem.Location {
	return e.location
}

func (e *End) SetLocation(l *problem.Location) {
	e.location = l
}

func (e *End) OperationTime() float64 {
	return 0.0
}

func (e *End) Name() string {
	return "end"
}

func (e *End) ArrTime() float64 {
	return e.arrTime
}

func (e *End) EndTime() float64 {
	return e.endTime
}

func (e *End) SetArrTime(arrTime float64) {
	e.arrTime = arrTime
}

func (e *End) SetEndTime(endTime float64) {
	e.endTime = endTime
}

func (e *End) Duplicate() problem.TourActivity {
	res := &End{
		BaseActivity:                      e.BaseActivity,
		location:                          e.location,
		theoreticalEarliestOperationStart: e.theoreticalEarliestOperationStart,
		theoreticalLatestOperationStart:   e.theoreticalLatestOperationStart,
		endTime:                           e.endTime,
		arrTime:                           e.arrTime,
	}
	res.SetIndex(-2)
	return res
}

func (e *End) Size() *problem.Capacity {
	return e.capacity
}

func (e *End) String() string {
	return fmt.Sprintf("[type=%s][location=%s][twStart=%.2f][twEnd=%.2f]",
		e.Name(), e.location.String(), e.theoreticalEarliestOperationStart, e.theoreticalLatestOperationStart)
}

func (e *End) Copy() *End {
	res := &End{
		location:                          e.Location(),
		theoreticalEarliestOperationStart: e.TheoreticalEarliestOperationStartTime(),
		theoreticalLatestOperationStart:   e.TheoreticalLatestOperationStartTime(),
		arrTime:                           e.ArrTime(),
		endTime:                           e.EndTime(),
		capacity:                          problem.NewCapacity(make([]int, 1)),
	}

	res.SetIndex(-2)
	return res
}
