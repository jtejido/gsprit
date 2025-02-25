package algorithm

import (
	"gsprit/problem/solution"
)

type SearchStrategyModule interface {
	RunAndGetSolution(*solution.VehicleRoutingProblemSolution) *solution.VehicleRoutingProblemSolution
	Name() string
	AddModuleListener(moduleListener SearchStrategyModuleListener)
}
