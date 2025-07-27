package backtesting

import (
	"fmt"
	"math"
	"swing-trader/internal/types"
	"swing-trader/pkg/strategy"
	"time"
)

// Engine handles the backtesting execution
type Engine struct {
	config   types.BacktestConfig
	strategy *strategy.BBRSIStrategy
}

// NewEngine creates a new backtesting engine
func NewEngine(config types.BacktestConfig) *Engine {
	return &Engine{
		config:   config,
		strategy: strategy.NewBBRSIStrategy(config.StrategyConfig),
	}
}

// Run executes the backtest and returns results
func (e *Engine) Run(data []types.StockData) (*types.BacktestResult, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("no data provided for backtesting")
	}

	// Generate trading signals
	signals := e.strategy.GenerateSignals(data)
	
	// Execute trades based on signals
	trades, err := e.executeTrades(signals, data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute trades: %w", err)
	}

	// Calculate comprehensive results
	result := e.calculateResults(trades, data)
	
	return result, nil
}

// executeTrades processes signals and simulates trade execution
func (e *Engine) executeTrades(signals []types.Signal, data []types.StockData) ([]types.Trade, error) {
	var trades []types.Trade
	var openTrades []types.Trade
	availableCapital := e.config.InitialCapital
	tradeID := 1

	// Create a map for quick data lookup by date
	dataMap := make(map[time.Time]types.StockData)
	for _, d := range data {
		dataMap[d.Date] = d
	}

	for _, signal := range signals {
		switch signal.Type {
		case "BUY":
			if len(openTrades) == 0 { // Only open one position at a time for simplicity
				shares := e.strategy.CalculatePositionSize(availableCapital, signal.Price, e.config.RiskManagementConfig)
				if shares > 0 {
					// Apply slippage and fees
					entryPrice := signal.Price * (1 + e.config.Slippage)
					tradeFee := float64(shares) * entryPrice * e.config.TradeFee
					totalCost := float64(shares)*entryPrice + tradeFee

					if totalCost <= availableCapital {
						trade := types.Trade{
							ID:         fmt.Sprintf("T%d", tradeID),
							EntryDate:  signal.Date,
							EntryPrice: entryPrice,
							Quantity:   shares,
							Status:     "open",
							StopLoss:   e.strategy.GetStopLossPrice(entryPrice),
							TakeProfit: e.strategy.GetTakeProfitPrice(entryPrice),
						}
						openTrades = append(openTrades, trade)
						availableCapital -= totalCost
						tradeID++
					}
				}
			}

		case "SELL":
			// Close all open positions on sell signal
			for i := range openTrades {
				exitPrice := signal.Price * (1 - e.config.Slippage)
				tradeFee := float64(openTrades[i].Quantity) * exitPrice * e.config.TradeFee
				proceeds := float64(openTrades[i].Quantity)*exitPrice - tradeFee
				
				openTrades[i].ExitDate = &signal.Date
				openTrades[i].ExitPrice = &exitPrice
				openTrades[i].Status = "closed"
				openTrades[i].ProfitLoss = proceeds - (float64(openTrades[i].Quantity) * openTrades[i].EntryPrice)
				
				availableCapital += proceeds
				trades = append(trades, openTrades[i])
			}
			openTrades = nil
		}

		// Check stop loss and take profit for open trades
		openTrades = e.checkStopLossAndTakeProfit(openTrades, signal, &trades, &availableCapital)
	}

	// Close any remaining open trades at the end
	if len(openTrades) > 0 && len(data) > 0 {
		lastPrice := data[len(data)-1].Close
		lastDate := data[len(data)-1].Date
		
		for i := range openTrades {
			exitPrice := lastPrice * (1 - e.config.Slippage)
			tradeFee := float64(openTrades[i].Quantity) * exitPrice * e.config.TradeFee
			proceeds := float64(openTrades[i].Quantity)*exitPrice - tradeFee
			
			openTrades[i].ExitDate = &lastDate
			openTrades[i].ExitPrice = &exitPrice
			openTrades[i].Status = "closed"
			openTrades[i].ProfitLoss = proceeds - (float64(openTrades[i].Quantity) * openTrades[i].EntryPrice)
			
			trades = append(trades, openTrades[i])
		}
	}

	return trades, nil
}

