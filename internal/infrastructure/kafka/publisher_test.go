package kafka_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/fetch-server/internal/infrastructure/kafka"
	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/stretchr/testify/assert"
)

var (
	pub  *kafka.Publisher
)

func SetupPublisher() {
	var err error

	pub, err = kafka.NewPublisher(&resolver.ConfigMap{
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

	fmt.Println("do not cache")

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
	fmt.Println("do not cache plzzzz")

	const topic = "default-topic"

	var data = &kafka.StockAggregate{Average: 1.0,Min: 1.0, Max: 1.0}
	var dataBatch = []*kafka.StockAggregate{
		{Average: 1.0,Min: 1.1, Max: 1.2},
		{Average: 1.2,Min: 1.3, Max: 1.4},
		{Average: 1.4,Min: 1.5, Max: 1.6},
	}

	SetupPublisher()
	defer TeardownPublisher()

	t.Run("SendToExistingTopic", func(t *testing.T) {
		err := pub.SendData(topic, data)
		assert.NoError(t, err)
	})

	t.Run("SendToNonExistingTopic", func(t *testing.T) {
		t.Skip("Skip this test because auto.create.topics.enable is default true, want false")
		err := pub.SendData("non-existent-topic", data)
		assert.Error(t, err)
	})

	t.Run("SendDataBatch", func(t *testing.T) {
		err := pub.SendDataBatch(topic, dataBatch)
		assert.NoError(t, err)
	})

	time.Sleep(1 * time.Second)
}
