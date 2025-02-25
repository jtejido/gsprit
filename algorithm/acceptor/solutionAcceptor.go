package acceptor

import "gsprit/problem/solution"

type SolutionAcceptor interface {
	AcceptSolution(solutions []*solution.VehicleRoutingProblemSolution, newSolution *solution.VehicleRoutingProblemSolution) bool
}
