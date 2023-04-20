package kafkapub

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/bson"
)


func (p *Producer) SendData(topic string, data StockAggregate) error {

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

	for idx := range batch {

		bsonData, err := bson.Marshal(batch[idx])

		if err != nil {
			return err
		}

		msgChan <- &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          bsonData,
		}
	}

	return nil
}