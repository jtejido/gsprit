package algorithm

import (
	"gsprit/problem"
	"gsprit/problem/cost"
	"gsprit/problem/solution/route"
	"gsprit/problem/vrp"
)

type VehicleRoutingProblem interface {
	Activities(job problem.Job) []problem.AbstractActivity
	ActivityCosts() cost.VehicleRoutingActivityCosts
	AllLocations() []*problem.Location
	FleetSize() vrp.FleetSize
	InitialVehicleRoutes() []*route.VehicleRoute
	JobActivityFactory() func(problem.Job) []problem.AbstractActivity
	Jobs() map[string]problem.Job
	JobsInclusiveInitialJobsInRoutes() map[string]problem.Job
	JobsWithLocation() []problem.Job
	NuActivities() int
	String() string
	TransportCosts() cost.VehicleRoutingTransportCosts
	Types() []problem.VehicleType
	Vehicles() []problem.Vehicle
}
