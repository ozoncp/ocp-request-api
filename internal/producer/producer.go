package producer

import (
	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
)

type Producer interface {
	Send(msg ...EventMsg)
	Close() error
}

// NewProducer Returns new kafka producer
func NewProducer(topic string, kafkaProducer sarama.SyncProducer) Producer {
	p := &producer{topic: topic, kafkaProducer: kafkaProducer}
	return p
}

type producer struct {
	topic         string
	kafkaProducer sarama.SyncProducer
}

// Send sends a batch of message to Kafka broker
func (p *producer) Send(msgs ...EventMsg) {
	if len(msgs) == 0 {
		return
	}
	preped := make([]*sarama.ProducerMessage, 0, len(msgs))
	for _, m := range msgs {
		preped = append(
			preped,
			&sarama.ProducerMessage{
				Topic:     p.topic,
				Partition: -1,
				Value:     m,
			},
		)
	}
	err := p.kafkaProducer.SendMessages(preped)
	if err != nil {
		log.Error().Msgf("failed to send messages to Kafka: %v", err)
	}
}

// Close closes makes sure all send requests are completed
// and closes producer and underlying client.
func (p *producer) Close() error {
	return p.kafkaProducer.Close()
}
