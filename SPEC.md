# Stock app

## Goal

Create a stock backtesting app in Golang that uses historical stock data (found in the csv files in `historic_data`) along with a strategy and risk management parameters to simulate trades and calculate profit/loss.

## Plan

Create separate packages the following parts:

- [ ] **Data**: Read the historical stock data from the CSV files and provide it in a format that can be used by the strategy.
- [ ] **Strategy**: Implement the strategy logic that will use the data to make buy/sell decisions.
- [ ] **Risk Management**: Implement risk management logic to manage the trades and calculate profit/loss.
- [ ] **Backtesting**: Implement the backtesting logic that will run the strategy against the historical data and provide results.
- [ ] **UI**: Create a user interface to display the results of the backtesting and allow users to configure the strategy and risk management parameters.
- [ ] **Entry**: Create an entry point, inside `cmd/main.go`, that will run the backtesting with the provided parameters as cli flags and return the results.

## Data Structures

### StockData

```go
type StockData struct {
	Date        time.Time
	Open        float64
	High        float64
	Low         float64
	Close       float64
	Volume      int64
	AdjustedClose float64
}
```

### Trade

```go
type Trade struct {
	ID          string
	EntryDate   time.Time
	ExitDate    time.Time
	EntryPrice  float64
	ExitPrice   float64
	Quantity    int64
	ProfitLoss   float64
	Status      string // "open", "closed", "cancelled"
}

type TradeResult struct {
	TotalTrades int64
	TotalProfitLoss float64
	WinningTrades int64
	LosingTrades int64
	WinRate float64 // percentage of winning trades
}
```

### StrategyConfig

```go
type StrategyConfig struct {
	BuyThreshold  float64 // e.g. RSI threshold for buying
	SellThreshold float64 // e.g. RSI threshold for selling
	StopLoss      float64 // e.g. percentage for stop loss
	TakeProfit    float64 // e.g. percentage for take profit
	InitialCapital float64 // e.g. starting capital for the backtest
}

type RiskManagementConfig struct {
	MaxDrawdown float64 // e.g. maximum drawdown percentage
	PositionSize float64 // e.g. percentage of capital to risk per trade
}
```
### BacktestResult

```go
type BacktestResult struct {
	Trades       []Trade
	TotalProfitLoss float64
	WinRate      float64 // percentage of winning trades
	TotalTrades  int64
	WinningTrades int64
	LosingTrades int64
	AverageWin   float64 // average profit of winning trades
	AverageLoss  float64 // average loss of losing trades
	TotalTradesDuration time.Duration // total duration of all trades
	TotalTradesCount int64 // total number of trades
	TotalTradesProfitLoss float64 // total profit/loss of all trades
	TotalTradesWinRate float64 // win rate of all trades
	TotalTradesAverageWin float64 // average win of all trades
	TotalTradesAverageLoss float64 // average loss of all trades
}
```

### BacktestConfig

```go
type BacktestConfig struct {
	StockDataPath string // path to the CSV files with historical stock data
	StrategyConfig StrategyConfig // configuration for the strategy
	RiskManagementConfig RiskManagementConfig // configuration for risk management
	StartDate time.Time // start date for the backtest
	EndDate time.Time // end date for the backtest
	InitialCapital float64 // initial capital for the backtest
	TradeFee float64 // fee per trade, e.g. 0.001 for 0.1%
	Slippage float64 // slippage percentage, e.g. 0.001 for 0.1%
}
```

## Backtesting Data

Here's an example of the historical stock data in CSV format:

```csv
Date,Open,High,Low,Close,AdjClose,Volume
Jul 2 2025,209.08,213.34,208.14,212.44,212.44,66327031
Jul 1 2025,206.67,210.19,206.14,207.82,207.82,78788900
Jun 30 2025,202.01,207.39,199.26,205.17,205.17,91912800
```

## Strategy Implementation

The strategy should be tuneable via the `StrategyConfig` struct, allowing users to set parameters like buy/sell thresholds, stop loss, and take profit levels.

There are 2 indicators that we will use to begine with. Bollinger Bands and RSI.
If a stock's price is below the lower Bollinger Band and RSI is below a certain threshold, we will buy.

Once a position has been opened, there needs to be a stop loss and take profit level set. That way, during the backtest, we can simulate the trade being closed at those levels.

Split the calculations for Bollinger Bands and RSI into separate packages inside a `pkg` directory, so they can be reused. Include unit tests for these calculations.

## Risk Management Implementation

The risk management config should also be set via the `RiskManagementConfig` struct, allowing users to set parameters like maximum drawdown and position size.

## Error Handling

Using the standard Go error handling, return errors from functions where appropriate. Use custom error types if needed to provide more context.

If there is an error in the backtesting process, return a descriptive error message that includes the context of the error.
