package message

import "github.com/google/uuid"

type SendUUID struct {
	UUID string
}

func NewSendUUID(uuid uuid.UUID) *SendUUID {
	var s SendUUID
	s.UUID = uuid.String()
	return &s
}
