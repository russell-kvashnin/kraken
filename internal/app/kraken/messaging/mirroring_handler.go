package messaging

import (
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/node"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/rpc/client"
	"github.com/streadway/amqp"
)

// File mirroring message
type MirroringMessage struct {
	NodeId   string
	ShortUrl string
}

// Mirroring message handler
type MirroringMessageHandler struct {
	repo   node.Repository
	client *client.MirroringClient
}

// Mirroring message handler constructor
func NewMirroringMessageHandler(repo node.Repository, client *client.MirroringClient) *MirroringMessageHandler {
	handler := new(MirroringMessageHandler)
	handler.repo = repo
	handler.client = client

	return handler
}

// Handle file mirroring message
func (handler *MirroringMessageHandler) Handle(d amqp.Delivery) error {
	message := &MirroringMessage{}
	err := json.Unmarshal(d.Body, message)
	if err != nil {
		return err
	}

	n, err := handler.repo.Get(message.NodeId)
	if err != nil {
		return err
	}

	err = handler.client.Download(n.Mirroring.Endpoint, message.ShortUrl)
	if err != nil {
		return err
	}

	return nil
}
