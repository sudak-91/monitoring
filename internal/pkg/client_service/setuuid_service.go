package clientservice

import (
	"github.com/google/uuid"
	serverupdates "github.com/sudak-91/monitoring/pkg/server_updates"
)

func (m *ClientService) SetUUIDService(newUUID string) error {
	id, err := uuid.Parse(newUUID)
	if err != nil {
		return nil
	}
	var ChangeUUID serverupdates.ChangeUUID
	ChangeUUID.NewID = id
	ChangeUUID.OldID = m.UUID
	m.serverUpdate <- ChangeUUID
	m.UUID = id
	m.IsUUIDTemp = false

	return nil
}
