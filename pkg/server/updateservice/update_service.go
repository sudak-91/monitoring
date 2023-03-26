package updateservice

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/sudak-91/monitoring/pkg/server"
)

type ChangeUUID struct {
	OldID uuid.UUID
	NewID uuid.UUID
}

type UpdateService struct {
	ctx        context.Context
	updateData chan any
	server     *server.Server
}

func NewUpdateService(ctx context.Context, updateData chan any, server *server.Server) *UpdateService {
	var u UpdateService
	u.ctx = ctx
	u.updateData = updateData
	u.server = server
	return &u
}

func (s *UpdateService) Update() {
	log.Println("Update service start")
	for {
		select {
		case data := <-s.updateData:
			s.router(data)
		case <-s.ctx.Done():
			return

		}
	}
}

func (s *UpdateService) router(data any) {
	switch v := data.(type) {
	case ChangeUUID:
		err := s.changeUUID(v.OldID, v.NewID)
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *UpdateService) changeUUID(oldUUID uuid.UUID, newUUID uuid.UUID) error {
	s.server.Users[newUUID] = s.server.Users[oldUUID]
	delete(s.server.Users, oldUUID)
	return nil
}
