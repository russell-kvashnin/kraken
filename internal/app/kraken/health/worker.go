package health

import (
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/node"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
	"time"
)

// Health check worker
type Worker struct {
	ticker *time.Ticker

	cfg     config.NodeConfig
	repo    node.Repository
	coll    *LoadStatsCollector
	handler *kerr.Handler
}

// Worker info struct for chan
type WorkerInfo struct {
	Error error
}

// Health check worker constructor
func NewWorker(cfg config.NodeConfig, coll *LoadStatsCollector, repo node.Repository, handler *kerr.Handler) *Worker {
	worker := new(Worker)
	worker.cfg = cfg
	worker.coll = coll
	worker.repo = repo
	worker.handler = handler

	return worker
}

// Collect and update node health data
func (worker *Worker) work() error {
	n, err := worker.repo.Get(worker.cfg.Id)
	if err != nil {
		return err
	}

	stats := worker.coll.CollectHealthData()

	n.Health.CPUStats = stats.CPUStats
	n.Health.DiskUsage = stats.DiskUsage
	n.Health.NetworkStats = stats.NetworkStats
	n.LastSeen = time.Now()

	err = worker.repo.Update(n)
	if err != nil {
		return err
	}

	return nil
}

// Configure worker
func (worker *Worker) Configure() {
	worker.ticker = time.NewTicker(worker.cfg.HealthTicker)
}

// Work work
func (worker *Worker) Run() {

	for range worker.ticker.C {
		err := worker.work()
		if err != nil {
			worker.handler.Handle(err)
		}
	}
}

// Stop worker
func (worker *Worker) Stop() {
	worker.ticker.Stop()
}
