package command

import (
	"bytes"
	"encoding/gob"
	"log"

	message "github.com/sudak-91/monitoring/pkg/message/command"
)

func (c *Command) GetUUID() error {
	var GetUUIDCommand message.Command
	GetUUIDCommand.GetUUID = &message.GetUUID{}
	GetUUIDCommand.GetUUID.Info = "get_uuid"
	var data bytes.Buffer

	encoder := gob.NewEncoder(&data)
	err := encoder.Encode(GetUUIDCommand)
	if err != nil {
		return err
	}
	log.Printf("Send data: %v", data)
	return c.client.Requst(data.Bytes())
}
