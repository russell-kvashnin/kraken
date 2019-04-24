package startup

import (
	"github.com/russell-kvashnin/kraken/internal/app/kraken/health"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/messaging"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/rpc"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/web"
)

// Stage dependencies
type ConfigureServersStageDeps struct {
	Api       *web.Server
	Health    *health.Worker
	Mirroring *messaging.MessageConsumer
	Producer  *messaging.MessageProducer
	Rpc       *rpc.Server
}

// Configuring api server, system workers, etc
type ConfigureServersStage struct {
	api       *web.Server
	health    *health.Worker
	mirroring *messaging.MessageConsumer
	producer  *messaging.MessageProducer
	rpc       *rpc.Server
}

// Stage constructor
func NewConfigureServersStages(d ConfigureServersStageDeps) *ConfigureServersStage {
	stage := new(ConfigureServersStage)
	stage.api = d.Api
	stage.health = d.Health
	stage.mirroring = d.Mirroring
	stage.producer = d.Producer
	stage.rpc = d.Rpc

	return stage
}

// Execute stage
func (stage *ConfigureServersStage) Execute() error {
	stage.api.ConfigureRoutes()

	stage.health.Configure()

	err := stage.mirroring.Configure()
	if err != nil {
		return err
	}

	err = stage.producer.Configure()
	if err != nil {
		return err
	}

	err = stage.rpc.Configure()
	if err != nil {
		return err
	}

	return nil
}
