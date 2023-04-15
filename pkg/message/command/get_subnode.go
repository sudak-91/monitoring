package message

type GetSubNode struct {
	ParentID            string
	ParentNodeID        uint32
	ParentNodeNamespace uint16
}

func GetSubNodeCommande(ParentID string, NodeId uint32, NodeNamespase uint16) Command {
	var Command Command
	Command.GetSubNode = &GetSubNode{ParentNodeID: NodeId, ParentNodeNamespace: NodeNamespase, ParentID: ParentID}
	return Command
}
