package fetcher



func (a *StockFetchAdapter) FetchDomesticStock(stock string) error {
	return a.b.SubscribeStockAggs(stock)
}

func (a *StockFetchAdapter) FetchForeignStock(stock string) error {
	return a.p.SubscribeStockAggs(stock)
}