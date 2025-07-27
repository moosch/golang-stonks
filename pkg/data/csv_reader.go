package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"swing-trader/internal/types"
	"time"
)

// LoadStockDataFromCSV reads historical stock data from a CSV file
func LoadStockDataFromCSV(filePath string) ([]types.StockData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV data: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	// Skip header row if present
	startIndex := 0
	if len(records) > 0 && records[0][0] == "Date" {
		startIndex = 1
	}

	var stockData []types.StockData
	for i := startIndex; i < len(records); i++ {
		record := records[i]
		
		// Skip empty lines or lines with insufficient data
		if len(record) == 0 || (len(record) == 1 && record[0] == "") {
			continue
		}
		
		if len(record) < 7 {
			return nil, fmt.Errorf("invalid CSV format at row %d: expected 7 columns, got %d", i+1, len(record))
		}

		// Parse date - trying common formats
		var date time.Time
		dateFormats := []string{
			"Jan 2 2006",
			"2006-01-02",
			"01/02/2006",
			"1/2/2006",
		}
		
		dateStr := record[0]
		for _, format := range dateFormats {
			if d, err := time.Parse(format, dateStr); err == nil {
				date = d
				break
			}
		}
		
		if date.IsZero() {
			return nil, fmt.Errorf("failed to parse date %s at row %d", dateStr, i+1)
		}

		open, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse open price at row %d: %w", i+1, err)
		}

		high, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse high price at row %d: %w", i+1, err)
		}

		low, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse low price at row %d: %w", i+1, err)
		}

		close, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse close price at row %d: %w", i+1, err)
		}

		adjClose, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse adjusted close price at row %d: %w", i+1, err)
		}

		var volume int64
		if record[6] == "-" || record[6] == "" {
			volume = 0 // Set volume to 0 for missing data
		} else {
			var err error
			volume, err = strconv.ParseInt(record[6], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse volume at row %d: %w", i+1, err)
			}
		}

		stockData = append(stockData, types.StockData{
			Date:          date,
			Open:          open,
			High:          high,
			Low:           low,
			Close:         close,
			AdjustedClose: adjClose,
			Volume:        volume,
		})
	}

	// Sort data chronologically (oldest first)
	sort.Slice(stockData, func(i, j int) bool {
		return stockData[i].Date.Before(stockData[j].Date)
	})

	return stockData, nil
}

// FilterDataByDateRange filters stock data by start and end dates
func FilterDataByDateRange(data []types.StockData, startDate, endDate time.Time) []types.StockData {
	var filteredData []types.StockData
	
	for _, record := range data {
		if (record.Date.Equal(startDate) || record.Date.After(startDate)) &&
		   (record.Date.Equal(endDate) || record.Date.Before(endDate)) {
			filteredData = append(filteredData, record)
		}
	}
	
	return filteredData
}
