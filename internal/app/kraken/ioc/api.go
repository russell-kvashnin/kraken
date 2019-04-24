package ioc

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/file"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/web"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/web/actions"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
	"github.com/russell-kvashnin/kraken/internal/pkg/log"
	"go.uber.org/fx"
)

// Rest api module
// Provides api specific stuff: web server, actions, etc
var ApiModule = fx.Options(
	fx.Provide(RouterProvider),
	fx.Provide(ApiServerProvider),
	fx.Provide(PrepareUploadRequestProvider),
	fx.Provide(UploadFileActionProvider),
	fx.Provide(PrepareDownloadRequestProvider),
	fx.Provide(DownloadActionProvider),
)

// Chi router provider
func RouterProvider() *chi.Mux {
	router := chi.NewRouter()

	return router
}

// Api server params
type ApiServerParams struct {
	fx.In

	Lc fx.Lifecycle

	Cfg        config.WebConfig
	Log        *log.Logger
	R          *chi.Mux
	ErrHandler *kerr.Handler

	PrepareUploadRequest   *actions.PrepareUploadRequest
	PrepareDownloadRequest *actions.PrepareDownloadRequest

	UploadAction   *actions.UploadAction
	DownloadAction *actions.DownloadAction
}

func ApiServerProvider(p ApiServerParams) *web.Server {
	d := web.ServerDeps{
		Router:     p.R,
		ErrHandler: p.ErrHandler,

		Middlewares: struct {
			PrepareUploadRequest   *actions.PrepareUploadRequest
			PrepareDownloadRequest *actions.PrepareDownloadRequest
		}{PrepareUploadRequest: p.PrepareUploadRequest, PrepareDownloadRequest: p.PrepareDownloadRequest},

		Actions: struct {
			Upload   *actions.UploadAction
			Download *actions.DownloadAction
		}{Upload: p.UploadAction, Download: p.DownloadAction},
	}

	server := web.NewServer(d)

	p.Lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go server.Run(p.Cfg)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := server.Stop(ctx)
			if err != nil {
				return err
			}

			return nil
		},
	})

	return server
}

// Prepare upload request middleware
func PrepareUploadRequestProvider(cfg config.NodeConfig) *actions.PrepareUploadRequest {
	middleware := actions.NewPrepareUploadRequest(cfg)

	return middleware
}

// Upload action provider
func UploadFileActionProvider(handler *file.UploadFileHandler) *actions.UploadAction {
	return actions.NewUploadAction(handler)
}

// Prepare download request middleware provider
func PrepareDownloadRequestProvider() *actions.PrepareDownloadRequest {
	return actions.NewPrepareDownloadRequest()
}

// Download action provider
func DownloadActionProvider(handler *file.DownloadFileQueryHandler) *actions.DownloadAction {
	return actions.NewDownloadAction(handler)
}
