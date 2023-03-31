package messageservice

import (
	"context"
	"log"

	"github.com/google/uuid"
	message "github.com/sudak-91/monitoring/pkg/message/update"
	"github.com/sudak-91/wasmhtml/cookie"
)

type MessageService struct {
	messageServiceChan chan interface{}
}

func NewMessageService(ctx context.Context, messageServiceChan chan interface{}) *MessageService {
	var s MessageService
	s.messageServiceChan = messageServiceChan
	return &s
}

func (s *MessageService) Update(data message.Update, uuid *uuid.UUID, cookie *cookie.Cookie) {
	switch {
	case data.SendUUID != nil:
		go s.SendUUIDService(*data.SendUUID, uuid, cookie)
		return
	case data.SendOpcNodes != nil:
		for _, v := range data.SendOpcNodes.Nodes.Nodes {
			log.Println(v.Name)
			log.Println(v.ID)
		}
		s.messageServiceChan <- data.SendOpcNodes
	}

}
