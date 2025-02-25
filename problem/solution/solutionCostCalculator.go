package solution

type SolutionCostCalculator interface {
	Costs(solution *VehicleRoutingProblemSolution) float64
}
