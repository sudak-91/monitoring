package messageservice

import (
	"log"

	"github.com/google/uuid"
	message "github.com/sudak-91/monitoring/pkg/message/update"
	"github.com/sudak-91/wasmhtml/cookie"
)

func (m *MessageService) SendUUIDService(data message.SendUUID, clientUuid *uuid.UUID, cookie *cookie.Cookie) {
	ClientUUID := data.UUID
	bUUID, err := uuid.Parse(ClientUUID)
	if err != nil {
		log.Println(err)
	}
	*clientUuid = bUUID
	cookie.AddCoookie("UUID", ClientUUID)
}
