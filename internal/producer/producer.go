package producer

import (
	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

const (
	producerBatchSize = 10000
	sendMessagesEvery = time.Second
)

type Producer interface {
	Send(msg EventMsg)
	Close()
}

// NewProducer Returns new kafka producer
func NewProducer(topic string, kafkaProducer sarama.SyncProducer) Producer {
	p := &producer{topic: topic,
		kafkaProducer: kafkaProducer,
		waitGroup:     &sync.WaitGroup{}}
	p.init()
	return p
}

type producer struct {
	topic         string
	kafkaProducer sarama.SyncProducer
	queue         chan EventMsg
	waitGroup     *sync.WaitGroup
}

// init initializes producer
func (p *producer) init() {
	p.queue = make(chan EventMsg, producerBatchSize)
	ticker := time.NewTicker(sendMessagesEvery)
	go func() {
		p.waitGroup.Add(1)
		batch := make([]EventMsg, 0, producerBatchSize)
		for {
			select {
			case msg, ok := <-p.queue:
				if !ok {
					p.sendMessages(batch)
					batch = batch[:0]
					p.waitGroup.Done()
					return
				} else {
					batch = append(batch, msg)
				}
				if len(batch) >= producerBatchSize {
					p.sendMessages(batch)
					batch = batch[:0]
				}
			case <-ticker.C:
				p.sendMessages(batch)
				batch = batch[:0]
			}
		}
	}()
}

// Send sends message to Kafka broker
func (p *producer) Send(msg EventMsg) {
	p.queue <- msg
}

// Close makes sure all messages are sent to Kafka and stops internal producer goroutine
func (p *producer) Close() {
	close(p.queue)
	p.waitGroup.Wait()
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
