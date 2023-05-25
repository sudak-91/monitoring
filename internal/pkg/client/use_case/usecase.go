package clientusecase

import (
	"github.com/google/uuid"
	cliententity "github.com/sudak-91/monitoring/internal/pkg/client/entity"

	"github.com/sudak-91/wasmhtml/cookie"
)

type UserUseCase struct {
	client cliententity.Clienter
}

func NewUserUseCase(client cliententity.Clienter) *UserUseCase {
	return &UserUseCase{client: client}
}

func (i *UserUseCase) InitUUID() error {
	cookie := cookie.NewCookie()
	result, err := cookie.GetValue("UUID")
	if err != nil {
		return err
	}
	uuid, err := uuid.Parse(result)
	if err != nil {
		return err
	}
	i.client.SetUUID(&uuid)
	return nil
}
