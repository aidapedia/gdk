package time

import (
	"testing"
	"time"
)

func TestTime_NumberMonthInRange(t *testing.T) {
	// testing code
	start := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)
	result := NumberMonthInRange(start, end)
	if result != 2 {
		t.Errorf("NumberMonthInRange() = %v, want %v", result, 2)
	}
}
