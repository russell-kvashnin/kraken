package file

import (
	"bytes"
	"github.com/asaskevich/EventBus"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/fs"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/messaging"
	"io"
	"mime/multipart"
)

// Upload handler request
type UploadFileCommand struct {
	File       multipart.File
	Header     *multipart.FileHeader
	Tags       []string
	UploadRoot string
	NodeId     string
}

// Update file command
type UpdateFileCommand struct {
	Path string
}

// Handler for upload command
type UploadFileHandler struct {
	fs       *fs.FileService
	repo     Repository
	producer *messaging.MessageProducer
	bus      EventBus.Bus
}

// Upload file handler dependencies
type UploadFileHandlerDeps struct {
	Fs       *fs.FileService
	Repo     Repository
	Producer *messaging.MessageProducer
	Bus      EventBus.Bus
}

// Upload handler constructor
func NewUploadFileHandler(d UploadFileHandlerDeps) *UploadFileHandler {
	handler := new(UploadFileHandler)
	handler.fs = d.Fs
	handler.repo = d.Repo
	handler.producer = d.Producer
	handler.bus = d.Bus

	return handler
}

// handle command
func (handler *UploadFileHandler) Handle(cmd UploadFileCommand) (File, error) {
	var (
		err error
	)

	model := File{}
	model.Upload(cmd)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, cmd.File); err != nil {
		return File{}, err
	}

	filePath, err := handler.fs.Save(buf.Bytes(), cmd.Header.Size, model.ShortUrl)
	if err != nil {
		return File{}, err
	}

	upCmd := UpdateFileCommand{
		Path: filePath,
	}
	model.Update(upCmd)

	err = handler.repo.Store(model)
	if err != nil {
		return File{}, err
	}

	ev := UploadedEvent{
		NodeId:   cmd.NodeId,
		ShortUrl: model.ShortUrl,
	}

	handler.bus.Publish(UploadedEventAlias, ev)

	return model, nil
}
