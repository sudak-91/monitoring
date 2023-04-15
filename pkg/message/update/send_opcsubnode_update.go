package message

import opcuaservice "github.com/sudak-91/monitoring/pkg/server/opcua_service"

type OPCSubNode struct {
	Parent        string
	OrganizesNode []opcuaservice.OPCNode
	ComponentNode []opcuaservice.OPCNode
	PropertyNode  []opcuaservice.OPCNode
}

func NewOPCSubNodeUpdate(parent string) OPCSubNode {

	return OPCSubNode{Parent: parent}
}

func (n *OPCSubNode) AddOrganizeNode(node opcuaservice.OPCNode) {
	n.OrganizesNode = append(n.OrganizesNode, node)
}

func (n *OPCSubNode) AddComponentNode(node opcuaservice.OPCNode) {
	n.ComponentNode = append(n.ComponentNode, node)
}

func (n *OPCSubNode) AddPropertyNode(node opcuaservice.OPCNode) {
	n.PropertyNode = append(n.PropertyNode, node)
}

func (n *OPCSubNode) GetUpdate() Update {
	var Update Update
	Update.OPCSubNode = n
	return Update
}
