package clientservice

import "github.com/google/uuid"

func (m *ClientService) SetUUIDHandle(UUID string) error {
	newUUID, err := uuid.Parse(UUID)
	if err != nil {
		return err
	}
	if err = m.clientList.ChangeUUID(newUUID, m.client.UUID); err != nil {
		return err
	}
	m.client.UUID = newUUID
	return nil
}
