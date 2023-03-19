package client

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"

	"github.com/google/uuid"
	message "github.com/sudak-91/monitoring/pkg/message/update"
	"nhooyr.io/websocket"
)

type Client struct {
	ws            *websocket.Conn
	updateMessage chan message.Update
	ctx           context.Context
	UUID          uuid.UUID
}

func NewClient(ctx context.Context) *Client {
	var c Client

	c.ctx = ctx
	c.updateMessage = make(chan message.Update, 5)
	return &c
}

func (c *Client) Run() {
	ctx, cancel := context.WithCancel(c.ctx)
	conn, _, err := websocket.Dial(ctx, "ws://localhost:8000/ws", nil)
	if err != nil {
		cancel()
		log.Println(err)
		return
	}
	c.ws = conn

	for {
		messageType, data, err := c.ws.Read(ctx)
		if err != nil {
			log.Println(err)
			cancel()
			return
		}
		select {
		default:
			log.Println(messageType)
			var decodeData message.Update
			reader := bytes.NewReader(data)
			decoder := gob.NewDecoder(reader)
			err = decoder.Decode(&decodeData)
			if err != nil {
				log.Println(err)
				continue
			}
			c.updateMessage <- decodeData
		case <-ctx.Done():

		}
	}
}
