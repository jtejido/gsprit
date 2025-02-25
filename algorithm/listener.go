package algorithm

import (
	"gsprit/problem/solution"
	"gsprit/problem/vrp"
)

type SearchStrategyModuleListener interface {
	VehicleRoutingAlgorithmListener
}

type VehicleRoutingAlgorithmListener interface {
}

type SearchStrategyListener interface {
	VehicleRoutingAlgorithmListener
}

type AlgorithmEndsListener interface {
	VehicleRoutingAlgorithmListener
	InformAlgorithmEnds(problem *vrp.VehicleRoutingProblem, solutions []*solution.VehicleRoutingProblemSolution)
}

type AlgorithmStartsListener interface {
	VehicleRoutingAlgorithmListener
	InformAlgorithmStarts(problem *vrp.VehicleRoutingProblem, algorithm VehicleRoutingAlgorithm, solutions []*solution.VehicleRoutingProblemSolution)
}

type IterationStartsListener interface {
	VehicleRoutingAlgorithmListener
	InformIterationStarts(i int, problem *vrp.VehicleRoutingProblem, solutions []*solution.VehicleRoutingProblemSolution)
}

type IterationEndsListener interface {
	VehicleRoutingAlgorithmListener
	InformIterationEnds(i int, problem *vrp.VehicleRoutingProblem, solutions []*solution.VehicleRoutingProblemSolution)
}

type StrategySelectedListener interface {
	InformSelectedStrategy(discoveredSolution *DiscoveredSolution, vehicleRoutingProblem *vrp.VehicleRoutingProblem, vehicleRoutingProblemSolutions []*(solution.VehicleRoutingProblemSolution))
}

type VehicleRoutingAlgorithm interface {
	SearchSolutions() ([]*solution.VehicleRoutingProblemSolution, error)
}
