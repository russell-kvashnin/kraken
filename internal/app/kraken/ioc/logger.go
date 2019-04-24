package ioc

import (
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	"github.com/russell-kvashnin/kraken/internal/pkg/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger module providers
var LoggerModule = fx.Options(
	fx.Provide(LoggerProvider),
)

// Logger provider
func LoggerProvider(cfg config.AppConfig) (*log.Logger, error) {
	logCfg := zap.Config{
		Encoding: "json",

		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{cfg.LogFile},
		ErrorOutputPaths: []string{cfg.LogFile},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "timestamp",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	zlog, err := logCfg.Build()
	if err != nil {
		return nil, err
	}

	err = zlog.Sync()
	if err != nil {
		return nil, err
	}

	return log.NewLogger(zlog.Sugar()), nil
}
