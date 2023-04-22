package opcuaservice

import (
	"context"
	"fmt"
	"log"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/id"
	"github.com/gopcua/opcua/ua"
)

type OPCNode struct {
	Name      string
	ID        uint32
	Namespace uint16
}

func (opc *OPCUAService) GetRootNodes() (OPCUAObjectData, error) {
	uid, err := ua.ParseNodeID("ns=0;i=84")
	if err != nil {
		panic(err)
	}
	node := opc.OPCLient.Node(uid)
	nodesList, err := node.ReferencedNodesWithContext(context.Background(), id.Organizes, ua.BrowseDirectionForward, ua.NodeClassAll, true)
	if err != nil {
		panic(err)
	}
	var Nodes OPCUAObjectData
	for _, v := range nodesList {
		var node OPCNode
		log.Println(v.ID.Namespace())
		log.Println(v.ID.IntID())
		node.ID = v.ID.IntID()
		node.Namespace = v.ID.Namespace()
		name, err := v.BrowseName()
		if err != nil {
			fmt.Println(err.Error())
			return OPCUAObjectData{}, err
		}
		node.Name = name.Name
		fmt.Println(name.Name)
		fmt.Printf("Roots node list is %v\n", v.ID)
		Nodes.Nodes = append(Nodes.Nodes, node)

	}
	return Nodes, nil
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
func CreateNode(v *opcua.Node) (OPCNode, error) {
	var node OPCNode
	node.ID = v.ID.IntID()
	node.Namespace = v.ID.Namespace()
	nodeName, err := v.BrowseName()
	if err != nil {
		return OPCNode{}, err
	}
	node.Name = nodeName.Name
	return node, nil
}
