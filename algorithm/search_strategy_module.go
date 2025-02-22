package algorithm

import (
	"gsprit/algorithm/listener"
	"gsprit/problem/solution"
)

type SearchStrategyModule interface {
	RunAndGetSolution() *solution.VehicleRoutingProblemSolution
	Name() string
	AddModuleListener(moduleListener listener.SearchStrategyModuleListener)
}
