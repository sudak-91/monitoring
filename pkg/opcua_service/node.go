package opcuaservice

import (
	"fmt"
	"log"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/id"
	"github.com/gopcua/opcua/ua"
	"github.com/spf13/viper"
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
		NodeType, err := opc.GetNodeDataType(v)
		if err == nil {
			log.Printf("AddNodeType: %d", NodeType)
			node.NodeType = NodeType
		}
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
		NodeType, err := opc.GetNodeDataType(v)
		if err == nil {
			log.Printf("AddNodeType: %d", NodeType)
			node.NodeType = NodeType
		}
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
		NodeType, err := opc.GetNodeDataType(v)
		if err == nil {
			node.NodeType = NodeType
		}
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
	attrs, err := node.Attributes(ua.AttributeIDDataType, ua.AttributeIDValue, ua.AttributeIDNodeClass)

	if err != nil {
		return nil, err
	}
	return attrs, nil
}

func (opc *OPCUAService) GetNodeDataType(node *opcua.Node) (uint32, error) {
	att, err := node.Attributes(ua.AttributeIDDataType)
	if err != nil {
		log.Printf("[GetNode|Error]:%s\n", err.Error())
		return 0, err

	}
	if viper.GetBool("Debug") {
		log.Println(len(att))
	}
	if len(att) == 0 {
		return 0, nil
	}

	if att[0].Status == ua.StatusOK {
		return att[0].Value.NodeID().IntID(), nil
	}

	return 0, nil

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

func (opc *OPCUAService) GetNodeBySID(ns uint16, sid string) *opcua.Node {
	nodeID := ua.NewStringNodeID(ns, sid)
	node := opc.OPCLient.Node(nodeID)
	return node
}
