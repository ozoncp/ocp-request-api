package producer

import (
	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	producerBatchSize = 10000
	sendMessagesEvery = time.Second
)

type Producer interface {
	Send(msg EventMsg)
}

func NewProducer(topic string, kafkaProducer sarama.SyncProducer) Producer {
	p := &producer{topic: topic, kafkaProducer: kafkaProducer}
	p.Init()
	return p
}

type producer struct {
	topic         string
	kafkaProducer sarama.SyncProducer
	queue         chan EventMsg
}

func (p *producer) Init() {
	p.queue = make(chan EventMsg, producerBatchSize)
	ticker := time.NewTicker(sendMessagesEvery)
	// TODO For simplicity we don't have "ensure all sent before shutdown" mechanism. Will add later.
	go func() {
		batch := make([]EventMsg, 0, producerBatchSize)
		for {
			select {
			case msg, ok := <-p.queue:
				if !ok {
					p.sendMessages(batch)
					batch = batch[:0]
					return
				} else {
					batch = append(batch, msg)
				}
				if len(batch) >= producerBatchSize {
					p.sendMessages(batch)
				}
			case <-ticker.C:
				p.sendMessages(batch)
				batch = batch[:0]
			}
		}
	}()
}

func (p *producer) Send(msg EventMsg) {
	p.queue <- msg
}

func (p *producer) sendMessages(msgs []EventMsg) {

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
