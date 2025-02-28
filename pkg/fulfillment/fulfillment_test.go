package fulfillment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateShipmentPlan(t *testing.T) {
	tests := []struct {
		name          string
		orderUnits    int
		boxSizes      []int
		wantExtra     int
		wantBoxes     int
		wantBoxCounts map[int]int
	}{
		{
			name:          "Order=1",
			orderUnits:    1,
			boxSizes:      []int{250, 500, 1000, 2000, 5000},
			wantExtra:     249,
			wantBoxes:     1,
			wantBoxCounts: map[int]int{250: 1},
		},
		{
			name:          "Order=250",
			orderUnits:    250,
			boxSizes:      []int{250, 500, 1000, 2000, 5000},
			wantExtra:     0,
			wantBoxes:     1,
			wantBoxCounts: map[int]int{250: 1},
		},
		{
			name:          "Order=501",
			orderUnits:    501,
			boxSizes:      []int{250, 500, 1000, 2000, 5000},
			wantExtra:     249,
			wantBoxes:     2,
			wantBoxCounts: map[int]int{500: 1, 250: 1},
		},
		{
			name:          "EdgeCase=500000 with 23,31,53",
			orderUnits:    500000,
			boxSizes:      []int{23, 31, 53},
			wantExtra:     0,
			wantBoxes:     9438,
			wantBoxCounts: map[int]int{23: 2, 31: 7, 53: 9429},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPlan := CalculateShipmentPlan(tt.orderUnits, tt.boxSizes)

			assert.Equal(t, tt.wantExtra, gotPlan.ExtraUnits, "ExtraUnits mismatch")
			assert.Equal(t, tt.wantBoxes, gotPlan.TotalBoxes, "TotalBoxes mismatch")
			assert.Equal(t, tt.wantBoxCounts, gotPlan.BoxCounts, "BoxCounts mismatch")
		})
	}
}
