package clientservice

import (
	"context"
	"log"

	"github.com/sudak-91/monitoring/pkg/message"
	update "github.com/sudak-91/monitoring/pkg/message/update"
	"nhooyr.io/websocket"
)

func (cs *ClientService) GetNodeDescriptionHandler(ns uint16, sid string) error {
	node := cs.opcuaService.GetNodeBySID(ns, sid)
	DataType, Description, err := cs.opcuaService.GetNodeDescription(node)
	if err != nil {
		log.Printf("[GetNodeDescription]|%s", err.Error())
		return err
	}
	update := update.NewNodeDescriptionUpdate(DataType, Description)
	data, err := message.EncodeData(update)
	if err != nil {
		return err
	}
	err = cs.client.Conn.Write(context.TODO(), websocket.MessageBinary, data)
	if err != nil {
		return err
	}
	return nil
}
