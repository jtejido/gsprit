package vra

import (
	"container/heap"
	"gsprit/algorithm"
	"gsprit/problem/solution"
	"gsprit/problem/vrp"
)

type Priority int

const (
	HIGH Priority = iota
	MEDIUM
	LOW
)

type PrioritizedVRAListener struct {
	priority Priority
	l        algorithm.VehicleRoutingAlgorithmListener
	index    int
}

func newPrioritizedVRAListener(p Priority, l algorithm.VehicleRoutingAlgorithmListener) *PrioritizedVRAListener {
	return &PrioritizedVRAListener{
		priority: p,
		l:        l,
	}
}

type priorityQueue []*PrioritizedVRAListener

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	if pq[i].priority == pq[j].priority {
		return i < j
	}
	return pq[i].priority < pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PrioritizedVRAListener)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

type treeSet struct {
	pq priorityQueue
}

func newTreeSet() *treeSet {
	ts := &treeSet{}
	heap.Init(&ts.pq)
	return ts
}

func (ts *treeSet) Add(listener *PrioritizedVRAListener) {
	heap.Push(&ts.pq, listener)
}

func (ts *treeSet) Get() *PrioritizedVRAListener {
	if ts.pq.Len() == 0 {
		return nil
	}
	return heap.Pop(&ts.pq).(*PrioritizedVRAListener)
}

func (ts *treeSet) Remove(listener *PrioritizedVRAListener) {
	for i, item := range ts.pq {
		if item == listener {
			heap.Remove(&ts.pq, i)
		}
	}
}

func (ts *treeSet) All() []*PrioritizedVRAListener {
	listeners := make([]*PrioritizedVRAListener, len(ts.pq))
	copy(listeners, ts.pq)
	return listeners
}

type VehicleRoutingAlgorithmListeners struct {
	algorithmListeners *treeSet
}

func NewVehicleRoutingAlgorithmListeners() *VehicleRoutingAlgorithmListeners {
	return &VehicleRoutingAlgorithmListeners{
		algorithmListeners: newTreeSet(),
	}
}

func (v *VehicleRoutingAlgorithmListeners) AddListener(listener algorithm.VehicleRoutingAlgorithmListener, priority Priority) {
	v.algorithmListeners.Add(newPrioritizedVRAListener(priority, listener))
}

func (v *VehicleRoutingAlgorithmListeners) AddListenerWithDefaultPriority(listener algorithm.VehicleRoutingAlgorithmListener) {
	v.AddListener(listener, LOW)
}

func (v *VehicleRoutingAlgorithmListeners) Remove(listener *PrioritizedVRAListener) {
	v.algorithmListeners.Remove(listener)
}

func (v *VehicleRoutingAlgorithmListeners) AlgorithmListeners() []algorithm.VehicleRoutingAlgorithmListener {
	c := make([]algorithm.VehicleRoutingAlgorithmListener, 0)
	for _, i := range v.algorithmListeners.All() {
		c = append(c, i.l)
	}

	return c
}

func (v *VehicleRoutingAlgorithmListeners) AlgorithmEnds(problem *vrp.VehicleRoutingProblem, solutions []*solution.VehicleRoutingProblemSolution) {
	for _, l := range v.algorithmListeners.All() {
		if listener, ok := l.l.(algorithm.AlgorithmEndsListener); ok {
			listener.InformAlgorithmEnds(problem, solutions)
		}
	}
}

func (v *VehicleRoutingAlgorithmListeners) IterationEnds(iteration int, problem *vrp.VehicleRoutingProblem, solutions []*solution.VehicleRoutingProblemSolution) {
	for _, l := range v.algorithmListeners.All() {
		if listener, ok := l.l.(algorithm.IterationEndsListener); ok {
			listener.InformIterationEnds(iteration, problem, solutions)
		}
	}
}

func (v *VehicleRoutingAlgorithmListeners) IterationStarts(iteration int, problem *vrp.VehicleRoutingProblem, solutions []*solution.VehicleRoutingProblemSolution) {
	for _, l := range v.algorithmListeners.All() {
		if listener, ok := l.l.(algorithm.IterationStartsListener); ok {
			listener.InformIterationStarts(iteration, problem, solutions)
		}
	}
}

func (v *VehicleRoutingAlgorithmListeners) AlgorithmStarts(problem *vrp.VehicleRoutingProblem, alg *VehicleRoutingAlgorithm, solutions []*solution.VehicleRoutingProblemSolution) {
	for _, l := range v.algorithmListeners.All() {
		if listener, ok := l.l.(algorithm.AlgorithmStartsListener); ok {
			listener.InformAlgorithmStarts(problem, alg, solutions)
		}
	}
}

func (v *VehicleRoutingAlgorithmListeners) SelectedStrategy(discoveredSolution *algorithm.DiscoveredSolution, problem *vrp.VehicleRoutingProblem, solutions []*solution.VehicleRoutingProblemSolution) {
	for _, l := range v.algorithmListeners.All() {
		if listener, ok := l.l.(algorithm.StrategySelectedListener); ok {
			listener.InformSelectedStrategy(discoveredSolution, problem, solutions)
		}
	}
}

func (v *VehicleRoutingAlgorithmListeners) Add(listener *PrioritizedVRAListener) {
	v.algorithmListeners.Add(listener)
}

func (v *VehicleRoutingAlgorithmListeners) AddAll(listeners []*PrioritizedVRAListener) {
	for _, l := range listeners {
		v.algorithmListeners.Add(l)
	}
}
