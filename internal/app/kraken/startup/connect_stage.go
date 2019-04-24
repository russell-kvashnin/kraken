package startup

import (
	"github.com/russell-kvashnin/kraken/internal/app/kraken/messaging"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/persistence/mongo"
	"github.com/russell-kvashnin/kraken/internal/pkg/log"
)

// Establish connection stage
// Prepare mongo db connection, rabbit, etc
type PrepareConnectionsStage struct {
	log    *log.Logger
	mongo  *mongo.Mongo
	rabbit *messaging.Rabbit
}

// Constructor
func NewPrepareConnectionsStage(log *log.Logger, mongo *mongo.Mongo, rabbit *messaging.Rabbit) *PrepareConnectionsStage {
	stage := new(PrepareConnectionsStage)
	stage.log = log
	stage.mongo = mongo
	stage.rabbit = rabbit

	return stage
}

// Execute stage
func (stage *PrepareConnectionsStage) Execute() error {
	var err error

	if err = stage.connectMongo(); err != nil {
		return err
	}

	if err = stage.connectRabbit(); err != nil {
		return err
	}

	return nil
}

// Obtain mongodb connection
func (stage *PrepareConnectionsStage) connectMongo() error {
	cChan := make(chan mongo.ConnectInfo)
	var (
		isSuccess bool
		lastErr   error
	)

	go stage.mongo.Connect(cChan)

	for c := range cChan {
		switch c.Done {
		case false:
			lastErr = c.Error

			stage.log.Error(lastErr)
		case true:
			isSuccess = true
			break
		}
	}

	if !isSuccess {
		return lastErr
	}

	return lastErr
}

// Obtain rabbitmq connection
func (stage *PrepareConnectionsStage) connectRabbit() error {
	cChan := make(chan messaging.ConnectInfo)
	var (
		isSuccess bool
		lastErr   error
	)

	go stage.rabbit.Connect(cChan)

	for c := range cChan {
		switch c.Done {
		case false:
			lastErr = c.Error

			stage.log.Error(lastErr)
		case true:
			isSuccess = true
			break
		}
	}

	if !isSuccess {
		return lastErr
	}

	return lastErr
}
