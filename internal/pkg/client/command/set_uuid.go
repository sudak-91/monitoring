package command

import (
	"github.com/sudak-91/monitoring/pkg/message"
	commandMessage "github.com/sudak-91/monitoring/pkg/message/command"
)

func (c *Command) SetUUID(uuid string) error {
	var Command commandMessage.Command
	Command.SetUUID = commandMessage.NewSetUUID(uuid)
	data, err := message.EncodeData(Command)
	if err != nil {
		return err
	}
	err = c.client.Requst(data)
	if err != nil {
		return err
	}
	return nil
}
