package service

import (
	"github.com/russell-kvashnin/kraken/internal/app/kraken/fs"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/file"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
	"github.com/russell-kvashnin/kraken/pkg/api/v1"
	"google.golang.org/grpc"
)

const (
	ErrorDomain         = "MIRRORING_DOMAIN"
	FileStreamErrorCode = "FILE_STREAM_ERROR"
)

// Mirroring gRPC service
type MirroringService struct {
	grpc *grpc.Server

	repo file.Repository
	fs   *fs.FileService
}

// Constructor
func NewMirroringService(repo file.Repository, fs *fs.FileService) *MirroringService {
	srv := new(MirroringService)
	srv.repo = repo
	srv.fs = fs

	return srv
}

// File download handler
func (srv *MirroringService) Download(request *v1.MirroringRequest, stream v1.FileMirroring_DownloadServer) error {
	model, err := srv.repo.Get(request.ShortUrl)
	if err != nil {
		return err
	}

	f, err := srv.fs.Read(model.Path)
	if err != nil {
		return err
	}

	response := &v1.FileChunk{
		Data: *f,
	}

	err = stream.Send(response)
	if err != nil {
		details := make(map[string]string)
		details["shortUrl"] = request.ShortUrl

		e := kerr.NewErr(kerr.ErrLvlError, ErrorDomain, FileStreamErrorCode, err, details)

		return e
	}

	err = f.Unmap()
	if err != nil {
		return err
	}

	return nil
}
