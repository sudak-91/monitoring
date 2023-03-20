package clientservice

import (
	"log"

	command "github.com/sudak-91/monitoring/pkg/message/command"
)

func (m *ClientService) Update(data command.Command) {
	switch {
	case data.GetUUID != nil:
		err := m.GetUUIDService()
		if err != nil {
			log.Println(err)
		}
		log.Println("GetUUID")
	case data.SetUUID != nil:
		log.Println("setUUID")
		log.Println(m.UUID)
		err := m.SetUUIDService(data.SetUUID.UUID)
		if err != nil {
			log.Println(err)
		}
		log.Println(m.UUID)
		log.Println("finish")

	}
}
