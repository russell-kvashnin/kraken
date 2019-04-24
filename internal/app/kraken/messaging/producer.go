package messaging

import (
	"context"
	"github.com/json-iterator/go"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
	"github.com/streadway/amqp"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	MessageContentType     = "application/json"
	MessageContentEncoding = "utf-8"

	ProduceErrorCode = "FAILED_TO_PRODUCE_MESSAGE"
)

// RabbitMQ message producer
type MessageProducer struct {
	rabbit  *Rabbit
	channel *amqp.Channel

	exchange string
}

// Message producer constructor
func NewMessageProducer(rabbit *Rabbit) *MessageProducer {
	producer := new(MessageProducer)
	producer.rabbit = rabbit

	return producer
}

// Configure producer
func (producer *MessageProducer) Configure() error {
	var err error

	producer.channel, err = producer.rabbit.GetChannel()
	if err != nil {
		return err
	}

	return nil
}

// Produce message
func (producer *MessageProducer) Produce(key string, message interface{}) error {
	msgBytes, err := json.Marshal(&message)
	if err != nil {
		return err
	}

	pub := amqp.Publishing{
		ContentType:     MessageContentType,
		ContentEncoding: MessageContentEncoding,
		Body:            msgBytes,
	}

	err = producer.channel.Publish(
		producer.exchange,
		key,
		false,
		false,
		pub,
	)

	if err != nil {
		details := make(map[string]string)
		details["exchange"] = producer.exchange
		details["routing_key"] = key
		details["message"] = string(msgBytes)

		e := kerr.NewErr(kerr.ErrLvlError, ErrorDomain, ProduceErrorCode, err, nil)

		return e
	}

	return nil
}

// Shut down producer
func (producer *MessageProducer) Shutdown(ctx context.Context) error {
	if producer.channel == nil {
		return nil
	}

	err := producer.channel.Close()
	if err != nil {
		return err
	}

	return nil
}
