package messaging

import (
	"context"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
	"github.com/streadway/amqp"
)

const (
	DeliveryErrorCode     = "DELIVERY_ERROR"
	ConsumerStopErrorCode = "CONSUMER_STOP_ERROR"
)

// Message handler interface
type MessageHandler interface {
	Handle(amqp.Delivery) error
}

// Message consumer server
type MessageConsumer struct {
	rabbit  *Rabbit
	channel *amqp.Channel
	handler MessageHandler

	delivery <-chan amqp.Delivery

	queue string
}

// Message consumer constructor
func NewMessageConsumer(queue string, rabbit *Rabbit, handler *MirroringMessageHandler) *MessageConsumer {
	consumer := new(MessageConsumer)
	consumer.queue = queue
	consumer.rabbit = rabbit
	consumer.handler = handler

	return consumer
}

// Run consumer
func (consumer *MessageConsumer) Configure() error {
	var err error

	consumer.channel, err = consumer.rabbit.GetChannel()
	if err != nil {
		return err
	}

	delivery, err := consumer.channel.Consume(
		consumer.queue,
		"mirroring",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		e := kerr.NewErr(kerr.ErrLvlError, ErrorDomain, DeliveryErrorCode, err, nil)

		return e
	}

	consumer.delivery = delivery

	return nil
}

// Run consumer server
func (consumer *MessageConsumer) Run() {
	for d := range consumer.delivery {
		go consumer.handler.Handle(d)
	}
}

// Stop consuming server
func (consumer *MessageConsumer) Stop(ctx context.Context) error {
	if consumer.channel == nil {
		return nil
	}

	err := consumer.channel.Close()
	if err != nil {
		e := kerr.NewErr(kerr.ErrLvlError, ErrorDomain, ConsumerStopErrorCode, err, nil)

		return e
	}

	return nil
}
