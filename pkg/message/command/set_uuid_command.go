package message

type SetUUID struct {
	UUID string
}

func SetUUIDCommand(uuid string) Command {
	var cmd Command
	cmd.SetUUID = &SetUUID{UUID: uuid}
	return cmd
}
