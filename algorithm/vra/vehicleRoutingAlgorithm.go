package vra

import (
	"fmt"
	"gsprit/algorithm"

	"gsprit/algorithm/termination"
	"gsprit/problem"
	"gsprit/problem/solution"
	"gsprit/problem/vrp"
	"log"
	"maps"
	"sync"
	"time"
)

type TerminationManager struct {
	terminationCriteria []termination.PrematureAlgorithmTermination
}

func newTerminationManager() *TerminationManager {
	return &TerminationManager{
		terminationCriteria: make([]termination.PrematureAlgorithmTermination, 0),
	}
}

func (m *TerminationManager) AddTermination(termination termination.PrematureAlgorithmTermination) {
	m.terminationCriteria = append(m.terminationCriteria, termination)
}

func (m *TerminationManager) IsPrematureBreak(discoveredSolution *algorithm.DiscoveredSolution) bool {
	for _, termination := range m.terminationCriteria {
		if termination.IsPrematureBreak(discoveredSolution) {
			return true
		}
	}
	return false
}

type Counter struct {
	sync.Mutex
	name        string
	counter     int
	nextCounter int
}

func newCounter(name string) *Counter {
	return &Counter{
		name:        name,
		nextCounter: 1,
	}
}

func (c *Counter) IncCounter() {
	c.Lock()
	defer c.Unlock()

	c.counter++
	if c.counter >= c.nextCounter {
		c.nextCounter *= 2
	}
}

// Reset resets the counter values.
func (c *Counter) Reset() {
	c.Lock()
	defer c.Unlock()

	c.counter = 0
	c.nextCounter = 1
}

type VehicleRoutingAlgorithm struct {
	counter               *Counter
	problem               *vrp.VehicleRoutingProblem
	searchStrategyManager *algorithm.SearchStrategyManager
	algoListeners         *VehicleRoutingAlgorithmListeners
	initialSolutions      []*solution.VehicleRoutingProblemSolution
	maxIterations         int
	terminationManager    *TerminationManager
	bestEver              *solution.VehicleRoutingProblemSolution
	objectiveFunction     solution.SolutionCostCalculator
}

func NewVehicleRoutingAlgorithm(problem *vrp.VehicleRoutingProblem, searchStrategyManager *algorithm.SearchStrategyManager) *VehicleRoutingAlgorithm {
	return &VehicleRoutingAlgorithm{
		counter:               newCounter("iterations "),
		maxIterations:         100,
		terminationManager:    newTerminationManager(),
		initialSolutions:      make([]*solution.VehicleRoutingProblemSolution, 0),
		searchStrategyManager: searchStrategyManager,
	}
}

func (a *VehicleRoutingAlgorithm) AddInitialSolution(solution *solution.VehicleRoutingProblemSolution) error {
	solution = solution.Copy()
	if err := a.verifyAndAdaptSolution(solution); err != nil {
		return err
	}
	a.initialSolutions = append(a.initialSolutions, solution)
	return nil
}

func (a *VehicleRoutingAlgorithm) verifyAndAdaptSolution(solution *solution.VehicleRoutingProblemSolution) error {
	jobsNotInSolution := maps.Clone(a.problem.Jobs())
	for _, job := range solution.UnassignedJobs() {
		delete(jobsNotInSolution, job.Id())
	}

	for _, route := range solution.Routes() {
		for _, job := range route.TourActivities().Jobs() {
			delete(jobsNotInSolution, job.Id())
		}
		if route.Vehicle().Index() == 0 {
			return fmt.Errorf("vehicle used in initial solution has no index. probably a vehicle is used that has not been added to the " +
				" the VehicleRoutingProblem. only use vehicles that have already been added to the problem")
		}
		for _, act := range route.Activities() {
			if act.Index() == 0 {
				return fmt.Errorf("act in initial solution has no index. activities are created and associated to their job in VehicleRoutingProblem\n." +
					" thus if you build vehicle-routes use the jobActivityFactory from vehicle routing problem like that \n" +
					" VehicleRoute.Builder.newInstance(knownVehicle).setJobActivityFactory(vrp.getJobActivityFactory).addService(..)....build() \n" +
					" then the activities that are created to build the route are identical to the ones used in VehicleRoutingProblem")

			}
		}
	}
	unassignedJobs := solution.UnassignedJobs()
	for _, job := range jobsNotInSolution {
		unassignedJobs = append(unassignedJobs, job)
	}
	solution.SetUnassignedJobs(unassignedJobs)
	solution.SetCost(a.objectiveFunction.Costs(solution))
	return nil
}

func (a *VehicleRoutingAlgorithm) SetPrematureAlgorithmTermination(prematureAlgorithmTermination termination.PrematureAlgorithmTermination) {
	a.terminationManager = newTerminationManager()
	a.terminationManager.AddTermination(prematureAlgorithmTermination)
}

func (a *VehicleRoutingAlgorithm) AddTerminationCriterion(terminationCriterion termination.PrematureAlgorithmTermination) {
	a.terminationManager.AddTermination(terminationCriterion)
}

func (a *VehicleRoutingAlgorithm) SearchStrategyManager() *algorithm.SearchStrategyManager {
	return a.searchStrategyManager
}

