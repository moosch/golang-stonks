package types

import "time"

// StockData represents a single day's stock data
type StockData struct {
	Date          time.Time
	Open          float64
	High          float64
	Low           float64
	Close         float64
	Volume        int64
	AdjustedClose float64
}

// Trade represents a single trade with entry and exit information
type Trade struct {
	ID         string
	EntryDate  time.Time
	ExitDate   *time.Time // Pointer to handle open trades
	EntryPrice float64
	ExitPrice  *float64 // Pointer to handle open trades
	Quantity   int64
	ProfitLoss float64
	Status     string // "open", "closed", "cancelled"
	StopLoss   float64
	TakeProfit float64
}

// TradeResult provides summary statistics for a collection of trades
type TradeResult struct {
	TotalTrades     int64
	TotalProfitLoss float64
	WinningTrades   int64
	LosingTrades    int64
	WinRate         float64 // percentage of winning trades
}

// StrategyConfig holds the configuration for the trading strategy
type StrategyConfig struct {
	BuyThreshold   float64 // RSI threshold for buying (e.g., 30)
	SellThreshold  float64 // RSI threshold for selling (e.g., 70)
	StopLoss       float64 // percentage for stop loss (e.g., 0.05 for 5%)
	TakeProfit     float64 // percentage for take profit (e.g., 0.10 for 10%)
	InitialCapital float64 // starting capital for the backtest
	RSIPeriod      int     // period for RSI calculation (typically 14)
	BBPeriod       int     // period for Bollinger Bands (typically 20)
	BBStdDev       float64 // standard deviation multiplier for Bollinger Bands (typically 2.0)
}

// RiskManagementConfig holds risk management parameters
type RiskManagementConfig struct {
	MaxDrawdown  float64 // maximum drawdown percentage (e.g., 0.20 for 20%)
	PositionSize float64 // percentage of capital to risk per trade (e.g., 0.02 for 2%)
}

// BacktestResult contains comprehensive results from a backtest
type BacktestResult struct {
	Trades                    []Trade
	TotalProfitLoss          float64
	WinRate                  float64
	TotalTrades              int64
	WinningTrades            int64
	LosingTrades             int64
	AverageWin               float64
	AverageLoss              float64
	MaxDrawdown              float64
	MaxDrawdownDuration      time.Duration
	TotalReturn              float64
	AnnualizedReturn         float64
	SharpeRatio              float64
	StartDate                time.Time
	EndDate                  time.Time
	InitialCapital           float64
	FinalCapital             float64
}

// BacktestConfig holds all configuration for running a backtest
type BacktestConfig struct {
	StockDataPath        string
	StrategyConfig       StrategyConfig
	RiskManagementConfig RiskManagementConfig
	StartDate            time.Time
	EndDate              time.Time
	InitialCapital       float64
	TradeFee             float64 // fee per trade, e.g. 0.001 for 0.1%
	Slippage             float64 // slippage percentage, e.g. 0.001 for 0.1%
}

// BollingerBands represents Bollinger Bands values
type BollingerBands struct {
	Upper  float64
	Middle float64 // Simple Moving Average
	Lower  float64
}

// Signal represents a trading signal
type Signal struct {
	Date   time.Time
	Type   string  // "BUY", "SELL", "HOLD"
	Price  float64
	Reason string
}
