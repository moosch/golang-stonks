package indicators

import (
    "math"
    "swing-trader/internal/types"
)

// CalculateBollingerBands calculates the Bollinger Bands for given stock data
func CalculateBollingerBands(data []types.StockData, period int, stdDevMultiplier float64) (bands []types.BollingerBands) {
    for i := range data {
        sum := 0.0
        sqSum := 0.0
        
        if i >= period-1 {
            for j := 0; j < period; j++ {
                sum += data[i-j].Close
                sqSum += math.Pow(data[i-j].Close, 2)
            }

            // Calculate mean
            mean := sum / float64(period)

            // Calculate standard deviation
            variance := (sqSum / float64(period)) - math.Pow(mean, 2)
            stdDev := math.Sqrt(variance)

            // Append the Bollinger Bands for this point
            upper := mean + (stdDevMultiplier * stdDev)
            lower := mean - (stdDevMultiplier * stdDev)
            bands = append(bands, types.BollingerBands{
                Upper:  upper,
                Middle: mean,
                Lower:  lower,
            })
        } else {
            // Append nil for the first points where the period is not reached
            bands = append(bands, types.BollingerBands{})
        }
    }

    return bands
}

