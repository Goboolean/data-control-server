package polygon

import "github.com/polygon-io/client-go/websocket/models"


type Receiver interface {
	OnReceivePolygonStockAggs(models.EquityAgg) error
}

type Fetcher interface {
	Close() error
	SubscribeStockAggs(stock string) error
}