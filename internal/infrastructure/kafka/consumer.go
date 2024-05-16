package kafka

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/Ulqiora/Route256Project/internal/config"
)

type Consumer struct {
	brokers        []string
	SingleConsumer sarama.Consumer
}

func NewConsumer(configApp config.KafkaConfig) (*Consumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = false
	saramaConfig.Consumer.Offsets.AutoCommit.Enable = true
	saramaConfig.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumer(configApp.Hosts, saramaConfig)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		brokers:        configApp.Hosts,
		SingleConsumer: consumer,
	}, err
}
