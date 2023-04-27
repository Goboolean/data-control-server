package config


type OptionType int

const (
	StockRelay OptionType = iota
	StockReal
	StockPersistance
)