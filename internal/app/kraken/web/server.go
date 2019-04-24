package web

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/web/actions"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/web/handler"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"

	"net/http"
)

const (
	ErrorDomain = "API_SERVER"

	ListenErrCode = "HTTP_LISTEN_ERROR"
)

// Web server
type Server struct {
	http       *http.Server
	router     *chi.Mux
	errHandler *kerr.Handler

	Middlewares struct {
		PrepareUploadRequest   *actions.PrepareUploadRequest
		PrepareDownloadRequest *actions.PrepareDownloadRequest
	}

	Actions struct {
		Upload   http.Handler
		Download http.Handler
	}
}

// Web server dependencies
type ServerDeps struct {
	Router     *chi.Mux
	ErrHandler *kerr.Handler

	Middlewares struct {
		PrepareUploadRequest   *actions.PrepareUploadRequest
		PrepareDownloadRequest *actions.PrepareDownloadRequest
	}

	Actions struct {
		Upload   *actions.UploadAction
		Download *actions.DownloadAction
	}
}

// Web server constructor
func NewServer(d ServerDeps) *Server {
	server := new(Server)
	server.errHandler = d.ErrHandler
	server.router = d.Router
	server.Middlewares.PrepareUploadRequest = d.Middlewares.PrepareUploadRequest
	server.Middlewares.PrepareDownloadRequest = d.Middlewares.PrepareDownloadRequest
	server.Actions.Upload = handler.NewHttpHandler(d.ErrHandler, d.Actions.Upload.Execute)
	server.Actions.Download = handler.NewHttpHandler(d.ErrHandler, d.Actions.Download.Execute)

	return server
}

// Configure web server
func (server *Server) ConfigureRoutes() {
	server.router.Use(middleware.RequestID)
	server.router.Use(middleware.RealIP)
	server.router.Use(middleware.Logger)
	server.router.Use(middleware.Recoverer)

	server.router.
		With(server.Middlewares.PrepareUploadRequest.Execute).
		Method("POST", "/upload", server.Actions.Upload)

	server.router.
		With(server.Middlewares.PrepareDownloadRequest.Execute).
		Method("GET", "/{shortUrl}", server.Actions.Download)
}

// Run server
func (server *Server) Run(cfg config.WebConfig) {
	srvString := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	server.http = &http.Server{Addr: srvString, Handler: server.router}

	err := server.http.ListenAndServe()
	if err != nil {
		e := kerr.NewErr(kerr.ErrLvlFatal, ErrorDomain, ListenErrCode, err, nil)

		server.errHandler.Handle(e)
	}
}

// Stop server
func (server *Server) Stop(ctx context.Context) error {
	return server.http.Shutdown(ctx)
}
