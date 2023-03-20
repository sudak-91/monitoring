package clientservice

import (
	"bytes"
	"context"
	"encoding/gob"

	message "github.com/sudak-91/monitoring/pkg/message/update"
	"nhooyr.io/websocket"
)

func (m *ClientService) GetUUIDService() error {
	var update message.Update
	update.SendUUID = message.NewSendUUID(m.UUID)
	var data bytes.Buffer
	encoder := gob.NewEncoder(&data)
	err := encoder.Encode(update)
	if err != nil {
		return err
	}
	err = m.Conn.Write(context.TODO(), websocket.MessageBinary, data.Bytes())
	if err != nil {
		return err
	}
	return nil

}
