package log

import "go.uber.org/zap"

// Logger wrapper
type Logger struct {
	*zap.SugaredLogger
}

// Constructor
func NewLogger(zap *zap.SugaredLogger) *Logger {
	logger := Logger{zap}

	return &logger
}
