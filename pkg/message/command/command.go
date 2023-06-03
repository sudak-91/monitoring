package message

type Command struct {
	GetUUID            *GetUUID
	SetUUID            *SetUUID
	GetOpcUaNode       *GetOpcUaNode
	GetSubNode         *GetSubNode
	GetNodeDescription *GetNodeDescription
}
