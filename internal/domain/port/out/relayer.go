package out



type RelayerPort interface {
	FetchDomesticStock(string) error
	FetchInternationalStock(string) error
}