package activity

import (
	"fmt"
	"gsprit/problem"
)

type TimeWindow struct {
	start float64
	end   float64
}

// NewTimeWindow creates a new TimeWindow instance
func NewTimeWindow(start, end float64) (*TimeWindow, error) {
	if start < 0.0 || end < 0.0 {
		return nil, fmt.Errorf("neither time window start nor end must be < 0.0: [start=%.2f][end=%.2f]", start, end)
	}
	if end < start {
		return nil, fmt.Errorf("time window end cannot be smaller than its start: [start=%.2f][end=%.2f]", start, end)
	}
	return &TimeWindow{start: start, end: end}, nil
}

// Larger returns true if the current TimeWindow is larger than the given TimeWindow
func (tw *TimeWindow) Larger(other problem.TimeWindow) bool {
	return (tw.end - tw.start) > (other.End() - other.Start())
}

// String returns a string representation of the TimeWindow
func (tw *TimeWindow) String() string {
	return fmt.Sprintf("[start=%.2f][end=%.2f]", tw.start, tw.end)
}

func (tw *TimeWindow) Start() float64 {
	return tw.start
}

func (tw *TimeWindow) End() float64 {
	return tw.end
}

func (tw *TimeWindow) Equals(other problem.TimeWindow) bool {
	if other == nil {
		return false
	}
	return tw.start == other.Start() && tw.end == other.End()
}
