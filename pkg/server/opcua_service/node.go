package opcuaservice

import (
	"fmt"
	"log"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/id"
	"github.com/gopcua/opcua/ua"
	update "github.com/sudak-91/monitoring/pkg/message/update"
)

var (
	consoleColor = "\033[1;34m%s\033[0m"
)

func (opc *OPCUAService) GetNodes(ns uint16, id uint32, sid string) (update.OPCNode, error) {
	var (
		nodeID   *ua.NodeID
		NodeList update.OPCNode
	)
	if id != 0 {
		nodeID = ua.NewNumericNodeID(ns, id)
	} else {
		nodeID = ua.NewStringNodeID(ns, sid)
	}
	node := opc.OPCLient.Node(nodeID)
	organizesNodes, err := opc.GetOrganizesNodes(node)
	if err != nil {
		log.Println(err.Error())
	}
	for _, v := range organizesNodes {
		node := CreateNode(v)
		fmt.Printf(consoleColor, "[OrganizesNode]")
		log.Printf("Namespace = %d\t ID = %d\t ID(String) = %s\t Name = %s\n", node.Namespace, node.IID, node.SID, node.Name)
		NodeList.AddOrganizeNode(node)

	}

	componentNodes, err := opc.GetHasComponentNodes(node)
	if err != nil {
		log.Println(err.Error())

	}
	for _, v := range componentNodes {
		node := CreateNode(v)
		fmt.Printf(consoleColor, "[ComponentNode]")
		log.Printf("Namespace = %d\t ID = %d\t ID(String) = %s\t Name = %s\n", node.Namespace, node.IID, node.SID, node.Name)
		NodeList.AddComponentNode(node)
	}
	propertyNodes, err := opc.GetHasPropertyNodes(node)
	if err != nil {
		log.Println(err.Error())
	}
	for _, v := range propertyNodes {
		node := CreateNode(v)
		fmt.Printf(consoleColor, "[PropertyNode]")
		log.Printf("Namespace = %d\t ID = %d\t ID(String) = %s\t Name = %s\n", node.Namespace, node.IID, node.SID, node.Name)
		NodeList.AddPropertyNode(node)
	}
	return NodeList, nil
}

func (opc *OPCUAService) GetOrganizesNodes(node *opcua.Node) ([]*opcua.Node, error) {
	OrganizesNodeList, err := node.ReferencedNodes(id.Organizes, ua.BrowseDirectionForward, ua.NodeClassAll, true)
	if err != nil {
		return nil, err
	}
	return OrganizesNodeList, nil
}

func (opc *OPCUAService) GetHasComponentNodes(node *opcua.Node) ([]*opcua.Node, error) {
	HasComponentNodeList, err := node.ReferencedNodes(id.HasComponent, ua.BrowseDirectionForward, ua.NodeClassAll, true)
	if err != nil {
		return nil, err
	}
	return HasComponentNodeList, nil
}

func (opc *OPCUAService) GetHasPropertyNodes(node *opcua.Node) ([]*opcua.Node, error) {
	HasPropertyNodeList, err := node.ReferencedNodes(id.HasProperty, ua.BrowseDirectionForward, ua.NodeClassAll, true)
	if err != nil {
		return nil, err
	}
	return HasPropertyNodeList, nil
}

func (opc *OPCUAService) GetDataValuesNode(node *opcua.Node) ([]*ua.DataValue, error) {
	attrs, err := node.Attributes(ua.AttributeIDDataType, ua.AttributeIDValue)
	if err != nil {
		return nil, err
	}
	return attrs, nil
}
func CreateNode(v *opcua.Node) update.NodeDef {
	var node update.NodeDef
	node.IID = v.ID.IntID()
	node.SID = v.ID.StringID()
	node.Namespace = v.ID.Namespace()
	nodeName, err := v.BrowseName()
	if err == nil {
		node.Name = nodeName.Name
	}
	return node
}
