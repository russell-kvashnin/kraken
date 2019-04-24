package startup

import (
	"github.com/russell-kvashnin/kraken/internal/app/kraken/health"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/messaging"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/node"
	"time"
)

// Node registration stage
type RegisterNodeStage struct {
	cfg    config.NodeConfig
	coll   *health.LoadStatsCollector
	repo   node.Repository
	rabbit *messaging.Rabbit
}

// Dependencies struct
type RegisterNodeStageDeps struct {
	Cfg       config.NodeConfig
	Collector *health.LoadStatsCollector
	Repo      node.Repository
	Rabbit    *messaging.Rabbit
}

// Constructor
func NewRegisterNodeStage(d RegisterNodeStageDeps) *RegisterNodeStage {
	stage := new(RegisterNodeStage)
	stage.cfg = d.Cfg
	stage.coll = d.Collector
	stage.repo = d.Repo
	stage.rabbit = d.Rabbit

	return stage
}

// Execute registration stage
func (stage *RegisterNodeStage) Execute() error {
	err := stage.registerNode()
	if err != nil {
		return err
	}

	err = stage.registerExchange()
	if err != nil {
		return err
	}

	return nil
}

// Collect node health data, register node in DB
func (stage *RegisterNodeStage) registerNode() error {
	n := node.NewNode(stage.cfg)

	stats := stage.coll.CollectHealthData()
	n.Health = node.Health{
		NetworkStats: stats.NetworkStats,
		DiskUsage:    stats.DiskUsage,
		CPUStats:     stats.CPUStats,
	}
	n.LastSeen = time.Now()

	err := stage.repo.Register(n)
	if err != nil {

		return err
	}

	return nil
}

// Register cluster exchange, node queues
func (stage *RegisterNodeStage) registerExchange() error {
	ch, err := stage.rabbit.GetChannel()
	if err != nil {
		return err
	}

	// declare cluster exchange
	err = ch.ExchangeDeclare(
		stage.cfg.ClusterId,
		"direct",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {

		return err
	}

	_, err = ch.QueueDeclare(
		stage.cfg.Id,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		stage.cfg.Id,
		stage.cfg.Id,
		stage.cfg.ClusterId,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.Close()
	if err != nil {
		return err
	}

	return nil
}
