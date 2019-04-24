package ioc

import (
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
	"github.com/russell-kvashnin/kraken/internal/pkg/log"
	"go.uber.org/fx"
)

// Error & error handling module
var ErrorModule = fx.Options(
	fx.Provide(ErrorHandlerProvider),
)

// Error handler provider
func ErrorHandlerProvider(log *log.Logger) *kerr.Handler {
	handler := kerr.NewErrorHandler(log)

	return handler
}
