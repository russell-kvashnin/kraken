package ioc

import (
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/file"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/node"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/persistence/mongo"
	"go.uber.org/fx"
)

// Persistence module
// Provide MongoDB service
var PersistenceModule = fx.Options(
	fx.Provide(MongoProvider),
	fx.Provide(FileRepositoryProvider),
	fx.Provide(NodeRepositoryProvider),
)

// Mongo service provider
func MongoProvider(cfg config.MongoConfig) *mongo.Mongo {
	return mongo.NewMongo(cfg)
}

// Provider for file repository mongodb implementation
func FileRepositoryProvider(m *mongo.Mongo) file.Repository {
	repo := mongo.NewFileRepository(m)

	return repo
}

// Provider for node repository mongodb implementation
func NodeRepositoryProvider(m *mongo.Mongo, cfg config.NodeConfig) node.Repository {
	repo := mongo.NewNodeRepository(m, cfg)

	return repo
}
