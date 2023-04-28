package inport

type ConfiguratorPort interface {
	SetStockRelayableTrue(string) error
	SetStockRelayableFalse(string) error
	SetStockStoreableTrue(string) error
	SetStockStoreableFalse(string) error
	SetStockTransmittableTrue(string) error
	SetStockTransmittableFalse(string) error
}
