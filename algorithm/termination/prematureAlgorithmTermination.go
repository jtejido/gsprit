package termination

import "gsprit/algorithm"

type PrematureAlgorithmTermination interface {
	IsPrematureBreak(discoveredSolution *algorithm.DiscoveredSolution) bool
}
