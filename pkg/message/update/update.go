package message

import "github.com/google/uuid"

type Update struct {
	NewConnection *NewConnection
}

type NewConnection struct {
	UUID uuid.UUID
}