// checkStopLossAndTakeProfit checks if any open trades should be closed due to stop loss or take profit
func (e *Engine) checkStopLossAndTakeProfit(openTrades []types.Trade, signal types.Signal, trades *[]types.Trade, availableCapital *float64) []types.Trade {
	var remainingTrades []types.Trade

	for _, trade := range openTrades {
		closed := false
		
		// Check stop loss
		if signal.Price <= trade.StopLoss {
			exitPrice := signal.Price * (1 - e.config.Slippage)
			tradeFee := float64(trade.Quantity) * exitPrice * e.config.TradeFee
			proceeds := float64(trade.Quantity)*exitPrice - tradeFee
			
			trade.ExitDate = &signal.Date
			trade.ExitPrice = &exitPrice
			trade.Status = "closed"
			trade.ProfitLoss = proceeds - (float64(trade.Quantity) * trade.EntryPrice)
			
			*availableCapital += proceeds
			*trades = append(*trades, trade)
			closed = true
		} else if signal.Price >= trade.TakeProfit {
			// Check take profit
			exitPrice := signal.Price * (1 - e.config.Slippage)
			tradeFee := float64(trade.Quantity) * exitPrice * e.config.TradeFee
			proceeds := float64(trade.Quantity)*exitPrice - tradeFee
			
			trade.ExitDate = &signal.Date
			trade.ExitPrice = &exitPrice
			trade.Status = "closed"
			trade.ProfitLoss = proceeds - (float64(trade.Quantity) * trade.EntryPrice)
			
			*availableCapital += proceeds
			*trades = append(*trades, trade)
			closed = true
		}

		if !closed {
			remainingTrades = append(remainingTrades, trade)
		}
	}

	return remainingTrades
}

// calculateResults computes comprehensive backtest results
func (e *Engine) calculateResults(trades []types.Trade, data []types.StockData) *types.BacktestResult {
	result := &types.BacktestResult{
		Trades:         trades,
		InitialCapital: e.config.InitialCapital,
		StartDate:      data[0].Date,
		EndDate:        data[len(data)-1].Date,
	}

	// Calculate basic metrics
	var totalPL float64
	var winningTrades, losingTrades int64
	var totalWinAmount, totalLossAmount float64

	for _, trade := range trades {
		totalPL += trade.ProfitLoss
		if trade.ProfitLoss > 0 {
			winningTrades++
			totalWinAmount += trade.ProfitLoss
		} else if trade.ProfitLoss < 0 {
			losingTrades++
			totalLossAmount += math.Abs(trade.ProfitLoss)
		}
	}

	result.TotalTrades = int64(len(trades))
	result.WinningTrades = winningTrades
	result.LosingTrades = losingTrades
	result.TotalProfitLoss = totalPL
	result.FinalCapital = e.config.InitialCapital + totalPL

	if result.TotalTrades > 0 {
		result.WinRate = float64(winningTrades) / float64(result.TotalTrades) * 100
	}

	if winningTrades > 0 {
		result.AverageWin = totalWinAmount / float64(winningTrades)
	}

	if losingTrades > 0 {
		result.AverageLoss = totalLossAmount / float64(losingTrades)
	}

	// Calculate total return
	result.TotalReturn = (result.FinalCapital - result.InitialCapital) / result.InitialCapital * 100

	// Calculate annualized return
	years := result.EndDate.Sub(result.StartDate).Hours() / (24 * 365.25)
	if years > 0 && result.FinalCapital > 0 && result.InitialCapital > 0 {
		result.AnnualizedReturn = (math.Pow(result.FinalCapital/result.InitialCapital, 1/years) - 1) * 100
	}

	// Calculate max drawdown (simplified)
	result.MaxDrawdown = e.calculateMaxDrawdown(trades)

	return result
}

// calculateMaxDrawdown calculates the maximum drawdown during the backtest period
func (e *Engine) calculateMaxDrawdown(trades []types.Trade) float64 {
	if len(trades) == 0 {
		return 0
	}

	peak := e.config.InitialCapital
	maxDrawdown := 0.0
	runningCapital := e.config.InitialCapital

	for _, trade := range trades {
		runningCapital += trade.ProfitLoss
		
		if runningCapital > peak {
			peak = runningCapital
		}
		
		drawdown := (peak - runningCapital) / peak * 100
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}
	}

	return maxDrawdown
}
