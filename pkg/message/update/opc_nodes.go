package message

type OPCNodes struct {
	Nodes *OPCNode
}

func NewSendOpcNodes(nodes *OPCNode) *OPCNodes {
	var s OPCNodes
	s.Nodes = nodes
	return &s
}
