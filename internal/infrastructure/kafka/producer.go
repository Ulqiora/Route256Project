package kafka

import (
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/Ulqiora/Route256Project/internal/config"
)

type BrokersIP []string

type Producer struct {
	brokerIP          BrokersIP
	producerSamara    sarama.AsyncProducer
	errorLogManager   *slog.Logger
	successLogManager *slog.Logger
}

func newAsyncProducer(config config.KafkaConfig) (sarama.AsyncProducer, error) {
	asyncProducerConfig := sarama.NewConfig()
	// выставление алгоритма распределения сообщений по партициям
	asyncProducerConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	// от какого количества партици ждем ответа о получении сообщения
	asyncProducerConfig.Producer.RequiredAcks = sarama.WaitForAll

	// при отправке сообщенений возвращаются сохраняется информация об успехе
	asyncProducerConfig.Producer.Return.Successes = true
	// тоже самое, но для ошибочной отправке
	asyncProducerConfig.Producer.Return.Errors = true

	asyncProducer, err := sarama.NewAsyncProducer(config.Hosts, asyncProducerConfig)
	if err != nil {
		return nil, err
	}
	return asyncProducer, nil
}

func NewProducer(config config.KafkaConfig, errorLogger, successLogger *slog.Logger) (*Producer, error) {
	asyncProducer, err := newAsyncProducer(config)
	if err != nil {
		return nil, fmt.Errorf("error with async kafka-producer: %w", err)
	}
	producer := &Producer{
		brokerIP:          config.Hosts,
		producerSamara:    asyncProducer,
		errorLogManager:   errorLogger,
		successLogManager: successLogger,
	}
	return producer, nil
}

func (k *Producer) SendAsyncMessage(message *sarama.ProducerMessage) {
	k.producerSamara.Input() <- message
}

func (k *Producer) Run() {
	go func() {
		// Error и Retry топики можно использовать при получении ошибки
		for e := range k.producerSamara.Errors() {
			msg := e.Error()
			k.errorLogManager.Info(msg)
		}
	}()
	// фиксация ошибок
	go func() {
		for m := range k.producerSamara.Successes() {
			msg := fmt.Sprintf(`Async message successfully sent with key: %s, partition: %d, offset: %d`,
				m.Key,
				m.Partition,
				m.Offset,
			)
			k.successLogManager.Info(msg)
		}
	}()
}

func (k *Producer) Close() error {
	err := k.producerSamara.Close()
	if err != nil {
		return fmt.Errorf("error close producer: %w", err)
	}
	return nil
}
