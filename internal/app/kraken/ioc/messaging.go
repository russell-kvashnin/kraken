package ioc

import (
	"context"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/messaging"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/node"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/rpc/client"
	"go.uber.org/fx"
)

// RabbitMQ messaging module
// Provides consumers, producers, etc
var MessagingModule = fx.Options(
	fx.Provide(RabbitMQWrapperProvider),
	fx.Provide(MessageConsumersProvider),
	fx.Provide(MessageProducerProvider),
	fx.Provide(MirroringMessageHandlerProvider),
)

// RabbitMQ wrapper provider
func RabbitMQWrapperProvider(cfg config.RabbitConfig) *messaging.Rabbit {
	return messaging.NewRabbit(cfg)
}

// Dependencies for message consumer
type MessageConsumersParams struct {
	fx.In

	Lc fx.Lifecycle

	Rabbit  *messaging.Rabbit
	Config  config.MirroringConfig
	Handler *messaging.MirroringMessageHandler
}

// Message consumers providers result
type MessageConsumersResult struct {
	fx.Out

	MirroringConsumer *messaging.MessageConsumer `name:"mirroring_consumer"`
}

// RabbitMQ message consumers provider
func MessageConsumersProvider(p MessageConsumersParams) MessageConsumersResult {
	mirroring := messaging.NewMessageConsumer(p.Config.Queue, p.Rabbit, p.Handler)

	p.Lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go mirroring.Run()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := mirroring.Stop(ctx)
			if err != nil {
				return err
			}

			return nil
		},
	})

	return MessageConsumersResult{
		MirroringConsumer: mirroring,
	}
}

// Dependencies for message producer
type MessageProducerParams struct {
	fx.In

	Lc fx.Lifecycle

	Rabbit *messaging.Rabbit
}

// RabbitMQ message producer provider
func MessageProducerProvider(p MessageProducerParams) *messaging.MessageProducer {
	producer := messaging.NewMessageProducer(p.Rabbit)

	p.Lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			err := producer.Shutdown(ctx)
			if err != nil {
				return err
			}

			return nil
		},
	})

	return messaging.NewMessageProducer(p.Rabbit)
}

// MirroringMessageHandler provider
func MirroringMessageHandlerProvider(repo node.Repository, client *client.MirroringClient) *messaging.MirroringMessageHandler {
	handler := messaging.NewMirroringMessageHandler(repo, client)

	return handler
}
