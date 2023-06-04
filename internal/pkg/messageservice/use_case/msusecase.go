package msusecase

import (
	"context"

	cliententity "github.com/sudak-91/monitoring/internal/pkg/client/entity"
	"github.com/sudak-91/monitoring/pkg/message"
	command "github.com/sudak-91/monitoring/pkg/message/command"
	update "github.com/sudak-91/monitoring/pkg/message/update"
)

type MessagerUseCase interface {
	Update([]byte) error
}

type MessagePresenter interface {
	SendCommand(context.Context, []byte) error
}

type ServerMessager interface {
	GetUUIDFromServer(context.Context) ([]byte, error)
	SetUUID(context.Context) ([]byte, error)
}

type MSUseCase struct {
	ctx            context.Context
	client         cliententity.Clienter
	messageChannel chan any
}

func NewMSUseCase(ctx context.Context, client cliententity.Clienter, messageChannel chan any) *MSUseCase {
	return &MSUseCase{
		ctx:            ctx,
		client:         client,
		messageChannel: messageChannel,
	}
}

func (s *MSUseCase) Update(data []byte) error {
	update, err := message.Decode[update.Update](data)
	if err != nil {
		return err
	}
	s.update(update)
	return nil
}

func (s *MSUseCase) update(update update.Update) {
	switch {
	case update.SendUUID != nil:
		s.SendUUID(*update.SendUUID, s.client)
	case update.OpcNodes != nil:
		s.OpcNodesUpdate(update.OpcNodes)
	case update.OPCSubNode != nil:
		s.OpcSubNodesupdate(update.OPCSubNode)
	case update.NodeDescription != nil:
		s.OpcNodeDescription(update.NodeDescription)

	}
}

func (s *MSUseCase) GetUUIDFromServer(ctx context.Context) ([]byte, error) {
	cmd := command.Command{}
	var a command.GetUUID
	cmd.GetUUID = &a
	data, err := message.EncodeData(cmd)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *MSUseCase) SetUUID(context.Context) ([]byte, error) {
	cmd := command.Command{}
	var a command.SetUUID
	sUUID, err := s.client.GetUUID()
	if err != nil {
		return nil, err
	}
	UUID := sUUID.String()
	a.UUID = UUID
	cmd.SetUUID = &a
	data, err := message.EncodeData(cmd)
	if err != nil {
		return nil, err
	}
	return data, nil
}
