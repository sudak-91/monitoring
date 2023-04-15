package message

// Get users info
type GetUUID struct {
	Info string
}

func GetUUIDCommand() Command {
	var cmd Command
	cmd.GetUUID = &GetUUID{Info: "GetUUID"}
	return cmd
}
