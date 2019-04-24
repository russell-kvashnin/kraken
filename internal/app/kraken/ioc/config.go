package ioc

import (
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Application configurations module
// Provides configuration structs for all application services
var ConfigModule = fx.Options(
	fx.Provide(AppConfigProvider),
	fx.Provide(NodeConfigProvider),
	fx.Provide(FSConfigProvider),
	fx.Provide(WebConfigProvider),
	fx.Provide(RpcConfigProvider),
	fx.Provide(MirroringConfigProvider),
	fx.Provide(MongoConfigProvider),
	fx.Provide(RabbitConfigProvider),
)

// Global application config provider
func AppConfigProvider() (config.AppConfig, error) {
	cfg := config.AppConfig{}

	var configPath string

	pflag.String("config", config.DefaultPath, "Configuration file location")
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return config.AppConfig{}, err
	}

	configPath = viper.GetString("config")

	viper.SetConfigType(config.Format)
	viper.SetConfigFile(configPath)

	err = viper.ReadInConfig()
	if err != nil {
		e := kerr.NewErr(kerr.ErrLvlFatal, config.ErrorDomain, config.ReadErrorCode, err, nil)

		return cfg, e
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		e := kerr.NewErr(kerr.ErrLvlFatal, config.ErrorDomain, config.ParseErrorCode, err, nil)

		return cfg, e
	}

	return cfg, nil
}

// Node config provider
func NodeConfigProvider(cfg config.AppConfig) (config.NodeConfig, error) {
	return cfg.GetNodeConfig()
}

// FS config provider
func FSConfigProvider(cfg config.AppConfig) config.FSConfig {
	return cfg.GetFSConfig()
}

// Web server config provider
func WebConfigProvider(cfg config.AppConfig) config.WebConfig {
	return cfg.GetWebConfig()
}

// gRPC server config provider
func RpcConfigProvider(cfg config.AppConfig) config.RpcConfig {
	return cfg.GetRPCConfig()
}

// Mirroring config provider
func MirroringConfigProvider(cfg config.AppConfig) config.MirroringConfig {
	return cfg.GetMirroringConfig()
}

// MongoDB config provider
func MongoConfigProvider(cfg config.AppConfig) config.MongoConfig {
	return cfg.GetMongoConfig()
}

// RabbitMQ config provider
func RabbitConfigProvider(cfg config.AppConfig) config.RabbitConfig {
	return cfg.GetRabbitConfig()
}
