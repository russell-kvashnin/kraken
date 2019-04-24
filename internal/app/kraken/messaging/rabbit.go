package messaging

import (
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
	"github.com/streadway/amqp"
	"time"
)

// Mongo error codes
const (
	ErrorDomain = "RABBITMQ"

	DialErrCode            = "RABBIT_DIAL_ERROR"
	NotConnectedErrCode    = "RABBIT_NOT_CONNECTED"
	ObtainChannelErrorCode = "RABBIT_FAILED_TO_OBTAIN_CHANNEL"
)

// RabbitMQ wrapper
type Rabbit struct {
	cfg  config.RabbitConfig
	conn *amqp.Connection
}

type ConnectInfo struct {
	Error error
	Done  bool
}

// Rabbit wrapper constructor
func NewRabbit(cfg config.RabbitConfig) *Rabbit {
	rabbit := new(Rabbit)
	rabbit.cfg = cfg

	return rabbit
}

// Obtain rabbitmq connection
func (rabbit *Rabbit) Connect(c chan ConnectInfo) {
	defer close(c)

	var (
		conn     *amqp.Connection
		err      error
		attempts int
	)

	uri := rabbit.cfg.GetRabbitMqUri()

	for {
		conn, err = amqp.Dial(uri)
		attempts++

		if err != nil {
			details := make(map[string]string)
			details["uri"] = uri

			e := kerr.NewErr(kerr.ErrLvlFatal, ErrorDomain, DialErrCode, err, details)
			c <- ConnectInfo{
				Error: e,
				Done:  false,
			}

			if rabbit.cfg.Reconnect == false {
				break
			}

			timeout := rabbit.cfg.ReconnectTimeout * time.Second
			time.Sleep(timeout)

			continue
		}

		break
	}

	rabbit.conn = conn

	c <- ConnectInfo{
		Error: nil,
		Done:  true,
	}
}

// Open's new channel
func (rabbit *Rabbit) GetChannel() (*amqp.Channel, error) {
	if rabbit.conn == nil {
		e := kerr.NewErr(kerr.ErrLvlError, ErrorDomain, NotConnectedErrCode, nil, nil)

		return nil, e
	}

	ch, err := rabbit.conn.Channel()
	if err != nil {
		e := kerr.NewErr(kerr.ErrLvlFatal, ErrorDomain, ObtainChannelErrorCode, err, nil)

		return nil, e
	}

	return ch, nil
}
