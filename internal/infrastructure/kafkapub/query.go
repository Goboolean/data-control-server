package kafkapub

import (
	"github.com/Shopify/sarama"
	"go.mongodb.org/mongo-driver/bson"
)



func (p *Publisher) InsertStockList(topic string, batch []StockAggregate) error {

	message := &sarama.ProducerMessage{Topic: topic}


	messageList := []*sarama.ProducerMessage{};

	for idx := range batch {

		bsonData, err := bson.Marshal(batch[idx])

		if err != nil {
			return err
		}

		message = &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(bsonData),
		}

		messageList = append(messageList, message)
	}

	if err := p.conn.SendMessages(messageList); err != nil {
		return err
	}

	return nil
}