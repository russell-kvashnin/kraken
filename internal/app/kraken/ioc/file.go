package ioc

import (
	"github.com/asaskevich/EventBus"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/fs"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/messaging"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/file"
	"go.uber.org/fx"
)

// File domain module
var FileModule = fx.Options(
	fx.Provide(UploadFileHandlerProvider),
	fx.Provide(DownloadFileProvider),
)

// Upload file handler dependencies
type UploadFileHandlerParams struct {
	fx.In

	Fs       *fs.FileService
	Repo     file.Repository
	Producer *messaging.MessageProducer
	Bus      EventBus.Bus
}

// Upload file handler provider
func UploadFileHandlerProvider(p UploadFileHandlerParams) *file.UploadFileHandler {
	d := file.UploadFileHandlerDeps{
		Fs:       p.Fs,
		Repo:     p.Repo,
		Producer: p.Producer,
		Bus:      p.Bus,
	}

	handler := file.NewUploadFileHandler(d)

	return handler
}

// Download file handler provider
func DownloadFileProvider(repo file.Repository, fs *fs.FileService) *file.DownloadFileQueryHandler {
	handler := file.NewDownloadFileHandler(repo, fs)

	return handler
}
