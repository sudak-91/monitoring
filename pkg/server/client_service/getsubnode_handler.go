package clientservice

import (
	"github.com/gopcua/opcua/ua"
	"github.com/sudak-91/monitoring/pkg/message"
	update "github.com/sudak-91/monitoring/pkg/message/update"
	opcuaservice "github.com/sudak-91/monitoring/pkg/server/opcua_service"
	"nhooyr.io/websocket"
)

func (cs *ClientService) GetSubNodeHandle(parrentId string, nodeid uint32, nodeNamespace uint16) error {
	subNodes := update.NewOPCSubNodeUpdate(parrentId)
	NodeID := ua.NewNumericNodeID(nodeNamespace, nodeid)
	Node := cs.opcuaService.OPCLient.Node(NodeID)
	OrganizesNodeList, err := cs.opcuaService.GetOrganizesNodes(Node)
	if err != nil {
		return err
	}
	for _, v := range OrganizesNodeList {
		node, err := opcuaservice.CreateNode(v)
		if err != nil {
			return err
		}
		subNodes.AddOrganizeNode(node)
	}
	HasComponentNodeList, err := cs.opcuaService.GetHasComponentNodes(Node)
	if err != nil {
		return err
	}
	for _, v := range HasComponentNodeList {
		node, err := opcuaservice.CreateNode(v)
		if err != nil {
			return err
		}
		subNodes.AddComponentNode(node)
	}
	HasPropertyNodeList, err := cs.opcuaService.GetHasPropertyNodes(Node)
	if err != nil {
		return err
	}
	for _, v := range HasPropertyNodeList {
		node, err := opcuaservice.CreateNode(v)
		if err != nil {
			return err
		}
		subNodes.AddPropertyNode(node)
	}
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
