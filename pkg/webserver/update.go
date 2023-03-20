package webserver

import (
	"log"

	serverupdates "github.com/sudak-91/monitoring/pkg/server_updates"
)

func (s *Server) Update() {
	for {
		select {
		case data := <-s.updateData:
			s.router(data)
		case <-s.ctx.Done():
			return

		}
	}
}

func (s *Server) router(data any) {
	switch v := data.(type) {
	case serverupdates.ChangeUUID:
		err := s.changeUUID(v.OldID, v.NewID)
		if err != nil {
			log.Println(err)
		}
	}
}
