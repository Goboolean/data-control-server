package broker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/protobuf/proto"
)

// produce.Flush should be finished within this time
var defaultFlushTimeout = time.Second * 3

type Publisher struct {
	producer *kafka.Producer
}

func NewPublisher(c *resolver.ConfigMap) (*Publisher, error) {

	host, err := c.GetStringKey("HOST")
	if err != nil {
		return nil, err
	}

	port, err := c.GetStringKey("PORT")
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf("%s:%s", host, port)

	config := &kafka.ConfigMap{
		"bootstrap.servers":   address,
		"acks":                1,    // 0 if no response is required, 1 if only leader response is required, -1 if all in-sync replicas' response is required
		"go.delivery.reports": true, // Delivery reports (on delivery success/failure) will be sent on the Producer.Events() channel
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	return &Publisher{producer: producer}, nil
}

// It should be called before program ends to free memory
func (p *Publisher) Close() {
	p.producer.Close()
}

// Check if connection to kafka is alive
func (p *Publisher) Ping(ctx context.Context) error {

	// It requires ctx to be deadline set, otherwise it will return error
	// It will return error if there is no response within deadline
	deadline, ok := ctx.Deadline()

	if !ok {
		return fmt.Errorf("timeout setting on ctx required")
	}

	remaining := time.Until(deadline)

	_, err := p.producer.GetMetadata(nil, true, int(remaining.Milliseconds()))
	return err
}

func (p *Publisher) SendData(topic string, data *StockAggregate) error {

	topic = packTopic(topic)

	binaryData, err := proto.Marshal(data)

	if err != nil {
		return err
	}

	if err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, // Send data to every partition
		Value:          binaryData,
	}, nil); err != nil {
		return err
	}

	go func() {
		for e := range p.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
					log.Printf("binary data: %v\n", ev.Value)
				}
			}
		}
	}()

	if remain := p.producer.Flush(int(defaultFlushTimeout.Milliseconds())); remain != 0 {
		return fmt.Errorf("%d events is stil remaining: ", remain)
	}

	return nil
}

func (p *Publisher) SendDataBatch(topic string, batch []*StockAggregate) error {

	topic = packTopic(topic)

	msgChan := p.producer.ProduceChannel()

	topic = packTopic(topic)

	for idx := range batch {
		binaryData, err := proto.Marshal(batch[idx])

		if err != nil {
			return err
		}

		msgChan <- &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          binaryData,
		}
	}

	go func() {
		for e := range p.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Fatalf("Delivery failed: %v\n", ev.TopicPartition.Error)
					log.Fatalf("binary data: %v\n", ev.Value)
				}
			}
		}
	}()

	if remain := p.producer.Flush(int(defaultFlushTimeout.Milliseconds())); remain != 0 {
		return fmt.Errorf("%d events is stil remaining: ", remain)
	}

	return nil
}
