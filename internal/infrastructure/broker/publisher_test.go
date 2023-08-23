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

var (
	pub  *broker.Publisher
	data = &broker.StockAggregate{}
)

func SetupPublisher() {
	var err error

	pub, err = broker.NewPublisher(&resolver.ConfigMap{
		"HOST": os.Getenv("KAFKA_HOST"),
		"PORT": os.Getenv("KAFKA_PORT"),
	})
	if err != nil {
		panic(err)
	}
}

func TeardownPublisher() {
	pub.Close()
}

func TestPublisher(t *testing.T) {

	SetupPublisher()
	defer TeardownPublisher()

	t.Run("Ping", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancelFunc()

		if err := pub.Ping(ctx); err != nil {
			t.Errorf("Ping() failed: %v", err)
		}
	})
}

func Test_SendData(t *testing.T) {

	const topic = "default-topic"

	SetupPublisher()
	defer TeardownPublisher()

	t.Run("SendToExistingTopic", func(t *testing.T) {
		err := pub.SendData(topic, data)
		assert.NoError(t, err)
	})

	t.Run("SendToNonExistingTopic", func(t *testing.T) {
		t.Skip("Skip this test because auto.create.topics.enable may be true")
		err := pub.SendData("non-existent-topic", data)
		assert.Error(t, err)
	})

	t.Run("SendDataBatch", func(t *testing.T) {
		var dataBatch = []*broker.StockAggregate{
			{}, {}, {},
		}

		err := pub.SendDataBatch(topic, dataBatch)
		assert.NoError(t, err)
	})

	time.Sleep(1 * time.Second)
}
