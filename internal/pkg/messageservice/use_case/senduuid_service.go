package msusecase

import (
	"log"

	"github.com/google/uuid"
	cliententity "github.com/sudak-91/monitoring/internal/pkg/client/entity"
	update "github.com/sudak-91/monitoring/pkg/message/update"
	"github.com/sudak-91/wasmhtml/cookie"
)

func (m *MSUseCase) SendUUID(data update.SendUUID, client cliententity.Clienter) {
	cookie := cookie.NewCookie()
	ClientUUID := data.UUID
	bUUID, err := uuid.Parse(ClientUUID)
	if err != nil {
		log.Println(err)
	}
	client.SetUUID(&bUUID)
	cookie.AddCoookie("UUID", ClientUUID)
}
