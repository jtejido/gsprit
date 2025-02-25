package solution

import (
	"fmt"
	"gsprit/problem"
	"gsprit/problem/solution/route"
)

type VehicleRoutingProblemSolution struct {
	routes         []*route.VehicleRoute
	unassignedJobs []problem.Job
	cost           float64
}

// copyOf creates a deep copy of a given solution
func (solution *VehicleRoutingProblemSolution) Copy() *VehicleRoutingProblemSolution {
	var newRoutes []*route.VehicleRoute
	for _, route := range solution.routes {
		newRoutes = append(newRoutes, route.Copy())
	}

	newUnassignedJobs := make([]problem.Job, len(solution.unassignedJobs))
	copy(newUnassignedJobs, solution.unassignedJobs)

	return &VehicleRoutingProblemSolution{
		routes:         newRoutes,
		unassignedJobs: newUnassignedJobs,
		cost:           solution.cost,
	}
}

// newVehicleRoutingProblemSolution creates a new solution with given routes and cost
func NewVehicleRoutingProblemSolution(routes []*route.VehicleRoute, cost float64) *VehicleRoutingProblemSolution {
	return &VehicleRoutingProblemSolution{
		routes: routes,
		cost:   cost,
	}
}

// newVehicleRoutingProblemSolutionWithJobs creates a new solution with routes, unassigned jobs, and cost
func NewVehicleRoutingProblemSolutionWithJobs(routes []*route.VehicleRoute, unassignedJobs []problem.Job, cost float64) *VehicleRoutingProblemSolution {
	return &VehicleRoutingProblemSolution{
		routes:         routes,
		unassignedJobs: unassignedJobs,
		cost:           cost,
	}
}

// getRoutes returns the collection of vehicle routes
func (v *VehicleRoutingProblemSolution) Routes() []*route.VehicleRoute {
	return v.routes
}

// getCost returns the cost of the solution
func (v *VehicleRoutingProblemSolution) Cost() float64 {
	return v.cost
}

// setCost sets the cost of the solution
func (v *VehicleRoutingProblemSolution) SetCost(cost float64) {
	v.cost = cost
}

// getUnassignedJobs returns jobs that are not assigned to any vehicle route
func (v *VehicleRoutingProblemSolution) UnassignedJobs() []problem.Job {
	return v.unassignedJobs
}

// getUnassignedJobs returns jobs that are not assigned to any vehicle route
func (v *VehicleRoutingProblemSolution) SetUnassignedJobs(jobs []problem.Job) {
	v.unassignedJobs = jobs
}

// String returns a string representation of the solution
func (v *VehicleRoutingProblemSolution) String() string {
	return fmt.Sprintf("[cost=%.2f][routes=%d][unassigned=%d]", v.cost, len(v.routes), len(v.unassignedJobs))
}

func BestOf(solutions []*VehicleRoutingProblemSolution) *VehicleRoutingProblemSolution {
	var best *VehicleRoutingProblemSolution
	for _, s := range solutions {
		if best == nil {
			best = s
		} else if s.Cost() < best.Cost() {
			best = s
		}
	}
	return best
}
