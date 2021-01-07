package eddsarama

import (
	"github.com/Shopify/sarama"

	"github.com/stackus/edat/log"
)

type ProducerOption func(*Producer)

func WithProducerSerializer(serializer Serializer) ProducerOption {
	return func(producer *Producer) {
		producer.serializer = serializer
	}
}

func WithProducerConfig(fn func(cfg *sarama.Config)) ProducerOption {
	return func(producer *Producer) {
		fn(producer.config)
	}
}

func WithProducerLogger(logger log.Logger) ProducerOption {
	return func(producer *Producer) {
		producer.logger = logger
	}
}
