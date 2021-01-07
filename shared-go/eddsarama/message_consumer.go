package eddsarama

import (
	"context"

	"github.com/Shopify/sarama"

	"github.com/stackus/edat/log"
	"github.com/stackus/edat/msg"
)

type messageConsumer struct {
	ctx        context.Context
	marshaller Serializer
	consumer   func(ctx context.Context, message msg.Message) error
	logger     log.Logger
}

func (messageConsumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (messageConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (c *messageConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case saramaMsg, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			message, err := c.marshaller.Deserialize(saramaMsg)
			if err != nil {
				c.logger.Error("message failed to unmarshal", log.Error(err))
				return err
			}

			err = c.consumer(c.ctx, message)
			if err != nil {
				c.logger.Error("consumer returned an error", log.Error(err))
				break
			}
			sess.MarkMessage(saramaMsg, "")
		case <-c.ctx.Done():
			// handle graceful shutdowns initiated externally
			c.logger.Trace("context cancelled; shutting down claims consumer")
			return nil
		}
	}
}
