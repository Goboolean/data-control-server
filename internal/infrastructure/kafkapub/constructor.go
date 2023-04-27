package kafkapub

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



type Producer struct {
	producer *kafka.Producer
}



var instance *Producer

func init() {
	config := &kafka.ConfigMap{
		"bootstrap.servers":  KAFKA_ADDR,
		"security.protocol": "SASL_SSL",
		"sasl.mechanism":    "PLAIN",
		"sasl.username":     KAFKA_USER,
		"sasl.password":     KAFKA_PASS,
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		panic(err)
	}

	instance = &Producer{producer: producer}
}

func New() *Producer {
	return instance
}

func Close() error {
	instance.producer.Close()
	return nil
}
