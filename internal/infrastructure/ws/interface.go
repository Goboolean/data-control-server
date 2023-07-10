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
	// Subscribing several topic at once is allowed, but atomicity is not guaranteed.
	SubscribeStockAggs(...string) error
	// Unscribing several topic at once is allowed, but atomicity is not guaranteed.
	UnsubscribeStockAggs(...string) error
	Close() error
	Ping() error
	// It is need to distinguish between different fetchers,
	// and to subscribe/unsubscribe stock on appropraite platform.
	PlatformName() string
}

// Receiver is an interface for adapter that .
// An adapter that implements Receiver is given as an argument to Fetcher constructor.
type Receiver interface {
	OnReceiveStockAggs(*StockAggregate) error
}


// A stock aggs structure that every implementation shares.
type StockAggregate struct {
	EventType string	
	Average   float64
	Min       float64
	Max       float64
	Start     float64
	End       float64

	StartTime int64
	EndTime   int64
}
