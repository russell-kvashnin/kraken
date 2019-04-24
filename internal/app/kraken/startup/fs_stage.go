package startup

import (
	"github.com/russell-kvashnin/kraken/internal/app/kraken/fs"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	"os"
)

// Upload dir permissions
const UploadDirPerms = 0755

// Prepare file system stage
// Create necessary dirs, etc
type PrepareFileSystemStage struct {
	fs  *fs.FileService
	cfg config.FSConfig
}

// Constructor
func NewPrepareFileSystemStage(fs *fs.FileService, cfg config.FSConfig) *PrepareFileSystemStage {
	stage := new(PrepareFileSystemStage)
	stage.fs = fs
	stage.cfg = cfg

	return stage
}

// Execute stage
func (stage *PrepareFileSystemStage) Execute() error {
	if _, err := os.Stat(stage.cfg.UploadDir); os.IsNotExist(err) {

		err := stage.fs.Mkdir(stage.cfg.UploadDir, UploadDirPerms)
		if err != nil {
			return err
		}
	}

	return nil
}
