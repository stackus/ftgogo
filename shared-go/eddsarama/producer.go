package eddsarama

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"

	"github.com/stackus/edat/log"
	"github.com/stackus/edat/msg"
)

// Producer implements msg.Producer
type Producer struct {
	serializer Serializer
	producer   sarama.SyncProducer
	config     *sarama.Config
	logger     log.Logger
}

var _ msg.Producer = (*Producer)(nil)

// NewProducer constructs a new Producer
func NewProducer(brokers []string, clientID string, options ...ProducerOption) (*Producer, error) {
	config := sarama.NewConfig()
	config.ClientID = clientID
	config.Producer.Return.Successes = true

	d := &Producer{
		producer:   nil,
		config:     config,
		serializer: DefaultSerializer,
		logger:     log.DefaultLogger,
	}

	for _, option := range options {
		option(d)
	}

	producer, err := sarama.NewSyncProducer(brokers, d.config)
	if err != nil {
		return nil, err
	}

	d.producer = producer

	return d, nil
}

// Send implements msg.Producer.Send
func (p *Producer) Send(_ context.Context, channel string, message msg.Message) error {
	logger := p.logger.Sub(
		log.String("Channel", channel),
	)

	if _, err := message.Headers().GetRequired(msg.MessageID); err != nil {
		return err
	}

	saramaMsg, err := p.serializer.Serialize(channel, message)
	if err != nil {
		logger.Error("failed to marshal message", log.Error(err))
		return fmt.Errorf("message could not be marshalled")
	}

	if _, _, err = p.producer.SendMessage(saramaMsg); err != nil {
		return err
	}

	logger.Trace("message sent to kafka")

	return nil
}

// Close implements msg.Producer.Close
func (p *Producer) Close(context.Context) error {
	p.logger.Trace("closing message destination")
	err := p.producer.Close()
	if err != nil {
		p.logger.Error("error closing message destination", log.Error(err))
	}
	return err
}
