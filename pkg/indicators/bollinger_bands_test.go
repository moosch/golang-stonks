package indicators

import (
	"math"
	"swing-trader/internal/types"
	"testing"
	"time"
)

func TestCalculateBollingerBands(t *testing.T) {
	// Create test data with known values
	testData := []types.StockData{
		{Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Close: 100.0},
		{Date: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), Close: 102.0},
		{Date: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC), Close: 101.0},
		{Date: time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC), Close: 103.0},
		{Date: time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC), Close: 105.0},
	}

	period := 3
	stdDevMultiplier := 2.0

	bands := CalculateBollingerBands(testData, period, stdDevMultiplier)

	// First two points should be empty (not enough data)
	if bands[0].Middle != 0 || bands[1].Middle != 0 {
		t.Errorf("Expected first two points to be zero, got %v, %v", bands[0], bands[1])
	}

	// Test the third point (index 2)
	// Mean of 100, 102, 101 = 101
	expectedMean := 101.0
	if math.Abs(bands[2].Middle-expectedMean) > 0.001 {
		t.Errorf("Expected middle band at index 2 to be %f, got %f", expectedMean, bands[2].Middle)
	}

	// Standard deviation calculation
	// Variance = ((100-101)^2 + (102-101)^2 + (101-101)^2) / 3 = (1 + 1 + 0) / 3 = 2/3
	// StdDev = sqrt(2/3) ≈ 0.816
	expectedStdDev := math.Sqrt(2.0 / 3.0)
	expectedUpper := expectedMean + (stdDevMultiplier * expectedStdDev)
	expectedLower := expectedMean - (stdDevMultiplier * expectedStdDev)

	if math.Abs(bands[2].Upper-expectedUpper) > 0.001 {
		t.Errorf("Expected upper band at index 2 to be %f, got %f", expectedUpper, bands[2].Upper)
	}

	if math.Abs(bands[2].Lower-expectedLower) > 0.001 {
		t.Errorf("Expected lower band at index 2 to be %f, got %f", expectedLower, bands[2].Lower)
	}
}

func TestCalculateBollingerBandsEmptyData(t *testing.T) {
	bands := CalculateBollingerBands([]types.StockData{}, 20, 2.0)
	if len(bands) != 0 {
		t.Errorf("Expected empty result for empty data, got %d bands", len(bands))
	}
}

func TestCalculateBollingerBandsInsufficientData(t *testing.T) {
	testData := []types.StockData{
		{Close: 100.0},
		{Close: 102.0},
	}

	bands := CalculateBollingerBands(testData, 5, 2.0)
	
	// Should return bands for each data point, but all should be zero
	if len(bands) != len(testData) {
		t.Errorf("Expected %d bands, got %d", len(testData), len(bands))
	}

	for i, band := range bands {
		if band.Middle != 0 || band.Upper != 0 || band.Lower != 0 {
			t.Errorf("Expected zero values for insufficient data at index %d, got %v", i, band)
		}
	}
}
