package ioc

import (
	"github.com/asaskevich/EventBus"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/fs"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/health"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/messaging"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/file"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/node"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/persistence/mongo"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/rpc"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/startup"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/web"
	"github.com/russell-kvashnin/kraken/internal/pkg/log"
	"go.uber.org/fx"
)

// Startup scenario providers
var StartupModule = fx.Options(
	fx.Provide(PrepareConnectionsStageProvider),
	fx.Provide(PrepareFileSystemStageProvider),
	fx.Provide(RegisterEventsStageProvider),
	fx.Provide(RegisterNodeStageProvider),
	fx.Provide(ConfigureServersStageProvider),
	fx.Provide(ScenarioProvider),
)

// Prepare connections stage provider
func PrepareConnectionsStageProvider(log *log.Logger, mongo *mongo.Mongo, rabbit *messaging.Rabbit) *startup.PrepareConnectionsStage {
	return startup.NewPrepareConnectionsStage(log, mongo, rabbit)
}

// Prepare file system stage provider
func PrepareFileSystemStageProvider(fs *fs.FileService, cfg config.FSConfig) *startup.PrepareFileSystemStage {
	stage := startup.NewPrepareFileSystemStage(fs, cfg)

	return stage
}

// Register node stage dependencies
type RegisterNodeStageParams struct {
	fx.In

	Cfg    config.NodeConfig
	Coll   *health.LoadStatsCollector
	Repo   node.Repository
	Rabbit *messaging.Rabbit
}

// Register node stage
func RegisterNodeStageProvider(p RegisterNodeStageParams) *startup.RegisterNodeStage {
	d := startup.RegisterNodeStageDeps{
		Cfg:       p.Cfg,
		Collector: p.Coll,
		Repo:      p.Repo,
		Rabbit:    p.Rabbit,
	}

	return startup.NewRegisterNodeStage(d)
}

// Servers configuration stage parameters
type ConfigureServersStageParams struct {
	fx.In

	Api       *web.Server
	Health    *health.Worker
	Mirroring *messaging.MessageConsumer `name:"mirroring_consumer"`
	Producer  *messaging.MessageProducer
	Rpc       *rpc.Server
}

// Servers configuration stage provider
func ConfigureServersStageProvider(p ConfigureServersStageParams) *startup.ConfigureServersStage {
	d := startup.ConfigureServersStageDeps{
		Api:       p.Api,
		Health:    p.Health,
		Mirroring: p.Mirroring,
		Producer:  p.Producer,
		Rpc:       p.Rpc,
	}

	return startup.NewConfigureServersStages(d)
}

// Register events stage dependencies
type RegisterEventsStageParams struct {
	fx.In

	Bus                  EventBus.Bus
	UploadedEventHandler file.UploadedEventHandler
}

// Register events stage provider
func RegisterEventsStageProvider(p RegisterEventsStageParams) *startup.RegisterEventsStage {
	d := startup.RegisterEventsStageDeps{
		Bus:                 p.Bus,
		UploadedFileHandler: p.UploadedEventHandler,
	}

	return startup.NewRegisterEventsStage(d)
}

// Scenario provider dependencies set
type ScenarioProviderParams struct {
	fx.In

	Log *log.Logger

	Cs *startup.PrepareConnectionsStage
	Fs *startup.PrepareFileSystemStage
	Es *startup.RegisterEventsStage
	Rs *startup.RegisterNodeStage
	Sc *startup.ConfigureServersStage
}

// Scenario provider
func ScenarioProvider(p ScenarioProviderParams) *startup.Scenario {
	deps := startup.ScenarioDeps{
		Log: p.Log,
		Cs:  p.Cs,
		Es:  p.Es,
		Fs:  p.Fs,
		Rs:  p.Rs,
		Sc:  p.Sc,
	}

	return startup.NewScenario(deps)
}
