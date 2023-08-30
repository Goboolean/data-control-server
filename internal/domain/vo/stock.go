package vo



type StockAggregate struct {
	EventType string
	Time      int64
	Open      float64
	Closed    float64
	Max       float64
	Min       float64
	Volume	  float64
}

type StockAggsMeta struct {
	StockID  string
	Platform string
	Symbol   string
}


type StockAggregateForm struct {
	StockAggregate
	StockAggsMeta
}