package clientservice

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/sudak-91/monitoring/pkg/client"
	"nhooyr.io/websocket"
)

// Service for connected users

type ClientService struct {
	Mutex                        sync.RWMutex
	clientToServiceChan          chan interface{} //message from client
	commandToOpcUaControllerChan chan any
	//updateChan          <-chan any
	Users map[uuid.UUID]*client.Client
	ctx   context.Context
}

func (cs *ClientService) AddClient(client *client.Client) {
	cs.Mutex.Lock()
	defer cs.Mutex.Unlock()
	cs.Users[client.UUID] = client
}

func (cs *ClientService) ChangeUUID(newUUID uuid.UUID, oldUUID uuid.UUID) error {
	cs.Mutex.Lock()
	defer cs.Mutex.Unlock()
	_, ok := cs.Users[oldUUID]
	if ok {
		cs.Users[newUUID] = cs.Users[oldUUID]
		delete(cs.Users, oldUUID)
		return nil
	}
	return fmt.Errorf("Map hasn't key: %s", oldUUID.String())

}

func NewClientService(ctx context.Context, commandToOpcUaControllerChan chan any) *ClientService {
	var cs ClientService
	cs.ctx = ctx
	cs.Users = make(map[uuid.UUID]*client.Client)
	cs.clientToServiceChan = make(chan interface{}, 5)
	//cs.updateChan = updateChan
	cs.commandToOpcUaControllerChan = commandToOpcUaControllerChan

	return &cs

}

func (cs *ClientService) NewClient(ws *websocket.Conn) {
	client := client.NewClient(ws, cs.clientToServiceChan, cs.commandToOpcUaControllerChan, cs.ctx)
	cs.AddClient(client)
	go client.Run()
}

func (cs *ClientService) Run() {
	for {
		select {
		case data := <-cs.clientToServiceChan:
			cs.router(data)
		default:

		}
	}

}

func (cs *ClientService) router(data any) {
	switch v := data.(type) {
	case client.ChangeUUID:
		err := cs.setUUIDHandle(v.NewUUID, v.OldUUID)
		if err != nil {
			log.Println(err.Error())
			return
		}
	default:
		log.Println("Default")
	}
}

func (cs *ClientService) setUUIDHandle(newUUID uuid.UUID, oldUUID uuid.UUID) error {
	if err := cs.ChangeUUID(newUUID, oldUUID); err != nil {
		return err
	}
	return nil
}
