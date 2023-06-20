package message

type Updates interface {
	SendUUID()
	SendOpcNodes()
}

// OPCNodeType
type NodeDef struct {
	ChildNode OPCNode
	Name      string //OPC Node browser name
	IID       uint32 //OPC ID ex:"i=%d"
	SID       string //OPC ID ex:"s=%s"
	Namespace uint16 //OPC Namespace ex:"n=%s"
	NodeType  uint32
	DataType  string
}
type SubNodes struct {
	Parent string
	Nodes  OPCNode
}
type OPCNode struct {
	OrganizesNode []NodeDef
	ComponentNode []NodeDef
	PropertyNode  []NodeDef
}

type Update struct {
	SendUUID        *SendUUID
	OpcNodes        *OPCNodes
	OPCSubNode      *SubNodes
	NodeDescription *NodeDescription
}

func (n *OPCNode) AddOrganizeNode(node NodeDef) {
	n.OrganizesNode = append(n.OrganizesNode, node)
}

func (n *OPCNode) AddComponentNode(node NodeDef) {
	n.ComponentNode = append(n.ComponentNode, node)
}

func (n *OPCNode) AddPropertyNode(node NodeDef) {
	n.PropertyNode = append(n.PropertyNode, node)
}
