package startup

import (
	"github.com/asaskevich/EventBus"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/file"
)

// Registering event handlers
type RegisterEventsStage struct {
	bus EventBus.Bus

	handlers struct {
		file.UploadedEventHandler
	}
}

// Register events stage dependencies
type RegisterEventsStageDeps struct {
	Bus EventBus.Bus

	UploadedFileHandler file.UploadedEventHandler
}

// Register events stage constructor
func NewRegisterEventsStage(d RegisterEventsStageDeps) *RegisterEventsStage {
	stage := new(RegisterEventsStage)
	stage.bus = d.Bus
	stage.handlers = struct {
		file.UploadedEventHandler
	}{
		d.UploadedFileHandler,
	}

	return stage
}

// Execute stage
func (stage *RegisterEventsStage) Execute() error {
	err := stage.bus.SubscribeAsync(
		file.UploadedEventAlias,
		stage.handlers.UploadedEventHandler.Handle,
		false,
	)
	if err != nil {
		return err
	}

	return nil
}
