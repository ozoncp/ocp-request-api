package producer

import (
	"github.com/Shopify/sarama"
	desc "github.com/ozoncp/ocp-request-api/pkg/ocp-request-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type EventType int

const (
	Create EventType = iota
	Read
	Update
	Delete
)

type EventMsg interface {
	sarama.Encoder
}

func NewEvent(requestId uint64, eventType EventType, err error) EventMsg {
	return &event{
		requestId: requestId,
		eventType: eventType,
		err:       err,
	}
}

type event struct {
	requestId uint64
	eventType EventType
	err       error
}

func (e *event) Encode() ([]byte, error) {
	message := &desc.RequestAPIEvent{
		RequestId: e.requestId,
	}
	if e.err != nil {
		message.Error = e.err.Error()
	}

	switch e.eventType {
	case Create:
		message.Event = desc.RequestAPIEvent_CREATE
	case Read:
		message.Event = desc.RequestAPIEvent_READ
	case Update:
		message.Event = desc.RequestAPIEvent_UPDATE
	case Delete:
		message.Event = desc.RequestAPIEvent_DELETE
	default:
		log.Panic().Msgf("unexpected event type: %v", e.eventType)
	}
	return proto.Marshal(message)
}

func (e *event) Length() int {
	// strange we have to encode again to get length...but leave for now as is
	data, _ := e.Encode()
	return len(data)
}
