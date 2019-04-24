package client

import (
	"context"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/fs"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
	"github.com/russell-kvashnin/kraken/pkg/api/v1"
	"google.golang.org/grpc"
)

const (
	ErrorDomain          = "GRPC_CLIENT_DOMAIN"
	DialServiceErrorCode = "GRPC_DIAL_SERVICE_ERROR"
	DownloadErrorCode    = "GRPC_DOWNLOAD_ERROR"
)

// gRPC mirroring service client
type MirroringClient struct {
	conn *grpc.ClientConn

	fs *fs.FileService
}

// Constructor
func NewMirroringClient(fs *fs.FileService) *MirroringClient {
	client := new(MirroringClient)
	client.fs = fs

	return client
}

// Download file
func (client *MirroringClient) Download(endpoint string, shortUrl string) error {
	var err error

	client.conn, err = grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		details := make(map[string]string)
		details["enpoint"] = endpoint
		details["shortUrl"] = shortUrl

		e := kerr.NewErr(kerr.ErrLvlError, ErrorDomain, DialServiceErrorCode, err, details)

		return e
	}
	defer client.conn.Close()

	c := v1.NewFileMirroringClient(client.conn)

	req := &v1.MirroringRequest{
		ShortUrl: shortUrl,
	}

	res, err := c.Download(context.Background(), req)
	if err != nil {
		details := make(map[string]string)
		details["enpoint"] = endpoint
		details["shortUrl"] = shortUrl

		e := kerr.NewErr(kerr.ErrLvlError, ErrorDomain, DownloadErrorCode, err, details)

		return e
	}

	stream, err := res.Recv()
	if err != nil {
		details := make(map[string]string)
		details["enpoint"] = endpoint
		details["shortUrl"] = shortUrl

		e := kerr.NewErr(kerr.ErrLvlError, ErrorDomain, DownloadErrorCode, err, details)

		return e
	}

	_, err = client.fs.Save(stream.Data, int64(len(stream.Data)), shortUrl)
	if err != nil {
		return err
	}

	return nil
}
