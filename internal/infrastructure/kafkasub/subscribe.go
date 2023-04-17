package gokafka

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/bson"
)



func (s *Subscriber) SubscribeStock() (<-chan StockAggregate, error) {

	consumer, err := kafka.NewConsumer(s.config)

	if err != nil {
		return nil, err
	}

	if err := consumer.Subscribe(s.topic, nil); err != nil {
		return nil, err
	}

	msgChan := make(chan StockAggregate)

	go func() {
		defer close(msgChan)
		for {
			msg, err := consumer.ReadMessage(-1)

			if err != nil {
				log.Fatal(err)
			} else {
				data := StockAggregate{}

				if err := bson.Unmarshal(msg.Value, &data); err != nil {
					log.Fatal(err)
				}
				
				msgChan <- data
			}
		}
	}()

	return msgChan, nil
}