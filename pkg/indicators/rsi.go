package indicators

import (
	"swing-trader/internal/types"
)

// CalculateRSI calculates the Relative Strength Index for given stock data
func CalculateRSI(data []types.StockData, period int) []float64 {
	if len(data) < period+1 {
		return make([]float64, len(data))
	}

	rsiValues := make([]float64, len(data))
	gains := make([]float64, len(data))
	losses := make([]float64, len(data))

	// Calculate price changes
	for i := 1; i < len(data); i++ {
		change := data[i].Close - data[i-1].Close
		if change > 0 {
			gains[i] = change
			losses[i] = 0
		} else {
			gains[i] = 0
			losses[i] = -change
		}
	}

	// Calculate initial average gain and loss
	var avgGain, avgLoss float64
	for i := 1; i <= period; i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	// Calculate RSI for the first valid point
	if avgLoss == 0 {
		rsiValues[period] = 100
	} else {
		rs := avgGain / avgLoss
		rsiValues[period] = 100 - (100 / (1 + rs))
	}

	// Calculate RSI for subsequent points using smoothed averages
	for i := period + 1; i < len(data); i++ {
		avgGain = (avgGain*float64(period-1) + gains[i]) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + losses[i]) / float64(period)

		if avgLoss == 0 {
			rsiValues[i] = 100
		} else {
			rs := avgGain / avgLoss
			rsiValues[i] = 100 - (100 / (1 + rs))
		}
	}

	return rsiValues
}
