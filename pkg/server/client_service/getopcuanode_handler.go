package clientservice

import (
	"github.com/sudak-91/monitoring/pkg/message"
	update "github.com/sudak-91/monitoring/pkg/message/update"
	"nhooyr.io/websocket"
)

func (cs *ClientService) getOpcUaNodeHandle() error {
	var upd update.Update
	data, err := cs.opcuaService.GetNodes(0, 84, "")
	if err != nil {
		return err
	}
	upd.SendOpcNodes = update.NewSendOpcNodes(&data)
	breq, err := message.EncodeData(upd)
	if err != nil {
		return err
	}
	if err = cs.client.Conn.Write(cs.ctx, websocket.MessageBinary, breq); err != nil {
		return err
	}
	return nil
}
