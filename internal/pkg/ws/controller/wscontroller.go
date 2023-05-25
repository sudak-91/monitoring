package wscontroller

import (
	"context"
	"log"

	msusecase "github.com/sudak-91/monitoring/internal/pkg/messageservice/use_case"
	wsusecase "github.com/sudak-91/monitoring/internal/pkg/ws/use_case"
)

type WSController struct {
	ws wsusecase.SocketUseCase
	ms msusecase.MessagerUseCase
}

func NewWSController(ws wsusecase.SocketUseCase, ms msusecase.MessagerUseCase) *WSController {
	return &WSController{
		ws: ws,
		ms: ms}
}

func (u *WSController) Run() error {
	//context may be
	for {
		_, data, err := u.ws.Response(context.TODO())
		if err != nil {
			return err
		}
		log.Println("[WSUseCase|Run] Responce new message")
		go func() {
			err = u.ms.Update(data)
			if err != nil {
				log.Println(err.Error())
			}
		}()
	}
}
