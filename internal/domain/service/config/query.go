package config

func (m *ConfigurationManager) SetStockRelayableTrue(stock string) error {
	return m.relayer.FetchStock(stock)
}

func (m *ConfigurationManager) SetStockRelayableFalse(stock string) error {
	return m.relayer.StopFetchingStock(stock)
}

func (m *ConfigurationManager) SetStockStoreableTrue(stock string) error {
	return m.persistence.SubscribeRelayer(stock)
}

func (m *ConfigurationManager) SetStockStoreableFalse(stock string) error {
	return m.persistence.UnsubscribeRelayer(stock)
}

func (m *ConfigurationManager) SetStockTransmittableTrue(stock string) error {
	return m.transmitter.SubscribeRelayer(stock)
}

func (m *ConfigurationManager) SetStockTransmittableFalse(stock string) error {
	return m.transmitter.SubscribeRelayer(stock)
}
