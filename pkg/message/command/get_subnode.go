package message

type GetSubNode struct {
	DOMParrentID string
	IID          uint32
	SID          string
	Namespace    uint16
}

func GetSubNodeCommande(ParentID string, ID uint32, NodeNamespase uint16, SID string) Command {
	var Command Command
	Command.GetSubNode = &GetSubNode{IID: ID, Namespace: NodeNamespase, DOMParrentID: ParentID, SID: SID}
	return Command
}
