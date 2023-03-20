package client

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"

	"github.com/google/uuid"
	clientService "github.com/sudak-91/monitoring/internal/pkg/client/messageservice"
	message "github.com/sudak-91/monitoring/pkg/message/update"
	"github.com/sudak-91/wasmhtml/cookie"
	"nhooyr.io/websocket"
)

type Client struct {
	ws             *websocket.Conn
	messageService *clientService.MessageService
	ctx            context.Context
	Cookie         *cookie.Cookie
	UUID           uuid.UUID
}

func NewClient(ctx context.Context, cookie *cookie.Cookie) *Client {
	var c Client
	c.ctx = ctx
	c.Cookie = cookie
	return &c
}

func (c *Client) Requst(data []byte) error {
	return c.ws.Write(c.ctx, websocket.MessageBinary, data)
}

func (c *Client) Run(done chan bool) {
	ctx, cancel := context.WithCancel(c.ctx)
	conn, _, err := websocket.Dial(ctx, "ws://localhost:8000/ws", nil)
	if err != nil {
		cancel()
		log.Println(err)
		return
	}
	c.ws = conn
	done <- true
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
			go c.messageService.Update(decodeData, &c.UUID, c.Cookie)
		case <-ctx.Done():

		}
	}
}
