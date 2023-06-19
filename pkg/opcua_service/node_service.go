package opcuaservice

import (
	"fmt"
	"log"

	"github.com/gopcua/opcua"
	"github.com/spf13/viper"
	update "github.com/sudak-91/monitoring/pkg/message/update"
)

func (opc *OPCUAService) getNodes(node *opcua.Node, level int) (organizeNodesResult []update.NodeDef, componentNodesResult []update.NodeDef, propertyNodesResult []update.NodeDef, Err error) {
	log.Printf("OPCUA Level: %d", viper.GetInt("MaxNodeLevel"))
	if level > 10 {
		return
	}
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
		orgNodesRes, compNodesRes, propNodesRes, err := opc.getNodes(v, level+1)
		if err != nil {
			Err = err
			continue
		}
		node.ChildNode.OrganizesNode = orgNodesRes
		node.ChildNode.ComponentNode = compNodesRes
		node.ChildNode.PropertyNode = propNodesRes
		organizeNodesResult = append(organizeNodesResult, node)

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
		orgNodesRes, compNodesRes, propNodesRes, err := opc.getNodes(v, level+1)
		if err != nil {
			Err = err
			continue
		}
		node.ChildNode.OrganizesNode = orgNodesRes
		node.ChildNode.ComponentNode = compNodesRes
		node.ChildNode.PropertyNode = propNodesRes
		componentNodesResult = append(componentNodesResult, node)
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
		orgNodesRes, compNodesRes, propNodesRes, err := opc.getNodes(v, level+1)
		if err != nil {
			Err = err
			continue
		}
		node.ChildNode.OrganizesNode = orgNodesRes
		node.ChildNode.ComponentNode = compNodesRes
		node.ChildNode.PropertyNode = propNodesRes
		propertyNodesResult = append(propertyNodesResult, node)

	}
	return
}
