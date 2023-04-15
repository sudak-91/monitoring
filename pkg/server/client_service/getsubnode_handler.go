package clientservice

import (
	"github.com/gopcua/opcua/id"
	"github.com/gopcua/opcua/ua"
	"github.com/sudak-91/monitoring/pkg/message"
	update "github.com/sudak-91/monitoring/pkg/message/update"
	opcuaservice "github.com/sudak-91/monitoring/pkg/server/opcua_service"
	"nhooyr.io/websocket"
)

func (cs *ClientService) GetSubNodeHandle(parrentId string, nodeid uint32, nodeNamespace uint16) error {
	subNodes := update.NewOPCSubNodeUpdate(parrentId)
	Node := ua.NewNumericNodeID(nodeNamespace, nodeid)
	parenNode := cs.opcuaService.OPCLient.Node(Node)
	OrganizesNodeList, err := parenNode.ReferencedNodes(id.Organizes, ua.BrowseDirectionForward, ua.NodeClassAll, true)
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
	HasComponentNodeList, err := parenNode.ReferencedNodes(id.HasComponent, ua.BrowseDirectionForward, ua.NodeClassAll, true)
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
	HasPropertyNodeList, err := parenNode.ReferencedNodes(id.HasProperty, ua.BrowseDirectionForward, ua.NodeClassAll, true)
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
