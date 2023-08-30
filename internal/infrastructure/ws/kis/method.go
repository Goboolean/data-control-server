package kis

import (
	"bytes"
	"errors"
	"io"
	"log"
	"time"

	"encoding/json"
	"net/http"
)


func (s *Subscriber) GetApprovalKey(Appkey string, Secretkey string) (string, error) {
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

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var res *getApprovalKeyResponse

	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	if res.ApprovalKey == "" {
		return res.ApprovalKey, errors.New("invalid request")
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
			break
		}


		_, message, err := s.conn.ReadMessage()
		if err != nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		// There are two types of response that can be catched here:
		// 1. StockAggs
		// 2. StockSubscriptionInfo
		// If it fails to convert to StockAggs, it may be StockSubscriptionInfo, therefore ignorable.

		agg, err := NewStockAggs(string(message))
		if err != nil {
			symbol, flag := parseToSubscriptionResponse(message)
			if !flag {
				log.Println("")
			} else {
				s.subscribed <- symbol
			}
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
				ApprovalKey: s.approvalKey,
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
				ApprovalKey: s.approvalKey,
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
