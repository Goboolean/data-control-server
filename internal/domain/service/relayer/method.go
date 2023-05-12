package relayer

import "github.com/Goboolean/stock-fetch-server/internal/domain/value"



func (m *RelayerManager) FetchStock(stock string) error {
	if err := m.store.storeStock(stock); err != nil {
		return err
	}

	if err := m.subscriber.fetchStock(stock); err != nil {
		return err
	}

	return nil
}


func (m *RelayerManager) StopFetchingStock(stock string) error {
	if err := m.store.unstoreStock(stock); err != nil {
		return err
	}

	if err := m.subscriber.unfetchStock(stock); err != nil {
		return err
	}

	return nil
}


func (m *RelayerManager) PlaceStockFormBatch(stock []value.StockAggregateForm) {
	for idx := range stock {
		m.pipe.PlaceOnStartPoint(stock[idx])
	}
}


func (m *RelayerManager) Subscribe(stock string) (<-chan []value.StockAggregate, error) {
	return m.pipe.GetEndpointChannel(stock)
}