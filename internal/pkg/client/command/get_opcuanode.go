package command

import (
	"log"

	"github.com/sudak-91/monitoring/pkg/message"
	command "github.com/sudak-91/monitoring/pkg/message/command"
)

func (c *Command) GetOpcUaNode() error {
	var m command.GetOpcUaNode
	m.Info = "get_node"
	var com command.Command
	com.GetOpcUaNode = &command.GetOpcUaNode{}
	com.GetOpcUaNode.Info = "get opcua"
	data, err := message.EncodeData(com)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(len(data))
	err = c.client.Requst(data)
	if err != nil {
		return err
	}
	return nil
}
