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

var conf *kafka.Configurator

func SetupConfigurator() {

	var err error
	conf, err = kafka.NewConfigurator(&resolver.ConfigMap{
		"HOST": os.Getenv("KAFKA_HOST"),
		"PORT": os.Getenv("KAFKA_PORT"),
	})
	if err != nil {
		panic(err)
	}
}

func TeardownConfigurator() {
	conf.Close()
}

func Test_Configurator(t *testing.T) {

	SetupConfigurator()
	defer TeardownConfigurator()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := conf.Ping(ctx); err != nil {
		t.Errorf("Ping() = %v", err)
	}
}

func Test_CreateTopic(t *testing.T) {

	const topic = "test-topic"

	SetupConfigurator()
	defer TeardownConfigurator()

	t.Run("CreateTopic", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := conf.CreateTopic(ctx, topic)
		assert.NoError(t, err)

		exists, err := conf.TopicExists(ctx, topic)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("CreateExitingTopic", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := conf.CreateTopic(ctx, "existing-topic")
		assert.Error(t, err)
	})

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := conf.DeleteTopic(ctx, topic)
		assert.NoError(t, err)
	})
}

func Test_DeleteTopic(t *testing.T) {

	const topic = "test-topic"

	SetupConfigurator()
	defer TeardownConfigurator()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := conf.CreateTopic(ctx, topic)
	assert.NoError(t, err)

	t.Run("DeleteTopic", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := conf.DeleteTopic(ctx, topic)
		assert.NoError(t, err)

		exists, err := conf.TopicExists(ctx, topic)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("DeleteNonExistingTopic", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := conf.DeleteTopic(ctx, "non-existent-topic")
		assert.Error(t, err)
	})
}

func Test_GetTopicList(t *testing.T) {

	SetupConfigurator()
	defer TeardownConfigurator()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	topicList, err := conf.GetTopicList(ctx)
	assert.NoError(t, err)

	fmt.Printf("Topic Count: %d\n", len(topicList))
	fmt.Printf("Topic List: \n")
	for _, topic := range topicList {
		fmt.Println(topic)
	}
}

func Test_DeleteALlTopics(t *testing.T) {

	SetupConfigurator()
	defer TeardownConfigurator()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := conf.DeleteAllTopics(ctx)
	assert.NoError(t, err)
}
