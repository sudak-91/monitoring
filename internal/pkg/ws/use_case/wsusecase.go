package wsusecase

import (
	"context"

	wsservice "github.com/sudak-91/monitoring/internal/pkg/ws/service"
	"nhooyr.io/websocket"
)

type SocketUseCase interface {
	Request(context.Context, []byte) error
	Response(context.Context) (websocket.MessageType, []byte, error)
}

type WSUseCase struct {
	ctx     context.Context
	service wsservice.Servicer
}

func NewWSUseCase(ctx context.Context, service wsservice.Servicer) *WSUseCase {
	return &WSUseCase{service: service, ctx: ctx}
}

func (u *WSUseCase) Response(ctx context.Context) (websocket.MessageType, []byte, error) {
	return u.service.Response(ctx)
}

func (u *WSUseCase) Request(ctx context.Context, data []byte) error {
	return u.service.Request(ctx, data)
}
