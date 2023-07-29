package kis

import (
	"bytes"
	"log"

	"encoding/json"
	"io/ioutil"
	"net/http"
)


func (s *Subscriber) getApprovalKey(Appkey string, Secretkey string) (string, error) {
	data := &getApprovalKeyReqeust{
		GrantType: "client_credentials",
		AppKey:    Appkey,
		SecretKey: Secretkey,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	response, err := http.Post("https://openapi.koreainvestment.com:9443/oauth2/Approval", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var res *getApprovalKeyResponse

	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	return res.ApprovalKey, nil
}


const (
	custtype            string = "P"
	tr_type_subscribe   string = "1"
	tr_type_unsubscribe string = "0"
)




func (s *Subscriber) run() {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:

		}

		_, message, err := s.conn.ReadMessage()
		if err != nil {
			if valid := isResponseValid(message); !valid {
				log.Println("Error while reading message")
			}
			continue
		}

		agg, err := NewStockAggs(string(message))
		if err != nil {
			log.Println("Error while converting to StockAggs")
			continue
		}

		data, err := agg.ToStockAggsDetail()
		if err != nil {
			log.Println("Error while converting to StockAggsDetail")
			continue
		}

		if err := s.r.OnReceiveStockAggs(data); err != nil {
			log.Println("Error in OnReceiveStockAggs")
		}
	}
}


func (s *Subscriber) SubscribeStockAggs(symbols ...string) error {
	for _, symbol := range symbols {
		req := &RequestJson{
			Header: HeaderJson{
				ApprovalKey: s.approval_key,
				Custtype:    custtype,
				TrType:      tr_type_subscribe,
				ContentType: "utf-8",
			},
			Body: RequestBodyJson{
				Input: RequestInputJson{
					TrId:  "HDFSCNT0",
					TrKey: symbol,
				},
			},
		}

		if err := s.conn.WriteJSON(req); err != nil {
			return err
		}

	}
	return nil
}


func (s *Subscriber) UnsubscribeStockAggs(symbols ...string) error {
	for _, symbol := range symbols {
		req := &RequestJson{
			Header: HeaderJson{
				ApprovalKey: s.approval_key,
				Custtype:    custtype,
				TrType:      tr_type_unsubscribe,
				ContentType: "utf-8",
			},
			Body: RequestBodyJson{
				Input: RequestInputJson{
					TrId:  "HDFSCNT0",
					TrKey: symbol,
				},
			},
		}

		if err := s.conn.WriteJSON(req); err != nil {
			return err
		}

	}
	return nil
}
