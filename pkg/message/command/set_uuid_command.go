package message

type SetUUID struct {
	UUID string
}

func NewSetUUID(uuid string) *SetUUID {
	var s SetUUID
	s.UUID = uuid
	return &s
}
