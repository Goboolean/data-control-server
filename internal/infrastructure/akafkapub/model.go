package akafkapub

type StockAggregate struct {
	StockID   string  `bson:"stockId"`
	EventType string  `bson:"eventType"`
	Avg       float64 `bson:"avg"`
	Min       float64 `bson:"min"`
	Max       float64 `bson:"max"`
	Start     float64 `bson:"start"`
	End       float64 `bson:"end"`

	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}
