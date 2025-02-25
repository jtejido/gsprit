package algorithm

import (
	"errors"
	"fmt"
	"math/rand/v2"
)

type SearchStrategyManager struct {
	searchStrategyListeners []SearchStrategyListener
	strategies              []*SearchStrategy
	weights                 []float64
	id2index                map[string]int
	random                  *rand.Rand
	sumWeights              float64
	strategyIndex           int
}

func (m *SearchStrategyManager) SetRandom(r *rand.Rand) {
	m.random = r
}
func (m *SearchStrategyManager) Weights() []float64 {
	c := make([]float64, len(m.weights))
	copy(c, m.weights)
	return c
}

func (m *SearchStrategyManager) Weight(strategyId string) float64 {
	return m.weights[m.id2index[strategyId]]
}

func (m *SearchStrategyManager) AddStrategy(strategy *SearchStrategy, weight float64) error {
	if strategy == nil {
		return fmt.Errorf("strategy is null. make sure adding a valid strategy")
	}
	if _, exists := m.id2index[strategy.Id()]; exists {
		return fmt.Errorf("strategyId %s already in use. replace strateId in your config file or code with a unique strategy id", strategy.Id())
	}
	if weight < 0.0 {
		return fmt.Errorf("weight is lower than zero")
	}
	m.id2index[strategy.Id()] = m.strategyIndex
	m.strategyIndex++
	m.strategies = append(m.strategies, strategy)
	m.weights = append(m.weights, weight)
	m.sumWeights += weight
	return nil
}

func (m *SearchStrategyManager) InformStrategyWeightChanged(strategyId string, weight float64) error {
	index, exists := m.id2index[strategyId]
	if !exists {
		return errors.New("strategy ID not found")
	}

	m.weights[index] = weight
	m.updateSumWeights()
	return nil
}

func (m *SearchStrategyManager) updateSumWeights() {
	var sum float64
	for _, w := range m.weights {
		sum += w
	}
	m.sumWeights = sum
}

func (m *SearchStrategyManager) RandomStrategy() (*SearchStrategy, error) {
	if m.random == nil {
		return nil, fmt.Errorf("randomizer is null. make sure you set random object correctly")
	}

	if len(m.strategies) == 0 {
		return nil, fmt.Errorf("no search strategies available")
	}

	randomValue := m.random.Float64()
	cumulativeProbability := 0.0
	for i, weight := range m.weights {
		cumulativeProbability += weight / m.sumWeights
		if randomValue < cumulativeProbability {
			return m.strategies[i], nil
		}
	}

	return m.strategies[len(m.strategies)-1], nil
}

func (m *SearchStrategyManager) AddSearchStrategyListener(strategyListener SearchStrategyListener) {
	m.searchStrategyListeners = append(m.searchStrategyListeners, strategyListener)
}

func (m *SearchStrategyManager) AddSearchStrategyModuleListener(listener SearchStrategyModuleListener) {
	for _, strategy := range m.strategies {
		strategy.AddModuleListener(listener)
	}
}
