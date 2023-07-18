package entity



type StockAggregate struct {
	EventType string	
	Average float64
	Min    float64
	Max    float64
	Start  float64
	End    float64
	StartTime int64
	EndTime   int64
}

type StockAggsMeta struct {
	Platform  string
	Symbol    string
}


type StockAggregateForm struct {
	StockID string
	StockAggregate
	StockAggsMeta
}