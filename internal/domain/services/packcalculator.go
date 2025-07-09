package services

import (
	"sort"

	"go-pack-calculator/internal/domain/errors"
)

type PackCalculatorService struct{}

func NewPackCalculatorService() *PackCalculatorService {
	return &PackCalculatorService{}
}

func (s *PackCalculatorService) CalculateOptimalPacks(itemsOrdered int, packSizes []int) (map[int]int, error) {
	if itemsOrdered <= 0 {
		return nil, errors.ErrInvalidItemsOrdered
	}
	if len(packSizes) == 0 {
		return nil, errors.ErrNoPackSizesAvailable
	}

	// Sort pack sizes in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(packSizes)))

	// Define a state structure for our dynamic programming approach
	type state struct {
		totalItems int
		packs      int
		prev       int
		used       int
	}

	// Initialize our dynamic programming table
	dp := make(map[int]state)
	dp[0] = state{totalItems: 0, packs: 0, prev: -1, used: -1}

	// Queue for BFS
	queue := []int{0}
	visited := make(map[int]bool)
	visited[0] = true

	// BFS to find all possible combinations
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// If we've reached or exceeded the target, no need to add more packs
		if current >= itemsOrdered {
			continue
		}

		currentState := dp[current]

		for _, size := range packSizes {
			next := current + size
			nextState := state{
				totalItems: currentState.totalItems + size,
				packs:      currentState.packs + 1,
				prev:       current,
				used:       size,
			}

			// Check if we've already visited this state or if it's worse than an existing solution
			if _, exists := dp[next]; !exists {
				dp[next] = nextState
				if !visited[next] {
					queue = append(queue, next)
					visited[next] = true
				}
			} else {
				// If we have a better solution (fewer total items or same items but fewer packs)
				existingState := dp[next]
				if nextState.totalItems < existingState.totalItems ||
					(nextState.totalItems == existingState.totalItems && nextState.packs < existingState.packs) {
					dp[next] = nextState
					if !visited[next] {
						queue = append(queue, next)
						visited[next] = true
					}
				}
			}
		}
	}

	// Find the best solution that meets or exceeds the order
	var bestTotal int = -1
	var bestPacks int = -1

	// First find all valid solutions (those that meet or exceed the order)
	validSolutions := make([]int, 0)
	for total := range dp {
		if total >= itemsOrdered {
			validSolutions = append(validSolutions, total)
		}
	}

	// No valid solution found
	if len(validSolutions) == 0 {
		return nil, errors.ErrInvalidPackSize
	}

	// Sort valid solutions by total items (ascending)
	sort.Ints(validSolutions)

	// Find the solution with minimum total items
	minTotalItems := dp[validSolutions[0]].totalItems
	minTotalSolutions := make([]int, 0)

	for _, total := range validSolutions {
		if dp[total].totalItems == minTotalItems {
			minTotalSolutions = append(minTotalSolutions, total)
		} else if dp[total].totalItems > minTotalItems {
			break // We've found all solutions with minimum total items
		}
	}

	// Among solutions with minimum total items, find the one with fewest packs
	bestTotal = minTotalSolutions[0]
	bestPacks = dp[bestTotal].packs

	for _, total := range minTotalSolutions {
		if dp[total].packs < bestPacks {
			bestTotal = total
			bestPacks = dp[total].packs
		}
	}

	// Backtrack to build the result
	result := make(map[int]int)
	current := bestTotal

	for current > 0 {
		st := dp[current]
		result[st.used]++
		current = st.prev
	}

	return result, nil
}
