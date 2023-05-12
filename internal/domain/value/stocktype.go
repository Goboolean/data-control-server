package value



type StockType int

const (
	Domestic StockType = iota
	International
	Invalid
)