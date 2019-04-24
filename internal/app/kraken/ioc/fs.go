package ioc

import (
	"github.com/russell-kvashnin/kraken/internal/app/kraken/fs"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	"go.uber.org/fx"
)

// File system module providers
var FSModule = fx.Options(
	fx.Provide(FileServiceProvider),
)

// File service provider
func FileServiceProvider(cfg config.FSConfig) *fs.FileService {
	fileService := fs.NewFileService(cfg)

	return fileService
}
