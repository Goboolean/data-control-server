package kafkapub

import (
	"errors"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/bson"
)


func (p *Producer) SendDataSync(topic string, data StockAggregate) error {

	bsonData, err := bson.Marshal(data)
	
	if err != nil {
		return err
	}

	if err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          bsonData,
	}, nil); err != nil {
		return err
	}

	return nil
}



func (p *Producer) SendDataBatch(topic string, batch []StockAggregate) error {

	msgChan := p.producer.ProduceChannel()

	bsonBatch := make([][]byte, len(batch))


	for idx := range batch {
		data, err := bson.Marshal(batch[idx])

		if err != nil {
			return err
		}

		bsonBatch[idx] = data
	}

	for idx := range batch {
		msgChan <- &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          bsonBatch[idx],
		}
	}

	if e := p.producer.Flush(int(time.Minute)); e != 0  {
		return errors.New("failed some flush")
	}

	return nil
}