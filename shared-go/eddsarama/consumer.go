package eddsarama

import (
	"context"
	"sync"

	"github.com/Shopify/sarama"

	"github.com/stackus/edat/log"
	"github.com/stackus/edat/msg"
)

// Consumer implements msg.Consumer
type Consumer struct {
	brokers    []string
	groupID    string
	config     *sarama.Config
	serializer Serializer
	listenerWg sync.WaitGroup
	logger     log.Logger
}

var _ msg.Consumer = (*Consumer)(nil)

// NewConsumer constructs a new Consumer
func NewConsumer(brokers []string, groupID string, options ...ConsumerOption) *Consumer {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky

	c := &Consumer{
		brokers:    brokers,
		groupID:    groupID,
		config:     config,
		serializer: DefaultSerializer,
		logger:     log.DefaultLogger,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

// Listen implements msg.Consumer.Listen
func (c *Consumer) Listen(ctx context.Context, channel string, consumer msg.ReceiveMessageFunc) error {
	logger := c.logger.Sub(log.String("Channel", channel))

	defer logger.Trace("stopped listening")

	group, err := sarama.NewConsumerGroup(c.brokers, c.groupID, c.config)
	if err != nil {
		return err
	}

	msgConsumer := &messageConsumer{
		ctx:        ctx,
		marshaller: c.serializer,
		consumer:   consumer,
		logger:     logger,
	}

	go func() {
		for {
			err = group.Consume(ctx, []string{channel}, msgConsumer)
			if err != nil {
				logger.Error("error while consuming channel", log.Error(err))
				return
			}
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	logger.Trace("listening")

	<-ctx.Done()

	err = group.Close()
	if err != nil {
		logger.Error("error while closing listener", log.Error(err))
		return err
	}

	return nil
}

// Close implements msg.Consumer.Close
func (c *Consumer) Close(context.Context) error {
	c.logger.Trace("closing message source")

	return nil
}
