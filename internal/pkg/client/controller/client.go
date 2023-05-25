package clientcontroller

import (
	"context"

	clientusecase "github.com/sudak-91/monitoring/internal/pkg/client/use_case"
	msusecase "github.com/sudak-91/monitoring/internal/pkg/messageservice/use_case"
	wsservice "github.com/sudak-91/monitoring/internal/pkg/ws/service"
)

type ClientController struct {
	ms      msusecase.ServerMessager
	usecase *clientusecase.UserUseCase
	ws      wsservice.Servicer
}

func NewClientController(usecase *clientusecase.UserUseCase, ms msusecase.ServerMessager, ws wsservice.Servicer) *ClientController {
	return &ClientController{
		ms:      ms,
		usecase: usecase,
		ws:      ws}
}

func (c *ClientController) InitUUID() error {
	return c.usecase.InitUUID()
}

func (c *ClientController) GetUUIDFromServer(ctx context.Context) error {
	data, err := c.ms.GetUUIDFromServer(ctx)
	if err != nil {
		return err
	}
	err = c.ws.Request(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientController) SetUUID(ctx context.Context) error {
	data, err := c.ms.GetUUIDFromServer(ctx)
	if err != nil {
		return err
	}
	return c.ws.Request(ctx, data)
}
