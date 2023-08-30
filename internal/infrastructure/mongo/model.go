package mongo


type StockAggregate struct {
	StockID string  `bson:"stockId"`
	Open    float64 `bson:"start"`
	Closed  float64 `bson:"end"`
	Min     float64 `bson:"min"`
	Max     float64 `bson:"max"`
	Time    int64   `bson:"time"`
}