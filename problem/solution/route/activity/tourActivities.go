package activity

import (
	"fmt"
	"gsprit/problem"
)

type TourActivities struct {
	tourActivities []problem.TourActivity
	jobs           map[problem.Job]bool
}

func NewTourActivities() *TourActivities {
	return &TourActivities{
		tourActivities: []problem.TourActivity{},
		jobs:           make(map[problem.Job]bool),
	}
}

func (ta *TourActivities) Activities() []problem.TourActivity {
	return ta.tourActivities
}

func (ta *TourActivities) IsEmpty() bool {
	return len(ta.tourActivities) == 0
}

func (ta *TourActivities) Jobs() []problem.Job {
	jobs := []problem.Job{}
	for job := range ta.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

func (ta *TourActivities) ServesJob(job problem.Job) bool {
	_, exists := ta.jobs[job]
	return exists
}

func (ta *TourActivities) RemoveJob(job problem.Job) bool {
	if !ta.ServesJob(job) {
		return false
	}
	delete(ta.jobs, job)
	filteredActivities := []problem.TourActivity{}
	for _, act := range ta.tourActivities {
		if ja, ok := act.(problem.JobActivity); ok && ja.Job() == job {
			continue
		}
		filteredActivities = append(filteredActivities, act)
	}
	ta.tourActivities = filteredActivities
	return true
}

func (ta *TourActivities) RemoveActivity(activity problem.TourActivity) bool {
	filteredActivities := []problem.TourActivity{}
	var jobToRemove problem.Job
	for _, act := range ta.tourActivities {
		if act == activity {
			if ja, ok := act.(problem.JobActivity); ok {
				jobToRemove = ja.Job()
			}
			continue
		}
		filteredActivities = append(filteredActivities, act)
	}
	ta.tourActivities = filteredActivities
	if jobToRemove != nil {
		ta.RemoveJob(jobToRemove)
	}
	return true
}

func (ta *TourActivities) AddActivity(index int, act problem.TourActivity) error {
	if index < 0 {
		return fmt.Errorf("insertionIndex cannot be negative")
	}
	if index > len(ta.tourActivities) {
		ta.tourActivities = append(ta.tourActivities, act)
	} else {
		ta.tourActivities = append(ta.tourActivities[:index], append([]problem.TourActivity{act}, ta.tourActivities[index:]...)...)
	}
	ta.AddJob(act)
	return nil
}

func (ta *TourActivities) AddActivityToEnd(act problem.TourActivity) error {
	for _, existingAct := range ta.tourActivities {
		if existingAct == act {
			return fmt.Errorf("activity already exists in tour")
		}
	}
	ta.tourActivities = append(ta.tourActivities, act)
	ta.AddJob(act)
	return nil
}

func (ta *TourActivities) AddJob(act problem.TourActivity) {
	if ja, ok := act.(problem.JobActivity); ok {
		ta.jobs[ja.Job()] = true
	}
}

func (ta *TourActivities) JobSize() int {
	return len(ta.jobs)
}

func (ta *TourActivities) String() string {
	return fmt.Sprintf("[nuOfActivities=%d]", len(ta.tourActivities))
}

// ReverseActivityIterator allows iteration over activities in reverse order
type ReverseActivityIterator struct {
	activities []problem.TourActivity
	index      int
}

func NewReverseActivityIterator(activities []problem.TourActivity) *ReverseActivityIterator {
	return &ReverseActivityIterator{
		activities: activities,
		index:      len(activities) - 1,
	}
}

func (it *ReverseActivityIterator) HasNext() bool {
	return it.index >= 0
}

func (it *ReverseActivityIterator) Next() problem.TourActivity {
	if !it.HasNext() {
		panic("No more elements")
	}
	act := it.activities[it.index]
	it.index--
	return act
}

func (ta *TourActivities) ReverseIterator() *ReverseActivityIterator {
	return NewReverseActivityIterator(ta.tourActivities)
}

func (ta *TourActivities) Copy() problem.TourActivities {
	res := &TourActivities{
		tourActivities: []problem.TourActivity{},
		jobs:           make(map[problem.Job]bool),
	}

	for _, tourAct := range ta.Activities() {
		newAct := tourAct.Duplicate()
		res.tourActivities = append(res.tourActivities, newAct)
		res.AddJob(newAct)
	}

	return res
}
