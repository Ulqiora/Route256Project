package grpc_service

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"homework/internal/config"
	"homework/internal/infrastructure/kafka"
	"homework/internal/service/broker_io"
)

func PrintMessage(message *sarama.ConsumerMessage) {
	pm := broker_io.RequestMessage{}
	err := json.Unmarshal(message.Value, &pm)
	if err != nil {
		fmt.Println("Consumer error", err)
	}
	fmt.Println("Received Key: ", string(message.Key), " Value: ", pm)
}

func ConfigureReceivers(config *config.Config, topics ...string) (*broker_io.KafkaReceiver, error) {
	consumer, err := kafka.NewConsumer(config.Kafka)
	if err != nil {
		return nil, err
	}
	handlers := map[string]broker_io.HandleFunc{
		"order":     PrintMessage,
		"pickpoint": PrintMessage,
	}

	var receiver *broker_io.KafkaReceiver
	receiver = broker_io.NewReceiver(consumer, handlers)
	for _, topic := range topics {
		err = receiver.Subscribe(topic)
		if err != nil {
			return nil, errors.Wrap(err, "receiver subscribe error")
		}
	}
	return receiver, nil
}
