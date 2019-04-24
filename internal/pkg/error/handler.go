package error

import (
	"github.com/russell-kvashnin/kraken/internal/pkg/log"
)

// Error handler
type Handler struct {
	log *log.Logger
}

// Constructor
func NewErrorHandler(log *log.Logger) *Handler {
	handler := new(Handler)
	handler.log = log

	return handler
}

// Handles error
func (handler *Handler) Handle(err error) {
	switch err.(type) {
	case Error:
		handler.handleCustomErr(err.(Error))
	default:
		handler.log.Errorw(err.Error(), "err", err)
	}
}

// Handle custom error
func (handler *Handler) handleCustomErr(err Error) {
	switch err.Level() {
	case ErrLvlFatal:
		handler.log.Named(err.Domain()).
			Fatalw(
				err.Original().Error(),
				"code", err.Code(),
				"original", err.Original(),
				"details", err.Details())
	case ErrLvlError:
		handler.log.Named(err.Domain()).
			Errorw(
				err.Original().Error(),
				"code", err.Code(),
				"original", err.Original(),
				"details", err.Details())
	case ErrLvlWarning:
		handler.log.Named(err.Domain()).
			Warnw(
				err.Original().Error(),
				"code", err.Code(),
				"original", err.Original(),
				"details", err.Details())
	}
}
