package messageservice

import (
	"context"
	"log"

	message "github.com/sudak-91/monitoring/pkg/message/update"
)

type MessageService struct {
	ctx           context.Context
	updateMessage chan message.Update

	updateService updateService
}

func NewMessageService(ctx context.Context, update chan message.Update) *MessageService {
	var s MessageService
	s.ctx = ctx
	s.updateMessage = update
	s.updateService = updateService{}

	return &s
}

func (s *MessageService) Update() {
	for {
		select {
		case data := <-s.updateMessage:
			s.updateService.router(data)
		case <-s.ctx.Done():
			log.Println("Service is down")
			break
		}
	}

}
