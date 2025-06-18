package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"strings"
	"time"
)

const (
	flushtimeout = 5000
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(address []string) *Producer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
	})
	if err != nil {
		log.Fatalln(err)
	}
	return &Producer{producer: producer}
}

func (p *Producer) Produce(msg, topic string) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value:         []byte(msg),
		Key:           nil,
		Timestamp:     time.Time{},
		TimestampType: 0,
		Opaque:        nil,
		Headers:       nil,
	}
	kafkaChan := make(chan kafka.Event)
	err := p.producer.Produce(kafkaMsg, kafkaChan)
	if err != nil {
		return err
	}
	e := <-kafkaChan
	switch e.(type) {
	case kafka.Error:
		return e.(kafka.Error)
	case *kafka.Message:
		return nil
	default:
		return UknownType
	}
}
func (p *Producer) Close() {
	p.producer.Flush(flushtimeout)
	p.producer.Close()

}
