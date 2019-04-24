package bus

import (
	"github.com/russell-kvashnin/kraken/internal/app/kraken/messaging"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/file"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/node"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
)

// file.Uploaded event handler
type FileUploadedHandler struct {
	repo       node.Repository
	producer   *messaging.MessageProducer
	errHandler *kerr.Handler
}

// File uploaded handler constructor
func NewFileUploadedHandler(repo node.Repository, producer *messaging.MessageProducer, errHandler *kerr.Handler) file.UploadedEventHandler {
	handler := new(FileUploadedHandler)
	handler.repo = repo
	handler.producer = producer
	handler.errHandler = errHandler

	return handler
}

// Handle file upload event
func (handler *FileUploadedHandler) Handle(event file.UploadedEvent) {
	mNode, err := handler.repo.GetMirroringNode()
	if err != nil || mNode.Id == "" {
		return
	}

	msg := messaging.MirroringMessage{
		NodeId:   event.NodeId,
		ShortUrl: event.ShortUrl,
	}

	err = handler.producer.Produce(mNode.Id, msg)
	if err != nil {
		handler.errHandler.Handle(err)
	}
}
