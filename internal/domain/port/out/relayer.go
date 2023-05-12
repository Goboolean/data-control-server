package out


type RelayerPort interface {
	FetchDomesticStock(stock string) error
	FetchInternationalStock(stock string) error
}