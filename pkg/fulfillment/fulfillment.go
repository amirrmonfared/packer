package fulfillment

import (
	"math"
	"sort"
)

type ShipmentPlan struct {
	ExtraUnits int
	TotalBoxes int
	BoxCounts  map[int]int // boxSize -> count of that box
}

type shipmentEntry struct {
	totalBoxes int
	prevBox    int
}

func CalculateShipmentPlan(orderUnits int, boxSizes []int) ShipmentPlan {
	if orderUnits <= 0 || len(boxSizes) == 0 {
		return ShipmentPlan{
			ExtraUnits: -1,
			BoxCounts:  map[int]int{},
		}
	}

	// Sort from smallest to largest
	sort.Ints(boxSizes)
	maxBox := boxSizes[len(boxSizes)-1]

	// We'll only compute up to (orderUnits + largest box),
	// because overshooting more than maxBox is unnecessary.
	limit := orderUnits + maxBox

	// dp[x] = best way to form exactly x units
	dp := make([]*shipmentEntry, limit+1)
	dp[0] = &shipmentEntry{totalBoxes: 0, prevBox: -1}

	for x := 1; x <= limit; x++ {
		best := &shipmentEntry{totalBoxes: math.MaxInt32, prevBox: -1}
		for _, size := range boxSizes {
			if size > x {
				break
			}
			if dp[x-size] != nil {
				candidate := dp[x-size].totalBoxes + 1
				if candidate < best.totalBoxes {
					best.totalBoxes = candidate
					best.prevBox = size
				}
			}
		}
		if best.prevBox != -1 {
			dp[x] = best
		}
	}

	bestPlan := ShipmentPlan{
		ExtraUnits: math.MaxInt32,
		TotalBoxes: math.MaxInt32,
		BoxCounts:  map[int]int{},
	}
	for x := orderUnits; x <= limit; x++ {
		if dp[x] == nil {
			continue
		}
		overshoot := x - orderUnits
		if overshoot < bestPlan.ExtraUnits {
			bestPlan = reconstructPlan(x, overshoot, dp)
		} else if overshoot == bestPlan.ExtraUnits {
			temp := reconstructPlan(x, overshoot, dp)
			if temp.TotalBoxes < bestPlan.TotalBoxes {
				bestPlan = temp
			}
		}
		if overshoot == 0 {
			break
		}
	}
	return bestPlan
}

func reconstructPlan(x, overshoot int, dp []*shipmentEntry) ShipmentPlan {
	boxCounts := make(map[int]int)
	total := 0
	cur := x

	for cur > 0 {
		entry := dp[cur]
		if entry == nil || entry.prevBox == -1 {
			break
		}
		boxCounts[entry.prevBox]++
		total++
		cur -= entry.prevBox
	}

	return ShipmentPlan{
		ExtraUnits: overshoot,
		TotalBoxes: total,
		BoxCounts:  boxCounts,
	}
}
