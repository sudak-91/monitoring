package server

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

type Server struct {
	Mutex sync.RWMutex
	Users map[uuid.UUID]*Client
	ctx   context.Context
}

type Client struct {
	Conn       *websocket.Conn
	UUID       uuid.UUID
	IsUUIDTemp bool
}

func NewServer(ctx context.Context) *Server {
	var Server Server
	Server.Users = make(map[uuid.UUID]*Client)
	Server.ctx = ctx
	return &Server
}
