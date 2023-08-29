package kafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/protobuf/proto"

	log "github.com/sirupsen/logrus"
)

// produce.Flush should be finished within this time
var defaultFlushTimeout = time.Second * 3

type Publisher struct {
	producer *kafka.Producer

	wg sync.WaitGroup
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
		"acks":                -1,    // 0 if no response is required, 1 if only leader response is required, -1 if all in-sync replicas' response is required
		"go.delivery.reports": true, // Delivery reports (on delivery success/failure) will be sent on the Producer.Events() channel
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	instance := &Publisher{
		producer: producer,
	}

	

	go func() {
		instance.wg.Add(1)
		defer instance.wg.Done()

		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.WithFields(log.Fields{
						"topic": *ev.TopicPartition.Topic,
						"data":  ev.Value,
						"info":   ev.TopicPartition.Error,
					}).Error(ErrFailedToDeliveryData)
				} else {
					log.WithFields(log.Fields{
						"topic": *ev.TopicPartition.Topic,
						"data":  ev.Value,
						"info":   ev.TopicPartition.String(),
					}).Info("data arrived")
				}
			}
		}
	}()

	return instance, nil
}

// It should be called before program ends to free memory
func (p *Publisher) Close() {
	p.wg.Done()
	p.wg.Wait()
	//p.producer.Close()
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

	topic = prefix + topic

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

	//if left := p.producer.Flush(1); left != 0 {
	//	return fmt.Errorf("failed to flush all data, %d messages left", left)
	//}

	log.WithField("topic", topic).Debug("data sent")

	return nil
}

func (p *Publisher) SendDataBatch(topic string, batch []*StockAggregate) error {

	topic = prefix + topic

	msgChan := p.producer.ProduceChannel()

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

	return nil
}
