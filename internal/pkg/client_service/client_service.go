package clientservice

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"

	"github.com/google/uuid"

	message "github.com/sudak-91/monitoring/pkg/message/command"
	"nhooyr.io/websocket"
)

type ClientService struct {
	Conn         *websocket.Conn
	UUID         uuid.UUID
	serverUpdate chan interface{}
	IsUUIDTemp   bool
	ctx          context.Context
}

func NewClientService(ctx context.Context, uuid uuid.UUID, conn *websocket.Conn, serverUpdate chan interface{}) *ClientService {
	var client ClientService
	client.ctx = ctx
	client.UUID = uuid
	client.IsUUIDTemp = true
	client.Conn = conn
	client.serverUpdate = serverUpdate
	return &client

}

func (w *ClientService) Run() {
mailoop:
	for {
		select {
		default:
			MessageType, data, err := w.Conn.Read(context.TODO())
			if err != nil {
				log.Println(err)
			}
			log.Printf("%s get message %v\n", w.UUID, MessageType)
			log.Printf("%v", data)

			rd := bytes.NewReader(data)
			dec := gob.NewDecoder(rd)
			var msg message.Command
			err = dec.Decode(&msg)
			if err != nil {
				log.Println(err)
				continue
			}
			go w.Update(msg)

		case <-w.ctx.Done():
			break mailoop
		}
	}
	log.Println("End")
}
