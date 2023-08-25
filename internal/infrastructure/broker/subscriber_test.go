package broker_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/fetch-server/internal/infrastructure/broker"
	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/stretchr/testify/assert"
)

var sub *broker.Subscriber

type SubscribeListenerImpl struct{
	ch chan<- *broker.StockAggregate
}

func (i *SubscribeListenerImpl) OnReceiveStockAggs(name string, data *broker.StockAggregate) {
	i.ch <- data
}

func NewSubscribeListener(ch chan *broker.StockAggregate) *SubscribeListenerImpl {
	return &SubscribeListenerImpl{
		ch: ch,
	}
}

var received = make(chan *broker.StockAggregate, 10)



func SetupSubscriber() {
	var err error

	sub, err = broker.NewSubscriber(&resolver.ConfigMap{
		"HOST":  os.Getenv("KAFKA_HOST"),
		"PORT":  os.Getenv("KAFKA_PORT"),
		"GROUP": "test",
	}, NewSubscribeListener(received))
	if err != nil {
		panic(err)
	}
}

func TeardownSubscriber() {
	sub.Close()
}

func Test_Subscriber(t *testing.T) {

	SetupSubscriber()
	defer TeardownSubscriber()

	t.Run("Ping", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancelFunc()

		err := sub.Ping(ctx)
		assert.NoError(t, err)
	})
}

func Test_Subscribe(t *testing.T) {

	var topic = "test-topic"

	SetupSubscriber()
	SetupPublisher()
	defer TeardownSubscriber()
	defer TeardownPublisher()

	t.Run("SubscribeNonExistantTopic", func(t *testing.T) {
		t.Skip("Skip this test because auto.create.topics.enable is default true, want false")
		err := sub.Subscribe("non-existent-topic")
		assert.Error(t, err)
	})

	t.Run("SubscribeExistantTopic", func(t *testing.T) {
		err := sub.Subscribe(topic)
		assert.NoError(t, err)

		err = pub.SendData(topic, &broker.StockAggregate{})
		assert.NoError(t, err)

		select {
		case <-time.After(3 * time.Second):
			assert.Fail(t, "timeout")
		case <-received:
			return
		}
	})
}



func Test_SubscribeSameGroup(t *testing.T) {

	var topic = "test-topic"
	var count = 5

	SetupPublisher()
	defer TeardownPublisher()

	chanA1 := make(chan *broker.StockAggregate)
	subA1, err := broker.NewSubscriber(&resolver.ConfigMap{
		"HOST":  os.Getenv("KAFKA_HOST"),
		"PORT":  os.Getenv("KAFKA_PORT"),
		"GROUP": "A",
	}, NewSubscribeListener(chanA1))
	if err != nil {
		panic(err)
	}
	defer subA1.Close()

	chanA2 := make(chan *broker.StockAggregate)
	subA2, err := broker.NewSubscriber(&resolver.ConfigMap{
		"HOST":  os.Getenv("KAFKA_HOST"),
		"PORT":  os.Getenv("KAFKA_PORT"),
		"GROUP": "A",
	}, NewSubscribeListener(chanA2))
	if err != nil {
		panic(err)
	}
	defer subA2.Close()

	t.Run("Subscribe", func(t *testing.T) {
		err := subA1.Subscribe(topic)
		assert.NoError(t, err)

		err = subA2.Subscribe(topic)
		assert.NoError(t, err)

		for i := 0; i < count; i++ {
			err := pub.SendData(topic, &broker.StockAggregate{})
			assert.NoError(t, err)
		}

		// both A1 A2 should receive at least ${count} messages
		for i := 0; i < count; i++ {
			select {
			case <-time.After(3 * time.Second):
				assert.Fail(t, "failed to receive all message")
				return
			case <-chanA1:
				continue
			case <-chanA2:
				continue
			}
		}

		// both A1 A2 should not receive any more messages
		select {
		case <-chanA1:
			assert.Fail(t, "received more than expected")			
		case <-chanA2:
			assert.Fail(t, "received more than expected")
		}
	})
}



func Test_SubscribeDifferentGroup(t *testing.T) {

	var topic = "test-topic"
	var count = 5

	chanA := make(chan *broker.StockAggregate)
	subA, err := broker.NewSubscriber(&resolver.ConfigMap{
		"HOST":  os.Getenv("KAFKA_HOST"),
		"PORT":  os.Getenv("KAFKA_PORT"),
		"GROUP": "A",
	}, NewSubscribeListener(chanA))
	if err != nil {
		panic(err)
	}
	defer subA.Close()

	chanB := make(chan *broker.StockAggregate)
	subB, err := broker.NewSubscriber(&resolver.ConfigMap{
		"HOST":  os.Getenv("KAFKA_HOST"),
		"PORT":  os.Getenv("KAFKA_PORT"),
		"GROUP": "B",
	}, NewSubscribeListener(chanB))
	if err != nil {
		panic(err)
	}
	defer subB.Close()

	t.Run("Subscribe", func(t *testing.T) {
		err := subA.Subscribe(topic)
		assert.NoError(t, err)

		err = subB.Subscribe(topic)
		assert.NoError(t, err)

		for i := 0; i < count; i++ {
			err := pub.SendData(topic, &broker.StockAggregate{})
			assert.NoError(t, err)
		}

		// A should receive at least ${count} messages
		for i := 0; i < count; i++ {
			select {
			case <-time.After(3 * time.Second):
				assert.Fail(t, "failed to receive all message")
				return
			case <-chanA:
				continue
			}
		}

		// A should not receive any more messages
		select {
		case <-chanA:
			assert.Fail(t, "received more than expected")
		}


		// B should receive at least ${count} messages
		for i := 0; i < count; i++ {
			select {
			case <-time.After(3 * time.Second):
				assert.Fail(t, "failed to receive all message")
				return
			case <-chanB:
				continue
			}
		}
		
		// B should not receive any more messages
		select {
		case <-chanB:
			assert.Fail(t, "received more than expected")
		}
	})
}