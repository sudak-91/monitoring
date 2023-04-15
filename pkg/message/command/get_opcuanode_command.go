package message

type GetOpcUaNode struct {
	Info string
}

func GetOpcUaNodeCommand() Command {
	var Command Command
	Command.GetOpcUaNode = &GetOpcUaNode{Info: "GetOPCUaNode"}
	return Command
}
