package buycycle




type Receiver interface {
	OnReceiveBuycycleStockAggs(StockAggregate) error
}

func (w *Client) SubscribeStock(stock string) error {
	return w.WriteJSON(stock);
}