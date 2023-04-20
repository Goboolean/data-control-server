package buycycle

type StockAggregate struct {
	StockID   string
	EventType string
	Average   float64
	Min       float64
	Max       float64
	Start     float64
	End       float64

	StartTime int64
	EndTime   int64
}