func (a *VehicleRoutingAlgorithm) SearchSolutions() ([]*solution.VehicleRoutingProblemSolution, error) {
	log.Printf("algorithm starts: [maxIterations=%d]", a.maxIterations)
	now := time.Now().UnixMilli()
	noIterationsThisAlgoIsRunning := a.maxIterations
	a.counter.Reset()
	solutions := append([]*solution.VehicleRoutingProblemSolution{}, a.initialSolutions...)
	a.algorithmStarts(a.problem, solutions)
	a.bestEver = solution.BestOf(solutions)
	a.logSolutions(solutions)
	log.Printf("iterations start")
	for i := 0; i < a.maxIterations; i++ {
		a.iterationStarts(i+1, a.problem, solutions)
		log.Printf("start iteration: %d", i)
		a.counter.IncCounter()
		strategy, err := a.searchStrategyManager.RandomStrategy()
		if err != nil {
			return nil, err
		}
		discoveredSolution, err := strategy.Run(a.problem, solutions)
		if err != nil {
			return nil, err
		}
		a.logDiscoveredSolution(discoveredSolution)
		a.memorizeIfBestEver(discoveredSolution)
		a.selectedStrategy(discoveredSolution, a.problem, solutions)
		if a.terminationManager.IsPrematureBreak(discoveredSolution) {
			log.Printf("premature algorithm termination at iteration %d", (i + 1))
			noIterationsThisAlgoIsRunning = (i + 1)
			break
		}
		a.iterationEnds(i+1, a.problem, solutions)
	}
	log.Printf("iterations end at %d iterations", noIterationsThisAlgoIsRunning)
	solutions = a.addBestEver(solutions)
	a.algorithmEnds(a.problem, solutions)
	log.Printf("took %.2f seconds", (float64(time.Now().UnixMilli()-now) / 1000.0))
	return solutions, nil
}

func (a *VehicleRoutingAlgorithm) addBestEver(solutions []*solution.VehicleRoutingProblemSolution) []*solution.VehicleRoutingProblemSolution {
	if a.bestEver != nil {
		solutions = append(solutions, a.bestEver)
	}

	return solutions
}

func (a *VehicleRoutingAlgorithm) logSolutions(solutions []*solution.VehicleRoutingProblemSolution) {
	for _, sol := range solutions {
		a.logSolution(sol)
	}
}

func (a *VehicleRoutingAlgorithm) logSolution(solution *solution.VehicleRoutingProblemSolution) {
	log.Printf("solution costs: %.2f", solution.Cost())
	for _, r := range solution.Routes() {
		var b string
		b += r.Vehicle().Id()
		b += " : "
		b += "[ "

		for _, act := range r.Activities() {
			if aa, ok := act.(problem.JobActivity); ok {
				b += aa.Job().Id()
				b += " "
			}
		}
		b += "]"
		log.Print(b)
	}
	var b string
	b += "unassigned : [ "
	for _, j := range solution.UnassignedJobs() {
		b += j.Id()
		b += " "
	}
	b += "]"
	log.Print(b)
}

func (a *VehicleRoutingAlgorithm) logDiscoveredSolution(discoveredSolution *algorithm.DiscoveredSolution) {
	log.Printf("discovered solution: %v", discoveredSolution)
	a.logSolution(discoveredSolution.Solution())
}

func (a *VehicleRoutingAlgorithm) memorizeIfBestEver(discoveredSolution *algorithm.DiscoveredSolution) {
	if discoveredSolution == nil {
		return
	}
	if a.bestEver == nil {
		a.bestEver = discoveredSolution.Solution()
	} else if discoveredSolution.Solution().Cost() < a.bestEver.Cost() {
		a.bestEver = discoveredSolution.Solution()
	}
}

func (a *VehicleRoutingAlgorithm) selectedStrategy(discoveredSolution *algorithm.DiscoveredSolution, problem *vrp.VehicleRoutingProblem, solutions []*solution.VehicleRoutingProblemSolution) {
	a.algoListeners.SelectedStrategy(discoveredSolution, problem, solutions)
}

func (a *VehicleRoutingAlgorithm) algorithmEnds(problem *vrp.VehicleRoutingProblem, solutions []*solution.VehicleRoutingProblemSolution) {
	a.algoListeners.AlgorithmEnds(problem, solutions)
}

func (a *VehicleRoutingAlgorithm) AlgorithmListeners() *VehicleRoutingAlgorithmListeners {
	return a.algoListeners
}

func (a *VehicleRoutingAlgorithm) AddListener(l algorithm.VehicleRoutingAlgorithmListener) {
	a.algoListeners.AddListenerWithDefaultPriority(l)
	if ssl, ok := l.(algorithm.SearchStrategyListener); ok {
		a.searchStrategyManager.AddSearchStrategyListener(ssl)
	}
	if ssml, ok := l.(algorithm.SearchStrategyModuleListener); ok {
		a.searchStrategyManager.AddSearchStrategyModuleListener(ssml)
	}
}

func (a *VehicleRoutingAlgorithm) iterationEnds(i int, problem *vrp.VehicleRoutingProblem, solutions []*solution.VehicleRoutingProblemSolution) {
	a.algoListeners.IterationEnds(i, problem, solutions)
}

func (a *VehicleRoutingAlgorithm) iterationStarts(i int, problem *vrp.VehicleRoutingProblem, solutions []*solution.VehicleRoutingProblemSolution) {
	a.algoListeners.IterationStarts(i, problem, solutions)
}

func (a *VehicleRoutingAlgorithm) algorithmStarts(problem *vrp.VehicleRoutingProblem, solutions []*solution.VehicleRoutingProblemSolution) {
	a.algoListeners.AlgorithmStarts(problem, a, solutions)
}

func (a *VehicleRoutingAlgorithm) SetMaxIterations(maxIterations int) {
	a.maxIterations = maxIterations
	log.Printf("set maxIterations to %d", a.maxIterations)
}

func (a *VehicleRoutingAlgorithm) MaxIterations() int {
	return a.maxIterations
}

func (a *VehicleRoutingAlgorithm) ObjectiveFunction() solution.SolutionCostCalculator {
	return a.objectiveFunction
}
