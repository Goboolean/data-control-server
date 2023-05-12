package buycycle


type Receiver interface {
	OnReceiveBuycycleStockAggs(StockAggregate) error
}

type Fetcher interface {
	Close() error
	SubscribeStockAggs(stock string) error
}