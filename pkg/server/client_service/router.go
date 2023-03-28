package clientservice

import (
	"log"

	"github.com/sudak-91/monitoring/pkg/message"
	command "github.com/sudak-91/monitoring/pkg/message/command"
	update "github.com/sudak-91/monitoring/pkg/message/update"
	"nhooyr.io/websocket"
)

func (m *ClientService) messageRouter(data command.Command) {
	switch {
	case data.GetUUID != nil:
		err := m.GetUUIDHandle()
		if err != nil {
			log.Println(err)
		}
		log.Println("GetUUID")
	case data.SetUUID != nil:
		log.Println("setUUID")
		log.Println(m.client.UUID)
		err := m.SetUUIDHandle(data.SetUUID.UUID)
		if err != nil {
			log.Println(err)
		}
		log.Println(m.client.UUID)
		log.Println("finish")
	case data.GetOpcUaNode != nil:
		log.Println("getOPCNode")
		log.Println(m.client.UUID)
		err := m.getOpcUaNodeHandle()
		if err != nil {
			log.Println(err)
		}
		log.Println("Command send")
	}
}

func (cs *ClientService) updateRouter(data any) {
	switch v := data.(type) {
	case update.SendOpcNodes:
		var upd update.Update
		upd.SendOpcNodes = &v
		data, err := message.EncodeData(upd)
		if err != nil {
			log.Println(err)
		}
		err = cs.client.Conn.Write(cs.ctx, websocket.MessageBinary, data)
		if err != nil {
			log.Println(err)
		}
		log.Println("sendOPC complete")
	}
}
