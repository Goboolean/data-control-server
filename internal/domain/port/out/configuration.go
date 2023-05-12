package out



type ConfiguratorPort interface {
	GetStockInfo() (int, error)
}