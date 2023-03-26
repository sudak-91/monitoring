package clientservice

import (
	"github.com/google/uuid"
	serverupdates "github.com/sudak-91/monitoring/pkg/server/updateservice"
)

func (m *ClientService) SetUUIDService(newUUID string) error {
	id, err := uuid.Parse(newUUID)
	if err != nil {
		return nil
	}
	var ChangeUUID serverupdates.ChangeUUID
	ChangeUUID.NewID = id
	ChangeUUID.OldID = m.client.UUID
	m.serverUpdate <- ChangeUUID
	m.client.UUID = id
	m.client.IsUUIDTemp = false

	return nil
}
