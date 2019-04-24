package ioc

import (
	"context"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/fs"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/file"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/rpc"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/rpc/client"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/rpc/service"
	"go.uber.org/fx"
)

// RPC module for node communication
var RpcModule = fx.Options(
	fx.Provide(RpcServerProvider),
	fx.Provide(MirroringServiceProvider),
	fx.Provide(MirroringServiceClientProvider),
)

// Rpc server dependencies
type RpcServerParams struct {
	fx.In

	Lc fx.Lifecycle

	Cfg       config.RpcConfig
	Mirroring *service.MirroringService
}

// Rpc server provider
func RpcServerProvider(p RpcServerParams) *rpc.Server {
	server := rpc.NewServer(p.Cfg, p.Mirroring)

	p.Lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go server.Run()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Stop()
		},
	})

	return server
}

// Mirroring service provider
func MirroringServiceProvider(repo file.Repository, fs *fs.FileService) *service.MirroringService {
	return service.NewMirroringService(repo, fs)
}

// Mirroring service client provider
func MirroringServiceClientProvider(fs *fs.FileService) *client.MirroringClient {
	return client.NewMirroringClient(fs)
}
