package clientservice

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"

	message "github.com/sudak-91/monitoring/pkg/message/command"
	"github.com/sudak-91/monitoring/pkg/server"
)

type ClientService struct {
	serverUpdate chan interface{} //chan for communicate with server
	client       *server.Client
	ctx          context.Context
}

func NewClientService(ctx context.Context, client *server.Client, serverUpdate chan any) *ClientService {
	var cs ClientService
	cs.ctx = ctx
	cs.client = client
	cs.serverUpdate = serverUpdate
	return &cs

}

func (cs *ClientService) Run() {
	ctx, cancel := context.WithCancel(cs.ctx)
mailoop:
	for {
		select {
		default:
			MessageType, data, err := cs.client.Conn.Read(ctx)
			if err != nil {
				log.Println(err)
				cancel()
				continue
			}
			log.Printf("%s get message %v\n", cs.client.UUID, MessageType)
			log.Printf("%v", data)

			rd := bytes.NewReader(data)
			dec := gob.NewDecoder(rd)
			var msg message.Command
			err = dec.Decode(&msg)
			if err != nil {
				log.Println(err)
				continue
			}
			go cs.Update(msg)

		case <-ctx.Done():
			break mailoop
		}
	}
	log.Println("End")
}
