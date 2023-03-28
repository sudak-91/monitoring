package updateservice

import (
	"context"
	"log"

	"github.com/google/uuid"
	update "github.com/sudak-91/monitoring/pkg/message/update"
	"github.com/sudak-91/monitoring/pkg/server"
)

type ChangeUUID struct {
	OldID uuid.UUID
	NewID uuid.UUID
}

type GetOpcUaNode struct {
	Info string
}

type UpdateService struct {
	ctx               context.Context
	clientChan        <-chan any
	opcuaChan         <-chan any
	updateToClienChan chan<- any
	updateToOpcUaChan chan<- any
	server            *server.Server
}

func NewUpdateService(ctx context.Context, clientChan <-chan any, opcuaChan <-chan any, updateToClientChan chan<- any, updateToOpcUaChan chan<- any, server *server.Server) *UpdateService {
	var u UpdateService
	u.ctx = ctx
	u.clientChan = clientChan
	u.opcuaChan = opcuaChan
	u.updateToClienChan = updateToClientChan
	u.updateToOpcUaChan = updateToOpcUaChan
	u.server = server
	return &u
}

func (s *UpdateService) Update() {
	log.Println("Update service start")
	for {
		select {
		case data := <-s.clientChan:
			go s.clientRouter(data)
		case data := <-s.opcuaChan:
			log.Println("Get data from opcUA Chan")
			go s.opcuaRouter(data)
		case <-s.ctx.Done():
			log.Println("Connection done")
			return

		}
	}
}

func (s *UpdateService) clientRouter(data any) {
	switch v := data.(type) {
	case ChangeUUID:
		err := s.changeUUID(v.OldID, v.NewID)
		if err != nil {
			log.Println(err)
		}
	case GetOpcUaNode:
		s.updateToOpcUaChan <- v
	}
}

func (s *UpdateService) opcuaRouter(data any) {
	switch v := data.(type) {
	case update.SendOpcNodes:
		log.Println("v")
		s.updateToClienChan <- v
		log.Println("sendOPCUa Node")
	}
}

func (s *UpdateService) changeUUID(oldUUID uuid.UUID, newUUID uuid.UUID) error {
	s.server.Users[newUUID] = s.server.Users[oldUUID]
	delete(s.server.Users, oldUUID)
	return nil
}
