package packs

import (
	"errors"
	"math"
	"sort"
)

type Result struct {
	ItemsOrdered int         `json:"items_ordered"`
	ItemsShipped int         `json:"items_shipped"`
	Packs        map[int]int `json:"packs"`       // packSize -> count
	TotalPacks   int         `json:"total_packs"` // sum of counts
}

// Solve calculates the optimal pack distribution:
// 1) cannot break packs
// 2) ship minimal items >= itemsOrdered
// 3) within that, use the fewest packs
func Solve(itemsOrdered int, packSizes []int) (Result, error) {
	if itemsOrdered <= 0 {
		return Result{}, errors.New("items must be > 0")
	}

	sizes, err := normalisePackSizes(packSizes)
	if err != nil {
		return Result{}, err
	}

	maxPack := sizes[len(sizes)-1]
	upper := itemsOrdered + maxPack

	// dp[x] = min number of packs to reach exactly x items
	const inf = math.MaxInt / 4
	dp := make([]int, upper+1)
	prev := make([]int, upper+1) // backpointer: previous sum
	used := make([]int, upper+1) // pack size used to reach sum

	for i := range dp {
		dp[i] = inf
		prev[i] = -1
		used[i] = 0
	}
	dp[0] = 0

	for sum := 1; sum <= upper; sum++ {
		for _, p := range sizes {
			if sum-p < 0 {
				continue
			}
			if dp[sum-p]+1 < dp[sum] {
				dp[sum] = dp[sum-p] + 1
				prev[sum] = sum - p
				used[sum] = p
			}
		}
	}

	// Find minimal shipped total >= itemsOrdered that is reachable.
	bestTotal := -1
	for t := itemsOrdered; t <= upper; t++ {
		if dp[t] < inf {
			bestTotal = t
			break // minimal shipped items first
		}
	}
	if bestTotal == -1 {
		return Result{}, errors.New("no solution found for given pack sizes")
	}

	// Reconstruct pack counts.
	counts := make(map[int]int)
	cur := bestTotal
	for cur > 0 {
		p := used[cur]
		if p == 0 || prev[cur] == -1 {
			// Shouldn't happen if dp is correct, but guard anyway.
			return Result{}, errors.New("failed to reconstruct solution")
		}
		counts[p]++
		cur = prev[cur]
	}

	return Result{
		ItemsOrdered: itemsOrdered,
		ItemsShipped: bestTotal,
		Packs:        counts,
		TotalPacks:   dp[bestTotal],
	}, nil
}

func normalisePackSizes(packSizes []int) ([]int, error) {
	if len(packSizes) == 0 {
		return nil, errors.New("pack_sizes must not be empty")
	}
	seen := make(map[int]struct{}, len(packSizes))
	sizes := make([]int, 0, len(packSizes))

	for _, p := range packSizes {
		if p <= 0 {
			return nil, errors.New("pack_sizes must all be > 0")
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		sizes = append(sizes, p)
	}

	sort.Ints(sizes)
	return sizes, nil
}
