package ws

// ws package is an infrastructure that fetches stock data from external sources.

// There are several implements:
// 1. Buycycle
//   a. kor stocks
// 2. Polygon
//   a. usa stocks
//   b. crypto
// 3. KIS
//	 a. kor stocks
//   b. usa stocks
// 4. Mock

// Fetcher is an infrastructure interface for receiving data.
// Every stock fetcher must implement this interface.
type Fetcher interface {
	SubscribeStockAggs(stock string) error
	UnsubscribeStockAggs(stock string) error
	Close() error
	Ping() error
}

// Receiver is an interface for adapter that .
// An adapter that implements Receiver is given as an argument to Fetcher constructor.
type Receiver interface {
	OnReceiveStockAggs(*StockAggregate) error
}