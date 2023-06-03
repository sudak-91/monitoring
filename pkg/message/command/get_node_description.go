package message

type GetNodeDescription struct {
	NS  uint16
	SID string
}

func GetNodeDescriptionCommand(ns uint16, sid string) Command {
	var command Command
	command.GetNodeDescription = &GetNodeDescription{NS: ns, SID: sid}
	return command
}
