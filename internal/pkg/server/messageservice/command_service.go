package messageservice

import (
	message "github.com/sudak-91/monitoring/pkg/message/command"
	"nhooyr.io/websocket"
)

type CommandService struct {
	connection *websocket.Conn
}

func NewCommandService(ws *websocket.Conn) *CommandService {
	var u CommandService
	u.connection = ws
	return &u
}

func (u *CommandService) Service(data message.Command) {
	switch {
	case data.GetInfo != nil:

	}
}
