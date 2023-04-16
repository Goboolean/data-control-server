package kafkapub

import (
	"github.com/Shopify/sarama"
	"go.mongodb.org/mongo-driver/bson"
)

func (p *Publisher) InsertStockBatch(topic string, batch []StockAggregate) error {

	message := &sarama.ProducerMessage{Topic: topic}

	messageBatch := []*sarama.ProducerMessage{}

	for idx := range batch {

		bsonData, err := bson.Marshal(batch[idx])

		if err != nil {
			return err
		}

		message = &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(bsonData),
		}

		messageBatch = append(messageBatch, message)
	}

	if err := p.conn.SendMessages(messageBatch); err != nil {
		return err
	}

	return nil
}
