package indicators

import (
	"swing-trader/internal/types"
	"testing"
	"time"
)

func TestCalculateRSI(t *testing.T) {
	testData := []types.StockData{
		{Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Close: 100.0},
		{Date: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), Close: 101.0},
		{Date: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC), Close: 102.0},
		{Date: time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC), Close: 103.0},
		{Date: time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC), Close: 104.0},
		{Date: time.Date(2023, 1, 6, 0, 0, 0, 0, time.UTC), Close: 105.0},
		{Date: time.Date(2023, 1, 7, 0, 0, 0, 0, time.UTC), Close: 106.0},
		{Date: time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC), Close: 107.0},
		{Date: time.Date(2023, 1, 9, 0, 0, 0, 0, time.UTC), Close: 108.0},
		{Date: time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC), Close: 109.0},
		{Date: time.Date(2023, 1, 11, 0, 0, 0, 0, time.UTC), Close: 110.0},
		{Date: time.Date(2023, 1, 12, 0, 0, 0, 0, time.UTC), Close: 111.0},
		{Date: time.Date(2023, 1, 13, 0, 0, 0, 0, time.UTC), Close: 112.0},
		{Date: time.Date(2023, 1, 14, 0, 0, 0, 0, time.UTC), Close: 113.0},
		{Date: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), Close: 114.0},
	}

	period := 14
	rsi := CalculateRSI(testData, period)

	// Check the length of RSI values
	if len(rsi) != len(testData) {
		t.Errorf("Expected RSI length %d, got %d", len(testData), len(rsi))
	}

	// Check the last RSI value
	expectedRSI := 100.0 // Upward trend
	lastRSI := rsi[len(rsi)-1]

	if lastRSI != expectedRSI {
		t.Errorf("Expected last RSI to be %.2f, got %.2f", expectedRSI, lastRSI)
	}
}
