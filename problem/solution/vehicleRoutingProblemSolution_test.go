package solution

import (
	"testing"

	"gsprit/problem"
	"gsprit/problem/solution/route"

	"github.com/stretchr/testify/assert"
)

// TestCreatingSolutionWithTwoRoutes checks if the solution contains the correct routes.
func TestCreatingSolutionWithTwoRoutes(t *testing.T) {
	r1 := &route.VehicleRoute{}
	r2 := &route.VehicleRoute{}

	sol := NewVehicleRoutingProblemSolution([]*route.VehicleRoute{r1, r2}, 0.0)
	assert.Equal(t, 2, len(sol.Routes()))
}

// TestSettingSolutionCosts checks if setting the cost updates correctly.
func TestSettingSolutionCosts(t *testing.T) {
	sol := NewVehicleRoutingProblemSolution([]*route.VehicleRoute{}, 10.0)
	assert.InDelta(t, 10.0, sol.Cost(), 0.01)
}

// TestChangingSolutionCosts checks if the cost can be changed after creation.
func TestChangingSolutionCosts(t *testing.T) {
	sol := NewVehicleRoutingProblemSolution([]*route.VehicleRoute{}, 10.0)
	sol.SetCost(20.0)
	assert.InDelta(t, 20.0, sol.Cost(), 0.01)
}

type MockJob struct {
	id         string
	index      int
	size       *problem.Capacity
	skills     *problem.Skills
	name       string
	priority   int
	maxTimeInV float64
	activities []problem.Activity
	jobType    problem.JobType
}

func (mj *MockJob) Id() string                      { return mj.id }
func (mj *MockJob) Index() int                      { return mj.index }
func (mj *MockJob) SetIndex(index int)              { mj.index = index }
func (mj *MockJob) Size() *problem.Capacity         { return mj.size }
func (mj *MockJob) RequiredSkills() *problem.Skills { return mj.skills }
func (mj *MockJob) Name() string                    { return mj.name }
func (mj *MockJob) Priority() int                   { return mj.priority }
func (mj *MockJob) MaxTimeInVehicle() float64       { return mj.maxTimeInV }
func (mj *MockJob) Activities() []problem.Activity  { return mj.activities }
func (mj *MockJob) JobType() problem.JobType        { return mj.jobType }
func (mj *MockJob) String() string                  { return mj.Id() }

// TestSizeOfBadJobs checks if unassigned jobs are counted correctly.
func TestSizeOfBadJobs(t *testing.T) {
	badJob := &MockJob{id: "bad-job-1"}
	badJobs := []problem.Job{badJob}

	sol := NewVehicleRoutingProblemSolutionWithJobs([]*route.VehicleRoute{}, badJobs, 10.0)
	assert.Equal(t, 1, len(sol.UnassignedJobs()))
}

// TestSizeOfBadJobs_2 checks if adding unassigned jobs later updates correctly.
func TestSizeOfBadJobs_2(t *testing.T) {
	badJob := &MockJob{id: "bad-job-1"}
	badJobs := []problem.Job{badJob}

	sol := NewVehicleRoutingProblemSolutionWithJobs([]*route.VehicleRoute{}, nil, 10.0)
	unassigned := sol.UnassignedJobs()
	unassigned = append(unassigned, badJobs...)

	assert.Equal(t, 1, len(unassigned))
}

// TestBadJobsCorrectness checks if unassigned jobs are retrieved correctly.
func TestBadJobsCorrectness(t *testing.T) {
	badJob := &MockJob{id: "bad-job-1"}
	badJobs := []problem.Job{badJob}

	sol := NewVehicleRoutingProblemSolutionWithJobs([]*route.VehicleRoute{}, badJobs, 10.0)
	assert.Equal(t, badJob, sol.UnassignedJobs()[0])
}

// TestBadJobsCorrectness_2 checks if adding unassigned jobs later works correctly.
func TestBadJobsCorrectness_2(t *testing.T) {
	badJob := &MockJob{id: "bad-job-1"}
	badJobs := []problem.Job{badJob}

	sol := NewVehicleRoutingProblemSolutionWithJobs([]*route.VehicleRoute{}, nil, 10.0)
	unassigned := sol.UnassignedJobs()
	unassigned = append(unassigned, badJobs...)
	assert.Equal(t, badJob, unassigned[0])
}
