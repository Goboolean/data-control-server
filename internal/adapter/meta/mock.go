package meta

import (
	"github.com/Goboolean/fetch-server.v1/internal/domain/port"
	"github.com/Goboolean/fetch-server.v1/internal/domain/port/out"
	"github.com/Goboolean/fetch-server.v1/internal/domain/vo"
)

var data = map[string]vo.StockAggsMeta{
	"stock.facebook.usa": {
		StockID:  "stock.facebook.usa",
		Platform: "mock",
		Symbol:   "FB",
	},

	"stock.apple.usa": {
		StockID:  "stock.apple.usa",
		Platform: "mock",
		Symbol:   "AAPL",
	},

	"stock.amazon.usa": {
		StockID:  "stock.amazon.usa",
		Platform: "mock",
		Symbol:   "AMZN",
	},

	"stock.netflix.usa": {
		StockID:  "stock.netflix.usa",
		Platform: "mock",
		Symbol:   "NFLX",
	},

	"stock.google.usa": {
		StockID:  "stock.google.usa",
		Platform: "mock",
		Symbol:   "GOOG",
	},
}

type MockAdapter struct{}

func NewMockAdapter() out.StockMetadataPort {
	return &MockAdapter{}
}

func (a *MockAdapter) CheckStockExists(tx port.Transactioner, stockId string) (bool, error) {
	_, ok := data[stockId]
	return ok, nil
}

func (a *MockAdapter) GetStockMetadata(tx port.Transactioner, stockId string) (vo.StockAggsMeta, error) {
	return data[stockId], nil
}

func (a *MockAdapter) GetAllStockMetadata(tx port.Transactioner) ([]vo.StockAggsMeta, error) {
	var stocks []vo.StockAggsMeta
	for _, stock := range data {
		stocks = append(stocks, stock)
	}
	return stocks, nil
}
