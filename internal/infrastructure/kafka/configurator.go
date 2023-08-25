package kafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
)

// Configurator has a role for making and deleting topic, checking topic exists, and getting topic list.
type Configurator struct {
	AdminClient *kafka.AdminClient
}

var (
	conf *Configurator
	once sync.Once
)

func NewConfigurator(c *resolver.ConfigMap) (*Configurator, error) {

	host, err := c.GetStringKey("HOST")
	if err != nil {
		return nil, err
	}

	port, err := c.GetStringKey("PORT")
	if err != nil {
		return nil, err
	}

	debug, err := c.GetBoolKeyOptional("DEBUG")
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf("%s:%s", host, port)

	config := &kafka.ConfigMap{
		"bootstrap.servers": address,
	}

	if debug {
		config.SetKey("debug", "security, broker")
	}

	admin, err := kafka.NewAdminClient(config)

	if err != nil {
		return nil, err
	}

	return &Configurator{AdminClient: admin}, nil
}

// It should be called before program ends to free memory
func (c *Configurator) Close() {
	//c.AdminClient.Close()
}

// Check if connection to kafka is alive
func (c *Configurator) Ping(ctx context.Context) error {

	// It requires ctx to be deadline set, otherwise it will return error
	// It will return error if there is no response within deadline
	deadline, ok := ctx.Deadline()
	if !ok {
		return ErrTimeoutRequired
	}

	_, err := c.AdminClient.GetMetadata(nil, true, int(time.Until(deadline).Milliseconds()))
	return err
}

// Create a topic
func (c *Configurator) CreateTopic(ctx context.Context, topic string) error {

	// It returns error when topic already exists
	topic = packTopic(topic)

	exists, err := c.TopicExists(ctx, topic)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	topicInfo := kafka.TopicSpecification{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	result, err := c.AdminClient.CreateTopics(ctx, []kafka.TopicSpecification{topicInfo})
	if err != nil {
		return err
	}

	if err := result[0].Error; err.Code() != kafka.ErrNoError {
		return fmt.Errorf(err.String())
	}

	return nil
}

// Delete a topic
func (c *Configurator) DeleteTopic(ctx context.Context, topic string) error {

	// It returns error when the topic does not exist
	topic = packTopic(topic)

	result, err := c.AdminClient.DeleteTopics(ctx, []string{topic})

	if err != nil {
		return errors.Wrap(ErrFatalWhileDeletingTopic, err.Error())
	}

	if err := result[0].Error; err.Code() != kafka.ErrNoError {
		return errors.Wrap(ErrTrivalWhileDeletingTopic, err.String())
	}

	return nil
}

// delete all topic
func (c *Configurator) DeleteAllTopics(ctx context.Context) error {
	topicList, err := c.GetTopicList(ctx)
	if err != nil {
		return err
	}

	result, err := c.AdminClient.DeleteTopics(ctx, topicList)
	if err != nil {
		return err
	}

	if err := result[0].Error; err.Code() != kafka.ErrNoError {
		return fmt.Errorf(err.String())
	}

	return nil
}

// Check if given topic exists
func (c *Configurator) TopicExists(ctx context.Context, topic string) (bool, error) {
	topic = packTopic(topic)

	deadline, ok := ctx.Deadline()
	if !ok {
		return false, ErrTimeoutRequired
	}

	metadata, err := c.AdminClient.GetMetadata(nil, true, int(time.Until(deadline).Milliseconds()))
	if err != nil {
		return false, err
	}

	_, exists := metadata.Topics[topic]
	return exists, nil
}

// Get all existing topic list as a string slice
func (c *Configurator) GetTopicList(ctx context.Context) ([]string, error) {

	deadline, ok := ctx.Deadline()
	if !ok {
		return nil, ErrTimeoutRequired
	}

	metadata, err := c.AdminClient.GetMetadata(nil, true, int(time.Until(deadline).Milliseconds()))
	if err != nil {
		return nil, err
	}

	topicList := make([]string, 0)

	for topic := range metadata.Topics {
		if len(topic) > 0 {
			topicList = append(topicList, topic)
		}
	}

	return topicList, nil
}
