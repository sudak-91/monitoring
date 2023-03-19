package messageservice

import (
	"context"
	"log"

	message "github.com/sudak-91/monitoring/pkg/message/command"
	"nhooyr.io/websocket"
)

type MessageService struct {
	connection     *websocket.Conn
	ctx            context.Context
	commandService *CommandService
	commandMessage chan message.Command
}

func NewMessageService(ctx context.Context, ws *websocket.Conn) *MessageService {
	var m MessageService
	m.ctx = ctx
	m.connection = ws
	m.commandMessage = make(chan message.Command, 5)
	m.commandService = NewCommandService(ws)
	return &m
}

func (m *MessageService) Run() {
	for {
		select {
		case data, ok := <-m.commandMessage:
			if !ok {
				log.Println("Server errror read command")
				continue
			}
			go m.commandService.Service(data)
		case <-m.ctx.Done():
			log.Println("Message Service is done")
			break
		}
	}
}
