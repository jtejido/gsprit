package activity

import (
	"fmt"
	"gsprit/problem"
)

type TimeWindows interface {
	Add(timeWindow problem.TimeWindow) error
	TimeWindows() []problem.TimeWindow
}

type TimeWindowsImpl struct {
	timeWindows []problem.TimeWindow
}

func NewTimeWindows() *TimeWindowsImpl {
	return &TimeWindowsImpl{timeWindows: []problem.TimeWindow{}}
}

func (tw *TimeWindowsImpl) Add(timeWindow problem.TimeWindow) error {
	for _, existingTW := range tw.timeWindows {
		if (timeWindow.Start() > existingTW.Start() && timeWindow.Start() < existingTW.End()) ||
			(timeWindow.End() > existingTW.Start() && timeWindow.End() < existingTW.End()) ||
			(timeWindow.Start() <= existingTW.Start() && timeWindow.End() >= existingTW.End()) {
			return fmt.Errorf("time-windows cannot overlap each other. overlap: %v, %v", existingTW, timeWindow)
		}
	}
	tw.timeWindows = append(tw.timeWindows, timeWindow)
	return nil
}

func (tw *TimeWindowsImpl) TimeWindows() []problem.TimeWindow {
	return tw.timeWindows
}

func (tw *TimeWindowsImpl) String() string {
	result := ""
	for _, t := range tw.timeWindows {
		result += fmt.Sprintf("[timeWindow=%v]", t)
	}
	return result
}
