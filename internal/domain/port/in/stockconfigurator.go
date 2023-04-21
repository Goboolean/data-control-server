package inport


type StockConfiguratorPort interface {
	SetStockRelayableTrue(string) error
	SetStockRelayableFalse(string) error
	SetStockStoreableTrue(string) error
	SetStockStoreableFalse(string) error
	SetStockTransmitableTrue(string) error
	SetStockTransmitableFalse(string) error
}