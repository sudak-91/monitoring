package messageservice

import (
	"context"

	"github.com/google/uuid"
	message "github.com/sudak-91/monitoring/pkg/message/update"
	"github.com/sudak-91/wasmhtml/cookie"
)

type MessageService struct {
}

func NewMessageService(ctx context.Context, update chan message.Update) *MessageService {
	var s MessageService
	return &s
}

func (s *MessageService) Update(data message.Update, uuid *uuid.UUID, cookie *cookie.Cookie) {
	switch {
	case data.SendUUID != nil:
		s.SendUUIDService(*data.SendUUID, uuid, cookie)
		return
	}

}
