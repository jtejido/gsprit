package algorithm

import (
	"fmt"
	"gsprit/algorithm/acceptor"
	"gsprit/algorithm/selector"
	"gsprit/problem/solution"
	"gsprit/problem/vrp"
)

type DiscoveredSolution struct {
	solution   *solution.VehicleRoutingProblemSolution
	accepted   bool
	strategyId string
}

func (d *DiscoveredSolution) Solution() *solution.VehicleRoutingProblemSolution {
	return d.solution
}

func newDiscoveredSolution(solution *solution.VehicleRoutingProblemSolution, accepted bool, strategyId string) *DiscoveredSolution {
	return &DiscoveredSolution{
		solution:   solution,
		accepted:   accepted,
		strategyId: strategyId,
	}
}

func (s *DiscoveredSolution) String() string {
	return fmt.Sprintf("[strategyId=%s][solution=%v][accepted=%v]", s.strategyId, s.solution, s.accepted)
}

type SearchStrategy struct {
	searchStrategyModules  []SearchStrategyModule
	solutionSelector       selector.SolutionSelector
	solutionCostCalculator solution.SolutionCostCalculator
	solutionAcceptor       acceptor.SolutionAcceptor
	id                     string
	name                   string
}

func NewSearchStrategy(id string, solutionSelector selector.SolutionSelector, solutionAcceptor acceptor.SolutionAcceptor, solutionCostCalculator solution.SolutionCostCalculator) *SearchStrategy {
	return &SearchStrategy{
		id:                     id,
		solutionSelector:       solutionSelector,
		solutionAcceptor:       solutionAcceptor,
		solutionCostCalculator: solutionCostCalculator,
	}
}

func (s *SearchStrategy) Id() string {
	return s.id
}

func (s *SearchStrategy) Name() string {
	return s.name
}

func (s *SearchStrategy) SetName(n string) {
	s.name = n
}

func (s *SearchStrategy) SearchStrategyModules() []SearchStrategyModule {
	c := make([]SearchStrategyModule, len(s.searchStrategyModules))
	copy(c, s.searchStrategyModules)
	return c
}

func (s *SearchStrategy) String() string {
	return fmt.Sprintf("searchStrategy [#modules=%d][selector=%v][acceptor=%v]", len(s.searchStrategyModules), s.solutionSelector, s.solutionAcceptor)
}

func (s *SearchStrategy) Run(vrp *vrp.VehicleRoutingProblem, solutions []*solution.VehicleRoutingProblemSolution) (*DiscoveredSolution, error) {
	solution := s.solutionSelector.SelectSolution(solutions)
	if solution == nil {
		return nil, fmt.Errorf("solution is nil. check solutionSelector to return an appropriate solution. " +
			"figure out whether you start with an initial solution. either you set it manually by algorithm.AddInitialSolution(...)" +
			" or let the algorithm create an initial solution for you. then add the <construction>...</construction> xml-snippet to your algorithm's config file")
	}
	lastSolution := solution.Copy()
	for _, module := range s.searchStrategyModules {
		lastSolution = module.RunAndGetSolution(lastSolution)
	}
	costs := s.solutionCostCalculator.Costs(lastSolution)
	lastSolution.SetCost(costs)
	solutionAccepted := s.solutionAcceptor.AcceptSolution(solutions, lastSolution)
	return newDiscoveredSolution(lastSolution, solutionAccepted, s.Id()), nil
}

func (s *SearchStrategy) AddModule(module SearchStrategyModule) error {
	if module == nil {
		return fmt.Errorf("module to be added is null")
	}
	s.searchStrategyModules = append(s.searchStrategyModules, module)
	return nil
}

func (s *SearchStrategy) AddModuleListener(moduleListener SearchStrategyModuleListener) {
	for _, module := range s.searchStrategyModules {
		module.AddModuleListener(moduleListener)
	}
}
