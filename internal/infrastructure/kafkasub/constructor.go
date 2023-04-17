package gokafka

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)


var (
	KAFKA_HOST = os.Getenv("KAFKA_HOST")
	KAFKA_PORT = os.Getenv("KAFKA_PORT")
	KAFKA_USER = os.Getenv("KAFKA_USER")
	KAFKA_PASS = os.Getenv("KAFKA_PASS")
	KAFKA_ADDR = fmt.Sprintf("%s:%s", KAFKA_HOST, KAFKA_PORT)
)



type Subscriber struct {
	topic string

	config *kafka.ConfigMap
}



func NewChannel(stock string) *Subscriber {

	config := &kafka.ConfigMap{
		"bootstrap.servers":  KAFKA_ADDR,
		"security.protocol": "SASL_SSL",
		"sasl.mechanism":    "PLAIN",
		"sasl.username":     KAFKA_USER,
		"sasl.password":     KAFKA_PASS,
		"auto.offset.reset":        "earliest",
		"socket.keepalive.enable":  true,
	}

	return &Subscriber{
		topic: stock,
		config: config,
	}
}
