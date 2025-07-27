package visualization

import (
	"fmt"
	"os"
	stockTypes "swing-trader/internal/types"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// TradeMarker represents a trade entry or exit point for visualization
type TradeMarker struct {
	Date  string
	Price float64
	Type  string // "BUY" or "SELL"
	ID    string
}

// GenerateKLineChartWithTrades creates a candlestick chart with trade markers
func GenerateKLineChartWithTrades(stockData []stockTypes.StockData, trades []stockTypes.Trade, title, filePath string) error {
	// Prepare data for candlestick chart
	dates := make([]string, len(stockData))
	klineData := make([]opts.KlineData, len(stockData))

	for i, data := range stockData {
		dates[i] = data.Date.Format("2006-01-02")
		klineData[i] = opts.KlineData{
			Value: [4]float64{data.Open, data.Close, data.Low, data.High},
		}
	}

	// Create candlestick chart
	kline := charts.NewKLine()
	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: fmt.Sprintf("%s - Stock Price with Trades", title),
		}),
	)

	kline.SetXAxis(dates).AddSeries("Stock Price", klineData)

	// Save the chart
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer f.Close()

	return kline.Render(f)
}

// GenerateAccountBalanceChart creates a line chart showing account balance over time
func GenerateAccountBalanceChart(stockData []stockTypes.StockData, trades []stockTypes.Trade, initialCapital float64, title, filePath string) error {
	// Calculate account balance over time
	dates, balances := calculateAccountBalance(stockData, trades, initialCapital)

	// Create line chart
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: fmt.Sprintf("%s - Account Balance Over Time", title),
		}),
	)

	lineItems := make([]opts.LineData, len(balances))
	for i, balance := range balances {
		lineItems[i] = opts.LineData{Value: balance}
	}

	line.SetXAxis(dates).AddSeries("Account Balance", lineItems)

	// Save the chart
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer f.Close()

	return line.Render(f)
}

// generateTradeMarkers creates scatter plot data for trade entry and exit points
func generateTradeMarkers(stockData []stockTypes.StockData, trades []stockTypes.Trade) ([]opts.ScatterData, []opts.ScatterData) {
	// Create a map for quick date lookup
	dateToIndex := make(map[string]int)
	for i, data := range stockData {
		dateToIndex[data.Date.Format("2006-01-02")] = i
	}

	var buyMarkers []opts.ScatterData
	var sellMarkers []opts.ScatterData

	for _, trade := range trades {
		// Add buy marker
		buyDate := trade.EntryDate.Format("2006-01-02")
		if idx, exists := dateToIndex[buyDate]; exists {
			buyMarkers = append(buyMarkers, opts.ScatterData{
				Value:  []interface{}{idx, trade.EntryPrice},
				Symbol: "triangle",
				SymbolSize: 15,
			})
		}

		// Add sell marker if trade is closed
		if trade.ExitDate != nil && trade.ExitPrice != nil {
			sellDate := trade.ExitDate.Format("2006-01-02")
			if idx, exists := dateToIndex[sellDate]; exists {
				sellMarkers = append(sellMarkers, opts.ScatterData{
					Value:  []interface{}{idx, *trade.ExitPrice},
					Symbol: "triangle",
					SymbolSize: 15,
				})
			}
		}
	}

	return buyMarkers, sellMarkers
}

// calculateAccountBalance computes the account balance over time
func calculateAccountBalance(stockData []stockTypes.StockData, trades []stockTypes.Trade, initialCapital float64) ([]string, []float64) {
	dates := make([]string, len(stockData))
	balances := make([]float64, len(stockData))

	// Create trade events map
	tradeEvents := make(map[string]float64) // date -> P&L
	for _, trade := range trades {
		if trade.ExitDate != nil {
			exitDate := trade.ExitDate.Format("2006-01-02")
			tradeEvents[exitDate] += trade.ProfitLoss
		}
	}

	currentBalance := initialCapital
	for i, data := range stockData {
		dateStr := data.Date.Format("2006-01-02")
		dates[i] = dateStr

		// Check if there's a trade completion on this date
		if pnl, exists := tradeEvents[dateStr]; exists {
			currentBalance += pnl
		}

		balances[i] = currentBalance
	}

	return dates, balances
}

