package mongo

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/file"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/node"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
)

// Mongo collection for nodes
const NodeCollectionName = "nodes"

// Node repository implementation
type NodeRepository struct {
	ConcurrentRepository

	nodeId    string
	clusterId string
}

// Node repository constructor
func NewNodeRepository(mongo *Mongo, cfg config.NodeConfig) *NodeRepository {
	repo := new(NodeRepository)
	repo.mongo = mongo
	repo.dbName = mongo.GetCurrentDBName()
	repo.cName = NodeCollectionName

	repo.nodeId = cfg.Id
	repo.clusterId = cfg.ClusterId

	return repo
}

// Get node from database by nodeId
func (repo *NodeRepository) Get(nodeId string) (node.Node, error) {
	var (
		model node.Node
		e     error
	)

	err := repo.C(func(c *mgo.Collection) {
		err := c.Find(
			bson.M{
				"node_id": nodeId,
			},
		).One(&model)

		if err != nil {
			details := make(map[string]string)
			details["node_id"] = model.Id
			details["cluster_id"] = model.ClusterId

			e = kerr.NewErr(kerr.ErrLvlError, file.ErrorDomain, node.NotExistsErrorCode, err, details)
		}
	})

	if err != nil {
		return node.Node{}, err
	}

	return model, e
}

// Update node information
func (repo *NodeRepository) Update(model node.Node) error {
	var (
		e error
	)

	err := repo.C(func(c *mgo.Collection) {
		err := c.Update(
			bson.M{
				"node_id": model.Id,
			},
			&model)

		if err != nil {
			details := make(map[string]string)
			details["node_id"] = model.Id
			details["cluster_id"] = model.ClusterId

			e = kerr.NewErr(kerr.ErrLvlError, file.ErrorDomain, node.UpdateErrorCode, err, details)
		}
	})

	if err != nil {
		return err
	}

	return e
}

// Register current node in nodes collection
func (repo *NodeRepository) Register(model node.Node) error {
	var (
		e error
	)

	err := repo.C(func(c *mgo.Collection) {
		_, err := c.Upsert(
			bson.M{
				"node_id": model.Id,
			},
			model,
		)

		if err != nil {
			details := make(map[string]string)
			details["node_id"] = model.Id
			details["cluster_id"] = model.ClusterId

			e = kerr.NewErr(kerr.ErrLvlFatal, file.ErrorDomain, node.RegistrationErrorCode, err, details)
		}
	})

	if err != nil {
		return err
	}

	return e
}

// Returns current node document
func (repo *NodeRepository) CurrentNode() (node.Node, error) {
	return repo.Get(repo.nodeId)
}

// Get mirroring node
func (repo *NodeRepository) GetMirroringNode() (node.Node, error) {
	var (
		model node.Node
		e     error
	)

	err := repo.C(func(c *mgo.Collection) {
		err := c.Find(
			bson.M{
				"node_id": bson.M{
					"$ne": repo.nodeId,
				},
				"cluster_id": repo.clusterId,
			},
		).One(&model)

		if err != nil {
			e = kerr.NewErr(kerr.ErrLvlError, ErrorDomain, node.NotExistsErrorCode, err, nil)
		}
	})

	if err != nil {
		return model, err
	}

	return model, nil
}
