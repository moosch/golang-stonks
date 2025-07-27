package strategy

import (
	"swing-trader/internal/types"
	"swing-trader/pkg/indicators"
)

// BBRSIStrategy implements the Bollinger Bands + RSI strategy
type BBRSIStrategy struct {
	config types.StrategyConfig
}

// NewBBRSIStrategy creates a new Bollinger Bands + RSI strategy
func NewBBRSIStrategy(config types.StrategyConfig) *BBRSIStrategy {
	return &BBRSIStrategy{
		config: config,
	}
}

// GenerateSignals generates buy/sell signals based on Bollinger Bands and RSI
func (s *BBRSIStrategy) GenerateSignals(data []types.StockData) []types.Signal {
	if len(data) < s.config.BBPeriod || len(data) < s.config.RSIPeriod {
		return []types.Signal{}
	}

	// Calculate indicators
	bollingerBands := indicators.CalculateBollingerBands(data, s.config.BBPeriod, s.config.BBStdDev)
	rsiValues := indicators.CalculateRSI(data, s.config.RSIPeriod)

	var signals []types.Signal
	
	// Start from the maximum of the two periods to ensure both indicators are valid
	startIndex := s.config.BBPeriod
	if s.config.RSIPeriod > s.config.BBPeriod {
		startIndex = s.config.RSIPeriod
	}

	for i := startIndex; i < len(data); i++ {
		signal := s.evaluatePosition(data[i], bollingerBands[i], rsiValues[i])
		if signal.Type != "HOLD" {
			signals = append(signals, signal)
		}
	}

	return signals
}

// evaluatePosition evaluates whether to buy, sell, or hold based on current conditions
func (s *BBRSIStrategy) evaluatePosition(stockData types.StockData, bb types.BollingerBands, rsi float64) types.Signal {
	signal := types.Signal{
		Date:  stockData.Date,
		Price: stockData.Close,
		Type:  "HOLD",
	}

	// Buy signal: price is below lower Bollinger Band AND RSI is below buy threshold
	if stockData.Close < bb.Lower && rsi < s.config.BuyThreshold {
		signal.Type = "BUY"
		signal.Reason = "Price below lower BB and RSI oversold"
		return signal
	}

	// Sell signal: RSI is above sell threshold (overbought)
	if rsi > s.config.SellThreshold {
		signal.Type = "SELL"
		signal.Reason = "RSI overbought"
		return signal
	}

	return signal
}

// CalculatePositionSize calculates the number of shares to buy based on available capital and risk management
func (s *BBRSIStrategy) CalculatePositionSize(availableCapital, currentPrice float64, riskConfig types.RiskManagementConfig) int64 {
	// Calculate position size based on risk percentage
	riskAmount := availableCapital * riskConfig.PositionSize
	
	// Calculate shares based on stop loss risk
	stopLossPrice := currentPrice * (1 - s.config.StopLoss)
	riskPerShare := currentPrice - stopLossPrice
	
	if riskPerShare <= 0 {
		return 0
	}
	
	shares := int64(riskAmount / riskPerShare)
	
	// Ensure we don't exceed available capital
	totalCost := float64(shares) * currentPrice
	if totalCost > availableCapital {
		shares = int64(availableCapital / currentPrice)
	}
	
	return shares
}

// GetStopLossPrice calculates the stop loss price for a given entry price
func (s *BBRSIStrategy) GetStopLossPrice(entryPrice float64) float64 {
	return entryPrice * (1 - s.config.StopLoss)
}

// GetTakeProfitPrice calculates the take profit price for a given entry price
func (s *BBRSIStrategy) GetTakeProfitPrice(entryPrice float64) float64 {
	return entryPrice * (1 + s.config.TakeProfit)
}
