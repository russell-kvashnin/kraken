package startup

import (
	"github.com/russell-kvashnin/kraken/internal/pkg/log"
)

// Startup scenario
type Scenario struct {
	log *log.Logger

	stages struct {
		cs *PrepareConnectionsStage
		fs *PrepareFileSystemStage
		es *RegisterEventsStage
		rs *RegisterNodeStage
		sc *ConfigureServersStage
	}
}

// Startup scenario dependencies
type ScenarioDeps struct {
	Log *log.Logger

	Cs *PrepareConnectionsStage
	Fs *PrepareFileSystemStage
	Es *RegisterEventsStage
	Rs *RegisterNodeStage
	Sc *ConfigureServersStage
}

// Constructor
func NewScenario(d ScenarioDeps) *Scenario {
	s := new(Scenario)
	s.log = d.Log
	s.stages.cs = d.Cs
	s.stages.fs = d.Fs
	s.stages.es = d.Es
	s.stages.rs = d.Rs
	s.stages.sc = d.Sc

	return s
}

// Constructor
func (s *Scenario) Execute() error {
	s.log.Named("startup").Infow("Executing application startup process")

	err := s.stages.cs.Execute()
	if err != nil {
		return err
	}

	err = s.stages.fs.Execute()
	if err != nil {
		return err
	}

	err = s.stages.es.Execute()
	if err != nil {
		return err
	}

	err = s.stages.rs.Execute()
	if err != nil {
		return err
	}

	err = s.stages.sc.Execute()
	if err != nil {
		return err
	}

	s.log.Named("startup").Infow("Application ready")

	return nil
}
