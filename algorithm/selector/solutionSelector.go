package selector

import (
	"gsprit/problem/solution"
)

type SolutionSelector interface {
	SelectSolution(solutions []*solution.VehicleRoutingProblemSolution) *solution.VehicleRoutingProblemSolution
}
