package eddsarama

import (
	"github.com/Shopify/sarama"

	"github.com/stackus/edat/msg"
)

var DefaultSerializer = SaramaSerializer{}

type Serializer interface {
	Serialize(channel string, message msg.Message) (*sarama.ProducerMessage, error)
	Deserialize(message *sarama.ConsumerMessage) (msg.Message, error)
}

type SaramaSerializer struct{}

func (SaramaSerializer) Serialize(channel string, message msg.Message) (*sarama.ProducerMessage, error) {
	msg := &sarama.ProducerMessage{
		Topic:   channel,
		Headers: make([]sarama.RecordHeader, 0, len(message.Headers())),
		Value:   sarama.ByteEncoder(message.Payload()),
	}

	for key, value := range message.Headers() {
		msg.Headers = append(msg.Headers, sarama.RecordHeader{
			Key:   []byte(key),
			Value: []byte(value),
		})
	}

	return msg, nil
}

func (SaramaSerializer) Deserialize(message *sarama.ConsumerMessage) (msg.Message, error) {
	var id string

	headers := make(map[string]string, len(message.Headers))

	for _, header := range message.Headers {
		if string(header.Key) == msg.MessageID {
			id = string(header.Value)
		} else {
			headers[string(header.Key)] = string(header.Value)
		}
	}

	return msg.NewMessage(message.Value, msg.WithMessageID(id), msg.WithHeaders(headers)), nil
}
