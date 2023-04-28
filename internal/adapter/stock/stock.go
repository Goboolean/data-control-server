package stock




type StockAdapter struct {}

var instance = &StockAdapter{}

func NewStockAdapter() *StockAdapter {
	return instance
}

