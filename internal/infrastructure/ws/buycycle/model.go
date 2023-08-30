package buycycle

import (
	"strconv"

	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/ws"
)

type HeaderJson struct {
	Uuid string `json:"uuid"`
	Code string `json:"szTrCode"`
}

type RequestBodyJson struct {
	TrName string `json:"trName"`
	BNext  bool   `json:"bNext"`
	Query  struct {
		Shcode string `json:"shcode"`
	} `json:"query"`
}

type RequestJson struct {
	Header HeaderJson      `json:"header"`
	Body   RequestBodyJson `json:"body"`
}

type ResponseBodyJson struct {
	OutBlock []StockDetail `json:"outblock"`
}

type StockDetail struct {
	Mdchecnt   string `json:"mdchecnt"`
	Sign       string `json:"sign"`
	Mschecnt   string `json:"mschecnt"`
	Mdvolume   string `json:"mdvolume"`
	W_avrg     string `json:"w_avrg"`
	Cpower     string `json:"cpower"`
	Offerho    string `json:"offerho"`
	Cvolume    string `json:"cvolume"`
	High       string `json:"high"`
	Bidho      string `json:"bidho"`
	Low        string `json:"low"`
	Price      string `json:"price"`
	Cgubun     string `json:"cgubun"`
	Value      string `json:"value"`
	Change     string `json:"change"`
	Shcode     string `json:"shcode"`
	Chetime    string `json:"chetime"`
	Opentime   string `json:"opentime"`
	Lowtime    string `json:"lowtime"`
	Volume     string `json:"volume"`
	Drate      string `json:"drate"`
	Hightime   string `json:"hightime"`
	Jinivolume string `json:"jinivolume"`
	Msvolume   string `json:"msvolume"`
	Open       string `json:"open"`
	Status     string `json:"status"`
}

func (s *StockDetail) ToStockAggsDetail() (*ws.StockAggsDetail, error) {

	avg, err := strconv.ParseFloat(s.W_avrg, 64)
	if err != nil {
		return nil, err
	}

	min, err := strconv.ParseFloat(s.Low, 64)
	if err != nil {
		return nil, err
	}

	max, err := strconv.ParseFloat(s.High, 64)
	if err != nil {
		return nil, err
	}

	start, err := strconv.ParseFloat(s.Open, 64)
	if err != nil {
		return nil, err
	}

	end, err := strconv.ParseFloat(s.Price, 64)
	if err != nil {
		return nil, err
	}

	startTime, err := strconv.ParseInt(s.Opentime, 10, 64)
	if err != nil {
		return nil, err
	}

	endTime, err := strconv.ParseInt(s.Chetime, 10, 64)
	if err != nil {
		return nil, err
	}

	return &ws.StockAggsDetail{
		Average:   avg,
		Min:       min,
		Max:       max,
		Start:     start,
		End:       end,
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

type ResponseJson struct {
	Header HeaderJson       `json:"header"`
	Body   ResponseBodyJson `json:"body"`
}
