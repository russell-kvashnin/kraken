package ioc

import (
	"context"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/health"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/node"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
	"go.uber.org/fx"
)

// Node health module
// Provides utilization metrics, etc
var HealthModule = fx.Options(
	fx.Provide(LoadStatsCollectorProvider),
	fx.Provide(WorkerProvider),
)

// Load statistic collector provider
func LoadStatsCollectorProvider(errHandler *kerr.Handler) *health.LoadStatsCollector {
	collector := health.NewLoadStatsCollector(errHandler)

	return collector
}

// Worker dependencies structure
type WorkerParams struct {
	fx.In

	Lc fx.Lifecycle

	Cfg     config.NodeConfig
	Coll    *health.LoadStatsCollector
	Repo    node.Repository
	Handler *kerr.Handler
}

// Health check worker provider
func WorkerProvider(p WorkerParams) *health.Worker {
	worker := health.NewWorker(p.Cfg, p.Coll, p.Repo, p.Handler)

	p.Lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go worker.Run()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			worker.Stop()

			return nil
		},
	})

	return worker
}
