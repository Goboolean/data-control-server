package kafka_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/fetch-server/internal/infrastructure/kafka"
	_ "github.com/Goboolean/fetch-server/internal/util/env"
	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/sirupsen/logrus"
)

// Issue on kafka broker : each test inturrupts the other test.
// As the confluent-kafka-go library is a wrapper of the C library, it acts as a singleton,
// which means calling instance.Close() may close the connection for other tests,
// resulting in next instance.Close() call to close already closed connection.
// This issue should be fixed by configuring broker libraries singleton
// and only the last instance.Close() call should clear the resources.
// Temporary solution is to replace instance.Close() to do nothing.

func SetUp() {

	// Topic Policy:
	// As many services(configurator, publisher, subscriber) and many tests
	// in same package, we need to make sure that topic information does not interrupt the other test.

	const (
		existingTopic    = "existing-topic" // it is assured to be exist
		nonExistentTopic = "non-existent-topic" // it is assured to be not exist
		testTopic        = "test-topic" // it does not exist at first state, and it is assured to be deleted at last state
		defaultTopic     = "default-topic" // exists at first state, and it is assured to be exist at last state
	)

	conf, err := kafka.NewConfigurator(&resolver.ConfigMap{
		"HOST": os.Getenv("KAFKA_HOST"),
		"PORT": os.Getenv("KAFKA_PORT"),
	})

	if err != nil {
		panic(err)
	}
	defer conf.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Verify that "existing-topic" exists, if not, create it
	exists, err := conf.TopicExists(ctx, existingTopic)
	if err != nil {
		panic(err)
	}

	if !exists {
		if err := conf.CreateTopic(ctx, existingTopic); err != nil {
			panic(err)
		}
	}

	// Verify that "non-existent-topic" does not exist, if it does, delete it
	exists, err = conf.TopicExists(ctx, nonExistentTopic)
	if err != nil {
		panic(err)
	}

	if exists {
		if err := conf.DeleteTopic(ctx, nonExistentTopic); err != nil {
			panic(err)
		}
	}

	// Verify that "test-topic" does not exist, if it does, delete it
	exists, err = conf.TopicExists(ctx, testTopic)
	if err != nil {
		panic(err)
	}

	if exists {
		if err := conf.DeleteTopic(ctx, testTopic); err != nil {
			panic(err)
		}
	}

	// Verify that "default-topic" exist, if not, create it
	exists, err = conf.TopicExists(ctx, defaultTopic)
	if err != nil {
		panic(err)
	}

	if !exists {
		if err := conf.CreateTopic(ctx, defaultTopic); err != nil {
			panic(err)
		}
	}

}

func TestMain(m *testing.M) {

	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	SetUp()
	code := m.Run()
	os.Exit(code)
}
