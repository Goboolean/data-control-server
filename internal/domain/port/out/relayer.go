package out


type RelayerPort interface {
	FetchDomesticStock(stock string) error
	FetchForeignStock(stock string) error
}