package broker_io

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/Ulqiora/Route256Project/internal/infrastructure/kafka"
	"github.com/gofrs/uuid"
)

type SenderKafka struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaSender(producer *kafka.Producer, topic string) *SenderKafka {
	return &SenderKafka{
		producer,
		topic,
	}
}

func (s *SenderKafka) SendMessage(message RequestMessage) {
	kafkaMsg, err := s.buildMessage(message)
	if err != nil {
		fmt.Println("Send message marshal error", err)
	}

	s.producer.SendAsyncMessage(kafkaMsg)
}

func (s *SenderKafka) buildMessage(message RequestMessage) (*sarama.ProducerMessage, error) {
	msg, err := json.Marshal(message)

	if err != nil {
		fmt.Println("Send message marshal error", err)
		return nil, err
	}

	return &sarama.ProducerMessage{
		Topic:     s.topic,
		Value:     sarama.ByteEncoder(msg),
		Partition: -1,
		Key:       sarama.StringEncoder(fmt.Sprint(uuid.NewV4())),
	}, nil
}
