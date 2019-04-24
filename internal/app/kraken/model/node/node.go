package node

import (
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	"time"
)

// Node error codes
const (
	ErrorDomain = "NODE_DOMAIN"

	RegistrationErrorCode = "NODE_REGISTRATION_ERROR"
	NotExistsErrorCode    = "NODE_DOES_NOT_EXISTS"
	UpdateErrorCode       = "NODE_UPDATE_ERROR_CODE"
)

// Node document
type Node struct {
	Id        string `bson:"node_id"`
	ClusterId string `bson:"cluster_id"`

	Health    Health    `bson:"health"`
	Mirroring Mirroring `bson:"mirroring"`

	LastSeen time.Time `bson:"last_seen"`
}

// Node health data
type Health struct {
	NetworkStats float64 `bson:"network_stats"`
	DiskUsage    float64 `bson:"disk_usage"`
	CPUStats     float64 `bson:"cpu_stats"`
}

// Node mirroring settings
type Mirroring struct {
	QueueName string `bson:"queue_name"`
	Endpoint  string `bson:"endpoint"`
}

// Node document constructor
func NewNode(cfg config.NodeConfig) Node {
	node := new(Node)
	node.Id = cfg.Id
	node.ClusterId = cfg.ClusterId
	node.Mirroring.QueueName = cfg.Mirroring.Queue
	node.Mirroring.Endpoint = cfg.Mirroring.Endpoint

	return *node
}

// Node model repository interface
type Repository interface {
	Get(nodeId string) (Node, error)
	Update(node Node) error
	Register(node Node) error
	GetMirroringNode() (Node, error)
}
