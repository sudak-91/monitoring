package clientservice

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"

	message "github.com/sudak-91/monitoring/pkg/message/command"
	"github.com/sudak-91/monitoring/pkg/server"
)

// Service for connected users

type ClientService struct {
	clientChan chan<- interface{} //chan for communicate with server
	updateChan <-chan any
	client     *server.Client
	ctx        context.Context
}

func NewClientService(ctx context.Context, client *server.Client, clientChan chan<- any, updateChan <-chan any) *ClientService {
	var cs ClientService
	cs.ctx = ctx
	cs.client = client
	cs.clientChan = clientChan
	cs.updateChan = updateChan
	return &cs

}

func (cs *ClientService) Run() {
	ctx, cancel := context.WithCancel(cs.ctx)
	go func() {
		log.Println("Start read message")
		for {
			MessageType, data, err := cs.client.Conn.Read(ctx)
			if err != nil {
				log.Println(err)
				cancel()
				break
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
			go cs.messageRouter(msg)
		}
	}()
	log.Println("Start Update Chan service")
mailoop:
	for {
		select {
		case data := <-cs.updateChan:
			log.Println("Get data from updateChan")
			go cs.updateRouter(data)
		case <-ctx.Done():
			log.Println("Connection is odne")
			break mailoop
		default:

		}
	}
	log.Println("End")
}
