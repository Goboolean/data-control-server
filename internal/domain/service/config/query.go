package config

func (m *ConfigurationManager) SetStockRelayableTrue(stock string) error {
	return m.relayer.SubscribeWebsocket(stock)
}

func (m *ConfigurationManager) SetStockRelayableFalse(stock string) error {
	return m.relayer.UnsubscribeWebsocket(stock)
}

func (m *ConfigurationManager) SetStockStoreableTrue(stock string) error {
	return m.persistence.SubscribeRelayer(stock)
}

func (m *ConfigurationManager) SetStockStoreableFalse(stock string) error {
	return m.persistence.UnsubscribeRelayer(stock)
}


func (m *ConfigurationManager) SetStockTransmittableTrue(string) error {
}

func (m *ConfigurationManager) SetStockTransmittableFalse(string) error {
}
