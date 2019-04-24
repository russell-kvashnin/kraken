package file

import (
	"github.com/edsrzf/mmap-go"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/fs"
)

// Download file query
type DownloadFileQuery struct {
	ShortUrl string
}

// Handler for download file query
type DownloadFileQueryHandler struct {
	repo Repository
	fs   *fs.FileService
}

// Download file handler constructor
func NewDownloadFileHandler(repo Repository, fs *fs.FileService) *DownloadFileQueryHandler {
	handler := new(DownloadFileQueryHandler)
	handler.repo = repo
	handler.fs = fs

	return handler
}

// Handle query
func (handler *DownloadFileQueryHandler) Handle(query DownloadFileQuery) (*mmap.MMap, Meta, error) {
	doc, err := handler.repo.Get(query.ShortUrl)
	if err != nil {
		return nil, Meta{}, err

	}
	f, err := handler.fs.Read(doc.Path)
	if err != nil {
		return nil, Meta{}, err
	}

	return f, doc.Meta, nil
}
