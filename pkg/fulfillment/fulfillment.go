package fulfillment

import (
	"math"
	"sort"
)

// ShipmentPlan represents the final outcome of the packing algorithm.
// ExtraUnits is how many items we overshoot (if positive)
// or how many items short we are (if negative) but in this logic,
// we only overshoot, so it's generally >= 0.
// TotalBoxes is how many total boxes were used.
// BoxCounts maps each boxSize -> the number of that box used.
type ShipmentPlan struct {
	ExtraUnits int
	TotalBoxes int
	BoxCounts  map[int]int // boxSize -> count of that box
}

// shipmentEntry is our DP state storage:
// totalBoxes stores how many boxes were needed to reach exactly x units.
// prevBox tells us which box size was used last to form x units.
type shipmentEntry struct {
	totalBoxes int
	prevBox    int
}

// CalculateShipmentPlan computes how to fulfill the given orderUnits using
// the provided boxSizes. The algorithm strictly uses Dynamic Programming
// to guarantee the optimal result:
//
// 1) Minimize leftover (ExtraUnits).
// 2) Among solutions with the same leftover, minimize the number of boxes.
//
// This means for large inputs, it may be more memory/CPU intensive,
// but you'll always get the best combination that meets the requirements.
func CalculateShipmentPlan(orderUnits int, boxSizes []int) ShipmentPlan {
	// Quick sanity checks:
	// If orderUnits is zero or negative, or we have no boxSizes,
	// return a plan indicating an invalid request (ExtraUnits = -1).
	if orderUnits <= 0 || len(boxSizes) == 0 {
		return ShipmentPlan{
			ExtraUnits: -1,
			BoxCounts:  map[int]int{},
		}
	}

	// Sort boxSizes ascending so we can iterate from smallest to largest if needed.
	sort.Ints(boxSizes)
	maxBox := boxSizes[len(boxSizes)-1]

	// We'll create a DP array up to (orderUnits + maxBox) to allow for overshoot up to maxBox.
	// Overshooting more than maxBox wouldn't produce fewer leftover items anyway.
	limit := orderUnits + maxBox

	// dp[x] will point to the "best" way (fewest boxes) to form exactly x units.
	// If dp[x] is nil, we haven't found a way to form x units yet.
	dp := make([]*shipmentEntry, limit+1)

	// Base case: dp[0] => 0 boxes needed, no previous box.
	dp[0] = &shipmentEntry{totalBoxes: 0, prevBox: -1}

	// Build up the table from 1..limit.
	for x := 1; x <= limit; x++ {
		best := &shipmentEntry{totalBoxes: math.MaxInt32, prevBox: -1}
		// For each boxSize, check if we can form x-size, then add 1 box.
		for _, size := range boxSizes {
			if size > x {
				break // No need to continue if box size is bigger than x.
			}
			if dp[x-size] != nil {
				candidate := dp[x-size].totalBoxes + 1
				// Keep track of the option that uses fewer total boxes.
				if candidate < best.totalBoxes {
					best.totalBoxes = candidate
					best.prevBox = size
				}
			}
		}
		// If we found at least one valid way to form x, store it in dp[x].
		if best.prevBox != -1 {
			dp[x] = best
		}
	}

	// Now, we want to pick the best plan among dp[orderUnits], dp[orderUnits+1], ... dp[orderUnits+maxBox].
	// "Best" means minimal leftover first (overshoot = x - orderUnits), then minimal total boxes.
	bestPlan := ShipmentPlan{
		ExtraUnits: math.MaxInt32, // Start with a large leftover
		TotalBoxes: math.MaxInt32,
		BoxCounts:  map[int]int{},
	}

	// Check each possible x from orderUnits up to limit.
	for x := orderUnits; x <= limit; x++ {
		if dp[x] == nil {
			continue // No way to form x exactly, skip.
		}
		overshoot := x - orderUnits
		// If the leftover is smaller, take it.
		if overshoot < bestPlan.ExtraUnits {
			bestPlan = reconstructPlan(x, overshoot, dp)
		} else if overshoot == bestPlan.ExtraUnits {
			// If leftover is the same, choose the plan with fewer boxes.
			temp := reconstructPlan(x, overshoot, dp)
			if temp.TotalBoxes < bestPlan.TotalBoxes {
				bestPlan = temp
			}
		}
		// If overshoot == 0, we've found a perfect match. Can't do better leftover than zero.
		if overshoot == 0 {
			break
		}
	}

	return bestPlan
}

// reconstructPlan walks backwards through our dp array starting at x,
// using prevBox to see which box was chosen last, until we reach 0.
// This reconstructs the distribution of box counts (BoxCounts).
func reconstructPlan(x, overshoot int, dp []*shipmentEntry) ShipmentPlan {
	boxCounts := make(map[int]int)
	total := 0
	cur := x

	for cur > 0 {
		entry := dp[cur]
		// If entry is nil or prevBox = -1, we've reached an impossible state, break out.
		if entry == nil || entry.prevBox == -1 {
			break
		}
		// Increase the count for this box size.
		boxCounts[entry.prevBox]++
		total++
		// Move cur backward by the size of the last box used.
		cur -= entry.prevBox
	}

	// overshoot is how many items we exceed the exact orderUnits (or 0 if perfect).
	return ShipmentPlan{
		ExtraUnits: overshoot,
		TotalBoxes: total,
		BoxCounts:  boxCounts,
	}
}
