package webserver

import (
	"context"
	"log"

	"nhooyr.io/websocket"
)

type WebUser struct {
	Conn   *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

func (w *WebUser) Run(uuid string) {
mailoop:
	for {
		select {
		default:
			MessageType, data, err := w.Conn.Read(context.TODO())
			if err != nil {
				log.Println(err)
				w.cancel()
			}
			log.Printf("%s get message %v", uuid, MessageType)
			log.Println(data)
		case <-w.ctx.Done():
			break mailoop
		}
	}
	log.Println("End")
}
