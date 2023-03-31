package clients

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

type Client struct {
	Conn       *websocket.Conn
	UUID       uuid.UUID
	IsUUIDTemp bool
}

type ClientList struct {
	Mutex sync.RWMutex
	Users map[uuid.UUID]*Client
}

func NewClientList() *ClientList {
	var Server ClientList
	Server.Users = make(map[uuid.UUID]*Client)
	return &Server
}

func (cl *ClientList) ChangeUUID(newUUID uuid.UUID, oldUUID uuid.UUID) error {
	_, ok := cl.Users[oldUUID]
	if ok {
		cl.Users[newUUID] = cl.Users[oldUUID]
		delete(cl.Users, oldUUID)
		return nil
	}
	return fmt.Errorf("Map hasn't key: %s", oldUUID.String())

}
