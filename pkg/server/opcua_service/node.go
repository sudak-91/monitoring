package opcuaservice

import "github.com/gopcua/opcua"

type OPCNode struct {
	Name      string
	ID        uint32
	Namespace uint16
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
