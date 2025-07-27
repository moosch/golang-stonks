package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"swing-trader/internal/types"
	"swing-trader/pkg/backtesting"
	"swing-trader/pkg/data"
	"time"
)

func main() {
	// Define command line flags
	var (
		dataPath       = flag.String("data", "", "Path to CSV file with historical stock data")
		startDate      = flag.String("start", "", "Start date for backtest (YYYY-MM-DD)")
		endDate        = flag.String("end", "", "End date for backtest (YYYY-MM-DD)")
		initialCapital = flag.Float64("capital", 10000.0, "Initial capital for backtesting")
		buyThreshold   = flag.Float64("buy-rsi", 30.0, "RSI threshold for buying (oversold)")
		sellThreshold  = flag.Float64("sell-rsi", 70.0, "RSI threshold for selling (overbought)")
		stopLoss       = flag.Float64("stop-loss", 0.05, "Stop loss percentage (e.g., 0.05 for 5%)")
		takeProfit     = flag.Float64("take-profit", 0.10, "Take profit percentage (e.g., 0.10 for 10%)")
		positionSize   = flag.Float64("position-size", 0.02, "Position size as percentage of capital (e.g., 0.02 for 2%)")
		maxDrawdown    = flag.Float64("max-drawdown", 0.20, "Maximum drawdown percentage (e.g., 0.20 for 20%)")
		tradeFee       = flag.Float64("trade-fee", 0.001, "Trade fee percentage (e.g., 0.001 for 0.1%)")
		slippage       = flag.Float64("slippage", 0.001, "Slippage percentage (e.g., 0.001 for 0.1%)")
		rsiPeriod      = flag.Int("rsi-period", 14, "RSI calculation period")
		bbPeriod       = flag.Int("bb-period", 20, "Bollinger Bands calculation period")
		bbStdDev       = flag.Float64("bb-stddev", 2.0, "Bollinger Bands standard deviation multiplier")
	)
	flag.Parse()

	// Validate required flags
	if *dataPath == "" {
		log.Fatal("Data path is required. Use -data flag to specify CSV file path.")
	}

	// Parse dates
	var start, end time.Time
	var err error
	
	if *startDate != "" {
		start, err = time.Parse("2006-01-02", *startDate)
		if err != nil {
			log.Fatalf("Invalid start date format: %v", err)
		}
	}
	
	if *endDate != "" {
		end, err = time.Parse("2006-01-02", *endDate)
		if err != nil {
			log.Fatalf("Invalid end date format: %v", err)
		}
	}

	// Load stock data
	fmt.Printf("Loading stock data from %s...\n", *dataPath)
	stockData, err := data.LoadStockDataFromCSV(*dataPath)
	if err != nil {
		log.Fatalf("Failed to load stock data: %v", err)
	}

	fmt.Printf("Loaded %d data points\n", len(stockData))

	// Filter data by date range if specified
	if !start.IsZero() || !end.IsZero() {
		if start.IsZero() {
			start = stockData[0].Date
		}
		if end.IsZero() {
			end = stockData[len(stockData)-1].Date
		}
		stockData = data.FilterDataByDateRange(stockData, start, end)
		fmt.Printf("Filtered to %d data points between %s and %s\n", 
			len(stockData), start.Format("2006-01-02"), end.Format("2006-01-02"))
	}

	if len(stockData) == 0 {
		log.Fatal("No data available for the specified date range")
	}

	// Create backtest configuration
	config := types.BacktestConfig{
		StockDataPath:  *dataPath,
		InitialCapital: *initialCapital,
		TradeFee:       *tradeFee,
		Slippage:       *slippage,
		StartDate:      stockData[0].Date,
		EndDate:        stockData[len(stockData)-1].Date,
		StrategyConfig: types.StrategyConfig{
			BuyThreshold:   *buyThreshold,
			SellThreshold:  *sellThreshold,
			StopLoss:       *stopLoss,
			TakeProfit:     *takeProfit,
			InitialCapital: *initialCapital,
			RSIPeriod:      *rsiPeriod,
			BBPeriod:       *bbPeriod,
			BBStdDev:       *bbStdDev,
		},
		RiskManagementConfig: types.RiskManagementConfig{
			MaxDrawdown:  *maxDrawdown,
			PositionSize: *positionSize,
		},
	}

	// Run backtest
	fmt.Println("Running backtest...")
	engine := backtesting.NewEngine(config)
	result, err := engine.Run(stockData)
	if err != nil {
		log.Fatalf("Backtest failed: %v", err)
	}

	// Display results
	printResults(result)
}

// printResults displays the backtest results in a formatted way
func printResults(result *types.BacktestResult) {
	separator := strings.Repeat("=", 60)
	fmt.Println("\n" + separator)
	fmt.Println("BACKTEST RESULTS")
	fmt.Println(separator)
	
	fmt.Printf("Period: %s to %s\n", 
		result.StartDate.Format("2006-01-02"), 
		result.EndDate.Format("2006-01-02"))
	
	fmt.Println("\nCapital:")
	fmt.Printf("  Initial Capital:    $%.2f\n", result.InitialCapital)
	fmt.Printf("  Final Capital:      $%.2f\n", result.FinalCapital)
	fmt.Printf("  Total P&L:          $%.2f\n", result.TotalProfitLoss)
	fmt.Printf("  Total Return:       %.2f%%\n", result.TotalReturn)
	fmt.Printf("  Annualized Return:  %.2f%%\n", result.AnnualizedReturn)
	
	fmt.Println("\nTrade Statistics:")
	fmt.Printf("  Total Trades:       %d\n", result.TotalTrades)
	fmt.Printf("  Winning Trades:     %d\n", result.WinningTrades)
	fmt.Printf("  Losing Trades:      %d\n", result.LosingTrades)
	fmt.Printf("  Win Rate:           %.1f%%\n", result.WinRate)
	
	if result.AverageWin > 0 {
		fmt.Printf("  Average Win:        $%.2f\n", result.AverageWin)
	}
	if result.AverageLoss > 0 {
		fmt.Printf("  Average Loss:       $%.2f\n", result.AverageLoss)
	}
	
	fmt.Println("\nRisk Metrics:")
	fmt.Printf("  Max Drawdown:       %.2f%%\n", result.MaxDrawdown)
	
	if len(result.Trades) > 0 {
		fmt.Println("\nRecent Trades:")
		count := 5
		if len(result.Trades) < count {
			count = len(result.Trades)
		}
		
		for i := len(result.Trades) - count; i < len(result.Trades); i++ {
			trade := result.Trades[i]
			var exitDate string
			if trade.ExitDate != nil {
				exitDate = trade.ExitDate.Format("2006-01-02")
			} else {
				exitDate = "Open"
			}
			
			fmt.Printf("  %s: Entry %s @$%.2f -> Exit %s @$%.2f | P&L: $%.2f\n",
				trade.ID,
				trade.EntryDate.Format("2006-01-02"),
				trade.EntryPrice,
				exitDate,
				func() float64 {
					if trade.ExitPrice != nil {
						return *trade.ExitPrice
					}
					return 0
				}(),
				trade.ProfitLoss)
		}
	}
	
	fmt.Println(separator)
}
