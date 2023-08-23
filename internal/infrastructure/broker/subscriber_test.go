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

type SubscribeListenerImpl struct{}

var stockChan = make(chan *broker.StockAggregate)

func (i *SubscribeListenerImpl) OnReceiveStockAggs(name string, data *broker.StockAggregate) {
	stockChan <- data
}

func SetupSubscriber() {
	var err error

	sub, err = broker.NewSubscriber(&resolver.ConfigMap{
		"HOST":  os.Getenv("KAFKA_HOST"),
		"PORT":  os.Getenv("KAFKA_PORT"),
		"GROUP": "test",
	}, context.Background(), &SubscribeListenerImpl{})
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

	type args struct {
		topic string
		data  *broker.StockAggregate
	}

	tests := []struct {
		name string
		args args
		want *broker.StockAggregate
	}{
		{
			name: "",
			args: args{
				topic: topic,
				data: &broker.StockAggregate{
					Average: 1234,
				},
			},
			want: &broker.StockAggregate{
				Average: 1234,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
			defer cancel()

			err := sub.Subscribe(topic)
			assert.NoError(t, err)

			err = pub.SendData(tt.args.topic, tt.args.data)
			assert.NoError(t, err)

			select {
			case <-ctx.Done():
				t.Errorf("timeout: failed to receive data")
			case got := <-stockChan:
				assert.Equal(t, tt.want, got)
			}
		})
	}

	t.Run("SubscribeNonExistantTopic", func(t *testing.T) {
		err := sub.Subscribe("non-existent-topic")
		assert.Error(t, err)
	})
}
