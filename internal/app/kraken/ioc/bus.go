package ioc

import (
	"github.com/asaskevich/EventBus"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/bus"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/messaging"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/file"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/node"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
	"go.uber.org/fx"
)

// Application events module
var BusModule = fx.Options(
	fx.Provide(EventBusProvider),
	fx.Provide(FileUploadedHandlerProvider),
)

// Event bus provider
func EventBusProvider() EventBus.Bus {
	return bus.NewEventBus()
}

// File uploaded handler dependencies
type FileUploadedHandlerParams struct {
	fx.In

	Repo       node.Repository
	Producer   *messaging.MessageProducer
	ErrHandler *kerr.Handler
}

// File uploaded event handler provider
func FileUploadedHandlerProvider(p FileUploadedHandlerParams) file.UploadedEventHandler {
	return bus.NewFileUploadedHandler(p.Repo, p.Producer, p.ErrHandler)
}
