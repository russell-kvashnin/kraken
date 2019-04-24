package presenter

import (
	"github.com/edsrzf/mmap-go"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/file"
	"net/http"
)

// Write file response from mmap
func FileResponse(w http.ResponseWriter, f *mmap.MMap, meta file.Meta) error {
	w.Header().Set("Content-Type", meta.FileType)

	_, err := w.Write(*f)
	if err != nil {
		return err
	}

	return nil
}
