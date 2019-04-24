package config

import (
	"fmt"
	"time"
)

const (
	DefaultPath = "/etc/kraken/kraken.conf"

	Format = "toml"

	ErrorDomain = "CONFIG_DOMAIN"

	ReadErrorCode  = "CONFIG_READ_ERROR"
	ParseErrorCode = "CONFIG_PARSE_ERROR"
)

// Raw config
type AppConfig struct {
	NodeId       string `mapstructure:"node_id"`
	ClusterId    string `mapstructure:"cluster_id"`
	HealthTicker string `mapstructure:"health_ticker"`

	UploadDir string `mapstructure:"upload_dir"`
	LogFile   string `mapstructure:"log_file"`

	WebHost string `mapstructure:"web_host"`
	WebPort int    `mapstructure:"web_port"`

	RpcHost string `mapstructure:"rpc_host"`
	RpcPort int    `mapstructure:"rpc_port"`

	MirroringQueue    string `mapstructure:"mirroring_queue"`
	MirroringEndpoint string `mapstructure:"mirroring_endpoint"`

	MongoHost             string        `mapstructure:"mongo_host"`
	MongoPort             int           `mapstructure:"mongo_port"`
	MongoDbName           string        `mapstructure:"mongo_db_name"`
	MongoUser             string        `mapstructure:"mongo_user"`
	MongoPassword         string        `mapstructure:"mongo_password"`
	MongoMaxConnections   int           `mapstructure:"mongo_max_connections"`
	MongoReconnect        bool          `mapstructure:"mongo_reconnect"`
	MongoReconnectTimeout time.Duration `mapstructure:"mongo_reconnect_delay"`

	RabbitHost             string        `mapstructure:"rabbit_host"`
	RabbitPort             int           `mapstructure:"rabbit_port"`
	RabbitUser             string        `mapstructure:"rabbit_user"`
	RabbitPassword         string        `mapstructure:"rabbit_password"`
	RabbitReconnect        bool          `mapstructure:"rabbit_reconnect"`
	RabbitReconnectTimeout time.Duration `mapstructure:"rabbit_reconnect_timeout"`
}

// Extract node config
func (cfg AppConfig) GetNodeConfig() (NodeConfig, error) {
	duration, err := time.ParseDuration(cfg.HealthTicker)
	if err != nil {
		return NodeConfig{}, err
	}

	return NodeConfig{
		Id:           cfg.NodeId,
		ClusterId:    cfg.ClusterId,
		HealthTicker: duration,
		Mirroring: MirroringConfig{
			Queue:    cfg.MirroringQueue,
			Endpoint: cfg.MirroringEndpoint,
		},
	}, nil
}

// Extract FS config
func (cfg AppConfig) GetFSConfig() FSConfig {
	return FSConfig{
		UploadDir: cfg.UploadDir,
		LogFile:   cfg.LogFile,
	}
}

// Extract http server config
func (cfg AppConfig) GetWebConfig() WebConfig {
	return WebConfig{
		Host: cfg.WebHost,
		Port: cfg.WebPort,
	}
}

// Extract rpc server config
func (cfg AppConfig) GetRPCConfig() RpcConfig {
	return RpcConfig{
		Host: cfg.RpcHost,
		Port: cfg.RpcPort,
	}
}

// Extract mirroring config
func (cfg AppConfig) GetMirroringConfig() MirroringConfig {
	return MirroringConfig{
		Queue:    cfg.MirroringQueue,
		Endpoint: cfg.MirroringEndpoint,
	}
}

// Extract Mongo config
func (cfg AppConfig) GetMongoConfig() MongoConfig {
	return MongoConfig{
		Host:             cfg.MongoHost,
		Port:             cfg.MongoPort,
		DBName:           cfg.MongoDbName,
		User:             cfg.MongoUser,
		Password:         cfg.MongoPassword,
		MaxSessions:      cfg.MongoMaxConnections,
		Reconnect:        cfg.MongoReconnect,
		ReconnectTimeout: cfg.MongoReconnectTimeout,
	}
}

// Extract RabbitMQ config
func (cfg AppConfig) GetRabbitConfig() RabbitConfig {
	return RabbitConfig{
		Host:             cfg.RabbitHost,
		Port:             cfg.RabbitPort,
		User:             cfg.RabbitUser,
		Password:         cfg.RabbitPassword,
		Reconnect:        cfg.RabbitReconnect,
		ReconnectTimeout: cfg.RabbitReconnectTimeout,
	}
}

// Node configuration
type NodeConfig struct {
	Id           string
	ClusterId    string
	HealthTicker time.Duration
	Mirroring    MirroringConfig
}

// File system config
type FSConfig struct {
	UploadDir string
	LogFile   string
}

// Web config
type WebConfig struct {
	Host string
	Port int
}

// Rpc config
type RpcConfig struct {
	Host string
	Port int
}

// Compile Rpc listen address
func (cfg *RpcConfig) GetRpcListenAddress() string {
	var uri string

	uri = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	return uri
}

// Mirroring config
type MirroringConfig struct {
	Queue    string
	Endpoint string
}

// MongoDB config
type MongoConfig struct {
	Host             string
	Port             int
	DBName           string
	User             string
	Password         string
	MaxSessions      int
	Reconnect        bool
	ReconnectTimeout time.Duration
}

// Assemble MongoDB uri from config
func (cfg MongoConfig) GetMongoUri() string {
	var uri string

	if cfg.User == "" || cfg.Password == "" {
		uri = fmt.Sprintf("mongodb://%s:%d/%s", cfg.Host, cfg.Port, cfg.DBName)
	} else {
		uri = fmt.Sprintf(
			"mongodb://%s:%s@%s:%d/%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
		)
	}

	return uri
}

// RabbitMQ config
type RabbitConfig struct {
	Host             string
	Port             int
	User             string
	Password         string
	Reconnect        bool
	ReconnectTimeout time.Duration
}

// Assemble RabbitMQ uri from config
func (cfg RabbitConfig) GetRabbitMqUri() string {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d",
		cfg.User, cfg.Password, cfg.Host, cfg.Port,
	)

	return uri
}
