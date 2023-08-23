package broker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/protobuf/proto"
)

type SubscribeListener interface {
	OnReceiveStockAggs(name string, stock *StockAggregate)
}

type Subscriber struct {
	consumer *kafka.Consumer

	listener SubscribeListener

	ctx context.Context
	t   *topicManager
}

func NewSubscriber(c *resolver.ConfigMap, ctx context.Context, lis SubscribeListener) (*Subscriber, error) {

	host, err := c.GetStringKey("HOST")
	if err != nil {
		return nil, err
	}

	port, err := c.GetStringKey("PORT")
	if err != nil {
		return nil, err
	}

	group, err := c.GetStringKey("GROUP")
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf("%s:%s", host, port)

	config := &kafka.ConfigMap{
		"bootstrap.servers":       address,
		"auto.offset.reset":       "earliest",
		"socket.keepalive.enable": true,
		"group.id":                group,
	}

	consumer, err := kafka.NewConsumer(config)

	if err != nil {
		return nil, err
	}

	instance := &Subscriber{
		consumer: consumer,
		listener: lis,
		ctx:      ctx,
		t:        newTopicManager(consumer),
	}

	go instance.subscribeMessage(ctx)

	return instance, nil
}

func (s *Subscriber) subscribeMessage(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:

		}

		msg, err := s.consumer.ReadMessage(-1)

		if err != nil {
			log.Fatalf("err: failed to read received message: %v", err)

		} else {
			var data StockAggregate

			if err := proto.Unmarshal(msg.Value, &data); err != nil {
				log.Fatalf("err: failed to deserialize message: %v", err)
			}

			topic := unpackTopic(*msg.TopicPartition.Topic)

			s.listener.OnReceiveStockAggs(topic, &data)
		}
	}
}

func (s *Subscriber) Close() {
	s.consumer.Close()
}

func (s *Subscriber) Ping(ctx context.Context) error {
	deadline, ok := ctx.Deadline()

	if !ok {
		return fmt.Errorf("timeout setting on ctx required")
	}

	remaining := time.Until(deadline)

	_, err := s.consumer.GetMetadata(nil, true, int(remaining.Milliseconds()))
	return err
}

func (s *Subscriber) Subscribe(stock string) error {
	stock = packTopic(stock)

	if err := s.t.addTopic(stock); err != nil {
		return err
	}

	return s.t.renewSupscription()
}

func (s *Subscriber) Unsubscribe(stock string) error {
	stock = packTopic(stock)

	if err := s.t.deleteTopic(stock); err != nil {
		return err
	}

	return s.t.renewSupscription()
}

type topicManager struct {
	consumer *kafka.Consumer

	topicList map[string]struct{}
}

func newTopicManager(consumer *kafka.Consumer) *topicManager {
	return &topicManager{
		consumer:  consumer,
		topicList: make(map[string]struct{}),
	}
}

func (t *topicManager) addTopic(topic string) error {

	_, exists := t.topicList[topic]
	if exists {
		return errors.New("topic already exists")
	}

	t.topicList[topic] = struct{}{}
	return nil
}

func (t *topicManager) deleteTopic(topic string) error {
	topic = packTopic(topic)

	_, exists := t.topicList[topic]
	if !exists {
		return errors.New("topic does not exist")
	}

	delete(t.topicList, topic)
	return nil
}

func (t *topicManager) getTopicList() []string {
	topicList := make([]string, 0)

	for k := range t.topicList {
		topicList = append(topicList, k)
	}

	return topicList
}

func (t *topicManager) renewSupscription() error {
	return t.consumer.SubscribeTopics(t.getTopicList(), nil)
}
