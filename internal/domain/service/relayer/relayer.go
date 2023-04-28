package relayer

import "errors"



func (m *RelayerManager) SubscribeWebsocket(stock string) error {
	if m.LinedUpRelayer.isExist(stock) == true {
		return errors.New("stock is already subscribed")
	}

	m.LinedUpRelayer.openChannel(stock)
	return nil
}

func (m *RelayerManager) UnsubscribeWebsocket(stock string) error {
	if m.LinedUpRelayer.isExist(stock) == false {
		return errors.New("stock is already unsubscribed")
	}

	m.LinedUpRelayer.closeChannel(stock)
	return nil
}

func (m *RelayerManager) TransferRawToLinedUp() {

	go func() {
		data := <- m.RawRelayer.queue
		stock := data.StockID
	
		m.MiddleRelayer.Push(stock, &data.StockAggregate)
	
		batch, ok := m.MiddleRelayer.EmptyQueue(data.StockID)
	
		if ok {
			m.LinedUpRelayer.Push(stock, batch)
		}
	}()

}