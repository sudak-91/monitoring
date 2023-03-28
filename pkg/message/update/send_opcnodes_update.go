package message

type SendOpcNodes struct {
	Nodes []string
}

func NewSendOpcNodes(nodes []string) SendOpcNodes {
	var s SendOpcNodes
	s.Nodes = nodes
	return s
}
