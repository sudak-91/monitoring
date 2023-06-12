package client

import "github.com/google/uuid"

type ChangeUUID struct {
	OldUUID uuid.UUID
	NewUUID uuid.UUID
}
