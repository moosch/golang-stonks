# Swing Trading Backtesting System

A comprehensive backtesting application built in Go for swing trading strategies using Bollinger Bands and RSI indicators.

## Features

- **Bollinger Bands + RSI Strategy**: Combines Bollinger Bands and RSI indicators for swing trading signals
- **Comprehensive Backtesting**: Simulates trades with realistic fees, slippage, and risk management
- **Risk Management**: Includes stop-loss, take-profit, position sizing, and drawdown controls
- **Performance Metrics**: Detailed analysis including win rate, average win/loss, max drawdown, and annualized returns
- **Flexible Configuration**: Extensive command-line options to customize strategy parameters
- **CSV Data Support**: Reads historical stock data from CSV files

## Project Structure

```
swing-trader/
├── cmd/main.go                     # Main application entry point
├── internal/types/types.go         # Core data structures
├── pkg/
│   ├── indicators/                 # Technical indicators
│   │   ├── bollinger_bands.go     # Bollinger Bands calculation
│   │   ├── bollinger_bands_test.go # Bollinger Bands tests
│   │   ├── rsi.go                 # RSI calculation
│   │   └── rsi_test.go            # RSI tests
│   ├── data/                      # Data handling
│   │   └── csv_reader.go          # CSV file reader
│   ├── strategy/                  # Trading strategies
│   │   └── bb_rsi_strategy.go     # Bollinger Bands + RSI strategy
│   └── backtesting/               # Backtesting engine
│       └── engine.go              # Main backtesting logic
├── historic_data/                 # Historical stock data files
└── README.md                      # This file
```

## Strategy Overview

The backtesting system implements a **Bollinger Bands + RSI Strategy**:

### Buy Signals
- Stock price is **below the lower Bollinger Band** (indicating potential oversold condition)
- **AND** RSI is **below the buy threshold** (default: 30, confirming oversold condition)

### Sell Signals
- RSI is **above the sell threshold** (default: 70, indicating overbought condition)

### Risk Management
- **Stop Loss**: Automatically closes positions when losses reach a specified percentage (default: 5%)
- **Take Profit**: Automatically closes positions when profits reach a specified percentage (default: 10%)
- **Position Sizing**: Calculates position size based on available capital and risk tolerance
- **Single Position**: Only one position open at a time for simplicity

## Installation & Usage

### Build the Application

```bash
go build -o backtest cmd/main.go
```

### Basic Usage

```bash
./backtest -data historic_data/historic_AAPL.csv -capital 10000
```

### Advanced Usage with Custom Parameters

```bash
./backtest \
  -data historic_data/historic_AAPL.csv \
  -start 2020-01-01 \
  -end 2024-12-31 \
  -capital 10000 \
  -buy-rsi 25 \
  -sell-rsi 75 \
  -stop-loss 0.08 \
  -take-profit 0.15 \
  -position-size 0.03
```

## Command Line Options

### Required Parameters
- `-data`: Path to CSV file with historical stock data

### Basic Configuration
- `-capital`: Initial capital for backtesting (default: $10,000)
- `-start`: Start date for backtest (YYYY-MM-DD format)
- `-end`: End date for backtest (YYYY-MM-DD format)

### Strategy Parameters
- `-buy-rsi`: RSI threshold for buying (default: 30.0)
- `-sell-rsi`: RSI threshold for selling (default: 70.0)
- `-rsi-period`: RSI calculation period (default: 14)
- `-bb-period`: Bollinger Bands calculation period (default: 20)
- `-bb-stddev`: Bollinger Bands standard deviation multiplier (default: 2.0)

### Risk Management
- `-stop-loss`: Stop loss percentage (default: 0.05 = 5%)
- `-take-profit`: Take profit percentage (default: 0.10 = 10%)
- `-position-size`: Position size as percentage of capital (default: 0.02 = 2%)
- `-max-drawdown`: Maximum drawdown percentage (default: 0.20 = 20%)

### Trading Costs
- `-trade-fee`: Trade fee percentage (default: 0.001 = 0.1%)
- `-slippage`: Slippage percentage (default: 0.001 = 0.1%)

## CSV Data Format

The application expects CSV files with the following format:

```csv
Date,Open,High,Low,Close,AdjClose,Volume
Jul 2 2025,209.08,213.34,208.14,212.44,212.44,66327031
Jul 1 2025,206.67,210.19,206.14,207.82,207.82,78788900
...
```

### Supported Date Formats
- `Jan 2 2006` (e.g., "Jul 2 2025")
- `2006-01-02` (e.g., "2025-07-02")
- `01/02/2006` (e.g., "07/02/2025")
- `1/2/2006` (e.g., "7/2/2025")

## Example Results

```
============================================================
BACKTEST RESULTS
============================================================
Period: 2020-01-02 to 2024-12-31

Capital:
  Initial Capital:    $10000.00
  Final Capital:      $12251.46
  Total P&L:          $2251.46
  Total Return:       22.51%
  Annualized Return:  4.15%

Trade Statistics:
  Total Trades:       6
  Winning Trades:     5
  Losing Trades:      1
  Win Rate:           83.3%
  Average Win:        $514.31
  Average Loss:       $320.11

Risk Metrics:
  Max Drawdown:       2.66%

Recent Trades:
  T2: Entry 2021-02-25 @$121.11 -> Exit 2021-06-30 @$136.82 | P&L: $545.13
  T3: Entry 2022-09-30 @$138.34 -> Exit 2023-02-02 @$150.67 | P&L: $389.77
  ...
============================================================
```

## Testing

Run the unit tests for the indicators:

```bash
go test ./pkg/indicators/...
```

## Available Data Files

The `historic_data/` directory includes sample data for:
- AAPL (Apple Inc.)
- MSFT (Microsoft Corporation)
- PLTR (Palantir Technologies)
- BABA (Alibaba Group)
- AVGO (Broadcom Inc.)
- WMT (Walmart Inc.)

## Key Performance Metrics

- **Total Return**: Overall percentage return on investment
- **Annualized Return**: Compound annual growth rate (CAGR)
- **Win Rate**: Percentage of profitable trades
- **Average Win/Loss**: Average profit from winning trades vs. average loss from losing trades
- **Max Drawdown**: Maximum peak-to-trough decline in portfolio value
- **Sharpe Ratio**: Risk-adjusted return metric (calculated in BacktestResult but not yet displayed)

## Technical Indicators

### Bollinger Bands
- **Period**: 20 days (configurable)
- **Standard Deviation**: 2.0 (configurable)
- **Components**: Upper band, middle band (SMA), lower band

### RSI (Relative Strength Index)
- **Period**: 14 days (configurable)
- **Range**: 0-100
- **Overbought**: Above 70 (configurable)
- **Oversold**: Below 30 (configurable)

## Limitations

- **Single Asset**: Currently supports backtesting one stock at a time
- **Single Position**: Only one open position at a time
- **Daily Data**: Designed for daily timeframe data
- **No Dividends**: Does not account for dividend payments
- **No Stock Splits**: Uses adjusted close prices but doesn't handle splits explicitly

## Future Enhancements

Potential improvements that could be added:
- Portfolio-level backtesting with multiple assets
- Multiple concurrent positions
- Additional technical indicators (MACD, Stochastic, etc.)
- More sophisticated risk management (portfolio-level stops, correlation analysis)
- Web-based UI for easier configuration and visualization
- Database support for storing results
- More comprehensive performance analytics

## Contributing

1. Ensure all tests pass: `go test ./...`
2. Follow Go best practices and formatting
3. Add unit tests for new functionality
4. Update documentation as needed

## License

This project is open source and available under the MIT License.
