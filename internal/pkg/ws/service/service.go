package wsservice

import (
	"context"
	"log"

	"nhooyr.io/websocket"
)

type Requester interface {
	Request(context.Context, []byte) error
}
type Responcer interface {
	Response(context.Context) (websocket.MessageType, []byte, error)
}
type Servicer interface {
	Requester
	Responcer
}
type WSService struct {
	ws *websocket.Conn
}

func NewWSService(ctx context.Context) (*WSService, error) {
	log.Println("[WSService]Create WSService")
	conn, _, err := websocket.Dial(ctx, "ws://localhost:8000/ws", nil)
	if err != nil {
		return nil, err
	}
	var service WSService
	service.ws = conn
	return &service, nil
}

func (s *WSService) Response(ctx context.Context) (websocket.MessageType, []byte, error) {
	return s.ws.Read(ctx)
}

func (s *WSService) Request(ctx context.Context, data []byte) error {
	return s.ws.Write(ctx, websocket.MessageBinary, data)
}
