package kis

import (
	"log"
)

const custtype string = "P"
const tr_type string = "1"

func (s *Subscriber) run() {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:

		}
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			log.Println("Error while reading message")
			continue
		}

		receivedString := string(message)

		agg, err := ToStockDetail_USA(receivedString)
		if err != nil {
			log.Println("Error while converting to StockDetail_USA")
			continue
		}
		data, err := ToStockAggsDetail_USA(agg)
		if err != nil {
			log.Println("Error while converting to StockAggsDetail_USA")
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
				TrType:      tr_type,
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

func (s *Subscriber) UnsubscribeStockAggs(stocks ...string) error {
	return nil
}
