package kafkapub

import (
	"fmt"
	"os"
	"time"

	"github.com/Shopify/sarama"
)

var (
	KAFKA_HOST = os.Getenv("KAFKA_HOST")
	KAFKA_PORT = os.Getenv("KAFKA_PORT")
	KAFKA_USER = os.Getenv("KAFKA_USER")
	KAFKA_PASS = os.Getenv("KAFKA_PASS")
)



type Publisher struct {
	conn sarama.SyncProducer
}

var instance *Publisher



func New() *Publisher {

	if instance == nil {
		config := sarama.NewConfig()

		config.Net.SASL.Enable = true
		config.Net.SASL.User = KAFKA_USER
		config.Net.SASL.Password = KAFKA_PASS
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		config.Net.SASL.Version = sarama.SASLHandshakeV1
		config.Net.SASL.Handshake = true
		config.Net.TLS.Enable = false
	
		config.Producer.RequiredAcks = sarama.WaitForAll
		config.Producer.Retry.Max = 5
		config.Producer.Timeout = 5 * time.Second
	
		KAFKA_ADDR := fmt.Sprintf("%s:%s", KAFKA_HOST, KAFKA_PORT)
	
		producer, err := sarama.NewSyncProducer([]string{KAFKA_ADDR}, config)

		if err != nil {
			panic(err)
		}

		instance = &Publisher{conn: producer}
	}

	return instance
}
