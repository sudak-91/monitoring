package serverupdates

import "github.com/google/uuid"

type ChangeUUID struct {
	NewID uuid.UUID
	OldID uuid.UUID
}
