package vo



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
	StockID string
	Platform  string
	Symbol    string
}


type StockAggregateForm struct {
	StockAggregate
	StockAggsMeta
}