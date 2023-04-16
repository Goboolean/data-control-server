package akafkapub

import (
	"github.com/Shopify/sarama"
	"go.mongodb.org/mongo-driver/bson"
)


func (p *Publisher) InsertStockPub(topic string, batch []StockAggregate) error {

	for idx := range batch {

		bsonData, err := bson.Marshal(batch[idx])

		if err != nil {
			return err
		}

		message := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(bsonData),
		}

		p.conn.Input() <- message
	}

	if err := <-p.conn.Errors(); err != nil {
		return err
	}
	return nil
}
