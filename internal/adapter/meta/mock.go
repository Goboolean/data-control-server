package meta

import (
	"github.com/Goboolean/fetch-server/internal/domain/entity"
	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
)



var data = map[string] entity.StockAggsMeta{
	"stock.facebook.usa": {
		StockID: "stock.facebook.usa",
		Platform: "mock",
		Symbol: "FB",
	},

	"stock.apple.usa": {
		StockID: "stock.apple.usa",
		Platform: "mock",
		Symbol: "AAPL",
	},

	"stock.amazon.usa": {
		StockID: "stock.amazon.usa",
		Platform: "mock",
		Symbol: "AMZN",
	},

	"stock.netflix.usa": {
		StockID: "stock.netflix.usa",
		Platform: "mock",
		Symbol: "NFLX",
	},

	"stock.google.usa": {
		StockID: "stock.google.usa",
		Platform: "mock",
		Symbol: "GOOG",
	},

}


type MockAdapter struct {}

func NewMockAdapter() out.StockMetadataPort {
	return &MockAdapter{}
}


func (a *MockAdapter) CheckStockExists(tx port.Transactioner, stockId string) (bool, error) {
	_, ok := data[stockId]
	return ok, nil
}

func (a *MockAdapter) GetStockMetadata(tx port.Transactioner, stockId string) (entity.StockAggsMeta, error) {
	return data[stockId], nil
}

func (a *MockAdapter) GetAllStockMetadata(tx port.Transactioner) ([]entity.StockAggsMeta, error) {
	var stocks []entity.StockAggsMeta
	for _, stock := range data {
		stocks = append(stocks, stock)
	}
	return stocks, nil
}