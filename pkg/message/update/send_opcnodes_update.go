package message

type SendOpcNodes struct {
	Nodes *OPCNode
}

func NewSendOpcNodes(nodes *OPCNode) *SendOpcNodes {
	var s SendOpcNodes
	s.Nodes = nodes
	return &s
}
