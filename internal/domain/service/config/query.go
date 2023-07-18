package config

import "github.com/Goboolean/fetch-server/internal/domain/entity"

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

func (m *ConfigurationManager) GetStockConfiguration(stock string) (entity.StockConfiguration, error) {

}

func (m *ConfigurationManager) GetAllStockConfiguration() ([]entity.StockConfiguration, error) {
 // get all stock list
 // reflect all stock info to list
}