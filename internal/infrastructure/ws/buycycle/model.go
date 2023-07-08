package buycycle

type HeaderJson struct {
	Uuid string `json:""`
	Code string `json:""`
}

type BodyJson struct {
	
}

type StockAggregate struct {
	StockID   string  `json:""`
	EventType string  `json:""`
	Average   float64 `json:""`
	Min       float64 `json:""`
	Max       float64 `json:""`
	Start     float64 `json:""`
	End       float64 `json:""`

	StartTime int64  `json:""`
	EndTime   int64  `json:""`
}
