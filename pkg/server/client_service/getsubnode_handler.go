package clientservice

import (
	"github.com/sudak-91/monitoring/pkg/message"
	update "github.com/sudak-91/monitoring/pkg/message/update"
	"nhooyr.io/websocket"
)

func (cs *ClientService) GetSubNodeHandle(parrentId string, id uint32, sid string, nodeNamespace uint16) error {
	subNodes := update.NewOPCSubNodeUpdate(parrentId)
	OPCNodes, err := cs.opcuaService.GetNodes(nodeNamespace, id, sid)
	if err != nil {
		return err
	}
	subNodes.Nodes = OPCNodes
	update := subNodes.GetUpdate()
	data, err := message.EncodeData(update)
	if err != nil {
		return err
	}
	err = cs.client.Conn.Write(cs.ctx, websocket.MessageBinary, data)
	if err != nil {
		return err
	}
	return nil
}
