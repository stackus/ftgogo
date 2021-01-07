package eddsarama

import (
	"github.com/Shopify/sarama"

	"github.com/stackus/edat/log"
)

type ConsumerOption func(*Consumer)

func WithConsumerSerializer(serializer Serializer) ConsumerOption {
	return func(consumer *Consumer) {
		consumer.serializer = serializer
	}
}

func WithConsumerConfig(fn func(cfg *sarama.Config)) ConsumerOption {
	return func(consumer *Consumer) {
		fn(consumer.config)
	}
}

func WithConsumerLogger(logger log.Logger) ConsumerOption {
	return func(consumer *Consumer) {
		consumer.logger = logger
	}
}
